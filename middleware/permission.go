package middleware

import (
	"net/http"
	"strings"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PermissionCheck 创建权限校验中间件
// 用法: r.GET("/xxx", JWTAuth(), PermissionCheck(db, "device:view"), handler)
func PermissionCheck(db *gorm.DB, requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 超管拥有所有权限
		if IsSuperAdmin(c) {
			c.Next()
			return
		}

		// 获取用户权限列表
		perms := GetUserPermissions(c)
		if perms == nil {
			// 从数据库加载用户权限
			perms = LoadUserPermissionsFromDB(db, c)
		}

		// 检查是否有所需权限
		for _, perm := range perms {
			if perm == requiredPerm || perm == "*" {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "没有权限:" + requiredPerm,
		})
		c.Abort()
	}
}

// PermissionCheckAny 只要拥有任一权限即可通过
func PermissionCheckAny(db *gorm.DB, requiredPerms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsSuperAdmin(c) {
			c.Next()
			return
		}

		perms := GetUserPermissions(c)
		if perms == nil {
			perms = LoadUserPermissionsFromDB(db, c)
		}

		for _, required := range requiredPerms {
			for _, perm := range perms {
				if perm == required || perm == "*" {
					c.Next()
					return
				}
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "没有所需权限",
		})
		c.Abort()
	}
}

// PermissionCheckAll 必须拥有所有指定权限
func PermissionCheckAll(db *gorm.DB, requiredPerms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsSuperAdmin(c) {
			c.Next()
			return
		}

		perms := GetUserPermissions(c)
		if perms == nil {
			perms = LoadUserPermissionsFromDB(db, c)
		}

		permSet := make(map[string]bool)
		for _, p := range perms {
			permSet[p] = true
		}

		for _, required := range requiredPerms {
			if !permSet[required] && !permSet["*"] {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "缺少权限:" + required,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// GetUserPermissions 从 Gin Context 获取用户权限列表
func GetUserPermissions(c *gin.Context) []string {
	if val, exists := c.Get("permissions"); exists {
		if perms, ok := val.([]string); ok {
			return perms
		}
	}
	return nil
}

// SetUserPermissions 设置用户权限列表到 Context
func SetUserPermissions(c *gin.Context, perms []string) {
	c.Set("permissions", perms)
}

// LoadUserPermissionsFromDB 从数据库加载用户权限
func LoadUserPermissionsFromDB(db *gorm.DB, c *gin.Context) []string {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil
	}

	uid, ok := userID.(uint)
	if !ok {
		if uidFloat, ok := userID.(float64); ok {
			uid = uint(uidFloat)
		} else {
			return nil
		}
	}

	// 获取用户角色
	var user models.SysUser
	if err := db.First(&user, uid).Error; err != nil {
		return nil
	}

	// 通过角色获取权限
	var rolePerms []models.SysRolePermission
	db.Where("role_id = ?", user.RoleID).Find(&rolePerms)

	var permIDs []uint
	for _, rp := range rolePerms {
		permIDs = append(permIDs, rp.PermissionID)
	}

	if len(permIDs) == 0 {
		return nil
	}

	var perms []models.SysPermission
	db.Where("id IN ? AND status = 1", permIDs).Find(&perms)

	var result []string
	for _, p := range perms {
		if p.Permission != "" {
			result = append(result, p.Permission)
		}
	}

	// 缓存到 context
	SetUserPermissions(c, result)
	return result
}

// HasPermission 检查用户是否有指定权限
func HasPermission(c *gin.Context, required string) bool {
	if IsSuperAdmin(c) {
		return true
	}
	perms := GetUserPermissions(c)
	if perms == nil {
		return false
	}
	for _, p := range perms {
		if p == required || p == "*" {
			return true
		}
	}
	return false
}

// GetRolePermissions 获取角色的权限点列表（通过code）
func GetRolePermissions(db *gorm.DB, roleID uint) []string {
	var rolePerms []models.SysRolePermission
	db.Where("role_id = ?", roleID).Find(&rolePerms)

	var permIDs []uint
	for _, rp := range rolePerms {
		permIDs = append(permIDs, rp.PermissionID)
	}

	if len(permIDs) == 0 {
		return nil
	}

	var perms []models.SysPermission
	db.Where("id IN ? AND status = 1", permIDs).Find(&perms)

	var result []string
	for _, p := range perms {
		if p.Permission != "" {
			result = append(result, p.Permission)
		}
	}
	return result
}

// BuildPermissionTree 构建权限树
func BuildPermissionTree(perms []models.SysPermission, parentID uint) []map[string]interface{} {
	var tree []map[string]interface{}
	for _, p := range perms {
		if p.ParentID == parentID {
			node := map[string]interface{}{
				"id":         p.ID,
				"name":       p.Name,
				"code":       p.Code,
				"type":       p.Type,
				"path":       p.Path,
				"icon":       p.Icon,
				"sort":       p.Sort,
				"visible":    p.Visible,
				"permission": p.Permission,
				"status":     p.Status,
			}
			children := BuildPermissionTree(perms, p.ID)
			if len(children) > 0 {
				node["children"] = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// MatchPermissionByPrefix 根据前缀匹配权限（如 "device:*" 匹配所有设备权限）
func MatchPermissionByPrefix(perms []string, prefix string) bool {
	for _, p := range perms {
		if p == "*" {
			return true
		}
		// 精确匹配
		if p == prefix {
			return true
		}
		// 通配符匹配（如 device:* 匹配 device:view）
		if strings.HasSuffix(prefix, ":*") {
			base := strings.TrimSuffix(prefix, ":*")
			if strings.HasPrefix(p, base+":") {
				return true
			}
		}
	}
	return false
}
