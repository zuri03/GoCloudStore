/*
This package contains system wide constants. Keeping all of these constants in a single
location makes it easier to sychronize the value and name of each constant across the system
*/
package constants

const (
	MAX_CACHE_BUFFER_SIZE int = 1024 //Max amount of data that can be read to memory before it must be stored in a file
	TEMP_BUFFER_SIZE      int = 256

	GET_PROTOCOL     string = "GET"
	ERROR_PROTOCOL   string = "ERR"
	SEND_PROTOCOL    string = "SND"
	DELETE_PROTOCOL  string = "DEL"
	PROCEED_PROTOCOL string = "PRC"
	SUCCESS_PROTOCOL string = "SCS"
)
