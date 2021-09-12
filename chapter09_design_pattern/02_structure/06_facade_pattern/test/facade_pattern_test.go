package test

import (
	facade "go_advanced_code/chapter09_design_pattern/02_structure/06_facade_pattern"
	"testing"
	"time"
)

func Test_FacadePattern(t *testing.T) {
	//用户id
	iUserID := 1
	//礼物id
	iGiftID := 2

	// 预先存入1000积分
	e := facade.MockPointsService.SaveUserPoints(iUserID, 1000)
	if e != nil {
		t.Error(e)
		return
	}

	// 预先存入1个库存
	e = facade.MockInventoryService.SaveStock(iGiftID, 1)
	if e != nil {
		t.Error(e)
		return
	}

	// 开始构建发送请求
	request := &facade.GiftExchangeRequest{
		ID:         1,
		UserID:     iUserID,
		GiftID:     iGiftID,
		CreateTime: time.Now().Unix(),
	}

	// 调用服务
	e, sOrderNo := facade.MockGiftExchangeService.Exchange(request)
	if e != nil {
		t.Log(e)
	}
	t.Logf("shipping order no = %v", sOrderNo)
}
