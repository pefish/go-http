package go_http

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/pefish/go-application"
	"github.com/pefish/go-format"

	"github.com/pefish/go-json"
	"github.com/pefish/go-reflect"
	"net/http"
	"reflect"
	"time"
)

var Http = NewHttpRequester()

type HttpClass struct {
	RequestClient *gorequest.SuperAgent
}

type HttpRequestOptionFunc func(options *HttpRequestOption)

type HttpRequestOption struct {
	timeout time.Duration
}

var defaultHttpRequestOption = HttpRequestOption{
	timeout: 10 * time.Second,
}

func WithTimeout(timeout time.Duration) HttpRequestOptionFunc {
	return func(option *HttpRequestOption) {
		option.timeout = timeout
	}
}

func NewHttpRequester(opts ...HttpRequestOptionFunc) *HttpClass {
	option := defaultHttpRequestOption
	for _, o := range opts {
		o(&option)
	}
	return &HttpClass{
		RequestClient: gorequest.New().Timeout(option.timeout),
	}
}

type RequestParam struct {
	Url       string
	Params    interface{}
	Headers   map[string]interface{}
	BasicAuth *BasicAuth
}

func (this *HttpClass) MustPostJson(param RequestParam) interface{} {
	body, err := this.PostJson(param)
	if err != nil {
		panic(err)
	}
	return body
}

func (this *HttpClass) PostJson(param RequestParam) (interface{}, error) {
	body, err := this.PostJsonForString(param)
	if err != nil {
		return nil, err
	}
	return go_json.Json.Parse(body)
}

func (this *HttpClass) MustPostJsonForMap(param RequestParam) map[string]interface{} {
	body, err := this.PostJsonForMap(param)
	if err != nil {
		panic(err)
	}
	return body
}

func (this *HttpClass) PostJsonForMap(param RequestParam) (map[string]interface{}, error) {
	body, err := this.PostJsonForString(param)
	if err != nil {
		return nil, err
	}
	map_, err := go_json.Json.ParseToMap(body)
	if err != nil {
		return nil, err
	}
	return map_, nil
}

func (this *HttpClass) MustPostForMap(param RequestParam) map[string]interface{} {
	body, err := this.PostForMap(param)
	if err != nil {
		panic(err)
	}
	return body
}

func (this *HttpClass) PostForMap(param RequestParam) (map[string]interface{}, error) {
	body, err := this.PostForString(param)
	if err != nil {
		return nil, err
	}
	map_, err := go_json.Json.ParseToMap(body)
	if err != nil {
		return nil, err
	}
	return map_, nil
}

func (this *HttpClass) MustPostJsonForString(param RequestParam) string {
	body, err := this.PostForString(param)
	if err != nil {
		panic(err)
	}
	return body
}

func (this *HttpClass) PostJsonForString(param RequestParam) (string, error) {
	if param.Headers != nil {
		param.Headers[`Content-Type`] = `application/json`
	} else {
		param.Headers = map[string]interface{}{
			`Content-Type`: `application/json`,
		}
	}
	return this.PostForString(param)
}

func (this *HttpClass) MustPostForString(param RequestParam) string {
	body, err := this.PostForString(param)
	if err != nil {
		panic(err)
	}
	return body
}

func (this *HttpClass) PostForString(param RequestParam) (string, error) {
	_, body, err := this.Post(param)
	if err != nil {
		return ``, err
	}
	return body, nil
}

func (this *HttpClass) MustPostMultipartForMap(param PostMultipartParam) map[string]interface{} {
	body, err := this.PostMultipartForMap(param)
	if err != nil {
		panic(err)
	}
	return body
}

func (this *HttpClass) PostMultipartForMap(param PostMultipartParam) (map[string]interface{}, error) {
	body, err := this.PostMultipartForString(param)
	if err != nil {
		return nil, err
	}
	map_, err := go_json.Json.ParseToMap(body)
	if err != nil {
		return nil, err
	}
	return map_, nil
}

func (this *HttpClass) MustPostMultipartForString(param PostMultipartParam) string {
	body, err := this.PostMultipartForString(param)
	if err != nil {
		panic(err)
	}
	return body
}

func (this *HttpClass) PostMultipartForString(param PostMultipartParam) (string, error) {
	_, body, err := this.PostMultipart(param)
	if err != nil {
		return ``, err
	}
	return body, nil
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

func (this *HttpClass) MustPostMultipart(param PostMultipartParam) (*http.Response, string) {
	res, body, err := this.PostMultipart(param)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res, body
}

func (this *HttpClass) PostMultipart(param PostMultipartParam) (*http.Response, string, error) {
	request := this.RequestClient.Post(param.Url).Type("multipart")
	request.Debug = go_application.Application.Debug
	if param.Headers != nil {
		for key, value := range param.Headers {
			request.Set(key, go_reflect.Reflect.MustToString(value))
		}
	}
	if param.BasicAuth != nil {
		request = request.SetBasicAuth(param.BasicAuth.Username, param.BasicAuth.Password)
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

func (this *HttpClass) MustPostJsonForStruct(param RequestParam, struct_ interface{}) *http.Response {
	res, err := this.PostJsonForStruct(param, struct_)
	if err != nil {
		panic(err)
	}
	return res
}

func (this *HttpClass) PostJsonForStruct(param RequestParam, struct_ interface{}) (*http.Response, error) {
	if param.Headers != nil {
		param.Headers[`Content-Type`] = `application/json`
	} else {
		param.Headers = map[string]interface{}{
			`Content-Type`: `application/json`,
		}
	}
	return this.PostForStruct(param, struct_)
}

func (this *HttpClass) MustPostForStruct(param RequestParam, struct_ interface{}) *http.Response {
	res, err := this.PostForStruct(param, struct_)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res
}

func (this *HttpClass) PostForStruct(param RequestParam, struct_ interface{}) (*http.Response, error) {
	request := this.RequestClient.Post(param.Url)
	request.Debug = go_application.Application.Debug
	if param.Headers != nil {
		for key, value := range param.Headers {
			request.Set(key, go_reflect.Reflect.MustToString(value))
		}
	}
	if param.BasicAuth != nil {
		request = request.SetBasicAuth(param.BasicAuth.Username, param.BasicAuth.Password)
	}
	response, _, errs := request.Send(param.Params).EndStruct(struct_)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return response, nil
}

func (this *HttpClass) MustPost(param RequestParam) (*http.Response, string) {
	res, body, err := this.Post(param)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res, body
}

func (this *HttpClass) Post(param RequestParam) (*http.Response, string, error) {
	request := this.RequestClient.Post(param.Url)
	request.Debug = go_application.Application.Debug
	if param.Headers != nil {
		for key, value := range param.Headers {
			request.Set(key, go_reflect.Reflect.MustToString(value))
		}
	}
	if param.BasicAuth != nil {
		request = request.SetBasicAuth(param.BasicAuth.Username, param.BasicAuth.Password)
	}
	response, body, errs := request.Send(param.Params).End()
	if len(errs) > 0 {
		return nil, ``, errs[0]
	}
	return response, body, nil
}

func (this *HttpClass) MustGetForMap(param RequestParam) map[string]interface{} {
	body, err := this.GetForMap(param)
	if err != nil {
		panic(err)
	}
	return body
}

func (this *HttpClass) GetForMap(param RequestParam) (map[string]interface{}, error) {
	body, err := this.GetForString(param)
	if err != nil {
		return nil, err
	}
	map_, err := go_json.Json.ParseToMap(body)
	if err != nil {
		return nil, err
	}
	return map_, nil
}

func (this *HttpClass) MustGetForString(param RequestParam) string {
	body, err := this.GetForString(param)
	if err != nil {
		panic(err)
	}
	return body
}

func (this *HttpClass) GetForString(param RequestParam) (string, error) {
	_, body, err := this.Get(param)
	if err != nil {
		return ``, err
	}
	return body, nil
}

func (this *HttpClass) interfaceToUrlQuery(params interface{}) (string, error) {
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
			strParams += key + "=" + go_reflect.Reflect.MustToString(value) + "&"
		}
	} else if kind == reflect.Struct {
		return this.interfaceToUrlQuery(go_format.Format.StructToMap(params))
	} else if kind == reflect.Ptr {
		return this.interfaceToUrlQuery(reflect.ValueOf(params).Elem().Interface())
	} else {
		return ``, errors.New(`Params type error`)
	}
	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}
	return `?` + strParams, nil
}

func (this *HttpClass) MustGet(param RequestParam) (*http.Response, string) {
	res, body, err := this.Get(param)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res, body
}

func (this *HttpClass) Get(param RequestParam) (*http.Response, string, error) {
	urlParams, err := this.interfaceToUrlQuery(param.Params)
	if err != nil {
		return nil, ``, err
	}
	request := this.RequestClient.Get(param.Url + urlParams)
	request.Debug = go_application.Application.Debug
	if param.Headers != nil {
		for key, value := range param.Headers {
			request.Set(key, go_reflect.Reflect.MustToString(value))
		}
	}
	if param.BasicAuth != nil {
		request = request.SetBasicAuth(param.BasicAuth.Username, param.BasicAuth.Password)
	}
	response, body, errs := request.End()
	if len(errs) > 0 {
		return nil, ``, errs[0]
	}
	return response, body, nil
}

func (this *HttpClass) MustGetForStruct(param RequestParam, struct_ interface{}) *http.Response {
	res, err := this.GetForStruct(param, struct_)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res
}

func (this *HttpClass) GetForStruct(param RequestParam, struct_ interface{}) (*http.Response, error) {
	urlParams, err := this.interfaceToUrlQuery(param.Params)
	if err != nil {
		return nil, err
	}
	request := this.RequestClient.Get(param.Url + urlParams)
	request.Debug = go_application.Application.Debug
	if param.Headers != nil {
		for key, value := range param.Headers {
			request.Set(key, go_reflect.Reflect.MustToString(value))
		}
	}
	if param.BasicAuth != nil {
		request = request.SetBasicAuth(param.BasicAuth.Username, param.BasicAuth.Password)
	}
	response, _, errs := request.EndStruct(struct_)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return response, nil
}

func (this *HttpClass) MustPutForStruct(param RequestParam, struct_ interface{}) *http.Response {
	res, err := this.PutForStruct(param, struct_)
	if err != nil {
		panic(errors.New(fmt.Sprintf(`ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, err)))
	}
	return res
}

func (this *HttpClass) PutForStruct(param RequestParam, struct_ interface{}) (*http.Response, error) {
	request := this.RequestClient.Put(param.Url)
	request.Debug = go_application.Application.Debug
	if param.Headers != nil {
		for key, value := range param.Headers {
			request.Set(key, go_reflect.Reflect.MustToString(value))
		}
	}
	if param.BasicAuth != nil {
		request = request.SetBasicAuth(param.BasicAuth.Username, param.BasicAuth.Password)
	}
	response, _, errs := request.Send(param.Params).EndStruct(struct_)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return response, nil
}
