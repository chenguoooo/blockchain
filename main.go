package main

func main() {
	bc := NewBlockChain("cg")
	defer bc.db.Close()
	cli := CLI{bc}
	cli.Run()
}
