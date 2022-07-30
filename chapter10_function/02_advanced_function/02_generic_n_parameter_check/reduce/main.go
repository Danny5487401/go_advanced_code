package reduce

import "reflect"

func Reduce(slice, pairFunc, zero interface{}) interface{} {
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("reduce: wrong type, not slice")
	}

	length := sliceInType.Len()
	if length == 0 {
		return zero
	} else if length == 1 {
		return sliceInType.Index(0)
	}

	elemType := sliceInType.Type().Elem()
	fn := reflect.ValueOf(pairFunc)
	if !verifyFuncSignature(fn, elemType, elemType, elemType) {
		t := elemType.String()
		panic("reduce: function must be of type func(" + t + ", " + t + ") " + t)
	}

	var ins [2]reflect.Value
	ins[0] = sliceInType.Index(0)
	ins[1] = sliceInType.Index(1)
	out := fn.Call(ins[:])[0]
	for i := 2; i < length; i++ {
		ins[0] = out
		ins[1] = sliceInType.Index(i)
		out = fn.Call(ins[:])[0]
	}
	return out.Interface()
}

// 检验是否是函数
func verifyFuncSignature(fn reflect.Value, types ...reflect.Type) bool {

	//Check it is a function
	if fn.Kind() != reflect.Func {
		return false
	}
	// 函数的入参数量和出参 数量 检查
	// NumIn() - returns a function type's input parameter count.
	// NumOut() - returns a function type's output parameter count.
	if (fn.Type().NumIn() != len(types)-1) || (fn.Type().NumOut() != 1) {
		return false
	}

	// In() - returns the type of a function type's i'th input parameter.
	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}
	// Out() - returns the type of a function type's i'th output parameter.
	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}
	return true
}
