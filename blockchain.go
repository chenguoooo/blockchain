package main

//创建区块链，使用Block数组模拟
type BlockChain struct {
	Blocks []*Block
}

//实现创建区块链方法
func NewBlockChain() *BlockChain {
	//在创建的时候添加一个区块：创世块
	genesisBlock := NewBlock(genesisInfo, []byte{0x0000000000000000})

	bc := BlockChain{Blocks: []*Block{genesisBlock}}
	return &bc

}

//添加区块

func (bc *BlockChain) AddBlock(data string) {
	//1.创建一个区块

	//bc.Blocks的最后一个区块的Hash值就是当前新去爱的PrevBlockHash
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	prevHash := lastBlock.Hash

	block := NewBlock(data, prevHash)
	//2.添加到bc.Blocks数组中
	bc.Blocks = append(bc.Blocks, block)
}
