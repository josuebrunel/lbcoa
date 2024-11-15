package fizzbuzz

import (
	"context"
	"encoding/json"
	"fizzbuzz/pkg/assert"
	"net/http/httptest"
	"testing"
)

var testCases = []struct {
	name string
	path string
	code int
	err  string
	data string
}{
	{name: "WithMissingParams", path: "/", code: 400, err: ErrMissingParam.Error(), data: ""},
	{name: "WithMissingStrParam", path: "/?int1=3&int2=5&limit=20&str1=fizz", code: 400, err: ErrMissingParam.Error(), data: ""},
	{name: "WithInvalidIntParam", path: "/?int1=a&int2=5&limit=20&str1=fizz&str2=buzz", code: 400, err: ErrValMustBeInt.Error(), data: ""},
	{name: "WithInvalidLimitParam", path: "/?int1=3&int2=5&limit=-20&str1=fizz&str2=buzz", code: 400, err: ErrValMustBeGTZero.Error(), data: ""},
	{name: "WithValidParam", path: "/?int1=3&int2=5&limit=20&str1=fizz&str2=buzz", code: 200, err: "",
		data: "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,17,fizz,19,buzz"},
}

func TestFizzBuzzHanlder(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()
			t.Logf("Request: GET %v", tc.path)
			Handler(context.Background())(w, req)
			t.Logf("Response: %d %v", w.Code, w.Body.String())
			assert.Eq(t, w.Code, tc.code)
			resp := make(map[string]string)
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Eq(t, resp["error"], tc.err)
			assert.Eq(t, resp["data"], tc.data)
		})
	}
}
