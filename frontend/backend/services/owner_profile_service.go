package services

// OwnerProfile 主人画像
type OwnerProfile struct {
	ID           uint   `json:"id"`
	Nickname     string `json:"nickname"`
	Preferences  string `json:"preferences"`
	Interests    string `json:"interests"`
	ContactTime  string `json:"contact_time"`
}

// OwnerProfileService 主人画像服务
type OwnerProfileService struct{}

func NewOwnerProfileService() *OwnerProfileService {
	return &OwnerProfileService{}
}

func (s *OwnerProfileService) GetProfile(ownerID string) (*OwnerProfile, error) {
	return &OwnerProfile{ID: 1, Nickname: "主人"}, nil
}

func (s *OwnerProfileService) UpdateProfile(ownerID string, profile *OwnerProfile) error {
	return nil
}

func (s *OwnerProfileService) GetPreferences(ownerID string) (map[string]string, error) {
	return map[string]string{}, nil
}
