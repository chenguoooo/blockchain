package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	fmt.Println("My first lucky number is", rand.Intn(10))
	fmt.Println("My senond lucky number is", rand.Intn(10))
}
