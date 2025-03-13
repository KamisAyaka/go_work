package data

import (
	"Go-Minichain/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"strconv"
)

type UTXO struct {
	walletAddress string
	amount        int
	publicKeyHash []byte
	used          bool
}

func NewUTXO(address string, amount int, publicKey ecdsa.PublicKey) *UTXO {
	publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
	return &UTXO{
		walletAddress: address,
		amount:        amount,
		publicKeyHash: utils.Ripemd160Digest(utils.Sha256Digest(publicKeyBytes)),
		used:          false,
	}
}

func (utxo *UTXO) UnlockScript(sign []byte, publicKey ecdsa.PublicKey) bool {
	stack := make([][]byte, 0)
	// <sign> 签名入栈
	// <sign>
	stack = append(stack, sign)
	// <publicKey> 入栈
	// <sign> <publicKey>
	publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
	stack = append(stack, publicKeyBytes)
	// dup 复制栈顶元素
	// <sign> <publicKey> <publicKey>
	stack = append(stack, stack[len(stack)-1])
	// 弹出栈顶元素，并计算哈希后将其入栈
	// <sign><publicKey><publicKeyHash>
	data := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	stack = append(stack, utils.Ripemd160Digest(utils.Sha256Digest(data)))
	// 将先前保存的公钥哈希入栈
	// <sign><publicKey><publicKeyHash><publicKeyHash>
	stack = append(stack, utxo.publicKeyHash)
	// 比较栈顶两个哈希是否相同
	// <sign><publicKey>
	publicKeyHash1 := stack[len(stack)-1]
	publicKeyHash2 := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	if !bytes.Equal(publicKeyHash1, publicKeyHash2) {
		return false
	}
	// 验证签名是否正确，直接返回结果不入栈了
	publicKeyEncoded := stack[len(stack)-1]
	sign1 := stack[len(stack)-2]
	//stack = stack[:len(stack)-2]

	return utils.Verify(publicKeyEncoded, sign1, &publicKey)
}

func (utxo *UTXO) SetUsed() {
	if utxo.used {
		panic("UTXO already used")
	}
	utxo.used = true
}
func (utxo *UTXO) IsUsed() bool {
	return utxo.used
}

func (utxo *UTXO) GetWalletAddress() string {
	return utxo.walletAddress
}

func (utxo *UTXO) GetAmount() int {
	return utxo.amount
}

func (utxo *UTXO) GetPublicKeyHash() []byte {
	return utxo.publicKeyHash
}

func (utxo *UTXO) ToString() string {
	return "UTXO{" +
		"walletAddress=" + utxo.walletAddress + "," +
		"amount=" + strconv.Itoa(utxo.amount) + "," +
		"publicKeyHash=" + utils.Byte2HexString(utxo.publicKeyHash) +
		"}"
}

// UTXO2Bytes 将输入和输出的UTXO列表序列化为字节数组。
//
// 参数:
// - inUTXO: 输入的UTXO列表，表示未花费的交易输出（UTXO）集合。
// - outUTXO: 输出的UTXO列表，表示新的未花费交易输出集合。
//
// 返回值:
// - []byte: 将输入和输出的UTXO列表通过字符串格式化后转换为字节数组的结果。
//
// 注意: 该函数使用 fmt.Sprint 将 UTXO 列表直接转换为字符串并返回其字节形式，
// 可能不适用于需要精确序列化的场景。
func UTXO2Bytes(inUTXO []*UTXO, outUTXO []*UTXO) []byte {
	// 使用 fmt.Sprint 将输入和输出的 UTXO 列表格式化为字符串，并转换为字节数组。
	return []byte(fmt.Sprint(inUTXO) + fmt.Sprint(outUTXO))
}
