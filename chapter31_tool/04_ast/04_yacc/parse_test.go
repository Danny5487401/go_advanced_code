package jsonparser

import (
	"testing"
)

func TestParser(t *testing.T) {
	// 美国手机号码的格式通常为(XXX) XXX-XXXX，前三位数字代表区号，其余七位数字代表具体电话号码
	// 功能实现: 支持电话号码4085551212、408-555-1212解析成3个部分
	testcases := []struct {
		input   string
		wantErr string
		output  Result
	}{{
		input:   "4085551212",
		wantErr: "",
		output:  Result{area: "408", part1: "555", part2: "1212"},
	}, {
		input:   "408-555-1212",
		wantErr: "",
		output:  Result{area: "408", part1: "555", part2: "1212"},
	}, {
		input:   "408-5551212",
		wantErr: "syntax error",
	}}
	for _, tc := range testcases {
		v, err := Parse([]byte(tc.input))
		var gotErr string
		if err != nil {
			gotErr = err.Error()
		}
		if gotErr != tc.wantErr {
			t.Errorf("%s err: %v, want %v", tc.input, gotErr, tc.wantErr)
		}
		if v != tc.output {
			t.Errorf("%s: %v, want %v", tc.input, v, tc.output)
		}
	}
}
