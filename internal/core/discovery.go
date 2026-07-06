package core

import (
	"os"
	"path/filepath"
	"strings"
)

// SupportedExtensions lists all file extensions the app can handle.
var SupportedExtensions = map[string]bool{
	".docx": true, ".xlsx": true, ".docm": true, ".xlsm": true,
	".doc": true, ".xls": true,
	".odt": true,
	".html": true, ".htm": true,
	".rtf": true,
	".txt": true, ".xml": true,
}

// DiscoverFiles finds supported files from given paths.
func DiscoverFiles(paths []string, recursive bool, filterExts map[string]bool) []string {
	var result []string
	seen := make(map[string]bool)

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if info.IsDir() {
			files := walkDir(p, recursive, filterExts)
			for _, f := range files {
				if !seen[f] {
					seen[f] = true
					result = append(result, f)
				}
			}
		} else {
			ext := strings.ToLower(filepath.Ext(p))
			if SupportedExtensions[ext] && (filterExts == nil || filterExts[ext]) {
				abs, _ := filepath.Abs(p)
				if !seen[abs] {
					seen[abs] = true
					result = append(result, abs)
				}
			}
		}
	}
	return result
}

func walkDir(dir string, recursive bool, filterExts map[string]bool) []string {
	var result []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return result
	}
	for _, e := range entries {
		full := filepath.Join(dir, e.Name())
		if e.IsDir() {
			if recursive {
				result = append(result, walkDir(full, true, filterExts)...)
			}
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if SupportedExtensions[ext] && (filterExts == nil || filterExts[ext]) {
			abs, _ := filepath.Abs(full)
			result = append(result, abs)
		}
	}
	return result
}

// ParseMask parses a semicolon-separated mask string into a set of extensions.
func ParseMask(mask string) map[string]bool {
	if strings.TrimSpace(mask) == "" {
		return nil
	}
	result := make(map[string]bool)
	for _, part := range strings.Split(mask, ";") {
		part = strings.TrimSpace(part)
		part = strings.TrimLeft(part, "*")
		if part == "" {
			continue
		}
		if !strings.HasPrefix(part, ".") {
			part = "." + part
		}
		result[strings.ToLower(part)] = true
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
