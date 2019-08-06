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
