package micro

import (
	"github.com/avayayu/micro/code"
	"github.com/avayayu/micro/net/http"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

//CheckMysqlError 检查Mysql数据库错误 abort为true时将直接abort 否则返回构造好的 httpResponser
func CheckMysqlError(err error, response http.Response) {

	if err == gorm.ErrRecordNotFound {
		response.FlushHttpClientError(code.HTTPClientErroCode, "数据不存在", err)
	} else {
		response.FlushHttpInnerError(code.HTTPInnerErrorCode, "数据服务出错", err)
	}

}

//CheckMongoDBError 检查MongoDB数据库错误 abort为true时将直接abort 否则返回构造好的 httpResponser
func CheckMongoDBError(err error, response http.Response) {

	if err == mongo.ErrNilDocument || err == mongo.ErrNilCursor || err != mongo.ErrNoDocuments {
		response.FlushHttpClientError(code.HTTPClientErroCode, "数据不存在", err)
	} else {
		response.FlushHttpInnerError(code.HTTPInnerErrorCode, "数据服务出错", err)
	}

}
