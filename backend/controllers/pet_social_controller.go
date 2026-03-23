package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PetSocialController struct {
	DB *gorm.DB
}

func NewPetSocialController(db *gorm.DB) *PetSocialController {
	return &PetSocialController{DB: db}
}

func (ctrl *PetSocialController) RegisterRoutes(rg *gin.RouterGroup) {
	posts := rg.Group("/pet-social")
	{
		posts.GET("/posts", ctrl.GetPosts)
		posts.POST("/posts", ctrl.CreatePost)
		posts.GET("/posts/:id", ctrl.GetPost)
		posts.DELETE("/posts/:id", ctrl.DeletePost)
		posts.POST("/posts/:id/like", ctrl.LikePost)
		posts.DELETE("/posts/:id/like", ctrl.UnlikePost)
		posts.GET("/posts/:id/comments", ctrl.GetComments)
		posts.POST("/posts/:id/comments", ctrl.CreateComment)
		posts.DELETE("/comments/:id", ctrl.DeleteComment)

		posts.GET("/following", ctrl.GetFollowing)
		posts.POST("/follow", ctrl.Follow)
		posts.DELETE("/follow/:id", ctrl.Unfollow)
		posts.GET("/followers", ctrl.GetFollowers)

		posts.GET("/playdates", ctrl.GetPlaydates)
		posts.POST("/playdates", ctrl.CreatePlaydate)
		posts.GET("/playdates/:id", ctrl.GetPlaydate)
		posts.PUT("/playdates/:id", ctrl.UpdatePlaydate)
		posts.POST("/playdates/:id/join", ctrl.JoinPlaydate)
		posts.POST("/playdates/:id/cancel", ctrl.CancelPlaydate)
	}
}

func (ctrl *PetSocialController) GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	petID := c.Query("pet_id")
	userID := c.Query("user_id")

	var posts []models.Post
	query := ctrl.DB.Model(&models.Post{}).Where("is_public = ?", true)

	if petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if userID != "" {
		query = query.Where("author_id = ?", userID)
	}

	var total int64
	query.Count(&total)

	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *PetSocialController) CreatePost(c *gin.Context) {
	var req struct {
		PetID     uint     `json:"pet_id" binding:"required"`
		AuthorID  uint     `json:"author_id" binding:"required"`
		Content   string   `json:"content" binding:"required"`
		MediaURLs []string `json:"media_urls"`
		PostType  string   `json:"post_type"`
		Location  string   `json:"location"`
		IsPublic  bool     `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mediaURLsJSON, _ := json.Marshal(req.MediaURLs)

	post := models.Post{
		PetID:     req.PetID,
		AuthorID:  req.AuthorID,
		Content:   req.Content,
		MediaURLs: string(mediaURLsJSON),
		PostType:  req.PostType,
		Location:  req.Location,
		IsPublic:  req.IsPublic,
	}

	ctrl.DB.Create(&post)
	c.JSON(http.StatusOK, post)
}

func (ctrl *PetSocialController) GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := ctrl.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (ctrl *PetSocialController) DeletePost(c *gin.Context) {
	id := c.Param("id")

	ctrl.DB.Delete(&models.Post{}, id)
	ctrl.DB.Where("post_id = ?", id).Delete(&models.PostComment{})
	ctrl.DB.Where("post_id = ?", id).Delete(&models.PostLike{})

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (ctrl *PetSocialController) LikePost(c *gin.Context) {
	postID := c.Param("id")
	userID := c.Query("user_id")

	var like models.PostLike
	result := ctrl.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&like)
	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{"message": "already liked"})
		return
	}

	like = models.PostLike{
		PostID: toUint(postID),
		UserID: toUint(userID),
	}
	ctrl.DB.Create(&like)
	ctrl.DB.Model(&models.Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("like_count + 1"))

	c.JSON(http.StatusOK, gin.H{"message": "liked"})
}

func (ctrl *PetSocialController) UnlikePost(c *gin.Context) {
	postID := c.Param("id")
	userID := c.Query("user_id")

	ctrl.DB.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.PostLike{})
	ctrl.DB.Model(&models.Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)"))

	c.JSON(http.StatusOK, gin.H{"message": "unliked"})
}

func (ctrl *PetSocialController) GetComments(c *gin.Context) {
	postID := c.Param("id")
	var comments []models.PostComment

	ctrl.DB.Where("post_id = ?", postID).Order("created_at DESC").Find(&comments)
	c.JSON(http.StatusOK, comments)
}

func (ctrl *PetSocialController) CreateComment(c *gin.Context) {
	var req struct {
		UserID  uint   `json:"user_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postID := c.Param("id")
	comment := models.PostComment{
		PostID:  toUint(postID),
		UserID:  req.UserID,
		Content: req.Content,
	}

	ctrl.DB.Create(&comment)
	ctrl.DB.Model(&models.Post{}).Where("id = ?", postID).UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))

	c.JSON(http.StatusOK, comment)
}

func (ctrl *PetSocialController) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	var comment models.PostComment

	if err := ctrl.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	ctrl.DB.Delete(&comment)
	ctrl.DB.Model(&models.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)"))

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (ctrl *PetSocialController) GetFollowing(c *gin.Context) {
	userID := c.Query("user_id")
	var follows []models.Follow

	query := ctrl.DB.Model(&models.Follow{})
	if userID != "" {
		query = query.Where("follower_id = ?", userID)
	}

	query.Find(&follows)
	c.JSON(http.StatusOK, follows)
}

func (ctrl *PetSocialController) Follow(c *gin.Context) {
	var req struct {
		FollowerID uint   `json:"follower_id" binding:"required"`
		FolloweeID uint   `json:"followee_id" binding:"required"`
		FollowType string `json:"follow_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	follow := models.Follow{
		FollowerID: req.FollowerID,
		FolloweeID: req.FolloweeID,
		FollowType: req.FollowType,
	}

	ctrl.DB.Create(&follow)
	c.JSON(http.StatusOK, follow)
}

func (ctrl *PetSocialController) Unfollow(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.Follow{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}

func (ctrl *PetSocialController) GetFollowers(c *gin.Context) {
	userID := c.Query("user_id")
	var follows []models.Follow

	query := ctrl.DB.Model(&models.Follow{})
	if userID != "" {
		query = query.Where("followee_id = ?", userID)
	}

	query.Find(&follows)
	c.JSON(http.StatusOK, follows)
}

func (ctrl *PetSocialController) GetPlaydates(c *gin.Context) {
	var playdates []models.PetPlaydate
	ctrl.DB.Where("status != ?", "cancelled").Order("start_time ASC").Find(&playdates)
	c.JSON(http.StatusOK, playdates)
}

func (ctrl *PetSocialController) CreatePlaydate(c *gin.Context) {
	var req struct {
		OrganizerID uint      `json:"organizer_id" binding:"required"`
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		Location    string    `json:"location"`
		PetIDs      []uint    `json:"pet_ids"`
		StartTime   time.Time `json:"start_time" binding:"required"`
		EndTime     time.Time `json:"end_time" binding:"required"`
		MaxPets     int       `json:"max_pets"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	petIDsJSON, _ := json.Marshal(req.PetIDs)

	playdate := models.PetPlaydate{
		OrganizerID: req.OrganizerID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		PetIDs:      string(petIDsJSON),
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		MaxPets:     req.MaxPets,
		Status:      "pending",
	}

	ctrl.DB.Create(&playdate)
	c.JSON(http.StatusOK, playdate)
}

func (ctrl *PetSocialController) GetPlaydate(c *gin.Context) {
	id := c.Param("id")
	var playdate models.PetPlaydate

	if err := ctrl.DB.First(&playdate, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playdate not found"})
		return
	}

	c.JSON(http.StatusOK, playdate)
}

func (ctrl *PetSocialController) UpdatePlaydate(c *gin.Context) {
	id := c.Param("id")
	var playdate models.PetPlaydate

	if err := ctrl.DB.First(&playdate, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playdate not found"})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Location    string `json:"location"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&playdate).Updates(req)
	c.JSON(http.StatusOK, playdate)
}

func (ctrl *PetSocialController) JoinPlaydate(c *gin.Context) {
	id := c.Param("id")
	var playdate models.PetPlaydate

	if err := ctrl.DB.First(&playdate, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playdate not found"})
		return
	}

	ctrl.DB.Model(&playdate).Update("status", "confirmed")
	c.JSON(http.StatusOK, gin.H{"message": "joined", "playdate": playdate})
}

func (ctrl *PetSocialController) CancelPlaydate(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Model(&models.PetPlaydate{}).Where("id = ?", id).Update("status", "cancelled")
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

func toUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}
