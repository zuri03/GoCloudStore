/*
This package contains system wide constants. Keeping all of these constants in a single
location makes it easier to sychronize the value and name of each constant across the system
*/
package constants

const (
	MAX_CACHE_BUFFER_SIZE int = 1024 //Max amount of data that can be read to memory before it must be stored in a file
	TEMP_BUFFER_SIZE      int = 256

	//Messages sent by the client to a storage server describing the desired operation
	GET_PROTOCOL    string = "GET"
	SEND_PROTOCOL   string = "SND"
	DELETE_PROTOCOL string = "DEL"

	//Used in situations for entities to express an internal error to the other entity it is connected to
	ERROR_PROTOCOL string = "ERR"

	//Proceed is used in situations where the actions of the client and server need to be coordinated to ensure no error occurs
	PROCEED_PROTOCOL string = "PRC"
	//success message is used in situations where file data is not transferred E.g. file deletion
	SUCCESS_PROTOCOL string = "SCS"
)
