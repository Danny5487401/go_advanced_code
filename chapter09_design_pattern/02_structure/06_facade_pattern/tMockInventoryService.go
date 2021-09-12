package facade

//模拟实现库存管理服务

var MockInventoryService = newMockInventoryService()

type tMockInventoryService struct {
	mGoodsStock map[int]int
}

func newMockInventoryService() IInventoryService {
	return &tMockInventoryService{
		make(map[int]int, 16),
	}
}

func (me *tMockInventoryService) GetGift(id int) *GiftInfo {
	return NewGiftInfo(id, "mock gift", 100)
}

func (me *tMockInventoryService) GetStock(goodsID int) (error, int) {
	n, ok := me.mGoodsStock[goodsID]
	if ok {
		return nil, n
	} else {
		return nil, 0
	}
}

func (me *tMockInventoryService) SaveStock(goodsID int, num int) error {
	me.mGoodsStock[goodsID] = num
	return nil
}
