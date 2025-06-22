package types

type WsMessageTypes string

const (
	TypeMessage             WsMessageTypes = "message"
	TypeMessageNotification WsMessageTypes = "message_notification"
	TypeConnected           WsMessageTypes = "connected"
	TypeDisconnected        WsMessageTypes = "disconnected"
)
