package transfer

import "errors"

var (
	InstanceDoesNotExist     = errors.New("transfer instance does not exist")
	MessageSizeTooLarge      = errors.New("message size too large")
	MessageSizeZero          = errors.New("message size zero")
	StreamToPeerNotAllowed   = errors.New("stream to peer not allowed")
	InvalidOffset            = errors.New("invalid offset")
	InvalidFileSize          = errors.New("invalid file size")
	MetafileAlreadyExists    = errors.New("meta file already exists")
	AlreadyDownloadingFile   = errors.New("file already downloading")
	FileMetadataNotAvailable = errors.New("file metadata not available")
	FileNotDownloading       = errors.New("file not downloading")
	FileNotServing           = errors.New("file not serving")
	ErrorSendingMessage      = errors.New("error sending message")
	NotRunning               = errors.New("transfer not running")
)
