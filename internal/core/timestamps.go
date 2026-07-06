package core

import (
	"os"
	"time"
)

// ApplyTimestamps sets file system timestamps according to settings.
func ApplyTimestamps(path string, settings TimeSettings) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	modTime := info.ModTime()
	// On most systems, creation time is not easily accessible via Go stdlib.
	// We use modTime as a fallback for creation time.
	createTime := modTime

	newCreate := createTime
	newModify := modTime

	// Compute new creation time
	switch settings.CreationMode {
	case DateModeAbsolute:
		newCreate = settings.CreationAbsolute
	case DateModeOffset:
		newCreate = createTime.Add(
			time.Duration(settings.CreationOffsetD)*24*time.Hour +
				time.Duration(settings.CreationOffsetH)*time.Hour +
				time.Duration(settings.CreationOffsetM)*time.Minute,
		)
	case DateModeEqualCreateToModify:
		newCreate = modTime
	}

	// Compute new modify time
	switch settings.ModifyMode {
	case DateModeAbsolute:
		newModify = settings.ModifyAbsolute
	case DateModeOffset:
		newModify = modTime.Add(
			time.Duration(settings.ModifyOffsetD)*24*time.Hour +
				time.Duration(settings.ModifyOffsetH)*time.Hour +
				time.Duration(settings.ModifyOffsetM)*time.Minute,
		)
	case DateModeEqualModifyToCreate:
		newModify = createTime
	}

	// If creation time was changed and modify mode is "equal modify to create",
	// use the newly computed creation time
	if settings.ModifyMode == DateModeEqualModifyToCreate && settings.CreationMode != DateModeNone {
		newModify = newCreate
	}
	if settings.CreationMode == DateModeEqualCreateToModify && settings.ModifyMode != DateModeNone {
		newCreate = newModify
	}

	// os.Chtimes sets both access time and modification time.
	// For creation time on Linux, there's no direct API, so we set mod time.
	// The access time is set to the creation time value.
	if settings.CreationMode != DateModeNone || settings.ModifyMode != DateModeNone {
		atime := newCreate
		mtime := newModify
		return os.Chtimes(path, atime, mtime)
	}
	return nil
}
