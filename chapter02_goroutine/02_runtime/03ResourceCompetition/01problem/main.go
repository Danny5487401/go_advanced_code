package main

// 资源竞争
import (
	"fmt"
	"math/rand"
	"time"
)

//全局变量
var ticket = 10 // 假设10张票

func main() {
	/*
		4个goroutine，模拟4个售票口，4个子程序操作同一个共享数据。
	*/
	go saleTickets("售票口1") // g1,10
	go saleTickets("售票口2") // g2,10
	go saleTickets("售票口3") //g3,10
	go saleTickets("售票口4") //g4,10

	time.Sleep(5 * time.Second)
}

func saleTickets(name string) {
	rand.Seed(time.Now().UnixNano())
	//for i:=1;i<=100;i++{
	//	fmt.Println(name,"售出：",i)
	//}
	for { //ticket=1
		if ticket > 0 { //g1,g3,g2,g4
			//睡眠
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			// g1 ,g3, g2,g4
			fmt.Println(name, "售出：", ticket) // 1 , 0, -1 , -2
			ticket--                         //0 , -1 ,-2 , -3
		} else {
			fmt.Println(name, "售罄，没有票了。。")
			break
		}
	}
	fmt.Println("最终：", ticket)
}

/*
结果：
售票口1 售罄，没有票了。。
售票口3 售出： -1
售票口3 售罄，没有票了。。
售票口4 售出： -2
售票口4 售罄，没有票了。。

*/
