package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HospitalAppointment 宠物医院预约
type HospitalAppointment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AppointmentUUID string       `gorm:"type:varchar(64);uniqueIndex;not null" json:"appointment_uuid"`
	PetUUID       string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	PetName       string         `gorm:"type:varchar(32)" json:"pet_name"`
	HospitalName  string         `gorm:"type:varchar(128);not null" json:"hospital_name"`
	HospitalAddress string       `gorm:"type:varchar(256)" json:"hospital_address"`
	HospitalPhone string         `gorm:"type:varchar(32)" json:"hospital_phone"`
	Department    string         `gorm:"type:varchar(64)" json:"department"`            // 科室：内科/外科/皮肤科/眼科/牙科/急诊
	DoctorName    string         `gorm:"type:varchar(64)" json:"doctor_name"`
	AppointmentAt time.Time      `gorm:"type:timestamp;not null" json:"appointment_at"`  // 预约时间
	Reason        string         `gorm:"type:text" json:"reason"`                       // 就诊原因
	Status        string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/confirmed/completed/cancelled/no_show
	CancelReason  string         `gorm:"type:text" json:"cancel_reason"`
	CancelledAt   *time.Time     `json:"cancelled_at"`
	CheckInAt     *time.Time     `json:"check_in_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
	Diagnosis     string         `gorm:"type:text" json:"diagnosis"`                   // 诊断结果
	Prescription  string         `gorm:"type:text" json:"prescription"`               // 处方
	Cost          float64        `gorm:"type:decimal(10,2);default:0" json:"cost"`    // 费用
	Notes         string         `gorm:"type:text" json:"notes"`
	OwnerID       uint           `gorm:"index;not null" json:"owner_id"`
	HouseholdID   *uint          `gorm:"index" json:"household_id"`
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (HospitalAppointment) TableName() string {
	return "hospital_appointments"
}

// BeforeCreate 创建前自动生成 UUID
func (h *HospitalAppointment) BeforeCreate(tx *gorm.DB) error {
	if h.AppointmentUUID == "" {
		h.AppointmentUUID = uuid.New().String()
	}
	return nil
}
