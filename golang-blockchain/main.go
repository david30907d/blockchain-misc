package main

import (
	"fmt"
	"os"

	"github.com/david30907d/blockchain-misc/golang-blockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()
	iter := chain.Iterator()

	fmt.Println(iter.Next().Hash)
}
