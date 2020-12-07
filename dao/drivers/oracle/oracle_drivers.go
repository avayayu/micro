package oracle

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type OralceDriver struct {
	Configs *OralceConfigs
}

type OralceConfigs struct {
	URL                  string `json:"URL" yaml:"URL"`
	Port                 string `json:"Port" yarml:"Port"`
	UserName             string `json:"userName" yaml:"userName"`
	Password             string `json:"password" yaml:"password"`
	DBName               string `json:"dbName" yaml:"DBName"`
	OracleServiceName    string `json:"serviceName"`
	OracleLibPath        string `json:"libPath"`
	FullConnectionString string `json:"fullConnectionString"`
}

func (c *OralceConfigs) String() string {
	c.FullConnectionString = fmt.Sprintf("%s/%s@%s:%s/%s", c.UserName, c.Password, c.URL, c.Port, c.DBName)
	return c.FullConnectionString
}

func (d *OralceDriver) Connect() (*gorm.DB, *mongo.Client, error) {
	sqlFullConnection := d.Configs.String()
	client, err := gorm.Open(Open(sqlFullConnection), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return client, nil, nil

}

func (d *OralceConfigs) Type() uint8 {
	return 1
}
