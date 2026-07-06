package handlers

import (
	"os"
	"regexp"
	"strings"

	"github.com/Serge-Nook/topor/internal/core"
)

// RTFHandler handles .rtf files.
type RTFHandler struct{}

var rtfAuthorRe = regexp.MustCompile(`\{\\author\s+([^}]*)\}`)
var rtfOperatorRe = regexp.MustCompile(`\{\\operator\s+([^}]*)\}`)

func (h *RTFHandler) ReadMetadata(path string) (*MetadataInfo, error) {
	info := &MetadataInfo{}
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	info.ModifyTime = fi.ModTime()
	info.CreationTime = fi.ModTime()

	data, err := os.ReadFile(path)
	if err != nil {
		return info, nil
	}
	content := string(data)

	if m := rtfAuthorRe.FindStringSubmatch(content); len(m) > 1 {
		info.Author = strings.TrimSpace(m[1])
	}
	if m := rtfOperatorRe.FindStringSubmatch(content); len(m) > 1 {
		info.LastModBy = strings.TrimSpace(m[1])
	}

	return info, nil
}

func (h *RTFHandler) WriteMetadata(path string, author *core.AuthorSettings, timeSets *core.TimeSettings) error {
	if author != nil && author.Action != core.AuthorActionNone {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)

		switch author.Action {
		case core.AuthorActionDelete:
			content = rtfAuthorRe.ReplaceAllString(content, "")
			content = rtfOperatorRe.ReplaceAllString(content, "")
		case core.AuthorActionReplace:
			if author.NewAuthor != "" {
				if rtfAuthorRe.MatchString(content) {
					content = rtfAuthorRe.ReplaceAllString(content, `{\author `+author.NewAuthor+`}`)
				}
			}
			if author.NewLastModBy != "" {
				if rtfOperatorRe.MatchString(content) {
					content = rtfOperatorRe.ReplaceAllString(content, `{\operator `+author.NewLastModBy+`}`)
				}
			}
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	if timeSets != nil {
		return core.ApplyTimestamps(path, *timeSets)
	}
	return nil
}
