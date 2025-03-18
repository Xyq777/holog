package holog

import (
	"context"

	"github.com/gin-gonic/gin"
)

func FromGinContext(ctx context.Context) *logger {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if l, exists := ginCtx.Get("logger"); exists {
			return l.(*logger)
		}
	}
	return getGlobal()
}
