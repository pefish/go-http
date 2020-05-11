package go_http

import (
	"errors"
	"fmt"
	"github.com/pefish/go-format"
	"github.com/pefish/go-http/gorequest"

	"github.com/pefish/go-interface-logger"
	"github.com/pefish/go-reflect"
	"net/http"
	"reflect"
	"time"
)

type HttpClass struct {
	RequestClient *gorequest.SuperAgent
	logger        go_interface_logger.InterfaceLogger
	isDebug       bool
}

type HttpRequestOptionFunc func(options *HttpRequestOption)

type HttpRequestOption struct {
	timeout time.Duration
	logger  go_interface_logger.InterfaceLogger
	isDebug bool
}

var defaultHttpRequestOption = HttpRequestOption{
	timeout: 10 * time.Second,
}

func WithTimeout(timeout time.Duration) HttpRequestOptionFunc {
	return func(option *HttpRequestOption) {
		option.timeout = timeout
	}
}

func WithIsDebug(isDebug bool) HttpRequestOptionFunc {
	return func(option *HttpRequestOption) {
		option.isDebug = isDebug
	}
}

func WithLogger(logger go_interface_logger.InterfaceLogger) HttpRequestOptionFunc {
	return func(option *HttpRequestOption) {
		option.logger = logger
	}
}

func NewHttpRequester(opts ...HttpRequestOptionFunc) *HttpClass {
	option := defaultHttpRequestOption
	for _, o := range opts {
		o(&option)
	}
	return &HttpClass{
		RequestClient: gorequest.New().Timeout(option.timeout),
		logger:        option.logger,
		isDebug:       option.isDebug,
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
	request := httpInstance.RequestClient.Post(param.Url).Type("multipart")
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
	request := httpInstance.RequestClient.Post(param.Url)
	err := httpInstance.decorateRequest(request, param)
	if err != nil {
		return nil, err
	}
	response, _, errs := request.Send(param.Params).EndStruct(struct_)
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
	request := httpInstance.RequestClient.Post(param.Url)
	err := httpInstance.decorateRequest(request, param)
	if err != nil {
		return nil, ``, err
	}
	response, body, errs := request.Send(param.Params).End()
	if len(errs) > 0 {
		return nil, ``, errs[0]
	}
	return response, body, nil
}

func (httpInstance *HttpClass) decorateRequest(request *gorequest.SuperAgent, param RequestParam) error {
	request.Debug = httpInstance.isDebug
	if param.Headers != nil {
		for key, value := range param.Headers {
			str := go_reflect.Reflect.ToString(value)
			request.Set(key, str)
		}
	}
	if param.BasicAuth != nil {
		request = request.SetBasicAuth(param.BasicAuth.Username, param.BasicAuth.Password)
	}
	return nil
}

func (httpInstance *HttpClass) interfaceToUrlQuery(params interface{}) (string, error) {
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
		return httpInstance.interfaceToUrlQuery(go_format.Format.StructToMap(params))
	} else if kind == reflect.Ptr {
		return httpInstance.interfaceToUrlQuery(reflect.ValueOf(params).Elem().Interface())
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
	urlParams, err := httpInstance.interfaceToUrlQuery(param.Params)
	if err != nil {
		return nil, ``, err
	}
	request := httpInstance.RequestClient.Get(param.Url + urlParams)
	err = httpInstance.decorateRequest(request, param)
	if err != nil {
		return nil, ``, err
	}
	response, body, errs := request.End()
	if len(errs) > 0 {
		return nil, ``, errs[0]
	}
	return response, body, nil
}

func (httpInstance *HttpClass) MustGetForStruct(param RequestParam, struct_ interface{}) *http.Response {
	res, err := httpInstance.GetForStruct(param, struct_)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res
}

func (httpInstance *HttpClass) GetForStruct(param RequestParam, struct_ interface{}) (*http.Response, error) {
	urlParams, err := httpInstance.interfaceToUrlQuery(param.Params)
	if err != nil {
		return nil, err
	}
	request := httpInstance.RequestClient.Get(param.Url + urlParams)
	err = httpInstance.decorateRequest(request, param)
	if err != nil {
		return nil, err
	}
	response, _, errs := request.EndStruct(struct_)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return response, nil
}
