package cli

import (
	"fmt"
)

func printHelpMessage() {
	fmt.Printf(`
		cli:
			- This command starts a cli session that keeps you logged in and allows you to send and retreive multiple documents from the server
			- While in a cli session you can omit including your username and password in every command
		get: 
			- Retreives a file from the server and saves it to the current directory
			- Example:
				- With cli session: get [filename]
				- Without cli: get [username] [password] [filename]
		send: 
			- Sends a file in the current directory to the server
			- Example:
				- With cli session: send [filepath]
				- Without cli: send [username] [password] [filepath]
		send: 
			- Deletes a file on the server. DOES NOT DELETE THE FILE ON YOUR LOCAL COMPUTER
			- Example:
				- With cli session: delete [filename]
				- Without cli: delete [username] [password] [filename]
		quit:
			- If you are currently in a cli session this command will gracefully close the session
			- If you are not in a cli session this command does nothing
	`)
}
