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
	//fmt.Printf("创建区块链成功!\n")

}

func (cli *CLI) GetBalance(addr string) (total float64) {
	if !IsValidAddress(addr) {
		fmt.Printf("%s是无效地址！\n", addr)
		return
	}

	bc := NewBlockChain()
	if bc == nil {
		return
	}
	if bc != nil {
		defer bc.db.Close()
	}

	total = bc.GetBalance(addr)
	//fmt.Printf("%s的余额为：%f\n", addr, total)
	return
}

func (cli *CLI) PrintChain() {
	bc := NewBlockChain()
	if bc == nil {
		return
	}
	if bc != nil {
		defer bc.db.Close()
	}
	fmt.Printf("\n打印区块链！\n")

	it := bc.NewIterator()
	i := 0
	for {
		block := it.Next()
		fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")

		fmt.Printf("Version:%d\n", block.Version)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("MerkleRoot:%x\n", block.MerkleRoot)

		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp:%s\n", timeFormat)

		fmt.Printf("Difficulity:%d\n", block.Difficulity)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%s\n", string(block.Transactions[0].TXInputs[0].Pubkey))

		pow := NewProofOfWork(block)
		fmt.Printf("Isvalid:%v\n", pow.IsValid())
		i++
		fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
		if bytes.Equal(block.PrevBlockHash, []byte{}) {
			//fmt.Printf("区块链遍历结束！\n")

			break
		}

	}
	return

}

//
//func (cli *CLI) Send(from string, to string, amount float64) (tx *Transaction) {
//	if !IsValidAddress(from) {
//		fmt.Printf("from:%s是无效地址！\n", from)
//		return
//	}
//	if !IsValidAddress(to) {
//		fmt.Printf("to:%s是无效地址！\n", to)
//		return
//	}
//
//	bc := NewBlockChain()
//	if bc == nil {
//		return
//	}
//	if bc != nil {
//		defer bc.db.Close()
//	}
//	//创建普通交易
//	//tx = NewTransaction(from, to, amount, bc)
//	if tx != nil {
//		//fmt.Printf("\n创建交易成功!\n\n")
//		return tx
//	} else {
//		//fmt.Printf("发现无效交易，过滤!\n")
//	}
//	return
//}

func (cli *CLI) AddBlock(miner string, data string, txs []*Transaction) {
	if !IsValidAddress(miner) {
		fmt.Printf("miner:%s是无效地址！\n", miner)
		return
	}
	bc := NewBlockChain()
	if bc == nil {
		return
	}
	if bc != nil {
		defer bc.db.Close()
	}

	coinbase := NewCoinbaseTx(miner, data)
	txs1 := []*Transaction{}
	txs1 = append(txs1, coinbase)
	txs1 = append(txs1, txs...)

	bc.AddBlock(txs1)
	//fmt.Printf("挖矿成功!\n")
}

func (cli *CLI) CreateWallet() {

	ws := NewWallets()
	ws.CreateWallet()
	//address := ws.CreateWallet()
	//fmt.Printf("新的钱包地址为：%s\n", address)
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
	if bc != nil {
		defer bc.db.Close()
	}
	fmt.Printf("\n打印交易:")
	it := bc.NewIterator()

	for {
		block := it.Next()

		fmt.Printf("\n++++++++++++++++++++++++++++++++++新的区块++++++++++++++++++++++++++++++++++++++\n")

		for i, tx := range block.Transactions {
			fmt.Printf("--------------------tx%d----------------------%v\n", i, tx)
		}
		fmt.Printf("\n++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}

	}

}
