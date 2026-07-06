package handlers

import (
	"os"

	"github.com/Serge-Nook/topor/internal/core"
)

// PlainTextHandler handles .txt and .xml files.
// These formats have no internal metadata fields, so only filesystem
// timestamps are modified.
type PlainTextHandler struct{}

func (h *PlainTextHandler) ReadMetadata(path string) (*MetadataInfo, error) {
	info := &MetadataInfo{}
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	info.ModifyTime = fi.ModTime()
	info.CreationTime = fi.ModTime()
	return info, nil
}

func (h *PlainTextHandler) WriteMetadata(path string, _ *core.AuthorSettings, timeSets *core.TimeSettings) error {
	if timeSets != nil {
		return core.ApplyTimestamps(path, *timeSets)
	}
	return nil
}
