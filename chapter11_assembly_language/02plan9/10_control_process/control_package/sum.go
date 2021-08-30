package control_package

func Sum([]int64) int64

// 任意等差数列的和
func LoopAdd(cnt, v0, step int) int

/*
// 任意等差数列的和
	func LoopAdd(cnt, v0, step int) int {
		result := v0
		for i := 0; i < cnt; i++ {
			result += step
		}
		return result
	}
汇编思维
	func LoopAdd(cnt, v0, step int) int {
		var i = 0
		var result = 0

	LOOP_BEGIN:
		result = v0

	LOOP_IF:
		if i < cnt { goto LOOP_BODY }
		goto LOOP_END

	LOOP_BODY
		i = i+1
		result = result + step
		goto LOOP_IF

	LOOP_END:

		return result
	}
*/
