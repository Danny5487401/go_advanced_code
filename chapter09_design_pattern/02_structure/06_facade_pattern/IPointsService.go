package facade

// 模拟用户积分管理服务的接口

// 用户积分服务
type IPointsService interface {
	GetUserPoints(uid int) (error, int)
	SaveUserPoints(uid int, points int) error
}
