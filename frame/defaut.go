package frame

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xm-chentl/gocore/frame/handles"
)

type RouteHandle []string

func (r RouteHandle) Relative() (res string) {
	for _, v := range r {
		res += "/:" + v
	}
	return
}

func (r RouteHandle) Route(ctx *gin.Context) (res string) {
	vArr := make([]string, 0)
	for _, v := range r {
		vArr = append(vArr, ctx.Param(v))
	}
	res = "/" + strings.Join(vArr, "/")

	return
}

func NewHandle(method string, routeHandle RouteHandle) GinOption {
	relativePath := routeHandle.Relative()
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

			route := routeHandle.Route(ctx)
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
