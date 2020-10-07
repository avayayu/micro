//错误码分为HTTP错误码 和 内部错误码 HTTP错误可以导致前端axios抛出异常，所有的错误处理通过前端catch异常进行处理
//前端可以统一拦截捕获500错误，建立Sentry端收集捕获的500错误
//内部错误码用于内部错误的分类。事实上并不一定真的是错误
//内部错误码为2开头的为业务级错误 错误码为1开头的为系统级错误 1开头的错误为真正的错误，需要把他作为Bug处理掉
//内部错误码为2的为业务错误，用于让前端开发者能顺利的理解错误进行的分类

package code

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/pkg/errors"
)

var (
	_codes    = map[int]struct{}{} // register codes.
	_messages atomic.Value         // NOTE: stored map[int]string
)

//HTTPSuccess 状态码设定
const HTTPSuccess = 200

//HTTPClientErroCode 描述错误原因由客户引起的错误 比如登陆用户名/密码错误等行为
const HTTPClientErroCode = 400

//HTTPInnerErrorCode 描述本不应该存在却出现了的系统错误
const HTTPInnerErrorCode = 500

//HTTPNotRecognizedError 描述本不应该存在却出现了的系统错误
const HTTPNotRecognizedError = 500

//HTTPNonAuthoritative Head不存在Token信息
const HTTPNonAuthoritative = 203

//HTTPUnauthorizedError JWT解析错误 包括过期/Token解析错误/Token格式错误等等
const HTTPUnauthorizedError = 401

//Success 本次请求完完全全成功
const Success = 20000 //本次请求完完全全成功

var (
	//DBNotFound 数据库查询返回为空
	DBNotFound = add(20001, "数据库查询返回为空")

	//DBCreateError 数据库创建错误
	DBCreateError = add(20002, "数据库创建错误")

	//DBQueryError 数据库查询错误
	DBQueryError = add(20003, "数据库查询错误")
	//DBUpdateError 数据库更新错误
	DBUpdateError = add(20004, "数据库更新错误")
	//DBDeleteError 数据库删除错误
	DBDeleteError = add(20005, "数据库删除错误")

	//TypeConverionError 类型转化错误
	TypeConverionError = add(20010, "类型转化错误")

	//JWTErrorInvalid Token解析错误
	JWTErrorInvalid = add(20020, "Token解析错误")

	//JWTErrorExpired Token过期
	JWTErrorExpired = add(20021, "Token过期")

	//JWTErrorNotValidYet Token还没有效
	JWTErrorNotValidYet = add(20022, "Token还没有效")

	//JWTErrorMalformed Token格式错误
	JWTErrorMalformed = add(20023, "Token格式错误")
	//JWTErrorNotFound 请求没有携带Token信息
	JWTErrorNotFound = add(20024, "请求没有携带Token信息")
	//JWT Refresh Token Not Found
	JWTRefreshNotFound = add(20028, "请求必须携带refreshToken")

	VeriCodeGenFailure     = add(20025, "验证码生成错误")
	VeriCodeNotRight       = add(20026, "验证码校验错误")
	VeriCodeParamNotEnough = add(20027, "验证码参数不够")

	//RequestParamInCorrect 请求参数错误
	RequestParamInCorrect = add(20030, "请求参数错误")

	//AliOssError alioss相关错误
	AliOssError = add(20050, "alioss相关错误")
	AliSMSError = add(20051, "aliSMS相关错误")
	AliVODError = add(20052, "alivod相关错误")
	//OS ERROR 操作系统相关错误

	//CreateFileError 创建文件错误
	CreateFileError = add(20060, "创建文件错误")

	//ReportCaculate 报告计算错误
	ReportCaculate    = add(20070, "外骨骼报告计算错误")
	TrainRequestError = add(20071, "上传参数错误")

	//MemcacheError 缓存错误
	MemcacheError   = add(20080, "Memcache缓存错误")
	RedisCacheError = add(20081, "Rediscache缓存错误")

	//DeviceMonitor
	DeviceMonitor        = add(20090, "设备监测出错")
	DeviceAiderMockError = add(20091, "外骨骼设备模拟出错")

	//Router
	RouterContrunctError = add(20100, "计算路由数组错误")
)

var (
	//DBConnectLost 数据库连接丢失
	DBConnectLost = add(10001, "数据库连接丢失")
)



var (
	OK = add(0, "成功") // 正确

	NotModified        = add(-304, "木有改动")    // 木有改动
	TemporaryRedirect  = add(-307, "撞车跳转")    // 撞车跳转
	RequestErr         = add(-400, "请求错误")    // 请求错误
	Unauthorized       = add(-401, "未认证")     // 未认证
	AccessDenied       = add(-403, "访问权限不足")  // 访问权限不足
	NothingFound       = add(-404, "啥都木有")    // 啥都木有
	MethodNotAllowed   = add(-405, "不支持该方法")  // 不支持该方法
	Conflict           = add(-409, "冲突")      // 冲突
	Canceled           = add(-498, "客户端取消请求") // 客户端取消请求
	ServerErr          = add(-500, "服务器错误")   // 服务器错误
	ServiceUnavailable = add(-503, "服务暂不可用")  // 过载保护,服务暂不可用
	Deadline           = add(-504, "服务调用超时")  // 服务调用超时
	LimitExceed        = add(-509, "超出限制")    // 超出限制
)

// Codes code error interface which has a code & message.
type Codes interface {
	// sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	//Detail get error detail,it may be nil.
	Details() []interface{}
}

// A Code is an int error code spec.
type Code int

func (e Code) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Code) Code() int { return int(e) }

// Message return error message
func (e Code) Message() string {
	if cm, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := cm[e.Code()]; ok {
			return msg
		}
	}
	return e.Error()
}

// Details return details.
func (e Code) Details() []interface{} { return nil }

// Int parse code int to error.
func Int(i int) Code { return Code(i) }

// String parse code string to error.
func String(e string) Code {
	if e == "" {
		return OK
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return ServerErr
	}
	return Code(i)
}

// Cause cause from error to code.
func Cause(e error) Codes {
	if e == nil {
		return OK
	}
	ec, ok := errors.Cause(e).(Codes)
	if ok {
		return ec
	}
	return String(e.Error())
}

// Equal equal a and b by code int.
func Equal(a, b Codes) bool {
	if a == nil {
		a = OK
	}
	if b == nil {
		b = OK
	}
	return a.Code() == b.Code()
}

// EqualError equal error
func EqualError(code Codes, err error) bool {
	return Cause(err).Code() == code.Code()
}

func add(e int, message string) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = struct{}{}
	if cm, ok := _messages.Load().(map[int]string); ok {
		if _, ok := cm[e]; ok {
			return Int(e)
		} else {
			cm[e] = message
			Register(cm)
		}
	} else {
		cm := make(map[int]string)
		cm[e] = message
		Register(cm)
	}
	return Int(e)
}

func Register(cm map[int]string) {
	_messages.Store(cm)
}
