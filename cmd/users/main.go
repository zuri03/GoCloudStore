package main

import (
	"fmt"

	"github.com/zuri03/GoCloudStore/users"
)

func main() {
	fmt.Println("CREATING META DATA SERVER")
	users.InitServer()
}
