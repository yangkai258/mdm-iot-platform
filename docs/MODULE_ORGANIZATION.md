# 模块 PRD：组织架构 (Organization)

**版本：** V1.4
**模块负责人：** agentcp
**编制日期：** 2026-03-20

---

## 1. 概述

组织架构模块管理 MDM 中台的内部组织结构，包括公司、部门、岗位和员工信息，为系统管理和数据权限提供组织基础。

**业务目标：**
- 多公司架构支持
- 树形部门管理
- 岗位体系管理
- 员工信息管理
- 基准岗位库

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 | 前端入口/按钮 |
|------|------|--------|----------|--------------|
| 公司管理 | 公司CRUD、状态管理 | P0 | 人工 | 「新建公司」/「编辑」/「删除」按钮 |
| 部门管理 | 部门树CRUD、层级管理 | P0 | 人工 | 「新建部门」/「编辑」/「删除」按钮 |
| 岗位管理 | 岗位CRUD、部门归属 | P1 | 人工 | 「新建岗位」/「编辑」/「删除」按钮 |
| 员工管理 | 员工CRUD、关联用户 | P1 | 人工 | 「新建员工」/「编辑」/「删除」/「批量导入」/「导出」按钮 |
| 基准岗位 | 基准岗位库CRUD | P2 | 人工 | 「新建基准岗位」/「编辑」/「删除」按钮 |
| 批量操作 | 员工/部门批量导入导出 | P2 | 人工 | 「批量导入」/「导出」按钮 |
| 身份同步 | 第三方身份源（LDAP/SSO）同步 | P2 | 人工/定时 | 「同步」/「配置同步」按钮 |

---

## 3. 数据模型

### 3.1 公司表 (companies)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| company_code | string | 公司编码, unique |
| company_name | string | 公司名称 |
| short_name | string | 简称 |
| logo | string | Logo URL |
| province | string | 省 |
| city | string | 市 |
| district | string | 区 |
| address | string | 详细地址 |
| legal_person | string | 法人代表 |
| contact | string | 联系人 |
| phone | string | 联系电话 |
| email | string | 邮箱 |
| status | int | 1=正常 0=禁用 |
| sort | int | 排序 |
| remark | string | 备注 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 3.2 部门表 (departments)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| dept_code | string | 部门编码 |
| dept_name | string | 部门名称 |
| parent_id | uint | 上级部门, nullable |
| level | int | 层级 |
| path | string | 路径 /1/2/3/ |
| manager | string | 负责人 |
| phone | string | 联系电话 |
| email | string | 邮箱 |
| status | int | 1=正常 0=禁用 |
| sort | int | 排序 |
| company_id | uint | 所属公司, index |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

**部门树结构示例：**
```
北京总公司
├── 技术部
│   ├── 前端组
│   └── 后端组
├── 运营部
│   ├── 用户运营组
│   └── 内容运营组
└── 销售部
```

### 3.3 岗位表 (positions)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| pos_code | string | 岗位编码 |
| pos_name | string | 岗位名称 |
| category | string | 岗位类别 (技术/运营/销售等) |
| level | int | 级别 |
| dept_id | uint | 所属部门, index |
| company_id | uint | 所属公司, index |
| description | string | 岗位描述 |
| status | int | 1=正常 0=禁用 |
| sort | int | 排序 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 3.4 员工表 (employees)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| emp_code | string | 工号, unique |
| emp_name | string | 姓名 |
| gender | string | 性别 |
| birth_date | datetime | 出生日期 |
| phone | string | 手机号 |
| email | string | 邮箱 |
| id_card | string | 身份证号 |
| photo | string | 照片 |
| province | string | 省 |
| city | string | 市 |
| district | string | 区 |
| address | string | 详细地址 |
| dept_id | uint | 部门, index |
| position_id | uint | 岗位, index |
| company_id | uint | 公司, index |
| entry_date | datetime | 入职日期 |
| emp_status | int | 员工状态 在职/离职/退休 |
| status | int | 账号状态 |
| remark | string | 备注 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 3.5 基准岗位表 (standard_positions)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| sp_code | string | 编码 |
| sp_name | string | 名称 |
| category | string | 类别 |
| level | int | 级别 |
| description | string | 描述 |
| responsibility | text | 职责 |
| requirement | text | 任职要求 |
| skills | string | 技能要求 |
| status | int | 1=正常 0=禁用 |
| sort | int | 排序 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

---

## 4. 接口定义

### 4.1 公司管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/org/companies | 公司列表（Query: keyword, status, page, page_size） |
| POST | /api/v1/org/companies | 创建公司 |
| PUT | /api/v1/org/companies/:id | 更新公司 |
| DELETE | /api/v1/org/companies/:id | 删除公司 |

### 4.2 部门管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/org/departments | 部门列表（Query: company_id, parent_id, status） |
| GET | /api/v1/org/departments/tree | 部门树（Query: company_id） |
| POST | /api/v1/org/departments | 创建部门（dept_code, dept_name, parent_id, company_id, manager, phone, email, sort） |
| PUT | /api/v1/org/departments/:id | 更新部门 |
| DELETE | /api/v1/org/departments/:id | 删除部门（限制：有子部门时不能删除） |

### 4.3 岗位管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/org/positions | 岗位列表（Query: company_id, dept_id, category, status, page, page_size） |
| POST | /api/v1/org/positions | 创建岗位（pos_code, pos_name, category, level, dept_id, company_id, description, sort） |
| PUT | /api/v1/org/positions/:id | 更新岗位 |
| DELETE | /api/v1/org/positions/:id | 删除岗位 |

### 4.4 员工管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/org/employees | 员工列表（Query: keyword, company_id, dept_id, emp_status, page, page_size） |
| POST | /api/v1/org/employees | 创建员工（emp_code, emp_name, gender, phone, email, dept_id, position_id, company_id, entry_date等） |
| PUT | /api/v1/org/employees/:id | 更新员工 |
| DELETE | /api/v1/org/employees/:id | 删除员工 |
| POST | /api/v1/org/employees/import | 批量导入（multipart/form-data Excel文件） |
| GET | /api/v1/org/employees/export | 批量导出（Query: company_id, dept_id, emp_status） |

**导入响应：** total/success/failed/errors(row+message)

### 4.5 基准岗位

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/org/standard-positions | 基准岗位列表（Query: keyword, category, status, page, page_size） |
| POST | /api/v1/org/standard-positions | 创建基准岗位（sp_code, sp_name, category, level, description, responsibility, requirement, skills, sort） |
| PUT | /api/v1/org/standard-positions/:id | 更新基准岗位 |
| DELETE | /api/v1/org/standard-positions/:id | 删除基准岗位 |

---

## 5. 流程图

### 5.1 部门树构建流程

```
查询所有departments (WHERE company_id=:id)
    │
    ▼
按parent_id分组 → 递归构建树
    │
    buildTree(parentID=null):
        ├─► 找所有parent_id=parentID的部门
        ├─► 对每个部门: children=buildTree(dept.id), dept.children=children
        └─► 返回部门列表
```

### 5.2 员工入职流程

```
HR创建员工记录 ──► 自动生成emp_code（如果未指定）
    │
    ▼
系统自动创建关联的用户账号
    ├─► username=emp_code, password=默认密码, role_id=默认角色
    └─► 关联sys_user_ext(employee_id, dept_id, company_id)
    │
    ▼
员工完成入职
```

---

## 6. 模块联动

| 联动模块 | 联动方式 | 说明 |
|----------|----------|------|
| 系统管理 | 用户关联员工/部门 | sys_user_ext表 |
| 系统管理 | 数据权限范围 | 按company_id/dept_id隔离 |
| 会员管理 | 店铺归属公司 | store.company_id |
| 设备管理 | 运营人员管理设备权限 | 依赖组织架构 |

---

## 7. 验收标准

### P0 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 公司CRUD | 完整增删改查 | 调用各接口验证 |
| 部门树 | 正确返回树形结构 | 创建多级部门验证 |
| 部门删除限制 | 有子部门时不能删除 | 尝试删除父部门 |
| 员工CRUD | 完整增删改查 | 调用各接口验证 |

### P1 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 岗位CRUD | 完整增删改查 | 调用各接口验证 |
| 员工导入 | Excel正确解析，成功/失败分开统计 | 导入100条验证 |
| 员工导出 | 正确导出Excel | 导出后验证内容 |
| 基准岗位 | 完整CRUD | 调用各接口验证 |

### P2 验收标准

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 部门搜索 | 支持关键词搜索 | 搜索验证 |
| 员工筛选 | 支持多条件组合筛选 | 多条件验证 |

---

## 8. UI设计指引

### 8.1 页面结构
- **左侧菜单**：组织架构 → 公司管理 / 部门管理 / 岗位管理 / 员工管理 / 基准岗位
- **顶部区域**：左侧树形结构（公司→部门→岗位）可视化
- **中间区域**：右侧数据表格 + 顶部操作区
- **底部区域**：分页组件

### 8.2 组件选用
| 组件 | 用途 |
|------|------|
| a-tree | 左侧组织架构树形图（公司/部门/岗位层级）|
| a-table | 右侧数据列表（公司列表/部门列表/员工列表）|
| a-drawer | 创建/编辑公司、创建/编辑部门、员工详情 |
| a-modal | 删除确认、批量导入预览 |
| a-form | 公司表单、部门表单、岗位表单、员工表单 |
| a-upload | 员工Excel批量导入上传 |
| a-cascader | 部门级联选择（上级部门）|
| a-input-search | 公司名/部门名/员工名关键词搜索 |
| a-select | 状态筛选、公司筛选 |
| a-tag | 部门层级标签、员工状态标签 |
| a-button-group | 工具栏按钮组（导入/导出/新建）|

### 8.3 参考模板
```
┌──────────────────────────────────────────────────────────────┐
│  组织架构                                                      │
├────────────────┬───────────────────────────────────────────────┤
│  组织架构树    │  【公司列表 / 部门列表 / 员工列表 Tab切换】    │
│                │                                                │
│  北京总公司    │  ┌─────────────────────────────────────────┐  │
│   ├─技术部    │  │ [+新建公司]  [批量导入]  [导出]          │  │
│   │ ├─前端组 │  ├─────────────────────────────────────────┤  │
│   │ └─后端组 │  │ 公司名称   │编码│状态│部门数│操作        │  │
│   ├─运营部    │  │ 北京科技   │C001│正常│  5   │编辑删除    │  │
│   └─销售部    │  │ 上海分公司 │C002│正常│  3   │编辑删除    │  │
│                │  └─────────────────────────────────────────┘  │
│  上海分公司    │                                                │
│   └─运营部    │  【部门管理Tab】                              │
│               │  ┌─────────────────────────────────────────┐  │
│               │  │ 部门名称  │上级部门 │负责人│状态│操作    │  │
│               │  │ 技术部    │北京总司 │ 张三 │正常│编辑删除│  │
│               │  │ 前端组    │技术部   │ 李四 │正常│编辑删除│  │
│               │  └─────────────────────────────────────────┘  │
└───────────────┴───────────────────────────────────────────────┘
```

### 8.4 交互流程
```
公司管理
    ├── 左侧树选择公司 ──► 右侧显示该公司下的部门树
    ├── 点击「新建公司」──► a-drawer ──► 填写公司信息 ──► 创建
    └── 点击「编辑/删除」──► 操作公司信息

部门管理
    ├── 点击左侧部门树节点 ──► 右侧显示该部门详情+子部门
    ├── 点击「新建部门」──► a-drawer ──► 选择上级部门+填写信息 ──► 创建
    ├── 点击「编辑」──► a-drawer ──► 修改信息（上级部门不可变更）
    └── 点击「删除」──► 检查是否有子部门/员工 ──► 有则不可删除

员工管理
    ├── 点击「新建员工」──► a-drawer ──► 选择公司+部门+岗位+填写信息 ──► 创建
    ├── 点击「编辑」──► a-drawer ──► 修改信息 ──► 保存
    ├── 点击「批量导入」──► a-upload上传Excel ──► 预览确认 ──► 导入结果
    └── 点击「导出」──► 生成Excel下载
```

### 8.5 关键状态显示
- **组织树**：a-tree，展开/折叠动画，选中高亮
- **部门层级**：a-tag，数字=层级（1=总公司，2=二级部门...）
- **员工状态**：a-tag，在职=绿色，离职=灰色
- **公司状态**：a-tag，正常=绿色，禁用=灰色
- **导入结果**：成功行绿色背景，失败行红色背景+错误提示

---

## 附录 B. 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.0 | 2026-03-20 | agentcp | 初稿，基于代码调研 |
| V1.2 | 2026-03-20 | agentcp | 修订功能列表，补充触发方式和前端入口按钮列 |
| V1.4 | 2026-03-20 | agentcp | 重建文档结构，统一使用8章节格式，合并重复的八、九章节