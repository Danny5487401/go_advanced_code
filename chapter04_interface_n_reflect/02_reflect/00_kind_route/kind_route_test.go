package kind_route

import "testing"

func TestCollectUserInfo(t *testing.T) {
	//string
	CollectUserInfo("张三")
	// 姓名: 张三

	//Struct
	CollectUserInfo(User{
		Name:    "张三",
		Age:     20,
		Address: "北京市海淀区",
		Phone:   1234567,
	})
	//姓名: 张三
	//年龄: 20
	//住址: 北京市海淀区
	//电话: 1234567

	//Ptr
	CollectUserInfo(&User{
		Name:    "张三",
		Age:     20,
		Address: "北京市海淀区",
		Phone:   1234567,
	})
	//姓名: 张三
	//年龄: 20
	//住址: 北京市海淀区
	//电话: 1234567

	//Slice
	CollectUserInfo([]User{
		{
			Name:    "张三",
			Age:     20,
			Address: "北京市海淀区",
			Phone:   1234567,
		},
		{
			Name:    "李四",
			Age:     30,
			Address: "北京市朝阳区",
			Phone:   7654321,
		},
	})
	//姓名: 张三
	//年龄: 20
	//住址: 北京市海淀区
	//电话: 1234567
	//姓名: 李四
	//年龄: 30
	//住址: 北京市朝阳区
	//电话: 7654321

	//Array
	CollectUserInfo([2]User{
		{
			Name:    "张三",
			Age:     20,
			Address: "北京市海淀区",
			Phone:   1234567,
		},
		{
			Name:    "李四",
			Age:     30,
			Address: "北京市朝阳区",
			Phone:   7654321,
		},
	})
	//姓名: 张三
	//年龄: 20
	//住址: 北京市海淀区
	//电话: 1234567
	//姓名: 李四
	//年龄: 30
	//住址: 北京市朝阳区
	//电话: 7654321

	CollectUserInfo(map[int]*User{
		1: {
			Name:    "张三",
			Age:     20,
			Address: "北京市海淀区",
			Phone:   1234567,
		},
		2: {
			Name:    "李四",
			Age:     30,
			Address: "北京市朝阳区",
			Phone:   7654321,
		},
	})
	//用户ID: 1
	//姓名: 张三
	//年龄: 20
	//住址: 北京市海淀区
	//电话: 1234567
	//用户ID: 2
	//姓名: 李四
	//年龄: 30
	//住址: 北京市朝阳区
	//电话: 7654321
}
