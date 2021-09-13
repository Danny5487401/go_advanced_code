package facade

// 积分兑换礼品的接口, 该接口是为方便客户端调用的Facade接口

// 礼品兑换服务
type IGiftExchangeService interface {
	// 兑换礼品, 并返回物流单号
	Exchange(request *GiftExchangeRequest) (error, string)
}
