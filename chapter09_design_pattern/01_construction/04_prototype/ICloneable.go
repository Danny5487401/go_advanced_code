package prototype

// 定义克隆接口
type ICloneable interface {
	Clone() ICloneable
}
