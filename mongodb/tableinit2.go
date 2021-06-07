package mongodb

import (
	"gopkg.in/mgo.v2"
	"github.com/gaozhengxin/bridgeaudit/params"
)

var (
	collSummaryInfo *mgo.Collection
	collSummaryCollectionInfo *mgo.Collection
	collSummarys = make(map[string]*mgo.Collection)
)

func collSummary(tokenCfg *param.TokenConfig) *mgo.Collection {
	return collSummarys[tokenCfg.PairID]
}

// do this when reconnect to the database
func deinintCollections2() {
	collSummaryInfo = database.C(tbSummaryInfo)
	collSummaryCollectionInfo = database.C(tbSummaryCollectionInfo)
	for _, tk := range scanConfig.Tokens {
		collSummary(tk) = database.C(tbSummary(tk))
	}
}

func initCollections2(scanConfig *ScanConfig) {
	initCollection(tbSummaryInfo, collSummaryInfo)
	initCollection(tbSummaryCollectionInfo, collSummaryCollectionInfo)
	for _, tk := range scanConfig.Tokens {
		initCollection(tbSummary(tk), collSummary(tk))
	}
}