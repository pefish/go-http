package go_http

import (
	"fmt"
	"github.com/pefish/go-test-assert"
	"testing"
)


func TestHttpClass_interfaceToUrlQuery(t *testing.T) {
	type Test struct {
		A string `json:"a"`
		B string `json:"b"`
		C uint64 `json:"c"`
	}
	result, _ := interfaceToUrlQuery(Test{
		A: `11`,
		B: `22`,
		C: 123,
	})
	test.Equal(t, 16, len(result))
}

func TestHttpClass_Post(t *testing.T) {
	str := `{"side": "test", "chain": "test"}`
	requester := NewHttpRequester()
	_, body, err := requester.Post(RequestParam{
		Url:    `http://www.baidu.com`,
		Params: str,
	})
	test.Equal(t, nil, err)
	test.Equal(t, true, len(body) > 0)
}

func TestHttpClass_Proxy(t *testing.T) {
	requester := NewHttpRequester()
	//requester := NewHttpRequester(WithHttpProxy("http://127.0.0.1:1087"))
	_, body, err := requester.Get(RequestParam{
		Url:    `http://ip.me`,
		Headers: map[string]interface{}{
			"User-Agent": "curl/7.64.1",
		},
	})
	test.Equal(t, nil, err)
	fmt.Println(body)
}
