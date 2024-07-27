package gfile

import (
	"io"
	"os"
	"path/filepath"
	"sync"
)

func GetDirSize(path string) (int64, error) {
	size := int64(0)
	var wg sync.WaitGroup
	var mu sync.Mutex

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mu.Lock()
				size += info.Size()
				mu.Unlock()
			}()
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	wg.Wait()
	return size, nil
}

// CopyDir copies a directory from src to dst recursively.
func CopyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return CopyFile(path, destPath, info.Mode())
	})
}

// CopyFile copies a single file from src to dst and preserves file mode.
func CopyFile(src, dst string, mode os.FileMode) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return os.Chmod(dst, mode)
}
