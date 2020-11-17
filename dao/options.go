package dao

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"sync"

	"github.com/avayayu/micro/lib"
	"gorm.io/gorm"
)

var JSONColumn map[string]map[string]string
var orderMux sync.Mutex

//QueryOptions 分页排序
type QueryOptions struct {
	Order string
	Where string
	//条件可能是一个列表 比如 where 为 id in (?)
	Conditions    []interface{}
	PluckList     []string
	JoinTableList []string
	Ctx           context.Context
}

func (options *QueryOptions) WhereQuery(where string, conditions ...interface{}) *QueryOptions {
	options.Where = where
	options.Conditions = conditions
	return options
}

func (options *QueryOptions) Pluck(attrs ...string) *QueryOptions {
	options.PluckList = append(options.PluckList, attrs...)
	return options
}

func (options *QueryOptions) Joins(Table ...string) *QueryOptions {
	options.JoinTableList = append(options.JoinTableList, Table...)
	return options
}

func (options *QueryOptions) ParseQuery(session *gorm.DB) *gorm.DB {

	if options.Ctx != nil {
		session = session.WithContext(options.Ctx)
	}

	for _, table := range options.JoinTableList {
		session = session.Joins(table)
	}

	if options.Order != "" {
		session = session.Order(options.Order)
	}
	if options.Where != "" {
		session = session.Where(options.Where, options.Conditions...)
	}

	if len(options.PluckList) > 0 {
		var buf bytes.Buffer
		for _, col := range options.PluckList {
			buf.WriteString(col)
			buf.WriteByte(',')
		}
		cols := buf.String()
		cols = cols[1 : len(cols)-1]
		session = session.Select(cols)
	}
	return session
}

//OrderType 用于前端排序
type OrderType uint

const (
	_ OrderType = iota
	Asc
	Desc
)

func (orderType OrderType) String() string {
	if orderType == Asc {
		return "asc"
	} else if orderType == Desc {
		return "desc"
	} else {
		panic("not orderType")
	}
}

//OrderItem 排序项
type OrderItem struct {
	Column    string    `json:"column"`
	OrderType OrderType `json:"orderType"`
}

//Order 排序
type Order struct {
	Orders []OrderItem `json:"orders"`
}

func (order *Order) String() string {
	var buf bytes.Buffer

	for _, orderItem := range order.Orders {
		buf.WriteString(orderItem.Column)
		buf.WriteString(" ")
		buf.WriteString(orderItem.OrderType.String())
		buf.WriteString(",")
	}

	sqlBytes := buf.Bytes()
	if len(sqlBytes) > 0 {
		return string(sqlBytes[0 : len(sqlBytes)-1])
	} else {
		return ""
	}
}

//GetPageOrder 构造GetPageOrder
func (order *Order) GetPageOrder(models interface{}) []QueryOptions {
	if reflect.TypeOf(models).Kind() != reflect.Ptr {
		panic("models must be a ptr")
	}

	tableName := lib.GetTypeFullName(models)

	orderMux.Lock()
	defer orderMux.Unlock()
	if JSONColumn == nil {
		JSONColumn = make(map[string]map[string]string)
	}
	var mapData map[string]string = map[string]string{}
	if _, ok := JSONColumn[tableName]; !ok {
		rType := reflect.TypeOf(models).Elem()
		for i := 0; i < rType.NumField(); i++ {
			t := rType.Field(i)
			jsonKey := t.Tag.Get("json")
			if jsonKey == "-" {
				continue
			}
			column := t.Tag.Get("gorm")
			if column != "" {
				gormArr := strings.Split(column, ";")
				for _, field := range gormArr {
					if strings.Contains(field, "Column") {
						fieldArray := strings.Split(field, ":")
						mapData[jsonKey] = fieldArray[1]
					}
				}
			}
		}
		JSONColumn[tableName] = mapData
	}
	mapData = JSONColumn[tableName]
	var pageOrders []QueryOptions
	if order.Orders == nil {
		return nil
	}
	for _, orderItem := range order.Orders {
		pageOrder := QueryOptions{}
		if realColumn, ok := mapData[orderItem.Column]; ok {
			pageOrder.Order = realColumn + " " + orderItem.OrderType.String()
		}
		pageOrders = append(pageOrders, pageOrder)
	}

	return pageOrders
}
