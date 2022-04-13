package main

import (
	"fmt"
	"go_advanced_code/chapter17_dataStructure_n_algorithm/04_sort/sortByReflect"
	"sort"
)

func main() {
	// 1。不同结构体切片排序.根据相同字段排序
	type fruitList []interface{}
	list := fruitList{
		Fruit{"p苹果", 2, 3.3},
		Fruit{"x香蕉", 8, 4.55},
		Fruit{"j橘子", 5, 2.5},
		Fruit{"c橙子", 3, 6.05},
		User{1, "alice", true, 18},
	}

	sortByReflect.SortBodyByInt(list, "Name", "ASC")
	fmt.Println(list)

	// 2. 相同结构体切片排序
	users := []User{
		{1, "alice", true, 18},
		{2, "bob", false, 29},
		{3, "carl", true, 12},
		{4, "john", true, 16},
	}
	sort.Slice(users, func(i, j int) bool {
		// age从小到大
		return users[i].Age < users[j].Age
	})
	fmt.Println(users)

	// 3.对字典的键名升序排序
	m := map[string]interface{}{"id": 1, "name": "admin", "pid": 0, "age": 18}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	//Go的sort包已经为基本数据类型都实现了sort功能，其函数名的最后一个字母是s，表示sort之意。比如：Ints, Float64s, Strings，等等。
	//对应的，int、float64、string等这些基本数据类型对应的集合类似，则被命名为IntSlice、Float64Slice、StringSlice等。
	//底层调用StringSlice结构体方法
	sort.Strings(keys)
	for i, v := range keys {
		fmt.Println(i, v)
	}

	// 4. 对sort包内置好的IntSlice、Float64Slice、StringSlice分别进行排序
	// 对自定义int类型数组以内置的IntSlice进行排序
	arr := []int{2, 1, 6, 5, 3}
	intList := sort.IntSlice(arr)
	sort.Sort(intList)
	// intList := sort.IntSlice(arr), sort.Sort(intList) 等价于 sort.Ints(arr)
	fmt.Println(intList)

	// 5. 查找数据
	arr1 := []int{2, 1, 6, 5, 3}
	// 错误使用直接使用，应该先排序
	sort.Ints(arr1)
	fmt.Println(arr1)
	foundPosition := sort.SearchInts(arr1, 6)
	unfoundPosition := sort.SearchInts(arr1, 4)
	fmt.Println("找到元素索引位置是", foundPosition)
	fmt.Println("找不到元素索引位置是", unfoundPosition)

	// 6. sort.Search用于从一个已经排序的数组中找到某个值所对应的索引
	arrInts := []int{13, 35, 56, 79}
	findPos := sort.Search(len(arrInts), func(i int) bool {
		// 注意不是 arrInts[i] == 35
		return arrInts[i] >= 35
	})
	// 等价与下面sort.SearchInts
	// findPos := sort.SearchInts(arrInts,35)
	fmt.Println(findPos)

}

type Fruit struct {
	Name  string
	Count int
	Price float64
}

type User struct {
	Id     int
	Name   string
	Status bool
	Age    int
}
