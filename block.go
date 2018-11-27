package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"
)

//Block : Define object Block
type Block struct {
	// Block Header
	Timestamp      int64
	PrevBlockHash  []byte
	MerkleTreeRoot []byte
	Nonce          int
	// Block Data
	Data [][]byte
	// Block hash,  can be computed using header
	CurrentBlockHash []byte
}

//CreateBlock : Create new Block
func CreateBlock(dataString []string, prevBlockHash []byte) *Block {

	time.Sleep(1 * time.Second)

	var data [][]byte
	for i := 0; i < len(dataString); i++ {
		data = append(data, []byte(dataString[i]))
	}
	block := &Block{
		Timestamp:        time.Now().Unix(),
		PrevBlockHash:    prevBlockHash,
		MerkleTreeRoot:   MerkleTreeRoot(dataString),
		Data:             data,
		CurrentBlockHash: []byte{},
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run() //工作量证明过程，看func run中的具体操作

	block.CurrentBlockHash = hash[:]
	block.Nonce = nonce

	return block
}

var (
	maxNonce = math.MaxInt32 //nonce的最大值，整数64位里面最大的，可以说无限大了。
)

const targetBits = 8

//ProofOfWork : 工作量证明的结构，对block的计算要满足target
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

//NewProofOfWork : 把上面的区块传进来
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)                  //target设置为整数1
	target.Lsh(target, uint(256-targetBits)) //Lsh是一个移位操作，target本身为1，移位后前面的20个bits会变成0，
	// 即为20/4=5bytes

	pow := &ProofOfWork{block: b, target: target}

	return pow
}

//prepareData : Transfer the header elements in a Block to a single byte string header
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,       //固定值
			pow.block.MerkleTreeRoot,      //固定值
			IntToHex(pow.block.Timestamp), //固定值
			IntToHex(int64(targetBits)),   //固定值20
			IntToHex(int64(nonce)),        //只有nonce是变量，从0开始增加到maxInt32
		},
		[]byte{},
	)
	return data
}

//Run : Start creating Block by running Proof of Work
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < maxNonce {
		data := pow.prepareData(nonce) //prepareData里面算出来的Data，一大块包括timestamp,data,preblockhash

		hash = sha256.Sum256(data) //用sha256的方法计算出data的hash值
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:]) //把hash值转换为hash整数

		if hashInt.Cmp(pow.target) == -1 { //hashInt整数与target作对比，cmp一个内置对比函数
			break //对比成功就break
		} else {
			nonce++ //对比不成功就继续nonce++
		}
	}
	fmt.Print("\n")
	return nonce, hash[:]
}

//ValidateBlock : 校验算法，验证hash值不是随随便便出来的，是经过pow算法验证过的。
func (pow *ProofOfWork) ValidateBlock() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

//IntToHex : Rename function binary.Write() to IntToHex for easier to read
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
