package mongodb

import (
	"gopkg.in/mgo.v2"
	"github.com/jowenshaw/gethscan/params"
)

var (
	collSyncInfo *mgo.Collection
	collDeposits = make(map[string]*mgo.Collection)
	collRedeemeds = make(map[string]*mgo.Collection)
	collMints = make(map[string]*mgo.Collection)
	collBurns = make(map[string]*mgo.Collection)
)

func collDeposit(tokenCfg *param.TokenConfig) *mgo.Collection {
	return collDeposits[tokenCfg.PairID]
}

func collRedeemed(tokenCfg *param.TokenConfig) *mgo.Collection {
	return collRedeemeds[tokenCfg.PairID]
}

func collMint(tokenCfg *param.TokenConfig) *mgo.Collection {
	return collMints[tokenCfg.PairID]
}

func collBurn(tokenCfg *param.TokenConfig) *mgo.Collection {
	return collBurns[tokenCfg.PairID]
}

// do this when reconnect to the database
func deinintCollections() {
	collSyncInfo = database.C(tbSyncInfo)
	for _, tk := range scanConfig.Tokens {
		collDeposit(tk) = database.C(tbDeposit(tk))
		collRedeemed(tk) = database.C(tbRedeemed(tk))
		collMint(tk) = database.C(tbMint(tk))
		collBurn(tk) = database.C(tbBurn(tk))
	}
}

func initCollections(scanConfig *ScanConfig) {
	initCollection(tbSyncInfo, &collSyncInfo)
	for _, tk := range scanConfig.Tokens {
		initCollection(tbDeposit(tk), collDeposit(tk))
		initCollection(tbRedeemed(tk), collRedeemed(tk))
		initCollection(tbMint(tk), collMint(tk))
		initCollection(tbBurn(tk), collBurn(tk))
	}
}

func initCollection(table string, collection **mgo.Collection, indexKey ...string) {
	*collection = database.C(table)
	if len(indexKey) != 0 && indexKey[0] != "" {
		_ = (*collection).EnsureIndexKey(indexKey...)
	}
}
