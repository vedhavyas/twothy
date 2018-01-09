package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/vedhavyas/twothy"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 || strings.Contains(args[0], "help") {
		usage()
		return
	}

	cmd := args[0]
	result, err := twothy.ExecOp(cmd, args[1:]...)
	if err != nil {
		fmt.Printf("%v\n\n", err)
		usage()
		return
	}

	fmt.Print("\n" + result)
}

func usage() {
	fmt.Printf("Usage: %s command [arguments]\n", os.Args[0])
	fmt.Println("Commands and arguments:")
	w := tabwriter.NewWriter(os.Stdout, 7, 7, 0, '\t', tabwriter.AlignRight)
	fmt.Fprint(w, "\tconfigure:\tconfigures twothy\n")
	fmt.Fprint(w, "\tadd [name] [label] [key(base32)]:\tadds a new account with given info\n")
	fmt.Fprint(w, "\totp:\tgenerates otp for all accounts\n")
	fmt.Fprint(w, "\totp [name]:\tgenerates otp for accounts matching name\n")
	fmt.Fprint(w, "\totp [name] [label]:\tgenerates otp for accounts matching name and label\n")
	w.Flush()
	os.Exit(1)
}
