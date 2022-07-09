/*
策略模式
背景：商场打折是CRM系统中最常见的情况，商场随时会增加或者删除促销方式，而且促销方式非常多，我们需要针对各种情况做出不同的结果，
	例如打八折，我们就要根据实际价格计算出打八折的结果，满300返100等等。这个时候我们怎么做呢，大多数人没有接触设计模式时，
	都会想，只要用一个switch case语句就可以，其实这样确实可以，也能达到想要的结果，但是如果规则变得复杂了，增加和很多或者删除了很多，
	这样我们就很难维护了，如果有一百种促销方式，难道我们要在一百个case语句里写上各个算法吗?

*/

package main

import "fmt"

// 收银员
type CashSuper interface {
	AcceptCash(money float64) float64
}

// 正常收费
type Normal struct{}

func (n *Normal) AcceptCash(money float64) float64 {
	return money
}

// 打折
type CashRebate struct {
	Rebate float64
}

func (c *CashRebate) SetRebate(rebate float64) {
	c.Rebate = rebate
}

func (c *CashRebate) AcceptCash(money float64) float64 {
	return c.Rebate * money
}

// 满 x 返 y
type CashReturn struct {
	MoneyCondition float64
	MoneyReturn    float64
}

func (c *CashReturn) SetCashReturn(moneyCondition float64, moneyReturn float64) {
	c.MoneyCondition = moneyCondition
	c.MoneyReturn = moneyReturn
}
func (c *CashReturn) AcceptCash(money float64) float64 {
	if money >= c.MoneyCondition {
		moneyMinus := int(money / c.MoneyCondition)
		return money - float64(moneyMinus)*c.MoneyReturn
	}
	return money
}

type CashContext struct {
	Strategy CashSuper
}

func (c *CashContext) SetCashContext(t string) {
	switch t {
	case "正常收费":
		normal := new(Normal)
		c.Strategy = normal
	case "打八折":
		r := new(CashRebate)
		r.SetRebate(0.8)
		c.Strategy = r
	case "满300返100":
		r := new(CashReturn)
		r.SetCashReturn(200, 100)
		c.Strategy = r
	}
}

// 收钱
func (c *CashContext) GetMoney(money float64) float64 {
	return c.Strategy.AcceptCash(money)
}

func main() {
	cs := new(CashContext)

	cs.SetCashContext("打八折")
	result := cs.GetMoney(100)
	fmt.Println("100元打八折结果:", result)

	cs.SetCashContext("满300返100")
	result = cs.GetMoney(400)
	fmt.Println("400元满200返100结果:", result)
}
