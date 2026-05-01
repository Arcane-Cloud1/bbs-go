package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/locales"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"github.com/spf13/cast"

	"bbs-go/internal/cache"
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

func UserCurrent(ctx *gin.Context) {
	if !config.Instance.Installed {
		ctx.JSON(200, web.JsonSuccess())
		return
	}
	user := common.GetCurrentUser(ctx)
	if user != nil {
		ctx.JSON(200, web.JsonData(render.BuildUserProfile(user)))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserDetail(ctx *gin.Context) {
	userIdStr := ctx.Param("id")
	userId := idcodec.Decode(userIdStr)
	user := cache.UserCache.Get(userId)
	if user != nil && user.Status != constants.StatusDeleted {
		ctx.JSON(200, web.JsonData(render.BuildUserDetail(user)))
		return
	}
	ctx.JSON(200, web.JsonErrorMsg(locales.Get("user.not_found")))
}

func UserUpdate(ctx *gin.Context) {
	userIdStr := ctx.Param("id")
	userId := idcodec.Decode(userIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}
	if user.Id != userId {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("user.no_permission")))
		return
	}

	var (
		nickname    = strings.TrimSpace(ctx.PostForm("nickname"))
		homePage    = ctx.PostForm("homePage")
		description = ctx.PostForm("description")
		gender      = strings.TrimSpace(ctx.PostForm("gender"))
	)

	var (
		minLength = constants.NicknameMinLengthEnUS
		maxLength = constants.NicknameMaxLengthEnUS
	)
	if strings.EqualFold(string(config.Instance.Language), string(config.LanguageZhCN)) {
		minLength = constants.NicknameMinLengthZhCN
		maxLength = constants.NicknameMaxLengthZhCN
	}
	if nicknameLength := utf8.RuneCountInString(nickname); nicknameLength < minLength || nicknameLength > maxLength {
		ctx.JSON(200, web.JsonErrorMsg(locales.Getf("user.nickname_length_invalid", minLength, maxLength)))
		return
	}

	if strs.IsNotBlank(gender) {
		if gender != string(constants.GenderMale) && gender != string(constants.GenderFemale) {
			ctx.JSON(200, web.JsonErrorMsg(locales.Get("user.gender_error")))
			return
		}
	}

	err := services.UserService.Updates(user.Id, map[string]any{
		"nickname":    nickname,
		"home_page":   homePage,
		"description": description,
		"gender":      gender,
	})
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserUpdateAvatar(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}
	avatar := strings.TrimSpace(ctx.PostForm("avatar"))
	if len(avatar) == 0 {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("user.avatar_empty")))
		return
	}
	err := services.UserService.UpdateAvatar(user.Id, avatar)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserUpdateBackgroundImage(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}
	backgroundImage := ctx.PostForm("backgroundImage")
	if strs.IsBlank(backgroundImage) {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("user.upload_image_required")))
		return
	}
	if err := services.UserService.UpdateBackgroundImage(user.Id, backgroundImage); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserTopics(ctx *gin.Context) {
	userId := common.GetQueryID(ctx, "userId")
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)

	topics, cursor, hasMore := services.TopicService.GetUserTopics(userId, cursor)
	ctx.JSON(200, web.JsonCursorData(render.BuildSimpleTopicsGin(ctx, topics), strconv.FormatInt(cursor, 10), hasMore))
}

func UserArticles(ctx *gin.Context) {
	userId := common.GetQueryID(ctx, "userId")
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)

	articles, cursor, hasMore := services.ArticleService.GetUserArticles(userId, cursor)
	ctx.JSON(200, web.JsonCursorData(render.BuildSimpleArticles(articles), strconv.FormatInt(cursor, 10), hasMore))
}

func UserComments(ctx *gin.Context) {
	userId := common.GetQueryID(ctx, "userId")
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)

	comments, cursor, hasMore := services.CommentService.GetUserComments(userId, cursor)
	ctx.JSON(200, web.JsonCursorData(render.BuildComments(comments), strconv.FormatInt(cursor, 10), hasMore))
}

func UserFavorites(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)

	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	limit := 20
	var favorites []models.Favorite
	if cursor > 0 {
		favorites = services.FavoriteService.Find(sqls.NewCnd().Where("user_id = ? and id < ?",
			user.Id, cursor).Desc("id").Limit(20))
	} else {
		favorites = services.FavoriteService.Find(sqls.NewCnd().Where("user_id = ?", user.Id).Desc("id").Limit(limit))
	}

	hasMore := false
	if len(favorites) > 0 {
		cursor = favorites[len(favorites)-1].Id
		hasMore = len(favorites) >= limit
	}

	ctx.JSON(200, web.JsonCursorData(render.BuildFavorites(favorites), strconv.FormatInt(cursor, 10), hasMore))
}

func UserFans(ctx *gin.Context) {
	userId := common.GetQueryID(ctx, "userId")
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)

	fans, cursor, hasMore := services.UserFollowService.GetFans(userId, cursor)
	ctx.JSON(200, web.JsonCursorData(render.BuildUsers(fans), strconv.FormatInt(cursor, 10), hasMore))
}

func UserFollows(ctx *gin.Context) {
	userId := common.GetQueryID(ctx, "userId")
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)

	follows, cursor, hasMore := services.UserFollowService.GetFollows(userId, cursor)
	ctx.JSON(200, web.JsonCursorData(render.BuildUsers(follows), strconv.FormatInt(cursor, 10), hasMore))
}

func UserFollow(ctx *gin.Context) {
	userIdStr := ctx.Param("id")
	userId := idcodec.Decode(userIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	if err := services.UserFollowService.Follow(user.Id, userId); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserUnfollow(ctx *gin.Context) {
	userIdStr := ctx.Param("id")
	userId := idcodec.Decode(userIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	if err := services.UserFollowService.Unfollow(user.Id, userId); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserMessages(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	limit := 20
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)

	cnd := sqls.NewCnd().Eq("user_id", user.Id).Limit(limit).Desc("id")
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	list := services.MessageService.Find(cnd)

	var (
		nextCursor = cursor
		hasMore    = false
	)
	if len(list) > 0 {
		nextCursor = list[len(list)-1].Id
		hasMore = len(list) == limit
	}

	services.MessageService.MarkRead(user.Id)

	ctx.JSON(200, web.JsonCursorData(render.BuildMessages(list), cast.ToString(nextCursor), hasMore))
}

func UserReadMsg(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	msgId := params.FormValueInt64Default(ctx, "id", 0)
	if msgId > 0 {
		services.MessageService.MarkReadById(user.Id, msgId)
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserReadAllMsg(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	services.MessageService.MarkRead(user.Id)
	ctx.JSON(200, web.JsonSuccess())
}

func UserScoreLogs(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	limit := 20
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	cnd := sqls.NewCnd().Eq("user_id", user.Id).Limit(limit).Desc("id")
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	list := services.UserScoreLogService.Find(cnd)

	var (
		nextCursor = cursor
		hasMore    = false
	)
	if len(list) > 0 {
		nextCursor = list[len(list)-1].Id
		hasMore = len(list) == limit
	}

	ctx.JSON(200, web.JsonCursorData(list, cast.ToString(nextCursor), hasMore))
}

func UserExpLogs(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	limit := 20
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	cnd := sqls.NewCnd().Eq("user_id", user.Id).Limit(limit).Desc("id")
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	list := services.UserExpLogService.Find(cnd)

	var (
		nextCursor = cursor
		hasMore    = false
	)
	if len(list) > 0 {
		nextCursor = list[len(list)-1].Id
		hasMore = len(list) == limit
	}

	ctx.JSON(200, web.JsonCursorData(list, cast.ToString(nextCursor), hasMore))
}

func UserCheckinStatus(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	checkIn := services.CheckInService.GetToday(user.Id)
	ctx.JSON(200, web.JsonData(map[string]interface{}{
		"checked": checkIn != nil,
	}))
}

func UserCheckinLogs(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	limit := 20
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	cnd := sqls.NewCnd().Eq("user_id", user.Id).Limit(limit).Desc("id")
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	list := services.CheckInService.Find(cnd)

	var (
		nextCursor = cursor
		hasMore    = false
	)
	if len(list) > 0 {
		nextCursor = list[len(list)-1].Id
		hasMore = len(list) == limit
	}

	ctx.JSON(200, web.JsonCursorData(list, cast.ToString(nextCursor), hasMore))
}
