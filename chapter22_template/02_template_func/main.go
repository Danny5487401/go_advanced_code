package main

import (
	"os"
	"text/template"
)

// 模版函数

// 需求：给仓库不同材质物品数量加10个,然后进行打印

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Inventory struct {
	Material string // 材质
	Count    uint   // 数量
}
type Store struct {
	Fields []Inventory
}

func NewInventory(Fields []Inventory) Store {
	return Store{Fields: Fields}
}

func main() {
	sweaters := NewInventory([]Inventory{
		Inventory{Material: "wool", Count: 19},
		Inventory{Material: "wooltwo", Count: 20},
	})

	var Text = `
{{- range .Fields }}
   Material: {{.Material | handleString}} - Count:{{.Count | handleInt }}
{{- end }}
`
	tmpl, err := template.New("test").Funcs(template.FuncMap{"handleString": handleString, "handleInt": handleInt}).Parse(Text)
	CheckErr(err)
	err = tmpl.Execute(os.Stdout, sweaters)
	CheckErr(err)

}
func handleInt(number uint) uint {
	return number + 10
}
func handleString(field string) string {
	return "材质是: " + field
}
