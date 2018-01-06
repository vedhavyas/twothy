package twothy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// writeToFile write obj to file
func writeToFile(filePath string, obj interface{}) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	defer f.Close()
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshall the object: %v", err)
	}

	w := bufio.NewWriter(f)
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return w.Flush()
}
