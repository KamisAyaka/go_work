package main

import "fmt"

func main() {
	wallet := NewWallet()
	wallet1 := NewWallet()
	fmt.Printf("Wallet address: %s\n", wallet.GetAddress())
	fmt.Printf("Wallet1 address: %s\n", wallet1.GetAddress())
	bc := CreateBlockchain(string(wallet.GetAddress()))
	defer bc.db.Close()
	bc.getBalance(string(wallet.GetAddress()))
	bc.getBalance(string(wallet1.GetAddress()))
	bc.Send(string(wallet.GetAddress()), string(wallet1.GetAddress()), 10, "send 10 coins to wallet1", wallet)
	bc.Send(string(wallet.GetAddress()), string(wallet1.GetAddress()), 10, "send 10 coins to wallet1", wallet)
	bc.Send(string(wallet.GetAddress()), string(wallet1.GetAddress()), 10, "send 10 coins to wallet1", wallet)
	bc.getBalance(string(wallet.GetAddress()))
	bc.getBalance(string(wallet1.GetAddress()))
}
