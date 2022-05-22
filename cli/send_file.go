package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func getFileFromMemory(input []string) (*os.File, fs.FileInfo, error) {
	fileName := input[0]
	fmt.Printf("Found file => %s\n", fileName)
	fileExtension := filepath.Ext(fileName)
	if fileExtension != ".txt" && fileExtension != ".rtf" && fileExtension != ".pdf" {
		return nil, nil, fmt.Errorf("Not a text file")
	}

	fileMetaData, err := os.Stat(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting file metadata => %s\n", err.Error())
	}

	if fileMetaData.IsDir() {
		return nil, nil, fmt.Errorf("Cannot send directory to server => %s\n", err.Error())
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("File %s does not exist", fileName)
	}

	fmt.Printf("returing file of size = %d\n", fileMetaData.Size())
	return file, fileMetaData, nil
}
