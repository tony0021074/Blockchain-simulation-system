package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

func handleMsg(conn net.Conn, selfNodeChain Blockchain) {

	fmt.Printf("Node:	<%s> Connection established \n", conn.RemoteAddr().String())

	// Receive Message, maximum length is 8kB
	bufReceive := make([]byte, 8192)
	bufSend := make([]byte, 8192)
	_, err := conn.Read(bufReceive)
	if err == io.EOF {
		fmt.Println("Node:	Error reading:")
		fmt.Println("Node:	", err)
	}
	if err != nil {
		fmt.Println("Node:	Error reading:")
		fmt.Println("Node:	", err)
	}

	// Choose action depending on message header
	request := string(bytes.TrimRight(bufReceive, "\x00"))[0:5]
	payload := bytes.TrimRight(bufReceive, "\x00")[5:]

	if request == "addBK" {
		// "addBK":
		//	1. Return prevBlock to miner by conn.Write()
		fmt.Printf("Node:	<%s> Miner would like to add a block to blockchain\n", conn.RemoteAddr().String())
		selfNodeChain.LoadFromDB(selfNodeChain.UserID)
		bufSend = selfNodeChain.Blocks[len(selfNodeChain.Blocks)-1].CurrentBlockHash
		fmt.Printf("Node:	<%s> Return PrevBlockHash to Miner\n", conn.RemoteAddr().String())
		if bufSend != nil {
			_, err = conn.Write(bufSend)
		}

		//	2. Receive newBlock from miner by conn.Read()
		fmt.Printf("Node:	<%s> Waiting for new block\n", conn.RemoteAddr().String())
		_, err = conn.Read(bufReceive)
		var newBlock *Block
		err = json.Unmarshal(bytes.TrimRight(bufReceive, "\x00")[5:], &newBlock)

		//	3. Add the block to blockchain. Update Blockchain before adding
		selfNodeChain.LoadFromDB(selfNodeChain.UserID)
		pow := NewProofOfWork(newBlock)
		failFlag := false
		if pow.ValidateBlock() == false {
			failFlag = true
		} else {
			failFlag = !selfNodeChain.AddBlock(newBlock)
		}

		//  4. Return result to Miner. Need to LoadFromDB to update the block in RAM, to see if new block is added.
		selfNodeChain.LoadFromDB(selfNodeChain.UserID)
		if failFlag == false {
			bufSend = []byte("Success. Blockchain is updated.")
			_, err = conn.Write(bufSend)
			fmt.Println("Node:	Blockchain now:")
			selfNodeChain.PrintChain()
		} else {
			bufSend = []byte("Fail. Invalid Block.")
			_, err = conn.Write(bufSend)
		}

	} else {

		// Other request : getBC or getBK or getTX
		//
		//	getBC:		Return full block hash.
		//				i.e. getblocks() in Project Specification
		//
		//	getBK/	:	Return a block header / data if User provides an ID.
		//	getTX		In our case,
		//					ID is Merkle Tree Root for getTX;
		//					ID is currentBlockHash for getBK;
		//				i.e. getdata() in Project Specification, handle with payload.Type = "block" payload.Type = "tx" in ppt slide.
		result := handleInv(request, payload, conn, selfNodeChain)
		bufSend, _ = json.Marshal(result)
		_, err = conn.Write(bufSend)
		fmt.Printf("Node:	<%s> Return information to client.\n", conn.RemoteAddr().String())

	}

	conn.Close()
	fmt.Printf("Node:	<%s> Connection is closed\n", conn.RemoteAddr().String())

}

func handleInv(request string, payload []byte, conn net.Conn, selfNodeChain Blockchain) Blockchain {

	var resultChain Blockchain

	if request == "getBC" {
		// If "getBC" is detected, return a blockchain with headers only
		fmt.Printf("Node:	<%s> Client would like to retrive all block hashes\n", conn.RemoteAddr().String())
		selfNodeChain.LoadFromDB(selfNodeChain.UserID)
		resultChain.UserID = selfNodeChain.UserID
		for i := 0; i < len(selfNodeChain.Blocks); i++ {
			resultChain.Blocks = append(resultChain.Blocks, &Block{
				Timestamp:        selfNodeChain.Blocks[i].Timestamp,
				PrevBlockHash:    selfNodeChain.Blocks[i].PrevBlockHash,
				MerkleTreeRoot:   selfNodeChain.Blocks[i].MerkleTreeRoot,
				Nonce:            selfNodeChain.Blocks[i].Nonce,
				CurrentBlockHash: selfNodeChain.Blocks[i].CurrentBlockHash,
			})
		}
	} else if request == "getBK" {
		// If "getBK" is detected, return a blockchain with a single block == target block
		fmt.Printf("Node:	<%s> Client would like to check if a block exists\n", conn.RemoteAddr().String())
		fmt.Printf("Node:	<%s> The Block Hash is %x\n", conn.RemoteAddr().String(), payload)
		selfNodeChain.LoadFromDB(selfNodeChain.UserID)
		resultChain.UserID = selfNodeChain.UserID
		for i := 0; i < len(selfNodeChain.Blocks); i++ {
			if string(selfNodeChain.Blocks[i].CurrentBlockHash) == string(payload) {
				resultChain.Blocks = append(resultChain.Blocks, &Block{
					Timestamp:        selfNodeChain.Blocks[i].Timestamp,
					PrevBlockHash:    selfNodeChain.Blocks[i].PrevBlockHash,
					MerkleTreeRoot:   selfNodeChain.Blocks[i].MerkleTreeRoot,
					Nonce:            selfNodeChain.Blocks[i].Nonce,
					CurrentBlockHash: selfNodeChain.Blocks[i].CurrentBlockHash,
				})
				if len(resultChain.Blocks) > 0 {
					fmt.Printf("Node:	<%s> Target Block is found\n", conn.RemoteAddr().String())
					break
				}
			}
		}

	} else if request == "getTX" {
		// If "getTX" is detected, return a blockchain with a single block == target block
		fmt.Printf("Node:	<%s> Client would like to check if a data exists\n", conn.RemoteAddr().String())
		fmt.Printf("Node:	<%s> The Merkle Tree Root is %x\n", conn.RemoteAddr().String(), payload)
		selfNodeChain.LoadFromDB(selfNodeChain.UserID)
		resultChain.UserID = selfNodeChain.UserID
		for i := 0; i < len(selfNodeChain.Blocks); i++ {
			if string(selfNodeChain.Blocks[i].MerkleTreeRoot) == string(payload) {
				resultChain.Blocks = append(resultChain.Blocks, &Block{
					Timestamp:        selfNodeChain.Blocks[i].Timestamp,
					PrevBlockHash:    selfNodeChain.Blocks[i].PrevBlockHash,
					MerkleTreeRoot:   selfNodeChain.Blocks[i].MerkleTreeRoot,
					Nonce:            selfNodeChain.Blocks[i].Nonce,
					CurrentBlockHash: selfNodeChain.Blocks[i].CurrentBlockHash,
				})
				if len(resultChain.Blocks) > 0 {
					fmt.Printf("Node:	<%s> Target Block is found\n", conn.RemoteAddr().String())
					break
				}
			}
		}
	}
	if len(resultChain.Blocks) == 0 {
		fmt.Printf("Node:	<%s> No result\n", conn.RemoteAddr().String())
	}
	return resultChain
}
