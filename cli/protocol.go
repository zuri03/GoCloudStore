package cli

type Protocol string

const (
	GET_PROTOCOL    Protocol = "GET"
	ERROR_PROTOCOL  Protocol = "ERR"
	SEND_PROTOCOL   Protocol = "SND"
	DELETE_PROTOCOL Protocol = "DEL"
)
