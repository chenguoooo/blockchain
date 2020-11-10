package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

type Transactions struct {
	TransactionsMap map[string]*Transaction
}

func NewTransactions() *Transactions {
	var txs Transactions

	txs.TransactionsMap = make(map[string]*Transaction)
	if !txs.LoadFromFile() {
		fmt.Printf("加载交易失败！\n")
	}
	return &txs

}

func (txs *Transactions) CreateTransaction(tx *Transaction) {
	if tx == nil {
		return
	}
	txs.TransactionsMap[string(tx.Txid)] = tx
	res := txs.SaveToFile()
	if !res {
		fmt.Printf("创建钱包失败!\n")
		return
	}

	return
}

const TransactionName = "transaction.dat"

func (txs *Transactions) SaveToFile() bool {

	var buffer bytes.Buffer
	//将接口类型明确注册一下，否则gob编码失败
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(txs)
	if err != nil {
		fmt.Printf("交易序列化失败,err:%v\n", err)
	}
	content := buffer.Bytes()

	//func WriteFile(filename string, data []byte, perm os.FileMode) error {
	err = ioutil.WriteFile(TransactionName, content, 0600)
	if err != nil {
		fmt.Printf("交易创建失败！\n")
		return false
	}

	return true
}

func (txs *Transactions) LoadFromFile() bool {
	//判断文件是否存在
	if !IsFileExist(TransactionName) {
		//fmt.Printf("交易文件不存在，准备创建！\n")
		return true
	}
	//读取文件
	content, err := ioutil.ReadFile(TransactionName)
	if err != nil {
		return false
	}
	gob.Register(elliptic.P256())
	//gob解码
	decoder := gob.NewDecoder(bytes.NewReader(content))

	var transactions Transactions
	err = decoder.Decode(&transactions)

	if err != nil {
		fmt.Printf("err:%v\n", err)
		return false
	}
	//赋值给ws
	txs.TransactionsMap = transactions.TransactionsMap
	return true

}
func (txs *Transactions) ClearFile() {
	err := os.Remove(TransactionName)
	if err != nil {
		fmt.Printf("清空文件出错")
	}
}
