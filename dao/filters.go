package dao

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

//OrderType 用于前端排序
type FilterType uint

var filterMux sync.Mutex
var filtersColumn map[string]map[string]FilterType

const (
	_        FilterType = iota
	Category            //类别
	Vague               //模糊查询
)

//FilterItem 过滤项
type FilterItem struct {
	Column     string      `json:"key"`
	FilterType FilterType  `json:"filterType"`
	Value      interface{} `json:"value"`
}

func (item *FilterItem) WhereValue(models interface{}) (condition string, criterion interface{}, err error) {
	if reflect.TypeOf(models).Kind() != reflect.Ptr {
		panic("models must be a ptr")
	}

	var tableName string
	if _, ok := reflect.TypeOf(models).Elem().MethodByName("TableName"); !ok {
		panic("models do not have methods TableName()")
	} else {
		tableName = reflect.ValueOf(models).MethodByName("TableName").Call(nil)[0].String()
	}
	modelJSONGormMap(models, tableName)
	mapData := JSONColumn[tableName]

	var buf bytes.Buffer
	if item.FilterType == Category {
		buf.WriteString(mapData[item.Column])
		buf.WriteString(" ")
		switch reflect.TypeOf(item.Value).Kind() {
		case reflect.Array, reflect.Slice:
			slice := reflect.ValueOf(item.Value)
			if slice.Len() == 1 {
				buf.WriteString("=?")
			} else {
				buf.WriteString("in (?)")
			}

		case reflect.Float32, reflect.Float64, reflect.Bool, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			buf.WriteString("=?")
		default:
			err = errors.New("can not be fliters")
			return
		}

		criterion = item.Value

	} else {
		buf.WriteString(mapData[item.Column])
		buf.WriteString(" ")
		buf.WriteString("like")
		buf.WriteString(" ")
		buf.WriteString("(?)")

		criterion = fmt.Sprintf("%s%v%s", "%", item.Value, "%")

	}
	condition = buf.String()
	return
}

//Order 排序
type Filter struct {
	FilterItems []*FilterItem `json:"filters"`
}

func modelJSONGormMap(models interface{}, tableName string) {

	filterMux.Lock()
	defer filterMux.Unlock()
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
}

func RetreiveFilters(models interface{}) map[string]FilterType {
	if reflect.TypeOf(models).Kind() != reflect.Ptr {
		panic("models must be a ptr")
	}
	var tableName string
	if _, ok := reflect.TypeOf(models).Elem().MethodByName("TableName"); !ok {
		panic("models do not have methods TableName()")
	} else {
		tableName = reflect.ValueOf(models).MethodByName("TableName").Call(nil)[0].String()
	}
	filterMux.Lock()
	defer filterMux.Unlock()
	if filtersColumn == nil {
		filtersColumn = make(map[string]map[string]FilterType)
	}

	if tableFilters, ok := filtersColumn[tableName]; ok {
		return tableFilters
	} else {
		filtersColumn[tableName] = make(map[string]FilterType)
	}

	rType := reflect.TypeOf(models).Elem()
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		jsonKey := t.Tag.Get("json")
		if jsonKey == "-" {
			continue
		}
		column := t.Tag.Get("filters")
		if column != "" {
			gormArr := strings.Split(column, ";")
			for _, field := range gormArr {
				if strings.Contains(field, "type") {
					fieldArray := strings.Split(field, ":")

					if fieldArray[1] == "category" {
						filtersColumn[tableName][jsonKey] = Category
					} else if fieldArray[1] == "vague" {
						filtersColumn[tableName][jsonKey] = Vague
					}
				}
			}
		}
	}
	return filtersColumn[tableName]
}

//JudgeFilters 根据模型定义中的filters tag来判断该属性能不能参与筛选
func JudgeFilters(models interface{}, column string, filterType FilterType) (flag bool) {
	filtersJSONMap := RetreiveFilters(models)
	if _, ok := filtersJSONMap[column]; ok {
		if filterType == filtersJSONMap[column] {
			return true
		}
	}
	return
}
