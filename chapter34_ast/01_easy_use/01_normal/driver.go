package _1_normal

type Driver struct {
	Orders       int
	DrivingYears int
}

// 从第三方获取司机特征，json 表示
func getDriverRemote() []byte {
	return []byte(`{"orders":100000,"driving_years":18}`)
}

// 判断是否为老司机
func isOldDriver(d *Driver) bool {
	if d.Orders > 10000 && d.DrivingYears > 5 {
		return true
	}
	return false
}
