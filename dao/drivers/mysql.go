package drivers

import (
	"github.com/avayayu/micro/dao"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

type MysqlDriver struct{}

func (d *MysqlDriver) Connect(config *dao.DBConfigs) (*gorm.DB, *mongo.Client, error) {
	sqlFullConnection := config.String()
	client, err := gorm.Open(mysql.Open(sqlFullConnection), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, FullSaveAssociations: false})
	if err != nil {
		panic(err)
	}

	if config.MysqlOpenPrometheus {
		client.Use(prometheus.New(prometheus.Config{
			DBName:          "db1",                                         // 使用 `DBName` 作为指标 label
			RefreshInterval: uint32(config.MysqlPrometheusRefreshInterval), // 指标刷新频率（默认为 15 秒）
			PushAddr:        "prometheus pusher address",                   // 如果配置了 `PushAddr`，则推送指标
			StartServer:     true,                                          // 启用一个 http 服务来暴露指标
			HTTPServerPort:  uint32(config.MysqlPrometheusPort),            // 配置 http 服务监听端口，默认端口为 8080 （如果您配置了多个，只有第一个 `HTTPServerPort` 会被使用）
			MetricsCollector: []prometheus.MetricsCollector{
				&prometheus.MySQL{
					VariableNames: []string{"Threads_running"},
				},
			}, // 用户自定义指标
		}))
	}

	return client, nil, nil
}
