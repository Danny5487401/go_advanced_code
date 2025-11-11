package main

import (
	"fmt"
	"github.com/spf13/cast"
	"time"
)

func main() {
	dealBasicType()

	dealTime()

	dealSlice()
}

func dealBasicType() {
	// ToString
	fmt.Println(cast.ToString("danny"))            // danny
	fmt.Println(cast.ToString(8))                  // 8
	fmt.Println(cast.ToString(8.31))               // 8.31
	fmt.Println(cast.ToString([]byte("one time"))) // one time
	fmt.Println(cast.ToString(nil))                // ""

	var foo interface{} = "one more time"
	fmt.Println(cast.ToString(foo)) // one more time

	// ToInt
	fmt.Println(cast.ToInt(8))     // 8
	fmt.Println(cast.ToInt(8.31))  // 8
	fmt.Println(cast.ToInt("8"))   // 8
	fmt.Println(cast.ToInt(true))  // 1
	fmt.Println(cast.ToInt(false)) // 0

	var eight interface{} = 8
	fmt.Println(cast.ToInt(eight)) // 8
	fmt.Println(cast.ToInt(nil))   // 0
}

func dealTime() {

	now := time.Now()
	timestamp := 1579615973
	timeStr := "2020-01-21 22:13:48"

	fmt.Println(cast.ToTime(now))       // 2020-01-22 06:31:50.5068465 +0800 CST m=+0.000997701
	fmt.Println(cast.ToTime(timestamp)) // 2020-01-21 22:12:53 +0800 CST
	fmt.Println(cast.ToTime(timeStr))   // 2020-01-21 22:13:48 +0000 UTC

	d, _ := time.ParseDuration("1m30s")
	ns := 30000
	strWithUnit := "130s"
	strWithoutUnit := "130"

	fmt.Println(cast.ToDuration(d))              // 1m30s
	fmt.Println(cast.ToDuration(ns))             // 30µs
	fmt.Println(cast.ToDuration(strWithUnit))    // 2m10s
	fmt.Println(cast.ToDuration(strWithoutUnit)) // 130ns
}

func dealSlice() {
	sliceOfInt := []int{1, 3, 7}
	arrayOfInt := [3]int{8, 12}
	// ToIntSlice
	fmt.Println(cast.ToIntSlice(sliceOfInt)) // [1 3 7]
	fmt.Println(cast.ToIntSlice(arrayOfInt)) // [8 12 0]

	sliceOfInterface := []interface{}{1, 2.0, "danny"}
	sliceOfString := []string{"abc", "dj", "pipi"}
	stringFields := " abc  def hij   "
	any := interface{}(37)
	// ToStringSliceE
	fmt.Println(cast.ToStringSlice(sliceOfInterface)) // [1 2 danny]
	fmt.Println(cast.ToStringSlice(sliceOfString))    // [abc dj pipi]
	fmt.Println(cast.ToStringSlice(stringFields))     // [abc def hij]
	fmt.Println(cast.ToStringSlice(any))              // [37]
}
