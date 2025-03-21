package data

import (
	"Go-Minichain/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"strconv"
)

// UTXO 定义了一个未花费的交易输出（UTXO）结构体。
type UTXO struct {
	walletAddress string // 接收方的钱包地址
	amount        int    // 该 UTXO 所包含的金额
	publicKeyHash []byte // 接收方公钥的哈希值，用于验证所有权
	used          bool   // 该 UTXO 是否已经被使用
}

// NewUTXO 创建一个新的 UTXO 实例。
// 参数:
// - address: 接收方的钱包地址。
// - amount: 该 UTXO 所包含的金额。
// - publicKey: 接收方的公钥。
// 返回值:
// 返回一个指向新创建的 UTXO 实例的指针。
func NewUTXO(address string, amount int, publicKey ecdsa.PublicKey) *UTXO {
	publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
	return &UTXO{
		walletAddress: address,
		amount:        amount,
		publicKeyHash: utils.Ripemd160Digest(utils.Sha256Digest(publicKeyBytes)),
		used:          false,
	}
}

// UnlockScript 验证签名是否正确，并检查公钥哈希是否匹配。
// 参数:
// - sign: 签名数据。
// - publicKey: 公钥。
// 返回值:
// 返回布尔值，表示签名和公钥是否通过验证。
func (utxo *UTXO) UnlockScript(sign []byte, publicKey ecdsa.PublicKey) bool {
	stack := make([][]byte, 0)
	// 将签名压入栈中。
	stack = append(stack, sign)
	// 将公钥的字节序列压入栈中。
	publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
	stack = append(stack, publicKeyBytes)
	// 复制栈顶元素（公钥字节序列）。
	stack = append(stack, stack[len(stack)-1])
	// 弹出栈顶元素并计算其哈希值后重新压入栈中。
	data := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	stack = append(stack, utils.Ripemd160Digest(utils.Sha256Digest(data)))
	// 将先前保存的公钥哈希值压入栈中。
	stack = append(stack, utxo.publicKeyHash)
	// 比较栈顶两个哈希值是否相同。
	publicKeyHash1 := stack[len(stack)-1]
	publicKeyHash2 := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	if !bytes.Equal(publicKeyHash1, publicKeyHash2) {
		return false
	}
	// 验证签名是否正确。
	publicKeyEncoded := stack[len(stack)-1]
	sign1 := stack[len(stack)-2]
	return utils.Verify(publicKeyEncoded, sign1, &publicKey)
}

// SetUsed 标记该 UTXO 已被使用。
// 如果该 UTXO 已经被标记为已使用，则会触发 panic。
func (utxo *UTXO) SetUsed() {
	if utxo.used {
		panic("UTXO already used")
	}
	utxo.used = true
}

// IsUsed 检查该 UTXO 是否已被使用。
// 返回值:
// 返回布尔值，表示该 UTXO 是否已被使用。
func (utxo *UTXO) IsUsed() bool {
	return utxo.used
}

// GetWalletAddress 获取接收方的钱包地址。
// 返回值:
// 返回字符串类型的钱包地址。
func (utxo *UTXO) GetWalletAddress() string {
	return utxo.walletAddress
}

// GetAmount 获取该 UTXO 所包含的金额。
// 返回值:
// 返回整数类型的金额。
func (utxo *UTXO) GetAmount() int {
	return utxo.amount
}

// GetPublicKeyHash 获取接收方公钥的哈希值。
// 返回值:
// 返回字节数组类型的公钥哈希值。
func (utxo *UTXO) GetPublicKeyHash() []byte {
	return utxo.publicKeyHash
}

// ToString 将 UTXO 转换为字符串格式。
// 返回值:
// 返回字符串类型的 UTXO 表示。
func (utxo *UTXO) ToString() string {
	return "UTXO{" +
		"walletAddress=" + utxo.walletAddress + "," +
		"amount=" + strconv.Itoa(utxo.amount) + "," +
		"publicKeyHash=" + utils.Byte2HexString(utxo.publicKeyHash) +
		"}"
}

// UTXO2Bytes 将输入和输出的 UTXO 列表序列化为字节数组。
// 参数:
// - inUTXO: 输入的 UTXO 列表，表示未花费的交易输出集合。
// - outUTXO: 输出的 UTXO 列表，表示新的未花费交易输出集合。
// 返回值:
// 返回字节数组类型的序列化结果。
func UTXO2Bytes(inUTXO []*UTXO, outUTXO []*UTXO) []byte {
	return []byte(fmt.Sprint(inUTXO) + fmt.Sprint(outUTXO))
}
