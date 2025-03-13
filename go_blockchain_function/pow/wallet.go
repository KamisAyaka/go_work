package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"os"

	"golang.org/x/crypto/ripemd160"
)

const ChecksumLen = 4
const version = byte(0x00)
const walletFile = "wallet.dat"

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	return &Wallet{private, public}
}

// 创建钱包
func NewWalletByFile() *Wallet {
	//随机生成秘钥对
	var w Wallet
	w.LoadFromFile()

	return &w
}

// 生成秘钥对
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

// 获取公钥哈希
func HashPubKey(pubKey []byte) []byte {
	//1. 先hash一次
	publicSHA256 := sha256.Sum256(pubKey)
	//2. 计算  ripemd160
	RIPEMD160Hasher := ripemd160.New()
	RIPEMD160Hasher.Write(publicSHA256[:])

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// 生成校验和,输入为版本号+公钥哈希
func checksum(payload []byte) []byte {

	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:ChecksumLen]
}

// 获取地址
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// 保存钱包信息到文件
func (w Wallet) SaveToFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// 加载钱包文件
func (w *Wallet) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallet Wallet
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallet)
	if err != nil {
		log.Panic(err)
	}

	*w = wallet

	return nil
}

// 验证地址
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-ChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-ChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(actualChecksum, targetChecksum)
}
