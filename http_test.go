package go_http

import (
	"fmt"
	"testing"
)

func TestHttpClass_GetForString(t *testing.T) {
	fmt.Println(Http.GetForString(RequestParam{
		Url: `http://baidu.com`,
	}))
}
