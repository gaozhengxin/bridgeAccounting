package mongodb

import (
	"gopkg.in/mgo.v2"

	"github.com/gaozhengxin/bridgeAccounting/params"
)

var (
	collSyncInfo  *mgo.Collection
	collDeposits  = make(map[string]*mgo.Collection)
	collRedeemeds = make(map[string]*mgo.Collection)
	collMints     = make(map[string]*mgo.Collection)
	collBurns     = make(map[string]*mgo.Collection)
)

func collDeposit(tokenCfg *params.TokenConfig) *mgo.Collection {
	return collDeposits[tokenCfg.PairID]
}

func collRedeemed(tokenCfg *params.TokenConfig) *mgo.Collection {
	return collRedeemeds[tokenCfg.PairID]
}

func collMint(tokenCfg *params.TokenConfig) *mgo.Collection {
	return collMints[tokenCfg.PairID]
}

func collBurn(tokenCfg *params.TokenConfig) *mgo.Collection {
	return collBurns[tokenCfg.PairID]
}

// do this when reconnect to the database
func deinintCollections(scanConfig *params.ScanConfig) {
	collSyncInfo = database.C(tbSyncInfo)
	for _, tk := range scanConfig.Tokens {
		collDeposits[tk.PairID] = database.C(tbDeposit(tk))
		collRedeemeds[tk.PairID] = database.C(tbRedeemed(tk))
		collMints[tk.PairID] = database.C(tbMint(tk))
		collBurns[tk.PairID] = database.C(tbBurn(tk))
	}
}

func initCollections(scanConfig *params.ScanConfig) {
	initCollection(tbSyncInfo, collSyncInfo)
	for _, tk := range scanConfig.Tokens {
		initCollection(tbDeposit(tk), collDeposit(tk))
		initCollection(tbRedeemed(tk), collRedeemed(tk))
		initCollection(tbMint(tk), collMint(tk))
		initCollection(tbBurn(tk), collBurn(tk))
	}
}

func initCollection(table string, collection *mgo.Collection, indexKey ...string) {
	collection = database.C(table)
	if len(indexKey) != 0 && indexKey[0] != "" {
		_ = (collection).EnsureIndexKey(indexKey...)
	}
}
