package go_http

import (
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
	requester.GetRequestClient().BounceToRawString = true
	_, body, err := requester.Post(RequestParam{
		Url:    `http://www.baidu.com`,
		Params: str,
	})
	test.Equal(t, nil, err)
	test.Equal(t, true, len(body) > 0)
}
