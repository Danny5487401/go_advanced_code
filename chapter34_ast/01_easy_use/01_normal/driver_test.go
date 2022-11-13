package _1_normal

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	// 从第三方 RPC 获取一个司机的特征数据
	bs := getDriverRemote()
	var d Driver
	json.Unmarshal(bs, &d)
	fmt.Println(isOldDriver(&d))
}
