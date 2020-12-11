package mssql

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type MssqlDriver struct {
	Configs *MssqlConfigs
}

type MssqlConfigs struct {
	URL                  string `json:"URL" yaml:"URL"`
	Port                 string `json:"Port" yarml:"Port"`
	UserName             string `json:"userName" yaml:"userName"`
	Password             string `json:"password" yaml:"password"`
	DBName               string `json:"dbName" yaml:"DBName"`
	MongoIsReplicated    bool   `json:"isReplicated"`
	MongoReplicatedName  string `json:"replicatedName"`
	FullConnectionString string `json:"fullConnectionString"`
}

func (c *MssqlConfigs) String() string {
	//"sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
	c.FullConnectionString = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", c.UserName, c.Password, c.URL, c.Port, c.DBName)
	return c.FullConnectionString
}

func (d *MssqlDriver) Connect() (*gorm.DB, *mongo.Client, error) {
	sqlFullConnection := d.Configs.String()
	client, err := gorm.Open(sqlserver.Open(sqlFullConnection), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, FullSaveAssociations: false})
	if err != nil {
		panic(err)
	}

	return client, nil, nil
}
