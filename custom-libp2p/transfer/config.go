package transfer

const protocolID = "/file-drop/1.0.0"

const maxMessageSize = uint16(1 << 14)

const autoDownloadEnabled = true // Todo provide user with option

const (
	holePunchSyncStreamProtocolID = "/holepunch-sync-stream/1.0.0"
	holePunchRetries              = 3
	holePunchPacketSize           = 1 << 1
)

const (
	unAuthorizedHeader     = uint8(0)
	streamAuthorizedHeader = uint8(1)
)
