package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// defaultPage 默认页码
func defaultPage(ctx *gin.Context) int {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	return page
}

// defaultPageSize 默认每页数量
func defaultPageSize(ctx *gin.Context) int {
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageSize
}
