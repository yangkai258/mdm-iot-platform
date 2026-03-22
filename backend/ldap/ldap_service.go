package ldap

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"mdm-backend/utils"

	"github.com/go-ldap/ldap/v3"
)

// LDAPService LDAP/AD 连接服务
type LDAPService struct {
	Host         string `json:"host"`          // ldap://your-ad-server.com 或 ldaps://
	Port         int    `json:"port"`          // 389 / 636 (SSL)
	BaseDN       string `json:"base_dn"`        // dc=company,dc=com
	BindDN       string `json:"bind_dn"`        // cn=admin,dc=company,dc=com
	BindPassword string `json:"-"`              // 管理员密码（加密存储）
	UseSSL       bool   `json:"use_ssl"`        // 是否使用 SSL
	UseTLS       bool   `json:"use_tls"`        // 是否使用 STARTTLS
	UserFilter   string `json:"user_filter"`   // (objectClass=user)
	GroupFilter  string `json:"group_filter"`  // (objectClass=group)
	TenantID     string `json:"tenant_id"`      // 租户ID
}

// LDAPUser LDAP 用户结构
type LDAPUser struct {
	DN         string   `json:"dn"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	DisplayName string  `json:"display_name"`
	Groups     []string `json:"groups"`
	Department string   `json:"department"`
	JobTitle   string   `json:"job_title"`
	Phone      string   `json:"phone"`
}

// LDAPGroup LDAP 分组结构
type LDAPGroup struct {
	DN      string `json:"dn"`
	Name    string `json:"name"`
	Members int    `json:"members"`
}

// SyncResult 同步结果
type SyncResult struct {
	TotalUsers  int      `json:"total_users"`
	Added       int      `json:"added"`
	Updated     int      `json:"updated"`
	Skipped     int      `json:"skipped"`
	Errors      []string `json:"errors"`
	SyncedAt    time.Time `json:"synced_at"`
}

// TestResult 测试连接结果
type TestResult struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Server  string   `json:"server"`
	Users   int      `json:"users_found,omitempty"`
}

// Connect 建立 LDAP 连接
func (s *LDAPService) Connect() (*ldap.Conn, error) {
	var conn *ldap.Conn
	var err error

	if s.UseSSL {
		// LDAPS 连接 (SSL)
		conn, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port), &tls.Config{
			InsecureSkipVerify: false,
		})
	} else {
		// 普通 LDAP 连接
		conn, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	}

	if err != nil {
		return nil, fmt.Errorf("LDAP 连接失败: %w", err)
	}

	// STARTTLS
	if s.UseTLS && !s.UseSSL {
		err = conn.StartTLS(&tls.Config{
			InsecureSkipVerify: false,
		})
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("LDAP STARTTLS 失败: %w", err)
		}
	}

	// Bind (认证)
	err = conn.Bind(s.BindDN, s.BindPassword)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("LDAP Bind 失败: %w", err)
	}

	return conn, nil
}

// TestConnection 测试 LDAP 连接
func (s *LDAPService) TestConnection() (*TestResult, error) {
	conn, err := s.Connect()
	if err != nil {
		return &TestResult{
			Success: false,
			Message: err.Error(),
			Server:  fmt.Sprintf("%s:%d", s.Host, s.Port),
		}, nil
	}
	defer conn.Close()

	// 尝试搜索用户数量
	searchRequest := ldap.NewSearchRequest(
		s.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		s.UserFilter,
		[]string{"dn"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return &TestResult{
			Success: false,
			Message: fmt.Sprintf("搜索用户失败: %v", err),
			Server:  fmt.Sprintf("%s:%d", s.Host, s.Port),
		}, nil
	}

	return &TestResult{
		Success: true,
		Message: "连接成功",
		Server:  fmt.Sprintf("%s:%d", s.Host, s.Port),
		Users:   len(result.Entries),
	}, nil
}

// SearchUsers 搜索 LDAP 用户
func (s *LDAPService) SearchUsers(query string) ([]LDAPUser, error) {
	conn, err := s.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 构建搜索过滤器
	filter := s.UserFilter
	if query != "" {
		// 支持按 username, email, displayName 搜索
		filter = fmt.Sprintf("(&%s(|(sAMAccountName=*%s*)(mail=*%s*)(displayName=*%s*)))",
			s.UserFilter, query, query, query)
	}

	searchRequest := ldap.NewSearchRequest(
		s.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		100, 30, false,
		filter,
		[]string{"dn", "sAMAccountName", "mail", "displayName", "memberOf", "department", "title", "telephoneNumber"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("搜索用户失败: %w", err)
	}

	users := make([]LDAPUser, 0, len(result.Entries))
	for _, entry := range result.Entries {
		user := LDAPUser{
			DN:         entry.DN,
			Username:   entry.GetAttributeValue("sAMAccountName"),
			Email:      entry.GetAttributeValue("mail"),
			DisplayName: entry.GetAttributeValue("displayName"),
			Department: entry.GetAttributeValue("department"),
			JobTitle:   entry.GetAttributeValue("title"),
			Phone:      entry.GetAttributeValue("telephoneNumber"),
		}

		// 获取用户所属组
		groupEntries := entry.GetAttributeValues("memberOf")
		user.Groups = make([]string, 0, len(groupEntries))
		for _, g := range groupEntries {
			// 提取 CN=xxx,OU=xxx,...
			parts := strings.Split(g, ",")
			if len(parts) > 0 && strings.HasPrefix(parts[0], "CN=") {
				user.Groups = append(user.Groups, strings.TrimPrefix(parts[0], "CN="))
			}
		}

		users = append(users, user)
	}

	return users, nil
}

// SearchGroups 搜索 LDAP 分组
func (s *LDAPService) SearchGroups() ([]LDAPGroup, error) {
	conn, err := s.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	filter := s.GroupFilter
	if filter == "" {
		filter = "(objectClass=group)"
	}

	searchRequest := ldap.NewSearchRequest(
		s.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 60, false,
		filter,
		[]string{"dn", "cn", "description", "member"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("搜索分组失败: %w", err)
	}

	groups := make([]LDAPGroup, 0, len(result.Entries))
	for _, entry := range result.Entries {
		group := LDAPGroup{
			DN:      entry.DN,
			Name:    entry.GetAttributeValue("cn"),
			Members: len(entry.GetAttributeValues("member")),
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// Authenticate 验证 LDAP 用户密码
func (s *LDAPService) Authenticate(username, password string) (*LDAPUser, error) {
	// 先连接管理账号搜索用户 DN
	conn, err := s.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 搜索用户 DN
	filter := fmt.Sprintf("(&%s(sAMAccountName=%s))", s.UserFilter, ldap.EscapeFilter(username))
	searchRequest := ldap.NewSearchRequest(
		s.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1, 10, false,
		filter,
		[]string{"dn", "sAMAccountName", "mail", "displayName", "memberOf", "department", "title", "telephoneNumber"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("搜索用户失败: %w", err)
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	userDN := result.Entries[0].DN

	// 使用用户账号 Bind 验证密码
	userConn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return nil, fmt.Errorf("连接失败: %w", err)
	}
	defer userConn.Close()

	if s.UseTLS && !s.UseSSL {
		err = userConn.StartTLS(&tls.Config{InsecureSkipVerify: false})
		if err != nil {
			return nil, fmt.Errorf("STARTTLS 失败: %w", err)
		}
	}

	err = userConn.Bind(userDN, password)
	if err != nil {
		return nil, fmt.Errorf("密码验证失败")
	}

	// 密码验证成功，构建用户信息
	entry := result.Entries[0]
	user := &LDAPUser{
		DN:         entry.DN,
		Username:   entry.GetAttributeValue("sAMAccountName"),
		Email:      entry.GetAttributeValue("mail"),
		DisplayName: entry.GetAttributeValue("displayName"),
		Department: entry.GetAttributeValue("department"),
		JobTitle:   entry.GetAttributeValue("title"),
		Phone:      entry.GetAttributeValue("telephoneNumber"),
	}

	groupEntries := entry.GetAttributeValues("memberOf")
	user.Groups = make([]string, 0, len(groupEntries))
	for _, g := range groupEntries {
		parts := strings.Split(g, ",")
		if len(parts) > 0 && strings.HasPrefix(parts[0], "CN=") {
			user.Groups = append(user.Groups, strings.TrimPrefix(parts[0], "CN="))
		}
	}

	return user, nil
}

// GetEncryptedPassword 获取加密后的密码
func (s *LDAPService) GetEncryptedPassword() (string, error) {
	return utils.EncryptAES(s.BindPassword)
}

// SetDecryptedPassword 设置解密后的密码
func (s *LDAPService) SetDecryptedPassword(encrypted string) error {
	decrypted, err := utils.DecryptAES(encrypted)
	if err != nil {
		return err
	}
	s.BindPassword = decrypted
	return nil
}
