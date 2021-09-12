package facade

//模拟物流下单服务的接口

// 物流下单服务
type IShippingService interface {
	CreateShippingOrder(uid int, goodsID int) (error, string)
}
