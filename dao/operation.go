package dao

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	ztime "github.com/avayayu/micro/time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Create 创建某模型一行
func (query *QueryOptions) Create(model interface{}, createdBy string, value interface{}) error {

	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		panic("value must be ptr")
	}

	typ := reflect.TypeOf(value).Elem()
	switch typ.Kind() {
	case reflect.Struct:
		if _, ok := typ.FieldByName("CreatedBy"); ok {
			val := reflect.ValueOf(value).Elem().FieldByName("CreatedBy")
			val.SetString(createdBy)
		}
		return query.session.Omit(clause.Associations).Model(model).Create(value).Error
	case reflect.Slice, reflect.Array:
		sliceValue := reflect.ValueOf(value).Elem()
		for i := 0; i < sliceValue.Len(); i++ {
			v := sliceValue.Index(i)
			typV := v.Type()
			if typV.Kind() == reflect.Ptr {
				typV = typV.Elem()
				v = v.Elem()
			}

			if typV.Kind() != reflect.Struct {
				panic("element must be struct ")
			}

			if _, ok := typV.FieldByName("CreatedBy"); ok {
				val := v.FieldByName("CreatedBy")
				val.SetString(createdBy)
			}

			if err := query.session.Omit(clause.Associations).Model(model).Create(v.Addr().Interface()).Error; err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("slice or struct needed,but type of value is %s", reflect.TypeOf(value).Kind())
	}

	return nil

}

// Save 保存更新
func (query *QueryOptions) Save(value interface{}) error {
	return query.session.Omit(clause.Associations).Save(value).Error
}

//Count 根据querys中的where进行数量的统计
func (query *QueryOptions) Count(model interface{}) (count int64) {

	session := query.session.Model(model)
	session = session.Where(query.where, query.conditions...)
	session.Count(&count)
	return

}

// Updates 更新模型
func (query *QueryOptions) Updates(model interface{}, UpdatesBy string, value interface{}, filters ...interface{}) error {
	session := query.session.Omit(clause.Associations).Model(model)

	if query.where != "" {
		session = session.Where(query.where, query.conditions...)
	}

	if len(filters) > 0 {

		if len(filters)%2 != 0 {
			panic("filters must be odd")
		}

		for i := 0; i < len(filters); i += 2 {
			session = session.Where(filters[i], filters[i+1])
		}
	}

	if query.where == "" && len(filters) == 0 {
		log.Fatalln("updates data in no condition")
	}

	return session.Update("updated_by", UpdatesBy).Updates(value).Error
}

// First 符合条件的第一行
func (query *QueryOptions) First(model, out interface{}) (Found bool, err error) {
	op := query.session.Model(model)

	op = query.parseQuery(op)

	err = op.First(out).Error
	if err != nil {
		notFound := errors.Is(err, gorm.ErrRecordNotFound)
		if notFound {
			Found = false
			err = nil
		}
		return
	}
	return true, nil
}

func (query *QueryOptions) Raw(sql string, out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return query.session.Raw(sql).Scan(out).Error
}

// Find 根据条件查询到的数据
func (query *QueryOptions) Find(model, out interface{}) error {
	return query.parseQuery(query.session.Model(model)).Find(out).Error
}

//Update 更新单列数据
func (query *QueryOptions) Update(model, column string, value interface{}) error {
	return query.parseQuery(query.session.Model(model)).Update(column, value).Error
}

//FindToMap 将查询结果存放到map中，其中Column为作为key的列
//如果Column不是主键将会自动覆盖
func (query *QueryOptions) FindToMap(model, out interface{}, column string) error {
	typ := reflect.TypeOf(out)
	outValue := reflect.ValueOf(out).Elem()
	if typ.Kind() != reflect.Ptr {
		log.Fatal("out must be pointer")
		return errors.New("out must be pointer")
	}

	if typ.Elem().Kind() != reflect.Map {
		log.Fatal("out must be a pointer of golang Map")
		return errors.New("out must be pointer")
	}

	typ = typ.Elem().Elem()
	if (typ.Kind() != reflect.Ptr && typ.Kind() != reflect.Struct) || typ.Kind() == reflect.Ptr && typ.Elem().Kind() != reflect.Struct {
		log.Fatal("element of out map must be a struct")
		return errors.New("element of out map must be a struct")
	}
	// slice := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)

	slice := reflect.New(reflect.SliceOf(typ))
	sliceData := slice.Interface()

	if err := query.Find(model, sliceData); err != nil {
		return err
	}
	sliceValue := reflect.ValueOf(sliceData).Elem()
	for i := 0; i < sliceValue.Len(); i++ {
		value := sliceValue.Index(i)
		var key reflect.Value
		if value.Type().Kind() == reflect.Ptr {
			valueStruct := value.Elem()
			key = valueStruct.FieldByName(column)
		} else {
			key = value.FieldByName(column)
		}

		outValue.SetMapIndex(key, value)
	}

	return nil
}

// GetPage 从数据库中分页获取数据
func (query *QueryOptions) GetPage(model, out interface{}, pageIndex, pageSize int, totalCount *int64) error {
	var session *gorm.DB = query.parseQuery(query.session)

	err := session.Model(model).Count(totalCount).Error

	if err != nil {
		return err
	}
	if *totalCount == 0 {
		return nil
	}
	return session.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error
}

// GetPage 从数据库中分页获取数据
func (query *QueryOptions) GetPageWithFilters(parameter interface{}, filters *Filter, out interface{}, pageIndex, pageSize int, totalCount *int64) error {

	var session *gorm.DB = query.session

	if filters != nil {
		for _, item := range filters.FilterItems {
			condition, cri, err := item.WhereValue(parameter)
			if err != nil {
				return err
			}
			session = session.Where(condition, cri)
		}
	}
	session = query.parseQuery(session)
	model := parameter
	if rmodel, ok := parameter.(FilterModels); ok {
		model = rmodel.OrmModels()
	}

	// session = session.Order("updated_at desc").Order("created_at desc")

	err := session.Model(model).Count(totalCount).Error

	if err != nil {
		return err
	}
	if *totalCount == 0 {
		return nil
	}
	err = session.Model(model).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error

	return err
}

//GetPageByRaw 根据原始的sql进行分页查询
func (query *QueryOptions) GetPageByRaw(sql string, out interface{}, pageIndex, pageSize int, totalCount *int64, where ...interface{}) error {

	data := query.session.Raw(sql, where...)

	err := data.Count(totalCount).Error

	if err != nil {
		return err
	}
	if *totalCount == 0 {
		return nil
	}
	return data.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error

}

//PluckList 查询某表中的某一列 切片
func (query *QueryOptions) PluckList(model, out interface{}, fieldName string) error {
	return query.parseQuery(query.session.Model(model)).Pluck(fieldName, out).Error
}

func (query *QueryOptions) Delete(model interface{}, deletedBy string, filters ...interface{}) error {
	var op *gorm.DB = query.session.Model(model)

	if len(filters)%2 != 0 {
		panic("filters length must be even")
	}

	for i := 0; i < len(filters); i = i + 2 {
		op = op.Where(filters[i], filters[i+1])
	}
	return op.Updates(map[string]interface{}{"deleted_at": ztime.Now(), "deleted_by": deletedBy}).Error
}

//CheckError 检查错误是否为数据不存在
func CheckError(err error) (bool, error) {
	if err == gorm.ErrRecordNotFound {
		return true, nil
	}
	return false, err
}
