package twothy

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mitchellh/go-homedir"
)

// ExecOp executes the given operation
func ExecOp(cmd string, args ...string) (result string, err error) {
	switch cmd {
	case "configure":
		hd, err := homedir.Dir()
		if err != nil {
			return result, fmt.Errorf("failed to get user's home directory: %v", err)
		}

		c, err := Configure(hd)
		if err != nil {
			return result, fmt.Errorf("failed to configure twothy: %v", err)
		}

		return fmt.Sprintf("2FA accounts will be at: %s\n", c.AccountsFolder), nil

	case "add":
		c, err := GetConfig()
		if err != nil {
			return result, fmt.Errorf("failed to configure twothy: %v", err)
		}

		if len(args) < 3 {
			return result, fmt.Errorf("require Name, Label, and Key(base32) to add an account")
		}

		a := NewAccount(args[0], args[1], args[2])
		err = SaveAccount(c, a)
		if err != nil {
			return result, fmt.Errorf("failed to save account: %v", err)
		}

		otp, err := CreateOTP(a, time.Now().Unix())
		if err != nil {
			return result, fmt.Errorf("failed to generate otp for account %s@%s: %v", a.Name, a.Label, err)
		}

		return fmt.Sprintf("%s@%s: %d\n", a.Label, a.Name, otp), nil

	case "otp":
		fallthrough
	default:
		c, err := GetConfig()
		if err != nil {
			return result, fmt.Errorf("failed to configure twothy: %v", err)
		}

		var name, label string
		if len(args) >= 1 {
			name = args[0]
		}

		if len(args) >= 2 {
			label = args[1]
		}

		accounts, err := LoadAccounts(c, name, label)
		if len(accounts) < 1 {
			return fmt.Sprintf("No accounts associated with '%s' are found\n", name), nil
		}

		var b bytes.Buffer
		for _, a := range accounts {
			otp, err := CreateOTP(a, time.Now().Unix())
			if err != nil {
				return result, fmt.Errorf("failed to generate otp for account %s@%s: %v", a.Name, a.Label, err)
			}

			b.WriteString(fmt.Sprintf("%s@%s: %d\n", a.Label, a.Name, otp))
		}

		return b.String(), nil
	}
}
