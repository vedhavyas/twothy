package twofa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// configName for the configFile
const configName string = ".2fa.json"

// Config holds the path of the accounts
type Config struct {
	AccountsFolder string `json:"accounts_folder"`
}

// LoadConfig loads the 2fa Config if configured,
// Else returns an error
func LoadConfig(filePath string) (c Config, err error) {
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
