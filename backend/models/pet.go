package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// VaccinationRecord 疫苗记录
type VaccinationRecord struct {
	VaccineName   string    `json:"vaccine_name"`   // 疫苗名称
	VaccinationAt time.Time `json:"vaccination_at"` // 接种日期
	NextDueDate   time.Time `json:"next_due_date"`  // 下次接种日期
	VetClinic     string    `json:"vet_clinic"`      // 接种医院
	BatchNo       string    `json:"batch_no"`       // 疫苗批号
	Notes         string    `json:"notes"`          // 备注
}

// HealthRecord 健康档案
type HealthRecord struct {
	RecordType string    `json:"record_type"` // checkup/vaccination/medication/surgery
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	Hospital   string    `json:"hospital"`
	Doctor     string    `json:"doctor"`
	Diagnosis  string    `json:"diagnosis"`  // 诊断结果
	Treatment  string    `json:"treatment"` // 治疗方案
	Cost       float64   `json:"cost"`      // 费用
	Notes      string    `json:"notes"`
	Attachments []string `json:"attachments"` // 附件图片URLs
}

// Pet 宠物档案
type Pet struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	PetUUID            string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"pet_uuid"`
	PetName            string         `gorm:"type:varchar(32);not null" json:"pet_name"`
	Species            string         `gorm:"type:varchar(32)" json:"species"`            // dog/cat/bird/rabbit/other
	Breed              string         `gorm:"type:varchar(64)" json:"breed"`             // 品种
	Gender             string         `gorm:"type:varchar(16)" json:"gender"`           // male/female/unknown
	BirthDate          *time.Time     `json:"birth_date"`                             // 出生日期
	Weight             float64        `gorm:"type:decimal(5,2)" json:"weight"`          // kg
	Color              string         `gorm:"type:varchar(32)" json:"color"`           // 毛色
	Microchip          string         `gorm:"type:varchar(64)" json:"microchip"`       // 芯片号
	AvatarURL          string         `gorm:"type:varchar(512)" json:"avatar_url"`     // 头像URL
	OwnerID            uint           `gorm:"index;not null" json:"owner_id"`         // 主人用户ID
	HouseholdID        *uint          `gorm:"index" json:"household_id"`              // 所属家庭ID
	Status             string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/lost/adopted/deceased
	Description        string         `gorm:"type:text" json:"description"`            // 宠物简介
	VaccinationRecords pq.StringArray `gorm:"type:text[]" json:"vaccination_records"`  // 疫苗记录JSON数组
	HealthRecords      pq.StringArray `gorm:"type:text[]" json:"health_records"`       // 健康档案JSON数组
	TenantID           string         `gorm:"type:uuid;index" json:"tenant_id"`       // 租户ID
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Pet) TableName() string {
	return "pets"
}

// BeforeCreate 创建前自动生成 UUID
func (p *Pet) BeforeCreate(tx *gorm.DB) error {
	if p.PetUUID == "" {
		p.PetUUID = uuid.New().String()
	}
	return nil
}

// PetDeviceBinding 宠物-设备绑定关系
type PetDeviceBinding struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PetUUID     string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	DeviceID    string         `gorm:"type:varchar(64);index;not null" json:"device_id"`
	BindingType string         `gorm:"type:varchar(20);default:'primary'" json:"binding_type"` // primary/secondary
	IsActive    bool           `gorm:"type:boolean;default:true" json:"is_active"`
	BoundAt     time.Time      `gorm:"type:timestamp;default:now()" json:"bound_at"`
	UnboundAt   *time.Time     `json:"unbound_at"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TableName 指定表名
func (PetDeviceBinding) TableName() string {
	return "pet_device_bindings"
}

// HouseholdMember 家庭成员
type HouseholdMember struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	HouseholdID  uint           `gorm:"index;not null" json:"household_id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	Role         string         `gorm:"type:varchar(20);default:'member'" json:"role"` // owner/member/viewer
	Status       string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/inactive
	TenantID     string         `gorm:"type:uuid;index" json:"tenant_id"`
	InviteCode   string         `gorm:"type:varchar(20);uniqueIndex" json:"invite_code"` // 邀请码
	InvitedEmail string         `gorm:"type:varchar(255)" json:"invited_email"`
	InvitedBy    *uint          `json:"invited_by"` // 邀请人
	JoinedAt     time.Time      `gorm:"type:timestamp;default:now()" json:"joined_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (HouseholdMember) TableName() string {
	return "household_members"
}

// LostPet 失宠报告
type LostPet struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ReportUUID   string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"report_uuid"`
	PetUUID      string         `gorm:"type:varchar(64);index" json:"pet_uuid"`
	PetName      string         `gorm:"type:varchar(32)" json:"pet_name"`
	Species      string         `gorm:"type:varchar(32)" json:"species"`
	LastLocation JSON           `gorm:"type:jsonb" json:"last_location"` // {"lat": xx, "lng": xx, "address": "..."}
	LostTime     time.Time      `json:"lost_time"`
	Status       string         `gorm:"type:varchar(20);default:'searching'" json:"status"` // searching/found/closed
	Reward       string         `gorm:"type:varchar(256)" json:"reward"`                    // 悬赏金额/方式
	Contact      string         `gorm:"type:varchar(128)" json:"contact"`
	Description  string         `gorm:"type:text" json:"description"`
	PhotoURLs    pq.StringArray `gorm:"type:text[]" json:"photo_urls"` // 宠物照片
	Sightings    pq.StringArray `gorm:"type:text[]" json:"sightings"`  // 目击记录JSON数组
	ReporterID   uint           `gorm:"index;not null" json:"reporter_id"`
	SpreadRadius float64        `gorm:"type:decimal(5,2);default:10" json:"spread_radius_km"` // 扩散半径(km)
	TenantID     string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (LostPet) TableName() string {
	return "lost_pet_reports"
}

// BeforeCreate 创建前自动生成 UUID
func (l *LostPet) BeforeCreate(tx *gorm.DB) error {
	if l.ReportUUID == "" {
		l.ReportUUID = uuid.New().String()
	}
	return nil
}

// PetSighting 目击记录
type PetSighting struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ReportUUID    string         `gorm:"type:varchar(64);index;not null" json:"report_uuid"`
	Location      JSON           `gorm:"type:jsonb" json:"location"` // {"lat": xx, "lng": xx}
	SightingTime  time.Time      `json:"sighting_time"`
	Description   string         `gorm:"type:text" json:"description"`
	PhotoURL      string         `gorm:"type:varchar(512)" json:"photo_url"`
	ReporterName  string         `gorm:"type:varchar(64)" json:"reporter_name"`
	Contact       string         `gorm:"type:varchar(128)" json:"contact"`
	IsCredible    bool           `gorm:"type:boolean;default:true" json:"is_credible"`
	ReporterID    uint           `gorm:"index" json:"reporter_id"` // 上报人ID（可选）
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
}

// TableName 指定表名
func (PetSighting) TableName() string {
	return "pet_sightings"
}

// PetHealthReminder 宠物健康提醒
type PetHealthReminder struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PetUUID     string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	ReminderType string        `gorm:"type:varchar(32);not null" json:"reminder_type"` // vaccination/checkup/deworming/flea/grooming/other
	Title       string         `gorm:"type:varchar(128);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	ScheduledAt time.Time      `gorm:"type:timestamp;not null" json:"scheduled_at"`
	RepeatInterval string      `gorm:"type:varchar(32)" json:"repeat_interval"` // none/daily/weekly/monthly/yearly
	IsCompleted bool          `gorm:"type:boolean;default:false" json:"is_completed"`
	CompletedAt *time.Time     `json:"completed_at"`
	Notes       string         `gorm:"type:text" json:"notes"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PetHealthReminder) TableName() string {
	return "pet_health_reminders"
}

// PetCheckup 宠物体检记录
type PetCheckup struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PetUUID     string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	CheckupDate time.Time      `gorm:"type:date;not null" json:"checkup_date"`
	Hospital    string         `gorm:"type:varchar(128)" json:"hospital"`
	Doctor      string         `gorm:"type:varchar(64)" json:"doctor"`
	Weight      float64        `gorm:"type:decimal(5,2)" json:"weight"`
	Symptoms    string         `gorm:"type:text" json:"symptoms"`     // 主诉症状
	Diagnosis   string         `gorm:"type:text" json:"diagnosis"`    // 诊断结果
	Treatment   string         `gorm:"type:text" json:"treatment"`    // 治疗方案
	Prescription string        `gorm:"type:text" json:"prescription"` // 处方
	Cost        float64        `gorm:"type:decimal(10,2)" json:"cost"`
	NextDate    *time.Time     `json:"next_date"`                   // 下次体检日期
	Attachments pq.StringArray `gorm:"type:text[]" json:"attachments"` // 附件
	Notes       string         `gorm:"type:text" json:"notes"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PetCheckup) TableName() string {
	return "pet_checkups"
}
