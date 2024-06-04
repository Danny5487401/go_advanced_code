package gotest

// Add 方法用于演示go test使用
func Add(a int, b int) int {
	return a + b
}

func Sum(vals []int64) int64 {
	var total int64
	for _, val := range vals {
		if val%1e5 != 0 {

			total += val
		}
	}
	return total
}

func Reverse(s string) string {
	bs := []byte(s)
	length := len(bs)
	for i := 0; i < length/2; i++ {
		bs[i], bs[length-i-1] = bs[length-i-1], bs[i]
	}
	return string(bs)
}
