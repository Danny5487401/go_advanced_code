package sortByReflect

import (
	"reflect"
	"sort"
)

// 通用排序
// 结构体排序，必须重写数组Len() Swap() Less()函数
type bodyWrapper struct {
	Body []interface{}
	by   func(p, q interface{}) bool // 内部Less()函数会用到
}

// 数组长度Len()
func (acw bodyWrapper) Len() int {
	return len(acw.Body)
}

// 元素交换
func (acw bodyWrapper) Swap(i, j int) {
	acw.Body[i], acw.Body[j] = acw.Body[j], acw.Body[i]
}

// 比较函数，使用外部传入的by比较函数
func (acw bodyWrapper) Less(i, j int) bool {
	return acw.by(acw.Body[i], acw.Body[j])
}

// 按照 field 子段类型排序
func SortBodyByIntOrString(bodys []interface{}, field string, sortBy string) {
	sort.Sort(bodyWrapper{bodys, func(p, q interface{}) bool {
		i := reflect.ValueOf(p).FieldByName(field)
		j := reflect.ValueOf(q).FieldByName(field)
		if i.Kind() == reflect.Int {
			if sortBy == "ASC" {
				return i.Int() < j.Int()
			} else {
				return i.Int() > j.Int()
			}
		}
		if i.Kind() == reflect.String {
			if sortBy == "ASC" {
				return i.String() < j.String()
			} else {
				return i.String() > j.String()
			}
		}
		return true
	}})
}
