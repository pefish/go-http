package p_http

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pefish/go-application"
	"github.com/pefish/go-format"
	"github.com/pefish/go-json"
	"github.com/pefish/go-reflect"
	"github.com/pefish/gorequest"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"time"
)

type HttpClass struct {
	Timeout time.Duration
}

var Http = HttpClass{10 * time.Second}

func (this *HttpClass) SetTimeout(timeout time.Duration) {
	this.Timeout = timeout
}

func (this *HttpClass) PostJson(url string, params interface{}) interface{} {
	return p_json.Json.Parse(this.PostJsonForString(url, params))
}

func (this *HttpClass) PostJsonForMap(url string, params interface{}) map[string]interface{} {
	return p_json.Json.Parse(this.PostJsonForString(url, params)).(map[string]interface{})
}

func (this *HttpClass) PostForMap(url string, params interface{}, headers map[string]string) map[string]interface{} {
	return p_json.Json.Parse(this.PostForString(url, params, headers)).(map[string]interface{})
}

func (this *HttpClass) PostJsonForString(url string, params interface{}) string {
	request := gorequest.New()
	request.Debug = p_application.Application.Debug
	_, body, errs := request.Timeout(this.Timeout).Post(url).Set(`Content-Type`, `application/json`).Send(params).End()
	if len(errs) > 0 {
		panic(errors.New(fmt.Sprintf(`PostJsonForString ERROR!! url: %s, params: %v, error: %v`, url, params, errs[0])))
	}
	return body
}

func (this *HttpClass) PostForString(url string, params interface{}, headers map[string]string) string {
	request := gorequest.New()
	request.Debug = p_application.Application.Debug
	req := request.Timeout(this.Timeout).Post(url)
	if headers != nil {
		for key, value := range headers {
			req.Set(key, value)
		}
	}
	_, body, errs := req.Send(params).End()
	if len(errs) > 0 {
		panic(errors.New(fmt.Sprintf(`PostForString ERROR!! url: %s, params: %v, error: %v`, url, params, errs[0])))
	}
	return body
}

func (this *HttpClass) PostMultipartForMap(url string, params interface{}, files map[string][]BytesFileInfo, headers map[string]string) map[string]interface{} {
	body := this.PostMultipartForString(url, params, files, headers)
	return p_json.Json.ParseToMap(body)
}

func (this *HttpClass) PostMultipartForString(url string, params interface{}, files map[string][]BytesFileInfo, headers map[string]string) string {
	_, body := this.PostMultipart(url, params, files, headers)
	return body
}

type BytesFileInfo struct {
	Bytes    []byte
	FileName string
}

func (this *HttpClass) PostMultipart(url string, params interface{}, files map[string][]BytesFileInfo, headers map[string]string) (*http.Response, string) {
	request := gorequest.New()
	request.Debug = p_application.Application.Debug
	req := request.Post(url).Type("multipart")
	if headers != nil {
		for key, value := range headers {
			req.Set(key, value)
		}
	}
	for keyName, fileArr := range files {
		for _, file := range fileArr {
			req = req.SendFile(file.Bytes, file.FileName, keyName)
		}
	}
	response, body, errs := req.Send(params).End()
	if len(errs) > 0 {
		panic(errors.New(fmt.Sprintf(`PostMultipart ERROR!! url: %s, params: %v, error: %v`, url, params, errs[0])))
	}
	return response, body
}

func (this *HttpClass) Post(url string, params interface{}, headers map[string]string) (*http.Response, string) {
	request := gorequest.New()
	request.Debug = p_application.Application.Debug
	req := request.Timeout(this.Timeout).Post(url)
	if headers != nil {
		for key, value := range headers {
			req.Set(key, value)
		}
	}
	response, body, errs := req.Send(params).End()
	if len(errs) > 0 {
		panic(errors.New(fmt.Sprintf(`PostForString ERROR!! url: %s, params: %v, error: %v`, url, params, errs[0])))
	}
	return response, body
}

func (this *HttpClass) GetForMap(url string, headers map[string]string) map[string]interface{} {
	return p_json.Json.Parse(this.GetForString(url, headers)).(map[string]interface{})
}

func (this *HttpClass) GetWithParamsForMap(url string, params interface{}, headers map[string]string) map[string]interface{} {
	return p_json.Json.Parse(this.GetWithParamsForString(url, params, headers)).(map[string]interface{})
}

func (this *HttpClass) GetForString(url string, headers map[string]string) string {
	request := gorequest.New()
	request.Debug = p_application.Application.Debug
	req := request.Timeout(this.Timeout).Get(url)
	if headers != nil {
		for key, value := range headers {
			req.Set(key, value)
		}
	}
	_, body, errs := req.End()
	if len(errs) > 0 {
		panic(errors.New(fmt.Sprintf(`GetForString ERROR!! url: %s, error: %v`, url, errs[0])))
	}
	return body
}

func (this *HttpClass) GetWithParamsForString(url string, params interface{}, headers map[string]string) string {
	return this.GetForString(url+`?`+this.interfaceToUrlQuery(params), headers)
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
				strParams += (key + "=" + p_reflect.Reflect.ToString(value) + "&")
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
		relParams := p_format.Format.StructToMap(params)
		for key, value := range relParams {
			strParams += (key + "=" + p_reflect.Reflect.ToString(value) + "&")
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

func (this *HttpClass) GetWithParams(url string, params interface{}, headers map[string]string) (*http.Response, string) {
	return this.Get(url+`?`+this.interfaceToUrlQuery(params), headers)
}

func (this *HttpClass) Get(url string, headers map[string]string) (*http.Response, string) {
	request := gorequest.New()
	request.Debug = p_application.Application.Debug
	req := request.Timeout(this.Timeout).Get(url)
	if headers != nil {
		for key, value := range headers {
			req.Set(key, value)
		}
	}
	response, body, errs := req.End()
	if len(errs) > 0 {
		panic(errors.New(fmt.Sprintf(`GetForString ERROR!! url: %s, error: %v`, url, errs[0])))
	}
	return response, body
}

func (this *HttpClass) MultipartWithMap(url string, params map[string]string) (map[string]interface{}, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows; U; Windows NT 6.1; zh-CN; rv:1.9.2.6)")
	req.Header.Set("Charset", "UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return p_json.Json.Parse(string(buffer)).(map[string]interface{}), nil
}
