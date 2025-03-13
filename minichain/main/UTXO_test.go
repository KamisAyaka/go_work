package main

import (
	"Go-Minichain/consensus"
	"Go-Minichain/data"
	"Go-Minichain/utils"
	"testing"
)

func TestUTXOtest(t *testing.T) {
	UTXOtest()
}

func UTXOtest() {
	blockchain := data.NewBlockChain()
	transactionPool := data.NewTransactionPool(2, blockchain)

	miner := consensus.NewMinerNode(transactionPool, blockchain)
	transaction := GetOneTransaction(blockchain)
	transactionPool.Put(*transaction)
	miner.Run()

}

func GetOneTransaction(blockchain *data.BlockChain) *data.Transaction {
	// 获取账户列表
	accounts := blockchain.GetAccount()

	// 明确指定账户A（索引1）和账户B（索引2）
	aAccount := accounts[1]
	bAccount := accounts[2]

	// 获取账户A的可用UTXO
	aWalletAddress := aAccount.GetWalletAddress()
	aTrueUTXOs := blockchain.GetTrueUTXOs(aWalletAddress)

	// 构造输入UTXO（假设账户A有足够UTXO）
	inUTXOs := aTrueUTXOs[:1] // 取第一个可用UTXO

	// 构造输出UTXO：1000 BTC给B，剩余找零给A（假设输入UTXO总金额>=1000）
	outUTXOs := []*data.UTXO{
		data.NewUTXO(bAccount.GetWalletAddress(), 1000, bAccount.GetPublicKey()),
	}
	if inAmount := aAccount.GetAmount(inUTXOs); inAmount > 1000 {
		outUTXOs = append(outUTXOs,
			data.NewUTXO(aWalletAddress, inAmount-1000, aAccount.GetPublicKey()))
	}

	// 签名交易数据
	dataBytes := data.UTXO2Bytes(inUTXOs, outUTXOs)
	sign := utils.Signature(dataBytes, aAccount.GetPrivateKey())

	// 创建确定性的交易
	transaction := data.NewTransaction(inUTXOs, outUTXOs, sign, aAccount.GetPublicKey())
	blockchain.ProcessTransactionUTXO(inUTXOs, outUTXOs)

	return transaction
}
