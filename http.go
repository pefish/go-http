package go_http

import (
	"errors"
	"fmt"
	"github.com/pefish/go-application"
	"github.com/pefish/go-format"
	"github.com/pefish/go-json"
	"github.com/pefish/go-reflect"
	"github.com/pefish/gorequest"
	"net/http"
	"reflect"
	"time"
)

type HttpClass struct {
	timeout time.Duration
}

var Http = HttpClass{10 * time.Second}

func (this *HttpClass) SetTimeout(timeout time.Duration) {
	this.timeout = timeout
}

type RequestParam struct {
	url     string
	params  interface{}
	headers map[string]interface{}
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
	request := gorequest.New()
	request.Debug = go_application.Application.Debug
	req := request.Timeout(this.timeout).Post(param.url)
	if param.headers != nil {
		for key, value := range param.headers {
			req.Set(key, go_reflect.Reflect.ToString(value))
		}
	}
	_, body, errs := request.Set(`Content-Type`, `application/json`).Send(param.params).End()
	if len(errs) > 0 {
		panic(errors.New(fmt.Sprintf(`PostJsonForString ERROR!! url: %s, params: %v, error: %v`, param.url, param.params, errs[0])))
	}
	return body
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

type PostMultipartParam struct {
	url string
	params interface{}
	files map[string][]BytesFileInfo
	headers map[string]interface{}
}

func (this *HttpClass) PostMultipart(param PostMultipartParam) (*http.Response, string) {
	request := gorequest.New()
	request.Debug = go_application.Application.Debug
	req := request.Post(param.url).Type("multipart")
	if param.headers != nil {
		for key, value := range param.headers {
			req.Set(key, go_reflect.Reflect.ToString(value))
		}
	}
	for keyName, fileArr := range param.files {
		for _, file := range fileArr {
			req = req.SendFile(file.Bytes, file.FileName, keyName)
		}
	}
	response, body, errs := req.Send(param.params).End()
	if len(errs) > 0 {
		panic(errors.New(fmt.Sprintf(`PostMultipart ERROR!! url: %s, params: %v, error: %v`, param.url, param.params, errs[0])))
	}
	return response, body
}

func (this *HttpClass) Post(param RequestParam) (*http.Response, string) {
	request := gorequest.New()
	request.Debug = go_application.Application.Debug
	req := request.Timeout(this.timeout).Post(param.url)
	if param.headers != nil {
		for key, value := range param.headers {
			req.Set(key, go_reflect.Reflect.ToString(value))
		}
	}
	response, body, errs := req.Send(param.params).End()
	if len(errs) > 0 {
		panic(errors.New(fmt.Sprintf(`PostForString ERROR!! url: %s, params: %v, error: %v`, param.url, param.params, errs[0])))
	}
	return response, body
}

func (this *HttpClass) GetForMap(param RequestParam) map[string]interface{} {
	return go_json.Json.Parse(this.GetForString(param)).(map[string]interface{})
}

func (this *HttpClass) GetWithParamsForMap(param RequestParam) map[string]interface{} {
	return go_json.Json.Parse(this.GetWithParamsForString(param)).(map[string]interface{})
}

func (this *HttpClass) GetForString(param RequestParam) string {
	_, body := this.Get(param)
	return body
}

func (this *HttpClass) GetWithParamsForString(param RequestParam) string {
	return this.GetForString(param)
}

func (this *HttpClass) interfaceToUrlQuery(params interface{}) string {
	type_ := reflect.TypeOf(params)
	kind := type_.Kind()
	var strParams string
	if kind == reflect.Map {
		eleKind := type_.Elem().Kind()
		if eleKind == reflect.Interface {
			relParams := params.(map[string]interface{})
			for key, value := range relParams {
				strParams += (key + "=" + go_reflect.Reflect.ToString(value) + "&")
			}
		} else if eleKind == reflect.String {
			relParams := params.(map[string]string)
			for key, value := range relParams {
				strParams += (key + "=" + value + "&")
			}
		} else {
			panic(errors.New(`params type error`))
		}
	} else if kind == reflect.Struct {
		relParams := go_format.Format.StructToMap(params)
		for key, value := range relParams {
			strParams += (key + "=" + go_reflect.Reflect.ToString(value) + "&")
		}
	} else if kind == reflect.Ptr {
		return this.interfaceToUrlQuery(reflect.ValueOf(params).Elem().Interface())
	} else {
		panic(errors.New(`params type error`))
	}
	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}
	return strParams
}

func (this *HttpClass) GetWithParams(param RequestParam) (*http.Response, string) {
	return this.Get(param)
}

func (this *HttpClass) Get(param RequestParam) (*http.Response, string) {
	request := gorequest.New()
	request.Debug = go_application.Application.Debug
	req := request.Timeout(this.timeout).Get(param.url)
	if param.headers != nil {
		for key, value := range param.headers {
			req.Set(key, go_reflect.Reflect.ToString(value))
		}
	}
	response, body, errs := req.End()
	if len(errs) > 0 {
		panic(errors.New(fmt.Sprintf(`GetForString ERROR!! url: %s, error: %v`, param.url, errs[0])))
	}
	return response, body
}