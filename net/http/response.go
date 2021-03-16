package http

import (
	"reflect"

	"github.com/avayayu/micro/code"
	"github.com/avayayu/micro/net/http/render"
)

type Response interface {
	FlushHttpInnerError(code code.Code, msg string, err error)
	FlushHttpClientError(code code.Code, msg string, err error)
	FlushHttpResponse(datas ...interface{})
	FlushFailHttpResponse(reason error)
	SetPagesData(totalNum int64, curPageNum int, numPerPage int, v interface{})
}

/************************************/
/********Buffalo SPECIAL RESPONSE RENDERING ********/
/********AUTHOR:ZouHangYu ********/
/************************************/

//FlushHttpInnerError 返回内部错误
func (c *Context) FlushHttpInnerError(code code.Code, msg string, err error) {
	mode := c.engine.GetServerMode()

	httpResponse := render.NewInnerError(mode, code, msg, err)

	httpCode := 500
	c.WritedStatus = 500
	c.Render(httpCode, httpResponse)
}

//FlushHttpClientError 返回客户错误 是否返回错误细节取决于http服务模式
func (c *Context) FlushHttpClientError(code code.Code, msg string, err error) {
	mode := c.engine.GetServerMode()
	httpResponse := render.NewClientError(mode, code, msg, err)
	httpCode := 400
	c.WritedStatus = 400
	c.Render(httpCode, httpResponse)
}

//FlushHttpResponse 返回一次成功的http请求 datas为 key value格式 其中value一定要为指针！
func (c *Context) FlushHttpResponse(datas ...interface{}) {
	response := render.SuccessResponse()
	if len(datas)%2 != 0 {
		panic("http response must be key value ")
	}
	var index int = 0
	if len(datas) > 0 {
		for {
			if index >= len(datas) {
				break
			}
			if key, ok := datas[index].(string); !ok {
				panic("key must be string")
			} else {
				value := datas[index+1]
				response = response.Set(key, value)
			}
			index = index + 2
		}
	}
	httpResponse := render.NewhttpResponse(response)
	c.WritedStatus = 200
	c.Render(200, httpResponse)
}

//FlushFailHttpResponse 返回未出错但失败的请求的原因
func (c *Context) FlushFailHttpResponse(reason error) {
	response := render.FailResponse()
	response.Reason(reason)
	httpResponse := render.NewhttpResponse(response)
	c.WritedStatus = 200
	c.Render(200, httpResponse)
}

//SetPagesData 构造返回的分页参数以及数据
func (c *Context) SetPagesData(totalNum int64, curPageNum int, numPerPage int, v interface{}) {

	if v == nil {
		result := []interface{}{}
		c.FlushHttpResponse(totalNum, curPageNum, numPerPage, result)
		return
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		panic("v must be a pointer")
	}
	var totalPageNum int64
	totalPageNum = totalNum/(int64(numPerPage)) + 1
	c.FlushHttpResponse("perPage", &numPerPage, "page", &curPageNum, "totalPage", &totalPageNum, "totalNum", &totalNum, "dataList", v)
}
