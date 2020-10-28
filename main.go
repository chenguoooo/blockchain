package main

import (
	"fmt"
)

func main() {
	//fmt.Printf("hello world")
	//block := NewBlock(genesisInfo, []byte{0x0000000000000000})
	bc := NewBlockChain()
	bc.AddBlock("this is sec block")
	bc.AddBlock("this is third block")

	for i, block := range bc.Blocks {
		fmt.Printf("++++++++++++%d+++++++++++++++++\n", i)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%s\n", block.Data)
	}
}
