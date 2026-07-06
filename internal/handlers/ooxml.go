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

const corePropsPath = "docProps/core.xml"

const (
	nsDC      = "http://purl.org/dc/elements/1.1/"
	nsDCTerms = "http://purl.org/dc/terms/"
	nsCP      = "http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
	nsXSI     = "http://www.w3.org/2001/XMLSchema-instance"
)

// OOXMLHandler handles .docx, .xlsx, .docm, .xlsm files.
type OOXMLHandler struct{}

func (h *OOXMLHandler) ReadMetadata(path string) (*MetadataInfo, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open ZIP: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == corePropsPath {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}
			return parseCorePropsMeta(data)
		}
	}
	return &MetadataInfo{}, nil
}

func (h *OOXMLHandler) WriteMetadata(path string, author *core.AuthorSettings, timeSets *core.TimeSettings) error {
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

		if f.Name == corePropsPath {
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return err
			}
			modified := modifyCoreProps(content, author, timeSets)
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

func parseCorePropsMeta(data []byte) (*MetadataInfo, error) {
	info := &MetadataInfo{}
	d := xml.NewDecoder(bytes.NewReader(data))
	var inCreator, inLastModBy, inCreated, inModified bool

	for {
		tok, err := d.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Space == nsDC && t.Name.Local == "creator" {
				inCreator = true
			} else if t.Name.Space == nsCP && t.Name.Local == "lastModifiedBy" {
				inLastModBy = true
			} else if t.Name.Space == nsDCTerms && t.Name.Local == "created" {
				inCreated = true
			} else if t.Name.Space == nsDCTerms && t.Name.Local == "modified" {
				inModified = true
			}
		case xml.CharData:
			s := strings.TrimSpace(string(t))
			if inCreator {
				info.Author = s
				inCreator = false
			} else if inLastModBy {
				info.LastModBy = s
				inLastModBy = false
			} else if inCreated {
				if tm, err := time.Parse(time.RFC3339, s); err == nil {
					info.CreationTime = tm
				}
				inCreated = false
			} else if inModified {
				if tm, err := time.Parse(time.RFC3339, s); err == nil {
					info.ModifyTime = tm
				}
				inModified = false
			}
		case xml.EndElement:
			inCreator = false
			inLastModBy = false
			inCreated = false
			inModified = false
		}
	}
	return info, nil
}

func modifyCoreProps(data []byte, author *core.AuthorSettings, timeSets *core.TimeSettings) []byte {
	content := string(data)

	if author != nil && author.Action != core.AuthorActionNone {
		switch author.Action {
		case core.AuthorActionDelete:
			content = replaceXMLTag(content, nsDC, "creator", "")
			content = replaceXMLTag(content, nsCP, "lastModifiedBy", "")
		case core.AuthorActionReplace:
			if author.NewAuthor != "" {
				content = replaceXMLTag(content, nsDC, "creator", author.NewAuthor)
			}
			if author.NewLastModBy != "" {
				content = replaceXMLTag(content, nsCP, "lastModifiedBy", author.NewLastModBy)
			}
		}
	}

	if timeSets != nil {
		if timeSets.CreationMode == core.DateModeAbsolute {
			content = replaceXMLTag(content, nsDCTerms, "created", timeSets.CreationAbsolute.UTC().Format(time.RFC3339))
		}
		if timeSets.ModifyMode == core.DateModeAbsolute {
			content = replaceXMLTag(content, nsDCTerms, "modified", timeSets.ModifyAbsolute.UTC().Format(time.RFC3339))
		}
	}

	return []byte(content)
}

func replaceXMLTag(content, ns, localName, newValue string) string {
	// Find the tag with namespace prefix or full namespace
	// Try common patterns for OOXML core.xml
	prefixes := map[string][]string{
		nsDC:      {"dc:", ""},
		nsDCTerms: {"dcterms:", ""},
		nsCP:      {"cp:", ""},
	}

	for _, prefix := range prefixes[ns] {
		tagOpen := "<" + prefix + localName
		tagClose := "</" + prefix + localName + ">"

		startIdx := strings.Index(content, tagOpen)
		if startIdx < 0 {
			continue
		}
		// Find the end of the opening tag
		gtIdx := strings.Index(content[startIdx:], ">")
		if gtIdx < 0 {
			continue
		}
		contentStart := startIdx + gtIdx + 1
		endIdx := strings.Index(content[contentStart:], tagClose)
		if endIdx < 0 {
			continue
		}
		before := content[:contentStart]
		after := content[contentStart+endIdx:]
		content = before + newValue + after
		return content
	}
	return content
}
