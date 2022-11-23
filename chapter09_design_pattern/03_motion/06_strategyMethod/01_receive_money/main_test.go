package main

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	cs := new(CashContext)

	cs.SetCashContext("打八折")
	result := cs.GetMoney(100)
	fmt.Println("100元打八折结果:", result)

	cs.SetCashContext("满300返100")
	result = cs.GetMoney(400)
	fmt.Println("400元满200返100结果:", result)
}
