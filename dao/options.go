package dao

import (
	"bytes"
	"context"
	"fmt"
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
	order []string
	where string
	//条件可能是一个列表 比如 where 为 id in (?)
	conditions    []interface{}
	selectList    []string
	joinTableList []string
	preloadList   []string
	Ctx           context.Context
	session       *gorm.DB
}

func (options *QueryOptions) Model(model interface{}) Query {
	options.session = options.session.Model(model)
	return options
}

func (options *QueryOptions) WhereQuery(where string, conditions ...interface{}) Query {
	options.where = where
	options.conditions = conditions
	return options
}

func (options *QueryOptions) Select(attrs ...string) Query {
	options.selectList = append(options.selectList, attrs...)
	return options
}

func (options *QueryOptions) Joins(Table ...string) Query {
	options.joinTableList = append(options.joinTableList, Table...)
	return options
}

func (options *QueryOptions) Order(order ...string) Query {
	options.order = append(options.order, order...)
	return options
}

func (options *QueryOptions) PreLoad(Attrs ...string) Query {
	options.preloadList = append(options.preloadList, Attrs...)
	return options
}

//Filter 将对外的参数转化为我们可以识别的DB Where语句
func (options *QueryOptions) Filter(parameter interface{}, filter *Filter) Query {
	if filter != nil {
		for _, item := range filter.FilterItems {
			condition, cri, err := item.WhereValue(parameter)
			if err != nil {
				continue
			}
			options.session = options.session.Where(condition, cri)
		}
	}
	return options
}

//ParseOrder 将从请求中剥离的Order参数转化为认识的db的排序
//规则为 tag db:"Column:列名" 如果没有添加标记 则会自动跳过
func (options *QueryOptions) ParseOrder(parameter interface{}, order *Order) Query {
	if reflect.TypeOf(parameter).Kind() != reflect.Ptr {
		panic("models must be a ptr")
	}

	tableName := lib.GetTypeFullName(parameter)

	orderMux.Lock()
	defer orderMux.Unlock()
	if JSONColumn == nil {
		JSONColumn = make(map[string]map[string]string)
	}
	var mapData map[string]string = map[string]string{}
	if _, ok := JSONColumn[tableName]; !ok {
		rType := reflect.TypeOf(parameter).Elem()
		for i := 0; i < rType.NumField(); i++ {
			t := rType.Field(i)
			jsonKey := t.Tag.Get("json")
			if jsonKey == "-" {
				continue
			}
			column := t.Tag.Get("db")
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

	if order.Orders == nil {
		return nil
	}
	for _, orderItem := range order.Orders {
		if realColumn, ok := mapData[orderItem.Column]; ok {
			options.Order(fmt.Sprintf("%s %s", realColumn, orderItem.OrderType.String()))
		}
	}
	return options
}

func (options *QueryOptions) parseQuery(session *gorm.DB) *gorm.DB {

	if options.Ctx != nil {
		session = session.WithContext(options.Ctx)
	}

	for _, table := range options.joinTableList {
		session = session.Joins(table)
	}

	for _, attr := range options.preloadList {
		session = session.Preload(attr)
	}

	if len(options.order) != 0 {
		for _, order := range options.order {
			session = session.Order(order)
		}
	}
	if options.where != "" {
		session = session.Where(options.where, options.conditions...)
	}

	if len(options.selectList) > 0 {
		var buf bytes.Buffer
		for _, col := range options.selectList {
			buf.WriteString(col)
			buf.WriteByte(',')
		}
		cols := buf.String()
		cols = cols[0 : len(cols)-1]
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
func (order *Order) GetPageOrder(models interface{}, pageOrder *QueryOptions) *QueryOptions {
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
			column := t.Tag.Get("db")
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
	if order.Orders == nil {
		return nil
	}
	for _, orderItem := range order.Orders {
		pageOrder := QueryOptions{}
		if realColumn, ok := mapData[orderItem.Column]; ok {
			pageOrder.Order(fmt.Sprintf("%s %s", realColumn, orderItem.OrderType.String()))
		}
	}
	return pageOrder
}
