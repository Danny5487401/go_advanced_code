//go:generate mockgen -package dao -destination dao_mock.go -source dao.go

package dao

//dao和service中具体的代码
// dao.dao.go
type Search interface {
	GetNameByID(id int64) (string, error)
}
