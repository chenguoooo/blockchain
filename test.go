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

const walletsNum = 1
const blockNum = 1000
const transactionNum = 500

var mutexab sync.Mutex

func (test *Test) Run() {

	//随机数生成
	rand.Seed(time.Now().Unix())

	//创建钱包
	ws := NewWallets()
	ws.CreateWallet()
	fmt.Printf("创建钱包完成\n")

	addresses := ws.ListAddress()
	wslen := len(addresses)
	fmt.Printf("wslen:%d\n", wslen)

	//tps计数
	var timeSum int64
	var timeSumf float64
	var totaltps []int64
	timeSum = 0

	var ptimeSum int64
	var ptimeSumf float64
	var ptotaltps []int64
	ptimeSum = 0
	//生成区块链
	//addresses[rand.Intn(wslen)]
	bc := CreatBlockChain(addresses[rand.Intn(wslen)])
	bc = NewBlockChain()
	if bc != nil {
		defer bc.db.Close()
	}

	for i := 0; i < blockNum; i++ {
		addresses := ws.ListAddress()
		wslen := len(addresses)
		if i%3 == 0 {
			fmt.Printf("newwslen:%d\n", wslen)
			ws.CreateWallet()
		}

		//创建交易
		txs := []*Transaction{}
		for j := 0; j < wslen*2; j++ {
			from := addresses[rand.Intn(wslen)]
			to := addresses[rand.Intn(wslen)]
			//tx := NewTransaction(from, to, (rand.Float64())*float64(rand.Intn(10)), bc, ws, txs)
			tx := NewTransaction(from, to, 1, bc, ws, txs)
			if tx != nil {
				txs = append(txs, tx)
			}
		}

		fmt.Printf("创建交易成功!\n")
		fmt.Printf("txslen:%d\n", len(txs))

		//挖矿
		coinbase := NewCoinbaseTx(addresses[rand.Intn(wslen)], strconv.Itoa(i+1))
		txs2 := []*Transaction{}
		txs2 = append(txs2, coinbase)
		txs2 = append(txs2, txs...)

		var startTimeUnixNano int64
		var endTimeUnixNano int64
		var pstartTimeUnixNano int64
		var pendTimeUnixNano int64

		var validtxslen int64
		if i%2 == 0 {
			startTimeUnixNano = time.Now().UnixNano()
			validtxslen = bc.AddBlock(txs2)
			endTimeUnixNano = time.Now().UnixNano()
			timeSum = timeSum + endTimeUnixNano - startTimeUnixNano
			timeSumb := float64(endTimeUnixNano-startTimeUnixNano) / math.Pow(10, 9)
			timeSumf = float64(timeSum) / math.Pow(10, 9)
			tps := int64(float64(validtxslen) / timeSumb)
			totaltps = append(totaltps, tps)
			fmt.Printf("区块%d以成功创建!,耗时:%v s,tps:%d\n\n", i+2, timeSumb, tps)

		} else {
			pstartTimeUnixNano = time.Now().UnixNano()
			validtxslen = bc.AddBlockParal(txs2)
			pendTimeUnixNano = time.Now().UnixNano()
			ptimeSum = ptimeSum + pendTimeUnixNano - pstartTimeUnixNano
			ptimeSumb := float64(pendTimeUnixNano-pstartTimeUnixNano) / math.Pow(10, 9)
			ptimeSumf = float64(ptimeSum) / math.Pow(10, 9)
			ptps := int64(float64(validtxslen) / ptimeSumb)
			ptotaltps = append(ptotaltps, ptps)
			fmt.Printf("区块%d以成功创建!,耗时:%v s,tps:%d\n\n", i+2, ptimeSumb, ptps)

		}

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

	fmt.Printf("创建区块消耗%v\n", ptimeSumf)
	var ptotps int64
	ptotps = 0
	for _, t := range totaltps {
		ptotps += t
	}
	pavertps := ptotps / int64(len(ptotaltps))
	fmt.Printf("avertps:%v\n", pavertps)

}
