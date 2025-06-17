package types

type MessageType string

const (
	MessageText   MessageType = "text"
	MessageImage  MessageType = "image"
	MessageFile   MessageType = "file"
	MessageSystem MessageType = "system"
)
