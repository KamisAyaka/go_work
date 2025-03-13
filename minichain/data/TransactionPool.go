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

func (p *TransactionPool) GetNewTransaction() *Transaction {
	accounts := p.blockchain.GetAccount()
	var transaction *Transaction
	for {
		aAccount := accounts[rand.Intn(len(accounts))]
		bAccount := accounts[rand.Intn(len(accounts))]

		if aAccount == bAccount {
			continue
		}

		aWalletAddress := aAccount.GetWalletAddress()
		bWalletAddress := bAccount.GetWalletAddress()

		aTrueUTXOs := p.blockchain.GetTrueUTXOs(aWalletAddress)
		aAmount := aAccount.GetAmount(aTrueUTXOs)

		if aAmount == 0 {
			continue
		}

		txAmount := rand.Intn(aAmount) + 1
		inUTXOs := make([]*UTXO, 0)
		outUTXOs := make([]*UTXO, 0)

		publicKey := aAccount.GetPublicKey()
		publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
		aUnLockSign := utils.Signature(publicKeyBytes, aAccount.GetPrivateKey())
		//fmt.Println(utils.Verify(publicKeyBytes, aUnLockSign, &publicKey))

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
			continue
		}

		outUTXOs = append(outUTXOs, NewUTXO(bWalletAddress, txAmount, bAccount.GetPublicKey()))

		if inAmount > txAmount {
			outUTXOs = append(outUTXOs, NewUTXO(aWalletAddress, inAmount-txAmount, aAccount.GetPublicKey()))
		}

		data := UTXO2Bytes(inUTXOs, outUTXOs)
		sign := utils.Signature(data, aAccount.GetPrivateKey())
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
