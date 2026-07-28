package enums

type MessageType string

const (
	MessageText   MessageType = "TEXT"
	MessageImage  MessageType = "IMAGE"
	MessageFile   MessageType = "FILE"
	MessageSystem MessageType = "SYSTEM"
)
