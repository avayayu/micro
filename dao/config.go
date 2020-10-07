package dao

import "fmt"

type DBOptions struct {
	Mysql bool
	Mongo bool
}

type base struct {
	URL      string `json:"URL" yaml:"URL"`
	Port     string `json:"Port" yarml:"Port"`
	UserName string `json:"userName" yaml:"userName"`
	Password string `json:"password" yaml:"password"`
	DBName   string `json:"dbName" yaml:"DBName"`
}

type mysqlConfig struct {
	base
	FullConnectionString string `json:"fullConnection"`
}

type mongoConfig struct {
	base
	IsReplicated         bool   `json:"isReplicated"`
	ReplicatedName       string `json:"replicatedName"`
	FullConnectionString string `json:"fullConnection"`
}

func (c *mysqlConfig) String() string {
	c.FullConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.UserName, c.Password, c.URL, c.Port, c.DBName)
	return c.FullConnectionString
}

func (c *mongoConfig) String() string {
	if c.IsReplicated {
		c.FullConnectionString = fmt.Sprintf("mongodb://%s:%s/?replicaSet=%s", c.URL, c.Port, c.ReplicatedName)
	} else {
		c.FullConnectionString = fmt.Sprintf("mongodb://%s:%s", c.URL, c.Port)
	}
	return c.FullConnectionString
}
