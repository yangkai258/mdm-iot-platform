package models

import (
	"time"

	"gorm.io/gorm"
)

// FeatureGroup 功能分组
type FeatureGroup struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	GroupName  string         `gorm:"type:varchar(100);not null" json:"group_name"` // 分组名称
	GroupCode  string         `gorm:"type:varchar(50);uniqueIndex" json:"group_code"` // 分组编码
	Icon       string         `gorm:"type:varchar(100)" json:"icon"`                  // 分组图标
	Color      string         `gorm:"type:varchar(20)" json:"color"`                  // 分组颜色
	Sort       int            `gorm:"default:0" json:"sort"`                           // 排序
	Description string        `gorm:"type:text" json:"description"`                   // 描述
	Status     int            `gorm:"default:1" json:"status"`                        // 1=启用 0=禁用
	TenantID   string         `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Features   []FeatureItem  `gorm:"foreignKey:GroupID" json:"features,omitempty"`
}

func (FeatureGroup) TableName() string {
	return "feature_groups"
}

// FeatureItem 功能项
type FeatureItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	GroupID     *uint          `gorm:"index" json:"group_id"`                     // 分组ID
	FeatureName string         `gorm:"type:varchar(100);not null" json:"feature_name"` // 功能名称
	FeatureCode string         `gorm:"type:varchar(50);uniqueIndex" json:"feature_code"` // 功能编码
	Icon        string         `gorm:"type:varchar(100)" json:"icon"`             // 功能图标
	RoutePath   string         `gorm:"type:varchar(255)" json:"route_path"`       // 路由路径
	Component   string         `gorm:"type:varchar(255)" json:"component"`        // 组件路径
	ApiPaths    string         `gorm:"type:text" json:"api_paths"`                // 关联的API路径(JSON数组)
	Permission  string         `gorm:"type:varchar(100)" json:"permission"`     // 权限编码
	Sort        int            `gorm:"default:0" json:"sort"`                      // 排序
	Status      int            `gorm:"default:1" json:"status"`                   // 1=启用 0=禁用
	IsDefault   int            `gorm:"default:0" json:"is_default"`               // 1=默认选中
	Badge       string         `gorm:"type:varchar(50)" json:"badge"`              // 徽章(如"新","Beta")
	Description string         `gorm:"type:text" json:"description"`               // 描述
	TenantID    string         `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FeatureItem) TableName() string {
	return "feature_items"
}

// 请求/响应结构

// FeatureGroupRequest 创建/更新分组
type FeatureGroupRequest struct {
	GroupName   string `json:"group_name" binding:"required"`
	GroupCode   string `json:"group_code"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Sort        int    `json:"sort"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

// FeatureItemRequest 创建/更新功能项
type FeatureItemRequest struct {
	GroupID     *uint  `json:"group_id"`
	FeatureName string `json:"feature_name" binding:"required"`
	FeatureCode string `json:"feature_code"`
	Icon        string `json:"icon"`
	RoutePath   string `json:"route_path"`
	Component   string `json:"component"`
	ApiPaths    string `json:"api_paths"`
	Permission  string `json:"permission"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	IsDefault   int    `json:"is_default"`
	Badge       string `json:"badge"`
	Description string `json:"description"`
}

// FeatureSortRequest 排序请求
type FeatureSortRequest struct {
	Items []SortItem `json:"items" binding:"required"`
}

type SortItem struct {
	ID    uint   `json:"id"`
	Sort  int    `json:"sort"`
	GroupID *uint `json:"group_id"` // 跨分组拖拽时需要
}

// FeatureGroupTree 树形结构响应
type FeatureGroupTree struct {
	FeatureGroup
	Children []FeatureItem `json:"children"`
}
