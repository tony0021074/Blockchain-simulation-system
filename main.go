package main

import (
"fmt"
"strconv"
"./blockchain"
)

func main() {
	chain, _ := blockchain.LoadChain()
	if chain == nil {
		chain = blockchain.NewBlockchain()
	}

	chain.AddBlock("1st block")
	chain.AddBlock("2nd block")
	chain.AddBlock("3nd block")

	for _, block := range chain.Blocks {

		fmt.Printf("Previous Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.CurrentBlockHash)

		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

	}
}
