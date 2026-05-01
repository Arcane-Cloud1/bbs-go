package services

import (
	"strings"
	"time"

	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
)

func (s *userTokenService) GetCurrentUserIdGin(ctx *gin.Context) int64 {
	user := s.GetCurrentGin(ctx)
	if user != nil {
		return user.Id
	}
	return 0
}

func (s *userTokenService) GetCurrentGin(ctx *gin.Context) *models.User {
	token := s.GetUserTokenGin(ctx)
	userToken := cache.UserTokenCache.Get(token)
	if userToken == nil || userToken.Status == constants.StatusDeleted {
		return nil
	}
	if userToken.ExpiredAt <= dates.NowTimestamp() {
		return nil
	}
	user := cache.UserCache.Get(userToken.UserId)
	if user == nil || user.Status != constants.StatusOk {
		return nil
	}

	trySendUserLoginEventGin(ctx, user.Id)

	return user
}

func trySendUserLoginEventGin(ctx *gin.Context, userId int64) {
	if ctx == nil || userId <= 0 {
		return
	}
	ctxKeyDailyVisitSent := "daily_visit_sent"
	if _, exists := ctx.Get(ctxKeyDailyVisitSent); exists {
		return
	}
	ctx.Set(ctxKeyDailyVisitSent, true)

	if !cache.DailyVisitCache.TryMarkAndReturnIfNew(userId) {
		return
	}
	event.Send(event.UserLoginEvent{
		UserId:     userId,
		LoginTime:  dates.NowTimestamp(),
		IsNewLogin: false,
	})
}

func (s *userTokenService) CheckLoginGin(ctx *gin.Context) (*models.User, error) {
	user := s.GetCurrentGin(ctx)
	if user == nil {
		return nil, errs.NotLogin()
	}
	return user, nil
}

func (s *userTokenService) SignoutGin(ctx *gin.Context) error {
	token := s.GetUserTokenGin(ctx)
	userToken := repositories.UserTokenRepository.GetByToken(sqls.DB(), token)
	if userToken == nil {
		return nil
	}
	err := repositories.UserTokenRepository.UpdateColumn(sqls.DB(), userToken.Id, "status", constants.StatusDeleted)
	if err != nil {
		return err
	}
	ctx.SetCookie(constants.CookieTokenKey, "", -1, "/", "", false, true)
	return nil
}

func (s *userTokenService) GetUserTokenGin(ctx *gin.Context) string {
	if userToken := ctx.Query("userToken"); strs.IsNotBlank(userToken) {
		return userToken
	}
	if userToken, err := ctx.Cookie(constants.CookieTokenKey); err == nil && strs.IsNotBlank(userToken) {
		return userToken
	}
	return s.getUserTokenFromHeaderGin(ctx)
}

func (s *userTokenService) getUserTokenFromHeaderGin(ctx *gin.Context) string {
	if authorization := ctx.GetHeader("Authorization"); strs.IsNotBlank(authorization) {
		userToken, _ := strings.CutPrefix(authorization, "Bearer ")
		return userToken
	}
	return ctx.GetHeader("X-User-Token")
}
