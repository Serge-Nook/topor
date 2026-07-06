package utils

import (
	"fmt"
	"io"
	"os"
)

// CreateBackup creates a .bak copy of the file.
func CreateBackup(path string) error {
	bakPath := path + ".bak"
	src, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open source: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(bakPath)
	if err != nil {
		return fmt.Errorf("cannot create backup: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("cannot copy: %w", err)
	}
	return nil
}
