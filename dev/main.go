package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	token "github.com/gaozhengxin/bridgeAccount/dev/token"
)

func main() {
	decimal := getDecimal()
	fmt.Println(decimal)
	bal := getBalance()
	fmt.Println(bal)
}

func getDecimal() int {
	client, err := ethclient.Dial("https://rpcapi.fantom.network")
	if err != nil {
		panic(err)
	}
	//methodDecimal := common.FromHex("0x95d89b41")
	methodDecimal := common.FromHex("0x313ce567")
	toAddr := common.HexToAddress("0xd67de0e0a0fd7b15dc8348bb9be742f3c5850454")
	msg := ethereum.CallMsg {
		To: &toAddr,
		Data: methodDecimal[:],
	}
	bs, err := client.CallContract(context.Background(), msg, nil)
	decimal := new(big.Int).SetBytes(bs[0:32]).Uint64()
	return int(decimal)
}

func getBalance() *big.Int {
	client, err := ethclient.Dial("https://rpcapi.fantom.network")
	if err != nil {
		panic(err)
	}

	toAddr := common.HexToAddress("0xd67de0e0a0fd7b15dc8348bb9be742f3c5850454")

	instance, err := token.NewToken(toAddr, client)
	if err != nil {
        panic(err)
    }
	address := common.HexToAddress("0x0536806df512d6cdde913cf95c9886f65b1d3462")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
    if err != nil {
        panic(err)
    }
	return bal
}