/*
This package contains system wide constants. Keeping all of these constants in a single
location makes it easier to sychronize the value and name of each constant across the system
*/
package common

const (
	MAX_CACHE_BUFFER_SIZE int = 2048 //Max amount of data that can be read to memory before it must be stored in a file
	TEMP_BUFFER_SIZE      int = 1024

	//Used in situations for entities to express an internal error to the other entity it is connected to
	ERROR_PROTOCOL string = "ERR"

	//success message is used in situations where file data is not transferred E.g. file deletion
	SUCCESS_PROTOCOL string = "SCS"

	//Messages sent by the client to a storage server describing the desired operation
	GET_PROTOCOL    string = "GET"
	SEND_PROTOCOL   string = "SND"
	DELETE_PROTOCOL string = "DEL"

	//Proceed is used in situations where the actions of the client and server need to be coordinated to ensure no error occurs
	PROCEED_PROTOCOL string = "PRC"

	//Used in situations for entities to express an internal error to the other entity it is connected to
	ERROR_FRAME FrameType = -1

	//success message is used in situations where file data is not transferred E.g. file deletion
	SUCCESS_FRAME FrameType = 0

	//Messages sent by the client to a storage server describing the desired operation
	GET_FRAME    FrameType = 1
	SEND_FRAME   FrameType = 2
	DELETE_FRAME FrameType = 3

	//Proceed is used in situations where the actions of the client and server need to be coordinated to ensure no error occurs
	PROCEED_FRAME FrameType = 4
)
