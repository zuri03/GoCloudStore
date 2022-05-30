package main

import (
	"fmt"

	"github.com/zuri03/GoCloudStore/records"
)

func main() {
	fmt.Println("CREATING META DATA SERVER")
	keeper := records.InitRecordKeeper()
	records.InitServer(&keeper)
	fmt.Println("CREATED META DATA SERVER")
}
