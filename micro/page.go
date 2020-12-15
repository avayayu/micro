package micro

import (
	"errors"
	"reflect"

	"github.com/avayayu/micro/dao"
	"github.com/avayayu/micro/net/http"
)

//PagesQuery 根据请求与模型 进行资源的分页查询/排序/过滤
//parameter 如果不实现FilterModels接口则将直接使用该模型进行ORM访问
func PagesQuery(parameter interface{}, out interface{}, db dao.DAO, request http.HttpRequest, response http.Response, wheres ...*dao.QueryOptions) (totalCount int64, page int, perPage int, err error) {

	if reflect.TypeOf(parameter).Kind() != reflect.Ptr {
		err = errors.New("models must be ptr")
		return
	}
	var dataModel interface{}
	if rmodel, ok := parameter.(dao.FilterModels); ok {
		dataModel = rmodel.OrmModels()
	} else {
		dataModel = parameter
	}

	perPage, page, rawOrder, err := request.GetPageParameter()

	if err != nil {
		return
	}

	filters, err := request.GetPageFilters(parameter)
	if err != nil {
		return
	}

	var query *dao.QueryOptions
	if rawOrder != nil {
		query = rawOrder.GetPageOrder(parameter)
	}

	query = query.Filter(parameter, filters)

	if wheres != nil {
		wheres = append(wheres, query)
	}
	err = query.GetPageWithFilters(dataModel, filters, out, page, perPage, &totalCount, query)
	return
}

//PagesQueryRaw test2222
func PagesQueryRaw(rawSql string, out interface{}, db dao.DAO, request http.HttpRequest, response http.Response) (totalCount int64, page int, perPage int, err error) {
	perPage, page, _, err = request.GetPageParameter()

	if err != nil {
		return
	}

	err = db.NewQuery().GetPageByRaw(rawSql, out, page, perPage, &totalCount)
	// response.SetPagesData(totalCount, page, perPage, datas)
	return
}
