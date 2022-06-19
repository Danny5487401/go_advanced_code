package service

import "go_advanced_code/chapter28_test/02_mock/dao"

// service.service.go
type FindService struct {
	DB dao.Search
}

func (f *FindService) FindUser(id int64) string {
	name, _ := f.DB.GetNameByID(id)
	return name
}
