package code

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gogs.buffalo-robot.com/zouhy/micro/code/api"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func WrapError(err error) error {
	if err == nil {
		return nil
	}

	s := &spb.Status{
		Code:    int32(codes.Unknown),
		Message: err.Error(),
		Details: []*any.Any{
			{
				Value: []byte(stack()),
			},
		},
	}
	return status.FromProto(s).Err()
}

// Stack 获取堆栈信息
func stack() string {
	var pc = make([]uintptr, 20)
	n := runtime.Callers(3, pc)

	var build strings.Builder
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i] - 1)
		file, line := f.FileLine(pc[i] - 1)
		n := strings.Index(file, f.Name())
		if n != -1 {
			s := fmt.Sprintf(" %s:%d \n", file[n:], line)
			build.WriteString(s)
		}
	}
	return build.String()
}

func WrapBFRError(requestID string, code int32, err error) *api.Status {
	return &api.Status{
		Code:      code,
		Message:   err.Error(),
		RequestID: requestID,
		Details: []*any.Any{
			{
				Value: []byte(stack()),
			},
		},
	}

}

func WrapDetails(details []*any.Any) []zapcore.Field {
	fields := []zapcore.Field{}
	for index, detail := range details {
		data, _ := proto.Marshal(detail)
		field := zap.String(fmt.Sprintf("%d", index), string(data))
		fields = append(fields, field)
	}
	return fields
}
