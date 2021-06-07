package mongodb

import (
	"github.com/gaozhengxin/bridgeaudit/params"
)

/*
Swap events on src chain
1. transfer to MPC account (Deposit)
4. MPC transfer to account (Redeemed)

Swap events on dst chain
2. swapin (Mint)
3. swapout (Burn)
*/

const (
	tbSyncInfo string = "SyncInfo"
)

func tbDeposit(tokenCfg *param.TokenConfig) string {
	return "Deposit_" + tokenCfg.PairID
}

func tbRedeemed(tokenCfg *param.TokenConfig) string {
	return "Redeemed_" + tokenCfg.PairID
}

func tbMint(tokenCfg *param.TokenConfig) string {
	return "Mint_" + tokenCfg.PairID
}

func tbBurn(tokenCfg *param.TokenConfig) string {
	return "Burn_" + tokenCfg.PairID
}

type SyncInfo struct {
	SrcChainSyncedHeight int64 `bson:"src_synced_height"`
	SrcChainStartHeight int64 `bson:"src_start_height"`
	DstChainSyncedHeight int64 `bson:"dst_synced_height"`
	DstChainStartHeight int64 `bson:"dst_start_height"`
}

type SwapEvent struct {
	TxHash string `bson:"_id"`
	Timestamp int64 `bson:"timestamp"`
	BlockNumber int64 `bson:"block_number"`
	Amount string `bson:"amount"`
	FAmount float64 `bson:"famount"`
	User string `bson:"user"`
}
