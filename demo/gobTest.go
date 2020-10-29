package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//gob是go语言内置的编码包
//它可以对任意数据类型进行编码和解码
//编码时，先要创建编码器，编码器进行编码
//解码时，先要创建解码器，解码器进行解码

type Person struct {
	Name string
	Age  uint64
}

func main() {
	Jim := Person{
		Name: "Jim",
		Age:  19,
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&Jim)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("编码后的数据%x\n", buffer.Bytes())

	//传输中

	//解码，将字节六转换成Person结构
	var p1 Person

	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	err = decoder.Decode(&p1)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("解码后的数据%v\n", p1)
}
