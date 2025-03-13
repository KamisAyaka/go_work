package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	hash := sha256.Sum256([]byte("welcome to blockchain"))
	fmt.Printf("%x\n", hash)
}
