package main

import (
	"fmt"
	"os"
	"strconv"
)

const Usage = `
	./blockchain CreatBlockChain 地址		"创建区块链"
	./blockchain printChain				"打印区块链"
	./blockchain getBalance	地址			"获取地址的余额"
	./blockchain send FROM TO AMOUNT MINER DATA	"转账命令"
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
	case "CreatBlockChain":
		if len(cmds) != 3 {
			fmt.Printf(Usage)
			os.Exit(1)
		}
		fmt.Printf("创建区块链命令被调用\n")
		addr := cmds[2]
		cli.CreatBlockChain(addr)

	case "printChain":
		fmt.Printf("打印区块链命令被调用\n")
		cli.PrintChain()

	case "getBalance":
		fmt.Printf("获取余额命令被调用\n")
		cli.GetBalance(cmds[2])

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

	default:
		fmt.Printf("无效命令，请检查\n")
		fmt.Printf(Usage)

	}
}
