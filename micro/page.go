package micro

import (
	"errors"
	"reflect"

	"github.com/avayayu/micro/dao"
	"github.com/avayayu/micro/net/http"
)

//PagesQuery 根据请求与模型 进行资源的分页查询/排序/过滤
func PagesQuery(model interface{}, out interface{}, db dao.DAO, request http.HttpRequest, response http.Response, wheres ...dao.QueryOptions) (totalCount int64, page int, perPage int, err error) {

	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		err = errors.New("models must be ptr")
		return
	}
	var dataModel interface{}
	if rmodel, ok := model.(dao.FilterModels); ok {
		dataModel = rmodel.OrmModels()
	} else {
		dataModel = model
	}

	perPage, page, rawOrder, err := request.GetPageParameter()

	if err != nil {
		return
	}

	filters, err := request.GetPageFilters(model)
	if err != nil {
		return
	}

	var orders []dao.QueryOptions
	if rawOrder != nil {
		orders = rawOrder.GetPageOrder(model)
	}

	if wheres != nil {
		orders = append(orders, wheres...)
	}
	err = db.GetPageWithFilters(dataModel, filters, out, page, perPage, &totalCount, true, orders...)
	return
}

func PagesQueryRaw(rawSql string, out interface{}, db dao.DAO, request http.HttpRequest, response http.Response) (totalCount int64, page int, perPage int, err error) {
	perPage, page, _, err = request.GetPageParameter()

	if err != nil {
		return
	}

	err = db.GetPageByRaw(rawSql, out, page, perPage, &totalCount)
	// response.SetPagesData(totalCount, page, perPage, datas)
	return
}
