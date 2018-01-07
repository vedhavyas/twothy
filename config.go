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

const (
	// configName for the configFile
	configName = ".twothy.json"

	// accountsFilename for accounts
	accountsFolder = "twothy_accounts"
)

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
	fmt.Printf("I will create '%s' folder inside the given folder.\n", accountsFolder)
	fmt.Printf("If you are restoring accounts, provide path to '%s'.\n", accountsFolder)

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Path(%s):", homeDir)
	dir, err := reader.ReadString('\n')
	if err != nil {
		return config, fmt.Errorf("failed to scan user choice: %v", err)
	}

	dir = strings.TrimSpace(dir)
	if dir == "" {
		dir = homeDir
	}

	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	if !strings.HasSuffix(dir, accountsFolder+"/") {
		dir = dir + accountsFolder + "/"
	}

	err = os.MkdirAll(dir, 0766)
	if err != nil {
		return config, fmt.Errorf("failed to create %s: %v", dir, err)
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
