package handlers

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Serge-Nook/topor/internal/core"
)

// HTMLHandler handles .html and .htm files.
type HTMLHandler struct{}

var metaAuthorRe = regexp.MustCompile(`(?i)<meta\s+name\s*=\s*"author"\s+content\s*=\s*"([^"]*)"`)
var metaAuthorRe2 = regexp.MustCompile(`(?i)<meta\s+content\s*=\s*"([^"]*)"\s+name\s*=\s*"author"`)

func (h *HTMLHandler) ReadMetadata(path string) (*MetadataInfo, error) {
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

	if m := metaAuthorRe.FindStringSubmatch(content); len(m) > 1 {
		info.Author = m[1]
	} else if m := metaAuthorRe2.FindStringSubmatch(content); len(m) > 1 {
		info.Author = m[1]
	}

	return info, nil
}

func (h *HTMLHandler) WriteMetadata(path string, author *core.AuthorSettings, timeSets *core.TimeSettings) error {
	if author != nil && author.Action != core.AuthorActionNone {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)

		switch author.Action {
		case core.AuthorActionDelete:
			content = metaAuthorRe.ReplaceAllString(content, "")
			content = metaAuthorRe2.ReplaceAllString(content, "")
		case core.AuthorActionReplace:
			newTag := fmt.Sprintf(`<meta name="author" content="%s"`, author.NewAuthor)
			if metaAuthorRe.MatchString(content) {
				content = metaAuthorRe.ReplaceAllString(content, newTag)
			} else if metaAuthorRe2.MatchString(content) {
				content = metaAuthorRe2.ReplaceAllString(content, newTag)
			} else {
				// Insert meta author tag in head
				headIdx := strings.Index(strings.ToLower(content), "<head>")
				if headIdx >= 0 {
					insertAt := headIdx + 6
					content = content[:insertAt] + "\n" + newTag + `">` + content[insertAt:]
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


