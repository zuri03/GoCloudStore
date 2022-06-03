package main

import (
	"fmt"

	"github.com/zuri03/GoCloudStore/records"
)

func main() {
	fmt.Println("CREATING META DATA SERVER")
	keeper := records.InitRecordKeeper()
	fmt.Println("CREATED META DATA SERVER")
	records.InitServer(&keeper)
}
