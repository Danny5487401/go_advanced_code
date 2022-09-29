package kind_route

import (
	"fmt"
	"reflect"
)

func CollectUserInfo(param interface{}) {
	val := reflect.ValueOf(param)
	switch val.Kind() {
	case reflect.String:
		fmt.Println("姓名:", val.String())
	case reflect.Struct:
		fmt.Println("姓名:", val.FieldByName("Name"))
		fmt.Println("年龄:", val.FieldByName("Age"))
		fmt.Println("住址:", val.FieldByName("Address"))
		fmt.Println("电话:", val.FieldByName("Phone"))
	case reflect.Ptr:
		fmt.Println("姓名:", val.Elem().FieldByName("Name"))
		fmt.Println("年龄:", val.Elem().FieldByName("Age"))
		fmt.Println("住址:", val.Elem().FieldByName("Address"))
		fmt.Println("电话:", val.Elem().FieldByName("Phone"))
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			fmt.Println("姓名:", val.Index(i).FieldByName("Name"))
			fmt.Println("年龄:", val.Index(i).FieldByName("Age"))
			fmt.Println("住址:", val.Index(i).FieldByName("Address"))
			fmt.Println("电话:", val.Index(i).FieldByName("Phone"))
		}
	case reflect.Map:
		itr := val.MapRange()
		for itr.Next() {
			fmt.Println("用户ID:", itr.Key())
			fmt.Println("姓名:", itr.Value().Elem().FieldByName("Name"))
			fmt.Println("年龄:", itr.Value().Elem().FieldByName("Age"))
			fmt.Println("住址:", itr.Value().Elem().FieldByName("Address"))
			fmt.Println("电话:", itr.Value().Elem().FieldByName("Phone"))
		}
	default:
		panic("unsupport type !!!")
	}
}

type User struct {
	Name    string
	Age     int64
	Address string
	Phone   int64
}
