package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	fio "github.com/bass3t/fiobenchgo/internal"
	bench "github.com/bass3t/fiobenchgo/internal/bench"
	fiofs "github.com/bass3t/fiobenchgo/internal/filesystem"
)

func benchRead(ctx context.Context, rb *bench.Reader) {
	for readers := 1; readers <= 8; readers++ {
		for _, secSize := range fio.BenchSectionSizes {
			for _, blockSize := range fio.BenchBlockSizes {
				if secSize.Size >= blockSize.Size {
					params := fio.BenchParams{
						SecSize:    secSize.Size,
						SecWorkers: readers,
						BlockSize:  blockSize.Size}

					select {
					case <-ctx.Done():
						return
					default:
						readed, spend, err := rb.Start(ctx, params)

						if err != nil {
							panic("failed read of file")
						}

						fmt.Printf("Readers: %d Section: %s Block: %s Speed: %s\n", readers, secSize.Info, blockSize.Info, fio.PrintSpeed(readed, spend))
					}
				}
			}
		}
	}
}

func main() {
	path := os.Args[1]
	if !fiofs.IsFileExist(path) {
		panic("file not found: " + path)
	}

	size := fiofs.FileSize(path)
	fmt.Printf("Input size: %d\n", size)

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		cancel()
		done <- true
	}()

	benchRead(ctx, bench.NewReader(path, size))

	<-done
}
