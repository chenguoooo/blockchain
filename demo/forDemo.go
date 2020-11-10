package main

import (
	"fmt"
	"time"
)

func main() {
	a := 1

	for i := 0; i < 100; i++ {
		func() {
			fmt.Println(a)
			a++
			time.Sleep(1 * time.Second) //延时1s
		}()
	}
	time.Sleep(200 * time.Second) //延时1s
}
