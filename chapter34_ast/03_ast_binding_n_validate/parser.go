package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

var codeTemplate = `
func getReqStruct(r *http.Request) (*{{requestStructName}}, error) {
	r.ParseForm()
	var reqStruct = &{{requestStructName}}{}

	// bind data
	{{bindData}}j

	// bind partial example
	// reqStruct.{{fieldName}} =
	// {{transToFieldType}}(r.Form['{{fieldTagFormName}}'])


	if bindErr != nil {
		return nil, err
	}

	// validate data
	{{validateData}}

	// validate partial example
	// validateErr = validate(reqStruct.{{fieldName}}, validateStr)
	// if validateErr != nil
	// return nil, err


	return reqStruct, nil
}
`

func getTag(input string) []structTag {
	var out []structTag
	var tagStr = input
	tagStr = strings.Replace(tagStr, "`", "", -1)
	tagStr = strings.Replace(tagStr, "\"", "", -1)
	tagList := strings.Split(tagStr, " ")
	for _, val := range tagList {
		tmpArr := strings.Split(val, ":")
		st := structTag{}
		st.key = tmpArr[0]
		st.values = strings.Split(tmpArr[1], ",")
		out = append(out, st)
	}

	return out
}

type structTag struct {
	key    string
	values []string
}

func main() {
	fset := token.NewFileSet()
	// if the src parameter is nil, then will auto read the second filepath file
	f, _ := parser.ParseFile(fset, "chapter34_ast/03_ast_binding_n_validate/example/example.go", nil, parser.Mode(0))
	//	ast.Print(fset, f.Decls[0])

	filedList := f.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
	for i, v := range filedList {
		fieldName := v.Names[0].Name
		fieldType := v.Type.(*ast.Ident).Name
		tagList := getTag(v.Tag.Value)
		fmt.Printf("第%v个字段-->%v:%v ,tag列表: %v  \n", i, fieldName, fieldType, tagList)
	}

	requestStructName := f.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name
	fmt.Println(requestStructName)
}
