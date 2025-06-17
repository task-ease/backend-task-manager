package types

type WsMessageTypes string

const (
	TypeMessage WsMessageTypes = "message"
	TypeStatus  WsMessageTypes = "status"
)
