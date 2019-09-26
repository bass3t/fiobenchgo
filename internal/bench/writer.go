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

// Writer describe benchamarking process for write file
type Writer struct {
	path   string
	size   int64
	params fio.BenchParams
}

// NewWriter return new instance of Writer
func NewWriter(path string, size int64) *Writer {
	return &Writer{path: path, size: size}
}

func (bw *Writer) writeSection(ctx context.Context, w io.WriterAt, section fio.SectionInfo) (written int64) {
	blkSize := fio.MinInt64(bw.params.SecSize, bw.params.BlockSize)
	blk := make([]byte, blkSize)

	var err error
	for written < section.Size && err == nil {
		curOff := section.Offset + written

		// write one block
		var sz int64
		for sz < blkSize && err == nil {
			select {
			case <-ctx.Done():
				return
			default:
				var n int
				n, err = w.WriteAt(blk[n:], curOff+sz)
				sz += int64(n)
			}
		}

		written += int64(sz)
	}

	return
}

func (bw *Writer) writeSections(ctx context.Context, w io.WriterAt, input <-chan fio.SectionInfo) (written int64) {
	for section := range input {
		written += bw.writeSection(ctx, w, section)
	}
	return
}

// Start begin benchamrk process
func (bw *Writer) Start(ctx context.Context, params fio.BenchParams) (written int64, spend time.Duration, err error) {
	bw.params = params

	var f *os.File
	if f, err = os.Create(bw.path); err != nil {
		return
	}
	defer f.Close()

	fiofs.DropCaches()

	sections := fio.SplitToSections(bw.size, bw.params.SecSize)

	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(bw.params.SecWorkers)

	for i := 0; i < bw.params.SecWorkers; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&written, bw.writeSections(ctx, f, sections))
		}()
	}

	wg.Wait()
	f.Sync()

	endTime := time.Now()
	spend = endTime.Sub(startTime)
	return
}
