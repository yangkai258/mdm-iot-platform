package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OperationLog 操作日志中间件
func OperationLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除不需要记录操作的路径
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/v1/auth/") || path == "/health" {
			c.Next()
			return
		}
		
		// 只记录 /api 开头的写入请求
		if !strings.HasPrefix(path, "/api") ||
			(c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "DELETE") {
			c.Next()
			return
		}

		start := time.Now()

		// 读取请求体
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 记录响应
		blw := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer, statusCode: http.StatusOK}
		c.Writer = blw

		c.Next()

		// 获取用户信息
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		var userIDUint uint
		if userID != nil {
			switch v := userID.(type) {
			case int:
				userIDUint = uint(v)
			case uint:
				userIDUint = v
			case int64:
				userIDUint = uint(v)
			case float64:
				userIDUint = uint(v)
			}
		}
		var usernameStr string
		if username != nil {
			usernameStr = fmt.Sprintf("%v", username)
		}

		// 创建操作日志
		logEntry := models.SysOperationLog{
			UserID:    userIDUint,
			Username:  usernameStr,
			Module:    extractModule(c.Request.URL.Path),
			Operation: c.Request.Method + " " + c.Request.URL.Path,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			IP:        c.ClientIP(),
			Params:    string(bodyBytes),
			Result:    blw.body.String(),
			Status:    blw.statusCode,
			Duration:  int(time.Since(start).Milliseconds()),
		}

		if len(bodyBytes) > 2000 {
			logEntry.Params = string(bodyBytes[:2000]) + "...(truncated)"
		}

		// 解析路径获取资源类型和ID
		pathParts := strings.Split(c.Request.URL.Path, "/")
		resourceType := ""
		resourceID := uint(0)
		resourceName := ""
		if len(pathParts) >= 4 {
			resourceType = pathParts[3] // /api/v1/devices -> devices
			// 尝试从响应中提取 ID 和名称
			if blw.statusCode >= 200 && blw.statusCode < 300 {
				var resp map[string]interface{}
				if json.Unmarshal(blw.body.Bytes(), &resp) == nil {
					if data, ok := resp["data"].(map[string]interface{}); ok {
						if id, ok := data["id"].(float64); ok {
							resourceID = uint(id)
						}
						if name, ok := data["name"].(string); ok {
							resourceName = name
						} else if deviceID, ok := data["device_id"].(string); ok {
							resourceName = deviceID
						} else if memberName, ok := data["member_name"].(string); ok {
							resourceName = memberName
						} else if roleName, ok := data["role_name"].(string); ok {
							resourceName = roleName
						}
					}
				}
			}
		}

		// 确定操作类型
		action := "unknown"
		switch c.Request.Method {
		case "POST":
			action = "create"
		case "PUT", "PATCH":
			action = "update"
		case "DELETE":
			action = "delete"
		}

		// 映射资源类型名称（更友好的展示名）
		resourceTypeDisplay := resourceType
		switch resourceType {
		case "devices":
			resourceTypeDisplay = "device"
		case "members":
			resourceTypeDisplay = "member"
		case "roles":
			resourceTypeDisplay = "role"
		case "configs", "policy-configs":
			resourceTypeDisplay = "config"
		case "menus":
			resourceTypeDisplay = "menu"
		case "api-permissions":
			resourceTypeDisplay = "api_permission"
		case "permission-groups":
			resourceTypeDisplay = "permission_group"
		}

		// 创建活动日志（ActivityLog）
		activityLog := models.ActivityLog{
			UserID:       userIDUint,
			Username:     usernameStr,
			Action:       action,
			ResourceType: resourceTypeDisplay,
			ResourceID:   resourceID,
			ResourceName: resourceName,
			IP:           c.ClientIP(),
			UserAgent:    c.GetHeader("User-Agent"),
		}
		activityLog.SetDetails(map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"duration": time.Since(start).Milliseconds(),
		})

		// 异步写入日志
		go func() {
			if err := db.Create(&logEntry).Error; err != nil {
				log.Printf("Failed to create operation log: %v", err)
			}
			// 写入活动日志（仅记录关键的创建/更新/删除操作）
			if action != "unknown" && resourceType != "" {
				if err := db.Create(&activityLog).Error; err != nil {
					log.Printf("Failed to create activity log: %v", err)
				}
			}
		}()
	}
}

type bodyLogWriter struct {
	body       *bytes.Buffer
	statusCode int
	gin.ResponseWriter
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyLogWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func extractModule(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		return parts[2] // /api/v1/devices -> devices
	}
	return "unknown"
}
