package transfer

import "errors"

//ErrInstanceDoesNotExist   = errors.New("transfer instance does not exist")
//ErrStreamToPeerNotAllowed = errors.New("stream to peer not allowed")
//ErrInvalidOffset          = errors.New("invalid offset")
//InvalidFileSize           = errors.New("invalid file size")

var (
	ErrMetafileAlreadyExists    = errors.New("meta file already exists")
	ErrAlreadyDownloadingFile   = errors.New("file already downloading")
	ErrFileMetadataNotAvailable = errors.New("file metadata not available")
	ErrFileNotDownloading       = errors.New("file not downloading")
	ErrFileNotServing           = errors.New("file not serving")
	ErrSendingMessage           = errors.New("error sending message")
	ErrClientNotRunning         = errors.New("client not running")
	ErrServerNotRunning         = errors.New("server not running")
	ErrNotAllowedNode           = errors.New("peer node is not allowed")
	ErrForbidden                = errors.New("peer forbids this action")
)
