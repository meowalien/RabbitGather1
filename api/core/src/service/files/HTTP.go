package files

import (
	"context"
	"github.com/gin-gonic/gin"
)

type HTTP struct {

}

func (m *HTTP) Mount(ctx context.Context, engine *gin.Engine) error {
	router := engine.Group("/files")
	// 上傳檔案
	router.POST("/upload", m.upload)

	return nil
}

func (m *HTTP) upload(c *gin.Context) {

}