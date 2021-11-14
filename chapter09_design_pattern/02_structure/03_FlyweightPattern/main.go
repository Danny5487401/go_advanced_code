package main

import (
	"fmt"
	"sync"
)

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
