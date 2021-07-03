package memento

// 定义编辑器的备忘录, 也就是编辑器的内部状态数据模型, 同时也对应一个历史版本
import "time"

type tEditorMemento struct {
	title      string
	content    string
	createTime int64
}

func newEditorMemento(title string, content string) *tEditorMemento {
	return &tEditorMemento{
		title, content, time.Now().Unix(),
	}
}
