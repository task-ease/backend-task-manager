package types

type WsMessageTypes string

const (
	TypeMessage      WsMessageTypes = "message"
	TypeConnected    WsMessageTypes = "connected"
	TypeDisconnected WsMessageTypes = "disconnected"
)
