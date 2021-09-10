package prototype

import (
	"encoding/json"
	"strings"
)

/*
UserFactory实现了创建UserInfo的简单工厂.创建的过程本质是调用了用户原型的Clone方法.用户原型是从json配置加载的, 便于按需修改配置.
*/

// DefaultUserFactory 用户工厂的全局单例
var DefaultUserFactory IUserFactory = newUserFactory()

// IUserFactory 工厂函数
type IUserFactory interface {
	Create() *UserInfo
}
type tUserFactory struct {
	defaultUserInfo *UserInfo
}

func (t *tUserFactory) Create() *UserInfo {
	return t.defaultUserInfo.Clone().(*UserInfo)
}

// 创建用户工厂实例
func newUserFactory() *tUserFactory {
	reader := strings.NewReader(loadUserConfig())
	decoder := json.NewDecoder(reader)
	user := newEmptyUser()
	err := decoder.Decode(user)
	if err != nil {
		panic(err)
	}

	return &tUserFactory{
		defaultUserInfo: user,
	}
}

// 加载默认用户的属性配置
func loadUserConfig() string {
	return `{
    "ID": 0,
    "Name" : "新用户",
    "RoleList" : ["guest"]
}`
}
