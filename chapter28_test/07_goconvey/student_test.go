package _7_goconvey

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestMin(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given some students with scores", t, func() {
		students := []Student{
			{
				ID:    1,
				Name:  "join",
				Score: 12,
			},
			{
				ID:    2,
				Name:  "michelle",
				Score: 13,
			},
			{
				ID:    3,
				Name:  "kelly",
				Score: 5,
			},
		}
		initialMin := GetMinimumScore(students)

		Convey("The result of GetMinimumScore should be the minimum score", func() {
			So(5, ShouldEqual, initialMin)
		})

		Convey("When a new score becomes the minimum", func() {
			students[0].Score = 3

			Convey("The minimum score should be updated", func() {
				newMin := GetMinimumScore(students)
				So(newMin, ShouldEqual, students[0].Score)
			})
		})
	})
}

func TestGetSumScore(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given some students with scores", t, func() {
		students := []Student{
			{
				ID:    1,
				Name:  "join",
				Score: 12,
			},
			{
				ID:    2,
				Name:  "michelle",
				Score: 13,
			},
			{
				ID:    3,
				Name:  "kelly",
				Score: 5,
			},
		}
		initialSum := GetSumScore(students)

		Convey("The result of GetSumScore should be equal to the sum of scores", func() {
			So(30, ShouldEqual, initialSum)
		})

		Convey("When any of the score is incremented", func() {
			students[0].Score++

			Convey("The sum should be greater by one", func() {
				newSum := GetSumScore(students)
				So(newSum, ShouldEqual, initialSum+1)
			})
		})

		Convey("When any of the score is decremented", func() {
			students[1].Score--

			Convey("The sum should be less by one", func() {
				newSum := GetSumScore(students)
				So(newSum, ShouldEqual, initialSum-1)
			})
		})
	})
}

func TestMain(m *testing.M) {
	// convey在TestMain场景下的入口
	SuppressConsoleStatistics()
	result := m.Run()
	// convey在TestMain场景下的结果打印
	PrintConsoleStatistics()
	os.Exit(result)
}
