package gcfile

import (
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
