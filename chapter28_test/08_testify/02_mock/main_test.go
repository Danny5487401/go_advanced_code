package main

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

// mock对象
type MockCrawler struct {
	mock.Mock
}

func (m *MockCrawler) GetUserList() ([]*User, error) {
	// Called()会返回一个mock.Arguments对象，该对象中保存着返回的值
	args := m.Called()
	return args.Get(0).([]*User), args.Error(1)
}

func TestGetUserList(t *testing.T) {
	crawler := new(MockCrawler)
	var (
		mockUsersInfo = []*User{{"dj", 18},
			{"zhangsan", 20}}
	)

	// 这里指示调用GetUserList()方法的返回值分别为MockUsers和nil，返回值在上面的GetUserList()方法中被Arguments.Get(0)和Arguments.Error(1)获取
	crawler.On("GetUserList").Return(mockUsersInfo, nil)

	GetAndPrintUsers(crawler)

	// crawler.AssertExpectations(t)对 Mock 对象做断言。
	crawler.AssertExpectations(t)
}
