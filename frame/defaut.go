package frame

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHandle(method string, handle func(ctx *gin.Context)) GinOption {
	return func(g *gin.Engine) {
		g.Handle(method, "/:server/:action", func(ctx *gin.Context) {
			var err error
			resp := NewDefaultResponse()
			defer func() {
				if errRecover := recover(); errRecover != nil {
					err = NewError(500, fmt.Sprintf("%v", errRecover))
					// 栈数据
				}
				if v, ok := err.(*CustomError); ok {
					resp.Code = v.Code
					resp.Message = v.Message
				}
				ctx.JSON(http.StatusOK, resp)
			}()

			// 逻辑区
			handle(ctx)
		})
	}
}
