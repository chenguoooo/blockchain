package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

//wallets结构
//把地址和密钥对对应起来
//map[address1]->walletKeyPair1

type Wallets struct {
	WalletsMap map[string]*WalletKeyPair
}

//创建wallets，返回wallets实例
func NewWallets() *Wallets {
	var ws Wallets

	ws.WalletsMap = make(map[string]*WalletKeyPair)
	//把所有的钱包从本地加载出来
	if !ws.LoadFromFile() {
		fmt.Printf("加载数据失败!\n")
	}

	//把实例返回
	return &ws
}

//这个wallets是对外的，walletkeypair是对内的
//wallets调用walletkeypair

const WalletName = "wallet.dat"

func (ws *Wallets) CreateWallet() string {
	//调用Newwalletkeypair
	wallet := NewWalletKeyPair()
	//将返回的Newwalletkeypair添加到walletmap中
	address := wallet.GetAddress()

	ws.WalletsMap[address] = wallet
	//
	//var wsLocal Wallets
	//wsLocal.WalletsMap[address] = wallet
	//ws.WalletsMap = wsLocal.WalletsMap
	//保存到本地文件
	res := ws.SaveToFile()
	if !res {
		fmt.Printf("创建钱包失败!\n")
		return ""
	}

	return address
}

//保存钱包到文件
func (ws *Wallets) SaveToFile() bool {

	var buffer bytes.Buffer
	//将接口类型明确注册一下，否则gob编码失败
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		fmt.Printf("钱包序列化失败,err:%v\n", err)
	}
	content := buffer.Bytes()

	//func WriteFile(filename string, data []byte, perm os.FileMode) error {
	err = ioutil.WriteFile(WalletName, content, 0600)
	if err != nil {
		fmt.Printf("钱包创建失败！\n")
		return false
	}

	return true

}

func (ws *Wallets) LoadFromFile() bool {
	//判断文件是否存在
	if !IsFileExist(WalletName) {
		fmt.Printf("钱包文件不存在，准备创建！\n")
		return true
	}

	//读取文件

	content, err := ioutil.ReadFile(WalletName)

	if err != nil {
		return false
	}
	gob.Register(elliptic.P256())
	//gob解码
	decoder := gob.NewDecoder(bytes.NewReader(content))

	var wallets Wallets
	err = decoder.Decode(&wallets)

	if err != nil {
		fmt.Printf("err:%v\n", err)
		return false
	}
	//赋值给ws
	ws.WalletsMap = wallets.WalletsMap
	return true

}

func (ws *Wallets) ListAddress() []string {
	//遍历ws.WalletsMap结构返回key即可
	var addresses []string
	for address, _ := range ws.WalletsMap {
		addresses = append(addresses, address)
	}
	return addresses
}
