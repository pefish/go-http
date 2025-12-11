package go_http

import (
	"bytes"
	"encoding/json"
	"io"

	"net/http"
	"strings"
	"time"

	i_logger "github.com/pefish/go-interface/i-logger"
	t_logger "github.com/pefish/go-interface/t-logger"
	"github.com/pkg/errors"
)

type HttpType struct {
	httpClient *http.Client
}

var HttpInstance = New(10 * time.Second)

func New(timeout time.Duration) *HttpType {
	return &HttpType{
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,              // 全局最大空闲连接
				MaxIdleConnsPerHost: 100,              // 每个主机最大空闲连接
				IdleConnTimeout:     90 * time.Second, // 空闲连接超时
				DisableKeepAlives:   false,            // 不要关闭 KeepAlive
			},
			Timeout: timeout,
		},
	}
}

type RequestParams struct {
	Url     string
	Queries map[string]string
	Params  any
	Headers map[string]any
}

func (t *HttpType) PostForStruct(
	logger i_logger.ILogger,
	params *RequestParams,
	struct_ any,
) (res_ *http.Response, bodyBytes_ []byte, err error) {
	bodyBytes, _ := json.Marshal(params.Params)
	url := params.Url
	if params.Queries != nil {
		url = "?" + mapToUrlQuery(params.Queries)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	req.Header.Set("Content-Type", "application/json")
	if logger.Level() == t_logger.Level_DEBUG {
		logger.DebugF("http post url: %s, body: %s\n", url, string(bodyBytes))
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if logger.Level() == t_logger.Level_DEBUG {
		logger.DebugF("http post resp url: %s, body: %s\n", url, string(respBytes))
	}

	if err := json.Unmarshal(respBytes, &struct_); err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	return resp, respBytes, nil
}

func mapToUrlQuery(paramsMap map[string]string) string {
	if paramsMap == nil {
		return ``
	}
	strParams := make([]string, 0)
	for key, value := range paramsMap {
		strParams = append(strParams, key+"="+value)
	}
	return strings.Join(strParams, "&")
}

func (t *HttpType) GetForStruct(
	logger i_logger.ILogger,
	params *RequestParams,
	struct_ any,
) (res_ *http.Response, bodyBytes_ []byte, err_ error) {
	res, bodyBytes, err := t.getForBytes(logger, params)
	if err != nil {
		return nil, nil, err
	}

	if err := json.Unmarshal(bodyBytes, &struct_); err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	return res, bodyBytes, nil
}

func (t *HttpType) GetForString(
	logger i_logger.ILogger,
	params *RequestParams,
) (res *http.Response, bodyStr string, err error) {
	res, bodyBytes, err := t.getForBytes(logger, params)
	if err != nil {
		return nil, "", err
	}
	return res, string(bodyBytes), nil
}

func (t *HttpType) getForBytes(
	logger i_logger.ILogger,
	params *RequestParams,
) (res_ *http.Response, bodyBytes_ []byte, err_ error) {
	url := params.Url
	if params.Queries != nil {
		url = "?" + mapToUrlQuery(params.Queries)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	req.Header.Set("Content-Type", "application/json")

	if logger.Level() == t_logger.Level_DEBUG {
		logger.DebugF("http get url: %s\n", url)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if logger.Level() == t_logger.Level_DEBUG {
		logger.DebugF("http get resp url: %s, body: %s\n", url, string(respBytes))
	}

	return resp, respBytes, nil
}
