package main

import (
	"fmt"
	"go_advanced_code/chapter17_dataStructure_n_algrithm/04_sort/sortByReflect"
	"sort"
)

/*
一。sort包内部实现了四种基本的排序算法

	1.插入排序insertionSort
	2.归并排序symMerge
	3.堆排序heapSort
	4.快速排序quickSort
	// 插入排序
	func insertionSort(data Interface, a, b int)
	// 堆排序
	func heapSort(data Interface, a, b int)
	// 快速排序
	func quickSort(data Interface, a, b, maxDepth int)
	// 归并排序
	func symMerge(data Interface, a, m, b int)
	sort包内置的四种排序方法是不公开的，只能被用于sort包内部使用。因此，对数据集合排序时，
	不必考虑应当选择哪一种，只需要实现sort.Interface接口定义三个接口即可
二. type Interface interface{
	  Len() int //返回集合中的元素个数
	  Less(i,j int) bool//i>j 返回索引i和元素是否比索引j的元素小
	  Swap(i,j int)//交换i和j的值
	}
	// 这里其实隐含要求这个容器或数据集合是slice类型或Array类型。否则，没法按照索引号取值
	//逆序
	sort包提供了Reverse()方法，允许将数据按Less()定义的排序方式逆序排序，而无需修改Less()代码。
三。Go的sort包已经为基本数据类型都实现了sort功能，其函数名的最后一个字母是s，表示sort之意。比如：Ints, Float64s, Strings，等等。
*/

func main() {
	// 1。不同结构体比较使用
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

	//2.相同结构体使用
	users := []User{
		{1, "alice", true, 18},
		{2, "bob", false, 29},
		{3, "carl", true, 12},
		{4, "john", true, 16},
	}
	sort.Slice(users, func(i, j int) bool {
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
