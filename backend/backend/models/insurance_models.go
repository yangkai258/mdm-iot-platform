package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InsuranceProduct 保险产品
type InsuranceProduct struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ProductUUID    string         `gorm:"size:64;uniqueIndex;not null" json:"product_uuid"`
	Name           string         `gorm:"size:128;not null" json:"name"`
	Code           string         `gorm:"size:64;uniqueIndex;not null" json:"code"`
	CoverageType   string         `gorm:"size:32;not null" json:"coverage_type"`
	Provider       string         `gorm:"size:128" json:"provider"`
	CoverageAmount float64        `gorm:"type:numeric(12,2)" json:"coverage_amount"`
	Premium        float64        `gorm:"type:numeric(10,2)" json:"premium"`
	PremiumPeriod  string         `gorm:"size:20;default:'monthly'" json:"premium_period"`
	Deductible     float64        `gorm:"type:numeric(10,2);default:0" json:"deductible"`
	WaitPeriodDays int64          `gorm:"default:0" json:"wait_period_days"`
	CoverAgeMin    int64          `gorm:"default:0" json:"cover_age_min"`
	CoverAgeMax    int64          `gorm:"default:240" json:"cover_age_max"`
	BreedCodes     []string       `gorm:"type:text[]" json:"breed_codes"`
	SpeciesAllowed []string       `gorm:"type:text[]" json:"species_allowed"`
	Description    string         `gorm:"type:text" json:"description"`
	Terms          string         `gorm:"type:text" json:"terms"`
	Exclusions     string         `gorm:"type:text" json:"exclusions"`
	CoverageItems  []string       `gorm:"type:text[]" json:"coverage_items"`
	MaxClaimAmount float64        `gorm:"type:numeric(12,2);default:0" json:"max_claim_amount"`
	AnnualMaxClaim float64        `gorm:"type:numeric(12,2);default:0" json:"annual_max_claim"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	SortOrder      int64          `gorm:"default:0" json:"sort_order"`
	TenantID       string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// BeforeCreate 创建前自动生成 UUID
func (i *InsuranceProduct) BeforeCreate(tx *gorm.DB) error {
	if i.ProductUUID == "" {
		i.ProductUUID = uuid.New().String()
	}
	return nil
}

// InsuranceClaim 保险理赔
type InsuranceClaim struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ClaimUUID      string         `gorm:"size:64;uniqueIndex;not null" json:"claim_uuid"`
	ClaimNo        string         `gorm:"size:32;uniqueIndex;not null" json:"claim_no"`
	ProductUUID    string         `gorm:"size:64;not null;index" json:"product_uuid"`
	PetUUID        string         `gorm:"size:64;not null;index" json:"pet_uuid"`
	OwnerID        int64          `gorm:"not null;index" json:"owner_id"`
	IncidentDate   time.Time      `gorm:"type:date;not null" json:"incident_date"`
	IncidentType   string         `gorm:"size:32;not null" json:"incident_type"`
	IncidentDesc   string         `gorm:"type:text" json:"incident_desc"`
	HospitalName   string         `gorm:"size:128" json:"hospital_name"`
	Diagnosis      string         `gorm:"type:text" json:"diagnosis"`
	ClaimAmount    float64        `gorm:"type:numeric(12,2);not null" json:"claim_amount"`
	ApprovedAmount float64        `gorm:"type:numeric(12,2);default:0" json:"approved_amount"`
	RejectionReason string        `gorm:"type:text" json:"rejection_reason"`
	Status         string         `gorm:"size:20;default:'draft'" json:"status"`
	ReviewerID     int64          `gorm:"index" json:"reviewer_id"`
	ReviewNotes    string         `gorm:"type:text" json:"review_notes"`
	PaidAt         *time.Time     `json:"paid_at"`
	PolicyNo       string         `gorm:"size:64" json:"policy_no"`
	TenantID       string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Documents      []InsuranceClaimDocument `gorm:"foreignKey:ClaimUUID;references:ClaimUUID" json:"documents,omitempty"`
}

// BeforeCreate 创建前自动生成 UUID 和 ClaimNo
func (ic *InsuranceClaim) BeforeCreate(tx *gorm.DB) error {
	if ic.ClaimUUID == "" {
		ic.ClaimUUID = uuid.New().String()
	}
	if ic.ClaimNo == "" {
		ic.ClaimNo = "CLM" + time.Now().Format("20060102150405")
	}
	return nil
}

// InsuranceClaimDocument 理赔文档
type InsuranceClaimDocument struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DocUUID     string         `gorm:"size:64;uniqueIndex;not null" json:"doc_uuid"`
	ClaimUUID   string         `gorm:"size:64;not null;index" json:"claim_uuid"`
	DocType     string         `gorm:"size:32;not null" json:"doc_type"`
	FileName    string         `gorm:"size:256;not null" json:"file_name"`
	FileURL     string         `gorm:"size:512;not null" json:"file_url"`
	FileSize    int64          `json:"file_size"`
	MimeType    string         `gorm:"size:64" json:"mime_type"`
	Description string         `gorm:"size:256" json:"description"`
	IsVerified  bool           `gorm:"default:false" json:"is_verified"`
	VerifiedBy  int64          `json:"verified_by"`
	VerifiedAt  *time.Time     `json:"verified_at"`
	TenantID    string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// BeforeCreate 创建前自动生成 DocUUID
func (d *InsuranceClaimDocument) BeforeCreate(tx *gorm.DB) error {
	if d.DocUUID == "" {
		d.DocUUID = uuid.New().String()
	}
	return nil
}
