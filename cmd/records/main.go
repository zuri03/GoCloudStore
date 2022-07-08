package main

import (
	"fmt"
	"time"

	"github.com/zuri03/GoCloudStore/records"
)

func main() {
	fmt.Println("CREATING META DATA SERVER")
	keeper := records.InitRecordKeeper()
	fmt.Println("CREATED META DATA SERVER")
	users := records.UserClient{
		Timeout: time.Duration(time.Second * 10),
	}
	records.InitServer(&keeper, &users)
}
