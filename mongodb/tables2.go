package mongodb

import (
	"github.com/gaozhengxin/bridgeaudit/params"
)

const (
	tbSummaryInfo string = "SummaryInfo"
	tbSummaryCollectionInfo string = "SummaryCollectionInfo"
)

func tbSummary(tokenCfg *param.TokenConfig) string {
	return "CheckRange_" + tokenCfg.PairID
}

type Summary struct {
	Sequence int64 `bson:"_id"`
	AccDeposit float64
	AccMint float64
	AccBurn float64
	AccRedeemed float64
}

type SummaryInfo struct {
	Sequence int64 `bson:"_id"`
	Tag string `bson:"tag"` // for example, date
	SrcStartHeight int64 `bson:"src_start_height"`
	SrcEndHeight int64 `bson:"src_end_height"`
	DstStartHeight int64 `bson:"dst_start_height"`
	DstEndHeight int64 `bson:"dst_end_height"`
}

type SummaryCollectionInfo struct {
	LatestSequence int64 `bson:"_id"`
}