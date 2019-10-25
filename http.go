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

func (this *HttpClass) PostJson(param RequestParam) interface{} {
	return go_json.Json.Parse(this.PostJsonForString(param))
}

func (this *HttpClass) PostJsonForMap(param RequestParam) map[string]interface{} {
	return go_json.Json.Parse(this.PostJsonForString(param)).(map[string]interface{})
}

func (this *HttpClass) PostForMap(param RequestParam) map[string]interface{} {
	return go_json.Json.Parse(this.PostForString(param)).(map[string]interface{})
}

func (this *HttpClass) PostJsonForString(param RequestParam) string {
	if param.Headers != nil {
		param.Headers[`Content-Type`] = `application/json`
	} else {
		param.Headers = map[string]interface{}{
			`Content-Type`: `application/json`,
		}
	}
	return this.PostForString(param)
}

func (this *HttpClass) PostForString(param RequestParam) string {
	_, body := this.Post(param)
	return body
}

func (this *HttpClass) PostMultipartForMap(param PostMultipartParam) map[string]interface{} {
	body := this.PostMultipartForString(param)
	return go_json.Json.ParseToMap(body)
}

func (this *HttpClass) PostMultipartForString(param PostMultipartParam) string {
	_, body := this.PostMultipart(param)
	return body
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

func (this *HttpClass) PostMultipart(param PostMultipartParam) (*http.Response, string) {
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
		panic(errors.New(fmt.Sprintf(`PostMultipart ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, errs[0])))
	}
	return response, body
}

func (this *HttpClass) PostJsonForStruct(param RequestParam, struct_ interface{}) *http.Response {
	if param.Headers != nil {
		param.Headers[`Content-Type`] = `application/json`
	} else {
		param.Headers = map[string]interface{}{
			`Content-Type`: `application/json`,
		}
	}
	return this.PostForStruct(param, struct_)
}

func (this *HttpClass) PostForStruct(param RequestParam, struct_ interface{}) *http.Response {
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
		panic(errors.New(fmt.Sprintf(`PostForStruct ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, errs[0])))
	}
	return response
}

func (this *HttpClass) Post(param RequestParam) (*http.Response, string) {
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
		panic(errors.New(fmt.Sprintf(`Post ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, errs[0])))
	}
	return response, body
}

func (this *HttpClass) GetForMap(param RequestParam) map[string]interface{} {
	return go_json.Json.Parse(this.GetForString(param)).(map[string]interface{})
}

func (this *HttpClass) GetForString(param RequestParam) string {
	_, body := this.Get(param)
	return body
}

func (this *HttpClass) interfaceToUrlQuery(params interface{}) string {
	if params == nil {
		return ``
	}
	type_ := reflect.TypeOf(params)
	kind := type_.Kind()
	var strParams string
	if kind == reflect.Map {
		for key, value := range params.(map[string]interface{}) {
			strParams += key + "=" + go_reflect.Reflect.MustToString(value) + "&"
		}
	} else if kind == reflect.Struct {
		return this.interfaceToUrlQuery(go_format.Format.StructToMap(params))
	} else if kind == reflect.Ptr {
		return this.interfaceToUrlQuery(reflect.ValueOf(params).Elem().Interface())
	} else {
		panic(errors.New(`Params type error`))
	}
	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}
	return `?` + strParams
}

func (this *HttpClass) Get(param RequestParam) (*http.Response, string) {
	request := this.RequestClient.Get(param.Url + this.interfaceToUrlQuery(param.Params))
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
		panic(errors.New(fmt.Sprintf(`Get ERROR!! Url: %s, error: %v`, param.Url, errs[0])))
	}
	return response, body
}

func (this *HttpClass) GetForStruct(param RequestParam, struct_ interface{}) *http.Response {
	request := this.RequestClient.Get(param.Url + this.interfaceToUrlQuery(param.Params))
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
		panic(errors.New(fmt.Sprintf(`GetForStruct ERROR!! Url: %s, error: %v`, param.Url, errs[0])))
	}
	return response
}

func (this *HttpClass) PutForStruct(param RequestParam, struct_ interface{}) *http.Response {
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
		panic(errors.New(fmt.Sprintf(`PutForStruct ERROR!! Url: %s, Params: %v, error: %v`, param.Url, param.Params, errs[0])))
	}
	return response
}
