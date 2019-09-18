package go_http

import (
	"fmt"
	"testing"
)

func TestHttpClass_GetForString(t *testing.T) {
	fmt.Println(Http.GetForString(RequestParam{
		Url: `http://www.baidu.com`,
		Params: map[string]interface{}{
			`haha`: 56,
			`test`: `hdh`,
		},
	}))
}

func TestHttpClass_interfaceToUrlQuery(t *testing.T) {
	type Test struct {
		A string `json:"a"`
		B string `json:"b"`
		C uint64 `json:"c"`
	}
	result := Http.interfaceToUrlQuery(Test{
		A: `11`,
		B: `22`,
		C: 123,
	})
	fmt.Printf(`%#v`, result)

}
