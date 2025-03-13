package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到以太坊节点
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 合约地址
	contractAddress := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")

	// 私钥
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}

	// 交易选项
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(31337))
	auth.Value = big.NewInt(0) // 设置为 0，因为这是一个非支付交易

	// 初始化合约实例
	calldemo, err := NewCalldemo(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// 调用 setCount 函数
	tx, err := calldemo.SetCount(auth, big.NewInt(700))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	// 等待交易确认
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction mined: %v\n", receipt)
}
