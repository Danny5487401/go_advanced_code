package _7_goconvey

import "github.com/ahmetb/go-linq/v3"

type Student struct {
	ID    int64
	Name  string
	Age   int8
	Major string
	Score int32
}

// 返回这些学生的分数总和
func GetSumScore(students []Student) int64 {
	return linq.From(students).SelectT(
		func(s Student) int32 {
			return s.Score
		}).SumInts()
}

// 返回最低分
func GetMinimumScore(student []Student) int32 {
	return linq.From(student).SelectT(
		func(s Student) int32 {
			return s.Score
		}).Min().(int32)

}
