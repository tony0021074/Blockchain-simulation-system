package blockchain

import (
	"time"
)

type Block struct{
	Timestamp             int64
	Data                  []byte
	PrevBlockHash         []byte
	CurrentBlockHash      []byte
	Nonce                 int
	MerkleTreeRoot        []byte
}

func NewBlock(data string,prevBlockHash []byte) *Block {
	block := &Block{Timestamp: time.Now().Unix(), Data: []byte(data),PrevBlockHash: prevBlockHash, CurrentBlockHash: []byte{}}
	pow :=NewProofOfWork(block)
	nonce,hash :=pow.Run()//工作量证明过程，看func run中的具体操作

	block.CurrentBlockHash=hash[:]
	block.Nonce=nonce

	return block
}

//创世纪块
func NewGenesisBlock() *Block{
	return NewBlock("Genesis Blcok",[]byte{})//第一个值是被系统默认的data“Genesis Blcok”,不像NewBlock里面
	// 定义的，它没有preBlockHash，再看看newblock中往Genesis后
	//添加区块的过程，同时在里面还要进行工作量证明
}