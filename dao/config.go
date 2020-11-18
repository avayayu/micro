package dao

import "fmt"

type DBOptions struct {
	Mysql  bool
	Mongo  bool
	Oracle bool
}

type Base struct {
	URL      string `json:"URL" yaml:"URL"`
	Port     string `json:"Port" yarml:"Port"`
	UserName string `json:"userName" yaml:"userName"`
	Password string `json:"password" yaml:"password"`
	DBName   string `json:"dbName" yaml:"DBName"`
}

type MysqlConfig struct {
	Base
	FullConnectionString      string `json:"fullConnection"`
	OpenPrometheus            bool   `json:"openPrometheus"`
	PrometheusPort            int    `json:"prometheusPort"`
	PrometheusRefreshInterval int    `json:"prometheusRefreshInterval"`
}

type MongoConfig struct {
	Base
	IsReplicated         bool   `json:"isReplicated"`
	ReplicatedName       string `json:"replicatedName"`
	FullConnectionString string `json:"fullConnection"`
}

type OracleConfig struct {
	Base
	FullConnectionString string `json:"fullConnection"`
}

func (c *MysqlConfig) String() string {
	c.FullConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.UserName, c.Password, c.URL, c.Port, c.DBName)
	return c.FullConnectionString
}

func (c *OracleConfig) String() string {
	c.FullConnectionString = fmt.Sprintf("%s/%s@%s:%s/%s", c.UserName, c.Password, c.URL, c.Port, c.DBName)
	return c.FullConnectionString
}

func (c *MongoConfig) String() string {
	if c.IsReplicated {
		c.FullConnectionString = fmt.Sprintf("mongodb://%s/?replicaSet=%s", c.URL, c.ReplicatedName)
	} else {
		c.FullConnectionString = fmt.Sprintf("mongodb://%s:%s", c.URL, c.Port)
	}
	return c.FullConnectionString
}
