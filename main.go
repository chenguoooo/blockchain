package main

import (
	"fmt"
	"os"
)

func main() {
	//bc := NewBlockChain("cg")
	//defer bc.db.Close()
	//cli := CLI{bc}
	cli := CLI{}
	test := Test{}

	cmds := os.Args
	if len(cmds) < 2 {
		test.Run()
		fmt.Printf(Usage)
	} else {
		cli.Run()
	}

}
