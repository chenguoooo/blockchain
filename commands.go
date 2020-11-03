package main

import (
	"bytes"
	"fmt"
	"time"
)

//实现具体的命令
func (cli *CLI) CreatBlockChain(addr string) {
	if !IsValidAddress(addr) {
		fmt.Printf("%s是无效地址！\n", addr)
		return
	}

	bc := CreatBlockChain(addr)
	if bc != nil {
		defer bc.db.Close()
	}
	fmt.Printf("创建区块链成功!\n")

}

func (cli *CLI) GetBalance(addr string) {
	if !IsValidAddress(addr) {
		fmt.Printf("%s是无效地址！\n", addr)
		return
	}

	bc := NewBlockChain()
	if bc == nil {
		return
	}

	bc.GetBalance(addr)
}

func (cli *CLI) PrintChain() {
	bc := NewBlockChain()
	if bc == nil {
		return
	}

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
		fmt.Printf("Data:%s\n", block.Transactions[0].TXInputs[0].Pubkey)

		pow := NewProofOfWork(block)
		fmt.Printf("Isvalid:%v\n", pow.IsValid())

		if bytes.Equal(block.PrevBlockHash, []byte{}) {
			fmt.Printf("区块链遍历结束！\n")
			break
		}
	}

}

func (cli *CLI) Send(from string, to string, amount float64, miner string, data string) {
	if !IsValidAddress(from) {
		fmt.Printf("from:%s是无效地址！\n", from)
		return
	}
	if !IsValidAddress(to) {
		fmt.Printf("to:%s是无效地址！\n", to)
		return
	}
	if !IsValidAddress(miner) {
		fmt.Printf("miner:%s是无效地址！\n", miner)
		return
	}

	bc := NewBlockChain()
	if bc == nil {
		return
	}
	//1.创建挖矿交易
	coinbase := NewCoinbaseTx(miner, data)

	//2.创建普通交易
	tx := NewTransaction(from, to, amount, bc)

	txs := []*Transaction{coinbase}

	if tx != nil {
		txs = append(txs, tx)
	} else {
		fmt.Printf("发现无效交易，过滤!\n")
	}

	//3.添加到区块
	bc.AddBlock(txs)

	fmt.Printf("挖矿成功!\n")
}

func (cli *CLI) CreateWallet() {

	ws := NewWallets()
	address := ws.CreateWallet()

	fmt.Printf("新的钱包地址为：%s\n", address)
}

func (cli *CLI) ListAddresses() {
	ws := NewWallets()

	addresses := ws.ListAddress()
	for _, address := range addresses {
		fmt.Printf("address : %s\n", address)
	}
}

func (cli *CLI) PrintTx() {
	bc := NewBlockChain()
	if bc == nil {
		return
	}

	it := bc.NewIterator()

	for {
		block := it.Next()
		fmt.Printf("\n+++++++++++++++新的区块+++++++++++++++++++\n")
		for _, tx := range block.Transactions {
			fmt.Printf("tx:%v\n", tx)
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}

	}

}
