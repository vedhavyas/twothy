package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/vedhavyas/twothy"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		flag.PrintDefaults()
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

	cmd := args[0]
	result, err := twothy.ExecOp(cmd, args[1:]...)
	if err != nil {
		fmt.Printf("%s: failed due to: %v\n", cmd, err)
		os.Exit(1)
		return
	}

	fmt.Print(result)
}
