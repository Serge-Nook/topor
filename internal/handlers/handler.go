package handlers

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/Serge-Nook/topor/internal/core"
)

// MetadataInfo holds metadata extracted from a document.
type MetadataInfo struct {
	Author       string
	LastModBy    string
	CreationTime time.Time
	ModifyTime   time.Time
}

// Handler interface for reading/writing document metadata.
type Handler interface {
	ReadMetadata(path string) (*MetadataInfo, error)
	WriteMetadata(path string, author *core.AuthorSettings, timeSets *core.TimeSettings) error
}

// GetHandler returns the appropriate handler for a file path.
func GetHandler(path string) Handler {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".docx", ".xlsx", ".docm", ".xlsm":
		return &OOXMLHandler{}
	case ".odt":
		return &ODTHandler{}
	case ".doc", ".xls":
		return &OLEHandler{}
	case ".html", ".htm":
		return &HTMLHandler{}
	case ".rtf":
		return &RTFHandler{}
	case ".txt", ".xml":
		return &PlainTextHandler{}
	default:
		return nil
	}
}
