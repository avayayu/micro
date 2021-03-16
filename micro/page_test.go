package micro

import (
	"testing"

	"github.com/avayayu/micro/dao"
	"github.com/avayayu/micro/net/http"
)

func TestPagesQuery(t *testing.T) {
	type args struct {
		parameter interface{}
		out       interface{}
		db        dao.DAO
		request   http.HttpRequest
		response  http.Response
		query     *dao.QueryOptions
	}
	tests := []struct {
		name           string
		args           args
		wantTotalCount int64
		wantPage       int
		wantPerPage    int
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTotalCount, gotPage, gotPerPage, err := PagesQuery(tt.args.parameter, tt.args.out, tt.args.db, tt.args.request, tt.args.response, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("PagesQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTotalCount != tt.wantTotalCount {
				t.Errorf("PagesQuery() gotTotalCount = %v, want %v", gotTotalCount, tt.wantTotalCount)
			}
			if gotPage != tt.wantPage {
				t.Errorf("PagesQuery() gotPage = %v, want %v", gotPage, tt.wantPage)
			}
			if gotPerPage != tt.wantPerPage {
				t.Errorf("PagesQuery() gotPerPage = %v, want %v", gotPerPage, tt.wantPerPage)
			}
		})
	}
}
