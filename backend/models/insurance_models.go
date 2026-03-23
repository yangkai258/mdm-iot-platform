package models

import "time"

// InsuranceProduct 保险产品
type InsuranceProduct struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    Name          string    `json:"name"`
    Provider      string    `json:"provider"` // 保险公司名
    CoverageType  string    `json:"coverage_type"` // accident, illness, routine
    MonthlyPrice  float64   `json:"monthly_price"`
    CoverageLimit float64   `json:"coverage_limit"`
    Description   string    `json:"description"`
    Terms         string    `json:"terms"`
    IsActive      bool      `gorm:"default:true" json:"is_active"`
    CreatedAt     time.Time `json:"created_at"`
}

// InsurancePolicy 保险保单
type InsurancePolicy struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    UserID        uint      `gorm:"index" json:"user_id"`
    PetID         uint      `gorm:"index" json:"pet_id"`
    ProductID     uint      `json:"product_id"`
    PlanName      string    `json:"plan_name"`
    Insurer       string    `json:"insurer"`
    PolicyNumber  string    `gorm:"uniqueIndex" json:"policy_number"`
    Premium       float64   `json:"premium"`
    CoverageAmount float64   `json:"coverage_amount"`
    StartDate     time.Time `json:"start_date"`
    EndDate       time.Time `json:"end_date"`
    Status        string    `json:"status"` // active, expired, cancelled
    CreatedAt     time.Time `json:"created_at"`
}

// InsuranceClaim 保险理赔
type InsuranceClaim struct {
    ID           uint       `gorm:"primaryKey" json:"id"`
    PolicyID     uint       `gorm:"index" json:"policy_id"`
    ClaimNumber  string     `gorm:"uniqueIndex" json:"claim_number"`
    ClaimType    string     `json:"claim_type"`
    Amount       float64    `json:"amount"`
    Reason       string     `json:"reason"`
    Description  string     `json:"description"`
    Diagnosis    string     `json:"diagnosis"`
    Status       string     `json:"status"` // pending, approved, rejected, paid
    SubmittedAt  time.Time  `json:"submitted_at"`
    ProcessedAt *time.Time `json:"processed_at"`
    Attachments  string     `json:"attachments"` // JSON array of file URLs
}

// PetMedicalRecord 宠物病历
type PetMedicalRecord struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    PetID        uint      `gorm:"index" json:"pet_id"`
    RecordType   string    `json:"record_type"` // checkup, vaccination, surgery, prescription
    HospitalName string    `json:"hospital_name"`
    DoctorName   string    `json:"doctor_name"`
    Diagnosis    string    `json:"diagnosis"`
    Treatment    string    `json:"treatment"`
    Prescription string    `json:"prescription"`
    RecordDate   time.Time `json:"record_date" binding:"required"`
    Attachments  string    `json:"attachments"` // JSON array
    SyncStatus   string    `json:"sync_status"` // synced, pending
    CreatedAt    time.Time `json:"created_at"`
}

// VetAppointment 兽医预约
type VetAppointment struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    PetID         uint      `gorm:"index" json:"pet_id"`
    VetName       string    `json:"vet_name"`
    ClinicName    string    `json:"clinic_name"`
    AppointmentAt time.Time `json:"appointment_at"`
    Reason        string    `json:"reason"`
    Notes         string    `json:"notes"`
    Status        string    `json:"status"` // scheduled, completed, cancelled
    CreatedAt     time.Time `json:"created_at"`
}
