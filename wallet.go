package main

import (
	"base58"
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
