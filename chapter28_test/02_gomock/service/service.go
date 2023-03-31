package service

import "github.com/Danny5487401/go_advanced_code/chapter28_test/02_gomock/dao"

// service.service.go
type FindService struct {
	DB dao.Search
}

func (f *FindService) FindUser(id int64) string {
	name, _ := f.DB.GetNameByID(id)
	return name
}
