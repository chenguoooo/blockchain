# blockchain

## 创建区块，实现区块

``` go
package main

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	PrevBlockHash []byte //前区块哈希
	Hash          []byte //当前区块哈希
	Data          []byte //数据，目前使用字节流
}

const genesisInfo = "hello blockchain"

//创建区块，对Block的每个字段填充数据
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{}, //先填充为空，后续会填充数据
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
```

## 实现setHash函数

``` go
//为了生成区块哈希，实现一个简单的函数，来计算哈希值，没有随机值，没有难度值
func (block *Block) SetHash() {
	var data []byte
	data = append(data, block.PrevBlockHash...)
	data = append(data, block.Data...)

	hash /*[32]byte*/ := sha256.Sum256(data)
	block.Hash = hash[:]
}
```

## 在NewBlock中调用

``` go
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{}, //先填充为空，后续会填充数据
		Data:          []byte(data),
	}

	block.SetHash()

	return &block
}
```

## 区块链的定义及使用
``` go
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

func main() {
	//fmt.Printf("hello world")
	//block := NewBlock(genesisInfo, []byte{0x0000000000000000})
	bc := NewBlockChain()

	for _, block := range bc.Blocks {
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%s\n", block.Data)
	}
}
```
## 添加区块

``` go
func (bc *BlockChain) AddBlock(data string) {
	//1.创建一个区块

	//bc.Blocks的最后一个区块的Hash值就是当前新去爱的PrevBlockHash
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	prevHash := lastBlock.Hash

	block := NewBlock(data, prevHash)
	//2.添加到bc.Blocks数组中
	bc.Blocks = append(bc.Blocks, block)
}
```
## 重构代码
添加block.go和blockchain.go

## 更新补充区块字段
```go
type Block struct {
	Version uint64 //区块版本号

	PrevBlockHash []byte //前区块哈希

	MerkleRoot []byte //先填写为空

	TimeStamp uint64 //从1970.1.1至今

	Difficulity uint64 //挖矿难度

	Nonce uint64 //随机数

	Data []byte //数据，目前使用字节流

	Hash []byte //当前区块哈希,区块中本不存在

}
```

## 更新newblock函数
``` go
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:       00,
		PrevBlockHash: prevBlockHash,
		MerkleRoot:    []byte{},
		TimeStamp:     uint64(time.Now().Unix()),
		Difficulity:   10,       //先随便写
		Nonce:         10,       //先随便写
		Hash:          []byte{}, //先填充为空，后续会填充数据
		Data:          []byte(data),
	}

	block.SetHash()

	return &block
}
```
## 更新setHash函数
``` go
func (block *Block) SetHash() {
	var data []byte
	//uintToByte()将数字转成[]byte{},在utils.go实现
	data = append(data, uintToByte(block.Version)...)
	data = append(data, block.PrevBlockHash...)
	data = append(data, block.MerkleRoot...)
	data = append(data, uintToByte(block.TimeStamp)...)
	data = append(data, uintToByte(block.Difficulity)...)
	data = append(data, block.Data...)
	data = append(data, uintToByte(block.Nonce)...)

	hash /*[32]byte*/ := sha256.Sum256(data)
	block.Hash = hash[:]
}
```
## 添加空函数uintToByte
创建新文件utils.go,内容入下
``` go
package main

//这是一个工具函数文件

func uintToByte(num uint64) []byte {
	//TODO
	//具体实现后面再写
	return []byte{}

}

```
## 编码逻辑实现
```go
package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func uintToByte(num uint64) []byte {
	//使用binary.Write来进行编码
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()

}
```
## 使用bytes.join改写函数
``` go
func (block *Block) SetHash() {
	/*
		var data []byte
		//uintToByte()将数字转成[]byte{},在utils.go实现
		data = append(data, uintToByte(block.Version)...)
		data = append(data, block.PrevBlockHash...)
		data = append(data, block.MerkleRoot...)
		data = append(data, uintToByte(block.TimeStamp)...)
		data = append(data, uintToByte(block.Difficulity)...)
		data = append(data, block.Data...)
		data = append(data, uintToByte(block.Nonce)...)
	*/
	tmp := [][]byte{
		uintToByte(block.Version),
		block.PrevBlockHash,
		block.MerkleRoot,
		uintToByte(block.TimeStamp),
		uintToByte(block.Difficulity),
		block.Data,
		uintToByte(block.Nonce),
	}
	data := bytes.Join(tmp, []byte{})

	hash /*[32]byte*/ := sha256.Sum256(data)
	block.Hash = hash[:]
}
```
