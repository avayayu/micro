package rpc

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"

	nmd "gogs.buffalo-robot.com/zouhy/micro/net/metadata"
	"gogs.buffalo-robot.com/zouhy/micro/stat/sys/cpu"
	"google.golang.org/grpc"
	gmd "google.golang.org/grpc/metadata"
)

// recovery is a server interceptor that recovers from any panics.
func (s *RpcServer) recovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:rs]
				pl := fmt.Sprintf("grpc server panic: %v\n%v\n%s\n", req, rerr, buf)
				s.logger.Error(pl)
				// err = status.Errorf(codes.Unknown, ecode.ServerErr.Error())
			}
		}()
		resp, err = handler(ctx, req)
		return
	}
}

func (s *RpcServer) stats() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		var cpustat cpu.Stat
		cpu.ReadStat(&cpustat)
		if cpustat.Usage != 0 {
			trailer := gmd.Pairs([]string{nmd.CPUUsage, strconv.FormatInt(int64(cpustat.Usage), 10)}...)
			grpc.SetTrailer(ctx, trailer)
		}
		return
	}
}

func (c *RpcClient) recovery() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:rs]
				pl := fmt.Sprintf("grpc client panic: %v\n%v\n%v\n%s\n", req, reply, rerr, buf)
				fmt.Fprintf(os.Stderr, pl)
			}
		}()
		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}
}
