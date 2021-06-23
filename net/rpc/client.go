package rpc

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	ztime "gogs.buffalo-robot.com/zouhy/micro/time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	_once           sync.Once
	_defaultCliConf = &ClientConfig{
		Dial:              ztime.Duration(time.Second * 10),
		Timeout:           ztime.Duration(time.Millisecond * 250),
		Subset:            50,
		KeepAliveInterval: ztime.Duration(time.Second * 60),
		KeepAliveTimeout:  ztime.Duration(time.Second * 20),
	}
	_defaultClient *RpcClient
)

type ClientConfig struct {
	Dial    ztime.Duration
	Timeout ztime.Duration
	// Breaker                *breaker.Config
	Method                 map[string]*ClientConfig
	Clusters               []string
	Zone                   string
	Subset                 int
	NonBlock               bool
	KeepAliveInterval      ztime.Duration
	KeepAliveTimeout       ztime.Duration
	KeepAliveWithoutStream bool
}

type RpcClient struct {
	// breaker *breaker.Group
	conf     *ClientConfig
	mutex    sync.RWMutex
	opts     []grpc.DialOption
	handlers []grpc.UnaryClientInterceptor
}

type TimeoutCallOption struct {
	*grpc.EmptyCallOption
	Timeout ztime.Duration
}

func WithTimeoutCallOption(timeout ztime.Duration) *TimeoutCallOption {
	return &TimeoutCallOption{&grpc.EmptyCallOption{}, timeout}
}

func NewConn(target string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DefaultClient().Dial(context.Background(), target, opt...)
}
func NewClient(conf *ClientConfig, opt ...grpc.DialOption) *RpcClient {
	c := new(RpcClient)
	if err := c.SetConfig(conf); err != nil {
		panic(err)
	}
	// c.UseOpt(grpc.WithBalancerName(p2c.Name))
	c.UseOpt(opt...)
	return c
}

func (c *RpcClient) Use(handlers ...grpc.UnaryClientInterceptor) *RpcClient {
	finalSize := len(c.handlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("warden: client use too many handlers")
	}
	mergedHandlers := make([]grpc.UnaryClientInterceptor, finalSize)
	copy(mergedHandlers, c.handlers)
	copy(mergedHandlers[len(c.handlers):], handlers)
	c.handlers = mergedHandlers
	return c
}

func DefaultClient() *RpcClient {
	_once.Do(func() {
		_defaultClient = NewClient(nil)
	})
	return _defaultClient
}

// SetConfig hot reloads client config
func (c *RpcClient) SetConfig(conf *ClientConfig) (err error) {
	if conf == nil {
		conf = _defaultCliConf
	}
	if conf.Dial <= 0 {
		conf.Dial = ztime.Duration(time.Second * 10)
	}
	if conf.Timeout <= 0 {
		conf.Timeout = ztime.Duration(time.Millisecond * 250)
	}
	if conf.Subset <= 0 {
		conf.Subset = 50
	}
	if conf.KeepAliveInterval <= 0 {
		conf.KeepAliveInterval = ztime.Duration(time.Second * 60)
	}
	if conf.KeepAliveTimeout <= 0 {
		conf.KeepAliveTimeout = ztime.Duration(time.Second * 20)
	}

	// FIXME(maojian) check Method dial/timeout
	c.mutex.Lock()
	c.conf = conf
	// if c.breaker == nil {
	// 	c.breaker = breaker.NewGroup(conf.Breaker)
	// } else {
	// 	c.breaker.Reload(conf.Breaker)
	// }
	c.mutex.Unlock()
	return nil
}

func (c *RpcClient) UseOpt(opts ...grpc.DialOption) *RpcClient {
	c.opts = append(c.opts, opts...)
	return c
}

func (c *RpcClient) cloneOpts() []grpc.DialOption {
	dialOptions := make([]grpc.DialOption, len(c.opts))
	copy(dialOptions, c.opts)
	return dialOptions
}

func (c *RpcClient) dial(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	dialOptions := c.cloneOpts()
	if !c.conf.NonBlock {
		dialOptions = append(dialOptions, grpc.WithBlock())
	}
	dialOptions = append(dialOptions, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                time.Duration(c.conf.KeepAliveInterval),
		Timeout:             time.Duration(c.conf.KeepAliveTimeout),
		PermitWithoutStream: !c.conf.KeepAliveWithoutStream,
	}))
	dialOptions = append(dialOptions, opts...)

	// init default handler
	var handlers []grpc.UnaryClientInterceptor
	handlers = append(handlers, c.recovery())
	// handlers = append(handlers, clientLogging(dialOptions...))
	handlers = append(handlers, c.handlers...)
	// NOTE: c.handle must be a last interceptor.
	// handlers = append(handlers, c.handle())

	dialOptions = append(dialOptions, grpc.WithUnaryInterceptor(chainUnaryClient(handlers)))
	c.mutex.RLock()
	conf := c.conf
	c.mutex.RUnlock()
	if conf.Dial > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(conf.Dial))
		defer cancel()
	}
	// if u, e := url.Parse(target); e == nil {
	// 	v := u.Query()
	// 	for _, c := range c.conf.Clusters {
	// 		v.Add(naming.MetaCluster, c)
	// 	}
	// 	if c.conf.Zone != "" {
	// 		v.Add(naming.MetaZone, c.conf.Zone)
	// 	}
	// 	if v.Get("subset") == "" && c.conf.Subset > 0 {
	// 		v.Add("subset", strconv.FormatInt(int64(c.conf.Subset), 10))
	// 	}
	// 	u.RawQuery = v.Encode()
	// 	// 比较_grpcTarget中的appid是否等于u.path中的appid，并替换成mock的地址
	// 	// for _, t := range _grpcTarget {
	// 	// 	strs := strings.SplitN(t, "=", 2)
	// 	// 	if len(strs) == 2 && ("/"+strs[0]) == u.Path {
	// 	// 		u.Path = "/" + strs[1]
	// 	// 		u.Scheme = "passthrough"
	// 	// 		u.RawQuery = ""
	// 	// 		break
	// 	// 	}
	// 	// }
	// 	target = u.String()
	// }
	if conn, err = grpc.DialContext(ctx, target, dialOptions...); err != nil {
		fmt.Fprintf(os.Stderr, "warden client: dial %s error %v!", target, err)
	}
	err = errors.WithStack(err)
	return
}

func chainUnaryClient(handlers []grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	n := len(handlers)
	if n == 0 {
		return func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}

	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var (
			i            int
			chainHandler grpc.UnaryInvoker
		)
		chainHandler = func(ictx context.Context, imethod string, ireq, ireply interface{}, ic *grpc.ClientConn, iopts ...grpc.CallOption) error {
			if i == n-1 {
				return invoker(ictx, imethod, ireq, ireply, ic, iopts...)
			}
			i++
			return handlers[i](ictx, imethod, ireq, ireply, ic, chainHandler, iopts...)
		}

		return handlers[0](ctx, method, req, reply, cc, chainHandler, opts...)
	}
}

func (c *RpcClient) Dial(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append(opts, grpc.WithInsecure())
	return c.dial(ctx, target, opts...)
}

func (c *RpcClient) DialTLS(ctx context.Context, target string, file string, name string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	var creds credentials.TransportCredentials
	creds, err = credentials.NewClientTLSFromFile(file, name)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	return c.dial(ctx, target, opts...)
}
