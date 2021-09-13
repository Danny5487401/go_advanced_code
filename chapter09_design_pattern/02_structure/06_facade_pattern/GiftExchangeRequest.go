package facade

//积分兑换礼品请求
type GiftExchangeRequest struct {
	ID         int
	UserID     int
	GiftID     int
	CreateTime int64
}
