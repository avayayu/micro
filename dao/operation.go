package dao

import (
	"errors"
	"reflect"

	ztime "github.com/avayayu/micro/time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Create 创建某模型一行
func (db *DB) Create(model interface{}, createdBy string, value interface{}) error {

	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		panic("value must be ptr")
	}

	typ := reflect.TypeOf(value).Elem()
	if _, ok := typ.FieldByName("CreatedBy"); !ok {
		return errors.New("model is not a bfr micro models")
	} else {
		val := reflect.ValueOf(value).Elem().FieldByName("CreatedBy")
		val.SetString(createdBy)
	}
	return db.db.Omit(clause.Associations).Model(model).Create(value).Error

}

// Save 保存更新
func (db *DB) Save(value interface{}) error {
	return db.db.Omit(clause.Associations).Save(value).Error
}

//Count 根据querys中的where进行数量的统计
func (db *DB) Count(model interface{}, querys ...*QueryOptions) (count int64) {

	session := db.db.Model(model)
	for _, query := range querys {
		session = session.Where(query.where, query.conditions...)
	}

	session.Count(&count)
	return

}

// Updates 更新模型
func (db *DB) Updates(model interface{}, UpdatesBy string, value interface{}, filters ...interface{}) error {
	session := db.db.Omit(clause.Associations).Model(model)
	if len(filters) > 0 {

		if len(filters)%2 != 0 {
			panic("filters must be odd")
		}

		for i := 0; i < len(filters); i += 2 {
			session = session.Where(filters[i], filters[i+1])
		}
	} else {
		db.logger.Warn("updates data in no condition")
	}
	return session.Update("updated_by", UpdatesBy).Updates(value).Error
}

// First 符合条件的第一行
func (db *DB) First(model, out interface{}, options ...*QueryOptions) (Found bool, err error) {
	op := db.db.Model(model)
	for _, option := range options {
		op = option.ParseQuery(op)
	}
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

func (db *DB) Raw(sql string, out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return db.db.Raw(sql).Scan(out).Error
}

// Find 根据条件查询到的数据
func (db *DB) Find(model, out interface{}, options ...*QueryOptions) error {
	op := db.db.Model(model)
	for _, option := range options {
		op = option.ParseQuery(op)
	}
	return op.Scan(out).Error
}

// GetPage 从数据库中分页获取数据
func (db *DB) GetPage(model, where, out interface{}, pageIndex, pageSize int, totalCount *int64, options ...*QueryOptions) error {
	var data *gorm.DB

	if where != nil {
		data = data.Where(where)
	}
	if options != nil && len(options) > 0 {
		for _, option := range options {
			data = option.ParseQuery(data)
		}
	} else {
		data = data.Order("updated_at desc").Order("created_at desc")
	}
	err := data.Count(totalCount).Error

	if err != nil {
		return err
	}
	if *totalCount == 0 {
		return nil
	}
	return data.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error
}

// GetPage 从数据库中分页获取数据
func (db *DB) GetPageWithFilters(model interface{}, filters *Filter, out interface{}, pageIndex, pageSize int, totalCount *int64, options ...*QueryOptions) error {

	var data *gorm.DB = db.db

	if filters != nil {
		for _, item := range filters.FilterItems {
			condition, cri, err := item.WhereValue(model)
			if err != nil {
				return err
			}
			data = data.Where(condition, cri)
		}
	}

	for _, option := range options {
		data = option.ParseQuery(data)
	}

	data = data.Order("updated_at desc").Order("created_at desc")

	err := data.Count(totalCount).Error

	if err != nil {
		return err
	}
	if *totalCount == 0 {
		return nil
	}
	err = data.Model(model).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error

	return err
}

//GetPageByRaw 根据原始的sql进行分页查询
func (db *DB) GetPageByRaw(sql string, out interface{}, pageIndex, pageSize int, totalCount *int64, where ...interface{}) error {

	data := db.db.Raw(sql, where...)

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
func (db *DB) PluckList(model, where interface{}, out interface{}, fieldName string) error {
	return db.db.Model(model).Where(where).Pluck(fieldName, out).Error
}

func (db *DB) Delete(model interface{}, deletedBy string, filters ...interface{}) error {
	var op *gorm.DB = db.db.Model(model)

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

func (db *DB) NewTransaction() *Transactions {

	trans := Transactions{
		session: db.db.Session(&gorm.Session{SkipDefaultTransaction: true, FullSaveAssociations: false}),
	}
	return &trans
}

func (db *DB) AddSubTransaction(tran *Transactions, subT SubTransactions) *Transactions {
	tran.subTransactions = append(tran.subTransactions, subT)
	return tran
}

func (db *DB) ExecTrans(tran *Transactions) error {
	return tran.Run()
}
