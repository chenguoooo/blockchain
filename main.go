package main

import (
	"fmt"
	"time"
)

func main() {
	//fmt.Printf("hello world")
	//block := NewBlock(genesisInfo, []byte{0x0000000000000000})
	bc := NewBlockChain()
	bc.AddBlock("1_chen")
	bc.AddBlock("2_chen")

	for i, block := range bc.Blocks {
		fmt.Printf("++++++++++++%d+++++++++++++++++\n", i)
		fmt.Printf("Version:%d\n", block.Version)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("MerkleRoot:%x\n", block.MerkleRoot)

		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp:%s\n", timeFormat)

		fmt.Printf("Difficulity:%d\n", block.Difficulity)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%s\n", block.Data)

		pow := NewProofOfWork(block)
		fmt.Printf("Isvalid:%v\n", pow.IsValid())
	}

}
