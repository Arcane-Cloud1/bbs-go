package middleware

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/urls"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/web"
)

type pathRole struct {
	Pattern string
	Roles   []string
}

var (
	authCfg = []pathRole{
		{Pattern: "/api/admin/sys-config/**", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/user/create", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/user/update", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/topic-node/create", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/topic-node/update", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/tag/create", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/tag/update", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/**", Roles: []string{constants.RoleOwner, constants.RoleAdmin}},
	}
	antPathMatcher = urls.NewAntPathMatcher()
)

func AdminMiddleware(ctx *gin.Context) {
	roles := getPathRoles(ctx)

	if len(roles) == 0 {
		ctx.Next()
		return
	}

	user := common.GetCurrentUser(ctx)
	if user == nil {
		notLogin(ctx)
		return
	}
	if !user.HasAnyRole(roles...) {
		noPermission(ctx)
		return
	}

	ctx.Next()
}

func getPathRoles(ctx *gin.Context) []string {
	p := ctx.Request.URL.Path
	for _, pathRole := range authCfg {
		if antPathMatcher.Match(pathRole.Pattern, p) {
			return pathRole.Roles
		}
	}
	return nil
}

func notLogin(ctx *gin.Context) {
	ctx.JSON(200, web.JsonError(errs.NotLogin()))
	ctx.Abort()
}

func noPermission(ctx *gin.Context) {
	ctx.JSON(200, web.JsonError(errs.NoPermission()))
	ctx.Abort()
}
