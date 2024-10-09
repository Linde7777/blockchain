package main

import (
	"blockchain/blockchain"
	"blockchain/utils"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
	chain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - print the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.chain.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CommandLine) printChain() {
	iter := cli.chain.NewIterator()
	for {
		block := iter.Next()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()
	const (
		cmdKeyAdd   = "add"
		cmdKeyPrint = "print"
	)
	addBlockCMD := flag.NewFlagSet(cmdKeyAdd, flag.ExitOnError)
	printChainCMD := flag.NewFlagSet(cmdKeyPrint, flag.ExitOnError)
	addBlockData := addBlockCMD.String("block", "", "block data")

	switch os.Args[1] {
	case cmdKeyAdd:
		err := addBlockCMD.Parse(os.Args[2:])
		utils.HandlePanic(err)
	case cmdKeyPrint:
		err := printChainCMD.Parse(os.Args[2:])
		utils.HandlePanic(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCMD.Parsed() {
		if *addBlockData == "" {
			addBlockCMD.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCMD.Parsed() {
		cli.printChain()
	}
}

func main() {
	chain := blockchain.InitBlockChain()
	cli := CommandLine{chain: chain}
	cli.run()
}
