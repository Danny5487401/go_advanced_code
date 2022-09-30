## 拷贝切片

拷贝切片可以用copy方法
```go
func copy(dst, src []Type) int
```
实际上copy根据数据类型，最终会调用切片的runtime.slicecopy方法。
```go
func slicecopy(toPtr unsafe.Pointer, toLen int, fmPtr unsafe.Pointer, fmLen int, width uintptr) int {
	//如果源切片和目标切片长度为0,则直接返回0
	if fmLen == 0 || toLen == 0 {
		return 0
	}

	//根据源切片和目标切片的长度，以长度最小的切片进行拷贝
	n := fmLen
	if toLen < n {
		n = toLen
	}

	if width == 0 {
		return n
	}

	//拷贝的空间大小=长度 * 元素大小
	size := uintptr(n) * width
	if size == 1 { // common case worth about 2x to do here
		// TODO: is this still worth it with new memmove impl?
		//如果拷贝的空间大小 等于1,那么直接转化赋值
		*(*byte)(toPtr) = *(*byte)(fmPtr) // known to be a byte pointer
	} else {
		//如果拷贝的空间大小 大于1,则源切片中array的数据拷贝到目标切片的array
		memmove(toPtr, fmPtr, size)
	}
	return n
}
```

从源码可以看出在切片拷贝的时候，要预先定义切片的长度再进行拷贝，否则有可能拷贝失败。