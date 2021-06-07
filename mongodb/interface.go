package mongodb

import (
	"github.com/gaozhengxin/bridgeaudit/params"
)

type QueryAPI interface {
	BaseQueryAPI
	AuditQueryAPI
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
	GetDepositsByBlockRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error)
	GetDepositsByTimeRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error)
	GetDepositByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) ([]*SwapEvent, error)

	GetMint(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error)
	GetMintByBlockRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error)
	GetMintByTimeRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error)
	GetMintByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) ([]*SwapEvent, error)

	GetBurn(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error)
	GetBurnByBlockRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error)
	GetBurnByTimeRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error)
	GetBurnByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) ([]*SwapEvent, error)

	GetRedeemed(tokenCfg *param.TokenConfig, txhash string) (*SwapEvent, error)
	GetRedeemedByBlockRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error)
	GetRedeemedByTimeRange(tokenCfg *param.TokenConfig, start, end int64) ([]*SwapEvent, error)
	GetRedeemedByUserTimeRange(tokenCfg *param.TokenConfig, user string, start, end int64) ([]*SwapEvent, error)
}

type AuditAPI interface {
	BaseQueryAPI
	AuditQueryAPI
	AddCheckRange(tokenCfg *param.TokenConfig) error
	UpdateCheckRangeAccDeposit(tokenCfg *param.TokenConfig, value float64) error
	UpdateCheckRangeAccMint(tokenCfg *param.TokenConfig, value float64) error
	UpdateCheckRangeAccBurn(tokenCfg *param.TokenConfig, value float64) error
	UpdateCheckRangeAccRedeemed(tokenCfg *param.TokenConfig, value float64) error
	AddCheckRangeInfo() error
	UpdateCheckInfo(int64) error
}

type AuditQueryAPI interface {
	GetSheetCollectionInfo() (*SheetCollectionInfo, error)
	GetSheetInfo(sequence int64) (*SheetInfo, error)
	GetSheet(tokenCfg *param.TokenConfig, sequence int64) (*Sheet, error)
	GetSheetsBySequenceRange(tokenCfg *param.TokenConfig, start, end int64) ([]*Sheet, error)
}
