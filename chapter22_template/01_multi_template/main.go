package main

import (
	"os"
	"sort"
	"text/template"
)

// 多模板

// 需求：结构体方法

// 模版一：简单模版,读取文件
// 模版二：遍历模式
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
	tmpl, err := template.New("struct_method1").ParseFiles("chapter22_template/01_multi_template/struct_method1")
	CheckErr(err)

	data2 := map[string]interface{}{
		"ModelName": "A", //方法接受者名
		"TableName": "t1",
		"PKFields":  sort.StringSlice{"ID", "SubID"},
	}
	tmpl, err = tmpl.New("struct_method2").Parse(template2)
	CheckErr(err)

	// 选择并执行模版
	err = tmpl.ExecuteTemplate(os.Stdout, "struct_method1", data1)
	CheckErr(err)

	err = tmpl.ExecuteTemplate(os.Stdout, "struct_method2", data2)
	CheckErr(err)

}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
