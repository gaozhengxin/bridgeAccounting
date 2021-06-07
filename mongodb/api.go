package mongodb

import (
	"github.com/gaozhengxin/bridgeaudit/params"
)

var (
	maxListLength = 100
)

func NewQueryAPI() QueryAPI {
	return new(QueryAPIImpl)
}

func NewSyncAPI() SyncAPI {
	return new(SyncAPIImpl)
}

func NewAuditAPI() AccountingAPI {
	return new(AccountingAPIImpl)
}

type TxType int

const (
	TypeDeposit TxType = 0
	TypeMint TxType = 1
	TypeBurn TxType = 2
	TypeRedeemed TxType = 3
)

type SyncAPIImpl SyncAPI {
	*BaseQueryAPIImpl
}

type QueryAPIImpl MongoAPI {
	*BaseQueryAPIImpl
	*AccountingQueryAPI
}

type AuditAPIImpl AccountingAPI {
	*AuditQueryAPI
}

type BaseQueryAPIImpl struct {}

type AuditQueryAPI struct {}

func (*BaseQueryAPIImpl) GetSyncInfo() (*SyncInfo, error) {
	return nil. nil
}

func getSwapEvent(txtype TxType, tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error) {
	return nil, nil
}

func getSwapEventByBlockRange(txtype TxType, tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error) {
	return nil. nil
}

func getSwapEventByTimeRange(txtype TxType, tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error) {
	return nil, nil
}

func getSwapEventByUserTimeRange(txtype TxType, tokenCfg *param.TokenConfig, user string, start, end int64) ([]*SwapEvent, error) {
	return nil, nil
}

func (*BaseQueryAPIImpl) GetDeposit(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeDeposit, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetDepositsByBlockRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByBlockRange(TypeDeposit, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetDepositsByTimeRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByTimeRange(TypeDeposit, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetDepositByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByUserTimeRange(TypeDeposit, tokenCfg, user, start, end)
}

func (*BaseQueryAPIImpl) GetMint(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeMint, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetMintByBlockRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByBlockRange(TypeMint, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetMintByTimeRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByTimeRange(TypeMint, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetMintByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByUserTimeRange(TypeMint, tokenCfg, user, start, end)
}

func (*BaseQueryAPIImpl) GetBurn(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeBurn, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetBurnByBlockRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByBlockRange(TypeBurn, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetBurnByTimeRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByTimeRange(TypeBurn, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetBurnByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByUserTimeRange(TypeBurn, tokenCfg, user, start, end)
}

func (*BaseQueryAPIImpl) GetRedeemed(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeRedeemed, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetRedeemedByBlockRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByBlockRange(TypeRedeemed, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetRedeemedByTimeRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByTimeRange(TypeRedeemed, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetRedeemedByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) ([]*SwapEvent, error) {
	return getSwapEventByUserTimeRange(TypeRedeemed, tokenCfg, user, start, end)
}

func (*SyncAPIImpl) SetStartHeight(srcStartHeight, dstStartHeight int64) error {
	return nil
}

func (*SyncAPIImpl) UpdateSyncedHeight(srcSyncedHeight, dstSyncedHeight int64) error {
	return nil
}

func addSwapEvent(txtype TxType, tokenCfg *param.TokenConfig, data SwapEvent) error {
	return nil
}

func (*SyncAPIImpl) AddDeposit(tokenCfg *param.TokenConfig, data SwapEvent) error {
	return addSwapEvent(TypeDeposit, tokenCfg, data)
}

func (*SyncAPIImpl) AddMint(tokenCfg *param.TokenConfig, data SwapEvent) error {
	return addSwapEvent(TypeMint, tokenCfg, data)
}

func (*SyncAPIImpl) AddBurn(tokenCfg *param.TokenConfig, data SwapEvent) error {
	return addSwapEvent(TypeBurn, tokenCfg, data)
}

func (*SyncAPIImpl) AddRedeemed(tokenCfg *param.TokenConfig, data SwapEvent) error {
	return addSwapEvent(TypeRedeemed, tokenCfg, data)
}

func (*AuditQueryAPIImpl) GetSheetCollectionInfo() (*SheetCollectionInfo, error) {
	return nil, nil
}

func (*AuditQueryAPIImpl) GetSheetInfo(sequence int64) (*SheetInfo, error) {
	return nil, nil
}

func (*AuditQueryAPIImpl) GetSheet(tokenCfg *param.TokenConfig, sequence int64) (*Sheet, error) {
	return nil, nil
}

func (*AuditQueryAPIImpl) GetSheetsBySequenceRange(tokenCfg *param.TokenConfig, start, end int64) ([]*Sheet, error) {
	return nil, nil
}

func (*AuditAPIImpl) AddCheckRange(tokenCfg *param.TokenConfig) error {
	return nil
}

func (*AuditAPIImpl) UpdateCheckRangeAccDeposit(tokenCfg *param.TokenConfig, value float64) error {
	return nil
}

func (*AuditAPIImpl) UpdateCheckRangeAccMint(tokenCfg *param.TokenConfig, value float64) error {
	return nil
}

func (*AuditAPIImpl) UpdateCheckRangeAccBurn(tokenCfg *param.TokenConfig, value float64) error {
	return nil
}

func (*AuditAPIImpl) UpdateCheckRangeAccRedeemed(tokenCfg *param.TokenConfig, value float64) error {
	return nil
}

func (*AuditAPIImpl) AddCheckRangeInfo() error {
	return nil
}

func (*AuditAPIImpl) UpdateCheckInfo(int64) error {
	return nil
}
