package go_http

import (
	"encoding/json"
	"fmt"
	go_format "github.com/pefish/go-format"
	"github.com/pefish/go-http/gorequest"
	"github.com/pefish/go-logger"
	go_reflect "github.com/pefish/go-reflect"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"time"
)

type IHttp interface {
	MustPostMultipart(param PostMultipartParam) (*http.Response, string)
	PostMultipart(param PostMultipartParam) (*http.Response, string, error)
	MustPostForStruct(param RequestParam, struct_ interface{}) *http.Response
	PostForStruct(param RequestParam, struct_ interface{}) (*http.Response, error)
	MustPost(param RequestParam) (*http.Response, string)
	Post(param RequestParam) (*http.Response, string, error)
	MustPostForBytes(param RequestParam) (*http.Response, []byte)
	PostForBytes(param RequestParam) (*http.Response, []byte, error)
	MustGet(param RequestParam) (*http.Response, string)
	Get(param RequestParam) (*http.Response, string, error)
	MustGetForBytes(param RequestParam) (*http.Response, []byte)
	GetForBytes(param RequestParam) (*http.Response, []byte, error)
	MustGetForObject(param RequestParam, obj interface{}) *http.Response
	GetForObject(param RequestParam, obj interface{}) (*http.Response, error)
}

type HttpClass struct {
	timeout time.Duration
	logger  go_logger.InterfaceLogger
	httpProxy string
}

type HttpRequestOptionFunc func(options *HttpRequestOption)

type HttpRequestOption struct {
	timeout time.Duration
	logger  go_logger.InterfaceLogger
	httpProxy string
}

var defaultHttpRequestOption = HttpRequestOption{
	timeout: 10 * time.Second,
	logger:  go_logger.DefaultLogger,
}

func WithTimeout(timeout time.Duration) HttpRequestOptionFunc {
	return func(option *HttpRequestOption) {
		option.timeout = timeout
	}
}

func WithLogger(logger go_logger.InterfaceLogger) HttpRequestOptionFunc {
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
		timeout: option.timeout,
		logger:  option.logger,
		httpProxy: option.httpProxy,
	}
}

type RequestParam struct {
	Url       string
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

type PostMultipartParam struct {
	Url       string
	Params    interface{}
	Files     map[string][]BytesFileInfo
	Headers   map[string]interface{}
	BasicAuth *BasicAuth
}

func (httpInstance *HttpClass) MustPostMultipart(param PostMultipartParam) (*http.Response, string) {
	res, body, err := httpInstance.PostMultipart(param)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res, body
}

func (httpInstance *HttpClass) PostMultipart(param PostMultipartParam) (*http.Response, string, error) {
	requestClient := gorequest.New().Proxy(httpInstance.httpProxy).Timeout(httpInstance.timeout)
	requestClient.Method = gorequest.POST
	requestClient.Url = param.Url
	request := requestClient.Type("multipart")
	err := httpInstance.decorateRequest(request, RequestParam{
		Url:       param.Url,
		Params:    param.Params,
		Headers:   param.Headers,
		BasicAuth: param.BasicAuth,
	})
	if err != nil {
		return nil, ``, err
	}
	for keyName, fileArr := range param.Files {
		for _, file := range fileArr {
			request = request.SendFile(file.Bytes, file.FileName, keyName)
		}
	}
	response, body, errs := request.Send(param.Params).End()
	if len(errs) > 0 {
		return nil, ``, errs[0]
	}
	return response, body, nil
}

func (httpInstance *HttpClass) MustPostForStruct(param RequestParam, struct_ interface{}) *http.Response {
	res, err := httpInstance.PostForStruct(param, struct_)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res
}

func (httpInstance *HttpClass) PostForStruct(param RequestParam, struct_ interface{}) (*http.Response, error) {
	requestClient := gorequest.New().Proxy(httpInstance.httpProxy).Timeout(httpInstance.timeout)
	requestClient.Method = gorequest.POST
	requestClient.Url = param.Url
	err := httpInstance.decorateRequest(requestClient, param)
	if err != nil {
		return nil, err
	}
	response, _, errs := requestClient.Send(param.Params).EndStruct(struct_)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return response, nil
}

func (httpInstance *HttpClass) MustPost(param RequestParam) (*http.Response, string) {
	res, body, err := httpInstance.Post(param)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res, body
}

func (httpInstance *HttpClass) Post(param RequestParam) (*http.Response, string, error) {
	res, b, err := httpInstance.PostForBytes(param)
	if err != nil {
		return nil, "", err
	}
	return res, string(b), nil
}

func (httpInstance *HttpClass) MustPostForBytes(param RequestParam) (*http.Response, []byte) {
	res, bodyBytes, err := httpInstance.PostForBytes(param)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res, bodyBytes
}

func (httpInstance *HttpClass) PostForBytes(param RequestParam) (*http.Response, []byte, error) {
	requestClient := gorequest.New().Proxy(httpInstance.httpProxy).Timeout(httpInstance.timeout)
	requestClient.Method = gorequest.POST
	requestClient.Url = param.Url
	err := httpInstance.decorateRequest(requestClient, param)
	if err != nil {
		return nil, nil, err
	}
	response, bodyBytes, errs := requestClient.Send(param.Params).EndBytes()
	if len(errs) > 0 {
		return nil, nil, errs[0]
	}
	return response, bodyBytes, nil
}

func (httpInstance *HttpClass) decorateRequest(request *gorequest.SuperAgent, param RequestParam) error {
	request.Debug = httpInstance.logger.IsDebug()
	if param.Headers != nil {
		for key, value := range param.Headers {
			str := go_reflect.Reflect.ToString(value)
			request.Set(key, str)
		}
	}
	if param.BasicAuth != nil {
		request = request.SetBasicAuth(param.BasicAuth.Username, param.BasicAuth.Password)
	}
	if param.Params != nil {
		switch param.Params.(type) {
		case string:
			request.BounceToRawString = true
		}
	}
	return nil
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
			return ``, errors.New(fmt.Sprintf(`%F cannot cast to map[string]interface{}`, params))
		}
		for key, value := range paramsMap {
			str := go_reflect.Reflect.ToString(value)
			strParams += key + "=" + str + "&"
		}
	} else if kind == reflect.Struct {
		return interfaceToUrlQuery(go_format.Format.StructToMap(params))
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
			return "", err
		}
		return interfaceToUrlQuery(b)
	} else {
		return ``, errors.New(`Params type error`)
	}
	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}
	return `?` + strParams, nil
}

func (httpInstance *HttpClass) MustGet(param RequestParam) (*http.Response, string) {
	res, body, err := httpInstance.Get(param)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res, body
}

func (httpInstance *HttpClass) Get(param RequestParam) (*http.Response, string, error) {
	res, b, err := httpInstance.GetForBytes(param)
	if err != nil {
		return nil, "", err
	}
	return res, string(b), nil
}

func (httpInstance *HttpClass) MustGetForBytes(param RequestParam) (*http.Response, []byte) {
	res, bodyBytes, err := httpInstance.GetForBytes(param)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res, bodyBytes
}

func (httpInstance *HttpClass) GetForBytes(param RequestParam) (*http.Response, []byte, error) {
	requestClient := gorequest.New().Proxy(httpInstance.httpProxy).Timeout(httpInstance.timeout)
	urlParams, err := interfaceToUrlQuery(param.Params)
	if err != nil {
		return nil, nil, err
	}
	requestClient.Method = gorequest.GET
	requestClient.Url = param.Url + urlParams
	err = httpInstance.decorateRequest(requestClient, param)
	if err != nil {
		return nil, nil, err
	}
	response, bodyBytes, errs := requestClient.EndBytes()
	if len(errs) > 0 {
		return nil, nil, errs[0]
	}
	return response, bodyBytes, nil
}

func (httpInstance *HttpClass) MustGetForObject(param RequestParam, struct_ interface{}) *http.Response {
	res, err := httpInstance.GetForObject(param, struct_)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res
}

func (httpInstance *HttpClass) GetForObject(param RequestParam, struct_ interface{}) (*http.Response, error) {
	requestClient := gorequest.New().Proxy(httpInstance.httpProxy).Timeout(httpInstance.timeout)
	urlParams, err := interfaceToUrlQuery(param.Params)
	if err != nil {
		return nil, err
	}
	requestClient.Method = gorequest.GET
	requestClient.Url = param.Url + urlParams
	err = httpInstance.decorateRequest(requestClient, param)
	if err != nil {
		return nil, err
	}
	response, _, errs := requestClient.EndStruct(struct_)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return response, nil
}
