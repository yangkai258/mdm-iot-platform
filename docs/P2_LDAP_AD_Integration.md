# LDAP/AD 集成 PRD

## 1. 功能概述
LDAP/AD 集成功能允许企业管理员将现有的 LDAP（轻量目录访问协议）或 Active Directory 用户体系与 MDM 平台对接，实现统一身份认证和用户自动同步。

## 2. 页面布局与交互

### 页面路径
`/security/ldap` → `LDAPConfigView.vue`

### 页面布局
- Tab1：LDAP 配置
- Tab2：用户同步
- Tab3：分组-角色映射

### LDAP 配置表单（Tab1）
| 字段 | 类型 | 说明 |
|------|------|------|
| 配置名称 | Input | 必填 |
| LDAP服务器地址 | Input | 必填（ldap://xxx 或 ldaps://xxx）|
| 端口 | Input | 默认389/636（SSL）|
| 基准DN | Input | 必填（dc=company,dc=com）|
| 管理员DN | Input | 必填（cn=admin,dc=company,dc=com）|
| 管理员密码 | Password | 加密存储 |
| 使用SSL | Switch | LDAPS 636端口 |
| 使用TLS | Switch | STARTTLS |
| 用户过滤器 | Input | 默认 (objectClass=user) |
| 分组过滤器 | Input | 默认 (objectClass=group) |
| 同步间隔 | Input | 秒，默认3600 |
| 启用状态 | Switch | 是否启用 |

### 按钮
- 「测试连接」：在不保存配置的情况下测试LDAP连通性
- 「保存配置」
- 「同步用户」：手动触发一次全量同步

### 用户同步（Tab2）
- 同步结果统计：总用户数/新增/更新/跳过
- 用户列表：显示已同步的LDAP用户
- 操作：「查看」「移除」

### 分组-角色映射（Tab3）
- LDAP分组列表
- 分配角色下拉选择

## 3. API 契约

### 获取LDAP配置
- 路径：`GET /api/v1/ldap/config`
- 响应：
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "config_name": "企业LDAP",
    "host": "ldap.company.com",
    "port": 389,
    "base_dn": "dc=company,dc=com",
    "bind_dn": "cn=admin,dc=company,dc=com",
    "use_ssl": false,
    "use_tls": true,
    "user_filter": "(objectClass=user)",
    "group_filter": "(objectClass=group)",
    "sync_interval": 3600,
    "is_enabled": true,
    "last_sync_at": "2024-01-01T00:00:00Z",
    "status": "active",
    "description": ""
  },
  "message": "success"
}
```

### 更新LDAP配置
- 路径：`PUT /api/v1/ldap/config`
- 请求体：
```json
{
  "config_name": "企业LDAP",
  "host": "ldap.company.com",
  "port": 389,
  "base_dn": "dc=company,dc=com",
  "bind_dn": "cn=admin,dc=company,dc=com",
  "bind_password": "新密码",
  "use_ssl": false,
  "use_tls": true,
  "user_filter": "(objectClass=user)",
  "group_filter": "(objectClass=group)",
  "sync_interval": 3600,
  "is_enabled": true,
  "description": ""
}
```

### 测试LDAP连接
- 路径：`POST /api/v1/ldap/test`
- 请求体：同更新配置
- 响应：
```json
{
  "code": 0,
  "data": {
    "success": true,
    "message": "连接成功",
    "server_info": "LDAP Server v3.0"
  }
}
```

### 获取LDAP用户列表
- 路径：`GET /api/v1/ldap/users`
- 参数：`query`（搜索）, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "dn": "uid=zhangsan,ou=users,dc=company,dc=com",
        "username": "zhangsan",
        "email": "zhangsan@company.com",
        "display_name": "张三",
        "groups": ["cn=admins,ou=groups,dc=company,dc=com"]
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 同步LDAP用户
- 路径：`POST /api/v1/ldap/sync`
- 响应：
```json
{
  "code": 0,
  "data": {
    "total_users": 100,
    "added": 5,
    "updated": 10,
    "skipped": 85,
    "errors": [],
    "synced_at": "2024-01-01T00:00:00Z"
  },
  "message": "同步完成"
}
```

### 获取LDAP分组列表
- 路径：`GET /api/v1/ldap/groups`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "dn": "cn=admins,ou=groups,dc=company,dc=com",
        "name": "admins",
        "description": "管理员组"
      }
    ],
    "total": 10
  }
}
```

### 设置分组-角色映射
- 路径：`POST /api/v1/ldap/group-role-mapping`
- 请求体：
```json
{
  "ldap_group_dn": "cn=admins,ou=groups,dc=company,dc=com",
  "ldap_group_name": "admins",
  "role_id": 1
}
```

### 获取所有分组-角色映射
- 路径：`GET /api/v1/ldap/group-role-mappings`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "ldap_group_dn": "cn=admins,ou=groups,dc=company,dc=com",
        "ldap_group_name": "admins",
        "role_id": 1,
        "role_name": "超级管理员"
      }
    ],
    "total": 5
  }
}
```

## 4. 数据模型

### LDAPConfig（LDAP配置）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| config_name | string | 配置名称 |
| host | string | LDAP服务器地址 |
| port | int | 端口 |
| base_dn | string | 基准DN |
| bind_dn | string | 管理员DN |
| bind_password | string | 加密后的密码 |
| use_ssl | bool | 使用SSL |
| use_tls | bool | 使用TLS |
| user_filter | string | 用户过滤器 |
| group_filter | string | 分组过滤器 |
| sync_interval | int | 同步间隔(秒) |
| is_enabled | bool | 是否启用 |
| last_sync_at | datetime | 最后同步时间 |
| status | string | active/inactive/error |
| tenant_id | string | 租户ID |

### LDAPUserMapping（用户映射）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| ldap_dn | string | LDAP DN |
| local_user_id | uint | 本地用户ID |
| username | string | 用户名 |
| email | string | 邮箱 |
| display_name | string | 显示名称 |
| ldap_groups | string | JSON分组列表 |
| sync_status | string | synced/pending/removed |
| last_synced_at | datetime | 最后同步时间 |
| tenant_id | string | 租户ID |

### LDAPGroupRoleMapping（分组角色映射）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| ldap_group_dn | string | LDAP分组DN |
| ldap_group_name | string | LDAP分组名 |
| role_id | uint | 关联角色ID |
| role_name | string | 角色名称（冗余）|
| tenant_id | string | 租户ID |

## 5. 验收标准
- [ ] LDAP配置表单保存成功
- [ ] 测试连接功能正常（支持SSL/TLS）
- [ ] 启用开关控制LDAP同步
- [ ] 手动同步用户成功
- [ ] 同步结果统计准确（新增/更新/跳过）
- [ ] 分组-角色映射保存成功
- [ ] 定时同步功能（按sync_interval）
- [ ] 错误信息清晰展示
