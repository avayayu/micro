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
	_ "github.com/go-sql-driver/mysql" //加载mysql驱动给gorm
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/avayayu/micro/logging"
)

type DAO interface {
	GetPageWithFilters(model interface{}, filters *Filter, out interface{}, pageIndex, pageSize int, totalCount *int64, autoLoad bool, whereOrder ...PageWhereOrder) error
	GetPageByRaw(sql string, out interface{}, pageIndex, pageSize int, totalCount *int64, where ...interface{}) error
}

type DBConnection interface {
}

const _module = "dbManager"

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
func NewDatabase(options *DBOptions) *DB {
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

func (db *DB) Connect() *DB {
	//连接mysql
	if db.options.Mysql {
		if db.mysqlConfigs != nil {
			sqlFullConnection := db.mysqlConfigs.String()
			client, err := gorm.Open(mysql.Open(sqlFullConnection), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
			if err != nil {
				panic(err)
			}
			db.mysqlClient = client
		} else {
			dsn := "tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
			db.logger.Warn("mysqlConfig not set.Use the default DSN", zap.String("DSN", dsn))
			client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
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

func (db *DB) SetMysqlConfig(c *MysqlConfig) *DB {
	db.mysqlConfigs = c
	return db
}

func (db *DB) SetMongoConfig(c *MongoConfig) *DB {
	db.mongoConfigs = c
	return db
}

func (db *DB) GetMysql() *gorm.DB {
	return db.mysqlClient
}
