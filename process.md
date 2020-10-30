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

## 改写区块链，获取区块链实例
```go

//使用blot改写
type BlockChain struct {
	db *bolt.DB //句柄

	tail []byte //最后一个区块hash值
}

//实现创建区块链方法
func NewBlockChain() *BlockChain {
	//在创建的时候添加一个区块：创世块
	//genesisBlock := NewBlock(genesisInfo, []byte{0x0000000000000000})
	//bc := BlockChain{Blocks: []*Block{genesisBlock}}
	//return &bc

	//功能分析
	//1.获得数据库句柄，打开数据库，读写数据
	//判断是否有bucket，如果没有，创建bucket
	//写入创世块
	//写入lasthashkey这条数据
	//更新tail为最后一个区块的哈希
	//返回bc实例
	db, err := bolt.Open("blockChain.db", 0600, nil)
	//向数据库中写入数据
	//从数据库中读取数据

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	var tail []byte

	db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("blockBucket"))

		if b == nil {
			//如果b为空，说明该桶不存在，需要创建
			fmt.Printf("bucket不存在，准备创建!\n")
			b, err = tx.CreateBucket([]byte("blockBucket"))

			if err != nil {
				log.Panic(err)
			}

			//抽屉准备完毕，开始添加创世块
			genesisBlock := NewBlock(genesisInfo, []byte{})
			b.Put(genesisBlock.Hash, genesisBlock.toBytes() /*将区块序列化，转成字节流*/)
			b.Put([]byte("lastHashKey"), genesisBlock.Hash)

			tail = genesisBlock.Hash

		} else {
			//2.获取最后一个区块哈希值
			//填充给tail
			//返回bc实例
			tail = b.Get([]byte("lastHashKey"))
		}

		return nil

	})

	return &BlockChain{db, tail}

}
```
## gobTest.go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//gob是go语言内置的编码包
//它可以对任意数据类型进行编码和解码
//编码时，先要创建编码器，编码器进行编码
//解码时，先要创建解码器，解码器进行解码

type Person struct {
	Name string
	Age  uint64
}

func main() {
	Jim := Person{
		Name: "Jim",
		Age:  19,
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&Jim)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("编码后的数据%x\n", buffer.Bytes())

	//传输中

	//解码，将字节六转换成Person结构
	var p1 Person

	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	err = decoder.Decode(&p1)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("解码后的数据%v\n", p1)
}

## 编码解码区块
``` go

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

	fmt.Printf("解码传入的数据:%x\n", data)

	var block Block

	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
```

## 更新addblock
```go


//添加区块

func (bc *BlockChain) AddBlock(data string) {
	//创建一个区块
	bc.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("blockBucket"))

		if b == nil {
			//如果b为空，说明该桶不存在，需要创建
			fmt.Printf("bucket不存在，请检查!\n")
			os.Exit(1)
		}

		block := NewBlock(data, bc.tail)
		b.Put(block.Hash, block.Serialize() /*将区块序列化，转成字节流*/)
		b.Put([]byte("lastHashKey"), block.Hash)

		bc.tail = block.Hash

		return nil

	})

}
```
## 定义迭代器，创建迭代器
```go 
//定义一个区块链的迭代器，包含db，current
type BlockChainIterator struct {
	db      *bolt.DB
	current []byte
}

//创建迭代器，使用bc进行初始化

func (bc *BlockChain) NewIterator() *BlockChainIterator {

	return &BlockChainIterator{bc.db, bc.tail}

}
```

## Next实现
```go

func (it *BlockChainIterator) Next() *Block {

	var block Block

	it.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockBucketName))
		if b == nil {
			//如果b为空，说明该桶不存在，需要创建
			fmt.Printf("bucket不存在，请检查!\n")
			os.Exit(1)
		}

		//真正读取数据
		blockInfo /*block的字节流*/ := b.Get(it.current)
		block = *DeSerialize(blockInfo)

		it.current = block.PrevBlockHash

		return nil

	})
	return &block
}
```
## 使用迭代器，更新main函数
```go
package main

import (
	"bytes"
	"fmt"
	"time"
)

func main() {
	//fmt.Printf("hello world")
	//block := NewBlock(genesisInfo, []byte{0x0000000000000000})
	bc := NewBlockChain()
	defer bc.db.Close()
	bc.AddBlock("hello itcast!!!")

	it := bc.NewIterator()

	for {
		block := it.Next()
		fmt.Printf("+++++++++++++++++++++++++++++++++++++\n")

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

		if bytes.Equal(block.PrevBlockHash, []byte{}) {
			fmt.Printf("区块链遍历结束！\n")
			break
		}
	}

}
```
## 命令行demo
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	cmds := os.Args

	for i, cmd := range cmds {
		fmt.Printf("cmd[%d]:%s\n", i, cmd)
	}
}
```
## 定义CLI结构，run方法框架搭建
```go 

const Usage = `
	blockchain addBlock "xxxxx"	添加数据到区块链
	blockchain printChain		打印区块链
`

type CLI struct {
	bc *BlockChain
}

//给CLI提供一个方法，进行命令解析，从而执行调度
func (cli *CLI) Run() {
	cmds := os.Args

	if len(cmds) < 2 {
		fmt.Printf(Usage)
		os.Exit(1)
	}

	switch cmds[1] {
	case "addBlock":
		fmt.Printf("添加区块链命令被调用，数据：%s\n", cmds[2])
	case "printChain":
		fmt.Printf("打印区块链命令被调用\n")
	default:
		fmt.Printf("无效命令，请检查\n")
		fmt.Printf(Usage)

	}
}
```

## 改写主函数
```go
package main

func main() {
	bc := NewBlockChain()
	defer bc.db.Close()
	cli := CLI{bc}
	cli.Run()
}
```
## 更新run函数

``` go

//给CLI提供一个方法，进行命令解析，从而执行调度
func (cli *CLI) Run() {
	cmds := os.Args

	if len(cmds) < 2 {
		fmt.Printf(Usage)
		os.Exit(1)
	}

	switch cmds[1] {
	case "addBlock":
		fmt.Printf("添加区块链命令被调用，数据：%s\n", cmds[2])
		data := cmds[2]
		cli.AddBlock(data)
	case "printChain":
		fmt.Printf("打印区块链命令被调用\n")
		cli.PrintChain()
	default:
		fmt.Printf("无效命令，请检查\n")
		fmt.Printf(Usage)

	}
}

```
## 添加command.go
``` go

func (cli *CLI) PrintChain() {

	it := cli.bc.NewIterator()

	for {
		block := it.Next()
		fmt.Printf("+++++++++++++++++++++++++++++++++++++\n")

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

		if bytes.Equal(block.PrevBlockHash, []byte{}) {
			fmt.Printf("区块链遍历结束！\n")
			break
		}
	}

}
```

## 交易结构定义
``` go 
package main

type TXInput struct {
	TXID    []byte //交易id
	Index   int64  //ouput的索引
	Address string //解锁脚本，先用地址来模拟
}

type TXOutput struct {
	Value   float64 //转账金额
	Address string  //锁定脚本
}

type Transaction struct {
	Txid      []byte     //交易id
	TXInputs  []TXInput  //所有的inputs
	TXOutputs []TXOutput //所有的outputs
}

```

## settxid函数实现
```go

func (tx *Transaction) SetTXID() {

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(tx)

	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(buffer.Bytes())
	tx.Txid = hash[:]

}
```
## 挖矿交易实现
```go
//实现挖矿交易，
//特点：只有输出，没有有效的输入（不需要引用id，不需要索引，不需要签名）

//把挖矿的人传递进来，因为有奖励
func NewCoinbaseTx(miner string) *Transaction {
	//我们在后面的程序中，需要识别一个交易是否为coinbase，所以设置一些特殊值，用于判断
	//TODO
	inputs := []TXInput{TXInput{nil, -1, genesisInfo}}
	outputs := []TXOutput{TXOutput{12.5, miner}}

	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()
	
	return &tx
}
```
## 使用Transaction改写程序
1、改写block结构
2、根据提示修改：逐个文件处理
3、使用strings命令查看

## 模拟梅克尔根
```go
//比特币做hash，是对区块头做hash
//模拟梅克尔根，做一个简单的处理
func (block *Block) HashTransactions() {
	//交易的id就是交易的哈希值，所以可以将id拼接起来，整体做一个哈希运算，作为Merkleroot
	var hashes []byte
	for _, tx := range block.Transactions {
		txid /*[]byte*/ := tx.Txid
		hashes = append(hashes, txid...)
	}
	hash := sha256.Sum256(hashes)
	block.MerkleRoot = hash[:]

}
//在newblock中调用
```

## 查找余额
```go

func (bc *BlockChain) FindMyUtxos(address string) []TXOutput {
	//TODO
	fmt.Printf("FindMyUtxos\n")

	return []TXOutput{}
}

func (bc *BlockChain) GetBalance(address string) {
	utxos := bc.FindMyUtxos(address)
	var total = 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("%s的余额为：%f\n", address, total)
}
```
在cli中添加getbalance命令，调用getbalance函数

## 遍历交易输出
```go

func (bc *BlockChain) FindMyUtxos(address string) []TXOutput {
	//TODO
	fmt.Printf("FindMyUtxos\n")
	var UTXOs []TXOutput //返回的结构

	it := bc.NewIterator()

	//1.遍历账本
	for {
		block := it.Next()

		//2.遍历交易
		for _, tx := range block.Transactions {

			//3.遍历output
			for i, output := range tx.TXOutputs {

				//4.找到所有属于账户的output
				if address == output.Address {
					fmt.Printf("找到了属于%s的output，i:%d\n", address, i)
					UTXOs = append(UTXOs, output)
				}

			}
		}

		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("遍历区块链结束!\n")
			break
		}
	}

	return UTXOs
}
```
## 遍历交易的inputs
```go

	//1.遍历账本
	for {
		block := it.Next()

		//2.遍历交易
		for _, tx := range block.Transactions {
			//遍历交易输入：inputs
			for _, input := range tx.TXInputs {
				if input.Address == address {
					fmt.Printf("找到了消耗过的output！index:%d\n", input.Index)
					key := string(input.TXID)
					spentUTXOs[key] = append(spentUTXOs[key], input.Index)

				}
			}

			//3.遍历output
			for i, output := range tx.TXOutputs {

				//4.找到所有属于账户的output
				if address == output.Address {
					fmt.Printf("找到了属于%s的output，i:%d\n", address, i)
					UTXOs = append(UTXOs, output)
				}

			}
		}

		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("遍历区块链结束!\n")
			break
		}
	}
```

## 整体遍历过程
```go

func (bc *BlockChain) FindMyUtxos(address string) []TXOutput {
	//TODO
	fmt.Printf("FindMyUtxos\n")
	var UTXOs []TXOutput //返回的结构

	it := bc.NewIterator()

	//这是标识已经消耗过的utxo结构，key是交易id，value是这个id里面的output索引的数组
	spentUTXOs := make(map[string][]int64)

	//1.遍历账本
	for {
		block := it.Next()

		//2.遍历交易
		for _, tx := range block.Transactions {
			//遍历交易输入：inputs
			for _, input := range tx.TXInputs {
				if input.Address == address {
					fmt.Printf("找到了消耗过的output！index:%d\n", input.Index)
					key := string(input.TXID)
					spentUTXOs[key] = append(spentUTXOs[key], input.Index)

				}
			}
		OUTPUT:
			//3.遍历output
			for i, output := range tx.TXOutputs {
				key := string(tx.Txid)
				indexs := spentUTXOs[key]
				if len(indexs) != 0 {
					fmt.Printf("当前这笔交易中有被消耗过的output！\n")
					for _, j := range indexs {
						if int64(i) == j {
							fmt.Printf("i=j,当前的output已经被消耗过了，跳过不统计")
							continue OUTPUT
						}

					}

				}

				//4.找到所有属于账户的output
				if address == output.Address {
					fmt.Printf("找到了属于%s的output，i:%d\n", address, i)
					UTXOs = append(UTXOs, output)
				}

			}
		}

		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("遍历区块链结束!\n")
			break
		}
	}

	return UTXOs
}

func (bc *BlockChain) GetBalance(address string) {
	utxos := bc.FindMyUtxos(address)
	var total = 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("%s的余额为：%f\n", address, total)
}
```

## 创建普通交易
参数：
1、付款人
2、收款人
3、转账金额
4、bc

内部逻辑：
遍历账本，找到属于付款人的合适的金额，把这个outputs找到
如果找到钱不足以转账，创建交易失败
将outputs转成inputs
创建输出，创建一个属于收款人的output
如果有找零，创建属于付款人output
设置交易id
返回交易结构

### 代码
``` go

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {

	utxos := make(map[string][]int64)
	var resValue float64
	//遍历账本，找到属于付款人的合适的金额，把这个outputs找到
	utxos, resValue = bc.FindNeedUtxos(from, amount)

	//如果找到钱不足以转账，创建交易失败
	if resValue < amount {
		fmt.Printf("余额不足，交易失败!\n")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	//将outputs转成inputs
	for txid, indexes := range utxos {
		for _, i := range indexes {
			input := TXInput{[]byte(txid), i, from}
			inputs = append(inputs, input)

		}
	}
	//创建输出，创建一个属于收款人的output
	output := TXOutput{amount, to}
	outputs = append(outputs, output)

	//如果有找零，创建属于付款人output
	if resValue > amount {
		output1 := TXOutput{resValue - amount, from}
		outputs = append(outputs, output1)
	}
	//创建交易
	tx := Transaction{nil, inputs, outputs}
	//设置交易id
	tx.SetTXID()
	//返回交易结构
	return &tx

}
```
#  FindNeedUtxos函数
```go
func (bc *BlockChain) FindNeedUtxos(from string, amount float64) (map[string][]int64, float64) {
	//TODO 找到合理utxos集合
	utxos := make(map[string][]int64)
	return utxos, 0.0

}
```

返回交易结构

## 添加send命令

```go
case "send":
		fmt.Printf("转账命令被调用\n")
		//./blockchain send FROM TO AMOUNT MINER DATA	转账
		if len(cmds) != 7 {
			fmt.Printf("send命令发现无效参数，请检查\n")
			fmt.Printf(Usage)
			os.Exit(1)
		}
		from := cmds[2]
		to := cmds[3]
		amount, _ := strconv.ParseFloat(cmds[4], 64)
		miner := cmds[5]
		data := cmds[6]
		cli.Send(from, to, amount, miner, data)
```

## Send命令实现
``` go

func (cli *CLI) Send(from string, to string, amount float64, miner string, data string) {
	//创建挖矿交易
	//创建普通交易
	//添加到区块

	//1.创建挖矿者
	coinbase := NewCoinbaseTx(miner, data)

	//2.创建普通交易
	tx := NewTransaction(from, to, amount, cli.bc)

	txs := []*Transaction{coinbase}

	if tx != nil {
		txs = append(txs, tx)
	} else {
		fmt.Printf("发现无效交易，过滤!\n")
	}

	//3.添加到区块
	cli.bc.AddBlock(txs)

	fmt.Printf("挖矿成功!\n")
}
```
## 定义UTXOInfo
```go
type UTXOInfo struct {
	TXID   []byte   //交易id
	Index  int64    //output的索引值
	Output TXOutput //output本身
}
```
## 改写FindMyUtxos
```go

func (bc *BlockChain) FindMyUtxos(address string) []UTXOInfo {

	fmt.Printf("FindMyUtxos\n")
	//var UTXOs []TXOutput //
	var UTXOInfos []UTXOInfo //新的返回结构

	it := bc.NewIterator()

	//这是标识已经消耗过的utxo结构，key是交易id，value是这个id里面的output索引的数组
	spentUTXOs := make(map[string][]int64)

	//1.遍历账本
	for {
		block := it.Next()

		//2.遍历交易
		for _, tx := range block.Transactions {
			//遍历交易输入：inputs
			for _, input := range tx.TXInputs {
				if input.Address == address {
					fmt.Printf("找到了消耗过的output！index:%d\n", input.Index)
					key := string(input.TXID)
					spentUTXOs[key] = append(spentUTXOs[key], input.Index)

				}
			}
		OUTPUT:
			//3.遍历output
			for i, output := range tx.TXOutputs {
				key := string(tx.Txid)
				indexs := spentUTXOs[key]
				if len(indexs) != 0 {
					fmt.Printf("当前这笔交易中有被消耗过的output！\n")
					for _, j := range indexs {
						if int64(i) == j {
							fmt.Printf("i=j,当前的output已经被消耗过了，跳过不统计")
							continue OUTPUT
						}

					}

				}

				//4.找到所有属于账户的output
				if address == output.Address {
					fmt.Printf("找到了属于%s的output，i:%d\n", address, i)
					//UTXOs = append(UTXOs, output)
					utxoinfo := UTXOInfo{tx.Txid, int64(i), output}
					UTXOInfos = append(UTXOInfos, utxoinfo)
				}

			}
		}

		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("遍历区块链结束!\n")
			break
		}
	}

	return UTXOInfos
}
```

## 改写FindNeedUtxos
```go

func (bc *BlockChain) FindNeedUtxos(from string, amount float64) (map[string][]int64, float64) {

	needUtxos := make(map[string][]int64)
	var resValue float64 //统计的金额

	//复用Findmyutxo函数，这个函数已经包含所有信息
	utxoinfos := bc.FindMyUtxos(from)
	for _, utxoinfo := range utxoinfos {
		key := string(utxoinfo.TXID)
		needUtxos[key] = append(needUtxos[key], int64(utxoinfo.Index))
		resValue += utxoinfo.Output.Value

		//判断金额是否足够
		if resValue >= amount {
			break
		}

	}
	return needUtxos, resValue
}
```
## IsCoinbase实现
```go

func (tx *Transaction) IsCoinbase() bool {
	//特点：1、只有一个input 2、引用的id是nil 3、引用的索引是-1
	inputs := tx.TXInputs
	if len(inputs) == 1 && inputs[0].TXID == nil && inputs[0].Index == -1 {
		return true
	}

	return false

}
```
## 在遍历inputs时使用
``` go
//2.遍历交易
		for _, tx := range block.Transactions {
			//遍历交易输入：inputs

			if tx.IsCoinbase() == false {
				//如果不是coinbase，说明是普通交易，才有必要进行遍历
				for _, input := range tx.TXInputs {
					if input.Address == address {
						fmt.Printf("找到了消耗过的output！index:%d\n", input.Index)
						key := string(input.TXID)
						spentUTXOs[key] = append(spentUTXOs[key], input.Index)

					}
				}
			}
			
```
## 创建区块链函数 
```go

func CreatBlockChain(miner string) *BlockChain {

	//功能分析
	//1.获得数据库句柄，打开数据库，读写数据

	db, err := bolt.Open(blockChainName, 0600, nil)
	//向数据库中写入数据
	//从数据库中读取数据

	if err != nil {
		log.Panic(err)
	}

	//defer db.Close()

	var tail []byte

	db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockBucketName))

		b, err := tx.CreateBucket([]byte(blockBucketName))

		if err != nil {
			log.Panic(err)
		}

		//抽屉准备完毕，开始添加创世块
		//创世块中只有一个挖矿交易
		coinbase := NewCoinbaseTx(miner, genesisInfo)

		genesisBlock := NewBlock([]*Transaction{coinbase}, []byte{})
		b.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化，转成字节流*/)
		b.Put([]byte(lastHashKey), genesisBlock.Hash)

		//为了测试，把写入的数据读取出来，如果没问题，注释掉
		//blockInfo := b.Get(genesisBlock.Hash)
		//block := DeSerialize(blockInfo)
		//fmt.Printf("解码后的blcok数据:%s\n", block)

		tail = genesisBlock.Hash

		return nil

	})

	return &BlockChain{db, tail}

}
```
## 获取实例
```go

//返回区块链实例
func NewBlockChain() *BlockChain {

	//功能分析
	//1.获得数据库句柄，打开数据库，读写数据

	db, err := bolt.Open(blockChainName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	//defer db.Close()
	var tail []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))
		if b == nil {
			fmt.Printf("区块链bucket为空，请检查!\n")
			os.Exit(1)
		}
		tail = b.Get([]byte(lastHashKey))
		return nil
	})
	return &BlockChain{db, tail}

}
```

## 添加CreatBlockChain命令
```go
func (cli *CLI) CreatBlockChain(addr string) {
	bc := CreatBlockChain(addr)
	bc.db.Close()
	fmt.Printf("创建区块链成功!\n")

}
```
## 主函数
```go
func main() {
	//bc := NewBlockChain("cg")
	//defer bc.db.Close()
	//cli := CLI{bc}
	cli := CLI{}
	cli.Run()
}
```


## 在程序中调用NewBlockChain，例如
```go

func (cli *CLI) GetBalance(addr string) {
	bc := NewBlockChain()
	defer bc.db.Close()
	bc.GetBalance(addr)
}
```
## 判断文件是否存在
```go

//判断文件是否存在
func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
```
