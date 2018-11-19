package blockchain

type Blockchain struct{
	Blocks []*Block  //一个blocks数组里面放的全是block
}


//往Blocks数组里面新加block的方法
func (bc *Blockchain)AddBlock(data string){
	preBlock :=bc.Blocks[len(bc.Blocks)-1]
	newBlock :=NewBlock(data,preBlock.CurrentBlockHash)//NewBlock在block.go里面已经定义，newBlock里面存了preBlock的Hash.
	bc.Blocks =append(bc.Blocks,newBlock)
	SaveBlock(newBlock)
}


func NewBlockchain() *	Blockchain{
	genesisBlock := NewGenesisBlock()
	SaveBlock(genesisBlock)
	return  &Blockchain{[]*Block{genesisBlock}}												//传到block里面NewGenesisBlock函数
	                                                        // 中创建系统默认的Genesis block值
}
