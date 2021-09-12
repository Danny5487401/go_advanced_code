package facade

//模拟库存管理服务的接口

// 库存服务
type IInventoryService interface {
	GetGift(goodsID int) *GiftInfo
	GetStock(goodsID int) (error, int)
	SaveStock(goodsID int, num int) error
}
