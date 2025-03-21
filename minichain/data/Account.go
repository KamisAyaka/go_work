package data

import (
	"Go-Minichain/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/asn1"
)

// Account 表示一个账户，包含公钥和私钥。
// 公钥用于验证交易签名，私钥用于生成签名。
type Account struct {
	PublicKey  ecdsa.PublicKey   // 公钥，用于验证交易签名
	PrivateKey *ecdsa.PrivateKey // 私钥，用于生成交易签名
}

// NewAccount 创建并初始化一个新的Account实例。
// 使用utils包中的Secp256k1Generate函数生成公钥和私钥对。
// 返回值是一个指向Account结构的指针。
func NewAccount() *Account {
	privateKey, publicKey := utils.Secp256k1Generate()
	return &Account{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

// GetWalletAddress 生成并返回该账户的钱包地址。
// 钱包地址的生成过程如下：
// 1. 对公钥进行SHA-256哈希，然后对结果进行RIPEMD-160哈希。
// 2. 在哈希结果前添加版本前缀（通常为0x00）。
// 3. 对上述数据进行两次SHA-256哈希，并取前4个字节作为校验码。
// 4. 将版本前缀、公钥哈希和校验码组合在一起，并进行Base58编码。
func (a *Account) GetWalletAddress() string {
	publicKey := a.GetPublicKey()
	publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
	publicKeyHash := utils.Ripemd160Digest(utils.Sha256Digest(publicKeyBytes))
	data := make([]byte, 1+len(publicKeyHash))
	data = append(data, 0)                // 添加版本前缀
	data = append(data, publicKeyHash...) // 添加公钥哈希

	doubleHash := utils.Sha256Digest(utils.Sha256Digest(data))
	wallEncoded := make([]byte, 1+len(publicKeyHash)+4)
	wallEncoded = append(wallEncoded, 0)                 // 添加版本前缀
	wallEncoded = append(wallEncoded, publicKeyHash...)  // 添加公钥哈希
	wallEncoded = append(wallEncoded, doubleHash[:4]...) // 添加校验码

	b := utils.NewBase58Util()
	walletAddress := b.Encode(wallEncoded) // 进行Base58编码
	return walletAddress
}

// ToString 返回该账户的字符串表示形式。
// 包括公钥和私钥的十六进制编码。
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

// GetPublicKey 返回该账户的公钥。
func (a *Account) GetPublicKey() ecdsa.PublicKey {
	return a.PublicKey
}

// GetPrivateKey 返回该账户的私钥。
func (a *Account) GetPrivateKey() *ecdsa.PrivateKey {
	return a.PrivateKey
}

// GetAmount 计算并返回指定UTXO列表中属于该账户的总金额。
// 参数 trueUtxo 是一个UTXO指针切片，表示未花费的交易输出。
// 函数通过遍历trueUtxo列表，累加每个UTXO的金额来计算总金额。
func (a *Account) GetAmount(trueUtxo []*UTXO) int {
	amount := 0
	for _, utxo := range trueUtxo {
		amount += utxo.amount
	}
	return amount
}
