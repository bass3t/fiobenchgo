package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	fio "github.com/bass3t/fiobenchgo/internal"
	bench "github.com/bass3t/fiobenchgo/internal/bench"
)

func benchWrite(ctx context.Context, wb *bench.Writer) {
	for writers := 1; writers <= 8; writers++ {
		for _, secSize := range fio.BenchSectionSizes {
			for _, blockSize := range fio.BenchBlockSizes {
				if secSize.Size >= blockSize.Size {
					params := fio.BenchParams{
						SecSize:    secSize.Size,
						SecWorkers: writers,
						BlockSize:  blockSize.Size}

					select {
					case <-ctx.Done():
						return
					default:
						written, spend, err := wb.Start(ctx, params)

						if err != nil {
							panic("failed write of file")
						}

						fmt.Printf("Writers: %d Section: %s Block: %s Speed: %s\n", writers, secSize.Info, blockSize.Info, fio.PrintSpeed(written, spend))
					}
				}
			}
		}
	}
}

func main() {
	path := "test.bin"
	var size int64 = 512 * 1024 * 1024
	fmt.Printf("Output size: %d\n", size)

	defer os.Remove(path)

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		cancel()
		done <- true
	}()

	benchWrite(ctx, bench.NewWriter(path, size))

	<-done
}
