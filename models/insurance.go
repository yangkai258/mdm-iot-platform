package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ============ 保险产品 ============

// InsuranceProduct 保险产品
type InsuranceProduct struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ProductUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"product_uuid"`
	Name           string         `gorm:"type:varchar(128);not null" json:"name"`
	Code           string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"code"` // 产品代码
	CoverageType   string         `gorm:"type:varchar(32);not null" json:"coverage_type"`     // accident/illness/comprehensive/third_party
	Provider       string         `gorm:"type:varchar(128)" json:"provider"`                 // 保险公司
	CoverageAmount float64        `gorm:"type:decimal(12,2)" json:"coverage_amount"`          // 保额
	Premium        float64        `gorm:"type:decimal(10,2)" json:"premium"`                  // 保费
	PremiumPeriod  string         `gorm:"type:varchar(20);default:'monthly'" json:"premium_period"` // monthly/yearly
	Deductible     float64        `gorm:"type:decimal(10,2);default:0" json:"deductible"`    // 免赔额
	WaitPeriodDays int            `gorm:"type:int;default:0" json:"wait_period_days"`        // 等待期天数
	CoverAgeMin    int            `gorm:"type:int;default:0" json:"cover_age_min"`            // 最小承保年龄（月）
	CoverAgeMax    int            `gorm:"type:int;default:240" json:"cover_age_max"`          // 最大承保年龄（月）
	BreedCodes     pq.StringArray `gorm:"type:text[]" json:"breed_codes"`                   // 承保品种代码
	SpeciesAllowed pq.StringArray `gorm:"type:text[];default:{'dog','cat'}" json:"species_allowed"` // 允许的物种
	Description    string         `gorm:"type:text" json:"description"`
	Terms          string         `gorm:"type:text" json:"terms"`                           // 条款
	Exclusions     string         `gorm:"type:text" json:"exclusions"`                      // 免责条款
	CoverageItems  pq.StringArray `gorm:"type:text[]" json:"coverage_items"`                // 保障项目列表
	MaxClaimAmount float64        `gorm:"type:decimal(12,2);default:0" json:"max_claim_amount"` // 单次最高赔付
	AnnualMaxClaim float64        `gorm:"type:decimal(12,2);default:0" json:"annual_max_claim"`  // 年度最高赔付
	IsActive       bool           `gorm:"type:boolean;default:true" json:"is_active"`
	SortOrder      int            `gorm:"type:int;default:0" json:"sort_order"`
	TenantID       string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (InsuranceProduct) TableName() string {
	return "insurance_products"
}

// BeforeCreate 创建前自动生成 UUID
func (p *InsuranceProduct) BeforeCreate(tx *gorm.DB) error {
	if p.ProductUUID == "" {
		p.ProductUUID = uuid.New().String()
	}
	return nil
}

// ============ 理赔申请 ============

// InsuranceClaim 理赔申请
type InsuranceClaim struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ClaimUUID       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"claim_uuid"`
	ClaimNo         string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"claim_no"` // 理赔单号
	ProductUUID     string         `gorm:"type:varchar(64);index;not null" json:"product_uuid"`   // 关联保险产品
	PetUUID         string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`       // 关联宠物
	OwnerID         uint           `gorm:"index;not null" json:"owner_id"`                       // 申请人
	IncidentDate    time.Time      `gorm:"type:date;not null" json:"incident_date"`             // 出险日期
	IncidentType    string         `gorm:"type:varchar(32);not null" json:"incident_type"`       // accident/illness/liability/other
	IncidentDesc    string         `gorm:"type:text" json:"incident_desc"`                      // 出险描述
	HospitalName    string         `gorm:"type:varchar(128)" json:"hospital_name"`              // 就诊医院
	Diagnosis       string         `gorm:"type:text" json:"diagnosis"`                           // 诊断结果
	ClaimAmount     float64        `gorm:"type:decimal(12,2);not null" json:"claim_amount"`    // 申请金额
	ApprovedAmount  float64        `gorm:"type:decimal(12,2);default:0" json:"approved_amount"` // 核定金额
	RejectionReason string         `gorm:"type:text" json:"rejection_reason"`                  // 拒赔原因
	Status          string         `gorm:"type:varchar(20);default:'draft';" json:"status"`     // draft/submitted/under_review/approved/rejected/paid/closed
	ReviewerID      *uint          `gorm:"index" json:"reviewer_id"`                           // 审核人
	ReviewNotes     string         `gorm:"type:text" json:"review_notes"`                      // 审核备注
	PaidAt          *time.Time     `json:"paid_at"`                                            // 支付时间
	PolicyNo        string         `gorm:"type:varchar(64)" json:"policy_no"`                   // 保单号
	ClaimDocuments  []InsuranceClaimDocument `gorm:"foreignKey:ClaimUUID;references:ClaimUUID" json:"claim_documents"`
	TenantID        string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (InsuranceClaim) TableName() string {
	return "insurance_claims"
}

// BeforeCreate 创建前自动生成 UUID 和 理赔单号
func (c *InsuranceClaim) BeforeCreate(tx *gorm.DB) error {
	if c.ClaimUUID == "" {
		c.ClaimUUID = uuid.New().String()
	}
	if c.ClaimNo == "" {
		c.ClaimNo = generateClaimNo()
	}
	return nil
}

// generateClaimNo 生成理赔单号
func generateClaimNo() string {
	now := time.Now()
	return fmt.Sprintf("CL%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)
}

// ============ 理赔文档 ============

// InsuranceClaimDocument 理赔文档
type InsuranceClaimDocument struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DocUUID     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"doc_uuid"`
	ClaimUUID   string         `gorm:"type:varchar(64);index;not null" json:"claim_uuid"`
	DocType     string         `gorm:"type:varchar(32);not null" json:"doc_type"`      // invoice/receipt/diagnosis/prescription/photo/id_card/other
	FileName    string         `gorm:"type:varchar(256);not null" json:"file_name"`    // 原始文件名
	FileURL     string         `gorm:"type:varchar(512);not null" json:"file_url"`    // 文件访问URL
	FileSize    int64          `gorm:"type:bigint" json:"file_size"`                  // 文件大小（字节）
	MimeType    string         `gorm:"type:varchar(64)" json:"mime_type"`             // MIME类型
	Description string         `gorm:"type:varchar(256)" json:"description"`          // 文档描述
	IsVerified  bool           `gorm:"type:boolean;default:false" json:"is_verified"` // 是否已验证
	VerifiedBy  *uint          `json:"verified_by"`                                   // 验证人
	VerifiedAt  *time.Time     `json:"verified_at"`                                   // 验证时间
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (InsuranceClaimDocument) TableName() string {
	return "insurance_claim_documents"
}

// BeforeCreate 创建前自动生成 UUID
func (d *InsuranceClaimDocument) BeforeCreate(tx *gorm.DB) error {
	if d.DocUUID == "" {
		d.DocUUID = uuid.New().String()
	}
	return nil
}

// ============ 宠物健康档案 ============

// PetHealthRecord 宠物健康档案
type PetHealthRecord struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	RecordUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"record_uuid"`
	PetUUID       string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	RecordType    string         `gorm:"type:varchar(32);not null" json:"record_type"` // checkup/vaccination/medication/surgery/weight/dental/other
	RecordDate    time.Time      `gorm:"type:date;not null" json:"record_date"`        // 记录日期
	Title         string         `gorm:"type:varchar(128);not null" json:"title"`
	Hospital      string         `gorm:"type:varchar(128)" json:"hospital"`           // 医院
	Doctor        string         `gorm:"type:varchar(64)" json:"doctor"`              // 主治医生
	VetName       string         `gorm:"type:varchar(64)" json:"vet_name"`            // 兽医名
	Diagnosis     string         `gorm:"type:text" json:"diagnosis"`                 // 诊断
	Treatment     string         `gorm:"type:text" json:"treatment"`                 // 治疗方案
	Prescription  string         `gorm:"type:text" json:"prescription"`              // 处方
	Medications    pq.StringArray `gorm:"type:text[]" json:"medications"`             // 用药记录
	Cost          float64        `gorm:"type:decimal(10,2)" json:"cost"`              // 费用
	VaccineName   string         `gorm:"type:varchar(64)" json:"vaccine_name"`       // 疫苗名称
	NextDueDate   *time.Time     `json:"next_due_date"`                              // 下次预约/接种日期
	Weight        float64        `gorm:"type:decimal(5,2)" json:"weight"`           // 体重（kg）
	Notes         string         `gorm:"type:text" json:"notes"`                    // 备注
	Attachments   pq.StringArray `gorm:"type:text[]" json:"attachments"`             // 附件URL列表
	IsInsured     bool           `gorm:"type:boolean;default:false" json:"is_insured"` // 是否已关联保险
	InsuranceClaimUUID string    `gorm:"type:varchar(64)" json:"insurance_claim_uuid"` // 关联的理赔UUID
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PetHealthRecord) TableName() string {
	return "pet_health_records"
}

// BeforeCreate 创建前自动生成 UUID
func (r *PetHealthRecord) BeforeCreate(tx *gorm.DB) error {
	if r.RecordUUID == "" {
		r.RecordUUID = uuid.New().String()
	}
	return nil
}
