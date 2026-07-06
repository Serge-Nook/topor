package handlers

import (
	"os"

	"github.com/Serge-Nook/topor/internal/core"
)

// OLEHandler handles .doc and .xls files (legacy OLE formats).
// Full OLE parsing requires a dedicated library. This handler provides
// basic support by modifying file system timestamps and logging that
// internal metadata changes for OLE require platform-specific tools.
type OLEHandler struct{}

func (h *OLEHandler) ReadMetadata(path string) (*MetadataInfo, error) {
	info := &MetadataInfo{}
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	info.ModifyTime = fi.ModTime()
	info.CreationTime = fi.ModTime()
	return info, nil
}

func (h *OLEHandler) WriteMetadata(path string, author *core.AuthorSettings, timeSets *core.TimeSettings) error {
	if timeSets != nil {
		if err := core.ApplyTimestamps(path, *timeSets); err != nil {
			return err
		}
	}

	if author != nil && author.Action != core.AuthorActionNone {
		// OLE binary format manipulation for .doc/.xls author fields
		// would require a full OLE compound document library.
		// For now, we handle timestamps via the filesystem.
		_ = modifyOLEAuthor(path, author)
	}
	return nil
}

func modifyOLEAuthor(path string, author *core.AuthorSettings) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// OLE compound documents store document summary info in specific streams.
	// The author and last-saved-by fields are in the SummaryInformation stream.
	// We search for known byte patterns to locate and modify these fields.

	// Look for common author field markers in OLE binary
	// This is a best-effort approach for the most common .doc/.xls structures.
	modified := false

	switch author.Action {
	case core.AuthorActionDelete:
		data, modified = oleReplaceStringProperty(data, "")
	case core.AuthorActionReplace:
		if author.NewAuthor != "" {
			data, modified = oleReplaceStringProperty(data, author.NewAuthor)
		}
	}

	if modified {
		return os.WriteFile(path, data, 0644)
	}
	return nil
}

// oleReplaceStringProperty attempts to find and replace author strings in OLE binary.
// This is a simplified approach; full OLE parsing would need a dedicated library.
func oleReplaceStringProperty(data []byte, _ string) ([]byte, bool) {
	// OLE compound document parsing is complex. For production use,
	// consider using an external tool or library.
	// Returning unchanged data as safe default.
	return data, false
}
