package notification

import (
	"fmt"
	"net/smtp"
)

// EmailService 邮件服务
type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

func NewEmailService(host string, port int, username, password, from string, useTLS bool) *EmailService {
	return &EmailService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		UseTLS:   useTLS,
	}
}

// Send 发送邮件
func (s *EmailService) Send(to []string, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("no recipients")
	}
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.From, to[0], subject, body)
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	err := smtp.SendMail(addr, auth, s.From, to, []byte(msg))
	return err
}

// SendTemplate 发送模板邮件
func (s *EmailService) SendTemplate(to []string, subject string, vars map[string]string) error {
	body := "这是一封模板邮件"
	return s.Send(to, subject, body)
}

// TestConnection 测试邮件连接
func (s *EmailService) TestConnection() error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	err := smtp.SendMail(addr, auth, s.From, []string{s.From}, []byte("Test"))
	return err
}

// EmailServiceInterface 邮件服务接口
type EmailServiceInterface interface {
	Send(to []string, subject, body string) error
	SendTemplate(to []string, subject string, vars map[string]string) error
	TestConnection() error
}
