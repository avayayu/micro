package dao

import (
	"bytes"
	"reflect"
	"strings"
	"sync"

	"github.com/avayayu/micro/lib"
)

var JSONColumn map[string]map[string]string
var orderMux sync.Mutex

//PageWhereOrder 分页排序
type PageWhereOrder struct {
	Order string
	Where string
	//条件可能是一个列表 比如 where 为 id in (?)
	Value []interface{}
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
func (order *Order) GetPageOrder(models interface{}) []PageWhereOrder {
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
	var pageOrders []PageWhereOrder
	if order.Orders == nil {
		return nil
	}
	for _, orderItem := range order.Orders {
		pageOrder := PageWhereOrder{}
		if realColumn, ok := mapData[orderItem.Column]; ok {
			pageOrder.Order = realColumn + " " + orderItem.OrderType.String()
		}
		pageOrders = append(pageOrders, pageOrder)
	}

	return pageOrders
}
