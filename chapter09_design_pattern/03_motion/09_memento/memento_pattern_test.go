package memento

import (
	"testing"
)

func Test_MementoPattern(t *testing.T) {
	editor := NewMockEditor()

	// test save()
	editor.Title("唐诗")
	editor.Content("白日依山尽")
	editor.Save()

	editor.Title("唐诗 登鹳雀楼")
	editor.Content("白日依山尽, 黄河入海流. ")
	editor.Save()

	editor.Title("唐诗 登鹳雀楼 王之涣")
	editor.Content("白日依山尽, 黄河入海流。欲穷千里目, 更上一层楼。")
	editor.Save()

	// test show()
	editor.Show()

	// test undo()
	for {
		e := editor.Undo()
		if e != nil {
			break
		} else {
			editor.Show()
		}
	}

	// test redo()
	for {
		e := editor.Redo()
		if e != nil {
			break
		} else {
			editor.Show()
		}
	}
}
