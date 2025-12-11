package go_http

import (
	"encoding/json"
	"fmt"
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

func TestHttpClass_PostFormDataForStruct(t *testing.T) {
	var httpResult struct {
		Data struct {
			Id            string `json:"id"`
			Name          string `json:"name"`
			Cid           string `json:"cid"`
			Size          uint64 `json:"size"`
			NumberOfFiles uint64 `json:"number_of_files"`
			MimeType      string `json:"mime_type"`
		} `json:"data"`
		Error struct {
			Code    uint64 `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	tokenInfoBytes, _ := json.Marshal(map[string]any{
		"name":        "test token",
		"description": "this is a test token1",
		"image":       "ipfs://bafkreie5g4y4u3p5u3z7k3z7k3z7k3z7k3z7k3z7k3z7k3z7k3z7k3z7k3z7k",
	})
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.w2i1UXZURz4dDY02S8mDLDtc_VbrOQpSb-R6tqGUcdY"

	_, _, err := HttpInstance.PostFormDataForStruct(
		&i_logger.DefaultLogger,
		&RequestParams{
			Url: "https://uploads.pinata.cloud/v3/files",
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			},
			Params: map[string]any{
				"network": "public",
				"file": FileInfoType{
					FileName:  "1.txt",
					FileBytes: tokenInfoBytes,
				},
			},
		},
		&httpResult,
	)
	go_test_.Equal(t, nil, err)
	spew.Dump(httpResult)
}
