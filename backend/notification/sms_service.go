package notification

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// SMSService SMS 短信服务
type SMSService struct {
	provider  string // "aliyun" / "tencent"
	accessKey string
	secretKey string
	signName  string
	region    string
	client    *http.Client
}

// SMSConfig SMS 配置
type SMSConfig struct {
	Provider  string `json:"provider"`  // "aliyun" / "tencent"
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"` // 加密存储
	SignName  string `json:"sign_name"`
	Region    string `json:"region"` // 区域，默认 cn-hangzhou
}

// NewSMSService 创建 SMS 服务实例
func NewSMSService(cfg SMSConfig) *SMSService {
	return &SMSService{
		provider:  cfg.Provider,
		accessKey: cfg.AccessKey,
		secretKey: cfg.SecretKey,
		signName:  cfg.SignName,
		region:    cfg.Region,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send 发送短信
func (s *SMSService) Send(phones []string, templateCode string, params map[string]string) error {
	if len(phones) == 0 {
		return fmt.Errorf("手机号不能为空")
	}

	switch s.provider {
	case "aliyun":
		return s.sendAliyun(phones, templateCode, params)
	case "tencent":
		return s.sendTencent(phones, templateCode, params)
	default:
		return fmt.Errorf("不支持的短信服务商: %s", s.provider)
	}
}

// SendBatch 批量发送短信（相同模板，不同手机号）
func (s *SMSService) SendBatch(phones []string, templateCode string, params map[string]string) error {
	if len(phones) == 0 {
		return fmt.Errorf("手机号不能为空")
	}

	switch s.provider {
	case "aliyun":
		return s.sendBatchAliyun(phones, templateCode, params)
	case "tencent":
		// 腾讯云批量发送也是调用单个接口
		return s.sendTencent(phones, templateCode, params)
	default:
		return fmt.Errorf("不支持的短信服务商: %s", s.provider)
	}
}

// sendAliyun 阿里云短信发送
func (s *SMSService) sendAliyun(phones []string, templateCode string, params map[string]string) error {
	if len(phones) == 0 {
		return fmt.Errorf("手机号不能为空")
	}

	// 阿里云短信 API
	apiURL := fmt.Sprintf("https://dysmsapi.aliyuncs.com/?RegionId=%s", s.region)
	if s.region == "" {
		apiURL = "https://dysmsapi.aliyuncs.com/?RegionId=cn-hangzhou"
	}

	phoneStr := strings.Join(phones, ",")

	// 构造参数
	payload := url.Values{}
	payload.Set("AccessKeyId", s.accessKey)
	payload.Set("Action", "SendBatchSms")
	payload.Set("SignName", s.signName)
	payload.Set("TemplateCode", templateCode)
	payload.Set("PhoneNumberJson", fmt.Sprintf(`["%s"]`, phoneStr))

	// 模板变量
	if len(params) > 0 {
		paramJSON, _ := json.Marshal(params)
		payload.Set("TemplateParamJson", fmt.Sprintf(`[%s]`, string(paramJSON)))
	}

	// 生成签名
	signature := s.aliyunSign(payload)
	payload.Set("Signature", signature)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("阿里云短信请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if respCode, ok := result["Code"].(string); ok && respCode != "OK" {
		return fmt.Errorf("阿里云短信发送失败: %s - %v", respCode, result["Message"])
	}

	return nil
}

// sendBatchAliyun 阿里云批量短信
func (s *SMSService) sendBatchAliyun(phones []string, templateCode string, params map[string]string) error {
	return s.sendAliyun(phones, templateCode, params)
}

// sendTencent 腾讯云短信发送
func (s *SMSService) sendTencent(phones []string, templateCode string, params map[string]string) error {
	// 腾讯云 SecretId/SecretKey 签名方式
	secretId := s.accessKey

	// 腾讯云短信 API
	apiURL := "https://sms.tencentcloudapi.com/"

	// 构建正文参数
	type SMSPhone struct {
		PhoneNumber string `json:"PhoneNumber"`
	}
	type TemplateParam struct {
		Value string `json:"Value"`
	}

	phoneNumbers := make([]string, 0, len(phones))
	for _, p := range phones {
		// 腾讯云手机号格式需要带 +86
		if !strings.HasPrefix(p, "+") && !strings.HasPrefix(p, "86") {
			p = "+86" + p
		}
		phoneNumbers = append(phoneNumbers, p)
	}

	templateParams := make([]map[string]string, 0)
	for k, v := range params {
		templateParams = append(templateParams, map[string]string{k: v})
	}

	payload := map[string]interface{}{
		"Version":  "2021-01-11",
		"Action":   "SendSms",
		"Region":    s.region,
		"SecretId":  secretId,
		"SmsType":   0,
		"From":      s.signName,
		"SmsSdkAppId": secretId, // 通常 AppId
		"PhoneNumberSet": phoneNumbers,
		"TemplateId":     templateCode,
	}
	if len(params) > 0 {
		payload["TemplateParamSet"] = templateParams
	}

	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}

	// 签名
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	canonicalRequest := fmt.Sprintf("POST\n/\n\ncontent-type:application/json\nhost:sms.tencentcloudapi.com\n\ncontent-type;host\n%s",
		string(payloadBytes))
	h := sha1.New()
	h.Write([]byte(canonicalRequest))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TC-Action", "SendSms")
	req.Header.Set("X-TC-Version", "2021-01-11")
	req.Header.Set("X-TC-Timestamp", timestamp)
	req.Header.Set("X-TC-Region", s.region)
	req.Header.Set("X-TC-SecretId", secretId)
	req.Header.Set("X-TC-Signature", signature)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("腾讯云短信请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("腾讯云短信发送失败: HTTP %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// aliyunSign 生成阿里云 API 签名
func (s *SMSService) aliyunSign(params url.Values) string {
	// 1. 排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 拼接
	var canonicalized strings.Builder
	for _, k := range keys {
		canonicalized.WriteString(url.QueryEscape(k))
		canonicalized.WriteString("=")
		canonicalized.WriteString(url.QueryEscape(params.Get(k)))
		canonicalized.WriteString("&")
	}
	signedStr := canonicalized.String()
	signedStr = strings.TrimSuffix(signedStr, "&")

	// 3. HMAC-SHA1
	mac := hmac.New(sha1.New, []byte(s.secretKey+"&"))
	mac.Write([]byte("GET&%2F&"+url.QueryEscape(signedStr)))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
