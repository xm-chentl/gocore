package frame

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xm-chentl/gocore/frame/handles"
	"github.com/xm-chentl/gocore/iocex"
)

type NumberSegments int

func (r NumberSegments) Relative() (res string) {
	for i := 0; i < int(r); i++ {
		res += fmt.Sprintf("/:%d", i)
	}

	return
}

func (r NumberSegments) Route(ctx *gin.Context) (res string) {
	vArr := make([]string, 0)
	for i := 0; i < int(r); i++ {
		vArr = append(vArr, ctx.Param(strconv.Itoa(i)))
	}
	res = "/" + strings.Join(vArr, "/")

	return
}

func NewHandle(method string, num NumberSegments) GinOption {
	relativePath := num.Relative()
	return func(g *gin.Engine) {
		g.Handle(method, relativePath, func(ctx *gin.Context) {
			var err error
			resp := NewDefaultResponse()
			defer func() {
				if errRecover := recover(); errRecover != nil {
					// 栈数据
					err = NewError(500, fmt.Sprintf("%v", errRecover))
				}
				if v, ok := err.(*CustomError); ok {
					resp.Code = v.Code
					resp.Message = v.Message
				} else {
					resp.Code = 500
					resp.Message = err.Error()
				}
				ctx.JSON(http.StatusOK, resp)
			}()

			route := num.Route(ctx)
			if !handles.Has(route) {
				err = ErrHandleNotExist
				return
			}

			handleTemp := reflect.New(reflect.TypeOf(handles.Get(route)).Elem()).Interface()
			if method == http.MethodGet {
				if err = ctx.BindQuery(handleTemp); err != nil {
					err = ErrQueryParameterFailed
					return
				}
			}
			if err = ctx.Bind(handleTemp); err != nil {
				return
			}
			if err = iocex.Inject(handleTemp); err != nil {
				err = ErrInjectFailed
				return
			}

			handler, ok := handleTemp.(handles.Handler)
			if !ok {
				err = ErrInvalidHandle
				return
			}

			resp.Data, err = handler.Call(ctx)
		})
	}
}

type UpgradeFunction func(ctx *gin.Context, conn *websocket.Conn) error

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocket(relativePath string, handle UpgradeFunction) GinOption {
	return func(g *gin.Engine) {
		g.GET(relativePath, func(ctx *gin.Context) {
			conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
			if err != nil {
				return
			}
			handle(ctx, conn)
		})
	}
}

func NewRegisterPProf() GinOption {
	return func(g *gin.Engine) {
		RouteRegister(&g.RouterGroup)
	}
}

var DefaultPrefix = "/debug/pprof"

func RouteRegister(rg *gin.RouterGroup, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)

	prefixRouter := rg.Group(prefix)
	{
		prefixRouter.GET("/", pprofHandler(pprof.Index))
		prefixRouter.GET("/cmdline", pprofHandler(pprof.Cmdline))
		prefixRouter.GET("/profile", pprofHandler(pprof.Profile))
		prefixRouter.POST("/symbol", pprofHandler(pprof.Symbol))
		prefixRouter.GET("/symbol", pprofHandler(pprof.Symbol))
		prefixRouter.GET("/trace", pprofHandler(pprof.Trace))
		prefixRouter.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
		prefixRouter.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		prefixRouter.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		prefixRouter.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		prefixRouter.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
		prefixRouter.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
	}
}

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// Cors 跨域设置
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}
