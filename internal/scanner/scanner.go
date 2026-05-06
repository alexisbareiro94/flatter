package scanner

import (
	"os"
	"path/filepath"
	"sync"

	"flatter/internal/utils"
)

type Scanner struct{}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Scan(dirs []string) ([]string, error) {
	files := make([]string, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup
	fileChan := make(chan string, 1000)

	for _, dir := range dirs {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			walkDir(d, fileChan)
		}(dir)
	}

	go func() {
		wg.Wait()
		close(fileChan)
	}()

	for path := range fileChan {
		if utils.IsValidImage(path) {
			mu.Lock()
			files = append(files, path)
			mu.Unlock()
		}
	}

	return files, nil
}

func walkDir(dir string, fileChan chan string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			walkDir(path, fileChan)
		} else {
			fileChan <- path
		}
	}
}

func CountFiles(dirs []string) int {
	total := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, dir := range dirs {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			countDir(d, &mu, &total)
		}(dir)
	}

	wg.Wait()
	return total
}

func countDir(dir string, mu *sync.Mutex, counter *int) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			countDir(path, mu, counter)
		} else if utils.IsValidImage(entry.Name()) {
			mu.Lock()
			*counter++
			mu.Unlock()
		}
	}
}