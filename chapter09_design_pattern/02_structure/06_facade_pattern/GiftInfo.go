package facade

// 礼品信息实体
type GiftInfo struct {
	ID     int
	Name   string
	Points int
}

func NewGiftInfo(id int, name string, points int) *GiftInfo {
	return &GiftInfo{
		id, name, points,
	}
}
