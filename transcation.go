package main

import (
	"base58"
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type TXInput struct {
	TXID  []byte //交易id
	Index int64  //ouput的索引
	//Address string //解锁脚本，先用地址来模拟
	Signature []byte //交易签名
	Pubkey    []byte //公钥本身，不是公钥哈希
}

type TXOutput struct {
	Value float64 //转账金额
	//Address string  //锁定脚本
	PubKeyHash []byte //公钥哈希，不是公钥本身
}

//给定转账地址，得到这个地址的公钥哈希，完成对output的锁定
func (output *TXOutput) Lock(address string) {
	//address ->public key hash
	decodeInfo, _ := base58.Decode(address)
	pubKeyHash := decodeInfo[1 : len(decodeInfo)-4]
	output.PubKeyHash = pubKeyHash
}

func NewTXOutput(value float64, address string) TXOutput {
	output := TXOutput{Value: value}
	output.Lock(address)
	return output
}

type Transaction struct {
	Txid      []byte     //交易id
	TXInputs  []TXInput  //所有的inputs
	TXOutputs []TXOutput //所有的outputs
}

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

//实现挖矿交易，
//特点：只有输出，没有有效的输入（不需要引用id，不需要索引，不需要签名）

//挖矿奖励
const reward = 12.5

//把挖矿的人传递进来，因为有奖励
func NewCoinbaseTx(miner string, data string) *Transaction {

	inputs := []TXInput{TXInput{nil, -1, nil, []byte(data)}}
	//outputs := []TXOutput{TXOutput{12.5, miner}}

	output := NewTXOutput(reward, miner)
	outputs := []TXOutput{output}

	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	//特点：1、只有一个input 2、引用的id是nil 3、引用的索引是-1
	inputs := tx.TXInputs
	if len(inputs) == 1 && inputs[0].TXID == nil && inputs[0].Index == -1 {
		return true
	}

	return false

}

//内部逻辑：

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//1.打开钱包
	ws := NewWallets()

	wallet := ws.WalletsMap[from]

	if wallet == nil {
		fmt.Printf("%s的私钥不存在，交易创建失败!\n")
		return nil
	}

	//2.获取公钥私钥
	privateKey := wallet.PrivateKey //目前用不到，步骤三签名时使用
	publicKey := wallet.PublicKey

	publicKeyHash := HashPubKey(publicKey)

	utxos := make(map[string][]int64)
	var resValue float64
	//遍历账本，找到属于付款人的合适的金额，把这个outputs找到
	utxos, resValue = bc.FindNeedUtxos(publicKeyHash, amount)

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
			input := TXInput{[]byte(txid), i, nil, publicKey}
			inputs = append(inputs, input)

		}
	}
	//创建输出，创建一个属于收款人的output
	//output := TXOutput{amount, to}
	output := NewTXOutput(amount, to)
	outputs = append(outputs, output)

	//如果有找零，创建属于付款人output
	if resValue > amount {
		//output1 := TXOutput{resValue - amount, from}
		output1 := NewTXOutput(resValue-amount, from)
		outputs = append(outputs, output1)
	}
	//创建交易
	tx := Transaction{nil, inputs, outputs}
	//设置交易id
	tx.SetTXID()
	//把查找引用交易的环节放到Blockchain中去，同时在BlockChain进行调用签名

	//付款人在创建交易时，已经得到了所有引用的output的详细信息
	//但是不使用，因为矿工校验的时候，矿工没有这部分信息，矿工需要遍历账本找到所有引用交易
	//为了统一操作，所以再次查询，进行签名

	bc.SignTranscation(&tx, privateKey)

	//返回交易结构
	return &tx

}

//第一个参数是私钥
//第二个参数是这个交易的input所引用的所有交易
func (tx *Transaction) Sign(privkey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	fmt.Printf("对交易进行签名...\n")

	//TODO
}
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	fmt.Printf("对交易进行校验...\n")

	//TODO
	return true

}
