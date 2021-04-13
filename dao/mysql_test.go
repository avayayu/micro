package dao

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"gogs.bfr.com/zouhy/micro/dao/drivers/mysql"
	"gogs.bfr.com/zouhy/micro/lib"
	"gogs.bfr.com/zouhy/micro/models"
)

var dao DAO

type Role struct {
	models.Model
	RoleName   string `json:"roleName" gorm:"Column:role_name;index:role_rolename_index"`        //角色名
	Describe   string `json:"describe" gorm:"Column:describe;"`                                  //角色用途
	RoleStatus bool   `json:"roleStatus" gorm:"Column:role_status;index:role_role_status_index"` //启用停用
}

func (t Role) TableName() string {
	return "auth_role"
}

func TestMain(m *testing.M) {

	configs := mysql.MysqlConfigs{
		URL:                 "192.168.100.128",
		Port:                "33309",
		UserName:            "root",
		Password:            "bfr123123",
		DBName:              "cloudbrain_test",
		MysqlOpenPrometheus: false,
	}

	mysqlDrvicer := &mysql.MysqlDriver{
		Configs: &configs,
	}

	dao = NewDatabase(mysqlDrvicer)

	dao.AutoMigrate(&DeviceType{}, &DeviceFactory{}, &Device{})

	flag.Parse()
	exitCode := m.Run()

	// 退出
	os.Exit(exitCode)

}

func TestDB_Create(t *testing.T) {

	factory := &DeviceFactory{
		FactoryName: "布法罗",
		Comments:    "强",
	}

	sliceTest := []DeviceFactory{
		{
			FactoryName: lib.RandStringBytesMaskImprSrcUnsafe(6),
			Comments:    "强",
		},
		{
			FactoryName: lib.RandStringBytesMaskImprSrcUnsafe(6),
			Comments:    "强",
		},
		{
			FactoryName: lib.RandStringBytesMaskImprSrcUnsafe(6),
			Comments:    "强",
		},
	}

	type args struct {
		model     interface{}
		createdBy string
		value     interface{}
	}
	tests := []struct {
		name    string
		db      DAO
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "DeviceFactory",
			db:   dao,
			args: args{
				model:     &DeviceFactory{},
				createdBy: "Test",
				value:     factory,
			},
			wantErr: false,
		},
		{
			name: "devicType",
			db:   dao,
			args: args{
				model:     &DeviceType{},
				createdBy: "Test",
				value: &DeviceType{
					DeviceTypeName: "外骨骼",
					DeviceTypeCode: "Aider",
					FactoryID:      factory.ID,
				},
			},
			wantErr: false,
		},
		{
			name: "slice Test",
			db:   dao,
			args: args{
				model:     &DeviceFactory{},
				createdBy: "Test",
				value:     &sliceTest,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.NewQuery().Create(tt.args.model, tt.args.createdBy, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("DB.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_Updates(t *testing.T) {

	type args struct {
		model     interface{}
		UpdatesBy string
		value     interface{}
		filters   []interface{}
	}
	tests := []struct {
		name    string
		db      DAO
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "DeviceType",
			db:   dao,
			args: args{
				model:     &DeviceType{},
				UpdatesBy: "test",
				value:     map[string]interface{}{"device_type_code": "aider_01"},
				filters:   []interface{}{"id = ?", 1328580021266644992},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.NewQuery().Updates(tt.args.model, tt.args.UpdatesBy, tt.args.value, tt.args.filters...); (err != nil) != tt.wantErr {
				t.Errorf("DB.Updates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_First(t *testing.T) {
	type args struct {
		model   interface{}
		out     interface{}
		options []*QueryOptions
	}
	tests := []struct {
		name         string
		db           DAO
		args         args
		wantNotFound bool
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name: "DeviceType_found",
			db:   dao,
			args: args{
				model:   &DeviceType{},
				out:     &DeviceType{},
				options: []*QueryOptions{{where: "device_type_code=?", conditions: []interface{}{"aider_01"}}},
			},
			wantNotFound: true,
			wantErr:      false,
		},
		{
			name: "DeviceType_not_found",
			db:   dao,
			args: args{
				model:   &DeviceType{},
				out:     &DeviceType{},
				options: []*QueryOptions{{where: "device_type_code=?", conditions: []interface{}{"aider_02"}}},
			},
			wantNotFound: false,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotFound, err := tt.db.NewQuery().First(tt.args.model, tt.args.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.First() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotNotFound != tt.wantNotFound {
				t.Errorf("DB.First() = %v, want %v", gotNotFound, tt.wantNotFound)
			}
		})
	}
}

func TestDB_Raw(t *testing.T) {

	deviceType := []DeviceType{}

	type args struct {
		sql string
		out interface{}
	}
	tests := []struct {
		name    string
		db      DAO
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "deviceType",
			db:   dao,
			args: args{
				sql: `select * from device_type where device_type_code = 'aider_01'`,
				out: &deviceType,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.NewQuery().Raw(tt.args.sql, tt.args.out); (err != nil) != tt.wantErr && len(deviceType) == 1 {
				t.Errorf("DB.Raw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_Find(t *testing.T) {

	deviceType := []DeviceType{}
	roles := []Role{}
	type args struct {
		model   interface{}
		out     interface{}
		options []*QueryOptions
	}
	tests := []struct {
		name    string
		db      DAO
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "deviceType_Find",
			db:   dao,
			args: args{
				model:   &DeviceType{},
				out:     &deviceType,
				options: []*QueryOptions{{where: "device_type_name=?", conditions: []interface{}{"外骨骼"}}},
			},
			wantErr: false,
		},
		{
			name: "deviceType_Find_0",
			db:   dao,
			args: args{
				model:   &DeviceType{},
				out:     &deviceType,
				options: []*QueryOptions{{where: "device_type_name=?", conditions: []interface{}{"aider_02"}}},
			},
			wantErr: false,
		},
		{
			name: "Role_Find_0",
			db:   dao,
			args: args{
				model: &Role{},
				out:   &roles,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.NewQuery().Find(tt.args.model, tt.args.out); (err != nil) != tt.wantErr && (tt.name == "deviceType_Find" && len(deviceType) == 1) && (tt.name == "deviceType_Find_0" && len(deviceType) == 0) {
				t.Errorf("DB.Find() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				fmt.Println(tt.args.out)
			}
		})
	}
}

func TestQueryOptions_FindToMap(t *testing.T) {

	type args struct {
		model   interface{}
		out     interface{}
		column  string
		options []Query
	}
	outData3 := map[uint64]*Role{}
	outData1 := map[string]Role{}
	outData2 := map[string]*Role{}
	tests := []struct {
		name    string
		query   Query
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:  "map to Role id test",
			query: dao.NewQuery(),
			args: args{
				model:  &Role{},
				out:    &outData3,
				column: "ID",
			},
			wantErr: false,
		},
		{
			name:  "map to Role test",
			query: dao.NewQuery(),
			args: args{
				model:  &Role{},
				out:    &outData1,
				column: "RoleName",
			},
			wantErr: false,
		},
		{
			name:  "map to Role test",
			query: dao.NewQuery(),
			args: args{
				model:  &Role{},
				out:    &outData2,
				column: "RoleName",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.query.FindToMap(tt.args.model, tt.args.out, tt.args.column); (err != nil) != tt.wantErr {
				t.Errorf("QueryOptions.FindToMap() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				fmt.Println(tt.args.out)
			}
		})
	}
}

func TestQueryOptions_Like(t *testing.T) {

	type args struct {
		where Model
	}

	outData3 := []*DeviceFactory{}

	tests := []struct {
		name    string
		query   Query
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:  "map to Role id test",
			query: dao.NewQuery().Debug(),
			args: args{
				where: &DeviceFactory{
					FactoryName: "布法罗",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.query.Like(tt.args.where).Find(&DeviceFactory{}, &outData3); (err != nil) != tt.wantErr {
				t.Errorf("QueryOptions.FindToMap() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				fmt.Println(outData3)
			}
		})
	}
}
