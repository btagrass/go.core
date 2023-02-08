package typ

import (
	"os"
)

// 文件切片
type FileSlice []os.FileInfo

func (f FileSlice) Len() int {
	return len(f)
}

func (f FileSlice) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f FileSlice) Less(i, j int) bool {
	return f[i].ModTime().Before(f[j].ModTime())
}
