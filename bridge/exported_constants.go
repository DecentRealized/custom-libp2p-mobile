package custom_libp2p_bridge

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/database"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/p2p"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/transfer"
)

var exportedErrors = map[string]error{
	"ErrMethodNotImplemented":     ErrMethodNotImplemented,
	"ErrDatabaseRunning":          database.ErrDatabaseRunning,
	"ErrDatabaseStopped":          database.ErrDatabaseStopped,
	"ErrNotifierClosed":           notifier.ErrNotifierClosed,
	"ErrNodeDoesNotExist":         p2p.ErrNodeDoesNotExist,
	"ErrMetafileAlreadyExists":    transfer.ErrMetafileAlreadyExists,
	"ErrAlreadyDownloadingFile":   transfer.ErrAlreadyDownloadingFile,
	"ErrFileMetadataNotAvailable": transfer.ErrFileMetadataNotAvailable,
	"ErrFileNotDownloading":       transfer.ErrFileNotDownloading,
	"ErrFileNotServing":           transfer.ErrFileNotServing,
	"ErrSendingMessage":           transfer.ErrSendingMessage,
	"ErrClientNotRunning":         transfer.ErrClientNotRunning,
	"ErrServerNotRunning":         transfer.ErrServerNotRunning,
	"ErrNotAllowedNode":           transfer.ErrNotAllowedNode,
	"ErrForbidden":                transfer.ErrForbidden,
	"ErrSha256DoesNotMatch":       transfer.ErrSha256DoesNotMatch,
}
