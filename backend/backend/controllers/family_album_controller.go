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
func (ctrl *FamilyAlbumController) RegisterRoutes(api *gin.RouterGroup) {
	albums := api.Group("/family/albums")
	{
		albums.GET("", ctrl.ListAlbums)
		albums.GET("/:id", ctrl.GetAlbum)
		albums.POST("", ctrl.CreateAlbum)
		albums.PUT("/:id", ctrl.UpdateAlbum)
		albums.DELETE("/:id", ctrl.DeleteAlbum)
		albums.POST("/:id/photos", ctrl.UploadPhoto)
		albums.DELETE("/:id/photos/:photo_id", ctrl.DeletePhoto)
		albums.GET("/:id/photos", ctrl.ListPhotos)
	}
}

// ListAlbums GET /api/v1/family/albums - 相册列表
func (ctrl *FamilyAlbumController) ListAlbums(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)
	familyID := c.Query("family_id")
	keyword := c.Query("keyword")
	privacy := c.Query("privacy")

	query := ctrl.DB.Model(&models.FamilyAlbumContainer{})

	if familyID != "" {
		query = query.Where("family_id = ?", familyID)
	}
	if keyword != "" {
		query = query.Where("album_name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if privacy != "" {
		query = query.Where("privacy = ?", privacy)
	}

	var total int64
	query.Count(&total)

	var albums []models.FamilyAlbumContainer
	if err := query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&albums).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":       albums,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetAlbum GET /api/v1/family/albums/:id - 相册详情
func (ctrl *FamilyAlbumController) GetAlbum(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的相册ID"})
		return
	}

	var album models.FamilyAlbumContainer
	if err := ctrl.DB.First(&album, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相册不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": album})
}

// CreateAlbum POST /api/v1/family/albums - 创建相册
func (ctrl *FamilyAlbumController) CreateAlbum(c *gin.Context) {
	var req CreateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	album := models.FamilyAlbumContainer{
		FamilyID:    req.FamilyID,
		AlbumName:   req.AlbumName,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		CreatedBy:   req.CreatedBy,
		Privacy:     req.Privacy,
		Tags:        req.Tags,
	}

	if err := ctrl.DB.Create(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": album})
}

// UpdateAlbum PUT /api/v1/family/albums/:id - 更新相册
func (ctrl *FamilyAlbumController) UpdateAlbum(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的相册ID"})
		return
	}

	var album models.FamilyAlbumContainer
	if err := ctrl.DB.First(&album, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相册不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req UpdateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := make(map[string]interface{})
	if req.AlbumName != "" {
		updates["album_name"] = req.AlbumName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.CoverURL != "" {
		updates["cover_url"] = req.CoverURL
	}
	if req.Privacy != "" {
		updates["privacy"] = req.Privacy
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	if err := ctrl.DB.Model(&album).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctrl.DB.First(&album, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": album})
}

// DeleteAlbum DELETE /api/v1/family/albums/:id - 删除相册
func (ctrl *FamilyAlbumController) DeleteAlbum(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的相册ID"})
		return
	}

	// 删除相册下的所有照片
	ctrl.DB.Where("album_id = ?", id).Delete(&models.FamilyAlbumItem{})

	if err := ctrl.DB.Delete(&models.FamilyAlbumContainer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// UploadPhoto POST /api/v1/family/albums/:id/photos - 上传照片
func (ctrl *FamilyAlbumController) UploadPhoto(c *gin.Context) {
	albumID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的相册ID"})
		return
	}

	// 检查相册是否存在
	var album models.FamilyAlbumContainer
	if err := ctrl.DB.First(&album, albumID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "相册不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req UploadPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	photo := models.FamilyAlbumItem{
		AlbumID:      uint(albumID),
		UploaderID:   req.UploaderID,
		PhotoURL:     req.PhotoURL,
		ThumbnailURL: req.ThumbnailURL,
		Caption:      req.Caption,
		Width:        req.Width,
		Height:       req.Height,
		FileSize:     req.FileSize,
		LocationLat:  req.LocationLat,
		LocationLng:  req.LocationLng,
		LocationName: req.LocationName,
		Tags:         req.Tags,
		IsFeatured:   req.IsFeatured,
	}

	if req.TakenAt != "" {
		if t, err := time.Parse(time.RFC3339, req.TakenAt); err == nil {
			photo.TakenAt = &t
		}
	}

	if err := ctrl.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "上传失败"})
		return
	}

	// 更新相册照片数量
	ctrl.DB.Model(&album).Update("photo_count", gorm.Expr("photo_count + 1"))

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": photo})
}

// DeletePhoto DELETE /api/v1/family/albums/:id/photos/:photo_id - 删除照片
func (ctrl *FamilyAlbumController) DeletePhoto(c *gin.Context) {
	albumID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的相册ID"})
		return
	}

	photoID, err := strconv.ParseUint(c.Param("photo_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的照片ID"})
		return
	}

	var photo models.FamilyAlbumItem
	if err := ctrl.DB.Where("id = ? AND album_id = ?", photoID, albumID).First(&photo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "照片不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if err := ctrl.DB.Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	// 更新相册照片数量
	ctrl.DB.Model(&models.FamilyAlbumContainer{}).Where("id = ?", albumID).Update("photo_count", gorm.Expr("photo_count - 1"))

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ListPhotos GET /api/v1/family/albums/:id/photos - 获取相册照片列表
func (ctrl *FamilyAlbumController) ListPhotos(c *gin.Context) {
	albumID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的相册ID"})
		return
	}

	page := defaultPage(c)
	pageSize := defaultPageSize(c)
	keyword := c.Query("keyword")

	query := ctrl.DB.Model(&models.FamilyAlbumItem{}).Where("album_id = ?", albumID)

	if keyword != "" {
		query = query.Where("caption LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var photos []models.FamilyAlbumItem
	if err := query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":       photos,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ========== 请求结构体 ==========

type CreateAlbumRequest struct {
	FamilyID    uint     `json:"family_id"`
	AlbumName   string   `json:"album_name"`
	Description string   `json:"description"`
	CoverURL    string   `json:"cover_url"`
	CreatedBy   uint     `json:"created_by"`
	Privacy     string   `json:"privacy"`
	Tags        []string `json:"tags"`
}

type UpdateAlbumRequest struct {
	AlbumName   string   `json:"album_name"`
	Description string   `json:"description"`
	CoverURL    string   `json:"cover_url"`
	Privacy     string   `json:"privacy"`
	Tags        []string `json:"tags"`
}

type UploadPhotoRequest struct {
	UploaderID    uint    `json:"uploader_id"`
	PhotoURL     string  `json:"photo_url"`
	ThumbnailURL string  `json:"thumbnail_url"`
	Caption      string  `json:"caption"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	FileSize     int64   `json:"file_size"`
	TakenAt      string  `json:"taken_at"`
	LocationLat  float64 `json:"location_lat"`
	LocationLng  float64 `json:"location_lng"`
	LocationName string  `json:"location_name"`
	Tags         []string `json:"tags"`
	IsFeatured   bool     `json:"is_featured"`
}
