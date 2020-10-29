package main

import (
	"bolt"
	"fmt"
	"log"
	"os"
)

////创建区块链，使用Block数组模拟
//type BlockChain struct {
//	Blocks []*Block
//}

//使用blot改写
type BlockChain struct {
	db *bolt.DB //句柄

	tail []byte //最后一个区块hash值
}

const blockChainName = "blockChain.db"
const blockBucketName = "blockBucket"
const lastHashKey = "lastHashKey"

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

		if b == nil {
			//如果b为空，说明该桶不存在，需要创建
			fmt.Printf("bucket不存在，准备创建!\n")
			b, err = tx.CreateBucket([]byte(blockBucketName))

			if err != nil {
				log.Panic(err)
			}

			//抽屉准备完毕，开始添加创世块
			genesisBlock := NewBlock(genesisInfo, []byte{})
			b.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化，转成字节流*/)
			b.Put([]byte(lastHashKey), genesisBlock.Hash)

			//为了测试，把写入的数据读取出来，如果没问题，注释掉
			//blockInfo := b.Get(genesisBlock.Hash)
			//block := DeSerialize(blockInfo)
			//fmt.Printf("解码后的blcok数据:%s\n", block)

			tail = genesisBlock.Hash

		} else {
			//2.获取最后一个区块哈希值
			//填充给tail
			//返回bc实例
			tail = b.Get([]byte(lastHashKey))
		}

		return nil

	})

	return &BlockChain{db, tail}

}

//添加区块

func (bc *BlockChain) AddBlock(data string) {
	//创建一个区块
	bc.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockBucketName))

		if b == nil {
			//如果b为空，说明该桶不存在，需要创建
			fmt.Printf("bucket不存在，请检查!\n")
			os.Exit(1)
		}

		block := NewBlock(data, bc.tail)
		b.Put(block.Hash, block.Serialize() /*将区块序列化，转成字节流*/)
		b.Put([]byte(lastHashKey), block.Hash)

		bc.tail = block.Hash

		return nil

	})

}
