package service

import (
	"github.com/golang/mock/gomock"
	mock_dao "go_advanced_code/chapter28_test/02_mock/dao"
	"testing"
)

func TestFindUser(t *testing.T) {
	// 首先生成一个controller对象，然后注册到defer中
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	// 然后生成一个mockSearch对象，用来替代Search接口
	mockSearch := mock_dao.NewMockSearch(ctl)
	service := FindService{DB: mockSearch}
	// 操纵mockSearch对象的行为
	mockSearch.
		EXPECT().
		GetNameByID(int64(10)). //  参数
		Return("liangPin", nil) // 返回值
	mockSearch.EXPECT().GetNameByID(int64(20)).Return("zhangHeng", nil)

	// 测试过程
	if name1 := service.FindUser(10); name1 != "liangPin" {
		t.Error(name1)
	}

	if name2 := service.FindUser(20); name2 != "zhangHeng" {
		t.Error(name2)
	}
}
