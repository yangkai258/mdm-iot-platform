package middleware

import (
	"bytes"
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
		// 只记录 /api 开头的写入请求
		if !strings.HasPrefix(c.Request.URL.Path, "/api") || 
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
			userIDUint = uint(userID.(int))
		}

		// 创建操作日志
		logEntry := models.SysOperationLog{
			UserID:    userIDUint,
			Username:  username.(string),
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

		// 异步写入日志
		go func() {
			if err := db.Create(&logEntry).Error; err != nil {
				log.Printf("Failed to create operation log: %v", err)
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
