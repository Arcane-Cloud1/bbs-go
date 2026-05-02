package req

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/tidwall/gjson"
)

type CreateTopicForm struct {
	Type          constants.TopicType   `json:"type"`
	NodeId        int64                 `json:"nodeId"`
	Title         string                `json:"title"`
	Content       string                `json:"content"`
	ContentType   constants.ContentType `json:"contentType"`
	HideContent   string                `json:"hideContent"`
	Tags          []string              `json:"tags"`
	ImageList     []ImageDTO            `json:"imageList"`
	Vote          *CreateVoteForm       `json:"vote"`
	BountyScore   int                   `json:"bountyScore"`
	AttachmentIds []string              `json:"attachmentIds"`
	UserAgent     string                `json:"userAgent"`
	Ip            string                `json:"ip"`

	CaptchaId       string `json:"captchaId"`
	CaptchaCode     string `json:"captchaCode"`
	CaptchaProtocol int    `json:"captchaProtocol"`
}

type CreateVoteForm struct {
	Type      constants.VoteType     `json:"type"`
	Title     string                 `json:"title"`
	ExpiredAt int64                  `json:"expiredAt"`
	VoteNum   int                    `json:"voteNum"`
	Options   []CreateVoteOptionForm `json:"options"`
}

type CreateVoteOptionForm struct {
	Content string `json:"content"`
}

type VoteCastForm struct {
	VoteId    int64   `json:"voteId"`
	OptionIds []int64 `json:"optionIds"`
}

type EditTopicForm struct {
	NodeId        int64    `json:"nodeId"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	HideContent   string   `json:"hideContent"`
	Tags          []string `json:"tags"`
	AttachmentIds []string `json:"attachmentIds"`
}

type CreateArticleForm struct {
	Title       string
	Summary     string
	Content     string
	ContentType constants.ContentType
	Cover       *ImageDTO
	Tags        []string
	SourceUrl   string
}

type CreateCommentForm struct {
	EntityType string     `form:"entityType"`
	EntityId   int64      `form:"entityId"`
	Content    string     `form:"content"`
	ImageList  []ImageDTO `form:"imageList"`
	QuoteId    int64      `form:"quoteId"`
	UserAgent  string     `form:"userAgent"`
	Ip         string     `form:"ip"`
}

type ImageDTO struct {
	Url string `json:"url"`
}

func GetCreateTopicForm(ctx *gin.Context) CreateTopicForm {
	var form CreateTopicForm
	contentType := ctx.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := ctx.ShouldBindJSON(&form); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	} else {
		form = CreateTopicForm{
			Type:            constants.TopicType(getFormIntDefault(ctx, "type", int(constants.TopicTypeTopic))),
			NodeId:          getFormInt64Default(ctx, "nodeId", 0),
			Title:           strings.TrimSpace(ctx.PostForm("title")),
			Content:         strings.TrimSpace(ctx.PostForm("content")),
			ContentType:     constants.ContentType(ctx.PostForm("contentType")),
			HideContent:     strings.TrimSpace(ctx.PostForm("hideContent")),
			Tags:            ctx.PostFormArray("tags"),
			ImageList:       GetImageList(ctx, "imageList"),
			BountyScore:     getFormIntDefault(ctx, "bountyScore", 0),
			AttachmentIds:   ctx.PostFormArray("attachmentIds"),
			CaptchaId:       ctx.PostForm("captchaId"),
			CaptchaCode:     ctx.PostForm("captchaCode"),
			CaptchaProtocol: getFormIntDefault(ctx, "captchaProtocol", 0),
		}
	}

	if form.Type == constants.TopicTypeTweet {
		form.ContentType = constants.ContentTypeText
	}

	form.Ip = ctx.ClientIP()
	form.UserAgent = ctx.GetHeader("User-Agent")
	return form
}

func GetCreateCommentForm(ctx *gin.Context) CreateCommentForm {
	form := CreateCommentForm{
		EntityType: ctx.PostForm("entityType"),
		EntityId:   common.GetID(ctx, "entityId"),
		Content:    strings.TrimSpace(ctx.PostForm("content")),
		ImageList:  GetImageList(ctx, "imageList"),
		QuoteId:    getFormInt64Default(ctx, "quoteId", 0),
		UserAgent:  ctx.GetHeader("User-Agent"),
		Ip:         ctx.ClientIP(),
	}
	return form
}

func GetCreateArticleForm(ctx *gin.Context) CreateArticleForm {
	var (
		title   = ctx.PostForm("title")
		summary = ctx.PostForm("summary")
		content = ctx.PostForm("content")
		tags    = ctx.PostFormArray("tags")
		cover   = GetImageDTO(ctx, "cover")
	)
	return CreateArticleForm{
		Title:       title,
		Summary:     summary,
		Content:     content,
		ContentType: constants.ContentTypeMarkdown,
		Cover:       cover,
		Tags:        tags,
	}
}

func GetImageList(ctx *gin.Context, paramName string) []ImageDTO {
	imageListStr := ctx.PostForm(paramName)
	var imageList []ImageDTO
	if strs.IsNotBlank(imageListStr) {
		ret := gjson.Parse(imageListStr)
		if ret.IsArray() {
			for _, item := range ret.Array() {
				url := item.Get("url").String()
				imageList = append(imageList, ImageDTO{
					Url: url,
				})
			}
		}
	}
	return imageList
}

func GetImageDTO(ctx *gin.Context, paramName string) (img *ImageDTO) {
	str := ctx.PostForm(paramName)
	if strs.IsBlank(str) {
		return
	}
	if err := jsons.Parse(str, &img); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
	return
}

func getFormIntDefault(ctx *gin.Context, key string, defaultValue int) int {
	val := ctx.PostForm(key)
	if val == "" {
		return defaultValue
	}
	var result int
	for _, c := range val {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return defaultValue
		}
	}
	return result
}

func getFormInt64Default(ctx *gin.Context, key string, defaultValue int64) int64 {
	val := ctx.PostForm(key)
	if val == "" {
		return defaultValue
	}
	var result int64
	for _, c := range val {
		if c >= '0' && c <= '9' {
			result = result*10 + int64(c-'0')
		} else {
			return defaultValue
		}
	}
	return result
}
