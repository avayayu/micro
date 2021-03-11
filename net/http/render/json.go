package render

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
	ecode "gogs.bfr.com/zouhy/micro/code"
	"gogs.bfr.com/zouhy/micro/net/constants"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

// JSON common json struct.
type JSON struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TTL     int         `json:"ttl"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, obj interface{}) (err error) {
	var jsonBytes []byte
	writeContentType(w, jsonContentType)
	if jsonBytes, err = json.Marshal(obj); err != nil {
		err = errors.WithStack(err)
		return
	}
	if _, err = w.Write(jsonBytes); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// Render (JSON) writes data with json ContentType.
func (r JSON) Render(w http.ResponseWriter) error {
	// FIXME(zhoujiahui): the TTL field will be configurable in the future
	if r.TTL <= 0 {
		r.TTL = 1
	}
	return writeJSON(w, r)
}

// WriteContentType write json ContentType.
func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// MapJSON common map json struct.
type MapJSON map[string]interface{}

// Render (MapJSON) writes data with json ContentType.
func (m MapJSON) Render(w http.ResponseWriter) error {
	return writeJSON(w, m)
}

// WriteContentType write json ContentType.
func (m MapJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

/************************************/
/********Buffalo SPECIAL RESPONSE Type ********/
/********AUTHOR:ZouHangYu ********/
/************************************/

//Response 回复
type Response map[string]interface{}

//SUCCESS 接口成功
var success Response = (Response)(map[string]interface{}{"status": constants.HTTPSuccess})

//Fail 接口失败
var fail Response = (Response)(map[string]interface{}{"status": constants.HTTPFail})

//Success 复制一份success，并返回其指针
func SuccessResponse() *Response {
	successTemp := (Response)(map[string]interface{}{"status": constants.HTTPSuccess, "data": map[string]interface{}{}})
	// successCopy := success
	return &successTemp
}

//Fail 复制一份fail,并返回其指针
func FailResponse() *Response {
	failCopy := (Response)(map[string]interface{}{"status": constants.HTTPFail})
	return &failCopy
}

//Get 构造回复 value必须为指针
func (r *Response) Set(key string, value interface{}) *Response {

	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		panic("Set only accept ptr value")
	}

	if value == nil {
		return r
	}
	dataMap := (*r)["data"].(map[string]interface{})

	if _, ok := dataMap[key]; ok {
		panic("not duplicate key,please")
	}
	dataMap[key] = value

	return r
}

//Reason 构造回复 失败
func (r *Response) Reason(reason error) *Response {
	(*r)["info"] = reason.Error()
	return r
}

type httpResponse struct {
	HttpCode int        `json:"-"`    //http 头部code 500 404 200等
	Code     ecode.Code `json:"code"` //内部id
	Message  string     `json:"message"`
	Errors   []error    `json:"errors"` //内部errors 可能为空
	*Response
}

// Render (JSON) writes data with json ContentType.
func (r *httpResponse) Render(w http.ResponseWriter) error {
	// FIXME(zhoujiahui): the TTL field will be configurable in the future
	return writeHttpResponse(w, r)
}

func (r *httpResponse) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func writeHttpResponse(w http.ResponseWriter, obj interface{}) (err error) {
	var jsonBytes []byte
	writeContentType(w, jsonContentType)
	if jsonBytes, err = json.Marshal(obj); err != nil {
		err = errors.WithStack(err)
		return
	}
	if _, err = w.Write(jsonBytes); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//
// Middleware Error Handler in server package
//

// func (e *httpResponse) Code() int { return e.code.Code() }

func (e *httpResponse) Error() string {
	return "code is " + strconv.Itoa(e.Code.Code()) + " reason:" + e.Message + " codeMessage:" + e.Code.Message()
}

func (e *httpResponse) Details() []interface{} {
	return nil
}

// func (e *httpResponse) Message() string {
// 	return e.message
// }

func (e *httpResponse) MarshalJSON() (jsonstr []byte, err error) {

	if e.HttpCode != 200 {
		mapTemp := map[string]string{}
		mapTemp["code"] = strconv.Itoa(e.Code.Code())
		mapTemp["reason"] = e.Message
		mapTemp["innerReason"] = e.Code.Message()

		keyCount := 1
		for _, err := range e.Errors {
			if err == nil {
				continue
			}
			errKeyCode := "detailError_" + strconv.Itoa(keyCount)
			mapTemp[errKeyCode] = err.Error()
		}
		jsonstr, err = json.Marshal(mapTemp)

	} else {
		// responseData, _ := (*e.Response)["data"].(map[string]interface{})

		// if len(responseData) == 0 {
		// 	jsonstr, err = json.Marshal(map[string]string{"status": "success"})
		// } else if len(responseData) == 1 {
		// 	for _, value := range responseData {
		// 		jsonstr, err = json.Marshal(map[string]interface{}{"status": "success", "data": value})
		// 		break
		// 	}
		// } else {
		// 	jsonstr, err = json.Marshal(e.Response)
		// }
		jsonstr, err = json.Marshal(e.Response)
	}
	return
}

//AddErr 添加错误列表
func (e *httpResponse) AddErr(err ecode.Code) *httpResponse {
	e.Errors = append(e.Errors, err)
	return e
}

//NewhttpResponse 初始化一个httpResponse
func NewhttpResponse(res *Response) *httpResponse {

	httpResponse := &httpResponse{Code: ecode.OK, HttpCode: 200, Response: res}

	return httpResponse
}

//NewInnerError 内部错误
func NewInnerError(mode constants.ServerMode, infoCode ecode.Code, msg string, err error) *httpResponse {
	httpResponse := &httpResponse{Code: infoCode, HttpCode: 500}

	if err != nil && mode == constants.Dev {
		httpResponse.Errors = append(httpResponse.Errors, err)
	}
	httpResponse.Message = msg
	return httpResponse
}

func NewClientError(mode constants.ServerMode, infoCode ecode.Code, msg string, err error) *httpResponse {
	httpResponse := &httpResponse{Code: infoCode, HttpCode: 400}

	if err != nil && mode == constants.Dev {
		httpResponse.Errors = append(httpResponse.Errors, err)
	}

	httpResponse.Message = msg
	return httpResponse
}
