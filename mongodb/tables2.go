package mongodb

import (
	"github.com/jowenshaw/gethscan/params"
)

const (
	tbSheetInfo string = "SheetInfo"
	tbSheetCollectionInfo string = "SheetCollectionInfo"
)

func tbSheet(tokenCfg *param.TokenConfig) string {
	return "CheckRange_" + tokenCfg.PairID
}

type Sheet struct {
	Sequence int64 `bson:"_id"`
	AccDeposit float64
	AccMint float64
	AccBurn float64
	AccRedeemed float64
}

type SheetInfo struct {
	Sequence int64 `bson:"_id"`
	Tag string `bson:"tag"` // for example, date
	SrcStartHeight int64 `bson:"src_start_height"`
	SrcEndHeight int64 `bson:"src_end_height"`
	DstStartHeight int64 `bson:"dst_start_height"`
	DstEndHeight int64 `bson:"dst_end_height"`
}

type SheetCollectionInfo struct {
	LatestSequence int64 `bson:"_id"`
}