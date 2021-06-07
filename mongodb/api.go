package mongodb

import (
	"fmt"
	"strings"

	"github.com/gaozhengxin/bridgeaudit/params"
	"github.com/pkg/errors"

	"github.com/davecgh/go-spew/spew"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SwapEventIterImpl struct {
	*mgo.Iter
}

func (iter *SwapEventIterImpl) Next(dst *SwapEvent) bool {
	return iter.Iter.Next(dst)
}

type SummaryImpl struct {
	*mgo.Iter
}

func (iter *SummaryImpl) Next(dst *SwapEvent) bool {
	return iter.Iter.Next(dst)
}

func wrapError(err error, tag ...string) {
	return errors.Wrap(err, fmt.Sprintf("[mongo db] %s", tag))
}

func NewQueryAPI() QueryAPI {
	return new(QueryAPIImpl)
}

func NewSyncAPI() SyncAPI {
	return new(SyncAPIImpl)
}

func NewAuditAPI() AccountingAPI {
	return new(AccountingAPIImpl)
}

type TxType int8

const (
	TypeDeposit  TxType = iota
	TypeMint
	TypeBurn
	TypeRedeemed
)

func selectCollection(txtype TxType, tokenCfg *param.TokenConfig) (*mgo.Collection, error) {
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
	*AccountingQueryAPI
}

type AuditAPIImpl struct {
	*AuditQueryAPI
}

type BaseQueryAPIImpl struct{}

type AuditQueryAPI struct{}

func (*BaseQueryAPIImpl) GetSyncInfo() (*SyncInfo, error) {
	result := new(SyncInfo)
	err := collSyncInfo.FindId(bson.M{"_id": syncInfoID}).One(result)
	if err != nil {
		return nil, wrapError(err, "GetSyncInfo")
	}
	return result.nil
}

func getSwapEvent(txtype TxType, tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error) {
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
	tokenCfg *param.TokenConfig,
	start, end int64) (*SwapEventIter, error) {
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
	tokenCfg *param.TokenConfig,
	start, end int64) (*SwapEventIter, error) {
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
	tokenCfg *param.TokenConfig,
	user string,
	start, end int64) (*SwapEventIter, error) {
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

func (*BaseQueryAPIImpl) GetDeposit(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeDeposit, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetDepositsByBlockRange(
	tokenCfg *param.TokenConfig,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByBlockRange(TypeDeposit, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetDepositsByTimeRange(
	tokenCfg *param.TokenConfig,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByTimeRange(TypeDeposit, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetDepositByUserTimeRange(
	tokenCfg *param.TokenConfig,
	user string,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByUserTimeRange(TypeDeposit, tokenCfg, user, start, end)
}

func (*BaseQueryAPIImpl) GetMint(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeMint, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetMintByBlockRange(tokenCfg *param.TokenConfig,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByBlockRange(TypeMint, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetMintByTimeRange(
	tokenCfg *param.TokenConfig,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByTimeRange(TypeMint, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetMintByUserTimeRange(
	tokenCfg *param.TokenConfig,
	user string,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByUserTimeRange(TypeMint, tokenCfg, user, start, end)
}

func (*BaseQueryAPIImpl) GetBurn(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeBurn, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetBurnByBlockRange(tokenCfg *param.TokenConfig,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByBlockRange(TypeBurn, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetBurnByTimeRange(tokenCfg *param.TokenConfig,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByTimeRange(TypeBurn, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetBurnByUserTimeRange(tokenCfg *param.TokenConfig,
	user string,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByUserTimeRange(TypeBurn, tokenCfg, user, start, end)
}

func (*BaseQueryAPIImpl) GetRedeemed(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error) {
	return getSwapEvent(TypeRedeemed, tokenCfg, txhash)
}

func (*BaseQueryAPIImpl) GetRedeemedByBlockRange(tokenCfg *param.TokenConfig,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByBlockRange(TypeRedeemed, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetRedeemedByTimeRange(tokenCfg *param.TokenConfig,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByTimeRange(TypeRedeemed, tokenCfg, start, end)
}

func (*BaseQueryAPIImpl) GetRedeemedByUserTimeRange(tokenCfg *param.TokenConfig,
	user string,
	start, end int64) (*SwapEventIter, error) {
	return getSwapEventByUserTimeRange(TypeRedeemed, tokenCfg, user, start, end)
}

func (*SyncAPIImpl) SetStartHeight(srcStartHeight, dstStartHeight int64) error {
	info, err := collSyncInfo.UpsertID(
		bson.M{"_id": syncInfoID},
		bson.M{"$set": bson.M{"src_start_height": srcStartHeight, "dst_start_height": dstStartHeight}})
	if err != nil {
		return wrapError(err, "SetStartHeight", spew.Sprintf("%v", info))
	}
	return nil
}

func (*SyncAPIImpl) UpdateSyncedHeight(srcSyncedHeight, dstSyncedHeight int64) error {
	info, err := collSyncInfo.UpsertID(
		bson.M{"_id": syncInfoID},
		bson.M{"$set": bson.M{"src_synced_height": srcSyncedHeight, "dst_synced_height": dstSyncedHeight}})
	if err != nil {
		return wrapError(err, "UpdateSyncedHeight", spew.Sprintf("%v", info))
	}
	return nil
}

func addSwapEvent(txtype TxType, tokenCfg *param.TokenConfig, data SwapEvent) error {
	coll, err := selectCollection(txtype, tokenCfg)
	if err != nil {
		return nil, wrapError(err, "addSwapEvent", "selectCollection")
	}

	err = coll.Insert(data)
	if err != nil {
		return wrapError(err, "addSwapEvent")
	}
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

func (*AuditQueryAPIImpl) GetSummaryCollectionInfo() (*SummaryCollectionInfo, error) {
	result := new(SummaryCollectionInfo)
	err := collSummaryCollectionInfo.FindId(bson.M{"_id": summaryCollectionInfoID}).One(result)
	if err != nil {
		return nil, wrapError(err, "GetSummaryCollectionInfo")
	}
	return result, nil
}

func (*AuditQueryAPIImpl) GetSummaryInfo(sequence int64) (*SummaryInfo, error) {
	result := new(SummaryInfo)
	err := collSummaryInfo.FindId(bson.M{"_id": sequence}).One(result)
	if err != nil {
		return nil, wrapError(err, "GetSummaryInfo")
	}
	return result, nil
}

func (*AuditQueryAPIImpl) GetSummary(tokenCfg *param.TokenConfig, sequence int64) (*Summary, error) {
	coll := collSummarys(tokenCfg)
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

func (*AuditQueryAPIImpl) GetSummarysBySequenceRange(tokenCfg *param.TokenConfig, start, end int64) (*SummaryIter, error) {
	coll := collSummarys(tokenCfg)
	if coll == nil {
		return nil, wrapError(fmt.Errorf("collection not initiated, pairID: %v", tokenCfg.PairID), "GetSummarysBySequenceRange")
	}
	iter := coll.Find(bson.M{"_id": bson.M{"$gte": start, "$lt": end}}).Iter()
	summaryIter := &SummaryIterImpl{
		Iter: iter,
	}
	return summaryIter, nil
}

func (*AuditAPIImpl) AddSummary(tokenCfg *param.TokenConfig, summary *Summary) error {
	coll := collSummarys(tokenCfg)
	if coll == nil {
		return nil, wrapError(fmt.Errorf("collection not initiated, pairID: %v", tokenCfg.PairID), GetSummary)
	}
	return nil
}

func (*AuditAPIImpl) UpdateSummary(
	tokenCfg *param.TokenConfig,
	accDeposit,
	accMint,
	accBurn,
	accRedeemed float64) error {
	coll := collSummarys(tokenCfg)
	if coll == nil {
		return nil, wrapError(fmt.Errorf("collection not initiated, pairID: %v", tokenCfg.PairID), GetSummary)
	}
	return nil
}

func (*AuditAPIImpl) AddSummaryInfo(data *SummaryInfo) error {
	err := collSummaryInfo.Insert(data)
	if err != nil {
		return wrapError(err, "AddSummaryInfo")
	}
	return nil
}

func (*AuditAPIImpl) UpdateSummaryCollectionInfo(latestSequence int64) error {
	info, err := collSummaryCollectionInfo.UpsertID(
		bson.M{"_id": summaryCollectionInfoID},
		bson.M{"$set": latestSequence},
	)
	if err != nil {
		return wrapError(err, "UpdateSummaryCollectionInfo", spew.Sprintf("%v", info))
	}
	return nil
}
