package filesystem

import (
	"fmt"
	"os"
)

// DropCaches try flushing disk caches
func DropCaches() {
	f, err := os.OpenFile("/proc/sys/vm/drop_caches", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error clear caches")
		return
	}
	defer f.Close()

	f.Write([]byte("3"))
}
