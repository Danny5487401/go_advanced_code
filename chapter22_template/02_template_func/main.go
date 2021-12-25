package main

import (
	"os"
	"text/template"
)

// 模版函数

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	type Inventory struct {
		Material string
		Count    uint
	}
	type NewInventory struct {
		Fields []Inventory
	}
	sweaters := NewInventory{
		Fields: []Inventory{
			Inventory{Material: "wool", Count: 19},
			Inventory{Material: "wooltwo", Count: 20},
		}}

	var Text = `
{{range .Fields }}
   Material: {{.Material | handleString}} - Count:{{.Count | handleInt }}
{{ end }}
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
	return " string is: " + field
}
