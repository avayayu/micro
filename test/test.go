package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/avayayu/micro/net/http"
	"github.com/stretchr/testify/assert"
)

// ParseToStr 将map中的键值对输出成querystring形式
func ParseToStr(mp map[string]string) string {
	values := ""
	for key, val := range mp {
		values += "&" + key + "=" + val
	}
	temp := values[1:]
	values = "?" + temp
	return values
}

// Get 根据特定请求uri，发起get请求返回响应
func Get(t *testing.T, uri string, router *http.Engine) []byte {
	// 构造get请求
	req := httptest.NewRequest("GET", uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}

// PostForm 根据特定请求uri和参数param，以表单形式传递参数，发起post请求返回响应
func PostForm(t *testing.T, uri string, param map[string]string, router *http.Engine) []byte {
	// 构造post请求，表单数据以querystring的形式加在uri之后
	req := httptest.NewRequest("POST", uri+ParseToStr(param), nil)

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应handler接口
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}

// PostJson 根据特定请求uri和参数param，以Json形式传递参数，发起post请求返回响应
func PostJson(t *testing.T, uri string, param interface{}, router *http.Engine) []byte {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(param)

	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}