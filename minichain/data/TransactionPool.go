package data

import (
	"Go-Minichain/utils"
	"crypto/elliptic"
	"math/rand"
)

/**
 * 交易池
 */

type TransactionPool struct {
	transactions []Transaction
	capacity     int
	blockchain   *BlockChain
}

func NewTransactionPool(c int, chain *BlockChain) *TransactionPool {
	p := new(TransactionPool)
	p.capacity = c
	p.transactions = make([]Transaction, 0)
	p.blockchain = chain
	return p
}
func (p *TransactionPool) Put(transaction Transaction) {
	p.transactions = append(p.transactions, transaction)
}
func (p *TransactionPool) GetAll() []Transaction {
	transactions := p.transactions
	p.clear()
	return transactions
}
func (p *TransactionPool) clear() {
	p.transactions = make([]Transaction, 0)
}
func (p *TransactionPool) IsFull() bool {
	return len(p.transactions) >= p.capacity
}
func (p *TransactionPool) IsEmpty() bool {
	return len(p.transactions) == 0
}
func (p *TransactionPool) GetCapacity() int {
	return p.capacity
}

// GetNewTransaction 生成一个新的交易对象，随机选择两个账户进行交易。
// 该方法会确保交易的有效性，包括检查账户余额、解锁UTXO、签名交易等操作，
// 并将交易的输入和输出更新到区块链的UTXO池中。
//
// 参数:
// - p: 指向 TransactionPool 的指针，包含区块链和交易池的相关信息。
//
// 返回值:
// - *Transaction: 指向新创建的交易对象的指针。
func (p *TransactionPool) GetNewTransaction() *Transaction {
	accounts := p.blockchain.GetAccount() // 获取区块链中的所有账户
	var transaction *Transaction

	for {
		// 随机选择两个不同的账户作为交易的发送方和接收方
		aAccount := accounts[rand.Intn(len(accounts))]
		bAccount := accounts[rand.Intn(len(accounts))]

		if aAccount == bAccount {
			continue // 如果选中了相同的账户，则重新选择
		}

		// 获取发送方和接收方的钱包地址
		aWalletAddress := aAccount.GetWalletAddress()
		bWalletAddress := bAccount.GetWalletAddress()

		// 获取发送方的有效未花费交易输出（UTXO）并计算其总金额
		aTrueUTXOs := p.blockchain.GetTrueUTXOs(aWalletAddress)
		aAmount := aAccount.GetAmount(aTrueUTXOs)

		if aAmount == 0 {
			continue // 如果发送方没有可用余额，则重新选择账户
		}

		// 随机生成交易金额，并初始化输入和输出UTXO列表
		txAmount := rand.Intn(aAmount) + 1
		inUTXOs := make([]*UTXO, 0)
		outUTXOs := make([]*UTXO, 0)

		// 使用发送方的私钥对其公钥进行签名，用于解锁UTXO
		publicKey := aAccount.GetPublicKey()
		publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
		aUnLockSign := utils.Signature(publicKeyBytes, aAccount.GetPrivateKey())

		// 解锁发送方的UTXO，累积足够的输入金额以满足交易需求
		inAmount := 0
		for _, utxo := range aTrueUTXOs {
			if utxo.UnlockScript(aUnLockSign, aAccount.GetPublicKey()) {
				inAmount += utxo.GetAmount()
				inUTXOs = append(inUTXOs, utxo)
				if inAmount >= txAmount {
					break
				}
			}
		}

		if inAmount < txAmount {
			continue // 如果累积的输入金额不足，则重新选择账户
		}

		// 创建接收方的输出UTXO
		outUTXOs = append(outUTXOs, NewUTXO(bWalletAddress, txAmount, bAccount.GetPublicKey()))

		// 如果输入金额大于交易金额，则创建找零的输出UTXO
		if inAmount > txAmount {
			outUTXOs = append(outUTXOs, NewUTXO(aWalletAddress, inAmount-txAmount, aAccount.GetPublicKey()))
		}

		// 将输入和输出UTXO序列化，并使用发送方的私钥对交易数据进行签名
		data := UTXO2Bytes(inUTXOs, outUTXOs)
		sign := utils.Signature(data, aAccount.GetPrivateKey())

		// 创建交易对象，并更新区块链的UTXO池
		transaction = NewTransaction(inUTXOs, outUTXOs, sign, aAccount.GetPublicKey())
		p.blockchain.ProcessTransactionUTXO(inUTXOs, outUTXOs)
		break
	}

	return transaction
}

func (p *TransactionPool) Start() {
	go func() {
		for {
			if !p.IsFull() {
				transaction := p.GetNewTransaction()
				p.Put(*transaction)
			}
		}
	}()
}
