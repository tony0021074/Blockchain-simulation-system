package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)
var(
	maxNonce=math.MaxInt32 //nance的最大值，整数64位里面最大的，可以说无限大了。
)

const targetBits = 20

//工作量证明的结构，对block的计算要满足target
type ProofOfWork struct{
	block  *Block
	target *big.Int
}


//把上面的区块传进来
func NewProofOfWork(b *Block) *ProofOfWork{
	target :=big.NewInt(1)    //target设置为整数1
	target.Lsh(target,uint(256 - targetBits))//Lsh是一个移位操作，target本身为1，移位后前面的20个bits会变成0，
	                                         // 即为20/4=5bytes

	pow :=&ProofOfWork{block:b,target:target}

	return pow
}

func(pow *ProofOfWork) prepareData(nonce int) []byte{
	data :=bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash, //固定值
			pow.block.Data,            //固定值
			IntToHex(pow.block.Timestamp),//固定值
			IntToHex(int64(targetBits)),//固定值20
			IntToHex(int64(nonce)),//只有nonce是变量，从0开始增加到maxInt32
		},
		[]byte{},
		)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte){
	var hashInt big.Int
	var hash [32]byte
	nonce :=0

	fmt.Printf("Mining the block containing \"%s\"\n",pow.block.Data)
	for nonce <maxNonce{
		data  :=pow.prepareData(nonce)//prepareData里面算出来的Data，一大块包括timestamp,data,preblockhash

		hash =sha256.Sum256(data)//用sha256的方法计算出data的hash值
		fmt.Printf("\r%x",hash)
		hashInt.SetBytes(hash[:])//把hash值转换为hash整数

		if hashInt.Cmp(pow.target)==-1{//hashInt整数与target作对比，cmp一个内置对比函数
			break                    //对比成功就break
		} else {
			nonce ++                 //对比不成功就继续nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
	}



//校验算法，验证hash值不是随随便便出来的，是经过pow算法验证过的。
func (pow *ProofOfWork) Validate() bool{
	var hashInt big.Int

	data :=pow.prepareData(pow.block.Nonce)
	hash :=sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid :=hashInt.Cmp(pow.target)==-1

	return isValid
}


func IntToHex(num int64)  []byte{
	buff :=new(bytes.Buffer)
	err  :=binary.Write(buff,binary.BigEndian,num)
	if err !=nil{
		log.Panic(err)
	}
	return buff.Bytes()
}


func DataToHash(data []byte) []byte{
	hash :=sha256.Sum256(data)
	return hash[:]
}