package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gogs.buffalo-robot.com/zouhy/micro/dao"
	"gogs.buffalo-robot.com/zouhy/micro/dao/drivers/mysql"
	mmodels "gogs.buffalo-robot.com/zouhy/micro/models"
)

//DeviceFactory 设备生产商
type DeviceFactory struct {
	mmodels.Model
	FactoryName string `gorm:"Column:factory_name;type:varchar(100);not null"` //产商名称
	Comments    string `gorm:"Column:comments"`                                //备注信息
}

func (p *DeviceFactory) TableName() string {
	return "device_factory"
}

type Int64Str uint64

func (i Int64Str) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(i), 10))
}

func (i *Int64Str) UnmarshalJSON(b []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		value, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*i = Int64Str(value)
		return nil
	}

	// Fallback to number
	return json.Unmarshal(b, (*uint64)(i))
}

func main() {

	configs := mysql.MysqlConfigs{
		URL:                 "192.168.100.132",
		Port:                "33309",
		UserName:            "root",
		Password:            "bfr123123",
		DBName:              "cloudbrain_test",
		MysqlOpenPrometheus: false,
	}

	mysqlDrvicer := &mysql.MysqlDriver{
		Configs: &configs,
	}

	db := dao.NewDatabase(mysqlDrvicer)
	outDatas := []DeviceFactory{}
	if err := db.NewQuery().SelectModel(&DeviceFactory{}, "FactoryName", "CreatedAt").Find(&DeviceFactory{}, &outDatas); err != nil {
		panic(err)
	}

	fmt.Println(outDatas)
}
