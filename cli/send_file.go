package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func sendFileCommand(owner string, input []string, serverClient ServerClient) {
	fileName := input[0]
	file, fileInfo, err := getFileFromMemory(fileName)
	defer file.Close()
	if err != nil {
		fmt.Printf("Error getting file from memory: %s\n", err.Error())
		return
	}

	if fileInfo == nil || file == nil {
		fmt.Printf("Error could not find file %s\n", fileName)
		return
	}

	key := parseSendCommandInput(input[1:])

	//If no key is found then the default behavior is to set the key to filename
	if key == "" {
		key = fileInfo.Name()
	}

	err = serverClient.CreateFileRecord(owner, key, fileInfo.Name(), fileInfo.Size()) //For now just leave the key as the file name
	if err != nil {
		fmt.Printf("Error sending creating file record: %s\n", err.Error())
		return
	}

	if err := serverClient.SendFile(owner, fileInfo); err != nil {
		fmt.Println(err.Error())
	}
}

//This command parses any flags and returns the value assigned to that flag
//Currently the only availabe flag is -k or key if there are anymore flags added in the future parse them with this function
func parseSendCommandInput(input []string) string {
	fmt.Printf("Parsing input\n")
	key := ""

	for idx, part := range input {
		fmt.Printf("%d => %s\n", idx, part)
		//If -k is found then the next element in the input array should be the key
		if part == "-k" {
			key = input[idx+1]
		}
	}
	fmt.Printf("Key is now: %s\n", key)
	return key
}

func getFileFromMemory(fileName string) (*os.File, fs.FileInfo, error) {
	fileExtension := filepath.Ext(fileName)
	if fileExtension != ".txt" && fileExtension != ".rtf" && fileExtension != ".pdf" {
		return nil, nil, fmt.Errorf("Not a text file")
	}

	fileMetaData, err := os.Stat(fileName)
	if err != nil {
		return nil, nil, err
	}

	if fileMetaData.IsDir() {
		return nil, nil, fmt.Errorf("Cannot send directory to server => %s\n", err.Error())
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}

	return file, fileMetaData, nil
}
