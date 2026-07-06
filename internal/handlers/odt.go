package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Serge-Nook/topor/internal/core"
)

const odtMetaPath = "meta.xml"

const (
	nsMeta   = "urn:oasis:names:tc:opendocument:xmlns:meta:1.0"
	nsOfficeMeta = "urn:oasis:names:tc:opendocument:xmlns:office:1.0"
	nsDCOdt  = "http://purl.org/dc/elements/1.1/"
)

// ODTHandler handles .odt files.
type ODTHandler struct{}

func (h *ODTHandler) ReadMetadata(path string) (*MetadataInfo, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open ZIP: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == odtMetaPath {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}
			return parseODTMeta(data)
		}
	}
	return &MetadataInfo{}, nil
}

func (h *ODTHandler) WriteMetadata(path string, author *core.AuthorSettings, timeSets *core.TimeSettings) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		if f.Name == odtMetaPath {
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return err
			}
			modified := modifyODTMeta(content, author, timeSets)
			fw, err := w.Create(f.Name)
			if err != nil {
				return err
			}
			if _, err := fw.Write(modified); err != nil {
				return err
			}
		} else {
			fw, err := w.CreateHeader(&f.FileHeader)
			if err != nil {
				rc.Close()
				return err
			}
			if _, err := io.Copy(fw, rc); err != nil {
				rc.Close()
				return err
			}
			rc.Close()
		}
	}

	if err := w.Close(); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}

func parseODTMeta(data []byte) (*MetadataInfo, error) {
	info := &MetadataInfo{}
	d := xml.NewDecoder(bytes.NewReader(data))

	var currentTag string
	for {
		tok, err := d.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "initial-creator" || (t.Name.Space == nsDCOdt && t.Name.Local == "creator") {
				currentTag = t.Name.Local
			} else if t.Name.Local == "creation-date" {
				currentTag = "creation-date"
			} else if t.Name.Local == "date" && t.Name.Space == nsDCOdt {
				currentTag = "date"
			}
		case xml.CharData:
			s := strings.TrimSpace(string(t))
			switch currentTag {
			case "initial-creator", "creator":
				if info.Author == "" {
					info.Author = s
				} else {
					info.LastModBy = s
				}
			case "creation-date":
				if tm, err := time.Parse(time.RFC3339, s); err == nil {
					info.CreationTime = tm
				} else if tm, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
					info.CreationTime = tm
				}
			case "date":
				if tm, err := time.Parse(time.RFC3339, s); err == nil {
					info.ModifyTime = tm
				} else if tm, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
					info.ModifyTime = tm
				}
			}
			currentTag = ""
		case xml.EndElement:
			currentTag = ""
		}
	}
	return info, nil
}

func modifyODTMeta(data []byte, author *core.AuthorSettings, timeSets *core.TimeSettings) []byte {
	content := string(data)

	if author != nil && author.Action != core.AuthorActionNone {
		switch author.Action {
		case core.AuthorActionDelete:
			content = replaceSimpleTag(content, "meta:initial-creator", "")
			content = replaceSimpleTag(content, "dc:creator", "")
		case core.AuthorActionReplace:
			if author.NewAuthor != "" {
				content = replaceSimpleTag(content, "meta:initial-creator", author.NewAuthor)
				content = replaceSimpleTag(content, "dc:creator", author.NewAuthor)
			}
		}
	}

	if timeSets != nil {
		if timeSets.CreationMode == core.DateModeAbsolute {
			content = replaceSimpleTag(content, "meta:creation-date", timeSets.CreationAbsolute.UTC().Format(time.RFC3339))
		}
		if timeSets.ModifyMode == core.DateModeAbsolute {
			content = replaceSimpleTag(content, "dc:date", timeSets.ModifyAbsolute.UTC().Format(time.RFC3339))
		}
	}

	return []byte(content)
}

func replaceSimpleTag(content, tagName, newValue string) string {
	open := "<" + tagName
	close := "</" + tagName + ">"
	startIdx := strings.Index(content, open)
	if startIdx < 0 {
		return content
	}
	gtIdx := strings.Index(content[startIdx:], ">")
	if gtIdx < 0 {
		return content
	}
	contentStart := startIdx + gtIdx + 1
	endIdx := strings.Index(content[contentStart:], close)
	if endIdx < 0 {
		return content
	}
	return content[:contentStart] + newValue + content[contentStart+endIdx:]
}
