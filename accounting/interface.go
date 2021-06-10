package accounting

import (
	"github.com/gaozhengxin/bridgeaudit/params"
	"github.com/gaozhengxin/bridgeaudit/mongodb"
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
	GetSummarysByTimeRange(*params.TokenConfig, start, end int64) ([]*mongodb.SummaryInfo)
}