package data

import (
	"Go-Minichain/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/asn1"
)

type Account struct {
	PublicKey  ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

func NewAccount() *Account {
	privateKey, publicKey := utils.Secp256k1Generate()
	return &Account{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

func (a *Account) GetWalletAddress() string {
	// 公钥哈希：ripemd160(sha256(publicKey))
	publicKey := a.GetPublicKey()
	publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
	publicKeyHash := utils.Ripemd160Digest(utils.Sha256Digest(publicKeyBytes))
	// 0x0 + 公钥哈希
	data := make([]byte, 1+len(publicKeyHash))
	data = append(data, 0)
	data = append(data, publicKeyHash...)

	// 0x0 + 公钥哈希 + 校验(两次哈希取前四个字节)
	doubleHash := utils.Sha256Digest(utils.Sha256Digest(data))
	wallEncoded := make([]byte, 1+len(publicKeyHash)+4)
	wallEncoded = append(wallEncoded, 0)
	wallEncoded = append(wallEncoded, publicKeyHash...)
	wallEncoded = append(wallEncoded, doubleHash[:4]...)
	// base58编码得到钱包地址
	b := utils.NewBase58Util()
	walletAddress := b.Encode(wallEncoded)
	return walletAddress
}

func (a *Account) ToString() string {
	publicKeyBytes, err := asn1.Marshal(a.PublicKey)
	if err != nil {
		panic(err)
	}
	privateKeyBytes, err := asn1.Marshal(a.PrivateKey)
	if err != nil {
		panic(err)
	}
	return "Account{" +
		"publicKey=" + utils.Byte2HexString(publicKeyBytes) + "," +
		"privateKey=" + utils.Byte2HexString(privateKeyBytes) + "}"
}

func (a *Account) GetPublicKey() ecdsa.PublicKey {
	return a.PublicKey
}

func (a *Account) GetPrivateKey() *ecdsa.PrivateKey {
	return a.PrivateKey
}

func (a *Account) GetAmount(trueUtxo []*UTXO) int {
	amount := 0
	for _, utxo := range trueUtxo {
		amount += utxo.amount
	}
	return amount
}
