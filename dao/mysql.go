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
	typ := reflect.TypeOf(value)
	if _, ok := typ.FieldByName("CreatedBy"); !ok {
		return errors.New("model is not a bfr micro models")
	} else {
		val := reflect.ValueOf(value).FieldByName("CreatedBy")
		val.SetString(createdBy)
	}
	return db.mysqlClient.Omit(clause.Associations).Model(model).Create(value).Error
}

// Save 保存更新
func (db *DB) Save(value interface{}) error {
	return db.mysqlClient.Omit(clause.Associations).Save(value).Error
}

// Updates 更新模型
func (db *DB) Updates(model interface{}, UpdatesBy string, where interface{}, value interface{}) error {
	return db.mysqlClient.Omit(clause.Associations).Model(model).Where(where).Updates(value).Error
}

//DeleteByWhere 条件删除
func (db *DB) DeleteByWhere(model, where interface{}) (count int64, err error) {
	data := db.mysqlClient.Where(where).Omit(clause.Associations).Delete(model)

	if data.Error != nil {
		return
	}
	count = data.RowsAffected
	return
}

// DeleteByID 根据ID删除一行
func (db *DB) DeleteByID(model interface{}, id uint64) (count int64, err error) {
	data := db.mysqlClient.Where("id=?", id).Omit(clause.Associations).Delete(model)
	err = data.Error
	if err != nil {
		return 0, err
	}
	count = data.RowsAffected
	return count, nil
}

// DeleteByIDS 根据id批量删除
func (db *DB) DeleteByIDS(model interface{}, ids []uint64) (count int64, err error) {
	data := db.mysqlClient.Omit(clause.Associations).Where("id in (?)", ids).Delete(model)
	err = data.Error
	if err != nil {
		return 0, err
	}
	count = db.mysqlClient.RowsAffected
	return 0, nil
}

// FirstByID 查找第一个ID的数据
func (db *DB) FirstByID(out interface{}, id int) (notFound bool, err error) {
	err = db.mysqlClient.First(out, id).Error
	if err != nil {
		notFound = errors.Is(err, gorm.ErrRecordNotFound)
	}
	return
}

// First 符合条件的第一行
func (db *DB) First(model, where interface{}, out interface{}) (notFound bool, err error) {
	err = db.mysqlClient.Model(model).Where(where).First(out).Error
	if err != nil {
		notFound = errors.Is(err, gorm.ErrRecordNotFound)
		if notFound {
			err = nil
		}
		return
	}
	return
}

// Find 根据条件查询到的数据
func (db *DB) Find(model, where interface{}, out interface{}, orders ...string) error {
	data := db.mysqlClient.Model(model).Where(where)
	if len(orders) > 0 {
		for _, order := range orders {
			data = data.Order(order)
		}
	}
	return data.Find(out).Error
}

// Scan where为sql语句，model为自定义的结果模型数据结构 将结果扫描到数据结构中
func (db *DB) Scan(model, where interface{}, out interface{}) (notFound bool, err error) {
	err = db.mysqlClient.Model(model).Where(where).Scan(out).Error
	if err != nil {
		notFound = errors.Is(err, gorm.ErrRecordNotFound)
		return notFound, err
	}
	return false, nil
}

// ScanList where为sql  将结果扫描到数据结构中 并根据orders排序
func (db *DB) ScanList(model, where interface{}, out interface{}, orders ...string) error {
	data := db.mysqlClient.Model(model).Where(where)
	if len(orders) > 0 {
		for _, order := range orders {
			data = data.Order(order)
		}
	}
	return data.Scan(out).Error
}

// GetPage 从数据库中分页获取数据
func (db *DB) GetPage(model, where, out interface{}, pageIndex, pageSize int, totalCount *int64, autoLoad bool, whereOrder ...PageWhereOrder) error {
	var data *gorm.DB
	if autoLoad {
		data = db.mysqlClient.Preload(clause.Associations).Model(model)
	} else {
		data = db.mysqlClient.Model(model)
	}
	if where != nil {
		data = data.Where(where)
	}
	if whereOrder != nil && len(whereOrder) > 0 {
		for _, wo := range whereOrder {
			if wo.Order != "" {
				data = data.Order(wo.Order)
			}
			if wo.Where != "" {
				data = data.Where(wo.Where, wo.Value...)
			}
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
func (db *DB) GetPageWithFilters(model interface{}, filters *Filter, out interface{}, pageIndex, pageSize int, totalCount *int64, autoLoad bool, whereOrder ...PageWhereOrder) error {

	var data *gorm.DB = db.mysqlClient.Debug()

	if autoLoad {
		data = data.Preload(clause.Associations).Model(model)
	} else {
		data = data.Model(model)
	}
	if filters != nil {
		for _, item := range filters.FilterItems {
			condition, cri, err := item.WhereValue(model)
			if err != nil {
				return err
			}
			data = data.Where(condition, cri)
		}
	}
	if whereOrder != nil && len(whereOrder) > 0 {
		for _, wo := range whereOrder {
			if wo.Order != "" {
				data = data.Order(wo.Order)
			}
			if wo.Where != "" {
				data = data.Where(wo.Where, wo.Value...)
			}
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
	err = data.Model(model).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error

	return err
}

//GetPageByRaw 根据原始的sql进行分页查询
func (db *DB) GetPageByRaw(sql string, out interface{}, pageIndex, pageSize int, totalCount *int64, where ...interface{}) error {

	data := db.mysqlClient.Raw(sql, where...)

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
	return db.mysqlClient.Model(model).Where(where).Pluck(fieldName, out).Error
}

func (db *DB) Delete(model interface{}, deletedBy string, where interface{}, filters ...interface{}) error {
	var op *gorm.DB = db.mysqlClient.Model(model)

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
