package facade

//模拟实现物流下单服务

var MockShippingService = newMockShippingService()

type tMockShippingService struct {
}

func newMockShippingService() IShippingService {
	return &tMockShippingService{}
}

func (me *tMockShippingService) CreateShippingOrder(uid int, goodsID int) (error, string) {
	return nil, "shipping-order-666"
}
