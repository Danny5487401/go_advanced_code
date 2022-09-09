package factorymethod

//PlusOperatorFactory 是 PlusOperator 的工厂类
type PlusOperatorFactory struct {
}

//PlusOperator Operator 的实际加法实现
type PlusOperator struct {
	*OperatorBase
}

func (PlusOperatorFactory) Create() Operator {
	return &PlusOperator{
		OperatorBase: &OperatorBase{},
	}
}

//Result 获取结果
func (o PlusOperator) Result() int {
	return o.a + o.b
}
