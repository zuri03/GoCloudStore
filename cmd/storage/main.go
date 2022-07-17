package main

import (
	"fmt"

	"github.com/zuri03/GoCloudStore/storage"
)

func main() {
	storage.InitializeListener()
	fmt.Println("Exiting")
}
