package util

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/constants"
)

// PageQuery 分页查询参数。
type PageQuery struct {
	Page     int64
	PageSize int64
	Offset   int64
}

// ParsePageQuery 从 query 解析分页参数（page、page_size），带默认值与上限。
func ParsePageQuery(c *gin.Context) PageQuery {
	page := int64(constants.DefaultPage)
	pageSize := int64(constants.DefaultPageSize)
	if v, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64); err == nil && v > 0 {
		page = v
	}
	if v, err := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 64); err == nil && v > 0 {
		pageSize = v
	}
	if pageSize > constants.MaxPageSize {
		pageSize = constants.MaxPageSize
	}
	return PageQuery{Page: page, PageSize: pageSize, Offset: (page - 1) * pageSize}
}
