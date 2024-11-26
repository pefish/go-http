package go_http

import (
	"encoding/json"

	go_desensitize "github.com/pefish/go-desensitize"
	go_format "github.com/pefish/go-format"

	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/pefish/go-http/gorequest"
	i_logger "github.com/pefish/go-interface/i-logger"
	t_logger "github.com/pefish/go-interface/t-logger"
	"github.com/pkg/errors"
)

type IHttp interface {
	PostMultipart(params *PostMultipartParams) (res *http.Response, body string, err error)
	PostMultipartForStruct(params *PostMultipartParams, struct_ interface{}) (res *http.Response, bodyBytes []byte, err error)
	PostForStruct(params *RequestParams, struct_ interface{}) (res *http.Response, bodyBytes []byte, err error)
	PostForString(params *RequestParams) (res *http.Response, body string, err error)
	PostForBytes(params *RequestParams) (res *http.Response, bodyBytes []byte, err error)
	GetForString(params *RequestParams) (res *http.Response, body string, err error)
	GetForBytes(params *RequestParams) (res *http.Response, bodyBytes []byte, err error)
	GetForStruct(params *RequestParams, obj interface{}) (res *http.Response, bodyBytes []byte, err error)
	PostFormDataForStruct(params *RequestParams, struct_ interface{}) (res *http.Response, bodyBytes []byte, err error)
}

type HttpClass struct {
	timeout   time.Duration
	logger    i_logger.ILogger
	httpProxy string
}

type HttpRequestOptionFunc func(options *HttpRequestOption)

type HttpRequestOption struct {
	timeout   time.Duration
	logger    i_logger.ILogger
	httpProxy string
}

var defaultHttpRequestOption = HttpRequestOption{
	timeout: 10 * time.Second,
	logger:  &i_logger.DefaultLogger,
}

func WithTimeout(timeout time.Duration) HttpRequestOptionFunc {
	return func(option *HttpRequestOption) {
		option.timeout = timeout
	}
}

func WithLogger(logger i_logger.ILogger) HttpRequestOptionFunc {
	return func(option *HttpRequestOption) {
		option.logger = logger
	}
}

func WithHttpProxy(proxyUrl string) HttpRequestOptionFunc {
	return func(option *HttpRequestOption) {
		option.httpProxy = proxyUrl
	}
}

func NewHttpRequester(opts ...HttpRequestOptionFunc) IHttp {
	option := defaultHttpRequestOption
	for _, o := range opts {
		o(&option)
	}
	return &HttpClass{
		timeout:   option.timeout,
		logger:    option.logger,
		httpProxy: option.httpProxy,
	}
}

type RequestParams struct {
	Url       string
	Queries   map[string]string
	Params    interface{}
	Headers   map[string]interface{}
	BasicAuth *BasicAuth
}

type BytesFileInfo struct {
	Bytes    []byte
	FileName string
}

type BasicAuth struct {
	Username string
	Password string
}

type PostMultipartParams struct {
	RequestParams
	Files map[string][]BytesFileInfo
}

func (httpInstance *HttpClass) makeMultipartRequest(params *PostMultipartParams) *gorequest.SuperAgent {
	requestClient := gorequest.New(httpInstance.logger).
		Proxy(httpInstance.httpProxy).
		Timeout(httpInstance.timeout)
	requestClient = requestClient.Type(gorequest.TypeMultipart)
	httpInstance.decorateRequest(requestClient, &params.RequestParams, gorequest.POST)
	for keyName, fileArr := range params.Files {
		for _, file := range fileArr {
			requestClient = requestClient.SendFile(file.Bytes, file.FileName, keyName)
		}
	}
	return requestClient
}

func (httpInstance *HttpClass) PostMultipart(params *PostMultipartParams) (res *http.Response, body string, err error) {
	request := httpInstance.makeMultipartRequest(params)
	response, body, errs := request.Send(params.Params).End()
	if len(errs) > 0 {
		return nil, body, httpInstance.combineErrors(params.Url, params.Params, errs, body)
	}
	return response, body, nil
}

func (httpInstance *HttpClass) PostFormDataForStruct(params *RequestParams, struct_ interface{}) (res *http.Response, bodyBytes []byte, err error) {
	requestClient := gorequest.
		New(httpInstance.logger).
		Proxy(httpInstance.httpProxy).
		Timeout(httpInstance.timeout)
	requestClient.TargetType = gorequest.TypeForm
	httpInstance.decorateRequest(requestClient, params, gorequest.POST)
	response, bodyBytes, errs := requestClient.
		Send(params.Params).
		EndStruct(struct_)
	if len(errs) > 0 {
		return nil, bodyBytes, httpInstance.combineErrors(params.Url, params.Params, errs, string(bodyBytes))
	}
	return response, bodyBytes, nil
}

func (httpInstance *HttpClass) PostMultipartForStruct(params *PostMultipartParams, struct_ interface{}) (res *http.Response, bodyBytes []byte, err error) {
	request := httpInstance.makeMultipartRequest(params)
	response, bodyBytes, errs := request.Send(params.Params).EndStruct(struct_)
	if len(errs) > 0 {
		return nil, bodyBytes, httpInstance.combineErrors(params.Url, params.Params, errs, string(bodyBytes))
	}
	return response, bodyBytes, nil
}

func (httpInstance *HttpClass) PostForStruct(params *RequestParams, struct_ interface{}) (res *http.Response, bodyBytes []byte, err error) {
	requestClient := gorequest.New(httpInstance.logger).Proxy(httpInstance.httpProxy).Timeout(httpInstance.timeout)
	requestClient.TargetType = gorequest.TypeJSON
	httpInstance.decorateRequest(requestClient, params, gorequest.POST)
	response, bodyBytes, errs := requestClient.
		Send(params.Params).
		EndStruct(struct_)
	if len(errs) > 0 {
		return nil, bodyBytes, httpInstance.combineErrors(params.Url, params.Params, errs, string(bodyBytes))
	}
	return response, bodyBytes, nil
}

func (httpInstance *HttpClass) PostForString(params *RequestParams) (res *http.Response, body string, err error) {
	res, b, err := httpInstance.PostForBytes(params)
	if err != nil {
		return nil, "", err
	}
	return res, string(b), nil
}

func (httpInstance *HttpClass) PostForBytes(params *RequestParams) (res *http.Response, bodyBytes []byte, err error) {
	requestClient := gorequest.
		New(httpInstance.logger).
		Proxy(httpInstance.httpProxy).
		Timeout(httpInstance.timeout)
	requestClient.TargetType = gorequest.TypeJSON
	httpInstance.decorateRequest(requestClient, params, gorequest.POST)
	response, bodyBytes, errs := requestClient.
		Send(params.Params).
		EndBytes()
	if len(errs) > 0 {
		return nil, bodyBytes, httpInstance.combineErrors(params.Url, params.Params, errs, string(bodyBytes))
	}
	return response, bodyBytes, nil
}

func (httpInstance *HttpClass) decorateRequest(requestClient *gorequest.SuperAgent, params *RequestParams, method string) {
	requestClient.Method = method

	requestClient.Debug = httpInstance.logger.Level() == t_logger.Level_DEBUG

	urlParamsStr := ""
	if method == gorequest.GET {
		if params.Queries != nil {
			urlParamsStr = mapToUrlQuery(params.Queries)
		}
		if params.Params != nil {
			urlParams, err := interfaceToUrlQuery(params.Params)
			if err != nil {
				panic(err)
			}
			if urlParamsStr == "" {
				urlParamsStr = urlParams
			} else {
				urlParamsStr += "&" + urlParams
			}
		}
	} else if method == gorequest.POST {
		urlParamsStr = mapToUrlQuery(params.Queries)
	}
	requestClient.Url = params.Url
	if urlParamsStr != "" {
		requestClient.Url += "?" + urlParamsStr
	}

	if params.Headers != nil {
		for key, value := range params.Headers {
			str := go_format.ToString(value)
			requestClient.Set(key, str)
		}
	}
	if params.BasicAuth != nil {
		requestClient = requestClient.SetBasicAuth(params.BasicAuth.Username, params.BasicAuth.Password)
	}
	if params.Params != nil {
		switch params.Params.(type) {
		case string:
			requestClient.BounceToRawString = true
		}
	}
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

func interfaceToUrlQuery(params interface{}) (string, error) {
	if params == nil {
		return ``, nil
	}
	type_ := reflect.TypeOf(params)
	kind := type_.Kind()
	var strParams string
	if kind == reflect.Map {
		paramsMap, ok := params.(map[string]interface{})
		if !ok {
			return ``, errors.Errorf(`%F cannot cast to map[string]interface{}`, params)
		}
		for key, value := range paramsMap {
			str := go_format.ToString(value)
			strParams += key + "=" + str + "&"
		}
	} else if kind == reflect.Struct {
		return interfaceToUrlQuery(go_format.StructToMap(params))
	} else if kind == reflect.Ptr {
		return interfaceToUrlQuery(reflect.ValueOf(params).Elem().Interface())
	} else if kind == reflect.String {
		paramsStr := params.(string)
		if paramsStr == "" {
			return "", nil
		}
		b := make(map[string]interface{})
		err := json.Unmarshal([]byte(paramsStr), &b)
		if err != nil {
			return "", errors.Wrap(err, "Unmarshal error.")
		}
		return interfaceToUrlQuery(b)
	} else {
		return ``, errors.New(`Params type error`)
	}
	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}
	return strParams, nil
}

func (httpInstance *HttpClass) GetForString(params *RequestParams) (res *http.Response, body string, err error) {
	res, b, err := httpInstance.GetForBytes(params)
	if err != nil {
		return nil, "", err
	}
	return res, string(b), nil
}

func (httpInstance *HttpClass) GetForBytes(params *RequestParams) (res *http.Response, bodyBytes []byte, err error) {
	requestClient := gorequest.
		New(httpInstance.logger).
		Proxy(httpInstance.httpProxy).
		Timeout(httpInstance.timeout)
	httpInstance.decorateRequest(requestClient, params, gorequest.GET)
	response, bodyBytes, errs := requestClient.EndBytes()
	if len(errs) > 0 {
		return nil, bodyBytes, httpInstance.combineErrors(params.Url, params.Params, errs, string(bodyBytes))
	}
	return response, bodyBytes, nil
}

func (httpInstance *HttpClass) GetForStruct(params *RequestParams, struct_ interface{}) (res *http.Response, bodyBytes []byte, err error) {
	requestClient := gorequest.
		New(httpInstance.logger).
		Proxy(httpInstance.httpProxy).
		Timeout(httpInstance.timeout)
	httpInstance.decorateRequest(requestClient, params, gorequest.GET)
	response, bodyBytes, errs := requestClient.EndStruct(struct_)
	if len(errs) > 0 {
		return nil, bodyBytes, httpInstance.combineErrors(params.Url, params.Params, errs, string(bodyBytes))
	}
	return response, bodyBytes, nil
}

func (httpInstance *HttpClass) combineErrors(
	url string,
	params interface{},
	errs []error,
	bodyStr string,
) error {
	errStrs := make([]string, 0, len(errs))
	for _, err := range errs {
		errStrs = append(errStrs, err.Error())
	}
	if len(bodyStr) > 200 {
		bodyStr = bodyStr[:200]
	}
	return errors.Errorf(
		"Url: %s, Params: %s, Body: %s. -- %s",
		url,
		go_desensitize.Desensitize.MustDesensitizeToString(params),
		bodyStr,
		strings.Join(errStrs, " -> "),
	)
}
