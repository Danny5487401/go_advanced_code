package oldFs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// findTargetFile 查找dir目录下的所有文件，返回第一个文件名以ext为扩展名的文件内容
//
// 假设一定存在至少一个这样的文件
func FindExtFileGo115(dir string, ext string) ([]byte, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		fmt.Println(filepath.Ext(e.Name()))
		if filepath.Ext(e.Name()) == ext && !e.IsDir() {
			name := filepath.Join(dir, e.Name())
			// 其实可以一行代码返回，这里只是为了展示更多的io操作
			// return ioutil.ReadFile(name)
			f, err := os.Open(name)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			return ioutil.ReadAll(f)
		}
	}
	panic("never happen")
}
