package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=mdm_user password=mdm_password dbname=mdm_db port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// InitRedis 初始化 Redis 连接
func InitRedis() (*RedisClient, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	// 解析 redis://[user:password@]host:port[/db]
	// 默认值
	addr := "localhost:6379"
	password := ""
	db := 0

	// 去除 redis:// 前缀
	url := strings.TrimPrefix(redisURL, "redis://")

	// 检查是否有认证信息
	if idx := strings.Index(url, "@"); idx != -1 {
		authPart := url[:idx]
		url = url[idx+1:]
		if pIdx := strings.Index(authPart, ":"); pIdx != -1 {
			password = authPart[pIdx+1:]
		}
	}

	// 解析 host:port[/db]
	// 先按 / 分割，取第一部分作为 host:port
	hostPortPart := strings.Split(url, "/")[0]
	
	if strings.Contains(hostPortPart, ":") {
		// 有端口号
		parts := strings.SplitN(hostPortPart, ":", 2)
		addr = hostPortPart
		// 端口本身就是端口，不是DB
		_ = parts // addr已经是完整的 host:port
	} else {
		// 没有端口号，默认6379
		addr = hostPortPart + ":6379"
	}
	
	// 检查是否有 /db 后缀
	if strings.Contains(url, "/") {
		dbParts := strings.Split(url, "/")
		if len(dbParts) > 1 && dbParts[1] != "" {
			if dbNum, err := strconv.Atoi(dbParts[1]); err == nil {
				db = dbNum
			}
		}
	}

	return NewRedisClient(addr, password, db)
}

// RedisClient Redis 客户端封装
type RedisClient struct {
	client *redis.Client
}

// DeviceShadow 设备影子结构
type DeviceShadow struct {
	DeviceID      string     `json:"device_id"`
	IsOnline      bool       `json:"is_online"`
	BatteryLevel  int        `json:"battery_level"`
	CurrentMode   string     `json:"current_mode"`
	LastIP        string     `json:"last_ip"`
	LastHeartbeat *time.Time `json:"last_heartbeat"`
	DesiredConfig string     `json:"desired_config"`
}

// NewRedisClient 创建 Redis 客户端
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	return &RedisClient{client: client}, nil
}

// SetDeviceShadow 设置设备影子
func (r *RedisClient) SetDeviceShadow(deviceID string, shadow DeviceShadow, ttl time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("shadow:%s", deviceID)

	data, err := json.Marshal(shadow)
	if err != nil {
		return err
	}

	if ttl > 0 {
		return r.client.Set(ctx, key, data, ttl).Err()
	}
	return r.client.Set(ctx, key, data, 0).Err()
}

// GetDeviceShadow 获取设备影子
func (r *RedisClient) GetDeviceShadow(deviceID string) (*DeviceShadow, error) {
	ctx := context.Background()
	key := fmt.Sprintf("shadow:%s", deviceID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var shadow DeviceShadow
	if err := json.Unmarshal(data, &shadow); err != nil {
		return nil, err
	}

	return &shadow, nil
}

// GetAllShadowKeys 获取所有设备影子 Key
func (r *RedisClient) GetAllShadowKeys() ([]string, error) {
	ctx := context.Background()
	keys, err := r.client.Keys(ctx, "shadow:*").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// DelDeviceShadow 删除设备影子
func (r *RedisClient) DelDeviceShadow(deviceID string) error {
	ctx := context.Background()
	key := fmt.Sprintf("shadow:%s", deviceID)
	return r.client.Del(ctx, key).Err()
}

// Close 关闭连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}
