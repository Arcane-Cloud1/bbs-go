package api

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"database/sql"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/pkg/bbsurls"
	captcha2 "bbs-go/internal/pkg/captcha"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/services"
)

func LoginLogin(ctx *gin.Context) {
	var (
		captchaId          = ctx.PostForm("captchaId")
		captchaCode        = ctx.PostForm("captchaCode")
		captchaProtocol, _ = strconv.Atoi(ctx.PostForm("captchaProtocol"))
		username           = ctx.PostForm("username")
		password           = ctx.PostForm("password")
		redirect           = ctx.PostForm("redirect")
	)

	if captchaProtocol == 2 {
		if !captcha2.Verify(captchaId, captchaCode) {
			ctx.JSON(200, web.JsonError(errs.CaptchaError()))
			return
		}
	} else {
		if !captcha.VerifyString(captchaId, captchaCode) {
			ctx.JSON(200, web.JsonError(errs.CaptchaError()))
			return
		}
	}

	user, err := services.UserService.SignIn(username, password)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	if !user.IsOwnerOrAdmin() {
		if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
			ctx.JSON(200, web.JsonErrorMsg(locales.Get("auth.password_login_disabled")))
			return
		}
	}
	render.BuildLoginSuccessGin(ctx, user, redirect)
}

func LoginLogout(ctx *gin.Context) {
	err := services.UserTokenService.SignoutGin(ctx)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func LoginSignup(ctx *gin.Context) {
	var (
		captchaId          = ctx.PostForm("captchaId")
		captchaCode        = ctx.PostForm("captchaCode")
		captchaProtocol, _ = strconv.Atoi(ctx.PostForm("captchaProtocol"))
		email              = ctx.PostForm("email")
		username           = ctx.PostForm("username")
		password           = ctx.PostForm("password")
		rePassword         = ctx.PostForm("rePassword")
		nickname           = ctx.PostForm("nickname")
		redirect           = ctx.PostForm("redirect")
	)
	if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("auth.password_login_disabled")))
		return
	}
	if captchaProtocol == 2 {
		if !captcha2.Verify(captchaId, captchaCode) {
			ctx.JSON(200, web.JsonError(errs.CaptchaError()))
			return
		}
	} else {
		if !captcha.VerifyString(captchaId, captchaCode) {
			ctx.JSON(200, web.JsonError(errs.CaptchaError()))
			return
		}
	}
	user, err := services.UserService.SignUp(username, email, nickname, password, rePassword)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	render.BuildLoginSuccessGin(ctx, user, redirect)
}

func LoginSendEmailCode(ctx *gin.Context) {
	ctx.JSON(200, web.JsonSuccess())
}

func LoginResetPassword(ctx *gin.Context) {
	if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("auth.password_login_disabled")))
		return
	}
	var (
		token      = ctx.PostForm("token")
		password   = ctx.PostForm("password")
		rePassword = ctx.PostForm("rePassword")
	)
	if err := services.UserService.ResetPasswordByToken(token, password, rePassword); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func LoginGithub(ctx *gin.Context) {
	redirect := ctx.Query("redirect")
	bind := ctx.Query("bind") == "true"
	state := strs.UUID()

	loginConfig := services.SysConfigService.GetLoginConfig()
	if !loginConfig.GithubLogin.Enabled {
		ctx.JSON(200, web.JsonErrorMsg("GitHub登录未启用"))
		return
	}

	if bind && strs.IsBlank(redirect) {
		redirect = "/user/profile/account"
	}

	cache.GithubLoginStateCache.Put(state, &cache.GithubLoginStateData{
		Redirect: redirect,
		Bind:     bind,
	})

	redirectURI := bbsurls.AbsUrl("/api/login/github/callback")
	oauth := github.NewGithubOAuth(loginConfig.GithubLogin.ClientId, loginConfig.GithubLogin.ClientSecret, redirectURI)
	authURL := oauth.GetAuthURL(state)

	ctx.JSON(200, web.JsonData(gin.H{
		"clientId":    loginConfig.GithubLogin.ClientId,
		"authUrl":     authURL,
		"redirectUri": redirectURI,
		"state":       state,
		"redirect":    redirect,
	}))
}

func LoginGithubCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if strs.IsBlank(state) {
		ctx.JSON(200, web.JsonErrorMsg("state参数缺失"))
		return
	}

	data := cache.GithubLoginStateCache.Get(state)
	if data == nil {
		ctx.JSON(200, web.JsonErrorMsg("登录数据错误或已过期，请重新登录"))
		return
	}

	if !services.SysConfigService.GetLoginConfig().GithubLogin.Enabled {
		ctx.JSON(200, web.JsonErrorMsg("GitHub登录未启用"))
		return
	}

	if data.Bind {
		user, err := common.CheckLogin(ctx)
		if err != nil {
			ctx.JSON(200, web.JsonError(err))
			return
		}
		if err := services.ThirdUserService.BindGithub(user.Id, code, state); err != nil {
			ctx.JSON(200, web.JsonError(err))
			return
		}
		render.BuildLoginSuccessGin(ctx, user, data.Redirect)
	} else {
		user, err := services.ThirdUserService.LoginGithub(code, state)
		if err != nil {
			ctx.JSON(200, web.JsonError(err))
			return
		}
		render.BuildLoginSuccessGin(ctx, user, data.Redirect)
	}
}

func LoginGoogle(ctx *gin.Context) {
	redirect := ctx.Query("redirect")
	bind := ctx.Query("bind") == "true"
	state := strs.UUID()

	loginConfig := services.SysConfigService.GetLoginConfig()
	if !loginConfig.GoogleLogin.Enabled {
		ctx.JSON(200, web.JsonErrorMsg("Google登录未启用"))
		return
	}

	cache.GoogleLoginStateCache.Put(state, &cache.GoogleLoginStateData{
		Redirect: redirect,
		Bind:     bind,
	})

	redirectURI := bbsurls.AbsUrl("/api/login/google/callback")
	oauth := google.NewGoogleOAuth(loginConfig.GoogleLogin.ClientId, loginConfig.GoogleLogin.ClientSecret, redirectURI)
	authURL := oauth.GetAuthURL(state)

	ctx.JSON(200, web.JsonData(gin.H{
		"clientId":    loginConfig.GoogleLogin.ClientId,
		"authUrl":     authURL,
		"redirectUri": redirectURI,
		"state":       state,
		"redirect":    redirect,
	}))
}

func LoginGoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if strs.IsBlank(state) {
		ctx.JSON(200, web.JsonErrorMsg("state参数缺失"))
		return
	}

	data := cache.GoogleLoginStateCache.Get(state)
	if data == nil {
		ctx.JSON(200, web.JsonErrorMsg("登录数据错误或已过期，请重新登录"))
		return
	}

	if !services.SysConfigService.GetLoginConfig().GoogleLogin.Enabled {
		ctx.JSON(200, web.JsonErrorMsg("Google登录未启用"))
		return
	}

	user, err := services.ThirdUserService.LoginGoogle(code, state)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	render.BuildLoginSuccessGin(ctx, user, data.Redirect)
}

func LoginWechat(ctx *gin.Context) {
	redirect := ctx.Query("redirect")
	bind := ctx.Query("bind") == "true"
	state := strs.UUID()

	loginConfig := services.SysConfigService.GetLoginConfig()
	if !loginConfig.WeixinLogin.Enabled {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("auth.weixin_login_disabled")))
		return
	}

	cache.WxLoginStateCache.Put(state, &cache.WxLoginStateData{
		Redirect: redirect,
		Bind:     bind,
	})

	redirectURI := bbsurls.AbsUrl("/api/login/wechat/callback")
	if bind {
		redirectURI = bbsUrls.AbsUrl("/api/login/wechat/bind/callback")
	}

	ctx.JSON(200, web.JsonData(gin.H{
		"appid":        loginConfig.WeixinLogin.AppId,
		"scope":        "snsapi_login",
		"redirect_uri": redirectURI,
		"state":        state,
	}))
}

func LoginWechatCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	data := cache.WxLoginStateCache.Get(state)
	if data == nil {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("auth.login_data_error")))
		return
	}

	if !services.SysConfigService.GetLoginConfig().WeixinLogin.Enabled {
		ctx.JSON(200, web.JsonErrorMsg(locales.Get("auth.weixin_login_disabled")))
		return
	}

	user, err := services.ThirdUserService.LoginWeixin(code, state)
	if err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}

	render.BuildLoginSuccessGin(ctx, user, data.Redirect)
}
