package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FamilyAlbumController 家庭相册控制器
type FamilyAlbumController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册家庭相册路由
func (c *FamilyAlbumController) RegisterRoutes(api *gin.RouterGroup) {
	album := api.Group("/family-album")
	{
		album.GET("/containers", c.ListContainers)
		album.GET("/containers/:id", c.GetContainer)
		album.POST("/containers", c.CreateContainer)
		album.PUT("/containers/:id", c.UpdateContainer)
		album.DELETE("/containers/:id", c.DeleteContainer)

		album.GET("/items", c.ListItems)
		album.GET("/items/:id", c.GetItem)
		album.POST("/items", c.CreateItem)
		album.PUT("/items/:id", c.UpdateItem)
		album.DELETE("/items/:id", c.DeleteItem)
		album.POST("/items/:id/like", c.LikeItem)
		album.DELETE("/items/:id/like", c.UnlikeItem)
		album.POST("/items/:id/share", c.ShareItem)

		album.GET("/comments", c.ListComments)
		album.POST("/comments", c.CreateComment)
		album.DELETE("/comments/:id", c.DeleteComment)
	}
}

// ListContainers 获取相册容器列表
func (c *FamilyAlbumController) ListContainers(ctx *gin.Context) {
	var containers []models.FamilyAlbumContainer
	var total int64

	query := c.DB.Model(&models.FamilyAlbumContainer{})

	// 用户筛选
	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 可见性筛选
	if visibility := ctx.Query("visibility"); visibility != "" {
		query = query.Where("visibility = ?", visibility)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&containers).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      containers,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetContainer 获取相册容器详情
func (c *FamilyAlbumController) GetContainer(ctx *gin.Context) {
	id := ctx.Param("id")
	var container models.FamilyAlbumContainer
	if err := c.DB.First(&container, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相册不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 获取前5个相片项作为预览
	var previewItems []models.FamilyAlbumItem
	c.DB.Where("container_id = ?", id).Order("created_at DESC").Limit(5).Find(&previewItems)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"container": container,
		"preview":   previewItems,
	}})
}

// CreateContainer 创建相册容器
func (c *FamilyAlbumController) CreateContainer(ctx *gin.Context) {
	var container models.FamilyAlbumContainer
	if err := ctx.ShouldBindJSON(&container); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Create(&container).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": container})
}

// UpdateContainer 更新相册容器
func (c *FamilyAlbumController) UpdateContainer(ctx *gin.Context) {
	id := ctx.Param("id")
	var container models.FamilyAlbumContainer
	if err := c.DB.First(&container, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相册不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var updateData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		CoverURL    string `json:"cover_url"`
		Visibility  string `json:"visibility"`
		IsDefault   *bool  `json:"is_default"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if updateData.Name != "" {
		updates["name"] = updateData.Name
	}
	if updateData.Description != "" {
		updates["description"] = updateData.Description
	}
	if updateData.CoverURL != "" {
		updates["cover_url"] = updateData.CoverURL
	}
	if updateData.Visibility != "" {
		updates["visibility"] = updateData.Visibility
	}
	if updateData.IsDefault != nil {
		updates["is_default"] = *updateData.IsDefault
	}

	if err := c.DB.Model(&container).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&container, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": container})
}

// DeleteContainer 删除相册容器
func (c *FamilyAlbumController) DeleteContainer(ctx *gin.Context) {
	id := ctx.Param("id")
	var container models.FamilyAlbumContainer
	if err := c.DB.First(&container, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相册不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 删除关联的相片项
	c.DB.Where("container_id = ?", id).Delete(&models.FamilyAlbumItem{})

	if err := c.DB.Delete(&container).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ListItems 获取相册相片列表
func (c *FamilyAlbumController) ListItems(ctx *gin.Context) {
	var items []models.FamilyAlbumItem
	var total int64

	query := c.DB.Model(&models.FamilyAlbumItem{})

	// 容器筛选
	if containerID := ctx.Query("container_id"); containerID != "" {
		query = query.Where("container_id = ?", containerID)
	}

	// 用户筛选
	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 类型筛选
	if mediaType := ctx.Query("media_type"); mediaType != "" {
		query = query.Where("media_type = ?", mediaType)
	}

	// 分类筛选
	if category := ctx.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// 宠物筛选
	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      items,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetItem 获取相片详情
func (c *FamilyAlbumController) GetItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item models.FamilyAlbumItem
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相片不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 增加浏览次数
	c.DB.Model(&item).Update("view_count", item.ViewCount+1)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": item})
}

// CreateItem 创建相片
func (c *FamilyAlbumController) CreateItem(ctx *gin.Context) {
	var item models.FamilyAlbumItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 更新容器计数
	c.DB.Model(&models.FamilyAlbumContainer{}).Where("id = ?", item.ContainerID).
		Update("item_count", gorm.Expr("item_count + 1"))

	if err := c.DB.Create(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": item})
}

// UpdateItem 更新相片
func (c *FamilyAlbumController) UpdateItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item models.FamilyAlbumItem
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相片不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var updateData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Tags        string `json:"tags"`
		TakenAt     *time.Time `json:"taken_at"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if updateData.Title != "" {
		updates["title"] = updateData.Title
	}
	if updateData.Description != "" {
		updates["description"] = updateData.Description
	}
	if updateData.Category != "" {
		updates["category"] = updateData.Category
	}
	if updateData.Tags != "" {
		updates["tags"] = updateData.Tags
	}
	if updateData.TakenAt != nil {
		updates["taken_at"] = updateData.TakenAt
	}

	if err := c.DB.Model(&item).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&item, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": item})
}

// DeleteItem 删除相片
func (c *FamilyAlbumController) DeleteItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item models.FamilyAlbumItem
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相片不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 更新容器计数
	c.DB.Model(&models.FamilyAlbumContainer{}).Where("id = ?", item.ContainerID).
		Update("item_count", gorm.Expr("item_count - 1"))

	if err := c.DB.Delete(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// LikeItem 点赞相片
func (c *FamilyAlbumController) LikeItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item models.FamilyAlbumItem
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相片不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 检查是否已经点赞
	var existingLike models.FamilyAlbumLike
	var userID uint
	if err := ctx.ShouldBindBodyWith(&struct{ UserID uint `json:"user_id"` }{UserID: userID}, nil); err == nil {
		ctx.ShouldBindJSON(&struct{ UserID uint `json:"user_id"` }{})
	}
	if userID == 0 {
		userID = 1 // 默认用户
	}

	if err := c.DB.Where("photo_uuid = ? AND user_id = ?", item.UUID, userID).First(&existingLike).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "已经点赞过了"})
		return
	}

	// 创建点赞记录
	like := models.FamilyAlbumLike{
		PhotoUUID: item.UUID,
		UserID:   userID,
	}
	c.DB.Create(&like)

	// 更新点赞计数
	c.DB.Model(&item).Update("like_count", item.LikeCount+1)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// UnlikeItem 取消点赞
func (c *FamilyAlbumController) UnlikeItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item models.FamilyAlbumItem
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相片不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var userID uint = 1
	c.DB.Where("photo_uuid = ? AND user_id = ?", item.UUID, userID).Delete(&models.FamilyAlbumLike{})

	// 更新点赞计数
	if item.LikeCount > 0 {
		c.DB.Model(&item).Update("like_count", item.LikeCount-1)
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ShareItem 分享相片
func (c *FamilyAlbumController) ShareItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item models.FamilyAlbumItem
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相片不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 生成新的分享令牌
	item.ShareToken = generateToken()
	item.IsShared = true

	if err := c.DB.Model(&item).Updates(map[string]interface{}{
		"share_token": item.ShareToken,
		"is_shared":   true,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "分享失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"share_url":   "/share/album/" + item.ShareToken,
		"share_token": item.ShareToken,
	}})
}

// ListComments 获取评论列表
func (c *FamilyAlbumController) ListComments(ctx *gin.Context) {
	photoUUID := ctx.Query("photo_uuid")
	var comments []models.FamilyAlbumComment
	var total int64

	query := c.DB.Model(&models.FamilyAlbumComment{})
	if photoUUID != "" {
		query = query.Where("photo_uuid = ?", photoUUID)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      comments,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateComment 创建评论
func (c *FamilyAlbumController) CreateComment(ctx *gin.Context) {
	var comment models.FamilyAlbumComment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Create(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	// 更新评论计数
	var item models.FamilyAlbumItem
	if err := c.DB.Where("uuid = ?", comment.PhotoUUID).First(&item).Error; err == nil {
		c.DB.Model(&item).Update("comment_count", item.CommentCount+1)
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": comment})
}

// DeleteComment 删除评论
func (c *FamilyAlbumController) DeleteComment(ctx *gin.Context) {
	id := ctx.Param("id")
	var comment models.FamilyAlbumComment
	if err := c.DB.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	if err := c.DB.Delete(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	// 更新评论计数
	var item models.FamilyAlbumItem
	if err := c.DB.Where("uuid = ?", comment.PhotoUUID).First(&item).Error; err == nil {
		if item.CommentCount > 0 {
			c.DB.Model(&item).Update("comment_count", item.CommentCount-1)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func generateToken() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
