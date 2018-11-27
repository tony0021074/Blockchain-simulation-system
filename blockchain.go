package main

import (
	"fmt"
)

//Blockchain : Define object Blockchain
type Blockchain struct {
	UserID string
	Blocks []*Block
}

// AddBlock :	Add a block in blockchain - array of block.
func (bc *Blockchain) AddBlock(newBlock *Block) bool {
	preBlock := bc.Blocks[len(bc.Blocks)-1]
	if string(newBlock.PrevBlockHash) == string(preBlock.CurrentBlockHash) {
		bc.Blocks = append(bc.Blocks, newBlock)
		SaveBlock(newBlock, bc.UserID)
		return true
	}
	fmt.Println("Chain:	Failed to add block. Invalid Hash.")
	return false

}

// LoadFromDB :	Load Blockchain from Database. Initialise to Genesis Block if the blockchain is empty.
func (bc *Blockchain) LoadFromDB(userID string) {
	bc.UserID = userID
	bc.Blocks, _ = LoadChain(bc.UserID)
	if len(bc.Blocks) == 0 {
		genesisBlock := CreateBlock([]string{"New", "Genesis", "Block"}, []byte("00000000000000000000000000000000"))
		bc.Blocks = []*Block{genesisBlock}
		SaveBlock(genesisBlock, bc.UserID)
	}
}

// PrintChain :	Print all blocks in blockchain
func (bc *Blockchain) PrintChain() {
	for i := 0; i < len(bc.Blocks); i++ {
		fmt.Printf("Chain:	Block #%d	-PrevHash	%x\n", i, bc.Blocks[i].PrevBlockHash)
		fmt.Printf("Chain:	Block #%d	-CurrHash	%x\n", i, bc.Blocks[i].CurrentBlockHash)
		fmt.Printf("Chain:	Block #%d	-MerkleRoot	%x\n", i, bc.Blocks[i].MerkleTreeRoot)
		if len(bc.Blocks[i].Data) > 0 {
			fmt.Printf("Chain:	Block #%d	-Data		%s\n", i, bc.Blocks[i].Data)
		} else {
			fmt.Printf("Chain:	Block #%d	-Data		%s\n", i, "not disclose")
		}

	}
}

// ValidateChain :	Check if the chain is valid
func (bc *Blockchain) ValidateChain() bool {

	validFlag := true
	if len(bc.Blocks) == 0 {
		validFlag = false
	} else {

		// Check Genesis Block first:	Only Check CurrBlockHash is Valid
		if NewProofOfWork(bc.Blocks[0]).ValidateBlock() == false {
			validFlag = false
		}
		// Then check other Blocks :	Check CurrBlockHash is Valid & PrevBlockHash Matches
		if len(bc.Blocks) > 1 {
			for i := 1; i < len(bc.Blocks); i++ {
				if NewProofOfWork(bc.Blocks[i]).ValidateBlock() == false {
					validFlag = false
				}
				if string(bc.Blocks[i-1].CurrentBlockHash) != string(bc.Blocks[i].PrevBlockHash) {
					validFlag = false
				}
			}
		}
	}
	return validFlag
}
