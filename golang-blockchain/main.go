package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/david30907d/blockchain-misc/golang-blockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) AddBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Printf("Added %s to block\n", data)
}

func (cli *CommandLine) PrintChain() {
	iterator := cli.blockchain.Iterator()
	for {
		prevBlock := iterator.IterBackWard()
		fmt.Printf("Previous hash: %x\n", prevBlock.Hash)
		fmt.Printf("Data: %s\n", prevBlock.Data)
		fmt.Printf("Hash: %x\n", prevBlock.Hash)
		proof := blockchain.NewProof(prevBlock)
		fmt.Printf("Validate: %s\n", strconv.FormatBool(proof.Validate()))
		fmt.Println()
		if len(iterator.CurrentHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}
func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	// chain.AddBlock("new block1")
	defer chain.Database.Close()
	cli := CommandLine{chain}
	cli.run()
}
