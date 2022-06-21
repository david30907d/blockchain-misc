package main

import (
	"github.com/david30907d/blockchain-misc/golang-blockchain/blockchain"
)

func main() {
	blockchain.InitBlockChain()

	// chain.AddBlock("First Block after Genesis")
	// chain.AddBlock("Second Block after Genesis")
	// chain.AddBlock("Third Block after Genesis")

	// for idx, block := range chain.Blocks {
	// 	fmt.Printf("Previous Hash: %x\n", block.PrevHash)
	// 	fmt.Printf("Data in Block: %x\n", block.Data)
	// 	fmt.Printf("Hash: %x\n", block.Hash)
	// 	fmt.Printf("=======Block %d======\n", idx)
	// }
	// target := big.NewInt(1)
	// fmt.Println("%d", target)
	// target.Lsh(target, uint(2))
	// fmt.Println("%d", target)
}
