package accounting

import (
	"github.com/gaozhengxin/bridgeAccounting/params"
	"github.com/gaozhengxin/bridgeAccounting/mongodb"
)

type AccountingAPI interface {
	AccountingQueryAPI
	MakeSummaryInfo(tag string, srcStartHeight, srcEndHeight, dstStartHeight, dstEndHeight int64) (*mongodb.SummaryInfo, error)
	MakeSummary(*params.TokenConfig) (*mongodb.SummaryInfo, error)
}

type AccountingQueryAPI interface {
	GetSummaryInfo()
	GetSummaryInfoByTag()
	GetSummary(*params.TokenConfig)
	GetSummarysByTimeRange(tk *params.TokenConfig, start, end int64) ([]*mongodb.SummaryInfo)
}