package micro

import (
	"errors"
	"reflect"

	"github.com/avayayu/micro/dao"
	"github.com/avayayu/micro/net/http"
)

//PagesQuery 根据请求与模型 进行资源的分页查询/排序/过滤
func PagesQuery(model interface{}, out interface{}, db dao.DAO, request http.HttpRequest, response http.Response, wheres ...dao.PageWhereOrder) (totalCount int64, page int, perPage int, err error) {

	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		err = errors.New("models must be ptr")
		return
	}

	if rmodel, ok := model.(dao.FilterModels); ok {
		model = rmodel.OrmModels()
	}

	perPage, page, rawOrder, err := request.GetPageParameter()

	if err != nil {
		return
	}

	filters, err := request.GetPageFilters(model)
	if err != nil {
		return
	}

	var orders []dao.PageWhereOrder
	if rawOrder != nil {
		orders = rawOrder.GetPageOrder(model)
	}

	if wheres != nil {
		orders = append(orders, wheres...)
	}
	err = db.GetPageWithFilters(model, filters, out, page, perPage, &totalCount, true, orders...)
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
