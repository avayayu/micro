package drivers

import (
	"github.com/avayayu/micro/dao"
	"github.com/avayayu/micro/dao/drivers/oracle"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type OraceDriver struct{}

func (d *OraceDriver) Connect(config *dao.DBConfigs) (*gorm.DB, *mongo.Client, error) {

	if config.DBType != dao.ORACLE {
		panic("DBTYPE must be oracle")
	}

	sqlFullConnection := config.String()
	client, err := gorm.Open(oracle.Open(sqlFullConnection), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return client, nil, nil

}
