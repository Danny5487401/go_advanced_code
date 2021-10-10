package memento

import (
	"errors"
	"fmt"
)

// 发起人
// 虚拟的编辑器类, 实现IEditor接口
type tMockEditor struct {
	title    string
	content  string
	versions []*tEditorMemento
	index    int
}

func NewMockEditor() IEditor {
	return &tMockEditor{
		"", "", make([]*tEditorMemento, 0), 0,
	}
}

// 实现接口
func (me *tMockEditor) Title(title string) {
	me.title = title
}

func (me *tMockEditor) Content(content string) {
	me.content = content
}

func (me *tMockEditor) Save() {
	it := newEditorMemento(me.title, me.content)
	me.versions = append(me.versions, it)
	me.index = len(me.versions) - 1
}
func (me *tMockEditor) Undo() error {
	return me.load(me.index - 1)
}
func (me *tMockEditor) Redo() error {
	return me.load(me.index + 1)
}

func (me *tMockEditor) Show() {
	fmt.Printf("tMockEditor.Show, title=%s, content=%s\n", me.title, me.content)
}

// 加载配置
func (me *tMockEditor) load(i int) error {
	size := len(me.versions)
	if size <= 0 {
		return errors.New("no history versions")
	}

	if i < 0 || i >= size {
		return errors.New("no more history versions")
	}

	it := me.versions[i]
	me.title = it.title
	me.content = it.content
	me.index = i
	return nil
}
