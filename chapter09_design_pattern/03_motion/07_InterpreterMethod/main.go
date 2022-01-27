package main

import (
	"fmt"
	"strconv"
	"strings"
)

//解释器接口
type Node interface {
	Interpret() int //解释方法
}

//数据节点
type ValNode struct {
	val int // 存放1,2,3,+,-
}

func (vn *ValNode) Interpret() int {
	return vn.val
}

//=============加法节点=============
type AddNode struct {
	left, right Node
}

func (an *AddNode) Interpret() int {
	return an.left.Interpret() + an.right.Interpret()
}

//=============减法节点=============
type SubNode struct {
	left, right Node
}

func (an *SubNode) Interpret() int {
	return an.left.Interpret() - an.right.Interpret()
}

//=============解释对象=============
type Parser struct {
	exp   []string //表达式
	index int      //空格分割后的索引
	prev  Node     //前序节点
}

func (p *Parser) newValNode() Node {
	//执行数据操作
	v, _ := strconv.Atoi(p.exp[p.index])
	p.index++
	return &ValNode{val: v}
}
func (p *Parser) newAddNode() Node {
	//执行加法操作( + )
	p.index++
	return &AddNode{
		left:  p.prev,
		right: p.newValNode(),
	}
}
func (p *Parser) newSubNode() Node {
	//执行减法操作( - )
	p.index++
	return &SubNode{
		left:  p.prev,
		right: p.newValNode(),
	}
}
func (p *Parser) Result() Node { //返回结果
	return p.prev
}
func (p *Parser) Parse(exp string) { //对表达式进行解析
	p.exp = strings.Split(exp, " ") //通过空格分割
	for {
		if p.index >= len(p.exp) {
			return
		}
		switch p.exp[p.index] {
		case "+":
			p.prev = p.newAddNode()
		case "-":
			p.prev = p.newSubNode()
		default:
			p.prev = p.newValNode()

		}
	}
}

// 调用则
func main() {
	p := Parser{}
	p.Parse("1 + 2 + 30 - 4 + 10") //是通过空格进行解释的
	fmt.Println(p.Result().Interpret())
}
