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
	"github.com/avayayu/micro/logging"
	_ "github.com/go-sql-driver/mysql" //加载mysql驱动给gorm
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Driver interface {
	Connect(config *DBConfigs) (*gorm.DB, *mongo.Client, error)
}

type DAO interface {
	Connect() DAO
	SetConfig(c *DBConfigs) DAO
	AutoMigrate(models ...interface{}) error
	Create(model interface{}, createdBy string, value interface{}) error
	Updates(model interface{}, updatedBy string, value interface{}, filters ...interface{}) error
	Delete(model interface{}, deletedBy string, filters ...interface{}) error
	First(model, out interface{}, options ...*QueryOptions) (Found bool, err error)
	Find(model, out interface{}, options ...*QueryOptions) error
	Raw(sql string, out interface{}) error
	NewQuery() *QueryOptions
	NewTransaction() *Transactions
	AddSubTransaction(tran *Transactions, subT SubTransactions) *Transactions
	ExecTrans(tran *Transactions) error

	GetPage(model, where, out interface{}, pageIndex, pageSize int, totalCount *int64, autoLoad bool, options ...*QueryOptions) error
	GetPageWithFilters(model interface{}, filters *Filter, out interface{}, pageIndex, pageSize int, totalCount *int64, autoLoad bool, options ...*QueryOptions) error
	GetPageByRaw(sql string, out interface{}, pageIndex, pageSize int, totalCount *int64, where ...interface{}) error

	GetDB() *gorm.DB
	SetLogger(logger *zap.Logger)
}

//Database 数据库管理
type DB struct {
	logger      *zap.Logger
	config      *DBConfigs
	db          *gorm.DB
	mongoClient *mongo.Client
	driver      Driver
}

//NewDatabase 新建数据库连接
//如果path 文件不存在，那么重建数据结构
func NewDatabase(configs *DBConfigs) DAO {

	database := &DB{
		config: configs,
		logger: logging.NewSimpleLogger(),
	}
	return database
}

//SetLogger 设置外部logger
func (db *DB) SetLogger(logger *zap.Logger) {
	db.logger = logger
}

func (db *DB) SetConfig(c *DBConfigs) DAO {
	db.config = c
	return db
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

func (db *DB) NewQuery() *QueryOptions {
	return &QueryOptions{
		conditions:    []interface{}{},
		selectList:    []string{},
		joinTableList: []string{},
		order:         []string{},
	}
}
