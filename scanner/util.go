package scanner

import (
	"context"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/gaozhengxin/bridgeAccounting/params"
	"github.com/gaozhengxin/bridgeAccounting/mongodb"
)

func convertToMgoSwapEvent(swapEvent *SwapEvent, decimal int) *mongodb.SwapEvent {
	return &mongodb.SwapEvent{
		TxHash: strings.ToLower(swapEvent.TxHash.String()),
		BlockTime: swapEvent.BlockTime,
		BlockNumber: swapEvent.BlockNumber.Int64(),
		Amount: swapEvent.Amount.String(),
		FAmount: toFloat(swapEvent.Amount, decimal),
		User: strings.ToLower(swapEvent.User.String()),
	}
}

func (scanner *ethSwapScanner) cachedDecimal(tokenCfg *params.TokenConfig) int {
	if tokenCfg.Decimal != 0 {
		return tokenCfg.Decimal
	}
	if tokenCfg.TokenAddress == "" {
		// native
		return 18
	}
	methodDecimal := common.FromHex("0x313ce567")
	msg := ethereum.CallMsg {
		Data: methodDecimal[:],
	}
	*msg.To = common.HexToAddress(tokenCfg.TokenAddress)
	bs, err := scanner.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0
	}
	decimal := new(big.Int).SetBytes(bs[0:32]).Uint64()
	tokenCfg.Decimal = int(decimal)
	return int(decimal)
}

func toFloat(bigint *big.Int, decimal int) float64 {
	divider := new(big.Int).SetInt64(10)
	divider = new(big.Int).Exp(divider, big.NewInt(int64(decimal)), nil)
	quo := new(big.Int)
	rem := new(big.Int)
	quo, rem = quo.DivMod(bigint, divider, rem)
	return float64(quo.Int64()) + float64(rem.Int64())/math.Pow10(decimal)
}