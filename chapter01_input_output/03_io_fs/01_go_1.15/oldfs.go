package oldFs

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func FindExtFileGo115(dir string, ext string) ([]byte, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
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
	return nil, err
}
