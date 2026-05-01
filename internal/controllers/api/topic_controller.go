package api

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/spam"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/models"
	"bbs-go/internal/models/req"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/services"
)

func TopicNodeNavs(ctx *gin.Context) {
	nodes := append(
		getBuiltInNodes(),
		render.BuildNodes(services.TopicNodeService.GetTopLevelNodes())...,
	)
	ctx.JSON(200, web.JsonData(nodes))
}

func getBuiltInNodes() []resp.NodeResponse {
	return []resp.NodeResponse{
		{
			Id:   0,
			Name: locales.Get("topic.node.latest"),
			Logo: "/res/images/node_latest.png",
		},
		{
			Id:   -1,
			Name: locales.Get("topic.node.recommend"),
			Logo: "/res/images/node_recommend.png",
		},
		{
			Id:   -2,
			Name: locales.Get("topic.node.follow"),
			Logo: "/res/images/node_follow.png",
		},
	}
}

func TopicNodes(ctx *gin.Context) {
	topicType := constants.TopicType(params.FormValueIntDefault(ctx, "type", -1))
	var nodeList []models.TopicNode
	if topicType >= 0 {
		nodeList = services.TopicNodeService.GetNodesByTopicType(topicType)
	} else {
		nodeList = services.TopicNodeService.GetNodes()
	}
	nodes := render.BuildNodeTree(0, nodeList)
	ctx.JSON(200, web.JsonData(nodes))
}

func TopicNode(ctx *gin.Context) {
	nodeId, _ := strconv.ParseInt(ctx.Query("nodeId"), 10, 64)
	if nodeId <= 0 {
		for _, node := range getBuiltInNodes() {
			if node.Id == nodeId {
				ctx.JSON(200, web.JsonData(node))
				return
			}
		}
	}
	node := services.TopicNodeService.Get(nodeId)
	if node == nil {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("common.not_found")))
		return
	}
	ctx.JSON(200, web.JsonData(render.BuildNodeWithChildren(node)))
}

func TopicCreate(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	var form req.CreateTopicForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	form.Title = strings.TrimSpace(form.Title)
	form.Content = strings.TrimSpace(form.Content)
	form.HideContent = strings.TrimSpace(form.HideContent)
	if constants.IsTweetTopicType(form.Type) {
		form.ContentType = constants.ContentTypeText
	}
	form.Ip = ctx.ClientIP()
	form.UserAgent = ctx.GetHeader("User-Agent")

	if err := spam.CheckTopic(user, form); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	topic, err := services.TopicPublishService.Publish(user.Id, form)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(render.BuildSimpleTopic(topic)))
}

func TopicEdit(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("common.not_found")))
		return
	}
	if constants.IsTweetTopicType(topic.Type) {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("topic.type_not_supported")))
		return
	}

	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("topic.no_permission")))
		return
	}

	tags := services.TopicService.GetTopicTags(topicId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	attachments := render.BuildAttachmentResponses(services.AttachmentService.ListByTopicId(topicId), nil)

	ctx.JSON(200, web.NewEmptyRspBuilder().
		Put("id", idcodec.Encode(topic.Id)).
		Put("type", topic.Type).
		Put("nodeId", topic.NodeId).
		Put("title", topic.Title).
		Put("content", topic.Content).
		Put("contentType", topic.ContentType).
		Put("hideContent", topic.HideContent).
		Put("tags", tagNames).
		Put("attachments", attachments).
		JsonResult())
}

func TopicEditPost(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("common.not_found")))
		return
	}

	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("topic.no_permission")))
		return
	}

	var form req.EditTopicForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	form.Title = strings.TrimSpace(form.Title)
	form.Content = strings.TrimSpace(form.Content)
	form.HideContent = strings.TrimSpace(form.HideContent)

	err := services.TopicService.Edit(user.Id, topicId, form)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(render.BuildSimpleTopic(topic)))
}

func TopicDelete(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		ctx.JSON(200, web.JsonSuccess())
		return
	}

	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("topic.no_permission")))
		return
	}

	if err := services.TopicService.Delete(topicId, user.Id, ctx.Request); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func TopicRecommend(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	topicId := idcodec.Decode(topicIdStr)
	recommend, err := params.FormValueBool(ctx, "recommend")
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("topic.no_permission")))
		return
	}

	err = services.TopicService.SetRecommend(topicId, recommend)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func TopicDetail(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	topicId := idcodec.Decode(topicIdStr)
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status == constants.StatusDeleted {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("common.not_found")))
		return
	}

	user := common.GetCurrentUser(ctx)
	if topic.Status == constants.StatusReview {
		if user != nil {
			if topic.UserId != user.Id && !user.IsOwnerOrAdmin() {
				ctx.JSON(200, web.JsonErrorCode(403, locales.Get("topic.under_review")))
				return
			}
		} else {
			ctx.JSON(200, web.JsonErrorCode(403, locales.Get("topic.under_review")))
			return
		}
	}

	services.TopicService.IncrViewCount(topicId)
	ctx.JSON(200, web.JsonData(render.BuildTopicGin(ctx, topic)))
}

func TopicRecentLikes(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	topicId := idcodec.Decode(topicIdStr)
	likes := services.UserLikeService.Recent(constants.EntityTopic, topicId, 5)
	var users []resp.UserInfo
	for _, like := range likes {
		userInfo := render.BuildUserInfoDefaultIfNull(like.UserId)
		if userInfo != nil {
			users = append(users, *userInfo)
		}
	}
	ctx.JSON(200, web.JsonData(users))
}

func TopicRecent(ctx *gin.Context) {
	topics := services.TopicService.Find(sqls.NewCnd().Where("status = ?", constants.StatusOk).Desc("id").Limit(10))
	ctx.JSON(200, web.JsonData(render.BuildSimpleTopicsGin(ctx, topics)))
}

func TopicUserTopics(ctx *gin.Context) {
	userId := common.GetQueryID(ctx, "userId")
	if userId <= 0 {
		ctx.JSON(200, web.JsonErrorMsg("param: userId required"))
		return
	}
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	topics, cursor, hasMore := services.TopicService.GetUserTopics(userId, cursor)
	ctx.JSON(200, web.JsonCursorData(render.BuildSimpleTopicsGin(ctx, topics), strconv.FormatInt(cursor, 10), hasMore))
}

func TopicTopics(ctx *gin.Context) {
	var (
		cursor   = params.FormValueInt64Default(ctx, "cursor", 0)
		nodeId   = params.FormValueInt64Default(ctx, "nodeId", 0)
		qaStatus = strings.TrimSpace(ctx.Query("qaStatus"))
		sort     = strings.TrimSpace(ctx.Query("sort"))
		user     = common.GetCurrentUser(ctx)
	)
	if nodeId == constants.NodeIdFollow && user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}

	var temp []models.Topic
	if cursor <= 0 {
		stickyTopics := services.TopicService.GetStickyTopics(nodeId, 3, qaStatus)
		temp = append(temp, stickyTopics...)
	}
	topics, cursor, hasMore := services.TopicService.GetTopics(user, nodeId, cursor, qaStatus, sort)
	for _, topic := range topics {
		topic.Sticky = false
		temp = append(temp, topic)
	}
	list := common.Distinct(temp, func(t models.Topic) any {
		return t.Id
	})
	ctx.JSON(200, web.JsonCursorData(render.BuildSimpleTopicsGin(ctx, list), strconv.FormatInt(cursor, 10), hasMore))
}

func TopicAcceptAnswer(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	topicId := idcodec.Decode(topicIdStr)
	commentId := params.FormValueInt64Default(ctx, "commentId", 0)
	if commentId <= 0 {
		ctx.JSON(200, web.JsonErrorMsg("commentId is required"))
		return
	}
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}
	if err := services.TopicService.AcceptAnswer(topicId, commentId, user.Id, user.IsOwnerOrAdmin()); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func TopicUnacceptAnswer(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}
	if err := services.TopicService.UnacceptAnswer(topicId, user.Id, user.IsOwnerOrAdmin()); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func TopicTagTopics(ctx *gin.Context) {
	var (
		cursor = params.FormValueInt64Default(ctx, "cursor", 0)
		tagId  = params.FormValueInt64Default(ctx, "tagId")
	)
	topics, cursor, hasMore := services.TopicService.GetTagTopics(tagId, cursor)
	ctx.JSON(200, web.JsonCursorData(render.BuildSimpleTopicsGin(ctx, topics), strconv.FormatInt(cursor, 10), hasMore))
}

func TopicFavorite(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}
	err := services.FavoriteService.AddTopicFavorite(user.Id, topicId)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func TopicSticky(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(200, web.JsonError(errs.NotLogin()))
		return
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("topic.no_permission")))
		return
	}

	sticky := params.FormValueBoolDefault(ctx, "sticky", false)
	if err := services.TopicService.SetSticky(topicId, sticky); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func TopicHideContent(ctx *gin.Context) {
	topicId := common.GetQueryID(ctx, "topicId")
	var (
		exists      = false
		show        = false
		hideContent = ""
	)
	topic := services.TopicService.Get(topicId)
	if topic != nil && topic.Status == constants.StatusOk && strs.IsNotBlank(topic.HideContent) {
		exists = true
		if user := common.GetCurrentUser(ctx); user != nil {
			if user.Id == topic.UserId || services.CommentService.IsCommented(user.Id, constants.EntityTopic, topic.Id) {
				show = true
				hideContent = markdown.ToHTML(topic.HideContent)
			}
		}
	}
	ctx.JSON(200, web.JsonData(map[string]interface{}{
		"exists":  exists,
		"show":    show,
		"content": hideContent,
	}))
}
