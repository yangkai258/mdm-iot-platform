package models

import "time"

// Dataset 数据集
type Dataset struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    Name          string    `json:"name"`
    Description   string    `json:"description"`
    Category      string    `json:"category"` // behavior, emotion, health, vocalization
    Tags          string    `json:"tags"` // JSON array
    DataFormat    string    `json:"data_format"` // json, csv, parquet
    DataSize      int64     `json:"data_size"` // bytes
    RecordCount   int       `json:"record_count"` // 样本数量
    FileURL       string    `json:"file_url"`
    License       string    `json:"license"` // MIT, GPL, proprietary
    AccessLevel   string    `json:"access_level"` // public, restricted, private
    DownloadCount int       `gorm:"default:0" json:"download_count"`
    CitationCount int       `gorm:"default:0" json:"citation_count"`
    DOI           string    `json:"doi"`
    PublishedAt   time.Time `json:"published_at"`
    CreatedAt     time.Time `json:"created_at"`
}

// ResearchDatasetVersion 数据集版本
type ResearchDatasetVersion struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    DatasetID   uint      `gorm:"index" json:"dataset_id"`
    Version     string    `json:"version"` // v1.0, v1.1
    Changes     string    `json:"changes"`
    FileURL     string    `json:"file_url"`
    RecordCount int       `json:"record_count"`
    PublishedAt time.Time `json:"published_at"`
}

// ResearchProject 研究项目
type ResearchProject struct {
    ID             uint       `gorm:"primaryKey" json:"id"`
    Name           string     `json:"name"`
    Description    string     `json:"description"`
    OwnerID        uint       `gorm:"index" json:"owner_id"`
    DatasetIDs     string     `json:"dataset_ids"` // JSON array of dataset IDs
    Status         string     `json:"status"` // draft, active, completed, archived
    StartDate      time.Time  `json:"start_date"`
    EndDate       *time.Time `json:"end_date"`
    Findings       string     `json:"findings"` // 研究结论
    PublishedPaper string     `json:"published_paper"`
    CreatedAt      time.Time  `json:"created_at"`
}

// ExperimentRun 实验记录
type ExperimentRun struct {
    ID           uint       `gorm:"primaryKey" json:"id"`
    ProjectID    uint       `gorm:"index" json:"project_id"`
    Name         string     `json:"name"`
    Config       string     `json:"config"` // JSON 实验配置
    Results      string     `json:"results"` // JSON 实验结果
    Metrics      string     `json:"metrics"` // JSON 评估指标
    Status       string     `json:"status"` // running, completed, failed
    StartedAt    time.Time  `json:"started_at"`
    CompletedAt *time.Time `json:"completed_at"`
    ErrorMessage string     `json:"error_message"`
}

// ResearchCollaboration 学术合作
type ResearchCollaboration struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    ProjectID   uint      `gorm:"index" json:"project_id"`
    CollaboratorID uint   `json:"collaborator_id"`
    Role        string    `json:"role"` // researcher, reviewer, admin
    Status      string    `json:"status"` // pending, accepted, rejected
    InvitedAt   time.Time `json:"invited_at"`
    RespondedAt *time.Time `json:"responded_at"`
}
