package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ===================== Child Mode =====================

// ChildModeConfig Child mode configuration
type ChildModeConfig struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	UserID              uint       `gorm:"not null;index" json:"user_id"`
	DeviceID            string     `gorm:"type:varchar(36);index" json:"device_id"`
	IsEnabled           bool       `gorm:"default:false" json:"is_enabled"`
	ContentFilterLevel  string     `gorm:"type:varchar(20);default:'moderate'" json:"content_filter_level"` // strict/moderate/relaxed/none
	AllowedCategories   string     `gorm:"type:varchar(500)" json:"allowed_categories"`                     // comma-separated
	BlockedKeywords     string     `gorm:"type:text" json:"blocked_keywords"`                               // comma-separated
	DailyTimeLimit      int        `gorm:"default:0" json:"daily_time_limit"`                               // minutes, 0=unlimited
	SessionDuration     int        `gorm:"default:30" json:"session_duration"`                              // minutes per session
	BreakDuration       int        `gorm:"default:10" json:"break_duration"`                                // break minutes
	AllowedStartTime    string     `gorm:"type:varchar(10);default:'08:00'" json:"allowed_start_time"`    // HH:MM
	AllowedEndTime      string     `gorm:"type:varchar(10);default:'20:00'" json:"allowed_end_time"`      // HH:MM
	TodayUsedMinutes    int        `gorm:"default:0" json:"today_used_minutes"`
	WeekUsedMinutes     int        `gorm:"default:0" json:"week_used_minutes"`
	TotalUsedMinutes    int        `gorm:"default:0" json:"total_used_minutes"`
	TotalSessions       int        `gorm:"default:0" json:"total_sessions"`
	LastSessionAt       *time.Time `json:"last_session_at"`
	EmergencyContact    string     `gorm:"type:varchar(100)" json:"emergency_contact"`
	PinCode             string     `gorm:"type:varchar(10)" json:"pin_code"` // parent PIN
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (ChildModeConfig) TableName() string { return "child_mode_configs" }

// ===================== Elderly Care =====================

// ElderlyCareConfig Elderly care configuration
type ElderlyCareConfig struct {
	ID                    uint       `gorm:"primaryKey" json:"id"`
	UserID                uint       `gorm:"not null;index" json:"user_id"`
	DeviceID              string     `gorm:"type:varchar(36);index" json:"device_id"`
	IsEnabled             bool       `gorm:"default:false" json:"is_enabled"`
	HealthMonitorEnabled  bool       `gorm:"default:true" json:"health_monitor_enabled"`
	HeartRateAlertHigh    int        `gorm:"default:100" json:"heart_rate_alert_high"`   // bpm
	HeartRateAlertLow     int        `gorm:"default:50" json:"heart_rate_alert_low"`     // bpm
	ActivityGoal          int        `gorm:"default:6000" json:"activity_goal"`          // daily steps
	SleepMonitoring       bool       `gorm:"default:true" json:"sleep_monitoring"`
	MedicationReminders   bool       `gorm:"default:true" json:"medication_reminders"`
	MedicationTimes       string     `gorm:"type:varchar(500)" json:"medication_times"` // JSON array ["08:00","12:00"]
	MedicationNames       string     `gorm:"type:varchar(500)" json:"medication_names"`
	CompanionEnabled      bool       `gorm:"default:true" json:"companion_enabled"`
	InteractionFrequency  string     `gorm:"type:varchar(20);default:'normal'" json:"interaction_frequency"` // frequent/normal/sparse
	VoiceCallEnabled      bool       `gorm:"default:true" json:"voice_call_enabled"`
	FallDetectionEnabled  bool       `gorm:"default:true" json:"fall_detection_enabled"`
	EmergencyContactName  string     `gorm:"type:varchar(100)" json:"emergency_contact_name"`
	EmergencyContactPhone string     `gorm:"type:varchar(50)" json:"emergency_contact_phone"`
	EmergencyMessage     string     `gorm:"type:varchar(500)" json:"emergency_message"`
	TotalInteractions     int        `gorm:"default:0" json:"total_interactions"`
	LastInteractionAt     *time.Time `json:"last_interaction_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

func (ElderlyCareConfig) TableName() string { return "elderly_care_configs" }

// ===================== Family Album =====================

// FamilyPhoto Family album photo
type FamilyPhoto struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UUID         string         `gorm:"type:varchar(36);uniqueIndex" json:"uuid"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	DeviceID     string         `gorm:"type:varchar(36);index" json:"device_id"`
	PetID        uint           `gorm:"index" json:"pet_id"`
	Title        string         `gorm:"type:varchar(255)" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	PhotoURL     string         `gorm:"type:varchar(500)" json:"photo_url"`
	ThumbnailURL string         `gorm:"type:varchar(500)" json:"thumbnail_url"`
	FileSize     int64          `gorm:"default:0" json:"file_size"`    // bytes
	Width        int            `gorm:"default:0" json:"width"`         // pixels
	Height       int            `gorm:"default:0" json:"height"`        // pixels
	MimeType     string         `gorm:"type:varchar(50)" json:"mime_type"`
	IsShared     bool           `gorm:"default:false" json:"is_shared"`
	ShareToken   string         `gorm:"type:varchar(64);index" json:"share_token"`
	ShareExpiry  *time.Time     `json:"share_expiry"`
	Category     string         `gorm:"type:varchar(50)" json:"category"`  // daily/health/event/other
	Tags         string         `gorm:"type:varchar(255)" json:"tags"`     // comma-separated
	TakenAt      *time.Time    `json:"taken_at"`
	TakenLat     float64        `gorm:"type:decimal(10,6)" json:"taken_lat"`
	TakenLng     float64        `gorm:"type:decimal(10,6)" json:"taken_lng"`
	TakenAddr    string         `gorm:"type:varchar(255)" json:"taken_addr"`
	LikeCount    int            `gorm:"default:0" json:"like_count"`
	CommentCount int            `gorm:"default:0" json:"comment_count"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate generates UUID and share token before insert
func (f *FamilyPhoto) BeforeCreate(tx *gorm.DB) error {
	if f.UUID == "" {
		f.UUID = uuid.New().String()
	}
	if f.ShareToken == "" {
		f.ShareToken = uuid.New().String()
	}
	return nil
}

func (FamilyPhoto) TableName() string { return "family_photos" }

// FamilyAlbumComment Photo comment
type FamilyAlbumComment struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PhotoUUID  string         `gorm:"type:varchar(36);index" json:"photo_uuid"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	UserName   string         `gorm:"type:varchar(100)" json:"user_name"`
	Content    string         `gorm:"type:text" json:"content"`
	ParentID   uint           `gorm:"default:0" json:"parent_id"`
	LikeCount  int            `gorm:"default:0" json:"like_count"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FamilyAlbumComment) TableName() string { return "family_album_comments" }

// FamilyAlbumLike Photo like
type FamilyAlbumLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PhotoUUID string    `gorm:"type:varchar(36);index" json:"photo_uuid"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (FamilyAlbumLike) TableName() string { return "family_album_likes" }

// ===================== Pet Vaccination =====================

// PetVaccination Pet vaccination record
type PetVaccination struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	PetID          uint       `gorm:"not null;index" json:"pet_id"`
	UserID         uint       `gorm:"not null;index" json:"user_id"`
	VaccineName    string     `gorm:"type:varchar(100);not null" json:"vaccine_name"` // rabies/distemper/parvo/combo
	VaccineType    string     `gorm:"type:varchar(50)" json:"vaccine_type"`           // core/non-core
	LotNumber      string     `gorm:"type:varchar(100)" json:"lot_number"`
	Manufacturer   string     `gorm:"type:varchar(200)" json:"manufacturer"`
	InoculationDate time.Time `gorm:"not null" json:"inoculation_date"`
	InoculationAge string     `gorm:"type:varchar(50)" json:"inoculation_age"`
	InoculationSite string    `gorm:"type:varchar(100)" json:"inoculation_site"`
	Inoculator     string     `gorm:"type:varchar(100)" json:"inoculator"`
	VetClinic      string     `gorm:"type:varchar(200)" json:"vet_clinic"`
	NextDoseDate   *time.Time `json:"next_dose_date"`
	NextDoseMemo   string     `gorm:"type:varchar(255)" json:"next_dose_memo"`
	AdverseReactions string   `gorm:"type:text" json:"adverse_reactions"` // none/mild/severe
	AdverseDetail   string    `gorm:"type:text" json:"adverse_detail"`
	CertificateURL  string    `gorm:"type:varchar(500)" json:"certificate_url"`
	Remark         string    `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (PetVaccination) TableName() string { return "pet_vaccinations" }

// VaccinationReminder Vaccination reminder
type VaccinationReminder struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	VaccinationID  uint       `gorm:"not null;index" json:"vaccination_id"`
	PetID          uint       `gorm:"not null;index" json:"pet_id"`
	UserID         uint       `gorm:"not null;index" json:"user_id"`
	RemindAt       time.Time  `gorm:"not null" json:"remind_at"`
	RemindType     string     `gorm:"type:varchar(20)" json:"remind_type"` // push/sms/email/in_app
	IsSent         bool       `gorm:"default:false" json:"is_sent"`
	IsCompleted    bool       `gorm:"default:false" json:"is_completed"`
	Memo           string     `gorm:"type:varchar(255)" json:"memo"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (VaccinationReminder) TableName() string { return "vaccination_reminders" }

// ===================== Pet Diet Record =====================

// PetDietRecord Pet diet/food record
type PetDietRecord struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	PetID         uint       `gorm:"not null;index" json:"pet_id"`
	UserID        uint       `gorm:"not null;index" json:"user_id"`
	DeviceID      string     `gorm:"type:varchar(36);index" json:"device_id"`
	MealType      string     `gorm:"type:varchar(20);not null" json:"meal_type"`    // breakfast/lunch/dinner/snack/supplement
	FoodName      string     `gorm:"type:varchar(200)" json:"food_name"`
	FoodBrand     string     `gorm:"type:varchar(100)" json:"food_brand"`
	FoodType      string     `gorm:"type:varchar(50)" json:"food_type"`             // dry/wet/fresh/raw/treat
	Amount        float64    `gorm:"type:decimal(10,2)" json:"amount"`
	AmountUnit    string     `gorm:"type:varchar(20);default:'g'" json:"amount_unit"` // g/kg/cup/piece
	Calories      int        `gorm:"default:0" json:"calories"`
	FeedingMethod string     `gorm:"type:varchar(50)" json:"feeding_method"`         // auto/manual/mixed
	AutoFeederID  string     `gorm:"type:varchar(36)" json:"auto_feeder_id"`
	EatTime       *time.Time `json:"eat_time"`
	PlannedTime   *time.Time `json:"planned_time"`
	Duration      int        `gorm:"default:0" json:"duration"` // minutes
	Status        string     `gorm:"type:varchar(20);default:'completed'" json:"status"` // planned/completed/partial/skipped
	LeftOver      float64    `gorm:"type:decimal(10,2);default:0" json:"left_over"`
	AppetiteScore int        `gorm:"default:5" json:"appetite_score"` // 1-10
	HealthNote    string     `gorm:"type:text" json:"health_note"`
	PhotoURL      string     `gorm:"type:varchar(500)" json:"photo_url"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (PetDietRecord) TableName() string { return "pet_diet_records" }

// ===================== Pet Finder =====================

// PetFinderReport Pet finder report
type PetFinderReport struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PetID        uint           `gorm:"not null;index" json:"pet_id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	ReportType   string         `gorm:"type:varchar(20);not null" json:"report_type"` // lost/found/theft
	Title        string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	LastSeenAt   *time.Time     `json:"last_seen_at"`
	LastSeenLat  float64        `gorm:"type:decimal(10,6)" json:"last_seen_lat"`
	LastSeenLng  float64        `gorm:"type:decimal(10,6)" json:"last_seen_lng"`
	LastSeenAddr string         `gorm:"type:varchar(500)" json:"last_seen_addr"`
	ContactName  string         `gorm:"type:varchar(100)" json:"contact_name"`
	ContactPhone string         `gorm:"type:varchar(50)" json:"contact_phone"`
	Reward       float64        `gorm:"type:decimal(10,2)" json:"reward"`
	RewardMemo   string         `gorm:"type:varchar(255)" json:"reward_memo"`
	Photos       StringArray    `gorm:"type:varchar(500)[]" json:"photos"`
	Status       string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/found/closed
	ResolvedAt   *time.Time     `json:"resolved_at"`
	ViewCount    int            `gorm:"default:0" json:"view_count"`
	AlertRadius  float64        `gorm:"type:decimal(10,2);default:5.0" json:"alert_radius"` // km
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (PetFinderReport) TableName() string { return "pet_finder_reports" }

// PetFinderSighting Pet finder sighting report
type PetFinderSighting struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ReportID     uint       `gorm:"not null;index" json:"report_id"`
	ReporterID   uint       `gorm:"not null;index" json:"reporter_id"`
	SightedAt    *time.Time `json:"sighted_at"`
	SightedLat   float64    `gorm:"type:decimal(10,6)" json:"sighted_lat"`
	SightedLng   float64    `gorm:"type:decimal(10,6)" json:"sighted_lng"`
	SightedAddr  string     `gorm:"type:varchar(500)" json:"sighted_addr"`
	Description  string     `gorm:"type:text" json:"description"`
	PetStatus    string     `gorm:"type:varchar(50)" json:"pet_status"` // alone/scared/hungry/with_owner
	PhotoURL     string     `gorm:"type:varchar(500)" json:"photo_url"`
	ContactName  string     `gorm:"type:varchar(100)" json:"contact_name"`
	ContactPhone string     `gorm:"type:varchar(50)" json:"contact_phone"`
	IsVerified   bool       `gorm:"default:false" json:"is_verified"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (PetFinderSighting) TableName() string { return "pet_finder_sightings" }

// PetFinderAlert Pet finder alert/notification
type PetFinderAlert struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ReportID    uint       `gorm:"not null;index" json:"report_id"`
	AlertType   string     `gorm:"type:varchar(20)" json:"alert_type"`  // push/email/sms
	RecipientID uint       `gorm:"index" json:"recipient_id"`
	SentAt      *time.Time `json:"sent_at"`
	Status      string     `gorm:"type:varchar(20)" json:"status"`      // pending/sent/failed
	ErrorMsg    string     `gorm:"type:text" json:"error_msg"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (PetFinderAlert) TableName() string { return "pet_finder_alerts" }
