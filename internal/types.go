package internal

// SectionInfo define processing block of data
type SectionInfo struct {
	Offset int64
	Size   int64
}

// BenchParams define basic parameters for execution
type BenchParams struct {
	SecSize    int64
	SecWorkers int

	BlockSize int64
}

type SizeInfo struct {
	Size int64
	Info string
}

var BenchSectionSizes = []SizeInfo{
	{Size: 1 * 1024 * 1024, Info: "1 MB"},
	{Size: 4 * 1024 * 1024, Info: "4 MB"},
	{Size: 8 * 1024 * 1024, Info: "8 MB"},
	{Size: 16 * 1024 * 1024, Info: "16 MB"},
	{Size: 32 * 1024 * 1024, Info: "32 MB"},
	{Size: 64 * 1024 * 1024, Info: "64 MB"},
	{Size: 128 * 1024 * 1024, Info: "128 MB"},
}

var BenchBlockSizes = []SizeInfo{
	{Size: 512, Info: "512 B"},
	{Size: 1024, Info: "1 KB"},
	{Size: 4096, Info: "4 KB"},
	{Size: 16 * 1024, Info: "16 KB"},
	{Size: 64 * 1024, Info: "64 KB"},
	{Size: 1 * 1024 * 1024, Info: "1 MB"},
	{Size: 16 * 1024 * 1024, Info: "16 MB"},
	{Size: 32 * 1024 * 1024, Info: "32 MB"},
}
