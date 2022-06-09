package storage


import (
	"net"
)

func initListener() {
	//In the future the port will come from environment vars/ command line args
	listender, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Printf("Error in listener: %s\n", err.Error())
		return
	}

	for {
		fmt.Println("Waiting for connections")
		connection, err := listener.Accept()

		if err != nil {
			fmt.Printf("Error accepting connections => %s\n", err.Error())
			break
		}
	}
}