package prototype

type UserInfo struct {
	ID       int
	Name     string
	RoleList []string
}

func newEmptyUser() *UserInfo {
	return &UserInfo{}
}

// 克隆
func (u *UserInfo) Clone() ICloneable {
	roles := u.RoleList
	it := &UserInfo{
		u.ID, u.Name, make([]string, len(roles)),
	}

	for i, s := range roles {
		it.RoleList[i] = s
	}
	return it
}
