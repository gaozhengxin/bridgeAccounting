package mongodb

import (
	"gopkg.in/mgo.v2"

	"github.com/gaozhengxin/bridgeAccounting/params"
)

var (
	collSummaryInfo           *mgo.Collection
	collSummaryCollectionInfo *mgo.Collection
	collSummarys              = make(map[string]*mgo.Collection)
)

func collSummary(tokenCfg *params.TokenConfig) *mgo.Collection {
	return collSummarys[tokenCfg.PairID]
}

// do this when reconnect to the database
func deinintCollections2(scanConfig *params.ScanConfig) {
	collSummaryInfo = database.C(tbSummaryInfo)
	collSummaryCollectionInfo = database.C(tbSummaryCollectionInfo)
	for _, tk := range scanConfig.Tokens {
		collSummarys[tk.PairID] = database.C(tbSummary(tk))
	}
}

func initCollections2(scanConfig *params.ScanConfig) {
	initCollection(tbSummaryInfo, collSummaryInfo)
	initCollection(tbSummaryCollectionInfo, collSummaryCollectionInfo)
	for _, tk := range scanConfig.Tokens {
		initCollection(tbSummary(tk), collSummary(tk))
	}
}
