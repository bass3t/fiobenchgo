package bench

import (
	"context"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	fio "github.com/bass3t/fiobenchgo/internal"
	fiofs "github.com/bass3t/fiobenchgo/internal/filesystem"
)

// Reader describe benchamarking process for read file
type Reader struct {
	path   string
	size   int64
	params fio.BenchParams
}

// NewReader return new instance of Reader
func NewReader(path string, size int64) *Reader {
	return &Reader{path: path, size: size}
}

func (br *Reader) readSection(ctx context.Context, r io.ReaderAt, section fio.SectionInfo) (readed int64) {
	blkSize := fio.MinInt64(br.params.SecSize, br.params.BlockSize)
	blk := make([]byte, blkSize)

	var err error
	for readed < section.Size && err == nil {
		curOff := section.Offset + readed

		// read one block
		var sz int64
		for sz < blkSize && err == nil {
			select {
			case <-ctx.Done():
				return
			default:
				var n int
				n, err = r.ReadAt(blk[n:], curOff+sz)
				sz += int64(n)
			}
		}

		readed += int64(sz)
	}

	return
}

func (br *Reader) readSections(ctx context.Context, r io.ReaderAt, input <-chan fio.SectionInfo) (readed int64) {
	for section := range input {
		readed += br.readSection(ctx, r, section)
	}
	return
}

// Start begin benchamrk process
func (br *Reader) Start(ctx context.Context, params fio.BenchParams) (readed int64, spend time.Duration, err error) {
	br.params = params

	var f *os.File
	if f, err = os.Open(br.path); err != nil {
		return
	}
	defer f.Close()

	fiofs.DropCaches()

	sections := fio.SplitToSections(br.size, br.params.SecSize)

	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(br.params.SecWorkers)

	for i := 0; i < br.params.SecWorkers; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&readed, br.readSections(ctx, f, sections))
		}()
	}

	wg.Wait()

	endTime := time.Now()
	spend = endTime.Sub(startTime)
	return
}
