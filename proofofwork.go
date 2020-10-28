package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	block *Block

	//来存储哈希值，它内置一些方法cmp：比较方法
	//SetBytes：把bytes转成big.int类型
	//SetString:把string转成big.int类型
	target *big.Int //系统提供的，是固定的
}

const Bits = 16

func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	//写难度值，应该是推导，先写成固定，等会推导

	//固定难度值
	/*
		16进制字符串
		targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
		var bigIntTmp big.Int
		bigIntTmp.SetString(targetStr, 16)
		pow.target = &bigIntTmp
	*/

	//推导难度值
	// 0001000000000000000000000000000000000000000000000000000000000000
	//初始化
	// 0000000000000000000000000000000000000000000000000000000000000001
	//向左移动16次，256位
	//10000000000000000000000000000000000000000000000000000000000000000
	//向右移动4次，16位
	// 0001000000000000000000000000000000000000000000000000000000000000
	bitIntTmp := big.NewInt(1)
	//bitIntTmp.Lsh(bitIntTmp,256)
	//bitIntTmp.Rsh(bitIntTmp,16)

	bitIntTmp.Lsh(bitIntTmp, 256-Bits)

	pow.target = bitIntTmp

	return &pow

}

//这是pow的运算函数，为了获取挖矿的随机数，同时返回区块的哈希值
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//获取block数据
	//拼接nonce
	//sha256
	//与难度值比较
	var nonce uint64

	//block := pow.block

	var hash [32]byte

	for {

		fmt.Printf("%x\r", hash)

		hash = sha256.Sum256(pow.prepareData(nonce))

		//将hash数组类型转成big.int,然后与pow.target进行比较，需要引入局部变量
		var bigIntTmp big.Int
		bigIntTmp.SetBytes(hash[:])

		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if bigIntTmp.Cmp(pow.target) == -1 {
			//此时x<y,挖矿成功
			fmt.Printf("挖矿成功！nonce：%d,哈希值为：%x\n", nonce, hash)
			break
		} else {
			nonce++
		}
	}

	return hash[:], nonce

}

func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	block := pow.block

	tmp := [][]byte{
		uintToByte(block.Version),
		block.PrevBlockHash,
		block.MerkleRoot,
		uintToByte(block.TimeStamp),
		uintToByte(block.Difficulity),
		block.Data,
		uintToByte(nonce),
	}

	data := bytes.Join(tmp, []byte{})
	return data

}

//校验挖矿结果
func (pow *ProofOfWork) IsValid() bool {

	//block := pow.block
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)

	var tmp big.Int
	tmp.SetBytes(hash[:])

	return tmp.Cmp(pow.target) == -1

}
