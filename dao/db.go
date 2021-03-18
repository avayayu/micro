package dao

/*
 * @Author: 邹航宇
 * @Date: 2019-12-30 12:21:28
 * @LastEditTime : 2019-12-31 17:26:22
 * @LastEditors  : 邹航宇
 * @Description: 本模块
 * @输出一段不带属性的自定义信息
 */
import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gogs.bfr.com/zouhy/micro/logging"
	"gogs.bfr.com/zouhy/micro/models"
	"gorm.io/gorm"
)

type DAO interface {
	Connect() DAO
	AutoMigrate(models ...interface{}) error
	NewQuery() Query
	GetDB() *gorm.DB
	GetMongo() *mongo.Client
	SetLogger(logger *zap.Logger)
	NewTrainsactions() Transactions
}

//Transactions 事务封装接口
type Transactions interface {
	//提交
	Commit() error
	//执行
	Execute(sub func(query Query) error)
}

type Query interface {
	Model(model interface{}) Query
	Create(model interface{}, createdBy string, value interface{}) error
	Update(model interface{}, column string, value interface{}) error
	Updates(model interface{}, updatedBy string, value interface{}, filters ...interface{}) error
	Delete(model interface{}, deletedBy string, filters ...interface{}) error
	First(model, out interface{}) (Found bool, err error)
	Find(model, out interface{}) error
	Count(model interface{}) (count int64)
	Raw(sql string, out interface{}) error
	FindToMap(model, out interface{}, column string) error
	GetPage(model, out interface{}, pageIndex, pageSize int, totalCount *int64) error
	GetPageWithFilters(model interface{}, filters *Filter, out interface{}, pageIndex, pageSize int, totalCount *int64) error
	ParseOrder(parameter interface{}, order *Order) Query
	Filter(parameter interface{}, filter *Filter) Query
	WhereQuery(where string, conditions ...interface{}) Query
	Select(attrs ...string) Query
	Joins(Table ...string) Query
	Order(order ...string) Query
	PreLoad(Attrs ...string) Query
	PluckList(model interface{}, out interface{}, fieldName string) error
	CheckIDList(model interface{}, idList []models.Int64Str) error
}

//Database 数据库管理
type DB struct {
	logger      *zap.Logger
	db          *gorm.DB
	mongoClient *mongo.Client
	driver      Driver
	dbType      DBType
}

// type DBConfigs struct {
// 	logger
// }

//NewDatabase 新建数据库连接
//如果path 文件不存在，那么重建数据结构
func NewDatabase(driver Driver) DAO {

	database := &DB{
		logger: logging.NewSimpleLogger(),
	}

	db, m, err := driver.Connect()
	if err != nil {
		database.logger.Error("数据库连接错误", zap.Error(err))
	}
	database.db = db
	database.mongoClient = m

	return database
}

//SetLogger 设置外部logger
func (db *DB) SetLogger(logger *zap.Logger) {
	db.logger = logger
}

//GetMongo 检测是否连接没连接 再次连接
func (db *DB) GetMongo() *mongo.Client {
	return db.mongoClient
}

func (db *DB) Connect() DAO {
	//连接mysql

	return db
}

func (db *DB) GetDB() *gorm.DB {
	return db.db
}

func (db *DB) AutoMigrate(models ...interface{}) error {
	return db.db.AutoMigrate(models...)
}

func (db *DB) NewQuery() Query {
	return &QueryOptions{
		conditions:    []interface{}{},
		selectList:    []string{},
		joinTableList: []string{},
		order:         []string{},

		session: db.db.Session(&gorm.Session{SkipDefaultTransaction: true, FullSaveAssociations: false}),
	}
}

func (db *DB) NewTrainsactions() Transactions {

	query := &QueryOptions{
		conditions:    []interface{}{},
		selectList:    []string{},
		joinTableList: []string{},
		order:         []string{},

		session: db.db.Session(&gorm.Session{SkipDefaultTransaction: true, FullSaveAssociations: false}),
	}

	return &transactions{
		query:           query,
		subTransactions: []SubTransactions{},
	}
}
