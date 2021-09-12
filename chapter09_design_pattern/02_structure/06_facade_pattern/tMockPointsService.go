package facade

//模拟实现用户积分管理服务
import "errors"

var MockPointsService = newMockPointsService()

type tMockPointsService struct {
	mUserPoints map[int]int
}

func newMockPointsService() IPointsService {
	return &tMockPointsService{
		make(map[int]int, 16),
	}
}

func (me *tMockPointsService) GetUserPoints(uid int) (error, int) {
	n, ok := me.mUserPoints[uid]
	if ok {
		return nil, n
	} else {
		return errors.New("user not found"), 0
	}
}

func (me *tMockPointsService) SaveUserPoints(uid int, points int) error {
	me.mUserPoints[uid] = points
	return nil
}
