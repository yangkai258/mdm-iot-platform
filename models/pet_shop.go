package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductCategory 商品分类
type ProductCategory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(64);not null" json:"name"`
	ParentID    *uint          `gorm:"index" json:"parent_id"`
	IconURL     string         `gorm:"type:varchar(256)" json:"icon_url"`
	SortOrder   int            `gorm:"type:int;default:0" json:"sort_order"`
	Description string         `gorm:"type:text" json:"description"`
	IsActive    bool           `gorm:"type:boolean;default:true" json:"is_active"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TableName 指定表名
func (ProductCategory) TableName() string {
	return "product_categories"
}

// Product 商品
type Product struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ProductUUID   string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"product_uuid"`
	Name          string         `gorm:"type:varchar(128);not null" json:"name"`
	Brand         string         `gorm:"type:varchar(64)" json:"brand"`
	CategoryID    *uint          `gorm:"index" json:"category_id"`
	CategoryName  string         `gorm:"type:varchar(64)" json:"category_name"`
	Description   string         `gorm:"type:text" json:"description"`
	Price         float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64        `gorm:"type:decimal(10,2)" json:"original_price"`
	Stock         int            `gorm:"type:int;default:0" json:"stock"`
	Unit          string         `gorm:"type:varchar(16);default:'件'" json:"unit"` // 件/袋/盒/箱
	Images        string         `gorm:"type:text" json:"images"`                  // JSON 数组
	Specs         string         `gorm:"type:text" json:"specs"`                   // JSON 规格属性
	Tags          string         `gorm:"type:varchar(256)" json:"tags"`            // 逗号分隔
	Rating        float64        `gorm:"type:decimal(3,2);default:5.00" json:"rating"`
	SalesCount    int            `gorm:"type:int;default:0" json:"sales_count"`
	ViewCount     int            `gorm:"type:int;default:0" json:"view_count"`
	IsActive      bool           `gorm:"type:boolean;default:true" json:"is_active"`
	IsFeatured    bool           `gorm:"type:boolean;default:false" json:"is_featured"`
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "pet_products"
}

// BeforeCreate 创建前自动生成 UUID
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ProductUUID == "" {
		p.ProductUUID = uuid.New().String()
	}
	return nil
}

// OrderItem 订单项
type OrderItem struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	ProductUUID  string  `gorm:"type:varchar(64);index;not null" json:"product_uuid"`
	ProductName  string  `gorm:"type:varchar(128);not null" json:"product_name"`
	ProductImage string  `gorm:"type:varchar(256)" json:"product_image"`
	Price        float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity     int     `gorm:"type:int;not null" json:"quantity"`
	Specs        string  `gorm:"type:text" json:"specs"` // JSON 规格
}

// Order 订单
type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderUUID     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"order_uuid"`
	OrderNo       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"order_no"`
	PetUUID       string         `gorm:"type:varchar(64);index" json:"pet_uuid"` // 可选，关联宠物
	TotalAmount   float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	DiscountAmount float64       `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	PayAmount     float64        `gorm:"type:decimal(10,2);not null" json:"pay_amount"`
	PayMethod     string         `gorm:"type:varchar(32)" json:"pay_method"`      // wechat/alipay/card/cash
	PayStatus     string         `gorm:"type:varchar(20);default:'pending'" json:"pay_status"` // pending/paid/refunded/cancelled
	PayAt         *time.Time     `json:"pay_at"`
	Status        string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/confirmed/shipped/delivered/cancelled/completed
	ShippingAddress string       `gorm:"type:varchar(256)" json:"shipping_address"`
	ShippingName  string         `gorm:"type:varchar(64)" json:"shipping_name"`
	ShippingPhone string         `gorm:"type:varchar(32)" json:"shipping_phone"`
	ExpressNo     string         `gorm:"type:varchar(64)" json:"express_no"`
	Remarks       string         `gorm:"type:text" json:"remarks"`
	Items         string         `gorm:"type:text" json:"items"`                  // JSON 数组 OrderItem
	OwnerID       uint           `gorm:"index;not null" json:"owner_id"`
	HouseholdID   *uint          `gorm:"index" json:"household_id"`
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "pet_shop_orders"
}

// BeforeCreate 创建前自动生成 UUID 和订单号
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.OrderUUID == "" {
		o.OrderUUID = uuid.New().String()
	}
	if o.OrderNo == "" {
		o.OrderNo = generateOrderNo()
	}
	return nil
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	return time.Now().Format("20060102150405") + uuid.New().String()[:8]
}
