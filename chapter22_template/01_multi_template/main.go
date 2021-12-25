package main

import (
	"os"
	"sort"
	"text/template"
)

// 多模板
var template1 = `func (m *{{.ModelName}}) HMapKey() string {
	return fmt.Sprintf("{{.TableName}}:{{.EntityDBID}}:%v", m.{{.EntityID}})
}`

var template2 = `
func (m *{{.ModelName}}) PK() (pk string){
	fmtStr:="{{.TableName}}
	{{- range $i, $v := .PKFields -}}
		:%v
	{{- end -}}
	"
	return fmt.Sprintf(fmtStr
		{{- range $i, $v := .PKFields -}}
			, m.{{$v}}
		{{- end -}}
		)
}
`

func main() {

	data1 := map[string]interface{}{
		"ModelName":  "A", //方法接受者名
		"TableName":  "t1",
		"EntityDBID": "id",
		"EntityID":   "ID",
	}
	tmpl, err := template.New("test1").Parse(template1)
	CheckErr(err)

	data2 := map[string]interface{}{
		"ModelName": "A", //方法接受者名
		"TableName": "t1",
		"PKFields":  sort.StringSlice{"ID", "SubID"},
	}
	tmpl, err = tmpl.New("test2").Parse(template2)
	CheckErr(err)

	// 选择模版
	err = tmpl.ExecuteTemplate(os.Stdout, "test1", data1)
	CheckErr(err)

	err = tmpl.ExecuteTemplate(os.Stdout, "test2", data2)
	CheckErr(err)

}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
