# Map Key 类型

map中的key可以是任何的类型，只要它的值能比较是否相等，Go的语言规范已精确定义，Key的类型可以是：

- 布尔值
- 数字
- 字符串
- 指针
- 通道
- 接口类型
- 结构体
- 只包含上述类型的数组。

但不能是：

- slice
- map
- function

Key类型只要能支持==和!=操作符，即可以做为Key，当两个值==时，则认为是同一个Key。


## Go语言规范关于相等

- Pointer values are comparable. Two pointer values are equal if they point to the same variable or if both have value nil. Pointers to distinct zero-size variables may or may not be equal.当指针指向同一变量，或同为nil时指针相等，但指针指向不同的零值时可能不相等。
- Channel values are comparable. Two channel values are equal if they were created by the same call to make or if both have value nil.Channel当指向同一个make创建的或同为nil时才相等
- Interface values are comparable. Two interface values are equal if they have identical dynamic types and equal dynamic values or if both have value nil.从上面的例子我们可以看出，当接口有相同的动态类型并且有相同的动态值，或者值为都为nil时相等。要注意的是：参考理解Go Interface
- A value x of non-interface type X and a value t of interface type T are comparable when values of type X are comparable and X implements T. They are equal if t’s dynamic type is identical to X and t’s dynamic value is equal to x.如果一个是非接口类型X的变量x，也实现了接口T，与另一个接口T的变量t，只t的动态类型也是类型X，并且他们的动态值相同，则他们相等。
- Struct values are comparable if all their fields are comparable. Two struct values are equal if their corresponding non-blank fields are equal.结构体当所有字段的值相同，并且没有 相应的非空白字段时，则他们相等，
- Array values are comparable if values of the array element type are comparable. Two array values are equal if their corresponding elements are equal.两个数组只要他们包括的元素，每个元素的值相同，则他们相等。

注意：Go语言里是无法重载操作符的，struct是递归操作每个成员变量，struct也可以称为map的key，但如果struct的成员变量里有不能进行==操作的，例如slice，那么就不能作为map的key。