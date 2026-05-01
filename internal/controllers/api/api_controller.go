package api

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/services"
	"strconv"
	"strings"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	captcha2 "bbs-go/internal/pkg/captcha"
)

func ConfigConfigs(ctx *gin.Context) {
	configs := services.SysConfigService.GetConfigs()
	ctx.JSON(200, web.JsonData(configs))
}

func ConfigSys(ctx *gin.Context) {
	sysConfig := services.SysConfigService.GetSysConfig()
	ctx.JSON(200, web.JsonData(sysConfig))
}

func CaptchaGet(ctx *gin.Context) {
	captchaId := captcha.New()
	ctx.JSON(200, web.JsonData(gin.H{
		"captchaId": captchaId,
	}))
}

func CaptchaVerify(ctx *gin.Context) {
	var req struct {
		CaptchaId   string `json:"captchaId"`
		CaptchaCode string `json:"captchaCode"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	if !captcha.VerifyString(req.CaptchaId, req.CaptchaCode) {
		ctx.JSON(200, web.JsonError(errs.CaptchaError()))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UploadUpload(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	attachment, err := services.UploadService.Upload(file, header, user.Id)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	ctx.JSON(200, web.JsonData(gin.H{
		"id":       attachment.Id,
		"url":      attachment.Url,
		"fileName": attachment.FileName,
	}))
}

func ArticleTags(ctx *gin.Context) {
	tags := services.ArticleTagService.Find(sqls.NewCnd().Desc("id"))
	ctx.JSON(200, web.JsonData(render.BuildArticleTags(tags)))
}

func ArticleRecent(ctx *gin.Context) {
	articles := services.ArticleService.Find(sqls.NewCnd().Where("status = ?", constants.StatusOk).Desc("id").Limit(10))
	ctx.JSON(200, web.JsonData(render.BuildSimpleArticles(articles)))
}

func ArticleDetail(ctx *gin.Context) {
	articleIdStr := ctx.Param("id")
	articleId := idcodec.Decode(articleIdStr)
	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == constants.StatusDeleted {
		ctx.JSON(200, web.JsonErrorMsg("文章不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(render.BuildArticle(article)))
}

func TagTags(ctx *gin.Context) {
	tags := services.TagService.Find(sqls.NewCnd().Desc("id"))
	ctx.JSON(200, web.JsonData(render.BuildTags(tags)))
}

func TagCreate(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	name := strings.TrimSpace(ctx.PostForm("name"))
	description := ctx.PostForm("description")

	tag, err := services.TagService.Create(user.Id, name, description)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(render.BuildTag(tag)))
}

func TagDetail(ctx *gin.Context) {
	tagIdStr := ctx.Param("id")
	tagId := idcodec.Decode(tagIdStr)
	tag := services.TagService.Get(tagId)
	if tag == nil {
		ctx.JSON(200, web.JsonErrorMsg("标签不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(render.BuildTag(tag)))
}

func TagUpdate(ctx *gin.Context) {
	tagIdStr := ctx.Param("id")
	tagId := idcodec.Decode(tagIdStr)

	var form struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	err := services.TagService.Update(tagId, form.Name, form.Description)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func TagDelete(ctx *gin.Context) {
	tagIdStr := ctx.Param("id")
	tagId := idcodec.Decode(tagIdStr)

	err := services.TagService.Delete(tagId)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func CommentComments(ctx *gin.Context) {
	entityType := ctx.Query("entityType")
	entityId := params.FormValueInt64Default(ctx, "entityId", 0)
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)

	comments, cursor, hasMore := services.CommentService.GetComments(entityType, entityId, cursor)
	ctx.JSON(200, web.JsonCursorData(render.BuildComments(comments), strconv.FormatInt(cursor, 10), hasMore))
}

func CommentCreate(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	var form struct {
		EntityType      string `json:"entityType"`
		EntityId        int64  `json:"entityId"`
		Content         string `json:"content"`
		ContentType     string `json:"contentType"`
		QuoteId         int64  `json:"quoteId"`
		ImageList       string `json:"imageList"`
	}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	comment, err := services.CommentService.Publish(user.Id, form.EntityType, form.EntityId, form.Content, form.ContentType, form.QuoteId, form.ImageList)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(render.BuildComment(comment)))
}

func CommentDetail(ctx *gin.Context) {
	commentIdStr := ctx.Param("id")
	commentId := idcodec.Decode(commentIdStr)
	comment := services.CommentService.Get(commentId)
	if comment == nil {
		ctx.JSON(200, web.JsonErrorMsg("评论不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(render.BuildComment(comment)))
}

func CommentDelete(ctx *gin.Context) {
	commentIdStr := ctx.Param("id")
	commentId := idcodec.Decode(commentIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	err := services.CommentService.Delete(commentId, user.Id)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func FavoriteFavorites(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	favorites, cursor, hasMore := services.FavoriteService.GetUserFavorites(user.Id, cursor)
	ctx.JSON(200, web.JsonCursorData(render.BuildFavorites(favorites), strconv.FormatInt(cursor, 10), hasMore))
}

func FavoriteCreate(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	entityType := ctx.PostForm("entityType")
	entityId := params.FormValueInt64Default(ctx, "entityId", 0)

	err := services.FavoriteService.Add(user.Id, entityType, entityId)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func FavoriteDelete(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	favoriteIdStr := ctx.Param("id")
	favoriteId := idcodec.Decode(favoriteIdStr)

	err := services.FavoriteService.Delete(favoriteId, user.Id)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func LikeLike(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	entityType := ctx.Param("entityType")
	entityId := params.FormValueInt64Default(ctx, "entityId", 0)

	liked, err := services.UserLikeService.Like(user.Id, entityType, entityId)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(gin.H{"liked": liked}))
}

func LikeUnlike(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	entityType := ctx.Param("entityType")
	entityId := params.FormValueInt64Default(ctx, "entityId", 0)

	err := services.UserLikeService.Unlike(user.Id, entityType, entityId)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func LikeUsers(ctx *gin.Context) {
	entityType := ctx.Param("entityType")
	entityId := params.FormValueInt64Default(ctx, "entityId", 0)

	users := services.UserLikeService.Recent(entityType, entityId, 10)
	ctx.JSON(200, web.JsonData(render.BuildUsers(users)))
}

func CheckinCheckin(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	result, err := services.CheckInService.CheckIn(user.Id)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(result))
}

func CheckinRanking(ctx *gin.Context) {
	ranking := services.CheckInService.GetRanking()
	ctx.JSON(200, web.JsonData(ranking))
}

func LinkLinks(ctx *gin.Context) {
	links := services.LinkService.Find(sqls.NewCnd().Eq("status", constants.StatusOk).Desc("id"))
	ctx.JSON(200, web.JsonData(render.BuildLinks(links)))
}

func SearchSearch(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)

	if len(keyword) == 0 {
		ctx.JSON(200, web.JsonData(gin.H{"results": []interface{}{}, "cursor": 0, "hasMore": false}))
		return
	}

	results, cursor, hasMore := services.SearchService.Search(keyword, cursor)
	ctx.JSON(200, web.JsonCursorData(results, strconv.FormatInt(cursor, 10), hasMore))
}

func FansRecent(ctx *gin.Context) {
	userId := params.FormValueInt64Default(ctx, "userId", 0)
	fans := services.UserFollowService.GetRecentFans(userId, 10)
	ctx.JSON(200, web.JsonData(render.BuildUsers(fans)))
}

func FansFollowRecent(ctx *gin.Context) {
	userId := params.FormValueInt64Default(ctx, "userId", 0)
	follows := services.UserFollowService.GetRecentFollows(userId, 10)
	ctx.JSON(200, web.JsonData(render.BuildUsers(follows)))
}

func UserReportCreate(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	var form struct {
		ReportType string `json:"reportType"`
		EntityType string `json:"entityType"`
		EntityId   int64  `json:"entityId"`
		Reason     string `json:"reason"`
	}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	err := services.UserReportService.Create(user.Id, form.ReportType, form.EntityType, form.EntityId, form.Reason)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func TaskTasks(ctx *gin.Context) {
	tasks := services.TaskConfigService.GetAll()
	ctx.JSON(200, web.JsonData(render.BuildTasks(tasks)))
}

func TaskUserTasks(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	tasks := services.UserTaskLogService.GetUserTasks(user.Id)
	ctx.JSON(200, web.JsonData(tasks))
}

func BadgeBadges(ctx *gin.Context) {
	badges := services.BadgeService.Find(sqls.NewCnd().Desc("id"))
	ctx.JSON(200, web.JsonData(render.BuildBadges(badges)))
}

func BadgeUserBadges(ctx *gin.Context) {
	userId := params.FormValueInt64Default(ctx, "userId", 0)
	badges := services.UserBadgeService.GetUserBadges(userId)
	ctx.JSON(200, web.JsonData(render.BuildUserBadges(badges)))
}

func VoteDetail(ctx *gin.Context) {
	voteIdStr := ctx.Param("id")
	voteId, _ := strconv.ParseInt(voteIdStr, 10, 64)
	vote := services.VoteService.Get(voteId)
	if vote == nil {
		ctx.JSON(200, web.JsonErrorMsg("投票不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(render.BuildVoteGin(ctx, vote)))
}

func VoteOptions(ctx *gin.Context) {
	voteIdStr := ctx.Param("id")
	voteId, _ := strconv.ParseInt(voteIdStr, 10, 64)
	options := services.VoteOptionService.FindByVoteId(voteId)
	ctx.JSON(200, web.JsonData(render.BuildVoteOptions(options)))
}

func VoteVote(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	voteIdStr := ctx.Param("id")
	voteId, _ := strconv.ParseInt(voteIdStr, 10, 64)

	var form struct {
		OptionIds []int64 `json:"optionIds"`
	}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	err := services.VoteService.Vote(user.Id, voteId, form.OptionIds)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func VoteVoted(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	voteIdStr := ctx.Param("id")
	voteId, _ := strconv.ParseInt(voteIdStr, 10, 64)

	if user == nil {
		ctx.JSON(200, web.JsonData(gin.H{"voted": false}))
		return
	}

	record := services.VoteRecordService.GetBy(user.Id, voteId)
	ctx.JSON(200, web.JsonData(gin.H{"voted": record != nil}))
}

func VoteUsers(ctx *gin.Context) {
	voteIdStr := ctx.Param("id")
	voteId, _ := strconv.ParseInt(voteIdStr, 10, 64)

	users := services.VoteRecordService.GetVoteUsers(voteId, 10)
	ctx.JSON(200, web.JsonData(render.BuildUsers(users)))
}

func AttachmentDetail(ctx *gin.Context) {
	attachmentId := ctx.Param("id")
	attachment := services.AttachmentService.Get(attachmentId)
	if attachment == nil {
		ctx.JSON(200, web.JsonErrorMsg("附件不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(render.BuildAttachment(attachment)))
}
