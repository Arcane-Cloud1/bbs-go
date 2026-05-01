package api

import (
	"bbs-go/internal/install"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/locales"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/web"
)

func InstallInfo(ctx *gin.Context) {
	cfg := config.Instance
	ctx.JSON(200, web.JsonData(map[string]any{
		"installed": cfg.Installed,
	}))
}

func InstallTestDb(ctx *gin.Context) {
	var req install.DbConfigReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	if config.Instance.Installed {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("install.already_installed")))
		return
	}

	if err := install.TestDbConnection(req); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	ctx.JSON(200, web.JsonSuccess())
}

func InstallInstall(ctx *gin.Context) {
	var req install.InstallReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	if config.Instance.Installed {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("install.already_installed")))
		return
	}

	if err := install.Install(req); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	ctx.JSON(200, web.JsonData(true))
}
