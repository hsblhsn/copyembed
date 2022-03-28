package copyembed

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyDirectory copies the source directory from embed.FS to the dstDir of os.
func CopyDirectory(em embed.FS, srcDir, dstDir string) error {
	entries, err := em.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(dstDir, entry.Name())

		file, err := em.Open(sourcePath)
		if err != nil {
			return err
		}

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExists(destPath, 0o755); err != nil {
				return err
			}
			if err := CopyDirectory(em, sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(em, sourcePath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// Copy copies a single file from embed.FS to the dstFile on os.
func Copy(em embed.FS, srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := em.Open(srcFile)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.CopyBuffer(out, in, make([]byte, 4096))
	if err != nil {
		return err
	}

	return nil
}

func exists(filePath string) bool {
	if _, err := os.Open(filePath); err != nil {
		return false
	}
	return true
}

func createIfNotExists(dir string, perm os.FileMode) error {
	if exists(dir) {
		return nil
	}
	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}
	return nil
}
