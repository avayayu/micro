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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

type DAO interface {
	Connect() DAO
	SetMysqlConfig(c *MysqlConfig) DAO
	SetMongoConfig(c *MongoConfig) DAO
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
	GetMongo() *mongo.Client
	GetMysql() *gorm.DB
	SetLogger(logger *zap.Logger)
}

//Database 数据库管理
type DB struct {
	logger       *zap.Logger
	options      *DBOptions
	mysqlClient  *gorm.DB
	mongoClient  *mongo.Client
	mysqlConfigs *MysqlConfig
	mongoConfigs *MongoConfig
}

//NewDatabase 新建数据库连接
//如果path 文件不存在，那么重建数据结构
//configs里 如果处于Debug模式 那么连阿里云RDS外网服务 如果是生产环境则直接连阿里云RDS内网服务
func NewDatabase(options *DBOptions) DAO {
	if options == nil {
		options = &DBOptions{
			Mysql: true,
			Mongo: false,
		}
	}
	database := &DB{
		options: options,
		logger:  logging.NewSimpleLogger(),
	}
	return database
}

func newDatabase(options *DBOptions) *DB {
	if options == nil {
		options = &DBOptions{
			Mysql: true,
			Mongo: false,
		}
	}
	database := &DB{
		options: options,
		logger:  logging.NewSimpleLogger(),
	}
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
	if db.options.Mysql {
		if db.mysqlConfigs != nil {
			sqlFullConnection := db.mysqlConfigs.String()
			client, err := gorm.Open(mysql.Open(sqlFullConnection), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, FullSaveAssociations: false})
			if err != nil {
				panic(err)
			}
			db.mysqlClient = client
			if db.mysqlConfigs.OpenPrometheus {
				db.mysqlClient.Use(prometheus.New(prometheus.Config{
					DBName:          "db1",                                             // 使用 `DBName` 作为指标 label
					RefreshInterval: uint32(db.mysqlConfigs.PrometheusRefreshInterval), // 指标刷新频率（默认为 15 秒）
					PushAddr:        "prometheus pusher address",                       // 如果配置了 `PushAddr`，则推送指标
					StartServer:     true,                                              // 启用一个 http 服务来暴露指标
					HTTPServerPort:  uint32(db.mysqlConfigs.PrometheusPort),            // 配置 http 服务监听端口，默认端口为 8080 （如果您配置了多个，只有第一个 `HTTPServerPort` 会被使用）
					MetricsCollector: []prometheus.MetricsCollector{
						&prometheus.MySQL{
							VariableNames: []string{"Threads_running"},
						},
					}, // 用户自定义指标
				}))
			}
		} else {
			dsn := "tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
			db.logger.Warn("mysqlConfig not set.Use the default DSN", zap.String("DSN", dsn))
			client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, FullSaveAssociations: false})
			if err != nil {
				panic(err)
			}
			db.mysqlClient = client
		}
	}

	//连接mongodb

	if db.options.Mongo {
		if db.mongoConfigs != nil {
			sqlFullConnection := db.mongoConfigs.String()
			client, err := db.NewMongoClient(sqlFullConnection, db.mongoConfigs.UserName, db.mongoConfigs.Password)
			if err != nil {
				panic(err)
			}
			db.mongoClient = client
		} else {
			sqlFullConnection := "mongodb://127.0.0.1:27017"
			client, err := db.NewMongoClient(sqlFullConnection, "", "")
			if err != nil {
				panic(err)
			}
			db.mongoClient = client
		}
	}

	return db
}

func (db *DB) SetMysqlConfig(c *MysqlConfig) DAO {
	db.mysqlConfigs = c
	return db
}

func (db *DB) SetMongoConfig(c *MongoConfig) DAO {
	db.mongoConfigs = c
	return db
}

func (db *DB) GetMysql() *gorm.DB {
	return db.mysqlClient
}

func (db *DB) AutoMigrate(models ...interface{}) error {
	return db.mysqlClient.AutoMigrate(models...)
}

func (db *DB) NewQuery() *QueryOptions {
	return &QueryOptions{
		conditions:    []interface{}{},
		selectList:    []string{},
		joinTableList: []string{},
	}
}
