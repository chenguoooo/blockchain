package main

import "fmt"

type Block struct {
	PrevBlockHash []byte
	Hash          []byte
	Data          []byte
}

const genesisInfo = "hello world"

//创建区块，对Block的每个字段填充数据
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{}, //先填充为空，后续回填充数据
		Data:          []byte(data),
	}
	return &block
}

func main() {
	//fmt.Printf("hello world")
	block := NewBlock(genesisInfo, []byte{0x0000000000000000})
	fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
	fmt.Printf("Hash:%x\n", block.Hash)
	fmt.Printf("Data:%s\n", block.Data)
}
