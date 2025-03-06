package jsonparser

import (
	"bytes"
	"errors"
	"unicode"
)

// Result is the type of the parser result
type Result struct {
	area  string
	part1 string
	part2 string
}

// 解析的总入口。输入电话号码，输出为Result的解析结果。
func Parse(input []byte) (Result, error) {
	l := newLex(input)
	_ = yyParse(l) // yyParse(yylex) 实际实现是 yyNewParser().Parse(yylex).
	return l.result, l.err
}

// 手动定义词法解析器
type lex struct {
	input  *bytes.Buffer
	result Result
	err    error
}

func newLex(input []byte) *lex {
	return &lex{
		input: bytes.NewBuffer(input),
	}
}

// 词法分析器满足Lex接口。能够逐个读取token。
func (l *lex) Lex(lval *yySymType) int {
	b := l.nextb()
	if unicode.IsDigit(rune(b)) { // 如果是有效的电话号码数字，返回D类型，并赋值lval.ch。
		lval.ch = b
		return D
	}
	return int(b)
}

func (l *lex) nextb() byte {
	b, err := l.input.ReadByte()
	if err != nil {
		return 0
	}
	return b
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
}

func cat(bytes ...byte) string {
	var out string
	for _, b := range bytes {
		out += string(b)
	}
	return out
}
