package main

import (
	"fmt"
	"os"
)

const Usage = `
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	./blockchain CreatBlockChain 地址		"创建区块链"
	./blockchain printChain				"打印区块链"
	./blockchain getBalance	地址			"获取地址的余额"
	./blockchain send FROM TO AMOUNT 		"转账命令"
	./blockchain createWallet			"创建钱包"
	./blockchain ListAddresses			"打印钱包地址"
	./blockchain printTx				"打印交易"
	./blockchain addBlock Miner Data		"添加区块"
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
`

type CLI struct {
	//bc *BlockChain
	//CLI中不需要保存区块链实例了，所有名字在自己调用前，自己获取区块链实例
}

//给CLI提供一个方法，进行命令解析，从而执行调度
func (cli *CLI) Run() {

	cmds := os.Args
	if len(cmds) < 2 {
		fmt.Printf(Usage)
		os.Exit(1)
	}

	switch cmds[1] {
	case "h":
		fmt.Printf(Usage)

	case "CreatBlockChain":
		if len(cmds) != 3 {
			fmt.Printf(Usage)
			os.Exit(1)
		}
		fmt.Printf("创建区块链命令被调用\n")
		addr := cmds[2]
		cli.CreatBlockChain(addr)
		fmt.Printf(Usage)

	case "printChain":
		fmt.Printf("打印区块链命令被调用\n")
		cli.PrintChain()
		fmt.Printf(Usage)

	case "getBalance":
		fmt.Printf("获取余额命令被调用\n")
		total := cli.GetBalance(cmds[2])
		fmt.Printf("%s的余额为：%f\n", cmds[2], total)
		fmt.Printf(Usage)

	//case "send":
	//	fmt.Printf("转账命令被调用\n")
	//	//./blockchain send FROM TO AMOUNT 	转账
	//	if len(cmds) != 5 {
	//		fmt.Printf("send命令发现无效参数，请检查\n")
	//		fmt.Printf(Usage)
	//		os.Exit(1)
	//	}
	//from := cmds[2]
	//to := cmds[3]
	//amount, _ := strconv.ParseFloat(cmds[4], 64)
	//tx := cli.Send(bc,from, to, amount)
	//	txs := NewTransactions()
	//	txs.CreateTransaction(tx)
	//	fmt.Printf(Usage)
	//TODO

	case "createWallet":
		fmt.Printf("创建钱包命令被调用\n")
		cli.CreateWallet()
		fmt.Printf(Usage)
	case "ListAddresses":
		fmt.Printf("打印钱包地址命令被调用\n")
		cli.ListAddresses()
		fmt.Printf(Usage)
	case "printTx":
		fmt.Printf("打印交易命令被调用\n")
		cli.PrintTx()
		fmt.Printf(Usage)
	case "addBlock":
		//func (cli *CLI) AddBlock(miner string, data string,txs []*Transaction){
		miner := cmds[2]
		data := cmds[3]
		fmt.Printf("挖矿交易命令被调用\n")
		txs := NewTransactions()
		txs1 := []*Transaction{}
		for _, tx := range txs.TransactionsMap {
			txs1 = append(txs1, tx)
		}
		cli.AddBlock(miner, data, txs1) //要改
		txs.ClearFile()
		fmt.Printf(Usage)

	default:
		fmt.Printf("无效命令，请检查\n")
		fmt.Printf(Usage)

	}
}
