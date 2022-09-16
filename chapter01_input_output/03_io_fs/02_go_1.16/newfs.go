package newFs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// findTargetFile 查找dir目录下的所有文件，返回第一个文件名以ext为扩展名的文件内容
//
// 假设一定存在至少一个这样的文件
func FindExtFileGo116(dir string, ext string) ([]byte, error) {
	fsSys := os.DirFS(dir)                 // 以dir为根目录的文件系统，也就是说，后续所有的文件在这目录下
	entries, err := fs.ReadDir(fsSys, ".") // 读当前目录
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ext && !e.IsDir() {
			// 也可以一行代码返回
			// return fs.ReadFile(fsys, e.Name())
			f, err := fsSys.Open(e.Name()) // 文件名fullname `{dir}/{e.Name()}``
			if err != nil {
				return nil, err
			}
			defer f.Close()
			return io.ReadAll(f)
		}
	}
	panic("never happen")
}
