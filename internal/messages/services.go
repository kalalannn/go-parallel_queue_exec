package messages

const (
	NewWebSocketClient                 = "Added new WebSocket client"
	RemoveWebSocketClient              = "Removed WebSocket client"
	BroadcastFailed                    = "Failed to send update to client: %s, disconnecting\n"
	StopBroadcastUpdatesStart          = "Stopping broadcast updates..."
	StopBroadcastUpdatesEnd            = "Broadcast updates stopped."
	WaitForExecutorShutdownWithTimeout = "Waiting for executor shutdown (timeout: %s)."
	ExecutorShutdownFinished           = "Executor shutdown finished."
	ExecutorShutdownTimeout            = "Executor shutdown timeout."
)
