package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

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

const genesisInfo = "0_chen"

//创建区块，对Block的每个字段填充数据
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:       00,
		PrevBlockHash: prevBlockHash,
		MerkleRoot:    []byte{},
		TimeStamp:     uint64(time.Now().Unix()),
		Difficulity:   Bits, //先随便写
		//Nonce:         10,       //先随便写
		Hash: []byte{}, //先填充为空，后续会填充数据
		Data: []byte(data),
	}

	//block.SetHash()
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce

	return &block
}

//序列化，将区块转换成字节流
func (block *Block) Serialize() []byte {

	var buffer bytes.Buffer

	//定义编码器
	encoder := gob.NewEncoder(&buffer)

	//编码器对结构进行编码，一定要进行校验
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}
func DeSerialize(data []byte) *Block {

	//fmt.Printf("解码传入的数据:%x\n", data)

	var block Block

	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
