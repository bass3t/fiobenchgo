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

func writeSection(w io.WriterAt, section fio.SectionInfo) (written int64) {
	blkSize := fio.MinInt64(params.SecSize, params.BlockSize)
	blk := make([]byte, blkSize)

	var err error
	for written < section.Size && err == nil {
		curOff := section.Offset + written

		// write one block
		var sz int64
		for sz < blkSize && err == nil {
			var n int
			n, err = w.WriteAt(blk[n:], curOff+sz)
			sz += int64(n)
		}

		written += int64(sz)
	}

	return
}

func writeSections(w io.WriterAt, input <-chan fio.SectionInfo) (written int64) {
	for section := range input {
		written += writeSection(w, section)
	}
	return
}

func benchFile(path string, size int64) (written int64, spend time.Duration) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fiofs.DropCaches()

	sections := fio.SplitToSections(size, params.SecSize)

	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(params.SecWorkers)

	for i := 0; i < params.SecWorkers; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&written, writeSections(f, sections))
		}()
	}

	wg.Wait()
	f.Sync()

	endTime := time.Now()
	spend = endTime.Sub(startTime)
	return
}

func benchWrite(path string, size int64) {
	fmt.Printf("Output size: %d\n", size)

	for writers := 1; writers <= 8; writers++ {
		for _, secSize := range fio.BenchSectionSizes {
			for _, blockSize := range fio.BenchBlockSizes {
				if secSize.Size >= blockSize.Size {
					params = fio.BenchParams{
						SecSize:    secSize.Size,
						SecWorkers: writers,
						BlockSize:  blockSize.Size}

					written, spend := benchFile(path, size)

					if written != size {
						panic("failed write of file")
					}

					fmt.Printf("Writers: %d Section: %s Block: %s Speed: %s\n", writers, secSize.Info, blockSize.Info, fio.PrintSpeed(written, spend))
				}
			}
		}
	}
}

func main() {
	path := "test.bin"

	var size int64 = 512 * 1024 * 1024
	benchWrite(path, size)
}
