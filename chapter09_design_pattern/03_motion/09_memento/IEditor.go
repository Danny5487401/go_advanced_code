package memento

// 定义编辑器接口

type IEditor interface {
	Title(title string)     // 标题
	Content(content string) //内容
	Save()                  // 保存
	Undo() error            //后退
	Redo() error            // 前进

	Show() // 展示
}
