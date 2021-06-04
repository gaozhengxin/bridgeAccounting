package mongodb

import (
	"gopkg.in/mgo.v2"
	"github.com/jowenshaw/gethscan/params"
)

var (
	collSheetInfo *mgo.Collection
	collSheetCollectionInfo *mgo.Collection
	collSheets = make(map[string]*mgo.Collection)
)

func collSheet(tokenCfg *param.TokenConfig) *mgo.Collection {
	return collSheets[tokenCfg.PairID]
}

// do this when reconnect to the database
func deinintCollections2() {
	collSheetInfo = database.C(tbSheetInfo)
	collSheetCollectionInfo = database.C(tbSheetCollectionInfo)
	for _, tk := range scanConfig.Tokens {
		collSheet(tk) = database.C(tbSheet(tk))
	}
}

func initCollections2(scanConfig *ScanConfig) {
	initCollection(tbSheetInfo, collSheetInfo)
	initCollection(tbSheetCollectionInfo, collSheetCollectionInfo)
	for _, tk := range scanConfig.Tokens {
		initCollection(tbSheet(tk), collSheet(tk))
	}
}