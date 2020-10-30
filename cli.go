package main

import (
	"fmt"
	"os"
	"strconv"
)

const Usage = `
	./blockchain printChain	打印区块链
	./blockchain getBalance	地址获取地址的余额
	./blockchain send FROM TO AMOUNT MINER DATA	转账命令
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
		if len(cmds) != 3 {
			fmt.Printf(Usage)
			os.Exit(1)
		}
		fmt.Printf("添加区块链命令被调用，数据：%s\n", cmds[2])
		//
		//data := cmds[2]
		//cli.AddBlock(data)//TODO

	case "printChain":
		fmt.Printf("打印区块链命令被调用\n")
		cli.PrintChain()

	case "getBalance":
		fmt.Printf("获取余额命令被调用\n")
		cli.bc.GetBalance(cmds[2])

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
