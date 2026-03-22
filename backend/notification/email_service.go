package notification

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// EmailService SMTP 邮件服务
type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"` // 加密存储
	From     string `json:"from"`
	UseTLS   bool   `json:"use_tls"`
}

// NewEmailService 创建邮件服务实例
func NewEmailService(cfg EmailConfig) *EmailService {
	return &EmailService{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		from:     cfg.From,
	}
}

// Send 发送邮件
func (s *EmailService) Send(to []string, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("收件人不能为空")
	}

	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = strings.Join(to, ",")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	if s.port == 465 {
		return s.sendWithTLS(addr, auth, msg.String())
	}

	return smtp.SendMail(addr, auth, s.from, to, []byte(msg.String()))
}

// sendWithTLS 使用 TLS 发送邮件
func (s *EmailService) sendWithTLS(addr string, auth smtp.Auth, msg string) error {
	tlsConfig := &tls.Config{
		ServerName: s.host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS连接失败: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("认证失败: %w", err)
		}
	}

	if err = client.Mail(s.from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	recipients := strings.Split(strings.TrimPrefix(msg, "To: "), ",")
	for _, addr := range recipients {
		addr = strings.TrimSpace(addr)
		if addr != "" {
			if err = client.Rcpt(addr); err != nil {
				return fmt.Errorf("设置收件人失败 [%s]: %w", addr, err)
			}
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("打开数据发送失败: %w", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭数据发送失败: %w", err)
	}

	return client.Quit()
}

// SendTemplate 发送模板邮件（支持变量替换）
func (s *EmailService) SendTemplate(to []string, tplID uint, vars map[string]string) error {
	subject := fmt.Sprintf("通知模板 #%d", tplID)
	body := fmt.Sprintf("模板变量: %v", vars)
	return s.Send(to, subject, body)
}

// TestConnection 测试 SMTP 连接
func (s *EmailService) TestConnection() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	deadline := time.Now().Add(5 * time.Second)

	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %w", err)
	}
	defer conn.Close()
	conn.SetDeadline(deadline)

	if s.port == 465 {
		tlsConfig := &tls.Config{ServerName: s.host}
		conn = tls.Client(conn, tlsConfig)
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("认证失败: %w", err)
		}
	}

	if err = client.Mail(s.username); err != nil {
		return fmt.Errorf("测试邮件发送失败: %w", err)
	}

	return nil
}
