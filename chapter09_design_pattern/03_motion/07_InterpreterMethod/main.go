/*
解释器模式（Interpreter）:针对特定问题设计的一种解决方案。例如，匹配字符串的时候，由于匹配条件非常灵活，使得通过代码来实现非常不灵活
举个例子，针对以下的匹配条件：
1.以+开头的数字表示的区号和电话号码，如+861012345678；
2.以英文开头，后接英文和数字，并以.分隔的域名，如www.baidu.com；
3. 以/开头的文件路径，如/path/to/file.txt；

优点:

1.可扩展性比较好，灵活。

2.增加了新的解释表达式的方式。

3.易于实现文法

缺点:

1. 执行效率比较低，可利用场景比较少。

2. 对于复杂的文法比较难维护
*/

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
