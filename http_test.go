package go_http

import (
	"fmt"
	go_logger "github.com/pefish/go-logger"
	"github.com/pefish/go-test-assert"
	"testing"
	"time"
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
	requester := NewHttpRequester(WithLogger(go_logger.Logger.CloneWithLevel("debug")))
	_, body, err := requester.PostForString(RequestParam{
		Url:    `http://www.baidu.com`,
		Params: str,
	})
	test.Equal(t, nil, err)
	test.Equal(t, true, len(body) > 0)
}

func TestHttpClass_Proxy(t *testing.T) {
	//var client http.Client
	//req, err := http.NewRequest("GET", "http://ip.me", nil)
	//test.Equal(t, nil, err)
	//req.Header.Add("User-Agent", "curl/7.64.1")
	//resp, err := client.Do(req)
	//test.Equal(t, nil, err)
	//b, err := ioutil.ReadAll(resp.Body)
	//test.Equal(t, nil, err)
	//fmt.Println(2, string(b))

	requester := NewHttpRequester()
	//requester := NewHttpRequester(WithHttpProxy("http://127.0.0.1:1087"))
	_, body, err := requester.GetForString(RequestParam{
		Url: `https://ip.me`,
		Headers: map[string]interface{}{
			"User-Agent": "curl/7.64.1",
		},
	})
	test.Equal(t, nil, err)
	fmt.Println(1, body)

}

func TestHttpClass_GetForStruct(t *testing.T) {
	var userInfo struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	}
	_, body, err := NewHttpRequester(WithTimeout(10*time.Second)).GetForStruct(RequestParam{
		Url: fmt.Sprintf("https://api.zealy.io/communities/%s/users", "baselendfinance"),
		Params: map[string]interface{}{
			"ethAddress": "0x8f72B6E3DF451F61b4E152de696B2Ba4748cFcc0",
		},
		Headers: map[string]interface{}{
			"x-api-key": "5ac42cpXHE3R9cq7l0Nm-ng78CA",
		},
		BasicAuth: nil,
	}, &userInfo)
	//test.Equal(t, nil, err)
	fmt.Printf("zealy users fetch error. body: %s - %#v\n", string(body), err)
	//if err != nil {
	//	fmt.Printf("zealy users fetch error - %#v\n", err)
	//	//fmt.Println()
	//}
}
