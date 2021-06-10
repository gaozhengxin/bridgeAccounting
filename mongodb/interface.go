package mongodb

import (
	"github.com/gaozhengxin/bridgeaudit/params"
)

type QueryAPI interface {
	BaseQueryAPI
	AccountingQueryAPI
}

type SyncAPI interface {
	BaseQueryAPI
	SetStartHeight(srcStartHeight, dstStartHeight int64) error
	UpdateSyncedHeight(srcSyncedHeight, dstSyncedHeight int64) error
	AddDeposit(tokenCfg *param.TokenConfig, data SwapEvent) error
	AddMint(tokenCfg *param.TokenConfig, data SwapEvent) error
	AddBurn(tokenCfg *param.TokenConfig, data SwapEvent) error
	AddRedeemed(tokenCfg *param.TokenConfig, data SwapEvent) error
}

type BaseQueryAPI interface {
	GetSyncInfo() (*SyncInfo, error)
	GetDeposit(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error)
	GetDepositsByBlockRange(tokenCfg *param.TokenConfig, start, end int64) (SwapEventIter, error)
	GetDepositsByTimeRange(tokenCfg *param.TokenConfig, start, end int64) (SwapEventIter, error)
	GetDepositByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) (SwapEventIter, error)

	GetMint(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error)
	GetMintByBlockRange(tokenCfg *param.TokenConfig, start, end int64) (SwapEventIter, error)
	GetMintByTimeRange(tokenCfg *param.TokenConfig, start, end int64) (SwapEventIter, error)
	GetMintByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) (SwapEventIter, error)

	GetBurn(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error)
	GetBurnByBlockRange(tokenCfg *param.TokenConfig, start, end int64) (SwapEventIter, error)
	GetBurnByTimeRange(tokenCfg *param.TokenConfig, start, end int64) (SwapEventIter, error)
	GetBurnByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) (SwapEventIter, error)

	GetRedeemed(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error)
	GetRedeemedByBlockRange(tokenCfg *param.TokenConfig, start, end int64) (SwapEventIter, error)
	GetRedeemedByTimeRange(tokenCfg *param.TokenConfig, start, end int64) (SwapEventIter, error)
	GetRedeemedByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) (SwapEventIter, error)
}

type AccountingAPI interface {
	BaseQueryAPI
	AccountingQueryAPI
	AddSummary(tokenCfg *param.TokenConfig, summary *Summary) error
	UpdateSummary(tokenCfg *param.TokenConfig, accDeposit, accMint, accBurn, accRedeemed float64) error
	AddSummaryInfo(*SummaryInfo) error
	UpdateSummaryCollectionInfo(int64) error
}

type AccountingQueryAPI interface {
	GetSummaryCollectionInfo() (*SummaryCollectionInfo, error)
	GetSummaryInfo(sequence int64) (*SummaryInfo, error)
	GetSummary(tokenCfg *param.TokenConfig, sequence int64) (*Summary, error)
	GetSummarysBySequenceRange(tokenCfg *param.TokenConfig, start, end int64) (SummaryIter, error)
}

type SwapEventIter interface {
	Next(*SwapEvent) bool
}

type SummaryIter interface {
	Next(*Summary) bool
}
