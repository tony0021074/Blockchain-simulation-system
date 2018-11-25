package main

import (
	"./blockchain"
)

func main() {

	nodeID1 := "ABC"
	nodeID2 := "ABDEE"
	blockchain.SaveNode(nodeID1)
	blockchain.SaveNode(nodeID2)

}
