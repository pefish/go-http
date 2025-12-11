package go_http

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"

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
	Headers map[string]string
}

func (t *HttpType) PostForStruct(
	logger i_logger.ILogger,
	params *RequestParams,
	struct_ any,
) (res_ *http.Response, bodyBytes_ []byte, err error) {
	res, bodyBytes, err := t.Post(logger, params)
	if err != nil {
		return nil, nil, err
	}

	if err := json.Unmarshal(bodyBytes, &struct_); err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	return res, bodyBytes, nil
}

func (t *HttpType) Post(
	logger i_logger.ILogger,
	params *RequestParams,
) (res_ *http.Response, bodyBytes_ []byte, err error) {
	var bodyBytes []byte

	switch params.Params.(type) {
	case string:
		bodyBytes = []byte(params.Params.(string))
	case []byte:
		bodyBytes = params.Params.([]byte)
	default:
		bodyBytes, _ = json.Marshal(params.Params)
	}

	url := params.Url
	if params.Queries != nil {
		url = "?" + mapToUrlQuery(params.Queries)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	if params.Headers == nil {
		params.Headers = make(map[string]string)
	}
	params.Headers["Content-Type"] = "application/json"
	for headerKey, headerValue := range params.Headers {
		req.Header.Set(headerKey, headerValue)
	}
	if logger.Level() == t_logger.Level_DEBUG {
		logger.DebugF(
			`[HTTP POST] 
Url: %s
Headers: %v
Body: %s

`,
			url,
			params.Headers,
			string(bodyBytes),
		)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if logger.Level() == t_logger.Level_DEBUG {
		logger.DebugF(
			`[HTTP POST result] 
Url: %s
Body: %s

`,
			url,
			string(respBytes),
		)
	}

	return resp, respBytes, nil
}

type FileInfoType struct {
	FileName  string
	FileBytes []byte
}

func (t *HttpType) PostFormDataForStruct(
	logger i_logger.ILogger,
	params *RequestParams,
	struct_ any,
) (res_ *http.Response, bodyBytes_ []byte, err error) {
	res, respBytes, err := t.PostFormData(logger, params)
	if err != nil {
		return nil, nil, err
	}

	if err := json.Unmarshal(respBytes, &struct_); err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	return res, respBytes, nil
}

func (t *HttpType) PostFormData(
	logger i_logger.ILogger,
	params *RequestParams,
) (res_ *http.Response, bodyBytes_ []byte, err error) {
	url := params.Url
	if params.Queries != nil {
		url = "?" + mapToUrlQuery(params.Queries)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fields, ok := params.Params.(map[string]any)
	if !ok {
		return nil, nil, errors.New("params.Params is not map[string]any")
	}

	for fieldName, value := range fields {
		if fileInfo, ok := value.(FileInfoType); ok {
			// 如果是文件信息，则创建文件字段
			part, err := writer.CreateFormFile(fieldName, fileInfo.FileName)
			if err != nil {
				return nil, nil, errors.Wrap(err, "")
			}
			_, err = part.Write(fileInfo.FileBytes)
			if err != nil {
				return nil, nil, errors.Wrap(err, "")
			}
		} else {
			// 否则创建普通字段
			fieldValue, ok := value.(string)
			if !ok {
				return nil, nil, errors.New("form field value is not string")
			}
			if err := writer.WriteField(fieldName, fieldValue); err != nil {
				return nil, nil, errors.Wrap(err, "")
			}
		}
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	if params.Headers == nil {
		params.Headers = make(map[string]string)
	}
	params.Headers["Content-Type"] = writer.FormDataContentType()
	for headerKey, headerValue := range params.Headers {
		req.Header.Set(headerKey, headerValue)
	}
	if logger.Level() == t_logger.Level_DEBUG {
		logger.DebugF(
			`[HTTP POST] 
Url: %s
Headers: %v
Body: %s

`,
			url,
			params.Headers,
			body.String(),
		)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if logger.Level() == t_logger.Level_DEBUG {
		logger.DebugF(
			`[HTTP POST result] 
Url: %s
Body: %s

`,
			url,
			string(respBytes),
		)
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
	res, bodyBytes, err := t.Get(logger, params)
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
	res, bodyBytes, err := t.Get(logger, params)
	if err != nil {
		return nil, "", err
	}
	return res, string(bodyBytes), nil
}

func (t *HttpType) Get(
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
	if params.Headers == nil {
		params.Headers = make(map[string]string)
	}
	params.Headers["Content-Type"] = "application/json"
	for headerKey, headerValue := range params.Headers {
		req.Header.Set(headerKey, headerValue)
	}

	if logger.Level() == t_logger.Level_DEBUG {
		logger.DebugF(
			`[HTTP GET] 
Url: %s
Headers: %v

`,
			url,
			params.Headers,
		)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if logger.Level() == t_logger.Level_DEBUG {
		logger.DebugF(
			`[HTTP GET result] 
Url: %s
Body: %s

`,
			url,
			string(respBytes),
		)
	}

	return resp, respBytes, nil
}
