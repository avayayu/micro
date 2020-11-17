package dao

import (
	mmodels "github.com/avayayu/micro/models"
	ztime "github.com/avayayu/micro/time"
)

//DeviceType 云端设备类型信息
type DeviceType struct {
	mmodels.Model
	DeviceTypeName          string        `gorm:"Column:device_type_name;type:varchar(50);not null;UNIQUE_INDEX"` //设备类型名 官方名称
	DeviceTypeCode          string        `gorm:"Column:device_type_code"`                                        //设备类型助记码
	DeviceGovRegisterCode   string        `gorm:"Column:device_gov_register_code"`                                //设备类医疗器械注册编号
	DeviceGovProductionCode string        `gorm:"Column:device_gov_production_code"`                              //设备类Production 生产许可编号
	FactoryID               uint64        `gorm:"Column:factory_id;type:bigint"`                                  //产商
	IsReport                bool          `gorm:"Column:is_report"`                                               //是否生成报告
	IsStructured            bool          `gorm:"Column:is_structured"`                                           //是否生成结构化数据
	Age                     uint          `gorm:"Column:age"`                                                     //使用寿命
	Weight                  float64       `gorm:"Column:weight"`                                                  //重量
	MaximumLoad             float64       `gorm:"Column:maximum_load"`                                            //最大荷重
	IsWifi                  bool          `gorm:"Column:is_wifi"`                                                 //是否有wifi
	IsBluetooth             bool          `gorm:"Column:is_blue_tooth"`                                           //是否有蓝牙
	IsDesktop               bool          `gorm:"Column:is_desktop"`                                              //是否有桌面程序
	IsApp                   bool          `gorm:"Column:is_app"`                                                  //是否有APP
	Factory                 DeviceFactory `json:"factory" gorm:"references:factory_id;foreignkey:id"`             //厂商
	Comments                string        `gorm:"Column:comment"`                                                 //备注信息
}

//DeviceFactory 设备生产商
type DeviceFactory struct {
	mmodels.Model
	FactoryName string `gorm:"Column:factory_name;type:varchar(100);not null"` //产商名称
	Comments    string `gorm:"Column:comments"`                                //备注信息
}

func (p DeviceFactory) TableName() string {
	return "device_factory"
}

//Device .云端获取账户绑定的设备信息（设备名称、蓝牙地址）（可以放在登录接口里面）
type Device struct {
	mmodels.Model
	DeviceID               string              `json:"deviceID" gorm:"Column:device_id;type:varchar(256);not null;index:device_id_index"` //设备ID 外部ID 非系统生成
	DeviceToken            string              `json:"deviceToken" gorm:"Column:device_token"`                                            //设备Token
	DeviceTypeID           string              `json:"deviceTypeID" gorm:"Column:device_type_id;type:bigint"`                             //设备类型ID
	DeviceUsage            string              `json:"deviceUsage" gorm:"Column:device_usage"`                                            //设备用途
	HardwareProgramVersion string              `gorm:"Column:hardware_program_version"`                                                   //硬件驱动版本
	WIFIMAC                string              `gorm:"Column:wifi_mac"`                                                                   //WIFI MAC
	OutDate                ztime.Time          `json:"outDate" gorm:"Column:out_date"`                                                    //出库日期
	Comments               string              `json:"comments" gorm:"Column:comment"`                                                    //备注信息
	DeviceType             DeviceType          `gorm:"references:device_type_id;foreignkey:id"`
	Longitude              float64             `gorm:"Column:longitude;default:104.07224258178712" json:"longitude"` //经度
	Latitude               float64             `gorm:"Column:latitude;default:30.4919057496395" json:"latitude"`     //纬度
	Province               string              `gorm:"Column:province;index:province_index" json:"province"`
	City                   string              `gorm:"Column:city" json:"city"`
	Distinct               string              `gorm:"Column:distinct" json:"district"`
	Street                 string              `gorm:"Column:street" json:"street"`
	Bluetooth              []BluetoothDescribe `gorm:"references:id;foreignkey:device_id"` //多个蓝牙地址
}

//BluetoothDescribe 设备蓝牙模块描述
type BluetoothDescribe struct {
	mmodels.Model
	Name          string `gorm:"Column:name" json:"name"`
	Address       string `gorm:"Column:address" json:"address"`
	DeviceModelID string `gorm:"Column:device_id;type:bigint" json:"deviceID"`
	Usages        string `gorm:"Column:usage" json:"usages"`
}

//TableName 表名
func (p DeviceType) TableName() string {
	return "device_type"
}

//TableName 表名
func (p Device) TableName() string {
	return "device_device"
}

//TableName 表名
func (p BluetoothDescribe) TableName() string {
	return "device_blutooth_describe"
}
