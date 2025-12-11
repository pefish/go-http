package go_http

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	i_logger "github.com/pefish/go-interface/i-logger"
	go_test_ "github.com/pefish/go-test"
)

func TestHttpClass_GetForStruct(t *testing.T) {
	_, ip, err := HttpInstance.GetForString(&i_logger.DefaultLogger, &RequestParams{
		Url: "http://ifconfig.io",
	})
	go_test_.Equal(t, nil, err)
	spew.Dump(ip)
}

func TestHttpClass_PostForStruct(t *testing.T) {
	var httpResult struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}
	_, _, err := HttpInstance.PostForStruct(
		&i_logger.DefaultLogger,
		&RequestParams{
			Url: "http://tyo-sender.helius-rpc.com/fast",
			Params: map[string]any{
				"id":      "gsgfsf",
				"jsonrpc": "2.0",
				"method":  "sendTransaction",
				"params": []any{
					"gsfdgs",
					map[string]any{
						"encoding":      "base64",
						"skipPreflight": true,
						"maxRetries":    0,
					},
				},
			},
		},
		&httpResult,
	)
	go_test_.Equal(t, nil, err)
	spew.Dump(httpResult)
}
