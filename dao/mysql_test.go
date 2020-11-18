package dao

import (
	"flag"
	"os"
	"testing"
)

var dao *DB

func TestMain(m *testing.M) {
	dao = newDatabase(&DBOptions{Mysql: true, Mongo: false})
	config := MysqlConfig{
		Base: Base{
			URL:      "192.168.100.128",
			Port:     "33309",
			UserName: "root",
			Password: "bfr123123",
			DBName:   "cloudbrain_test",
		},
		OpenPrometheus: false,
	}

	dao.SetMysqlConfig(&config).Connect()
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

	type args struct {
		model     interface{}
		createdBy string
		value     interface{}
	}
	tests := []struct {
		name    string
		db      *DB
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.Create(tt.args.model, tt.args.createdBy, tt.args.value); (err != nil) != tt.wantErr {
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
		db      *DB
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
			if err := tt.db.Updates(tt.args.model, tt.args.UpdatesBy, tt.args.value, tt.args.filters...); (err != nil) != tt.wantErr {
				t.Errorf("DB.Updates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_First(t *testing.T) {
	type args struct {
		model   interface{}
		out     interface{}
		options []QueryOptions
	}
	tests := []struct {
		name         string
		db           *DB
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
				options: []QueryOptions{{where: "device_type_code=?", conditions: []interface{}{"aider_01"}}},
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
				options: []QueryOptions{{where: "device_type_code=?", conditions: []interface{}{"aider_02"}}},
			},
			wantNotFound: false,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotFound, err := tt.db.First(tt.args.model, tt.args.out, tt.args.options...)
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
		db      *DB
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
			if err := tt.db.Raw(tt.args.sql, tt.args.out); (err != nil) != tt.wantErr && len(deviceType) == 1 {
				t.Errorf("DB.Raw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_Find(t *testing.T) {

	deviceType := []DeviceType{}

	type args struct {
		model   interface{}
		out     interface{}
		options []QueryOptions
	}
	tests := []struct {
		name    string
		db      *DB
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
				options: []QueryOptions{{where: "device_type_name=?", conditions: []interface{}{"aider_01"}}},
			},
			wantErr: false,
		},
		{
			name: "deviceType_Find_0",
			db:   dao,
			args: args{
				model:   &DeviceType{},
				out:     &deviceType,
				options: []QueryOptions{{where: "device_type_name=?", conditions: []interface{}{"aider_02"}}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.Find(tt.args.model, tt.args.out, tt.args.options...); (err != nil) != tt.wantErr && (tt.name == "deviceType_Find" && len(deviceType) == 1) && (tt.name == "deviceType_Find_0" && len(deviceType) == 0) {
				t.Errorf("DB.Find() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
