package models

import "time"

// VoiceEmotionRecord 语音情绪记录
type VoiceEmotionRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PetID       uint      `gorm:"index" json:"pet_id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	AudioURL    string    `json:"audio_url"`
	Duration    int       `json:"duration"`
	EmotionType string   `json:"emotion_type"`
	Intensity   int       `json:"intensity"`
	Confidence  float64    `json:"confidence"`
	Transcript  string    `json:"transcript"`
	CreatedAt   time.Time `json:"created_at"`
}
