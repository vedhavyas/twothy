package main

import (
	"fmt"
	"log"

	"github.com/vedhavyas/twothy"
)

func main() {
	config, err := twothy.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config.AccountsFolder)
}
