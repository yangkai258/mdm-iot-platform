package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PetShopProduct 宠物商店产品
type PetShopProduct struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ProductUUID  string         `gorm:"size:64;uniqueIndex;not null" json:"product_uuid"`
	Name         string         `gorm:"size:200;not null" json:"name"`
	Category     string         `gorm:"size:50;not null" json:"category"` // food, toy, accessory, medicine, grooming
	Brand        string         `gorm:"size:100" json:"brand"`
	PetSpecies   []string       `gorm:"type:text[]" json:"pet_species"` // dog, cat, bird, fish
	Price        float64        `gorm:"type:numeric(10,2);not null" json:"price"`
	OriginalPrice float64       `gorm:"type:numeric(10,2)" json:"original_price"`
	Stock        int64          `gorm:"default:0" json:"stock"`
	SoldCount    int64          `gorm:"default:0" json:"sold_count"`
	ImageURL     string         `gorm:"size:500" json:"image_url"`
	Images       []string       `gorm:"type:text[]" json:"images"`
	Description  string         `gorm:"type:text" json:"description"`
	Spec         string         `gorm:"type:text" json:"spec"` // 规格
	Weight       float64        `gorm:"type:numeric(8,2)" json:"weight"`
	WeightUnit   string         `gorm:"size:10;default:'g'" json:"weight_unit"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	IsFeatured   bool           `gorm:"default:false" json:"is_featured"`
	Tags         []string       `gorm:"type:text[]" json:"tags"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (p *PetShopProduct) BeforeCreate(tx *gorm.DB) error {
	if p.ProductUUID == "" {
		p.ProductUUID = uuid.New().String()
	}
	return nil
}

// PetShopOrder 宠物商店订单
type PetShopOrder struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OrderUUID    string         `gorm:"size:64;uniqueIndex;not null" json:"order_uuid"`
	OrderNo      string         `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	UserID       int64          `gorm:"not null;index" json:"user_id"`
	PetUUID      string         `gorm:"size:64;index" json:"pet_uuid"`
	TotalAmount  float64        `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	DiscountAmount float64      `gorm:"type:numeric(12,2);default:0" json:"discount_amount"`
	PayAmount    float64        `gorm:"type:numeric(12,2);not null" json:"pay_amount"`
	PayMethod    string         `gorm:"size:20" json:"pay_method"`
	PayStatus    string         `gorm:"size:20;default:'pending'" json:"pay_status"`
	OrderStatus  string         `gorm:"size:20;default:'pending';index" json:"order_status"`
	ReceiverName string         `gorm:"size:50" json:"receiver_name"`
	ReceiverPhone string        `gorm:"size:20" json:"receiver_phone"`
	ReceiverAddress string      `gorm:"size:500" json:"receiver_address"`
	Remark       string         `gorm:"type:text" json:"remark"`
	PaidAt       *time.Time     `json:"paid_at"`
	ShippedAt    *time.Time     `json:"shipped_at"`
	DeliveredAt  *time.Time     `json:"delivered_at"`
	CanceledAt   *time.Time     `json:"canceled_at"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Items        []PetShopOrderItem `gorm:"foreignKey:OrderUUID;references:OrderUUID" json:"items,omitempty"`
}

func (p *PetShopOrder) BeforeCreate(tx *gorm.DB) error {
	if p.OrderUUID == "" {
		p.OrderUUID = uuid.New().String()
	}
	if p.OrderNo == "" {
		p.OrderNo = "PSO" + time.Now().Format("20060102150405")
	}
	return nil
}

// PetShopOrderItem 订单商品
type PetShopOrderItem struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ItemUUID     string         `gorm:"size:64;uniqueIndex;not null" json:"item_uuid"`
	OrderUUID    string         `gorm:"size:64;not null;index" json:"order_uuid"`
	ProductUUID  string         `gorm:"size:64;not null" json:"product_uuid"`
	ProductName  string         `gorm:"size:200" json:"product_name"`
	Quantity     int64          `gorm:"default:1" json:"quantity"`
	Price        float64        `gorm:"type:numeric(10,2);not null" json:"price"`
	TotalAmount  float64        `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	CreatedAt    time.Time      `json:"created_at"`
}

func (p *PetShopOrderItem) BeforeCreate(tx *gorm.DB) error {
	if p.ItemUUID == "" {
		p.ItemUUID = uuid.New().String()
	}
	return nil
}

// PetShopCart 购物车
type PetShopCart struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       int64     `gorm:"not null;index" json:"user_id"`
	PetUUID      string    `gorm:"size:64;index" json:"pet_uuid"`
	ProductUUID  string    `gorm:"size:64;not null" json:"product_uuid"`
	Quantity     int64     `gorm:"default:1" json:"quantity"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PetShopFavorite 收藏
type PetShopFavorite struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       int64     `gorm:"not null;index" json:"user_id"`
	ProductUUID  string    `gorm:"size:64;not null" json:"product_uuid"`
	CreatedAt    time.Time `json:"created_at"`
}
