package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	fio "github.com/bass3t/fiobenchgo/internal"
	fiofs "github.com/bass3t/fiobenchgo/internal/filesystem"
)

var params fio.BenchParams

func readSection(r io.ReaderAt, section fio.SectionInfo) (readed int64) {
	blkSize := fio.MinInt64(params.SecSize, params.BlockSize)
	blk := make([]byte, blkSize)

	var err error
	for readed < section.Size && err == nil {
		curOff := section.Offset + readed

		// read one block
		var sz int64
		for sz < blkSize && err == nil {
			var n int
			n, err = r.ReadAt(blk[n:], curOff+sz)
			sz += int64(n)
		}

		readed += int64(sz)
	}

	return
}

func readSections(r io.ReaderAt, input <-chan fio.SectionInfo) (readed int64) {
	for section := range input {
		readed += readSection(r, section)
	}
	return
}

func benchFile(path string) (readed int64, spend time.Duration) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	size := fiofs.FileSize(path)

	fiofs.DropCaches()

	sections := fio.SplitToSections(size, params.SecSize)

	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(params.SecWorkers)

	for i := 0; i < params.SecWorkers; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&readed, readSections(f, sections))
		}()
	}

	wg.Wait()

	endTime := time.Now()
	spend = endTime.Sub(startTime)
	return
}

func benchRead(path string) {
	if !fiofs.IsFileExist(path) {
		panic("file not found: " + path)
	}

	size := fiofs.FileSize(path)
	fmt.Printf("Input size: %d\n", size)

	for readers := 1; readers <= 8; readers++ {
		for _, secSize := range fio.BenchSectionSizes {
			for _, blockSize := range fio.BenchBlockSizes {
				if secSize.Size >= blockSize.Size {
					params = fio.BenchParams{
						SecSize:    secSize.Size,
						SecWorkers: readers,
						BlockSize:  blockSize.Size}

					readed, spend := benchFile(path)

					if readed != size {
						panic("failed read of file")
					}

					fmt.Printf("Readers: %d Section: %s Block: %s Speed: %s\n", readers, secSize.Info, blockSize.Info, fio.PrintSpeed(readed, spend))
				}
			}
		}
	}
}

func main() {
	benchRead(os.Args[1])
}
