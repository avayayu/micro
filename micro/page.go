package micro

import (
	"errors"
	"reflect"

	"github.com/avayayu/micro/dao"
	"github.com/avayayu/micro/net/http"
)

//PagesQuery 根据请求与模型 进行资源的分页查询/排序/过滤
func PagesQuery(models interface{}, out interface{}, db dao.DAO, request http.HttpRequest, response http.Response, wheres ...dao.PageWhereOrder) (totalCount int64, page int, perPage int, err error) {

	if reflect.TypeOf(models).Kind() != reflect.Ptr {
		err = errors.New("models must be ptr")
		return
	}

	perPage, page, rawOrder, err := request.GetPageParameter()

	if err != nil {
		return
	}

	filters, err := request.GetPageFilters(models)
	if err != nil {
		return
	}

	var orders []dao.PageWhereOrder
	if rawOrder != nil {
		orders = rawOrder.GetPageOrder(models)
	}

	if wheres != nil {
		orders = append(orders, wheres...)
	}

	// elemType := reflect.TypeOf(models)
	// fmt.Println(elemType.Name())
	// sliceType := reflect.SliceOf(elemType)
	// slice := reflect.MakeSlice(sliceType, 0, 0)

	// // Create a pointer to a slice value and set it to the slice
	// ptr := reflect.New(sliceType)
	// ptr.Elem().Set(slice)
	// datas = ptr.Interface()
	// fmt.Printf("s is %T\n", datas)
	// fmt.Println(datas)
	// if err = db.GetPageWithFilters(models, filters, datas, page, perPage, &totalCount, true, orders...); err != nil {
	// 	response.FlushHttpInnerError(code.DBQueryError, "查询患者列表出错", err)
	// 	return
	// }

	err = db.GetPageWithFilters(models, filters, out, page, perPage, &totalCount, true, orders...)
	// response.SetPagesData(totalCount, page, perPage, datas)
	return

}
