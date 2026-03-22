package notification

import (
	"fmt"
)

// SMSService 短信服务
type SMSService struct {
	Provider  string
	AccessKey string
	SecretKey string
	SignName  string
	Endpoint  string
}

func NewSMSService(provider, accessKey, secretKey, signName string) *SMSService {
	endpoint := "https://dysmsapi.aliyuncs.com"
	if provider == "tencent" {
		endpoint = "https://sms.tencentcloudapi.com"
	}
	return &SMSService{
		Provider:  provider,
		AccessKey: accessKey,
		SecretKey: secretKey,
		SignName:  signName,
		Endpoint:  endpoint,
	}
}

// Send 发送短信
func (s *SMSService) Send(phones []string, templateCode string, params map[string]string) error {
	if len(phones) == 0 {
		return fmt.Errorf("no phones")
	}
	// 模拟发送成功
	fmt.Printf("[SMS] Sending to %s, template: %s\n", phones[0], templateCode)
	return nil
}

// SendBatch 批量发送短信
func (s *SMSService) SendBatch(phones []string, templateCode string, params map[string]string) error {
	for _, phone := range phones {
		if err := s.Send([]string{phone}, templateCode, params); err != nil {
			return err
		}
	}
	return nil
}

// SMSServiceInterface 短信服务接口
type SMSServiceInterface interface {
	Send(phones []string, templateCode string, params map[string]string) error
	SendBatch(phones []string, templateCode string, params map[string]string) error
}
