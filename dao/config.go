package dao

import "fmt"

type DBType uint8

const (
	_ DBType = iota
	MYSQL
	ORACLE
	MONGO
	SQLITE
)

type DBConfigs struct {
	URL                            string `json:"URL" yaml:"URL"`
	Port                           string `json:"Port" yarml:"Port"`
	UserName                       string `json:"userName" yaml:"userName"`
	Password                       string `json:"password" yaml:"password"`
	DBName                         string `json:"dbName" yaml:"DBName"`
	DBType                         DBType `json:"dbType"`
	SQLITEPATH                     string `json:"SQLITEPATH"`
	FullConnectionString           string `json:"fullConnection"`
	MysqlOpenPrometheus            bool   `json:"openPrometheus"`
	MysqlPrometheusPort            int    `json:"prometheusPort"`
	MysqlPrometheusRefreshInterval int    `json:"prometheusRefreshInterval"`
	MongoIsReplicated              bool   `json:"isReplicated"`
	MongoReplicatedName            string `json:"replicatedName"`
	OracleServiceName              string `json:"serviceName"`
	OracleLibPath                  string `json:"libPath"`
}

func (c *DBConfigs) String() string {
	switch c.DBType {
	case MYSQL:
		c.FullConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.UserName, c.Password, c.URL, c.Port, c.DBName)
		return c.FullConnectionString
	case ORACLE:
		c.FullConnectionString = fmt.Sprintf("%s/%s@%s:%s/%s", c.UserName, c.Password, c.URL, c.Port, c.DBName)
		return c.FullConnectionString
	case MONGO:
		if c.MongoIsReplicated {
			c.FullConnectionString = fmt.Sprintf("mongodb://%s/?replicaSet=%s", c.URL, c.MongoReplicatedName)
		} else {
			c.FullConnectionString = fmt.Sprintf("mongodb://%s:%s", c.URL, c.Port)
		}
		return c.FullConnectionString
	case SQLITE:
		return c.SQLITEPATH
	}
	panic("DBType must be setting")

}
