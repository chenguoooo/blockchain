package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	strsA := []string{"hello", "world", "itcast"}

	strRes := strings.Join(strsA, "=")
	fmt.Println("strRes:", strRes)

	joinRes := bytes.Join([][]byte{[]byte("hello"), []byte("world"), []byte("itcast")}, []byte("="))
	fmt.Printf("joinRes:%s", joinRes)

}
