package _3_f

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter28_test/01_testing/gotest"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
	"unicode/utf8"
)

func FuzzSum(f *testing.F) {
	f.Add(10) // 添加种子语料(seed corpus)数据，Fuzzing底层可以根据种子语料数据自动生成随机测试数据
	f.Fuzz(func(t *testing.T, n int) {
		n %= 20
		var vals []int64
		var expect int64
		var buf strings.Builder
		buf.WriteString("\n")
		for i := 0; i < n; i++ {
			val := rand.Int63() % 1e6
			vals = append(vals, val)
			expect += val
			buf.WriteString(fmt.Sprintf("%d,\n", val))
		}
		assert.Equal(t, expect, gotest.Sum(vals), buf.String())
	})
}

func FuzzReverse(f *testing.F) {
	str_slice := []string{"abc", "bb"}
	for _, v := range str_slice {
		f.Add(v)
	}
	f.Fuzz(func(t *testing.T, str string) {
		rev_str1 := gotest.Reverse(str)
		rev_str2 := gotest.Reverse(rev_str1)
		if str != rev_str2 {
			t.Errorf("fuzz test failed. str:%s, rev_str1:%s, rev_str2:%s", str, rev_str1, rev_str2)
		}
		if utf8.ValidString(str) && !utf8.ValidString(rev_str1) {
			t.Errorf("reverse result is not utf8. str:%s, len: %d, rev_str1:%s", str, len(str), rev_str1)
		}
	})
}
