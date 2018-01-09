package twothy

import (
	"bufio"
	"fmt"
	"os"
)

func writeToFile(filePath string, data []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return w.Flush()
}
