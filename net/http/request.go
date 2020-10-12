package http

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/avayayu/micro/dao"
)


type HttpRequest interface {
	GetPageParameter() (int, int, *dao.Order, error)
	GetPageFilters(model interface{}) (*dao.Filter, error)
}

//GetPageParameter 从GET参数中获取分页参数以及排序参数
func (c *Context) GetPageParameter() (int, int, *dao.Order, error) {

	page := c.Query("page")
	perPage := c.Query("perPage")
	orderJSON := c.Query("orders")

	var order dao.Order
	if orderJSON != "" {
		if err := json.Unmarshal([]byte(orderJSON), &order); err != nil {
			return 0, 0, nil, err
		}

	}

	pageNum, err := strconv.Atoi(page)

	if err != nil {
		return 0, 0, nil, errors.New("page must be int number")
	}
	perPageInt, err := strconv.Atoi(perPage)

	if err != nil {
		return 0, 0, nil, errors.New("per_page must be int number")
	}
	return perPageInt, pageNum, &order, nil
}

//GetPageParameter 从GET参数中获取分页参数以及排序参数
func (c *Context) GetPageFilters(model interface{}) (*dao.Filter, error) {
	filtersJSON := c.QueryArray("filters[]")
	filters := dao.Filter{}

	if len(filtersJSON) == 0 {
		return nil, nil
	}

	for _, filter := range filtersJSON {
		filerItem := dao.FilterItem{}
		if err := json.Unmarshal([]byte(filter), &filerItem); err != nil {
			return nil, err
		}
		if !dao.JudgeFilters(model, filerItem.Column, filerItem.FilterType) {
			return nil, errors.New("the field can not used to filter")
		}
		filters.FilterItems = append(filters.FilterItems, &filerItem)

	}
	return &filters, nil
}
