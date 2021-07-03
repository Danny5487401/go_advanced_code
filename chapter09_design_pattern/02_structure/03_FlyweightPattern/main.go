package main

import (
	"fmt"
	"sync"
)

/*
享元模式：
	尝试3重用现有的同类对象，如果未找到匹配的对象，则创建新的对象。
优点：
	在有大量对象时，会造成内存溢出。我们把其中共同的部分抽象出来，如果有相同的业务请求，直接返回内存已有的对象，避免重复创建，大大减少对象的创建
缺点
	提高了系统复杂度，需要分离出外部状态和内部状态。而且外部状态具有固有化性质，不应该
	随着内部状态的变化而变化。
在享元模式中，有以下几种角色：
	抽象享元角色（Flyweight）:是所有的具体享元类的基类，为具体享元规范需要实现的公共接口，非享元的外部状态以参数的形式通过方法传入。
	具体享元（Concrete Flyweight）角色：实现抽象享元角色中所规定的接口。
	非享元（Unsharable Flyweight)角色：是不可以共享的外部状态，它以参数的形式注入具体享元的相关方法中。
	享元工厂（Flyweight Factory）角色：负责创建和管理享元角色。当客户对象请求一个享元对象时，享元工厂检査系统中是否存在符合要求的享元对象，
		如果存在则提供给客户；如果不存在的话，则创建一个新的享元对象
举例：
	棋子中分为白子（具体享元）和黑子（具体享元），它们有自己独有的信息，在棋盘上会有许多棋子，它们的位置不同（非享元）。
	可以看到，虽然我们”下了“200个棋子，实际上，我们只用了两个棋子元素（享元），而位置（非共享信息）以参数的形式注入具体享元的相关方法中即下棋，
	在实际应用中，我们简单打印的白子和黑子的信息可能是许多内容，通过这种方式，节约了资源。每一次只需创建对象的非共享信息即可
备注：
	缓存，线程池等也是同等思路
何时使用：
	系统中有大量对象
	这些对象消耗大量内存
	这些对象的状态大部分可以外部化
	这些对象可以按照内蕴状态分为很多组，当把外蕴对象从对象中剔除出来时，每一组对象都可以用一个对象来代替
	系统不依赖于这些对象身份，这些对象是不可分辨的
	如何解决：用唯一标示判断，如果在内存中有，则返回这个唯一标示码的对象

关键代码：用Map存储这些对象
注意事项：
	注意划分外部状态和内部状态，否则可能会引起线程安全问题
	这些类必须有一个工厂对象加以控制

*/
const (
	BLACK ChessColor = 0
	WHITE ChessColor = 1
)

//非享元信息类
type PositionInfo struct {
	X int
	Y int
}

//享元接口类
type Chess interface {
	XiaQi(positionInfo *PositionInfo)
}

// 黑棋
type BlackChess struct {
	Info string
}

func (chess *BlackChess) XiaQi(info *PositionInfo) {
	fmt.Println("下黑子，位置：", info.X, ",", info.Y)
	fmt.Println("棋子的信息", chess.Info)
}

// 白棋
type WhiteChess struct {
	Info string
}

func (chess *WhiteChess) XiaQi(info *PositionInfo) {
	fmt.Println("下黑子，位置：", info.X, ",", info.Y)
	fmt.Println("棋子的信息", chess.Info)
}

type ChessColor uint8

//享元模式工厂类
type ChessFactory struct {
	m map[ChessColor]Chess
}

//棋子工厂，单例模式
var once sync.Once
var chessFactory *ChessFactory

func NewChessFactory() *ChessFactory {
	once.Do(func() {
		chessFactory = new(ChessFactory)
		chessFactory.m = make(map[ChessColor]Chess, 2)
		blackChess := new(BlackChess)
		blackChess.Info = "所有黑棋共有的信息"
		chessFactory.m[BLACK] = blackChess
		whiteCHess := new(WhiteChess)
		whiteCHess.Info = "所有白棋共有的信息"
		chessFactory.m[WHITE] = whiteCHess
	})
	return chessFactory
}

func (factory *ChessFactory) GetChess(color ChessColor) Chess {
	return factory.m[color]
}

// 开始调用
func main() {
	factory := NewChessFactory()
	blackChess := factory.GetChess(BLACK)
	whiteChess := factory.GetChess(WHITE)

	for i := 0; i < 100; i++ {
		positionInfo := new(PositionInfo)
		positionInfo.X, positionInfo.Y = i, i
		blackChess.XiaQi(positionInfo)
		whiteChess.XiaQi(positionInfo)
	}
}
