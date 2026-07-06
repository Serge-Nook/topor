package core

import "time"

// DateMode defines how dates should be modified.
type DateMode int

const (
	DateModeNone DateMode = iota
	DateModeAbsolute
	DateModeOffset
	DateModeEqualCreateToModify
	DateModeEqualModifyToCreate
)

// AuthorAction defines what to do with author fields.
type AuthorAction int

const (
	AuthorActionNone AuthorAction = iota
	AuthorActionDelete
	AuthorActionReplace
)

// TimeSettings holds timestamp modification parameters.
type TimeSettings struct {
	CreationMode     DateMode
	CreationAbsolute time.Time
	CreationOffsetD  int
	CreationOffsetH  int
	CreationOffsetM  int

	ModifyMode     DateMode
	ModifyAbsolute time.Time
	ModifyOffsetD  int
	ModifyOffsetH  int
	ModifyOffsetM  int
}

// AuthorSettings holds author modification parameters.
type AuthorSettings struct {
	Action       AuthorAction
	NewAuthor    string
	NewLastModBy string
}

// ProcessingSettings combines all processing parameters.
type ProcessingSettings struct {
	Time         TimeSettings
	Author       AuthorSettings
	CreateBackup bool
}

// FileMetadata holds extracted metadata for a file.
type FileMetadata struct {
	Path         string
	Name         string
	Size         int64
	Format       string
	CreationTime time.Time
	ModifyTime   time.Time
	Author       string
	LastModBy    string
}

// ProcessingStatus represents the result status of processing a file.
type ProcessingStatus int

const (
	StatusSuccess ProcessingStatus = iota
	StatusError
	StatusSkipped
)

// ProcessingResult holds the result of processing a single file.
type ProcessingResult struct {
	Path    string
	Status  ProcessingStatus
	Message string
}
