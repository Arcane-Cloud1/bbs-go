package admin

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/repositories"
	"strconv"

	"bbs-go/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"github.com/spf13/cast"

	"bbs-go/internal/services"
)

func UserSynccount(ctx *gin.Context) {
	go func() {
		services.UserService.Scan(func(users []models.User) {
			for _, user := range users {
				topicCount := repositories.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).Eq("status", constants.StatusOk))
				commentCount := repositories.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).Eq("status", constants.StatusOk))
				_ = repositories.UserRepository.UpdateColumn(sqls.DB(), user.Id, "topic_count", topicCount)
				_ = repositories.UserRepository.UpdateColumn(sqls.DB(), user.Id, "comment_count", commentCount)
				cache.UserCache.Invalidate(user.Id)
			}
		})
	}()
	ctx.JSON(200, web.JsonSuccess())
}

func UserGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	t := services.UserService.Get(id)
	if t == nil {
		ctx.JSON(200, web.JsonErrorMsg("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ctx.JSON(200, web.JsonData(buildUserItem(t, true)))
}

func UserList(ctx *gin.Context) {
	list, paging := services.UserService.FindPageByCnd(params.NewPagedSqlCndGin(ctx,
		params.QueryFilter{
			ParamName: "id",
			Op:        params.Eq,
			ValueWrapper: func(origin string) string {
				if id := idcodec.Decode(origin); id > 0 {
					return cast.ToString(id)
				}
				return ""
			},
		},
		params.QueryFilter{
			ParamName: "nickname",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "email",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "username",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "type",
			Op:        params.Eq,
		},
	).Desc("id"))
	var itemList []map[string]interface{}
	for _, user := range list {
		itemList = append(itemList, buildUserItem(&user, false))
	}
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: itemList, Page: paging}))
}

func UserCreate(ctx *gin.Context) {
	username := ctx.PostForm("username")
	email := ctx.PostForm("email")
	nickname := ctx.PostForm("nickname")
	password := ctx.PostForm("password")

	user, err := services.UserService.SignUp(username, email, nickname, password, password)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(buildUserItem(user, true)))
}

func UserUpdate(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.PostForm("id"), 10, 64)
	_type, _ := strconv.Atoi(ctx.PostForm("type"))
	username := ctx.PostForm("username")
	email := ctx.PostForm("email")
	nickname := ctx.PostForm("nickname")
	avatar := ctx.PostForm("avatar")
	gender := ctx.PostForm("gender")
	homePage := ctx.PostForm("homePage")
	description := ctx.PostForm("description")
	roleIds := params.FormValueInt64Array(ctx, "roleIds")
	status, _ := strconv.Atoi(ctx.PostForm("status"))

	user := services.UserService.Get(id)
	if user == nil {
		ctx.JSON(200, web.JsonErrorMsg("entity not found"))
		return
	}

	user.Type = _type
	user.Username = sqls.SqlNullString(username)
	user.Email = sqls.SqlNullString(email)
	user.Nickname = nickname
	user.Avatar = avatar
	user.Gender = constants.Gender(gender)
	user.HomePage = homePage
	user.Description = description
	user.Status = status

	if err := services.UserService.Update(user); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.UserRoleService.UpdateUserRoles(user.Id, roleIds); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	user = services.UserService.Get(user.Id)
	ctx.JSON(200, web.JsonData(buildUserItem(user, true)))
}

func UserForbidden(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		ctx.JSON(200, web.JsonErrorMsg("无权限"))
		return
	}

	userId := params.FormValueInt64Default(ctx, "userId", 0)
	days := params.FormValueIntDefault(ctx, "days", 0)
	reason := ctx.PostForm("reason")

	if userId < 0 {
		ctx.JSON(200, web.JsonErrorMsg("请传入：userId"))
		return
	}

	if days == 0 {
		services.UserService.RemoveForbidden(user.Id, userId, ctx.Request)
	} else {
		if err := services.UserService.Forbidden(user.Id, userId, days, reason, ctx.Request); err != nil {
			ctx.JSON(200, web.JsonError(err))
			return
		}
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserUpdatePassword(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	oldPassword := ctx.PostForm("oldPassword")
	password := ctx.PostForm("password")
	rePassword := ctx.PostForm("rePassword")

	if err := services.UserService.UpdatePassword(user.Id, oldPassword, password, rePassword); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserResetPassword(ctx *gin.Context) {
	userId, _ := strconv.ParseInt(ctx.PostForm("userId"), 10, 64)
	if userId <= 0 {
		ctx.JSON(200, web.JsonErrorMsg("invalid param: userId"))
		return
	}

	newPassword, err := services.UserService.ResetPassword(userId)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	ctx.JSON(200, web.JsonData(gin.H{
		"password": newPassword,
	}))
}

func buildUserItem(user *models.User, buildRoleIds bool) map[string]interface{} {
	b := web.NewRspBuilder(user).
		Put("idEncode", idcodec.Encode(user.Id)).
		Put("roles", user.GetRoles()).
		Put("username", user.Username.String).
		Put("email", user.Email.String).
		Put("score", user.Score).
		Put("forbidden", user.IsForbidden())
	if buildRoleIds {
		b.Put("roleIds", services.UserRoleService.GetUserRoleIds(user.Id))
	}
	return b.Build()
}
