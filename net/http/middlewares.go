package http

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gogs.buffalo-robot.com/zouhy/micro/code"
	"gogs.buffalo-robot.com/zouhy/micro/models"
)

//TimeMiddleWare 接口访问延迟中间件
func TimeMiddleWare(logger *zap.Logger) HandlerFunc {
	return func(c *Context) {
		t := time.Now()

		c.Next()
		status := c.WritedStatus
		latency := time.Since(t)
		if c.RoutePath != "" {
			logger.Debug("接口延迟",
				zap.String("latency", strconv.Itoa(int(latency.Milliseconds()))+"ms"),
				zap.String("path", c.RoutePath),
				zap.String("METHOD", c.method),
				zap.String("visitor", c.Request.RemoteAddr),
				zap.Int("状态码", status))
		}
	}
}

//RequestIDMiddleware 给每个请求添加唯一id 用于链路跟踪
func RequestIDMiddleware(logger *zap.Logger) HandlerFunc {
	return func(c *Context) {
		requestID := uuid.NewV4().String()

		c.Set("requestID", requestID)

		if c.RoutePath != "" {
			logger.Debug("接受到请求",
				zap.String("METHOD", c.method),
				zap.String("path", c.RoutePath),
				zap.String("visitor", c.Request.RemoteAddr),
				zap.String("request", requestID))
		}
		c.Next()

	}
}

//SetGRPCMiddlewaregin 给每个请求添加唯一id 用于链路跟踪
func SetGRPCMiddlewaregin() HandlerFunc {

	return func(c *Context) {
		// path := c.FullPath()
		// flysnowRegexp := regexp.MustCompile(`^(\w+)v1/([\w])/(\w+)`)
		// service := regexp.
		// 	c.Set("grpc", server)
		c.Next()
	}
}

//GinMiddleware socketio middleware
func GinMiddleware(allowOrigins string) HandlerFunc {
	return func(c *Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

//=========================================================JWT Middwares======================================

//JWT AA
type JWT struct {
	SigningKey []byte
}

//JWTclaims 载荷
type JWTclaims struct {
	ID       models.Int64Str `json:"id"`
	UserName string          `json:"userName"`
	Email    string          `json:"email"`
	Phone    string          `json:"phone"`
	WxID     string          `json:"wxID"`
	// Role   string `json:"role"`
	jwt.StandardClaims
}

//JWT签名结构

//CreateToken 创建Token
func (j *JWT) CreateToken(claims JWTclaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*JWTclaims, error) {
	var claims *JWTclaims
	var err error

	at(time.Unix(0, 0), func() {
		token, err := jwt.ParseWithClaims(tokenString, &JWTclaims{}, func(token *jwt.Token) (interface{}, error) {
			return j.SigningKey, nil
		})
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				claims = nil
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					err = code.JWTErrorMalformed
					return
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					// Token is expired
					err = code.JWTErrorExpired
					return
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					err = code.JWTErrorNotValidYet
					return
				} else {
					err = code.JWTErrorInvalid
					return
				}
			}
		}
		if claimsInner, ok := token.Claims.(*JWTclaims); ok && token.Valid {
			claims = claimsInner
			err = nil
			return
		}
		claims = nil
		err = code.JWTErrorInvalid
		return
	})

	return claims, err

}

//RefreshToken 更新Token 从当前时间延续一个月的有效期
func (j *JWT) RefreshToken(tokenString string, tokenLast time.Duration) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &JWTclaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*JWTclaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(tokenLast).Unix()
		return j.CreateToken(*claims)
	}
	return "", code.JWTErrorInvalid
}

// JWTAuth 中间件，检查token
func JWTAuth() HandlerFunc {
	return func(c *Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.FlushHttpClientError(code.JWTErrorNotFound, "JWToken不存在", nil)
			c.Abort()
			return
		}
		j := new(JWT)
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil || claims == nil {
			if err == code.JWTErrorExpired {
				c.FlushHttpClientError(code.JWTErrorExpired, "JWToken已过期", err)
				c.Abort()
				// c.AbortWithStatusJSON(errors.HTTPUnauthorizedError, httpErr)
				return
			}
			c.FlushHttpClientError(code.JWTErrorInvalid, "JWToken解析错误", err)
			c.Abort()
			return
		} else {
			c.Set("claims", claims)
			c.Set("reqUserID", claims.ID)
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Next()
	}
}

//GenerateToken 为用户生成Token claim可以为空 只有id是必须的
func GenerateToken(id models.Int64Str, claim *JWTclaims, lastDuration time.Duration) (string, error) {
	j := &JWT{
		[]byte("bfr-cloud"),
	}
	claims := JWTclaims{
		id,
		claim.UserName,
		claim.Email,
		claim.Phone,
		claim.WxID,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),            //签名生效时间
			ExpiresAt: int64(time.Now().Add(lastDuration).Unix()), //签名过期时间 一个月
			Issuer:    "bfr-cloud",                                //签名发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}

//CrosHandler 简单开启所有的跨域功能
func CrosHandler() HandlerFunc {
	return func(context *Context) {
		method := context.Request.Method
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
		context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Writer.Header().Set("Access-Control-Max-Age", "172800")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		context.Set("content-type", "application/json")

		if method == "OPTIONS" {
			context.FlushHttpResponse()
		}

		//处理请求
		context.Next()
	}
}
