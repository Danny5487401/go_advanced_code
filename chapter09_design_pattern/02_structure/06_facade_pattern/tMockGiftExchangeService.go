package facade

//实现积分兑换礼品服务. 内部封装了积分服务, 库存服务和物流下单服务的调用.

import "errors"

type tMockGiftExchangeService struct {
}

func newMockGiftExchangeService() IGiftExchangeService {
	return &tMockGiftExchangeService{}
}

var MockGiftExchangeService = newMockGiftExchangeService()

// 模拟环境下未考虑事务提交和回滚
func (me *tMockGiftExchangeService) Exchange(request *GiftExchangeRequest) (error, string) {
	//1.查询礼物
	gift := MockInventoryService.GetGift(request.GiftID)
	if gift == nil {
		return errors.New("gift not found"), ""
	}

	//2。查询用户积分
	e, points := MockPointsService.GetUserPoints(request.UserID)
	if e != nil {
		return e, ""
	}
	if points < gift.Points {
		return errors.New("insufficient user points"), ""
	}

	//3.查询库存
	e, stock := MockInventoryService.GetStock(gift.ID)
	if e != nil {
		return e, ""
	}
	if stock <= 0 {
		return errors.New("insufficient gift stock"), ""
	}

	//4.扣减库存
	e = MockInventoryService.SaveStock(gift.ID, stock-1)
	if e != nil {
		return e, ""
	}
	//5。扣减积分
	e = MockPointsService.SaveUserPoints(request.UserID, points-gift.Points)
	if e != nil {
		return e, ""
	}

	//6.创建订单
	e, orderNo := MockShippingService.CreateShippingOrder(request.UserID, gift.ID)
	if e != nil {
		return e, ""
	}
	return nil, orderNo
}
