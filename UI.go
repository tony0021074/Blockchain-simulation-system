package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	// Set Host and Port
	// 1st argument is the node's own port, 2nd argumnet is nearby noce's port.
	serverHost := "localhost"
	serverPort := "3000"
	nearbyHost := "localhost"
	nearbyPort := "3001"
	if len(os.Args) == 2 {
		serverPort = os.Args[1]
	}
	if len(os.Args) == 3 {
		serverPort = os.Args[1]
		nearbyPort = os.Args[2]
	}
	userID := serverPort
	serverAddr, err := net.ResolveTCPAddr("tcp", serverHost+":"+serverPort)
	errorMsg(err)
	nearbyAddr, err := net.ResolveTCPAddr("tcp", nearbyHost+":"+nearbyPort)
	errorMsg(err)

	// Choose Function - Either be a nodecontroller, or a miner
	// **Becauses Peer2Peer model (not Server & Client model) is need if a node is miner and nodecontroller at the same time.
	// **Need to make the a TCP socket "Dial" and "Listen" in simultaneously.
	// **By default, Peer2Peer socket is not supported in golang. Need to use external library.
	fmt.Println("Self Node port", serverPort, "; Nearby Node port ", nearbyPort)
	fmt.Println("Enter 10 to become a Node Server")
	fmt.Println("Enter 20 to become a Miner")
	fmt.Println("Enter 30 to calculate a Merkle Tree Root")
	var input string
	fmt.Scanln(&input)
	switch input {
	case "10" /*	Node  Mode */ :

		// Choose Node's action - either update Blockchain, or Start Server Service
		fmt.Println("- Enter 11 to Start acting as a server")
		fmt.Scanln(&input)

		switch input {

		case "11" /*Node - As a server*/ :

			// Initialize by loading blockchain from Database
			var selfNodeChain Blockchain
			selfNodeChain.LoadFromDB(userID)
			fmt.Println("Node:	Blockchain at local Database:")
			selfNodeChain.PrintChain()

			// Listening from Miner
			listener, err := net.ListenTCP("tcp", serverAddr)
			errorMsg(err)
			fmt.Println("Node:	Server Listening on port", serverPort)

			// Create new socket if a connection is accepted
			// golang allows multiple connection by default (non-blocking)
			for {
				conn, err := listener.Accept()
				errorMsg(err)
				go handleMsg(conn, selfNodeChain)
			}

		}

	case "20" /* Miner Mode */ :
		// Connect to nearby node
		conn, err := net.DialTCP("tcp", serverAddr, nearbyAddr)
		errorMsg(err)
		fmt.Printf("Miner:	Connection %s <--> %s\n", serverAddr.String(), nearbyAddr.String())

		// Choose Miner's action - either mining, or check transaction data
		fmt.Println("- Enter 21 to Mine")
		fmt.Println("- Enter 22 to Retrive all Block Hashes at nearby node")
		fmt.Println("- Enter 23 to Check if a block exists among all nodes by using a block hash")
		fmt.Println("- Enter 24 to Check if a data  exists among all nodes by using a Merkle Tree Root")
		fmt.Scanln(&input)

		switch input {

		case "21" /*Miner - Mining*/ :
			// Request PrevBlockHash
			fmt.Println("Miner:	Request PrevBlockHash from Node")
			message := []byte("addBK")
			prevBlockHashFromNode := minerSendMsg(conn, message)
			fmt.Printf("Miner:	Received %x\n", prevBlockHashFromNode)

			// Get Data from user, and build a new Block
			dataToPack := minerGetData()
			fmt.Println("Miner:	...mining...")
			newBlock := CreateBlock(dataToPack, prevBlockHashFromNode)
			pow := NewProofOfWork(newBlock)
			if pow.ValidateBlock() == true {
				fmt.Println("Miner:	Success! Block information here:")
				minerPrintBlock(newBlock)
			}

			// Serialize block using "encoding/json", then add the action indicator
			fmt.Println("Miner:	Now send the Block to nearby node.")
			newBlockJSON, _ := json.Marshal(newBlock)
			message = bytes.Join([][]byte{[]byte("addBK"), newBlockJSON}, []byte{})
			fmt.Println("Miner:	Result - ", string(minerSendMsg(conn, message)))

			conn.Close()
			break

		case "22" /*Miner - Check Block Hashes*/ :
			// Request BlockChain
			fmt.Println("Miner:	Request Full Block Hashes from Node")
			message := []byte("getBC")
			blockHashesFromNode := minerSendMsg(conn, message)
			fmt.Printf("Miner:	Received Block Hashes\n")
			conn.Close()

			// BlockChain is in JSON. Need to decode.
			fmt.Printf("Miner:	...Decoding Block Hashes...\n")
			var blockHashes Blockchain
			err = json.Unmarshal(blockHashesFromNode, &blockHashes)

			// Print BlockChain
			fmt.Println("Miner:	Block Hashes from nearby Node")
			blockHashes.PrintChain()
			fmt.Println("Miner:	Is the Block Hashes valid? -", blockHashes.ValidateChain())
			break

		case "23" /*Miner - Check Single Block*/ :
			// Request BlockChain
			fmt.Print("Miner:	Please input the Block Hash here ")
			fmt.Scanln(&input)
			fmt.Printf("Miner:	Request the Block with Hashes %s\n", input)
			message, _ := hex.DecodeString(input)
			message = bytes.Join([][]byte{[]byte("getBK"), message}, []byte{})
			targetBlockFromNode := minerSendMsg(conn, message)
			fmt.Printf("Miner:	Received the Block\n")
			conn.Close()

			// BlockChain is in JSON. Need to decode.
			fmt.Printf("Miner:	...Decoding the Block...\n")
			var targetBlock Blockchain
			err = json.Unmarshal(targetBlockFromNode, &targetBlock)

			// Print BlockChain
			if len(targetBlock.Blocks) > 0 {
				fmt.Println("Miner:	Target Block is found")
				targetBlock.PrintChain()
			} else {
				fmt.Println("Miner:	Target Block is not found")
			}
			break

		case "24" /*Miner - Check Single Data*/ :
			// Request BlockChain
			fmt.Print("Miner:	Please input the Merkle Tree Root here ")
			fmt.Scanln(&input)
			fmt.Printf("Miner:	Request the Block with Merkle Tree Root %s\n", input)
			message, _ := hex.DecodeString(input)
			message = bytes.Join([][]byte{[]byte("getTX"), message}, []byte{})
			targetBlockFromNode := minerSendMsg(conn, message)
			fmt.Printf("Miner:	Received the Block\n")
			conn.Close()

			// BlockChain is in JSON. Need to decode.
			fmt.Printf("Miner:	...Decoding the Block...\n")
			var targetBlock Blockchain
			err = json.Unmarshal(targetBlockFromNode, &targetBlock)

			// Print BlockChain
			if len(targetBlock.Blocks) > 0 {
				fmt.Println("Miner:	Target Block is found")
				targetBlock.PrintChain()
			} else {
				fmt.Println("Miner:	Target Block is not found")
			}
			break
		}
		break

	case "30" /*Calculated Merkle Tree Root*/ :
		var dataRaw string
		var dataString []string
		fmt.Println("Tree:	Enter data to be used for Merkle Tree Calculation (seperated by ',')")
		fmt.Scan(&dataRaw)
		dataString = strings.Split(dataRaw, ",")
		fmt.Printf("Tree:	The Merkle Tree Root is %x\n", MerkleTreeRoot(dataString))
		break

	}
}

func errorMsg(err error) {
	if err != nil {
		fmt.Println("Connection Error:	", err)
		os.Exit(1)
	}
}
