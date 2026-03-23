package models

import "time"

// Subscription 订阅
type Subscription struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	PlanName    string    `json:"plan_name"`
	PlanType    string    `json:"plan_type"` // free, basic, pro, enterprise
	Price       float64   `json:"price"`
	Duration    int       `json:"duration"`   // 天数
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"` // active, expired, cancelled
	AutoRenew   bool      `gorm:"default:true" json:"auto_renew"`
	LastRenewAt *time.Time `json:"last_renew_at"`
	RenewCount  int       `gorm:"default:0" json:"renew_count"`
	CreatedAt   time.Time `json:"created_at"`
}
