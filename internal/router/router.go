package router

import (
	"bbs-go/internal/controllers/admin"
	"bbs-go/internal/controllers/api"
	"bbs-go/internal/middleware"
	"bbs-go/internal/pkg/config"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	conf := config.Instance

	r.Use(cors.New(cors.Config{
		AllowOrigins:     conf.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           600,
	}))

	r.Use(middleware.AttachmentMiddleware)

	r.Static("/res", "./res")
	r.Static("/admin", "./admin")
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "Not found"})
			return
		}
		c.File("./site/index.html")
	})
	r.StaticFS("/site", http.Dir("./site"))

	setupAPIRoutes(r)
	setupAdminRoutes(r)
}

func setupAPIRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.InstallMiddleware)
	apiGroup.Use(middleware.AuthMiddleware)

	install := apiGroup.Group("/install")
	{
		install.GET("/info", api.InstallInfo)
		install.POST("/test_db", api.InstallTestDb)
		install.POST("/install", api.InstallInstall)
	}

	topic := apiGroup.Group("/topic")
	{
		topic.GET("/node_navs", api.TopicNodeNavs)
		topic.GET("/nodes", api.TopicNodes)
		topic.GET("/node", api.TopicNode)
		topic.POST("/create", api.TopicCreate)
		topic.GET("/edit/:id", api.TopicEdit)
		topic.POST("/edit/:id", api.TopicEditPost)
		topic.POST("/delete/:id", api.TopicDelete)
		topic.POST("/recommend/:id", api.TopicRecommend)
		topic.GET("/:id", api.TopicDetail)
		topic.GET("/recentlikes/:id", api.TopicRecentLikes)
		topic.GET("/recent", api.TopicRecent)
		topic.GET("/user_topics", api.TopicUserTopics)
		topic.GET("/topics", api.TopicTopics)
		topic.POST("/accept_answer/:id", api.TopicAcceptAnswer)
		topic.POST("/unaccept_answer/:id", api.TopicUnacceptAnswer)
		topic.GET("/tag_topics", api.TopicTagTopics)
		topic.GET("/favorite/:id", api.TopicFavorite)
		topic.POST("/sticky/:id", api.TopicSticky)
		topic.GET("/hide_content", api.TopicHideContent)
	}

	article := apiGroup.Group("/article")
	{
		article.GET("/tags", api.ArticleTags)
		article.GET("/recent", api.ArticleRecent)
		article.GET("/:id", api.ArticleDetail)
	}

	login := apiGroup.Group("/login")
	{
		login.POST("", api.LoginLogin)
		login.POST("/logout", api.LoginLogout)
		login.POST("/signup", api.LoginSignup)
		login.POST("/send_email_code", api.LoginSendEmailCode)
		login.POST("/reset_password", api.LoginResetPassword)
		login.GET("/github", api.LoginGithub)
		login.GET("/github/callback", api.LoginGithubCallback)
		login.GET("/google", api.LoginGoogle)
		login.GET("/google/callback", api.LoginGoogleCallback)
		login.GET("/wechat", api.LoginWechat)
		login.GET("/wechat/callback", api.LoginWechatCallback)
	}

	user := apiGroup.Group("/user")
	{
		user.GET("/current", api.UserCurrent)
		user.GET("/:id", api.UserDetail)
		user.POST("/update", api.UserUpdate)
		user.POST("/update_avatar", api.UserUpdateAvatar)
		user.POST("/update_background_image", api.UserUpdateBackgroundImage)
		user.GET("/topics", api.UserTopics)
		user.GET("/articles", api.UserArticles)
		user.GET("/comments", api.UserComments)
		user.GET("/favorites", api.UserFavorites)
		user.GET("/fans", api.UserFans)
		user.GET("/follows", api.UserFollows)
		user.POST("/follow/:id", api.UserFollow)
		user.POST("/unfollow/:id", api.UserUnfollow)
		user.GET("/messages", api.UserMessages)
		user.POST("/read_msg", api.UserReadMsg)
		user.POST("/read_all_msg", api.UserReadAllMsg)
		user.GET("/score_logs", api.UserScoreLogs)
		user.GET("/exp_logs", api.UserExpLogs)
		user.GET("/checkin_status", api.UserCheckinStatus)
		user.GET("/checkin_logs", api.UserCheckinLogs)
	}

	tag := apiGroup.Group("/tag")
	{
		tag.GET("/tags", api.TagTags)
		tag.POST("/create", api.TagCreate)
		tag.GET("/:id", api.TagDetail)
		tag.POST("/update/:id", api.TagUpdate)
		tag.POST("/delete/:id", api.TagDelete)
	}

	comment := apiGroup.Group("/comment")
	{
		comment.GET("/comments", api.CommentComments)
		comment.POST("/create", api.CommentCreate)
		comment.GET("/:id", api.CommentDetail)
		comment.POST("/delete/:id", api.CommentDelete)
	}

	favorite := apiGroup.Group("/favorite")
	{
		favorite.GET("/favorites", api.FavoriteFavorites)
		favorite.POST("/create", api.FavoriteCreate)
		favorite.POST("/delete/:id", api.FavoriteDelete)
	}

	like := apiGroup.Group("/like")
	{
		like.POST("/:entityType/:entityId", api.LikeLike)
		like.DELETE("/:entityType/:entityId", api.LikeUnlike)
		like.GET("/:entityType/:entityId/users", api.LikeUsers)
	}

	checkin := apiGroup.Group("/checkin")
	{
		checkin.POST("", api.CheckinCheckin)
		checkin.GET("/ranking", api.CheckinRanking)
	}

	config := apiGroup.Group("/config")
	{
		config.GET("/configs", api.ConfigConfigs)
		config.GET("/sys", api.ConfigSys)
	}

	upload := apiGroup.Group("/upload")
	{
		upload.POST("", api.UploadUpload)
	}

	attachment := apiGroup.Group("/attachment")
	{
		attachment.GET("/:id", api.AttachmentDetail)
	}

	link := apiGroup.Group("/link")
	{
		link.GET("/links", api.LinkLinks)
	}

	captcha := apiGroup.Group("/captcha")
	{
		captcha.GET("", api.CaptchaGet)
		captcha.POST("/verify", api.CaptchaVerify)
	}

	search := apiGroup.Group("/search")
	{
		search.GET("", api.SearchSearch)
	}

	fans := apiGroup.Group("/fans")
	{
		fans.GET("/recent", api.FansRecent)
		fans.GET("/follow_recent", api.FansFollowRecent)
	}

	userReport := apiGroup.Group("/user-report")
	{
		userReport.POST("/create", api.UserReportCreate)
	}

	task := apiGroup.Group("/task")
	{
		task.GET("/tasks", api.TaskTasks)
		task.GET("/user_tasks", api.TaskUserTasks)
	}

	badge := apiGroup.Group("/badge")
	{
		badge.GET("/badges", api.BadgeBadges)
		badge.GET("/user_badges", api.BadgeUserBadges)
	}

	vote := apiGroup.Group("/vote")
	{
		vote.GET("/:id", api.VoteDetail)
		vote.GET("/:id/options", api.VoteOptions)
		vote.POST("/:id/vote", api.VoteVote)
		vote.GET("/:id/voted", api.VoteVoted)
		vote.GET("/:id/users", api.VoteUsers)
	}
}

func setupAdminRoutes(r *gin.Engine) {
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.InstallMiddleware)
	adminGroup.Use(middleware.AuthMiddleware)
	adminGroup.Use(middleware.AdminMiddleware)

	role := adminGroup.Group("/role")
	{
		role.GET("/:id", admin.RoleGet)
		role.Any("/list", admin.RoleList)
		role.POST("/create", admin.RoleCreate)
		role.POST("/update", admin.RoleUpdate)
		role.POST("/delete/:id", admin.RoleDelete)
	}

	menu := adminGroup.Group("/menu")
	{
		menu.GET("/:id", admin.MenuGet)
		menu.Any("/list", admin.MenuList)
		menu.POST("/create", admin.MenuCreate)
		menu.POST("/update", admin.MenuUpdate)
		menu.POST("/delete/:id", admin.MenuDelete)
		menu.POST("/move", admin.MenuMove)
	}

	apiGroup := adminGroup.Group("/api")
	{
		apiGroup.GET("/:id", admin.ApiGet)
		apiGroup.Any("/list", admin.ApiList)
		apiGroup.POST("/create", admin.ApiCreate)
		apiGroup.POST("/update", admin.ApiUpdate)
		apiGroup.POST("/delete/:id", admin.ApiDelete)
	}

	dictType := adminGroup.Group("/dict-type")
	{
		dictType.GET("/:id", admin.DictTypeGet)
		dictType.Any("/list", admin.DictTypeList)
		dictType.POST("/create", admin.DictTypeCreate)
		dictType.POST("/update", admin.DictTypeUpdate)
		dictType.POST("/delete/:id", admin.DictTypeDelete)
	}

	dict := adminGroup.Group("/dict")
	{
		dict.GET("/:id", admin.DictGet)
		dict.Any("/list", admin.DictList)
		dict.POST("/create", admin.DictCreate)
		dict.POST("/update", admin.DictUpdate)
		dict.POST("/delete/:id", admin.DictDelete)
	}

	emailLog := adminGroup.Group("/email-log")
	{
		emailLog.GET("/:id", admin.EmailLogGet)
		emailLog.Any("/list", admin.EmailLogList)
	}

	common := adminGroup.Group("/common")
	{
		common.GET("/configs", admin.CommonConfigs)
		common.POST("/set_config", admin.CommonSetConfig)
		common.GET("/summary", admin.CommonSummary)
	}

	user := adminGroup.Group("/user")
	{
		user.GET("/synccount", admin.UserSynccount)
		user.GET("/:id", admin.UserGet)
		user.Any("/list", admin.UserList)
		user.POST("/create", admin.UserCreate)
		user.POST("/update", admin.UserUpdate)
		user.POST("/forbidden", admin.UserForbidden)
		user.POST("/update_password", admin.UserUpdatePassword)
		user.POST("/reset_password", admin.UserResetPassword)
	}

	tag := adminGroup.Group("/tag")
	{
		tag.GET("/:id", admin.TagGet)
		tag.Any("/list", admin.TagList)
		tag.POST("/create", admin.TagCreate)
		tag.POST("/update", admin.TagUpdate)
		tag.POST("/delete/:id", admin.TagDelete)
	}

	article := adminGroup.Group("/article")
	{
		article.GET("/:id", admin.ArticleGet)
		article.Any("/list", admin.ArticleList)
		article.POST("/update", admin.ArticleUpdate)
		article.POST("/delete/:id", admin.ArticleDelete)
	}

	comment := adminGroup.Group("/comment")
	{
		comment.GET("/:id", admin.CommentGet)
		comment.Any("/list", admin.CommentList)
		comment.POST("/delete/:id", admin.CommentDelete)
	}

	favorite := adminGroup.Group("/favorite")
	{
		favorite.GET("/:id", admin.FavoriteGet)
		favorite.Any("/list", admin.FavoriteList)
	}

	articleTag := adminGroup.Group("/article-tag")
	{
		articleTag.GET("/:id", admin.ArticleTagGet)
		articleTag.Any("/list", admin.ArticleTagList)
		articleTag.POST("/create", admin.ArticleTagCreate)
		articleTag.POST("/update", admin.ArticleTagUpdate)
		articleTag.POST("/delete/:id", admin.ArticleTagDelete)
	}

	topic := adminGroup.Group("/topic")
	{
		topic.GET("/:id", admin.TopicGet)
		topic.Any("/list", admin.TopicList)
		topic.POST("/update", admin.TopicUpdate)
		topic.POST("/delete/:id", admin.TopicDelete)
	}

	topicNode := adminGroup.Group("/topic-node")
	{
		topicNode.GET("/:id", admin.TopicNodeGet)
		topicNode.Any("/list", admin.TopicNodeList)
		topicNode.POST("/create", admin.TopicNodeCreate)
		topicNode.POST("/update", admin.TopicNodeUpdate)
		topicNode.POST("/delete/:id", admin.TopicNodeDelete)
	}

	sysConfig := adminGroup.Group("/sys-config")
	{
		sysConfig.GET("/:id", admin.SysConfigGet)
		sysConfig.Any("/list", admin.SysConfigList)
		sysConfig.POST("/save", admin.SysConfigSave)
	}

	link := adminGroup.Group("/link")
	{
		link.GET("/:id", admin.LinkGet)
		link.Any("/list", admin.LinkList)
		link.POST("/create", admin.LinkCreate)
		link.POST("/update", admin.LinkUpdate)
		link.POST("/delete/:id", admin.LinkDelete)
	}

	userScoreLog := adminGroup.Group("/user-score-log")
	{
		userScoreLog.GET("/:id", admin.UserScoreLogGet)
		userScoreLog.Any("/list", admin.UserScoreLogList)
	}

	taskConfig := adminGroup.Group("/task-config")
	{
		taskConfig.GET("/:id", admin.TaskConfigGet)
		taskConfig.Any("/list", admin.TaskConfigList)
		taskConfig.POST("/create", admin.TaskConfigCreate)
		taskConfig.POST("/update", admin.TaskConfigUpdate)
		taskConfig.POST("/delete/:id", admin.TaskConfigDelete)
	}

	badge := adminGroup.Group("/badge")
	{
		badge.GET("/:id", admin.BadgeGet)
		badge.Any("/list", admin.BadgeList)
		badge.POST("/create", admin.BadgeCreate)
		badge.POST("/update", admin.BadgeUpdate)
		badge.POST("/delete/:id", admin.BadgeDelete)
	}

	levelConfig := adminGroup.Group("/level-config")
	{
		levelConfig.GET("/:id", admin.LevelConfigGet)
		levelConfig.Any("/list", admin.LevelConfigList)
		levelConfig.POST("/create", admin.LevelConfigCreate)
		levelConfig.POST("/update", admin.LevelConfigUpdate)
		levelConfig.POST("/delete/:id", admin.LevelConfigDelete)
	}

	userTaskLog := adminGroup.Group("/user-task-log")
	{
		userTaskLog.GET("/:id", admin.UserTaskLogGet)
		userTaskLog.Any("/list", admin.UserTaskLogList)
	}

	userExpLog := adminGroup.Group("/user-exp-log")
	{
		userExpLog.GET("/:id", admin.UserExpLogGet)
		userExpLog.Any("/list", admin.UserExpLogList)
	}

	userBadge := adminGroup.Group("/user-badge")
	{
		userBadge.GET("/:id", admin.UserBadgeGet)
		userBadge.Any("/list", admin.UserBadgeList)
	}

	operateLog := adminGroup.Group("/operate-log")
	{
		operateLog.GET("/:id", admin.OperateLogGet)
		operateLog.Any("/list", admin.OperateLogList)
	}

	userReport := adminGroup.Group("/user-report")
	{
		userReport.GET("/:id", admin.UserReportGet)
		userReport.Any("/list", admin.UserReportList)
		userReport.POST("/delete/:id", admin.UserReportDelete)
	}

	forbiddenWord := adminGroup.Group("/forbidden-word")
	{
		forbiddenWord.GET("/:id", admin.ForbiddenWordGet)
		forbiddenWord.Any("/list", admin.ForbiddenWordList)
		forbiddenWord.POST("/create", admin.ForbiddenWordCreate)
		forbiddenWord.POST("/update", admin.ForbiddenWordUpdate)
		forbiddenWord.POST("/delete/:id", admin.ForbiddenWordDelete)
	}

	vote := adminGroup.Group("/vote")
	{
		vote.GET("/:id", admin.VoteGet)
		vote.Any("/list", admin.VoteList)
		vote.POST("/create", admin.VoteCreate)
		vote.POST("/update", admin.VoteUpdate)
		vote.POST("/delete/:id", admin.VoteDelete)
	}

	voteOption := adminGroup.Group("/vote-option")
	{
		voteOption.GET("/:id", admin.VoteOptionGet)
		voteOption.Any("/list", admin.VoteOptionList)
		voteOption.POST("/create", admin.VoteOptionCreate)
		voteOption.POST("/update", admin.VoteOptionUpdate)
		voteOption.POST("/delete/:id", admin.VoteOptionDelete)
	}

	voteRecord := adminGroup.Group("/vote-record")
	{
		voteRecord.GET("/:id", admin.VoteRecordGet)
		voteRecord.Any("/list", admin.VoteRecordList)
	}

	userTaskEvent := adminGroup.Group("/user-task-event")
	{
		userTaskEvent.GET("/:id", admin.UserTaskEventGet)
		userTaskEvent.Any("/list", admin.UserTaskEventList)
	}
}
