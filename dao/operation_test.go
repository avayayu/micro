package dao

import (
	"context"
	"testing"

	"gogs.bfr.com/zouhy/micro/models"
	"gorm.io/gorm"
)

func TestQueryOptions_CheckIDList(t *testing.T) {
	type fields struct {
		order         []string
		where         string
		conditions    []interface{}
		selectList    []string
		joinTableList []string
		preloadList   []string
		Ctx           context.Context
		session       *gorm.DB
	}
	type args struct {
		model  interface{}
		idList []models.Int64Str
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := &QueryOptions{
				order:         tt.fields.order,
				where:         tt.fields.where,
				conditions:    tt.fields.conditions,
				selectList:    tt.fields.selectList,
				joinTableList: tt.fields.joinTableList,
				preloadList:   tt.fields.preloadList,
				Ctx:           tt.fields.Ctx,
				session:       tt.fields.session,
			}
			if err := query.CheckIDList(tt.args.model, tt.args.idList); (err != nil) != tt.wantErr {
				t.Errorf("QueryOptions.CheckIDList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
