package main

import "fmt"

// 案例

// 1. 准备数据
type Employee struct {
	Name     string
	Age      int
	Vacation int
	Salary   int
}

var list = []Employee{
	{"Hao", 44, 0, 8000},
	{"Bob", 34, 10, 5000},
	{"Alice", 23, 5, 9000},
	{"Jack", 26, 0, 4000},
	{"Tom", 48, 9, 7500},
	{"Marry", 29, 0, 6000},
	{"Mike", 32, 8, 4000},
}

// 2. 设计函数
// 计数
func EmployeeCount(list []Employee, fn func(e *Employee) bool) int {
	count := 0
	for i, _ := range list {
		if fn(&list[i]) {
			count += 1
		}
	}
	return count
}

// 过滤
func EmployeeFilter(list []Employee, fn func(e *Employee) bool) []Employee {
	var newList []Employee
	for i, _ := range list {
		if fn(&list[i]) {
			newList = append(newList, list[i])
		}
	}
	return newList
}

// 求和
func EmployeeSum(list []Employee, fn func(e *Employee) int) int {
	var sum = 0
	for i, _ := range list {
		sum += fn(&list[i])
	}
	return sum
}


// 3. 业务逻辑
func main()  {
	// 1）统计有多少员工大于40岁
	old := EmployeeCount(list, func(e *Employee) bool {
		return e.Age > 40
	})

	fmt.Printf("old people: %d\n", old)

	// 2 .列出有没有休假的员工
	vacation := EmployeeFilter(list, func(e *Employee) bool {
		return e.Vacation != 0
	})
	fmt.Printf("People HAS vacation: %+v\n", vacation)

	//3. 统计30岁以下员工的薪资总和
	youngerPay := EmployeeSum(list, func(e *Employee) int {
		if e.Age < 30 {
			return e.Salary
		}
		return 0
	})
	fmt.Printf("People who larger than 30 has total money: %+v\n", youngerPay)
}
