package filesystem

import "os"

// IsFileExist check is path exist and file
func IsFileExist(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		return !fi.IsDir()
	}

	return false
}

// FileSize return size of file or 0 in error
func FileSize(path string) int64 {
	if fi, err := os.Stat(path); err == nil {
		return fi.Size()
	}
	return 0
}
