package go_http

import (
	"fmt"
	"testing"
)


func TestHttpClass_interfaceToUrlQuery(t *testing.T) {
	type Test struct {
		A string `json:"a"`
		B string `json:"b"`
		C uint64 `json:"c"`
	}
	result, _ := NewHttpRequester(WithIsDebug(true)).interfaceToUrlQuery(Test{
		A: `11`,
		B: `22`,
		C: 123,
	})
	fmt.Printf(`%#v`, result)

}

func TestHttpClass_Post(t *testing.T) {
	str := `{"side": "withdraw", "chain": "Pmeer"}`
	requester := NewHttpRequester(WithIsDebug(true))
	requester.RequestClient.BounceToRawString = true
	requester.Post(RequestParam{
		Url:    `http://www.baidu.com`,
		Params: str,
	})
	//fmt.Println(body)
}
