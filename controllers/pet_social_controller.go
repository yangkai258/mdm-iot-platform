package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PetSocialController 宠物社交控制器
type PetSocialController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册宠物社交路由
func (c *PetSocialController) RegisterRoutes(api *gin.RouterGroup) {
	social := api.Group("/pet-social")
	{
		social.GET("/feed", c.GetFeed)
		social.POST("/posts", c.CreatePost)
		social.GET("/posts/:id", c.GetPost)
		social.POST("/posts/:id/like", c.LikePost)
		social.POST("/posts/:id/comment", c.CommentPost)
		social.POST("/follow/:pet_id", c.FollowPet)
		social.DELETE("/unfollow/:pet_id", c.UnfollowPet)
	}
}

// GetFeed 获取宠物动态
func (c *PetSocialController) GetFeed(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var posts []models.PetSocialPost
	var total int64

	query := c.DB.Model(&models.PetSocialPost{}).Where("is_public = ?", true)

	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      posts,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
		},
		"message": "success",
	})
}

// CreatePost 创建动态
func (c *PetSocialController) CreatePost(ctx *gin.Context) {
	var input struct {
		PetID     string `json:"pet_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Images    string `json:"images"`
		IsPublic  *bool  `json:"is_public"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	post := models.PetSocialPost{
		PetID:   input.PetID,
		Content: input.Content,
		Images:  input.Images,
	}

	if input.IsPublic != nil {
		post.IsPublic = *input.IsPublic
	} else {
		post.IsPublic = true
	}

	if err := c.DB.Create(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "data": post, "message": "success"})
}

// GetPost 获取动态详情
func (c *PetSocialController) GetPost(ctx *gin.Context) {
	id := ctx.Param("id")

	var post models.PetSocialPost
	if err := c.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "动态不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": post, "message": "success"})
}

// LikePost 点赞动态
func (c *PetSocialController) LikePost(ctx *gin.Context) {
	id := ctx.Param("id")

	var post models.PetSocialPost
	if err := c.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "动态不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	c.DB.Model(&post).Update("like_count", post.LikeCount+1)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "点赞成功"})
}

// CommentPost 评论动态
func (c *PetSocialController) CommentPost(ctx *gin.Context) {
	id := ctx.Param("id")

	var input struct {
		PetID    string `json:"pet_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var post models.PetSocialPost
	if err := c.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "动态不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	comment := models.PetSocialComment{
		PostID:  uint(post.ID),
		PetID:   input.PetID,
		Content: input.Content,
	}

	if err := c.DB.Create(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "评论失败"})
		return
	}

	c.DB.Model(&post).Update("comment_count", post.CommentCount+1)

	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "data": comment, "message": "评论成功"})
}

// FollowPet 关注宠物
func (c *PetSocialController) FollowPet(ctx *gin.Context) {
	followerID := ctx.Param("pet_id")

	var input struct {
		FollowingID string `json:"following_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查是否已关注
	var existing models.PetSocialFollow
	err := c.DB.Where("follower_id = ? AND following_id = ?", followerID, input.FollowingID).First(&existing).Error
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "已关注"})
		return
	}

	follow := models.PetSocialFollow{
		FollowerID:  followerID,
		FollowingID: input.FollowingID,
	}

	if err := c.DB.Create(&follow).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "关注失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "关注成功"})
}

// UnfollowPet 取消关注
func (c *PetSocialController) UnfollowPet(ctx *gin.Context) {
	unfollowerID := ctx.Param("pet_id")

	var input struct {
		FollowingID string `json:"following_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	result := c.DB.Where("follower_id = ? AND following_id = ?", unfollowerID, input.FollowingID).Delete(&models.PetSocialFollow{})
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "未关注"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "取消关注成功"})
}
