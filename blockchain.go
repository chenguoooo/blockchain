package main

import (
	"base58"
	"bolt"
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"sync"
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

func CreatBlockChain(miner string) *BlockChain {

	if IsFileExist(blockChainName) {
		fmt.Printf("区块链已经存在，不需要重复创建！\n")
		return nil
	}

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

//返回区块链实例
func NewBlockChain() *BlockChain {
	if !IsFileExist(blockChainName) {
		fmt.Printf("区块链不存在，请先创建！\n")
		return nil
	}

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

//添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) int64 {
	//矿工得到交易时，第一时间对交易进行验证
	//矿工如果不验证，即使挖矿成功，广播区块后，其他的验证矿工，仍然会校验每一笔交易
	validTXs := []*Transaction{}
	for _, tx := range txs {
		if bc.VerifyTransaction(tx) && bc.ValidTransaction(validTXs, tx) {
			//fmt.Printf("---该交易有效：%x\n", tx.Txid)
			validTXs = append(validTXs, tx)
		}
	}
	validTXslen := len(validTXs)
	fmt.Printf("txlen:%d\n", validTXslen)
	//创建一个区块
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))
		if b == nil {
			//如果b为空，说明该桶不存在，需要创建
			fmt.Printf("bucket不存在，请检查!\n")
			os.Exit(1)
		}
		block := NewBlock(validTXs, bc.tail)
		b.Put(block.Hash, block.Serialize() /*将区块序列化，转成字节流*/)
		b.Put([]byte(lastHashKey), block.Hash)
		bc.tail = block.Hash
		return nil
	})
	return int64(validTXslen)
}

//添加区块

func (bc *BlockChain) AddBlockParal(txs []*Transaction) int64 {
	//矿工得到交易时，第一时间对交易进行验证
	//矿工如果不验证，即使挖矿成功，广播区块后，其他的验证矿工，仍然会校验每一笔交易

	validTXs := []*Transaction{}
	txslen := len(txs)
	t := make(chan *Transaction, txslen)
	mutex := make(chan int, txslen)
	for _, tx := range txs {
		go func(tx1 *Transaction) {
			if bc.VerifyTransaction(tx1) {
				t <- tx1
			}
			mutex <- 1
		}(tx)
		//fmt.Printf("\n运行\n")
	}
	for i := 0; i < len(txs); i++ {
		<-mutex
	}
	tlen := len(t)
	for i := 0; i < tlen; i++ {
		tx := <-t
		if bc.ValidTransaction(validTXs, tx) {
			validTXs = append(validTXs, tx)
		}
	}
	validtxslen := len(validTXs)
	fmt.Printf("txlen:%d\n", validtxslen)
	//fmt.Printf("\n%d\n", len(validTXs))
	//创建一个区块
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))
		if b == nil {
			//如果b为空，说明该桶不存在，需要创建
			fmt.Printf("bucket不存在，请检查!\n")
			os.Exit(1)
		}
		block := NewBlock(validTXs, bc.tail)
		b.Put(block.Hash, block.Serialize() /*将区块序列化，转成字节流*/)
		b.Put([]byte(lastHashKey), block.Hash)
		bc.tail = block.Hash
		return nil
	})
	return int64(validtxslen)

}

//定义一个区块链的迭代器，包含db，current
type BlockChainIterator struct {
	db      *bolt.DB
	current []byte //当前所指向区块的哈希值
}

//创建迭代器，使用bc进行初始化

func (bc *BlockChain) NewIterator() *BlockChainIterator {

	return &BlockChainIterator{bc.db, bc.tail}

}

var mutexView sync.Mutex

//取出当前区块信息，并将指针指向prevblock
func (it *BlockChainIterator) Next() *Block {

	var block Block
	mutexView.Lock()
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
	mutexView.Unlock()

	return &block
}

type UTXOInfo struct {
	TXID   []byte   //交易id
	Index  int64    //output的索引值
	Output TXOutput //output本身
}

func (bc *BlockChain) FindMyUtxos(pubKeyHash []byte) []UTXOInfo {

	//fmt.Printf("FindMyUtxos\n")
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
			if tx.IsCoinbase() == false {
				//如果不是coinbase，说明是普通交易，才有必要进行遍历
				for _, input := range tx.TXInputs {
					//判断当前使用的input是否为目标地址所有
					if bytes.Equal(HashPubKey(input.Pubkey), pubKeyHash) {
						//fmt.Printf("找到了消耗过的output！index:%d\n", input.Index)
						key := string(input.TXID)
						spentUTXOs[key] = append(spentUTXOs[key], input.Index)
					}
				}
			}
		OUTPUT:
			//3.遍历output
			for i, output := range tx.TXOutputs {
				key := string(tx.Txid)
				indexs := spentUTXOs[key]
				if len(indexs) != 0 {
					//fmt.Printf("当前这笔交易中有被消耗过的output！\n")
					for _, j := range indexs {
						if int64(i) == j {
							//fmt.Printf("i=j,当前的output已经被消耗过了，跳过不统计")
							continue OUTPUT
						}
					}
				}
				//4.找到所有属于账户的output
				if bytes.Equal(pubKeyHash, output.PubKeyHash) {
					//fmt.Printf("找到了属于%s的output，i:%d\n", address, i)
					//UTXOs = append(UTXOs, output)
					utxoinfo := UTXOInfo{tx.Txid, int64(i), output}
					UTXOInfos = append(UTXOInfos, utxoinfo)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			//fmt.Printf("遍历区块链结束!\n")
			break
		}
	}

	return UTXOInfos
}

func (bc *BlockChain) GetBalance(address string) (total float64) {
	decodeInfo, _ := base58.Decode(address)
	pubKeyHash := decodeInfo[1 : len(decodeInfo)-4]

	utxoinfos := bc.FindMyUtxos(pubKeyHash)

	total = 0.0
	//所有的output都在utxoinfos内部
	//获取余额时，遍历utxoinfos获取output即可

	for _, utxo := range utxoinfos {
		total += utxo.Output.Value
	}
	//fmt.Printf("%s的余额为：%f\n", address, total)
	return total

}

//遍历账本，找到属于付款人的合适的金额，把这个outputs找到
//utxos, resValue = bc.FindNeedUtxos(from, amount)

func (bc *BlockChain) FindNeedUtxos(pubKeyHash []byte, amount float64) (map[string][]int64, float64) {

	needUtxos := make(map[string][]int64)
	var resValue float64 //统计的金额

	//复用Findmyutxo函数，这个函数已经包含所有信息
	utxoinfos := bc.FindMyUtxos(pubKeyHash)
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

func (bc *BlockChain) SignTranscation(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	//1.遍历账本找到所有引用交易
	prevTXs := make(map[string]Transaction)

	//遍历tx的inputs，通过id去查找所引用的交易
	for _, input := range tx.TXInputs {
		prevTx := bc.FindTransaction(input.TXID)

		if prevTx == nil {
			fmt.Printf("没有找到交易：%x\n", input.TXID)
		} else {
			//把找到的引用交易保存起来
			prevTXs[string(input.TXID)] = *prevTx
		}

	}

	tx.Sign(privateKey, prevTXs)

}

//矿工校验过程
//1.找到交易input所引用的交易prevTXs
//2.对交易进行校验

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	//校验的时候，如果是挖矿交易，直接返回true
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]Transaction)

	//遍历tx的inputs，通过id去查找所引用的交易
	for _, input := range tx.TXInputs {
		prevTx := bc.FindTransaction(input.TXID)

		if prevTx == nil {
			//fmt.Printf("没有找到交易：%x\n", input.TXID)
		} else {
			//把找到的引用交易保存起来
			prevTXs[string(input.TXID)] = *prevTx

		}

	}

	return tx.Verify(prevTXs)

}
func (bc *BlockChain) ValidTransaction(validTXs []*Transaction, tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	for _, tx1 := range validTXs {
		for _, txi1 := range tx1.TXInputs {
			for _, txi := range tx.TXInputs {
				if (bytes.Equal(txi1.TXID, txi.TXID)) && (txi1.Index == txi.Index) {
					return false
				}
			}
		}
	}
	return true

}

func (bc *BlockChain) FindTransaction(txid []byte) *Transaction {
	//遍历区块链的交易
	//通过对比id来识别

	it := bc.NewIterator()

	for {
		block := it.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.Txid, txid) {
				//fmt.Printf("找到了所引用的交易，id:%x\n", tx.Txid)
				return tx
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return nil
}
