package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

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

//把挖矿的人传递进来，因为有奖励
func NewCoinbaseTx(miner string, data string) *Transaction {
	//我们在后面的程序中，需要识别一个交易是否为coinbase，所以设置一些特殊值，用于判断
	//TODO
	inputs := []TXInput{TXInput{nil, -1, data}}
	outputs := []TXOutput{TXOutput{12.5, miner}}

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
