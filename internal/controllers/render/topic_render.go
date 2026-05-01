package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/common"
	html2 "bbs-go/internal/pkg/html"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/pkg/text"
	"bbs-go/internal/services"
	"html"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/arrays"
	"github.com/mlogclub/simple/common/strs"
)

func BuildTopicGin(ctx *gin.Context, topic *models.Topic) *resp.TopicResponse {
	rsp := _buildTopic(topic, true)
	if rsp == nil {
		return nil
	}

	if currentUser := common.GetCurrentUser(ctx); currentUser != nil {
		rsp.Liked = services.UserLikeService.Exists(currentUser.Id, constants.EntityTopic, topic.Id)
		rsp.Favorited = services.FavoriteService.IsFavorited(currentUser.Id, constants.EntityTopic, topic.Id)
	}

	if vote := services.VoteService.Get(topic.VoteId); vote != nil {
		rsp.Vote = BuildVoteGin(ctx, vote)
	}

	list := services.AttachmentService.ListByTopicId(topic.Id)
	if len(list) > 0 {
		var currentUser *models.User
		if u := common.GetCurrentUser(ctx); u != nil {
			currentUser = u
		}
		rsp.Attachments = BuildAttachmentResponses(list, currentUser)
	}

	return rsp
}

func BuildSimpleTopicsGin(ctx *gin.Context, topics []models.Topic) []resp.TopicResponse {
	if len(topics) == 0 {
		return nil
	}

	var likedTopicIds []int64
	if currentUser := common.GetCurrentUser(ctx); currentUser != nil {
		var topicIds []int64
		for _, topic := range topics {
			topicIds = append(topicIds, topic.Id)
		}
		likedTopicIds = services.UserLikeService.IsLiked(currentUser.Id, constants.EntityTopic, topicIds)
	}

	var responses []resp.TopicResponse
	for _, topic := range topics {
		item := BuildSimpleTopic(&topic)
		item.Liked = arrays.Contains(topic.Id, likedTopicIds)
		if vote := services.VoteService.Get(topic.VoteId); vote != nil {
			item.Vote = BuildVoteGin(ctx, vote)
		}
		responses = append(responses, *item)
	}
	return responses
}
