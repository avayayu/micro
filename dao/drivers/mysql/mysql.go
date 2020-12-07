package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" //加载mysql驱动给gorm
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

type MysqlDriver struct {
	Configs *MysqlConfigs
}

type MysqlConfigs struct {
	URL                            string `json:"URL" yaml:"URL"`
	Port                           string `json:"Port" yarml:"Port"`
	UserName                       string `json:"userName" yaml:"userName"`
	Password                       string `json:"password" yaml:"password"`
	DBName                         string `json:"dbName" yaml:"DBName"`
	FullConnectionString           string `json:"fullConnection"`
	MysqlOpenPrometheus            bool   `json:"openPrometheus"`
	MysqlPrometheusPort            int    `json:"prometheusPort"`
	MysqlPrometheusRefreshInterval int    `json:"prometheusRefreshInterval"`
}

func (c *MysqlConfigs) String() string {
	c.FullConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.UserName, c.Password, c.URL, c.Port, c.DBName)
	return c.FullConnectionString
}

func (d *MysqlDriver) Connect() (*gorm.DB, *mongo.Client, error) {
	sqlFullConnection := d.Configs.String()
	client, err := gorm.Open(mysql.Open(sqlFullConnection), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, FullSaveAssociations: false})
	if err != nil {
		panic(err)
	}

	if d.Configs.MysqlOpenPrometheus {
		client.Use(prometheus.New(prometheus.Config{
			DBName:          "db1",                                 // 使用 `DBName` 作为指标 label
			RefreshInterval: uint32(d.Configs.MysqlPrometheusPort), // 指标刷新频率（默认为 15 秒）
			PushAddr:        "prometheus pusher address",           // 如果配置了 `PushAddr`，则推送指标
			StartServer:     true,                                  // 启用一个 http 服务来暴露指标
			HTTPServerPort:  uint32(d.Configs.MysqlPrometheusPort), // 配置 http 服务监听端口，默认端口为 8080 （如果您配置了多个，只有第一个 `HTTPServerPort` 会被使用）
			MetricsCollector: []prometheus.MetricsCollector{
				&prometheus.MySQL{
					VariableNames: []string{"Threads_running"},
				},
			}, // 用户自定义指标
		}))
	}

	return client, nil, nil
}

func (d *MysqlDriver) Type() uint8 {
	return 1
}
