package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpperCaseHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	w := httptest.NewRecorder()
	UpperCaseHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "ABC" {
		t.Errorf("expected ABC got %v", string(data))
	}
}

func TestUpperCaseHandler1(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test api",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(UpperCaseHandle1))
			defer ts.Close()
			params := make(map[string]string)
			params["params"] = "paramsBody"
			paramsByte, _ := json.Marshal(params)
			resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(paramsByte))
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()

			t.Log(resp.StatusCode)
			if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				t.Error(string(body))
			}
		})
	}
}
