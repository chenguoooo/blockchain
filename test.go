package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Test struct {
}

const walletsNum = 500
const blockNum = 100
const transactionNum = 500

var mutexab sync.Mutex

func (test *Test) Run() {

	//生成钱包
	for i := 0; i < walletsNum; i++ {

		ws := NewWallets()
		ws.CreateWallet()

	}
	fmt.Printf("创建钱包完成\n")

	//随机数生成
	rand.Seed(time.Now().Unix())
	ws := NewWallets()
	addresses := ws.ListAddress()
	wslen := len(addresses)
	//addresses[rand.Intn(wslen)]
	var timeSum int64
	var timeSumf float64
	var totaltps []int64
	timeSum = 0
	bc := CreatBlockChain(addresses[rand.Intn(wslen)])
	if bc != nil {
		defer bc.db.Close()
	}
	//fmt.Printf("创建区块链成功!\n")
	//创建区块链
	//c2 := make(chan int, transactionNum)
	for i := 0; i < blockNum; i++ {
		//创建交易
		txs := NewTransactions()

		for j := 0; j < transactionNum; j++ {

			from := addresses[rand.Intn(wslen)]
			to := addresses[rand.Intn(wslen)]
			tx := NewTransaction(from, to, 1, bc, ws)
			txs.CreateTransaction(tx)

		}
		fmt.Printf("创建交易成功!\n")

		txs1 := []*Transaction{}
		for _, tx := range txs.TransactionsMap {
			txs1 = append(txs1, tx)
		}
		fmt.Printf("tx:%d\n", len(txs1))

		//挖矿cli.AddBlock(addresses[rand.Intn(wslen)], "111", txs1)
		coinbase := NewCoinbaseTx(addresses[rand.Intn(wslen)], strconv.Itoa(i+1))
		txs2 := []*Transaction{}
		txs2 = append(txs2, coinbase)
		txs2 = append(txs2, txs1...)

		timeUnixNano1 := time.Now().UnixNano()

		//validtxslen := bc.AddBlockParal(txs2)
		validtxslen := bc.AddBlock(txs2)

		timeUnixNano2 := time.Now().UnixNano()
		txs.ClearFile()

		timeSum = timeSum + timeUnixNano2 - timeUnixNano1
		timeSumb := float64(timeUnixNano2-timeUnixNano1) / math.Pow(10, 9)
		timeSumf = float64(timeSum) / math.Pow(10, 9)
		tps := int64(float64(validtxslen) / timeSumb)
		totaltps = append(totaltps, tps)
		fmt.Printf("区块%d以成功创建!,耗时:%v s,tps:%d\n", i+2, timeSumb, tps)

	}

	fmt.Printf("\n完成!\n\n")
	fmt.Printf("创建区块消耗%v\n", timeSumf)
	var totps int64
	totps = 0
	for _, t := range totaltps {
		totps += t
	}
	avertps := totps / int64(len(totaltps))
	fmt.Printf("avertps:%v\n", avertps)

}
