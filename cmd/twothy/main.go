package main

import (
	"fmt"
	"os"

	"github.com/vedhavyas/twothy"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		//TODOs
		os.Exit(1)
	}

	cmd := args[0]
	result, err := twothy.ExecOp(cmd, args[1:]...)
	if err != nil {
		fmt.Printf("%s: failed due to: %v\n", cmd, err)
		os.Exit(1)
		return
	}

	fmt.Print(result)
}
