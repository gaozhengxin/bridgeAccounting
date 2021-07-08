package mongodb

import (
	"github.com/gaozhengxin/bridgeAccounting/params"
)

type QueryAPI interface {
	BaseQueryAPI
	AccountingQueryAPI
}

type SyncAPI interface {
	BaseQueryAPI
	SetStartHeight(srcStartHeight, dstStartHeight int64) error
	UpdateSyncedHeight(srcSyncedHeight, dstSyncedHeight int64) error
	AddDeposit(tokenCfg *params.TokenConfig, data *SwapEvent) error
	AddMint(tokenCfg *params.TokenConfig, data *SwapEvent) error
	AddBurn(tokenCfg *params.TokenConfig, data *SwapEvent) error
	AddRedeemed(tokenCfg *params.TokenConfig, data *SwapEvent) error
}

type BaseQueryAPI interface {
	GetSyncInfo() (*SyncInfo, error)
	GetDeposit(tokenCfg *params.TokenConfig, txhash string) (*SwapEvent, error)
	GetDepositsByBlockRange(tokenCfg *params.TokenConfig, start, end int64) (SwapEventIter, error)
	GetDepositsByTimeRange(tokenCfg *params.TokenConfig, start, end int64) (SwapEventIter, error)
	GetDepositByUserTimeRange(tokenCfg *params.TokenConfig, user string, start, end int64) (SwapEventIter, error)

	GetMint(tokenCfg *params.TokenConfig, txhash string) (*SwapEvent, error)
	GetMintByBlockRange(tokenCfg *params.TokenConfig, start, end int64) (SwapEventIter, error)
	GetMintByTimeRange(tokenCfg *params.TokenConfig, start, end int64) (SwapEventIter, error)
	GetMintByUserTimeRange(tokenCfg *params.TokenConfig, user string, start, end int64) (SwapEventIter, error)

	GetBurn(tokenCfg *params.TokenConfig, txhash string) (*SwapEvent, error)
	GetBurnByBlockRange(tokenCfg *params.TokenConfig, start, end int64) (SwapEventIter, error)
	GetBurnByTimeRange(tokenCfg *params.TokenConfig, start, end int64) (SwapEventIter, error)
	GetBurnByUserTimeRange(tokenCfg *params.TokenConfig, user string, start, end int64) (SwapEventIter, error)

	GetRedeemed(tokenCfg *params.TokenConfig, txhash string) (*SwapEvent, error)
	GetRedeemedByBlockRange(tokenCfg *params.TokenConfig, start, end int64) (SwapEventIter, error)
	GetRedeemedByTimeRange(tokenCfg *params.TokenConfig, start, end int64) (SwapEventIter, error)
	GetRedeemedByUserTimeRange(tokenCfg *params.TokenConfig, user string, start, end int64) (SwapEventIter, error)
}

type AccountingAPI interface {
	BaseQueryAPI
	AccountingQueryAPI
	AddSummary(tokenCfg *params.TokenConfig, summary *Summary) error
	UpdateSummary(tokenCfg *params.TokenConfig, accDeposit, accMint, accBurn, accRedeemed float64) error
	AddSummaryInfo(*SummaryInfo) error
	UpdateSummaryCollectionInfo(int64) error
}

type AccountingQueryAPI interface {
	GetSummaryCollectionInfo() (*SummaryCollectionInfo, error)
	GetSummaryInfo(sequence int64) (*SummaryInfo, error)
	GetSummary(tokenCfg *params.TokenConfig, sequence int64) (*Summary, error)
	GetSummarysBySequenceRange(tokenCfg *params.TokenConfig, start, end int64) (SummaryIter, error)
}

type SwapEventIter interface {
	Next(*SwapEvent) bool
}

type SummaryIter interface {
	Next(*Summary) bool
}
