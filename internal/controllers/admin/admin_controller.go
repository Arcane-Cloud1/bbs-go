package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

func RoleGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	role := services.RoleService.Get(id)
	if role == nil {
		ctx.JSON(200, web.JsonErrorMsg("角色不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(role))
}

func RoleList(ctx *gin.Context) {
	list, paging := services.RoleService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func RoleCreate(ctx *gin.Context) {
	var role models.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.RoleService.Create(&role); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(role))
}

func RoleUpdate(ctx *gin.Context) {
	var role models.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.RoleService.Update(&role); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(role))
}

func RoleDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.RoleService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func MenuGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	menu := services.MenuService.Get(id)
	if menu == nil {
		ctx.JSON(200, web.JsonErrorMsg("菜单不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(menu))
}

func MenuList(ctx *gin.Context) {
	list, paging := services.MenuService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func MenuCreate(ctx *gin.Context) {
	var menu models.Menu
	if err := ctx.ShouldBindJSON(&menu); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.MenuService.Create(&menu); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(menu))
}

func MenuUpdate(ctx *gin.Context) {
	var menu models.Menu
	if err := ctx.ShouldBindJSON(&menu); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.MenuService.Update(&menu); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(menu))
}

func MenuDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.MenuService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func MenuMove(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.PostForm("id"), 10, 64)
	parentId, _ := strconv.ParseInt(ctx.PostForm("parentId"), 10, 64)
	if err := services.MenuService.Move(id, parentId); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func ApiGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	api := services.ApiService.Get(id)
	if api == nil {
		ctx.JSON(200, web.JsonErrorMsg("API不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(api))
}

func ApiList(ctx *gin.Context) {
	list, paging := services.ApiService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func ApiCreate(ctx *gin.Context) {
	var api models.Api
	if err := ctx.ShouldBindJSON(&api); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.ApiService.Create(&api); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(api))
}

func ApiUpdate(ctx *gin.Context) {
	var api models.Api
	if err := ctx.ShouldBindJSON(&api); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.ApiService.Update(&api); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(api))
}

func ApiDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.ApiService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func DictTypeGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	dictType := services.DictTypeService.Get(id)
	if dictType == nil {
		ctx.JSON(200, web.JsonErrorMsg("字典类型不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(dictType))
}

func DictTypeList(ctx *gin.Context) {
	list, paging := services.DictTypeService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func DictTypeCreate(ctx *gin.Context) {
	var dictType models.DictType
	if err := ctx.ShouldBindJSON(&dictType); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.DictTypeService.Create(&dictType); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(dictType))
}

func DictTypeUpdate(ctx *gin.Context) {
	var dictType models.DictType
	if err := ctx.ShouldBindJSON(&dictType); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.DictTypeService.Update(&dictType); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(dictType))
}

func DictTypeDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.DictTypeService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func DictGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	dict := services.DictService.Get(id)
	if dict == nil {
		ctx.JSON(200, web.JsonErrorMsg("字典不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(dict))
}

func DictList(ctx *gin.Context) {
	list, paging := services.DictService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func DictCreate(ctx *gin.Context) {
	var dict models.Dict
	if err := ctx.ShouldBindJSON(&dict); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.DictService.Create(&dict); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(dict))
}

func DictUpdate(ctx *gin.Context) {
	var dict models.Dict
	if err := ctx.ShouldBindJSON(&dict); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.DictService.Update(&dict); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(dict))
}

func DictDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.DictService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func EmailLogGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	log := services.EmailLogService.Get(id)
	if log == nil {
		ctx.JSON(200, web.JsonErrorMsg("邮件日志不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(log))
}

func EmailLogList(ctx *gin.Context) {
	list, paging := services.EmailLogService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func CommonConfigs(ctx *gin.Context) {
	configs := services.SysConfigService.GetConfigs()
	ctx.JSON(200, web.JsonData(configs))
}

func CommonSetConfig(ctx *gin.Context) {
	key := ctx.PostForm("key")
	value := ctx.PostForm("value")

	if err := services.SysConfigService.Set(key, value); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func CommonSummary(ctx *gin.Context) {
	summary := services.SysConfigService.GetSummary()
	ctx.JSON(200, web.JsonData(summary))
}

func TagGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	tag := services.TagService.Get(id)
	if tag == nil {
		ctx.JSON(200, web.JsonErrorMsg("标签不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(tag))
}

func TagList(ctx *gin.Context) {
	list, paging := services.TagService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func TagCreate(ctx *gin.Context) {
	var tag models.Tag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.TagService.Create(&tag); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(tag))
}

func TagUpdate(ctx *gin.Context) {
	var tag models.Tag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.TagService.Update(&tag); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(tag))
}

func TagDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.TagService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func ArticleGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	article := services.ArticleService.Get(id)
	if article == nil {
		ctx.JSON(200, web.JsonErrorMsg("文章不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(article))
}

func ArticleList(ctx *gin.Context) {
	list, paging := services.ArticleService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func ArticleUpdate(ctx *gin.Context) {
	var article models.Article
	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.ArticleService.Update(&article); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(article))
}

func ArticleDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.ArticleService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func CommentGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	comment := services.CommentService.Get(id)
	if comment == nil {
		ctx.JSON(200, web.JsonErrorMsg("评论不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(comment))
}

func CommentList(ctx *gin.Context) {
	list, paging := services.CommentService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func CommentDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.CommentService.DeleteByAdmin(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func FavoriteGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	favorite := services.FavoriteService.Get(id)
	if favorite == nil {
		ctx.JSON(200, web.JsonErrorMsg("收藏不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(favorite))
}

func FavoriteList(ctx *gin.Context) {
	list, paging := services.FavoriteService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func ArticleTagGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	tag := services.ArticleTagService.Get(id)
	if tag == nil {
		ctx.JSON(200, web.JsonErrorMsg("标签不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(tag))
}

func ArticleTagList(ctx *gin.Context) {
	list, paging := services.ArticleTagService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func ArticleTagCreate(ctx *gin.Context) {
	var tag models.ArticleTag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.ArticleTagService.Create(&tag); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(tag))
}

func ArticleTagUpdate(ctx *gin.Context) {
	var tag models.ArticleTag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.ArticleTagService.Update(&tag); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(tag))
}

func ArticleTagDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.ArticleTagService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func TopicGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	topic := services.TopicService.Get(id)
	if topic == nil {
		ctx.JSON(200, web.JsonErrorMsg("帖子不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(topic))
}

func TopicList(ctx *gin.Context) {
	list, paging := services.TopicService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func TopicUpdate(ctx *gin.Context) {
	var topic models.Topic
	if err := ctx.ShouldBindJSON(&topic); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.TopicService.Update(&topic); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(topic))
}

func TopicDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.TopicService.DeleteByAdmin(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func TopicNodeGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	node := services.TopicNodeService.Get(id)
	if node == nil {
		ctx.JSON(200, web.JsonErrorMsg("节点不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(node))
}

func TopicNodeList(ctx *gin.Context) {
	list, paging := services.TopicNodeService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func TopicNodeCreate(ctx *gin.Context) {
	var node models.TopicNode
	if err := ctx.ShouldBindJSON(&node); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.TopicNodeService.Create(&node); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(node))
}

func TopicNodeUpdate(ctx *gin.Context) {
	var node models.TopicNode
	if err := ctx.ShouldBindJSON(&node); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.TopicNodeService.Update(&node); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(node))
}

func TopicNodeDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.TopicNodeService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func SysConfigGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	config := services.SysConfigService.Get(id)
	if config == nil {
		ctx.JSON(200, web.JsonErrorMsg("配置不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(config))
}

func SysConfigList(ctx *gin.Context) {
	list, paging := services.SysConfigService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func SysConfigSave(ctx *gin.Context) {
	var config models.SysConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.SysConfigService.Save(&config); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(config))
}

func LinkGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	link := services.LinkService.Get(id)
	if link == nil {
		ctx.JSON(200, web.JsonErrorMsg("链接不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(link))
}

func LinkList(ctx *gin.Context) {
	list, paging := services.LinkService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func LinkCreate(ctx *gin.Context) {
	var link models.Link
	if err := ctx.ShouldBindJSON(&link); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.LinkService.Create(&link); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(link))
}

func LinkUpdate(ctx *gin.Context) {
	var link models.Link
	if err := ctx.ShouldBindJSON(&link); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.LinkService.Update(&link); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(link))
}

func LinkDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.LinkService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserScoreLogGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	log := services.UserScoreLogService.Get(id)
	if log == nil {
		ctx.JSON(200, web.JsonErrorMsg("记录不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(log))
}

func UserScoreLogList(ctx *gin.Context) {
	list, paging := services.UserScoreLogService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func TaskConfigGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	config := services.TaskConfigService.Get(id)
	if config == nil {
		ctx.JSON(200, web.JsonErrorMsg("任务配置不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(config))
}

func TaskConfigList(ctx *gin.Context) {
	list, paging := services.TaskConfigService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func TaskConfigCreate(ctx *gin.Context) {
	var config models.TaskConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.TaskConfigService.Create(&config); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(config))
}

func TaskConfigUpdate(ctx *gin.Context) {
	var config models.TaskConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.TaskConfigService.Update(&config); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(config))
}

func TaskConfigDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.TaskConfigService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func BadgeGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	badge := services.BadgeService.Get(id)
	if badge == nil {
		ctx.JSON(200, web.JsonErrorMsg("徽章不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(badge))
}

func BadgeList(ctx *gin.Context) {
	list, paging := services.BadgeService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func BadgeCreate(ctx *gin.Context) {
	var badge models.Badge
	if err := ctx.ShouldBindJSON(&badge); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.BadgeService.Create(&badge); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(badge))
}

func BadgeUpdate(ctx *gin.Context) {
	var badge models.Badge
	if err := ctx.ShouldBindJSON(&badge); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.BadgeService.Update(&badge); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(badge))
}

func BadgeDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.BadgeService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func LevelConfigGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	config := services.LevelConfigService.Get(id)
	if config == nil {
		ctx.JSON(200, web.JsonErrorMsg("等级配置不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(config))
}

func LevelConfigList(ctx *gin.Context) {
	list, paging := services.LevelConfigService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func LevelConfigCreate(ctx *gin.Context) {
	var config models.LevelConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.LevelConfigService.Create(&config); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(config))
}

func LevelConfigUpdate(ctx *gin.Context) {
	var config models.LevelConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.LevelConfigService.Update(&config); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(config))
}

func LevelConfigDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.LevelConfigService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func UserTaskLogGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	log := services.UserTaskLogService.Get(id)
	if log == nil {
		ctx.JSON(200, web.JsonErrorMsg("记录不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(log))
}

func UserTaskLogList(ctx *gin.Context) {
	list, paging := services.UserTaskLogService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func UserExpLogGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	log := services.UserExpLogService.Get(id)
	if log == nil {
		ctx.JSON(200, web.JsonErrorMsg("记录不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(log))
}

func UserExpLogList(ctx *gin.Context) {
	list, paging := services.UserExpLogService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func UserBadgeGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	userBadge := services.UserBadgeService.Get(id)
	if userBadge == nil {
		ctx.JSON(200, web.JsonErrorMsg("记录不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(userBadge))
}

func UserBadgeList(ctx *gin.Context) {
	list, paging := services.UserBadgeService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func OperateLogGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	log := services.OperateLogService.Get(id)
	if log == nil {
		ctx.JSON(200, web.JsonErrorMsg("操作日志不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(log))
}

func OperateLogList(ctx *gin.Context) {
	list, paging := services.OperateLogService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func UserReportGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	report := services.UserReportService.Get(id)
	if report == nil {
		ctx.JSON(200, web.JsonErrorMsg("举报不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(report))
}

func UserReportList(ctx *gin.Context) {
	list, paging := services.UserReportService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func UserReportDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.UserReportService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func ForbiddenWordGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	word := services.ForbiddenWordService.Get(id)
	if word == nil {
		ctx.JSON(200, web.JsonErrorMsg("敏感词不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(word))
}

func ForbiddenWordList(ctx *gin.Context) {
	list, paging := services.ForbiddenWordService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func ForbiddenWordCreate(ctx *gin.Context) {
	var word models.ForbiddenWord
	if err := ctx.ShouldBindJSON(&word); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.ForbiddenWordService.Create(&word); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(word))
}

func ForbiddenWordUpdate(ctx *gin.Context) {
	var word models.ForbiddenWord
	if err := ctx.ShouldBindJSON(&word); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.ForbiddenWordService.Update(&word); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(word))
}

func ForbiddenWordDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.ForbiddenWordService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func VoteGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	vote := services.VoteService.Get(id)
	if vote == nil {
		ctx.JSON(200, web.JsonErrorMsg("投票不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(vote))
}

func VoteList(ctx *gin.Context) {
	list, paging := services.VoteService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func VoteCreate(ctx *gin.Context) {
	var vote models.Vote
	if err := ctx.ShouldBindJSON(&vote); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.VoteService.Create(&vote); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(vote))
}

func VoteUpdate(ctx *gin.Context) {
	var vote models.Vote
	if err := ctx.ShouldBindJSON(&vote); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.VoteService.Update(&vote); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(vote))
}

func VoteDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.VoteService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func VoteOptionGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	option := services.VoteOptionService.Get(id)
	if option == nil {
		ctx.JSON(200, web.JsonErrorMsg("投票选项不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(option))
}

func VoteOptionList(ctx *gin.Context) {
	list, paging := services.VoteOptionService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func VoteOptionCreate(ctx *gin.Context) {
	var option models.VoteOption
	if err := ctx.ShouldBindJSON(&option); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.VoteOptionService.Create(&option); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(option))
}

func VoteOptionUpdate(ctx *gin.Context) {
	var option models.VoteOption
	if err := ctx.ShouldBindJSON(&option); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	if err := services.VoteOptionService.Update(&option); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonData(option))
}

func VoteOptionDelete(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := services.VoteOptionService.Delete(id); err != nil {
		ctx.JSON(200, web.JsonError(err))
		return
	}
	ctx.JSON(200, web.JsonSuccess())
}

func VoteRecordGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	record := services.VoteRecordService.Get(id)
	if record == nil {
		ctx.JSON(200, web.JsonErrorMsg("投票记录不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(record))
}

func VoteRecordList(ctx *gin.Context) {
	list, paging := services.VoteRecordService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}

func UserTaskEventGet(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	event := services.UserTaskEventService.Get(id)
	if event == nil {
		ctx.JSON(200, web.JsonErrorMsg("记录不存在"))
		return
	}
	ctx.JSON(200, web.JsonData(event))
}

func UserTaskEventList(ctx *gin.Context) {
	list, paging := services.UserTaskEventService.FindPageByCnd(params.NewPagedSqlCndGin(ctx).Desc("id"))
	ctx.JSON(200, web.JsonData(&web.PageResult{Results: list, Page: paging}))
}
