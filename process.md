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

## 添加NewProofOfWork.go
```go
package main

import "math/big"

type ProofOfWork struct {
	block *Block

	//来存储哈希值，它内置一些方法cmp：比较方法
	//SetBytes：把bytes转成big.int类型
	//SetString:把string转成big.int类型
	target *big.Int //系统提供的，是固定的
}

func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	//写难度值，应该是推导，先写成固定，等会推导

	//16进制字符串
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"

	var bigIntTmp big.Int
	bigIntTmp.SetString(targetStr, 16)

	pow.target = &bigIntTmp

	return &pow

}
```

## run函数实现
```go
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//获取block数据
	//拼接nonce
	//sha256
	//与难度值比较
	var nonce uint64

	//block := pow.block

	var hash [32]byte

	for {

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
```

## 使用pow更新newblock函数
```go
//创建区块，对Block的每个字段填充数据
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

	//block.SetHash()
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce

	return &block
}
```

## 校验挖矿是否有效
```go
//校验挖矿结果
func (pow *ProofOfWork) IsValid() bool {

	//block := pow.block
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)

	var tmp big.Int
	tmp.SetBytes(hash[:])

	return tmp.Cmp(pow.target) == -1

}
```

## 打印block字段
```go

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
```

## 使用bits推导难度值
``` go

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
```

## boltDemo
```go
package main

import (
	"bolt"
	"fmt"
	"log"
)

func main() {
	db, err := bolt.Open("test.db", 0600, nil)
	//向数据库中写入数据
	//从数据库中读取数据

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		//所有操作在这里
		b1 := tx.Bucket([]byte("bucketname1"))

		if b1 == nil {
			//如果b1为空，说明该桶不存在，需要创建
			b1, err = tx.CreateBucket([]byte("bucketname1"))
			if err != nil {
				log.Panic(err)
			}
		}
		//bucket已经创建完成，准备写入数据
		//写数据使用Put，读数据使用Get
		b1.Put([]byte("name1"), []byte("Lily"))
		if err != nil {
			fmt.Printf("写入数据失败name1：Lily\n")
		}
		b1.Put([]byte("name2"), []byte("Jim"))
		if err != nil {
			fmt.Printf("写入数据失败name2：Jim\n")
		}

		//读取数据
		name1 := b1.Get([]byte("name1"))
		name2 := b1.Get([]byte("name2"))
		name3 := b1.Get([]byte("name3"))

		fmt.Printf("name1:%s\n", name1)
		fmt.Printf("name2:%s\n", name2)
		fmt.Printf("name3:%s\n", name3)

		return nil

	})

}
```

