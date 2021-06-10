package scanner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/anyswap/CrossChain-Bridge/cmd/utils"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/fsn-dev/fsn-go-sdk/efsn/common"
	"github.com/fsn-dev/fsn-go-sdk/efsn/core/types"
	"github.com/fsn-dev/fsn-go-sdk/efsn/ethclient"
	ethereum "github.com/fsn-dev/fsn-go-sdk/efsn"
	"github.com/gaozhengxin/bridgeAccounting/params"
	"github.com/gaozhengxin/bridgeAccounting/mongodb"
	"github.com/gaozhengxin/bridgeAccounting/tools"
	"github.com/gaozhengxin/bridgeAccounting/accounting"
	"github.com/urfave/cli/v2"
)

var (
	scanReceiptFlag = &cli.BoolFlag{
		Name:  "scanReceipt",
		Usage: "scan transaction receipt instead of transaction",
	}

	startHeightFlag = &cli.Int64Flag{
		Name:  "start",
		Usage: "start height (start inclusive)",
		Value: -200,
	}

	timeoutFlag = &cli.Uint64Flag{
		Name:  "timeout",
		Usage: "timeout of scanning one block in seconds",
		Value: 300,
	}

	// StartCommand scan swaps on eth like blockchain, and do accounting
	StartCommand = &cli.Command{
		Action:    start,
		Name:      "start",
		Usage:     "scan cross chain swaps",
		ArgsUsage: " ",
		Description: `
scan cross chain swaps
`,
		Flags: []cli.Flag{
			utils.ConfigFileFlag,
		},
	}

	// 0. Deposit and 3. Redeemed
	transferFuncHash       = common.FromHex("0xa9059cbb")
	transferFromFuncHash   = common.FromHex("0x23b872dd")

	// 1. Mint
	swapinFuncHash = common.FromHex("0xec126c77")

	// 2. Burn
	addressSwapoutFuncHash = common.FromHex("0x628d6cba") // for ETH like `address` type address
	stringSwapoutFuncHash  = common.FromHex("0xad54056d") // for BTC like `string` type address

	// 0. Deposit and 3. Redeemed log, but also seen in 1. Mint
	transferLogTopic       = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

	// 1. Mint log
	swapInLogTopic = common.HexToHash("0x05d0634fe981be85c22e2942a880821b70095d84e152c3ea3c17a4e4250d9d61")

	// 2. Burn log
	addressSwapoutLogTopic = common.HexToHash("0x6b616089d04950dc06c45c6dd787d657980543f89651aec47924752c7d16c888")
	stringSwapoutLogTopic  = common.HexToHash("0x9c92ad817e5474d30a4378deface765150479363a897b0590fbb12ae9d89396b")
)

const (
	swapExistKeywords   = "mgoError: Item is duplicate"
	httpTimeoutKeywords = "Client.Timeout exceeded while awaiting headers"
)

type ethSwapScanner struct {
	gateway     string
	scanReceipt bool

	startHeightArgument int64

	endHeight    uint64
	stableHeight uint64
	jobCount     uint64

	processBlockTimeout time.Duration
	processBlockTimers  []*time.Timer

	client *ethclient.Client
	ctx    context.Context

	rpcInterval   time.Duration
	rpcRetryCount int
}

var (
	dbAPI mongodb.SyncAPI
/*
type SyncAPI interface {
	BaseQueryAPI
	SetStartHeight(srcStartHeight, dstStartHeight int64) error
	UpdateSyncedHeight(srcSyncedHeight, dstSyncedHeight int64) error
	AddDeposit(tokenCfg *param.TokenConfig, data SwapEvent) error
	AddMint(tokenCfg *param.TokenConfig, data SwapEvent) error
	AddBurn(tokenCfg *param.TokenConfig, data SwapEvent) error
	AddRedeemed(tokenCfg *param.TokenConfig, data SwapEvent) error
}
*/
)

func start(ctx *cli.Context) error {
	utils.SetLogger(ctx)
	cfg := params.LoadConfig(utils.GetConfigFilePath(ctx))
	go params.WatchAndReloadScanConfig()

	srcScanner := &ethSwapScanner{
		ctx:           context.Background(),
		rpcInterval:   1 * time.Second,
		rpcRetryCount: 3,
	}
	srcScanner.gateway = ctg.Gateway
	srcScanner.scanReceipt = cfg.ScanReceipt
	srcScanner.StartHeightArgument = cfg.StartHeightArgument
	// TODO
	srcScanner.endHeight = ctx.Uint64(utils.EndHeightFlag.Name)
	srcScanner.stableHeight = ctx.Uint64(utils.StableHeightFlag.Name)
	srcScanner.jobCount = ctx.Uint64(utils.JobsFlag.Name)
	srcScanner.processBlockTimeout = time.Duration(ctx.Uint64(timeoutFlag.Name)) * time.Second

	log.Info("get argument success",
		"gateway", scanner.gateway,
		"scanReceipt", scanner.scanReceipt,
		"start", startHeightArgument,
		"end", scanner.endHeight,
		"stable", scanner.stableHeight,
		"jobs", scanner.jobCount,
		"timeout", scanner.processBlockTimeout,
	)

	scanner.initClient()
	select {
		go srcScanner.run()
		go dstScanner.run()
		go accounting()
	}
	return nil
}

func (scanner *ethSwapScanner) initClient() {
	ethcli, err := ethclient.Dial(scanner.gateway)
	if err != nil {
		log.Fatal("ethclient.Dail failed", "gateway", scanner.gateway, "err", err)
	}
	log.Info("ethclient.Dail gateway success", "gateway", scanner.gateway)
	scanner.client = ethcli
}

func (scanner *ethSwapScanner) run() {
	scanner.processBlockTimers = make([]*time.Timer, scanner.jobCount+1)
	for i := 0; i < len(scanner.processBlockTimers); i++ {
		scanner.processBlockTimers[i] = time.NewTimer(scanner.processBlockTimeout)
	}

	wend := scanner.endHeight
	if wend == 0 {
		wend = scanner.loopGetLatestBlockNumber()
	}
	if scanner.startHeightArgument != 0 {
		var start uint64
		if scanner.startHeightArgument > 0 {
			start = uint64(scanner.startHeightArgument)
		} else if scanner.startHeightArgument < 0 {
			start = wend - uint64(-scanner.startHeightArgument)
		}
		scanner.doScanRangeJob(start, wend)
	}
	if scanner.endHeight == 0 {
		scanner.scanLoop(wend)
	}
}

func (scanner *ethSwapScanner) doScanRangeJob(start, end uint64) {
	log.Info("start scan range job", "start", start, "end", end, "jobs", scanner.jobCount)
	if scanner.jobCount == 0 {
		log.Fatal("zero count jobs specified")
	}
	if start >= end {
		log.Fatalf("wrong scan range [%v, %v)", start, end)
	}
	jobs := scanner.jobCount
	count := end - start
	step := count / jobs
	if step == 0 {
		jobs = 1
		step = count
	}
	wg := new(sync.WaitGroup)
	for i := uint64(0); i < jobs; i++ {
		from := start + i*step
		to := start + (i+1)*step
		if i+1 == jobs {
			to = end
		}
		wg.Add(1)
		go scanner.scanRange(i+1, from, to, wg)
	}
	if scanner.endHeight != 0 {
		wg.Wait()
	}
}

func (scanner *ethSwapScanner) scanRange(job, from, to uint64, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info(fmt.Sprintf("[%v] scan range", job), "from", from, "to", to)

	for h := from; h < to; h++ {
		scanner.scanBlock(job, h, false)
	}

	log.Info(fmt.Sprintf("[%v] scan range finish", job), "from", from, "to", to)
}

func (scanner *ethSwapScanner) scanLoop(from uint64) {
	stable := scanner.stableHeight
	log.Info("start scan loop job", "from", from, "stable", stable)
	for {
		latest := scanner.loopGetLatestBlockNumber()
		for h := from; h <= latest; h++ {
			scanner.scanBlock(0, h, true)
		}
		if from+stable < latest {
			from = latest - stable
		}
		time.Sleep(1 * time.Second)
	}
}

func (scanner *ethSwapScanner) loopGetLatestBlockNumber() uint64 {
	for { // retry until success
		header, err := scanner.client.HeaderByNumber(scanner.ctx, nil)
		if err == nil {
			log.Info("get latest block number success", "height", header.Number)
			return header.Number.Uint64()
		}
		log.Warn("get latest block number failed", "err", err)
		time.Sleep(scanner.rpcInterval)
	}
}

func (scanner *ethSwapScanner) loopGetTxReceipt(txHash common.Hash) (receipt *types.Receipt, err error) {
	for i := 0; i < 5; i++ { // with retry
		receipt, err = scanner.client.TransactionReceipt(scanner.ctx, txHash)
		if err == nil {
			return receipt, err
		}
		time.Sleep(scanner.rpcInterval)
	}
	return nil, err
}

func (scanner *ethSwapScanner) loopGetBlock(height uint64) (block *types.Block, err error) {
	blockNumber := new(big.Int).SetUint64(height)
	for i := 0; i < 5; i++ { // with retry
		block, err = scanner.client.BlockByNumber(scanner.ctx, blockNumber)
		if err == nil {
			return block, nil
		}
		log.Warn("get block failed", "height", height, "err", err)
		time.Sleep(scanner.rpcInterval)
	}
	return nil, err
}

func (scanner *ethSwapScanner) scanBlock(job, height uint64, cache bool) {
	block, err := scanner.loopGetBlock(height)
	if err != nil {
		return
	}
	blockHash := block.Hash().Hex()
	if cache && cachedBlocks.isScanned(blockHash) {
		return
	}
	log.Info(fmt.Sprintf("[%v] scan block %v", job, height), "hash", blockHash, "txs", len(block.Transactions()))

	scanner.processBlockTimers[job].Reset(scanner.processBlockTimeout)
SCANTXS:
	for i, tx := range block.Transactions() {
		select {
		case <-scanner.processBlockTimers[job].C:
			log.Warn(fmt.Sprintf("[%v] scan block %v timeout", job, height), "hash", blockHash, "txs", len(block.Transactions()))
			break SCANTXS
		default:
			log.Debug(fmt.Sprintf("[%v] scan tx in block %v index %v", job, height, i), "tx", tx.Hash().Hex())
			scanner.scanTransaction(tx)
		}
	}
	if cache {
		cachedBlocks.addBlock(blockHash)
	}
}

func (scanner *ethSwapScanner) scanTransaction(tx *types.Transaction) {
	if tx.To() == nil {
		return
	}
	txHash := tx.Hash().Hex()
	var receipt *types.Receipt
	if scanner.scanReceipt {
		r, err := scanner.loopGetTxReceipt(tx.Hash())
		if err != nil {
			log.Warn("get tx receipt error", "txHash", txHash, "err", err)
			return
		}
		receipt = r
	}

	for _, tokenCfg := range params.GetScanConfig().Tokens {
		swapTxType, swapEvent, verifyErr := scanner.verifyTransaction(tx, receipt, tokenCfg)
		if verifyErr != nil {
			log.Debug("verify tx failed", "txHash", txHash, "err", verifyErr)
		}

		mgoSwapEvent := convertToMgoSwapEvent(swapEvent, scanner.cachedDecimal(tokenCfg))

		var syncError error
		switch swapTxType {
		case TypeDeposit:
			syncError = dbAPI.AddDeposit(tokenCfg, mgoSwapEvent)
		case TypeMint:
			syncError = dbAPI.AddMint(tokenCfg, mgoSwapEvent)
		case TypeBurn:
			syncError = dbAPI.AddBurn(tokenCfg, mgoSwapEvent)
		case TypeRedeemed:
			syncError = dbAPI.AddRedeemed(tokenCfg, mgoSwapEvent)
		default:
			if verifyErr != nil {
				scanner.printVerifyError(txHash, verifyErr)
			}
		}
		if syncError != nil {
			log.Warn("Add swap event error", "swapTxType", swapTxType, "syncError", syncError)
		}
	}
}

type SwapTxType int8

const (
	TypeDeposit  SwapTxType = iota
	TypeMint
	TypeBurn
	TypeRedeemed
)

const TypeNull SwapTxType = -1

type SwapEvent struct {
	TxHash common.Hash
	BlockTime int64
	BlockNumber *big.Int
	Amount *big.Int
	User common.Address
}

func (scanner *ethSwapScanner) verifyTransaction(tx *types.Transaction, receipt *types.Receipt, tokenCfg *params.TokenConfig) (txType SwapTxType, swapData *SwapEvent, verifyErr error) {
	txTo := tx.To().Hex()
	cmpTxTo := tokenCfg.TokenAddress
	depositAddress := tokenCfg.DepositAddress

	if tokenCfg.CallByContract != "" {
		cmpTxTo = tokenCfg.CallByContract
		if receipt == nil {
			txHash := tx.Hash()
			r, err := scanner.loopGetTxReceipt(txHash)
			if err != nil {
				log.Warn("get tx receipt error", "txHash", txHash.Hex(), "err", err)
				return TypeNull, nil, nil
			}
			receipt = r
		}
	}

	switch {
	case depositAddress != "":
		if tokenCfg.IsNativeToken() {
			// Src chain, Deposit or Redeemed
			// TODO
			matched := strings.EqualFold(txTo, depositAddress)
			if matched {
				swapData = &SwapEvent{
					TxHash: strings.ToLower(tx.Hash().String()),
					BlockTime: getBlockTimestamp(receipt.BlockNumber.ToInt()),
					BlockNumber: receipt.BlockNumber,
					Amount: tx.Value(),
					User: tx.From(),
				}
				return TypeDeposit, swapData, nil
			}
			return TypeNull, nil, nil
		} else if strings.EqualFold(txTo, cmpTxTo) {
			swapData, verifyErr = scanner.verifyErc20SwapinTx(tx, receipt, tokenCfg)
			if verifyErr == tokens.ErrTxWithWrongReceiver {
				// swapin my have multiple deposit addresses for different bridges
				return TypeNull, nil, verifyErr
			}
			return TypeDeposit, swapData, verifyErr
		}
	case redeemAddress != "":
		// Dst chain, Mint or Burn
		// TODO
		switch {
		case !scanner.scanReceipt:
			if strings.EqualFold(txTo, cmpTxTo) {
				swapData, verifyErr = scanner.verifySwapoutTx(tx, receipt, tokenCfg)
				return TypeBurn, swapData, verifyErr
			}
		default:
			swapData = &SwapEvent{
				TxHash: strings.ToLower(tx.Hash().String()),
				BlockTime: getBlockTimestamp(receipt.BlockNumber.ToInt()),
				BlockNumber: receipt.BlockNumber,
				Amount: nil,
				User: common.Address{},
			}
			verifyErr = scanner.parseSwapoutTxLogs(receipt.Logs, tokenCfg, swapData)
			if verifyErr == nil {
				return TypeBurn, swapData, nil
			}
		}
	default:
	}
	return TypeNull, nil, verifyErr
}

func (scanner *ethSwapScanner) printVerifyError(txHash string, verifyErr error) {
	switch {
	case errors.Is(verifyErr, tokens.ErrTxFuncHashMismatch):
	case errors.Is(verifyErr, tokens.ErrTxWithWrongReceiver):
	case errors.Is(verifyErr, tokens.ErrTxWithWrongContract):
	case errors.Is(verifyErr, tokens.ErrTxNotFound):
	default:
		log.Debug("verify swap error", "txHash", txHash, "err", verifyErr)
	}
}

func (scanner *ethSwapScanner) getSwapoutFuncHashByTxType(txType string) []byte {
	switch strings.ToLower(txType) {
	case params.TxSwapout:
		return addressSwapoutFuncHash
	case params.TxSwapout2:
		return stringSwapoutFuncHash
	default:
		log.Errorf("unknown swapout tx type %v", txType)
		return nil
	}
}

func (scanner *ethSwapScanner) getLogTopicByTxType(txType string) common.Hash {
	switch strings.ToLower(txType) {
	case params.TxSwapin:
		return transferLogTopic
	case params.TxSwapout:
		return addressSwapoutLogTopic
	case params.TxSwapout2:
		return stringSwapoutLogTopic
	default:
		log.Errorf("unknown tx type %v", txType)
		return common.Hash{}
	}
}

func (scanner *ethSwapScanner) verifyErc20SwapinTx(tx *types.Transaction, receipt *types.Receipt, tokenCfg *params.TokenConfig) (swapData *SwapEvent, err error) {
	swapData = &SwapEvent{
		TxHash: strings.ToLower(tx.Hash().String()),
		BlockTime: getBlockTimestamp(receipt.BlockNumber.ToInt()),
		BlockNumber: receipt.BlockNumber,
		Amount: nil,
		User: common.Address{},
	}
	if receipt == nil {
		err = scanner.parseErc20SwapinTxInput(tx.Data(), tokenCfg.DepositAddress, swapData)
	} else {
		err = scanner.parseErc20SwapinTxLogs(receipt.Logs, tokenCfg, swapData)
	}
	return swapData, err
}

func (scanner *ethSwapScanner) verifySwapoutTx(tx *types.Transaction, receipt *types.Receipt, tokenCfg *params.TokenConfig) (swapData *SwapEvent, err error) {
	swapData = &SwapEvent{
		TxHash: strings.ToLower(tx.Hash().String()),
		BlockTime: getBlockTimestamp(receipt.BlockNumber.ToInt()),
		BlockNumber: receipt.BlockNumber,
		Amount: nil,
		User: common.Address{},
	}
	if receipt == nil {
		err = scanner.parseSwapoutTxInput(tx.Data(), tokenCfg.TxType, swapData)
	} else {
		err = scanner.parseSwapoutTxLogs(receipt.Logs, tokenCfg, swapData)
	}
	return swapData, err
}

func (scanner *ethSwapScanner) parseErc20SwapinTxInput(input []byte, depositAddress string, swapData *SwapEvent) error {
	if len(input) < 4 {
		return nil, tokens.ErrTxWithWrongInput
	}
	var receiver string
	funcHash := input[:4]
	switch {
	case bytes.Equal(funcHash, transferFuncHash):
		receiver = common.BytesToAddress(common.GetData(input, 4, 32)).Hex()
	case bytes.Equal(funcHash, transferFromFuncHash):
		receiver = common.BytesToAddress(common.GetData(input, 36, 32)).Hex()
	default:
		return tokens.ErrTxFuncHashMismatch
	}
	if !strings.EqualFold(receiver, depositAddress) {
		return tokens.ErrTxWithWrongReceiver
	}
	return nil
}

func (scanner *ethSwapScanner) parseErc20SwapinTxLogs(logs []*types.Log, tokenCfg *params.TokenConfig, swapData *SwapEvent) (err error) {
	targetContract := tokenCfg.TokenAddress
	depositAddress := tokenCfg.DepositAddress
	cmpLogTopic := scanner.getLogTopicByTxType(tokenCfg.TxType)

	for _, rlog := range logs {
		if rlog.Removed {
			continue
		}
		if !strings.EqualFold(rlog.Address.Hex(), targetContract) {
			continue
		}
		if len(rlog.Topics) != 3 || rlog.Data == nil {
			continue
		}
		if rlog.Topics[0] == cmpLogTopic {
			receiver := common.BytesToAddress(rlog.Topics[2][:]).Hex()
			if strings.EqualFold(receiver, depositAddress) {
				return nil
			}
			return tokens.ErrTxWithWrongReceiver
		}
	}
	return tokens.ErrDepositLogNotFound
}

func (scanner *ethSwapScanner) parseSwapoutTxInput(input []byte, txType string, swapData *SwapEvent) error {
	if len(input) < 4 {
		return tokens.ErrTxWithWrongInput
	}
	funcHash := input[:4]
	if bytes.Equal(funcHash, scanner.getSwapoutFuncHashByTxType(txType)) {
		return nil
	}
	return tokens.ErrTxFuncHashMismatch
}

func (scanner *ethSwapScanner) parseSwapoutTxLogs(logs []*types.Log, tokenCfg *params.TokenConfig, swapData *SwapEvent) (err error) {
	targetContract := tokenCfg.TokenAddress
	cmpLogTopic := scanner.getLogTopicByTxType(tokenCfg.TxType)

	for _, rlog := range logs {
		if rlog.Removed {
			continue
		}
		if !strings.EqualFold(rlog.Address.Hex(), targetContract) {
			continue
		}
		if len(rlog.Topics) != 2 || rlog.Data == nil {
			continue
		}
		if rlog.Topics[0] == cmpLogTopic {
			return nil
		}
	}
	return tokens.ErrSwapoutLogNotFound
}

type cachedSacnnedBlocks struct {
	capacity  int
	nextIndex int
	hashes    []string
}

var cachedBlocks = &cachedSacnnedBlocks{
	capacity:  100,
	nextIndex: 0,
	hashes:    make([]string, 100),
}

func (cache *cachedSacnnedBlocks) addBlock(blockHash string) {
	cache.hashes[cache.nextIndex] = blockHash
	cache.nextIndex = (cache.nextIndex + 1) % cache.capacity
}

func (cache *cachedSacnnedBlocks) isScanned(blockHash string) bool {
	for _, b := range cache.hashes {
		if b == blockHash {
			return true
		}
	}
	return false
}
