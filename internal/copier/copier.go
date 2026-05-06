package copier

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"flatter/internal/progress"
	"flatter/internal/utils"
)

type Copier struct {
	destDir string
	mode    string
	workers int
	total   int32
	done    int32
}

func New(destDir string, mode string, workers int) *Copier {
	return &Copier{
		destDir: destDir,
		mode:    mode,
		workers: workers,
	}
}

func (c *Copier) SetTotal(total int) {
	atomic.StoreInt32(&c.total, int32(total))
}

func (c *Copier) Copy(files []string) (int, error) {
	atomic.StoreInt32(&c.done, 0)
	progressBar := progress.New(c.workers)
	total := int32(len(files))
	atomic.StoreInt32(&c.total, total)

	fileChan := make(chan string, len(files))
	resultChan := make(chan string, len(files))

	for i := 0; i < c.workers; i++ {
		go c.worker(fileChan, resultChan)
	}

	go func() {
		for _, f := range files {
			fileChan <- f
		}
		close(fileChan)
	}()

	copied := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			done := atomic.LoadInt32(&c.done)
			progressBar.Update(int(done), int(total))
		}
	}()

	for range resultChan {
		copied++
		if int32(copied) == total {
			break
		}
	}

	progressBar.Finish()
	return copied, nil
}

func (c *Copier) worker(fileChan <-chan string, resultChan chan<- string) {
	for srcPath := range fileChan {
		srcPath := srcPath

		src, err := os.Open(srcPath)
		if err != nil {
			resultChan <- ""
			continue
		}
		src.Close()

		filename := utils.GetFilename(srcPath)

		destPath, err := utils.GetDestPath(c.destDir, filename, c.mode)
		if err != nil {
			resultChan <- ""
			continue
		}

		if destPath == "" {
			atomic.AddInt32(&c.done, 1)
			resultChan <- ""
			continue
		}

		if err := c.copyFile(srcPath, destPath); err != nil {
			resultChan <- ""
			continue
		}

		atomic.AddInt32(&c.done, 1)
		resultChan <- destPath
	}
}

func (c *Copier) copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func (c *Copier) CopySimple(files []string) (int, error) {
	progressBar := progress.New(c.workers)
	total := len(files)

	fileChan := make(chan string, len(files))
	resultChan := make(chan string, len(files))

	var wg sync.WaitGroup
	var done int32

	for i := 0; i < c.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for srcPath := range fileChan {
				filename := utils.GetFilename(srcPath)

				destPath, err := utils.GetDestPath(c.destDir, filename, c.mode)
				if err != nil || destPath == "" {
					atomic.AddInt32(&done, 1)
					resultChan <- ""
					continue
				}

				if err := c.copyFile(srcPath, destPath); err != nil {
					atomic.AddInt32(&done, 1)
					resultChan <- ""
					continue
				}

				atomic.AddInt32(&done, 1)
				resultChan <- destPath
			}
		}()
	}

	go func() {
		for _, f := range files {
			fileChan <- f
		}
		close(fileChan)
	}()

	copied := 0
	lastUpdate := time.Now()

	updateTick := time.NewTicker(50 * time.Millisecond)
	defer updateTick.Stop()

	go func() {
		for range updateTick.C {
			d := atomic.LoadInt32(&done)
			progressBar.Update(int(d), total)
			if int(d) >= total {
				return
			}
		}
	}()

	for range resultChan {
		copied++
		if copied >= total {
			break
		}
		if time.Now().Sub(lastUpdate) > 100*time.Millisecond {
			progressBar.Update(copied, total)
			lastUpdate = time.Now()
		}
	}

	wg.Wait()
	progressBar.Finish()

	fmt.Printf("\nCopy result: %d/%d files copied\n", copied, total)
	return copied, nil
}

func (c *Copier) Progress() (int, int) {
	done := atomic.LoadInt32(&c.done)
	total := atomic.LoadInt32(&c.total)
	return int(done), int(total)
}