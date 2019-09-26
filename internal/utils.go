package internal

import (
	"fmt"
	"time"
)

// SplitToSections return channel of SectionInfo for processing totalSize
func SplitToSections(totalSize, sectionSize int64) <-chan SectionInfo {
	out := make(chan SectionInfo)

	go func() {
		defer close(out)
		for offset := int64(0); offset < totalSize; offset += sectionSize {
			out <- SectionInfo{Offset: offset, Size: MinInt64(totalSize-offset, sectionSize)}
		}
	}()

	return out
}

// PrintSpeed output infrmation processing read/write speed
func PrintSpeed(count int64, d time.Duration) (str string) {
	count = count * 1000000
	speed := float64(count) / float64(d.Microseconds())
	speedMB := int(speed / (1024 * 1024))

	if speedMB > 0 {
		str = fmt.Sprintf("%d MB/s", speedMB)
	} else {
		speedKB := int(speed / 1024)
		if speedKB > 0 {
			str = fmt.Sprintf("%d KB/s", speedKB)
		} else {
			speedB := int(speed)
			str = fmt.Sprintf("%d B/s", speedB)
		}
	}
	return
}

// MinInt64 return minimum of a and b
func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
