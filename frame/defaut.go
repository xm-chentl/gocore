package frame

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xm-chentl/gocore/frame/handles"
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
				}
				ctx.JSON(http.StatusOK, resp)
			}()

			route := num.Route(ctx)
			fmt.Println("route >>> ", route)
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
