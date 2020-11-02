package main

import (
	"base58"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type WalletKeyPair struct {
	PrivateKey *ecdsa.PrivateKey

	//type PublicKey struct {
	//	elliptic.Curve
	//	X, Y *big.Int
	//}

	//将公钥的X,Y进行字节流拼接后传输
	PublicKey []byte
}

func NewWalletKeyPair() *WalletKeyPair {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	publicKeyRaw := privateKey.PublicKey
	publicKey := append(publicKeyRaw.X.Bytes(), publicKeyRaw.Y.Bytes()...)
	return &WalletKeyPair{PrivateKey: privateKey, PublicKey: publicKey}
}

func (w *WalletKeyPair) GetAddress() string {
	publicHash := HashPubKey(w.PublicKey)

	version := 0x00

	//21字节数据
	payload := append([]byte{byte(version)}, publicHash...)

	checksum := CheckSum(payload)

	//25字节数据
	payload = append(payload, checksum...)
	address := base58.Encode(payload)

	return address

}

func IsValidAddress(address string) bool {
	//1.将输入的地址进行解码得到25字节
	//2.取出前21个字节，进行checksum函数，得到checksum1
	//3.取出后四个字节，得到checksum2
	//4.比较checksum1和checksum2，如果地址相同有效，反之无效
	decodeInfo, _ := base58.Decode(address)

	if len(decodeInfo) != 25 {
		return false
	}

	payload := decodeInfo[0 : len(decodeInfo)-4]
	checksum1 := CheckSum(payload)
	checksum2 := decodeInfo[len(decodeInfo)-4:]

	return bytes.Equal(checksum1, checksum2)

}

func HashPubKey(pubKey []byte) []byte {
	hash := sha256.Sum256(pubKey)

	//创建一个hash160对象
	//想hash160中write数据
	//做hash运算
	rip160Hasher := ripemd160.New()
	_, err := rip160Hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}

	//Sum函数会把我们的结果与Sum参数append到一起，然后返回，我们传入nil，放置数据污染
	publicHash := rip160Hasher.Sum(nil)

	return publicHash
}

func CheckSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	//4字节校验码
	checksum := second[0:4]

	return checksum
}
