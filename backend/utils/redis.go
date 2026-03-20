package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
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

	// 使用 net/url 正确解析 redis:// URL
	parsedURL, err := url.Parse(redisURL)
	if err != nil {
		return nil, fmt.Errorf("Redis URL 解析失败: %w", err)
	}

	// 默认值
	addr := "localhost:6379"
	password := ""
	db := 0

	// 获取 host:port
	if parsedURL.Host != "" {
		addr = parsedURL.Host
	} else if parsedURL.Path != "" {
		// 没有 Host 时，Path 包含路径
		addr = parsedURL.Path + ":6379"
	}

	// 获取密码
	if parsedURL.User != nil {
		password, _ = parsedURL.User.Password()
	}

	// 获取数据库索引 (path 以 / 开头，需要去掉)
	path := parsedURL.Path
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	if path != "" {
		if dbNum, err := strconv.Atoi(path); err == nil {
			db = dbNum
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
	// 越狱/ROOT检测
	IsJailbroken bool    `json:"is_jailbroken"`
	RootStatus   string  `json:"root_status"` // normal, rooted, jailbroken
	// 地理位置
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
