package middleware

import (
	"bbs-go/internal/pkg/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/web"
)

func InstallMiddleware(ctx *gin.Context) {
	if config.Instance.Installed {
		ctx.Next()
		return
	}

	path := ctx.Request.URL.Path
	if strings.HasPrefix(path, "/api/install/") || path == "/api/config/configs" || path == "/api/user/current" {
		ctx.Next()
		return
	}

	ctx.JSON(200, web.JsonErrorCode(-1, "Please install first"))
	ctx.Abort()
}
