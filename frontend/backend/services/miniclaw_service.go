package services

// MiniClawService MiniClaw服务
type MiniClawService struct{}

func NewMiniClawService() *MiniClawService {
	return &MiniClawService{}
}

// GetLatestFirmware 获取最新固件
func (s *MiniClawService) GetLatestFirmware() (*FirmwareInfo, error) {
	return &FirmwareInfo{Version: "1.0.0"}, nil
}

// FirmwareInfo 固件信息
type FirmwareInfo struct {
	Version string `json:"version"`
	FileURL string `json:"file_url"`
	Size    int64  `json:"size"`
	MD5     string `json:"md5"`
}
