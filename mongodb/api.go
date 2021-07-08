package mongodb

import (
	"fmt"
	"strings"

	"github.com/gaozhengxin/bridgeAccounting/params"
	"github.com/pkg/errors"

	"github.com/davecgh/go-spew/spew"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SwapEventIterImpl struct {
	*mgo.Iter
}

func (iter SwapEventIterImpl) Next(dst *SwapEvent) bool {
	return iter.Iter.Next(dst)
}

type SummaryIterImpl struct {
	*mgo.Iter
}

func (iter *SummaryIterImpl) Next(dst *Summary) bool {
	return iter.Iter.Next(dst)
}

func wrapError(err error, tag ...string) error {
	return errors.Wrap(err, fmt.Sprintf("[mongo db] %s", tag))
}

func NewQueryAPI() QueryAPI {
	return new(QueryAPIImpl)
}

func NewSyncAPI() SyncAPI {
	return new(SyncAPIImpl)
}

func NewAccountingAPI() AccountingAPI {
	return new(AccountingAPIImpl)
}

type TxType int8

const (
	TypeDeposit  TxType = iota
	TypeMint
	TypeBurn
	TypeRedeemed
)

func selectCollection(txtype TxType, tokenCfg *params.TokenConfig) (*mgo.Collection, error) {
	var coll *mgo.Collection
	switch txtype {
	case TypeDeposit:
		coll = collDeposit(tokenCfg)
	case TypeMint:
		coll = collMint(tokenCfg)
	case TypeBurn:
		coll = collBurn(tokenCfg)
	case TypeRedeemed:
		coll = collRedeemed(tokenCfg)
	default:
		return nil, fmt.Errorf("invalid txtype: %v", txtype)
	}
	if coll == nil {
		return nil, fmt.Errorf("collection not initiated, txtype: %v, pairID: %v", txtype, tokenCfg.PairID)
	}
	return coll, nil
}

type SyncAPIImpl struct {
	*BaseQueryAPIImpl
}

type QueryAPIImpl struct {
	*BaseQueryAPIImpl
	*AccountingQueryAPIImpl
}

type AccountingAPIImpl struct {
	*BaseQueryAPIImpl
	*AccountingQueryAPIImpl
}

type BaseQueryAPIImpl struct{}

type AccountingQueryAPIImpl struct{}

func (*BaseQueryAPIImpl) GetSyncInfo() (*SyncInfo, error) {
	result := new(SyncInfo)
	err := collSyncInfo.FindId(bson.M{"_id": syncInfoID}).One(result)
	if err != nil {
		return nil, wrapError(err, "GetSyncInfo")
	}
	return result, nil
}

func getSwapEvent(txtype TxType, tokenCfg *params.TokenConfig, txhash string) (*SwapEvent, error) {
	coll, err := selectCollection(txtype, tokenCfg)
	if err != nil {
		return nil, wrapError(err, "getSwapEvent", "selectCollection")
	}

	txhash = strings.ToLower(txhash)
	result := new(SwapEvent)
	err = coll.FindId(bson.M{"_id": txhash}).One(result)
	if err != nil {
		return nil, wrapError(err, "getSwapEvent")
	}
	return result, nil
}

func getSwapEventByBlockRange(
	txtype TxType,
	tokenCfg *params.TokenConfig,
	start, end int64) (SwapEventIter, error) {
	coll, err := selectCollection(txtype, tokenCfg)
	if err != nil {
		return nil, wrapError(err, "getSwapEventByBlockRange", "selectCollection")
	}

	query := bson.M{"block_number": bson.M{"$gte": start, "$lt": end}}
	iter := coll.Find(query).Iter()
	swapEventIter := &SwapEventIterImpl{
		Iter: iter,
	}
	return swapEventIter, nil
}

func getSwapEventByTimeRange(
	txtype TxType,
	tokenCfg *params.TokenConfig,
	start, end int64) (SwapEventIter, error) {
	coll, err := selectCollection(txtype, tokenCfg)
	if err != nil {
		return nil, wrapError(err, "getSwapEventByTimeRange", "selectCollection")
	}

	query := bson.M{"block_time": bson.M{"$gte": start, "$lt": end}}
	iter := coll.Find(query).Iter()
	swapEventIter := &SwapEventIterImpl{
		Iter: iter,
	}
	return swapEventIter, nil
}

func getSwapEventByUserTimeRange(
	txtype TxType,
	tokenCfg *params.TokenConfig,
	user string,
	start, end int64) (SwapEventIter, error) {
	coll, err := selectCollection(txtype, tokenCfg)
	if err != nil {
		return nil, wrapError(err, "getSwapEventByUserTimeRange", "selectCollection")
	}

	user = strings.ToLower(user)

	query := bson.M{"user": user, "block_time": bson.M{"$gte": start, "$lt": end}}
	iter := coll.Find(query).Iter()
	swapEventIter := &SwapEventIterImpl{
		Iter: iter,
	}
	return swapEventIter, nil
}

func (*BaseQueryAPIImpl) GetDeposit(tokenCfg *params.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeDeposit, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetDepositsByBlockRange(
	tokenCfg *params.TokenConfig,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByBlockRange(TypeDeposit, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetDepositsByTimeRange(
	tokenCfg *params.TokenConfig,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByTimeRange(TypeDeposit, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetDepositByUserTimeRange(
	tokenCfg *params.TokenConfig,
	user string,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByUserTimeRange(TypeDeposit, tokenCfg, user, start, end)
}

func (*BaseQueryAPIImpl) GetMint(tokenCfg *params.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeMint, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetMintByBlockRange(tokenCfg *params.TokenConfig,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByBlockRange(TypeMint, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetMintByTimeRange(
	tokenCfg *params.TokenConfig,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByTimeRange(TypeMint, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetMintByUserTimeRange(
	tokenCfg *params.TokenConfig,
	user string,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByUserTimeRange(TypeMint, tokenCfg, user, start, end)
}

func (*BaseQueryAPIImpl) GetBurn(tokenCfg *params.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeBurn, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetBurnByBlockRange(tokenCfg *params.TokenConfig,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByBlockRange(TypeBurn, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetBurnByTimeRange(tokenCfg *params.TokenConfig,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByTimeRange(TypeBurn, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetBurnByUserTimeRange(tokenCfg *params.TokenConfig,
	user string,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByUserTimeRange(TypeBurn, tokenCfg, user, start, end)
}

func (*BaseQueryAPIImpl) GetRedeemed(tokenCfg *params.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeRedeemed, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetRedeemedByBlockRange(tokenCfg *params.TokenConfig,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByBlockRange(TypeRedeemed, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetRedeemedByTimeRange(tokenCfg *params.TokenConfig,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByTimeRange(TypeRedeemed, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetRedeemedByUserTimeRange(tokenCfg *params.TokenConfig,
	user string,
	start, end int64) (SwapEventIter, error) {
	return getSwapEventByUserTimeRange(TypeRedeemed, tokenCfg, user, start, end)
}

func (*SyncAPIImpl) SetStartHeight(srcStartHeight, dstStartHeight int64) error {
	info, err := collSyncInfo.UpsertId(
		bson.M{"_id": syncInfoID},
		bson.M{"$set": bson.M{"src_start_height": srcStartHeight, "dst_start_height": dstStartHeight}})
	if err != nil {
		return wrapError(err, "SetStartHeight", spew.Sprintf("%v", info))
	}
	return nil
}

func (*SyncAPIImpl) UpdateSyncedHeight(srcSyncedHeight, dstSyncedHeight int64) error {
	info, err := collSyncInfo.UpsertId(
		bson.M{"_id": syncInfoID},
		bson.M{"$set": bson.M{"src_synced_height": srcSyncedHeight, "dst_synced_height": dstSyncedHeight}})
	if err != nil {
		return wrapError(err, "UpdateSyncedHeight", spew.Sprintf("%v", info))
	}
	return nil
}

func addSwapEvent(txtype TxType, tokenCfg *params.TokenConfig, data *SwapEvent) error {
	coll, err := selectCollection(txtype, tokenCfg)
	if err != nil {
		return wrapError(err, "addSwapEvent", "selectCollection")
	}

	err = coll.Insert(data)
	if err != nil {
		return wrapError(err, "addSwapEvent")
	}
	return nil
}

func (*SyncAPIImpl) AddDeposit(tokenCfg *params.TokenConfig, data *SwapEvent) error {
	return addSwapEvent(TypeDeposit, tokenCfg, data)
}

func (*SyncAPIImpl) AddMint(tokenCfg *params.TokenConfig, data *SwapEvent) error {
	return addSwapEvent(TypeMint, tokenCfg, data)
}

func (*SyncAPIImpl) AddBurn(tokenCfg *params.TokenConfig, data *SwapEvent) error {
	return addSwapEvent(TypeBurn, tokenCfg, data)
}

func (*SyncAPIImpl) AddRedeemed(tokenCfg *params.TokenConfig, data *SwapEvent) error {
	return addSwapEvent(TypeRedeemed, tokenCfg, data)
}

func (*AccountingQueryAPIImpl) GetSummaryCollectionInfo() (*SummaryCollectionInfo, error) {
	result := new(SummaryCollectionInfo)
	err := collSummaryCollectionInfo.FindId(bson.M{"_id": summaryCollectionInfoID}).One(result)
	if err != nil {
		return nil, wrapError(err, "GetSummaryCollectionInfo")
	}
	return result, nil
}

func (*AccountingQueryAPIImpl) GetSummaryInfo(sequence int64) (*SummaryInfo, error) {
	result := new(SummaryInfo)
	err := collSummaryInfo.FindId(bson.M{"_id": sequence}).One(result)
	if err != nil {
		return nil, wrapError(err, "GetSummaryInfo")
	}
	return result, nil
}

func (*AccountingQueryAPIImpl) GetSummary(tokenCfg *params.TokenConfig, sequence int64) (*Summary, error) {
	coll := collSummarys[tokenCfg.PairID]
	if coll == nil {
		return nil, wrapError(fmt.Errorf("collection not initiated, pairID: %v", tokenCfg.PairID), "GetSummary")
	}
	result := new(Summary)
	err := coll.FindId(bson.M{"_id": sequence}).One(result)
	if err != nil {
		return nil, wrapError(err, "GetSummaryInfo")
	}
	return result, nil
}

func (*AccountingQueryAPIImpl) GetSummarysBySequenceRange(tokenCfg *params.TokenConfig, start, end int64) (SummaryIter, error) {
	coll := collSummarys[tokenCfg.PairID]
	if coll == nil {
		return nil, wrapError(fmt.Errorf("collection not initiated, pairID: %v", tokenCfg.PairID), "GetSummarysBySequenceRange")
	}
	iter := coll.Find(bson.M{"_id": bson.M{"$gte": start, "$lt": end}}).Iter()
	summaryIter := &SummaryIterImpl{
		Iter: iter,
	}
	return summaryIter, nil
}

func (*AccountingAPIImpl) AddSummary(tokenCfg *params.TokenConfig, summary *Summary) error {
	coll := collSummarys[tokenCfg.PairID]
	if coll == nil {
		return wrapError(fmt.Errorf("collection not initiated, pairID: %v", tokenCfg.PairID), "AddSummary")
	}
	return nil
}

func (*AccountingAPIImpl) UpdateSummary(
	tokenCfg *params.TokenConfig,
	accDeposit,
	accMint,
	accBurn,
	accRedeemed float64) error {
	coll := collSummarys[tokenCfg.PairID]
	if coll == nil {
		return wrapError(fmt.Errorf("collection not initiated, pairID: %v", tokenCfg.PairID), "UpdateSummary")
	}
	return nil
}

func (*AccountingAPIImpl) AddSummaryInfo(data *SummaryInfo) error {
	err := collSummaryInfo.Insert(data)
	if err != nil {
		return wrapError(err, "AddSummaryInfo")
	}
	return nil
}

func (*AccountingAPIImpl) UpdateSummaryCollectionInfo(latestSequence int64) error {
	info, err := collSummaryCollectionInfo.UpsertId(
		bson.M{"_id": summaryCollectionInfoID},
		bson.M{"$set": latestSequence},
	)
	if err != nil {
		return wrapError(err, "UpdateSummaryCollectionInfo", spew.Sprintf("%v", info))
	}
	return nil
}
