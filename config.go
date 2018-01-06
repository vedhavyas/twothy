package twothy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// configName for the configFile
const configName string = ".twothy.json"

// Config holds the path of the accounts
type Config struct {
	AccountsFolder string `json:"accounts_folder"`
}

// loadConfig loads the 2fa Config if configured,
// else returns an error
func loadConfig(filePath string) (c Config, err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c, fmt.Errorf("failed to load Config file: %v", err)
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return c, fmt.Errorf("failed to unmarshall Config: %v", err)
	}

	return c, nil
}

// configure configures twothy
func configure(homeDir string) (config Config, err error) {
	fmt.Println("Welcome to twothy!!")
	fmt.Println("Enter the path to store your 2FA accounts.")
	fmt.Println("I will create 'twothy_accounts' folder inside the given folder.")
	fmt.Println("If you are restoring accounts, provide root path to 'twothy_accounts'.")

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Path(%s):", homeDir)
	dir, err := reader.ReadString('\n')
	if err != nil {
		return config, fmt.Errorf("failed to scan user choice: %v", err)
	}

	if strings.TrimSpace(dir) == "" {
		dir = homeDir
	}
	config.AccountsFolder = dir
	err = writeToFile(fmt.Sprintf("%s/%s", homeDir, configName), config)
	if err != nil {
		return config, fmt.Errorf("failed to write twothy config file: %v", err)
	}

	return config, nil
}

// GetConfig returns twothy config
func GetConfig() (config Config, err error) {
	hd, err := homedir.Dir()
	if err != nil {
		return config, fmt.Errorf("failed to get user home directory: %v", err)
	}

	config, err = loadConfig(fmt.Sprintf("%s/%s", hd, configName))
	if err == nil {
		return config, nil
	}

	config, err = configure(hd)
	if err != nil {
		return config, fmt.Errorf("failed to configure twothy: %v", err)
	}

	return config, nil
}
