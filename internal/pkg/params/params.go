package params

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func FormValueGin(ctx *gin.Context, name string) string {
	return ctx.PostForm(name)
}

func FormValueIntDefaultGin(ctx *gin.Context, name string, defaultValue int) int {
	return cast.ToInt(ctx.DefaultPostForm(name, cast.ToString(defaultValue)))
}

func FormValueInt64DefaultGin(ctx *gin.Context, name string, defaultValue int64) int64 {
	return cast.ToInt64(ctx.DefaultPostForm(name, cast.ToString(defaultValue)))
}

func FormValueBoolDefaultGin(ctx *gin.Context, name string, defaultValue bool) bool {
	return cast.ToBool(ctx.DefaultPostForm(name, cast.ToString(defaultValue)))
}

func FormValueInt64ArrayGin(ctx *gin.Context, name string) []int64 {
	values := ctx.PostFormArray(name)
	var result []int64
	for _, v := range values {
		result = append(result, cast.ToInt64(v))
	}
	return result
}

func GetInt64Gin(ctx *gin.Context, name string) (int64, error) {
	return cast.ToInt64E(ctx.Query(name))
}

func NewPagedSqlCndGin(ctx *gin.Context, filters ...QueryFilter) *PagedSqlCnd {
	cnd := NewPagedSqlCnd(ctx, filters...)
	return cnd
}
