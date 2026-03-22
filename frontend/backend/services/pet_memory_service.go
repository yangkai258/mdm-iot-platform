package services

import (
	"time"
)

// Memory 记忆
type Memory struct {
	ID        uint      `json:"id"`
	DeviceID  string    `json:"device_id"`
	Type      string    `json:"type"` // short_term/long_term
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// PetMemoryService 宠物记忆服务
type PetMemoryService struct{}

func NewPetMemoryService() *PetMemoryService {
	return &PetMemoryService{}
}

// StoreShortTerm 存储短期记忆
func (s *PetMemoryService) StoreShortTerm(deviceID, content string) error {
	return nil
}

// StoreLongTerm 存储长期记忆
func (s *PetMemoryService) StoreLongTerm(deviceID, content string) error {
	return nil
}

// Retrieve 检索记忆
func (s *PetMemoryService) Retrieve(deviceID, memoryType string) []Memory {
	return []Memory{}
}
