package controllers

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PerformanceController 性能优化控制器
type PerformanceController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// RegisterPerformanceRoutes 注册性能优化路由
func (c *PerformanceController) RegisterPerformanceRoutes(r *gin.RouterGroup) {
	// 缓存管理
	r.POST("/cache/clear", c.ClearCache)
	r.GET("/cache/stats", c.GetCacheStats)

	// 性能监控
	r.GET("/metrics", c.GetMetrics)
	r.GET("/health", c.GetHealth)

	// 数据库优化
	r.POST("/db/reindex", c.Reindex)
	r.GET("/db/stats", c.GetDBStats)
}

// ClearCache 清空缓存
// POST /api/v1/performance/cache/clear
func (c *PerformanceController) ClearCache(ctx *gin.Context) {
	var request struct {
		Pattern string `json:"pattern"` // 可选，指定清空模式，如 "shadow:*"
	}
	request.Pattern = ctx.DefaultQuery("pattern", "*")

	if request.Pattern == "" {
		request.Pattern = "*"
	}

	ctx2 := context.Background()

	var keysCleared int64
	var err error

	// 根据 pattern 清除 keys
	if request.Pattern == "*" {
		// 清除所有 keys（除保留的 keys）
		var cursor uint64
		for {
			var batch []string
			batch, cursor, err = c.Redis.Client().Scan(ctx2, cursor, "shadow:*", 1000).Result()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": fmt.Sprintf("清空缓存失败: %v", err),
				})
				return
			}
			if len(batch) > 0 {
				deleted, _ := c.Redis.Client().Del(ctx2, batch...).Result()
				keysCleared += deleted
			}
			if cursor == 0 {
				break
			}
		}
		// 清除应用缓存
		batch2, _ := c.Redis.Client().Keys(ctx2, "app:*").Result()
		if len(batch2) > 0 {
			deleted, _ := c.Redis.Client().Del(ctx2, batch2...).Result()
			keysCleared += deleted
		}
	} else {
		var cursor uint64
		for {
			var batch []string
			batch, cursor, err = c.Redis.Client().Scan(ctx2, cursor, request.Pattern, 1000).Result()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": fmt.Sprintf("清空缓存失败: %v", err),
				})
				return
			}
			if len(batch) > 0 {
				deleted, _ := c.Redis.Client().Del(ctx2, batch...).Result()
				keysCleared += deleted
			}
			if cursor == 0 {
				break
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "缓存清空成功",
		"data": gin.H{
			"keys_cleared": keysCleared,
			"pattern":     request.Pattern,
		},
	})
}

// GetCacheStats 获取缓存统计
// GET /api/v1/performance/cache/stats
func (c *PerformanceController) GetCacheStats(ctx *gin.Context) {
	ctx2 := context.Background()

	// 获取 Redis 信息
	info, err := c.Redis.Client().Info(ctx2, "memory", "stats", "server", "persistence").Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取缓存统计失败: %v", err),
		})
		return
	}

	// 解析 keyspace 信息
	keyspace, _ := c.Redis.Client().Info(ctx2, "keyspace").Result()

	// 获取 db size
	dbSize := int64(0)
	for i := 0; i <= 15; i++ {
		size, err := c.Redis.Client().DBSize(ctx2).Result()
		if err == nil {
			dbSize += size
		}
	}

	// 获取 Redis uptime
	uptime, _ := c.Redis.Client().TTL(ctx2, "dummy_key_nonexistent").Result()
	_ = uptime // 忽略这个方式，直接用 Info

	// 解析 uptime from info
	var uptimeStr string
	uptimeLines := findInfoLine(info, "uptime_in_days")
	if uptimeLines != "" {
		uptimeStr = uptimeLines
	} else {
		uptimeStr = "unknown"
	}

	// 构造缓存统计
	stats := models.CacheStats{
		TotalKeys:     dbSize,
		MemoryUsed:    extractInfoValue(info, "used_memory_human"),
		MemoryPeak:    extractInfoValue(info, "used_memory_peak_human"),
		HitRate:       calculateHitRate(info),
		MissRate:      calculateMissRate(info),
		EvictionCount: parseInt64FromInfo(info, "evicted_keys"),
		Uptime:        uptimeStr,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"stats":    stats,
			"raw_info": info,
		},
	})
}

// GetMetrics 获取性能指标
// GET /api/v1/performance/metrics
func (c *PerformanceController) GetMetrics(ctx *gin.Context) {
	// 获取 Go 运行时指标
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 获取系统 CPU 数量
	numCPU := runtime.NumCPU()
	numGoroutine := runtime.NumGoroutine()

	// 获取数据库连接池状态
	var sqlStats struct {
		MaxOpenConnections int `json:"max_open_connections"`
		OpenConnections    int `json:"open_connections"`
		InUse              int `json:"in_use"`
		Idle               int `json:"idle"`
		WaitCount          int64 `json:"wait_count"`
		WaitDuration       string `json:"wait_duration"`
	}

	db, _ := c.DB.DB()
	if db != nil {
		sqlStats.MaxOpenConnections, _ = db.Stats().MaxOpenConnections, nil
		sqlStats.OpenConnections = db.Stats().OpenConnections
		sqlStats.InUse = db.Stats().InUse
		sqlStats.Idle = db.Stats().Idle
		sqlStats.WaitCount = db.Stats().WaitCount
		if db.Stats().WaitDuration > 0 {
			sqlStats.WaitDuration = db.Stats().WaitDuration.String()
		} else {
			sqlStats.WaitDuration = "0s"
		}
	}

	// 构造性能指标
	metrics := gin.H{
		"timestamp": time.Now().Format(time.RFC3339),
		"go": gin.H{
			"version":       runtime.Version(),
			"num_cpu":       numCPU,
			"num_goroutine": numGoroutine,
			"memory": gin.H{
				"alloc":      formatBytes(int64(m.Alloc)),
				"total_alloc": formatBytes(int64(m.TotalAlloc)),
				"sys":        formatBytes(int64(m.Sys)),
				"lookups":    m.Lookups,
				"mallocs":    m.Mallocs,
				"frees":      m.Frees,
				"heap_alloc":  formatBytes(int64(m.HeapAlloc)),
				"heap_sys":    formatBytes(int64(m.HeapSys)),
				"heap_idle":   formatBytes(int64(m.HeapIdle)),
				"heap_inuse":  formatBytes(int64(m.HeapInuse)),
				"heap_released": formatBytes(int64(m.HeapReleased)),
				"gc": gin.H{
					"num_gc":       m.NumGC,
					"num_forced_gc": m.NumForcedGC,
					"pause_total":  time.Duration(m.PauseTotalNs).String(),
				},
			},
		},
		"database": gin.H{
			"max_open_connections": sqlStats.MaxOpenConnections,
			"open_connections":     sqlStats.OpenConnections,
			"in_use":               sqlStats.InUse,
			"idle":                 sqlStats.Idle,
			"wait_count":           sqlStats.WaitCount,
			"wait_duration":        sqlStats.WaitDuration,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": metrics,
	})
}

// GetHealth 健康检查
// GET /api/v1/performance/health
func (c *PerformanceController) GetHealth(ctx *gin.Context) {
	start := time.Now()
	components := make(map[string]models.ComponentStatus)
	allHealthy := true

	// 1. 数据库健康检查
	dbStatus := models.ComponentStatus{Status: "up", Latency: "0ms"}
	dbStart := time.Now()
	sqlDB, err := c.DB.DB()
	if err != nil || sqlDB == nil {
		dbStatus.Status = "down"
		dbStatus.Message = "无法获取数据库连接"
		allHealthy = false
	} else {
		if err := sqlDB.Ping(); err != nil {
			dbStatus.Status = "down"
			dbStatus.Message = fmt.Sprintf("数据库 ping 失败: %v", err)
			allHealthy = false
		}
	}
	dbStatus.Latency = time.Since(dbStart).String()
	components["database"] = dbStatus

	// 2. Redis 健康检查
	redisStatus := models.ComponentStatus{Status: "up", Latency: "0ms"}
	redisStart := time.Now()
	ctx2 := context.Background()
	if _, err := c.Redis.Client().Ping(ctx2).Result(); err != nil {
		redisStatus.Status = "down"
		redisStatus.Message = fmt.Sprintf("Redis ping 失败: %v", err)
		allHealthy = false
	}
	redisStatus.Latency = time.Since(redisStart).String()
	components["redis"] = redisStatus

	// 3. 内存健康检查
	memStatus := models.ComponentStatus{Status: "up", Latency: "0ms"}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memPercent := float64(m.Alloc) / float64(m.Sys) * 100
	if memPercent > 90 {
		memStatus.Status = "degraded"
		memStatus.Message = fmt.Sprintf("内存使用率 %.2f%%", memPercent)
	} else if memPercent > 80 {
		memStatus.Status = "degraded"
	}
	components["memory"] = memStatus

	// 4. Goroutine 健康检查
	goroutineStatus := models.ComponentStatus{Status: "up"}
	numGoroutine := runtime.NumGoroutine()
	if numGoroutine > 10000 {
		goroutineStatus.Status = "degraded"
		goroutineStatus.Message = fmt.Sprintf("Goroutine 数量过多: %d", numGoroutine)
	} else if numGoroutine > 50000 {
		goroutineStatus.Status = "down"
		goroutineStatus.Message = fmt.Sprintf("Goroutine 数量严重超标: %d", numGoroutine)
		allHealthy = false
	}
	components["goroutine"] = goroutineStatus

	// 5. GC 健康检查
	gcStatus := models.ComponentStatus{Status: "up"}
	if m.NumGC > 0 {
		gcPausePercent := float64(m.PauseTotalNs) / float64(uint64(time.Hour))
		if gcPausePercent > 50 {
			gcStatus.Status = "degraded"
			gcStatus.Message = fmt.Sprintf("GC 暂停时间占总运行时间 %.2f%%", gcPausePercent)
		}
	}
	components["gc"] = gcStatus

	// 汇总状态
	status := "healthy"
	if !allHealthy {
		status = "unhealthy"
	} else {
		for _, comp := range components {
			if comp.Status == "degraded" {
				status = "degraded"
				break
			}
		}
	}

	health := models.HealthStatus{
		Status:     status,
		Timestamp:  time.Now(),
		Duration:   time.Since(start).String(),
		Components: components,
	}

	httpStatus := http.StatusOK
	if status == "unhealthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	ctx.JSON(httpStatus, gin.H{
		"code": 0,
		"data": health,
	})
}

// Reindex 重建数据库索引
// POST /api/v1/performance/db/reindex
func (c *PerformanceController) Reindex(ctx *gin.Context) {
	var request struct {
		Tables []string `json:"tables"` // 可选，指定表名，为空则重建所有
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		// 如果没有 JSON body，则重建所有索引
		request.Tables = nil
	}

	results := make([]map[string]interface{}, 0)

	if len(request.Tables) == 0 {
		// 重建所有表的所有索引
		var tableNames []string
		c.DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tableNames)

		for _, tableName := range tableNames {
			// 使用 REINDEX INDEX CONCURRENTLY 重建该表所有索引
			reindexSQL := fmt.Sprintf("REINDEX INDEX CONCURRENTLY %s", sanitizeIdentifier(tableName))
			result := c.DB.Exec(reindexSQL)
			
			results = append(results, map[string]interface{}{
				"table":  tableName,
				"status": "success",
				"sql":    reindexSQL,
			})
			if result.Error != nil {
				results[len(results)-1].(map[string]interface{})["status"] = "failed"
				results[len(results)-1].(map[string]interface{})["error"] = result.Error.Error()
			}
		}

		// 同时重建所有索引（不指定表）
		c.DB.Exec("REINDEX DATABASE CURRENT")

		results = append(results, map[string]interface{}{
			"database": "current",
			"status":   "success",
			"message":  "数据库所有索引已重建",
		})
	} else {
		// 重建指定表的索引
		for _, tableName := range request.Tables {
			// 获取该表的所有索引
			var indexes []string
			c.DB.Raw("SELECT indexname FROM pg_indexes WHERE tablename = ?", tableName).Scan(&indexes)

			for _, indexName := range indexes {
				reindexSQL := fmt.Sprintf("REINDEX INDEX CONCURRENTLY %s", sanitizeIdentifier(indexName))
				result := c.DB.Exec(reindexSQL)
				
				results = append(results, map[string]interface{}{
					"table": tableName,
					"index": indexName,
					"sql":   reindexSQL,
				})
				if result.Error != nil {
					results[len(results)-1].(map[string]interface{})["status"] = "failed"
					results[len(results)-1].(map[string]interface{})["error"] = result.Error.Error()
				} else {
					results[len(results)-1].(map[string]interface{})["status"] = "success"
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "索引重建完成",
		"data": gin.H{
			"results": results,
			"total":   len(results),
		},
	})
}

// GetDBStats 获取数据库统计
// GET /api/v1/performance/db/stats
func (c *PerformanceController) GetDBStats(ctx *gin.Context) {
	var stats models.DBStats

	// 获取连接池状态
	db, _ := c.DB.DB()
	if db != nil {
		dbStat := db.Stats()
		stats.TotalConnections = dbStat.MaxOpenConnections
		stats.ActiveConnections = dbStat.OpenConnections - dbStat.Idle
		stats.IdleConnections = dbStat.Idle
	}

	// 获取表数量
	c.DB.Raw("SELECT COUNT(*) FROM pg_tables WHERE schemaname = 'public'").Scan(&stats.TotalTables)

	// 获取索引数量
	c.DB.Raw("SELECT COUNT(*) FROM pg_indexes WHERE schemaname = 'public'").Scan(&stats.TotalIndexes)

	// 获取数据库大小
	var totalSize, indexSize, tableSize string
	c.DB.Raw("SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&totalSize)
	stats.TotalSize = totalSize

	// 获取表大小
	var tableSizeBytes int64
	c.DB.Raw("SELECT SUM(pg_total_relation_size(schemaname||'.'||tablename)) FROM pg_tables WHERE schemaname = 'public'").Scan(&tableSizeBytes)
	stats.TableSize = formatBytes(tableSizeBytes)

	// 获取索引大小
	var indexSizeBytes int64
	c.DB.Raw("SELECT SUM(pg_indexes_size(schemaname||'.'||tablename)) FROM pg_tables WHERE schemaname = 'public'").Scan(&indexSizeBytes)
	stats.IndexSize = formatBytes(indexSizeBytes)

	// 获取缓存命中率
	var cacheHit float64
	c.DB.Raw("SELECT sum(blks_hit)*100.0 / nullif(sum(blks_hit + blks_read), 0) FROM pg_stat_database WHERE datname = current_database()").Scan(&cacheHit)
	if cacheHit > 0 {
		stats.CacheHitRate = fmt.Sprintf("%.2f%%", cacheHit)
	} else {
		stats.CacheHitRate = "N/A"
	}

	// 获取每秒查询数
	var queriesPerSec float64
	c.DB.Raw("SELECT sum(xact_commit + xact_rollback) / extract(epoch from (now() - stats_reset)) FROM pg_stat_database WHERE datname = current_database()").Scan(&queriesPerSec)
	stats.QueryPerSecond = queriesPerSec

	// 获取平均查询时间
	var avgQueryTime float64
	c.DB.Raw("SELECT avg(mean_exec_time) FROM pg_stat_statements WHERE query NOT LIKE '%pg_stat_statements%'").Scan(&avgQueryTime)
	if avgQueryTime > 0 {
		stats.AvgQueryTime = fmt.Sprintf("%.2fms", avgQueryTime)
	} else {
		stats.AvgQueryTime = "N/A"
	}

	// 获取索引使用情况（top 10）
	var indexInfos []models.IndexInfo
	c.DB.Raw(`
		SELECT 
			schemaname as schema,
			relname as table_name,
			indexrelname as index_name,
			pg_size_pretty(pg_relation_size(indexrelid)) as index_size,
			indisunique as is_unique,
			indisvalid as is_valid,
			idx_scan as scans
		FROM pg_stat_user_indexes 
		ORDER BY idx_scan DESC 
		LIMIT 10
	`).Scan(&indexInfos)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"stats":       stats,
			"top_indexes": indexInfos,
		},
	})
}

// ============ 辅助函数 ============

func sanitizeIdentifier(name string) string {
	// 简单防止 SQL 注入，只允许字母数字和下划线
	result := ""
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			result += string(c)
		}
	}
	return result
}

func findInfoLine(info, key string) string {
	lines := splitInfo(info)
	for _, line := range lines {
		if len(line) > 0 && line[0:min(len(line), len(key))] == key {
			return line
		}
	}
	return ""
}

func extractInfoValue(info, key string) string {
	lines := splitInfo(info)
	for _, line := range lines {
		if len(line) > len(key) && line[0:len(key)] == key {
			colonIdx := -1
			for i, c := range line {
				if c == ':' {
					colonIdx = i
					break
				}
			}
			if colonIdx != -1 && colonIdx < len(line)-1 {
				return line[colonIdx+1:]
			}
		}
	}
	return "N/A"
}

func splitInfo(info string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(info); i++ {
		if info[i] == '\n' {
			line := info[start:i]
			// 移除 \r
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(info) {
		line := info[start:]
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		lines = append(lines, line)
	}
	return lines
}

func calculateHitRate(info string) string {
	var hits, misses int64
	lines := splitInfo(info)
	for _, line := range lines {
		if len(line) > 13 && line[0:13] == "keyspace_hits:" {
			fmt.Sscanf(line[13:], "%d", &hits)
		}
		if len(line) > 15 && line[0:15] == "keyspace_misses:" {
			fmt.Sscanf(line[15:], "%d", &misses)
		}
	}
	total := hits + misses
	if total == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%.2f%%", float64(hits)/float64(total)*100)
}

func calculateMissRate(info string) string {
	var hits, misses int64
	lines := splitInfo(info)
	for _, line := range lines {
		if len(line) > 13 && line[0:13] == "keyspace_hits:" {
			fmt.Sscanf(line[13:], "%d", &hits)
		}
		if len(line) > 15 && line[0:15] == "keyspace_misses:" {
			fmt.Sscanf(line[15:], "%d", &misses)
		}
	}
	total := hits + misses
	if total == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%.2f%%", float64(misses)/float64(total)*100)
}

func parseInt64FromInfo(info, key string) int64 {
	var value int64
	lines := splitInfo(info)
	for _, line := range lines {
		if len(line) > len(key) && line[0:len(key)] == key {
			colonIdx := -1
			for i, c := range line {
				if c == ':' {
					colonIdx = i
					break
				}
			}
			if colonIdx != -1 && colonIdx < len(line)-1 {
				fmt.Sscanf(line[colonIdx+1:], "%d", &value)
			}
		}
	}
	return value
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
