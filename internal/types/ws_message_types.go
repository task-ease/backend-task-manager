package types

type WsMessageTypes string

const (
	TypeMessage             WsMessageTypes = "MESSAGE"
	TypeMessageNotification WsMessageTypes = "MESSAGE_NOTIFICATION"
	TypeConnected           WsMessageTypes = "CONNECTED"
	TypeDisconnected        WsMessageTypes = "DISCONNECTED"
)
