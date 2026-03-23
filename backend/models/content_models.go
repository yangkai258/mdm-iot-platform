package models

import "time"

// ContentFile 文件库
type ContentFile struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    TenantID    uint      `gorm:"index" json:"tenant_id"`
    UploaderID   uint      `json:"uploader_id"`
    FileName    string    `json:"file_name"`
    FileType    string    `json:"file_type"` // image, video, audio, document
    FileSize    int64     `json:"file_size"` // bytes
    FileURL     string    `json:"file_url"`
    ThumbnailURL string    `json:"thumbnail_url"`
    Category    string    `json:"category"` // emoticon, action, voice, wallpaper
    Tags        string    `json:"tags"` // JSON array
    DownloadCount int     `gorm:"default:0" json:"download_count"`
    LikeCount   int       `gorm:"default:0" json:"like_count"`
    IsPublic    bool      `gorm:"default:false" json:"is_public"`
    CreatedAt   time.Time `json:"created_at"`
}

// ContentDistribution 内容分发记录
type ContentDistribution struct {
    ID          uint       `gorm:"primaryKey" json:"id"`
    ContentID   uint       `gorm:"index" json:"content_id"`
    TargetType  string     `json:"target_type"` // device, user, group
    TargetID    string     `json:"target_id"`
    Status      string     `json:"status"` // pending, sent, failed
    SentAt     *time.Time `json:"sent_at"`
    CreatedAt   time.Time  `json:"created_at"`
}

// AppPackage App应用包
type AppPackage struct {
    ID            uint       `gorm:"primaryKey" json:"id"`
    Name          string     `json:"name"`
    BundleID      string     `gorm:"uniqueIndex" json:"bundle_id"`
    Version       string     `json:"version"`
    Platform      string     `json:"platform"` // ios, android, harmonyos
    FileURL       string     `json:"file_url"`
    FileSize      int64      `json:"file_size"`
    MinOSVersion  string     `json:"min_os_version"`
    ReleaseNotes  string     `json:"release_notes"`
    Status        string     `json:"status"` // draft, pending_review, approved, rejected, deprecated
    ReviewComment string     `json:"review_comment"`
    CreatedAt     time.Time  `json:"created_at"`
    PublishedAt  *time.Time `json:"published_at"`
}

// AppInstall App安装记录
type AppInstall struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    AppID     uint      `gorm:"index" json:"app_id"`
    DeviceID  uint      `gorm:"index" json:"device_id"`
    Version   string    `json:"version"`
    InstallAt time.Time `json:"install_at"`
    Status    string    `json:"status"` // installed, updated, uninstalled
}
