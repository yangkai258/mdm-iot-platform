package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Firmware 固件
type Firmware struct {
	ID        uint   `json:"id"`
	Version   string `json:"version"`
	FileURL   string `json:"file_url"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"created_at"`
}

// MiniClawController MiniClaw控制器
type MiniClawController struct{}

// RegisterRoutes 注册路由
func (c *MiniClawController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/miniclaw/firmwares", c.ListFirmwares)
	r.POST("/miniclaw/firmwares", c.UploadFirmware)
	r.GET("/miniclaw/firmwares/:id", c.GetFirmware)
	r.DELETE("/miniclaw/firmwares/:id", c.DeleteFirmware)
	r.GET("/miniclaw/devices/:device_id/firmware", c.GetDeviceFirmware)
	r.PUT("/miniclaw/devices/:device_id/firmware", c.UpdateDeviceFirmware)
}

func (c *MiniClawController) ListFirmwares(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": []Firmware{}})
}

func (c *MiniClawController) UploadFirmware(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": Firmware{}})
}

func (c *MiniClawController) GetFirmware(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": Firmware{}})
}

func (c *MiniClawController) DeleteFirmware(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0})
}

func (c *MiniClawController) GetDeviceFirmware(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
}

func (c *MiniClawController) UpdateDeviceFirmware(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0})
}
