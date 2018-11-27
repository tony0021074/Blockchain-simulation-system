package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

func minerSendMsg(conn net.Conn, msg []byte) (reply []byte) {
	_, err := conn.Write(msg)
	if err != nil {
		fmt.Println("Miner:	Error writing:")
		fmt.Println("Miner:	", err)
		return
	}

	fmt.Println("Miner:	...sending message to nearby node")

	buf := make([]byte, 8192)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Miner:	...Error Reading:")
		fmt.Println("Miner:	...", err)
		return
	}

	fmt.Println("Miner:	...received message from nearby node")
	return bytes.TrimRight(buf, "\x00")
}

func minerGetData() []string {
	var dataRaw string
	var dataString []string
	fmt.Println("Miner:	Enter data to be packed in blockchain (seperated by ',')")
	fmt.Scan(&dataRaw)
	dataString = strings.Split(dataRaw, ",")
	return dataString
}

func minerPrintBlock(block *Block) {
	fmt.Printf("Miner:	 > TimeStamp		: %d\n", block.Timestamp)
	fmt.Printf("Miner:	 > PrevBlockHash	: %x\n", block.PrevBlockHash)
	fmt.Printf("Miner:	 > Merkle Root		: %x\n", block.MerkleTreeRoot)
	fmt.Printf("Miner:	 > Nonce		: %d\n", block.Nonce)
	fmt.Printf("Miner:	 > CurrBlockHash	: %x\n", block.CurrentBlockHash)
	if len(block.Data) > 0 {
		fmt.Printf("Miner:	 > Data			: %s\n", block.Data)
	} else {
		fmt.Printf("Miner:	 > Data			: Not Disclose\n")
	}
	return
}
