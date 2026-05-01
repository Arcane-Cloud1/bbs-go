package render

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/dates"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

func BuildLoginSuccessGin(ctx *gin.Context, user *models.User, redirect string) {
	if user == nil || user.Status != constants.StatusOk {
		ctx.JSON(200, gin.H{"code": -1, "msg": "用户不存在或已被禁用"})
		return
	}
	token, err := services.UserTokenService.Generate(user.Id)
	if err != nil {
		ctx.JSON(200, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	ctx.SetCookie(constants.CookieTokenKey, token, 365*24*3600, "/", "", false, true)

	event.Send(event.UserLoginEvent{
		UserId:     user.Id,
		LoginTime:  dates.NowTimestamp(),
		IsNewLogin: true,
	})

	cache.DailyVisitCache.MarkSentToday(user.Id)

	ctx.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"token":    token,
			"user":     BuildUserProfile(user),
			"redirect": redirect,
		},
	})
}
