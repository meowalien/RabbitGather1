package member

import (
	"context"
	"github.com/gin-gonic/gin"
)

type HTTP struct {
}

func (h *HTTP) Mount(ctx context.Context, engine *gin.Engine) error {
	return nil
}
