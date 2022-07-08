package main

import (
	"os"

	"github.com/david30907d/blockchain-misc/golang-blockchain/cli"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
