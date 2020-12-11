package sqlite

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDriver struct {
	Configs *SQLiteConfigs
}

type SQLiteConfigs struct {
	FilePath string
}

func (c *SQLiteConfigs) String() string {
	//"sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
	return c.FilePath
}

func (d *SQLiteDriver) Connect() (*gorm.DB, *mongo.Client, error) {
	sqlFullConnection := d.Configs.String()
	client, err := gorm.Open(sqlite.Open(sqlFullConnection), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, FullSaveAssociations: false})
	if err != nil {
		panic(err)
	}

	return client, nil, nil
}
