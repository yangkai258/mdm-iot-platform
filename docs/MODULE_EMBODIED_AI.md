# 模块 PRD：具身智能（MODULE_EMBODIED_AI）

**版本：** V1.0
**所属Phase：** Phase 3（Sprint 21-22）
**优先级：** P1
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

具身智能模块是宠物机器人感知物理世界、执行动作的核心能力，通过环境感知、空间认知、自主探索、动作模仿四大核心能力，实现宠物在物理空间的智能移动和交互，构建真实世界中的智能陪伴体验。

### 1.2 核心价值

- **真实交互**：宠物能感知环境并做出真实物理反应
- **自主行动**：宠物能自主探索环境，不需要人为干预
- **动作学习**：宠物能从人类示范中学习新动作

### 1.3 范围边界

**包含：**
- 环境感知（视觉/距离/触觉）
- 空间认知（地图构建/定位）
- 自主探索（导航/避障）
- 动作模仿（学习/执行）
- 具身AI决策引擎
- 具身AI安全边界

**不包含：**
- 情感计算（MODULE_AFFECTIVE_COMPUTING）
- 仿真测试（MODULE_SIMULATION）
- 设备端固件实现（固件层）
- 电机/舵机控制（硬件抽象层）

---

## 二、功能详情

### 2.1 环境感知

#### 2.1.1 视觉感知

| 功能 | 说明 |
|------|------|
| 物体识别 | 识别家庭常见物体（人/宠物/家具/玩具） |
| 物体检测 | 检测物体位置和边界框 |
| 场景分类 | 识别当前场景（客厅/卧室/厨房） |
| 深度估计 | 估计物体距离 |
| 动态物体追踪 | 追踪移动的人和宠物 |
| 人的姿态估计 | 识别人体关键点和姿态 |
| 手势识别 | 识别用户手势指令 |

#### 2.1.2 距离感知

| 功能 | 说明 |
|------|------|
| 障碍物检测 | 检测前方障碍物距离 |
| 空间测量 | 测量物体尺寸和距离 |
| 边界检测 | 检测楼梯/悬崖等危险边界 |
| 地面检测 | 判断地面状态（地毯/地板/瓷砖） |

#### 2.1.3 触觉感知

| 功能 | 说明 |
|------|------|
| 触摸检测 | 检测是否被触摸 |
| 触摸力度 | 估计触摸力度 |
| 触摸位置 | 确定触摸位置 |

#### 2.1.4 感知融合

| 功能 | 说明 |
|------|------|
| 多传感器融合 | 融合视觉/距离/触觉信息 |
| 感知增强 | 单一传感器失效时降级使用 |
| 感知置信度 | 给出感知结果的置信度 |

### 2.2 空间认知

#### 2.2.1 地图构建

| 功能 | 说明 |
|------|------|
| SLAM | 同步定位与地图构建 |
| 栅格地图 | 构建环境栅格地图 |
| 语义地图 | 添加物体语义的地图 |
| 地图更新 | 动态更新环境变化 |
| 多层地图 | 支持多层建筑地图 |

#### 2.2.2 定位

| 功能 | 说明 |
|------|------|
| 室内定位 | 室内精确定位（精度<10cm） |
| 朝向估计 | 估计设备朝向 |
| 地图匹配 | 将当前位置匹配到地图 |
| 多人定位 | 同时跟踪多人的位置 |

#### 2.2.3 空间理解

| 功能 | 说明 |
|------|------|
| 区域识别 | 识别房间/区域类型 |
| 可通行区域 | 分析可通行的路径 |
| 危险区域 | 标记楼梯/电器等危险区域 |
| 物体位置记忆 | 记住物体放置位置 |

### 2.3 自主探索

#### 2.3.1 导航

| 功能 | 说明 |
|------|------|
| 路径规划 | 全局路径规划 |
| 局部规划 | 局部路径动态调整 |
| 自主充电 | 没电时自主回充 |
| 定点移动 | 移动到指定位置 |
| 跟随模式 | 跟随指定人员移动 |

#### 2.3.2 避障

| 功能 | 说明 |
|------|------|
| 静态避障 | 避开静止障碍物 |
| 动态避障 | 避开移动障碍物 |
| 紧急制动 | 检测到危险时紧急停止 |
| 窄道通过 | 通过窄门/通道 |

#### 2.3.3 探索策略

| 功能 | 说明 |
|------|------|
| 全覆盖探索 | 探索所有可达区域 |
| 兴趣点探索 | 优先探索可能有趣的地方 |
| 高效探索 | 最少重复路径探索 |
| 自主学习 | 学习高效的探索策略 |

### 2.4 动作模仿

#### 2.4.1 动作学习

| 功能 | 说明 |
|------|------|
| 动作录制 | 录制人类示范动作 |
| 动作分割 | 将连续动作分割为基本动作 |
| 动作迁移 | 将人类动作转换为宠物动作 |
| 动作优化 | 优化动作使其更适合宠物执行 |
| 动作库 | 存储学习到的动作 |

#### 2.4.2 动作执行

| 功能 | 说明 |
|------|------|
| 动作播放 | 执行指定动作序列 |
| 动作参数调整 | 调整动作幅度/速度 |
| 动作融合 | 平滑连接多个动作 |
| 动作中断 | 接收新指令时中断当前动作 |

#### 2.4.3 动作库

| 功能 | 说明 |
|------|------|
| 内置动作 | 预置的基础动作集 |
| 自定义动作 | 用户录制的动作 |
| 动作分享 | 分享动作给其他用户 |
| 动作评分 | 动作质量评分 |

### 2.5 具身AI决策引擎

#### 2.5.1 感知-认知-决策闭环

| 阶段 | 功能 | 说明 |
|------|------|------|
| 感知 | 多模态感知融合 | 整合视觉/语音/触觉等感知 |
| 认知 | 场景理解和推理 | 理解当前场景和上下文 |
| 决策 | 行动规划 | 选择最优行动方案 |
| 执行 | 动作控制和反馈 | 执行动作并获取反馈 |

#### 2.5.2 决策策略

| 策略 | 说明 |
|------|------|
| 安全优先 | 任何情况下安全第一 |
| 任务导向 | 以完成任务为目标 |
| 交互优先 | 以用户交互体验为主 |
| 自主探索 | 以探索新环境为主 |

#### 2.5.3 决策上下文

| 上下文类型 | 说明 |
|------------|------|
| 用户状态 | 用户情绪/姿态/位置 |
| 环境状态 | 当前位置/时间/场景 |
| 宠物状态 | 宠物情绪/电量/模式 |
| 任务状态 | 当前任务进度 |

### 2.6 具身AI安全边界

#### 2.6.1 物理安全

| 功能 | 说明 |
|------|------|
| 碰撞防护 | 避免碰撞人或物 |
| 坠落防护 | 防止从高处坠落 |
| 夹手防护 | 防止夹伤用户 |
| 高温防护 | 远离高温物体 |
| 溺水防护 | 远离水源 |

#### 2.6.2 行为安全

| 功能 | 说明 |
|------|------|
| 动作速度限制 | 限制最大移动速度 |
| 力度限制 | 限制最大接触力度 |
| 禁区设置 | 用户可设置禁止进入区域 |
| 安全模式 | 紧急情况切换到安全模式 |

#### 2.6.3 应急处理

| 功能 | 说明 |
|------|------|
| 紧急停止 | 用户可一键紧急停止 |
| 异常恢复 | 异常状态自动恢复 |
| 远程控制 | 用户可远程接管控制 |
| 安全日志 | 记录所有安全相关事件 |

---

## 三、API接口定义

### 3.1 环境感知

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/embodied/:device_id/perception | 获取当前感知结果 |
| POST | /api/v1/embodied/:device_id/perception/visual | 上报视觉感知结果 |
| POST | /api/v1/embodied/:device_id/perception/depth | 上报深度感知结果 |
| POST | /api/v1/embodied/:device_id/perception/touch | 上报触觉感知结果 |

### 3.2 空间认知

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/embodied/:device_id/map | 获取地图 |
| POST | /api/v1/embodied/:device_id/map/update | 更新地图 |
| GET | /api/v1/embodied/:device_id/localization | 获取当前位置 |
| POST | /api/v1/embodied/:device_id/localization/calibrate | 定位校准 |

### 3.3 自主探索

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/embodied/:device_id/navigate | 导航到目标 |
| POST | /api/v1/embodied/:device_id/stop | 停止移动 |
| POST | /api/v1/embodied/:device_id/follow | 跟随目标 |
| GET | /api/v1/embodied/:device_id/explore/status | 探索状态 |
| POST | /api/v1/embodied/:device_id/explore/start | 开始探索 |

### 3.4 动作模仿

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/embodied/action-library | 动作库 |
| POST | /api/v1/embodied/action-library/record | 录制动作 |
| POST | /api/v1/embodied/action-library/:id/learn | 学习动作 |
| POST | /api/v1/embodied/:device_id/action/execute | 执行动作 |
| POST | /api/v1/embodied/:device_id/action/stop | 停止动作 |
| POST | /api/v1/embodied/action-library/:id/share | 分享动作 |

### 3.5 决策引擎

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/embodied/:device_id/decision/context | 获取决策上下文 |
| POST | /api/v1/embodied/:device_id/decision/strategy | 设置决策策略 |
| GET | /api/v1/embodied/:device_id/decision/logs | 决策日志 |

### 3.6 安全边界

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/embodied/:device_id/safety/zones | 获取禁区 |
| POST | /api/v1/embodied/:device_id/safety/zones | 设置禁区 |
| DELETE | /api/v1/embodied/:device_id/safety/zones/:id | 删除禁区 |
| POST | /api/v1/embodied/:device_id/safety/emergency-stop | 紧急停止 |
| GET | /api/v1/embodied/:device_id/safety/logs | 安全日志 |

---

## 四、数据库设计

### 4.1 环境地图表 (embodied_maps)

```sql
CREATE TABLE embodied_maps (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    map_type        VARCHAR(20) NOT NULL,         -- 'grid'/'semantic'/'topological'
    map_data        JSONB NOT NULL,               -- 地图数据
    resolution      DECIMAL(10,3),                -- 栅格分辨率
    size            JSONB,                        -- 地图尺寸
    version         INT DEFAULT 1,
    is_active       BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_embodied_maps_device ON embodied_maps(device_id, is_active);
```

### 4.2 空间位置表 (spatial_positions)

```sql
CREATE TABLE spatial_positions (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    map_id          BIGINT REFERENCES embodied_maps(id),
    position_x      DECIMAL(10,3),
    position_y      DECIMAL(10,3),
    position_z      DECIMAL(10,3),
    orientation     DECIMAL(10,3),
    confidence      DECIMAL(5,4),
    recorded_at     TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_spatial_positions_device ON spatial_positions(device_id, recorded_at DESC);
```

### 4.3 动作库表 (action_library)

```sql
CREATE TABLE action_library (
    id              BIGSERIAL PRIMARY KEY,
    action_name     VARCHAR(100) NOT NULL,
    action_type     VARCHAR(50) NOT NULL,         -- 'built-in'/'learned'/'custom'
    description     TEXT,
    duration_ms     INT,
    difficulty      VARCHAR(20),                   -- 'easy'/'medium'/'hard'
    tags            VARCHAR(100)[],
    motion_data     JSONB,                        -- 动作运动数据
    video_url       VARCHAR(500),
    thumbnail_url   VARCHAR(500),
    score           DECIMAL(5,2),
    creator_id      BIGINT REFERENCES users(id),
    is_public       BOOLEAN DEFAULT FALSE,
    downloads       INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.4 动作执行记录表 (action_executions)

```sql
CREATE TABLE action_executions (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    action_id       BIGINT REFERENCES action_library(id),
    execution_type  VARCHAR(20),                  -- 'triggered'/'scheduled'/'manual'
    start_time      TIMESTAMP NOT NULL,
    end_time        TIMESTAMP,
    status          VARCHAR(20),                  -- 'running'/'completed'/'interrupted'/'failed'
    parameters      JSONB,                        -- 执行参数
    interruption_reason TEXT,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_action_executions_device ON action_executions(device_id, start_time DESC);
```

### 4.5 禁区表 (safety_zones)

```sql
CREATE TABLE safety_zones (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    zone_type       VARCHAR(20) NOT NULL,         -- 'forbidden'/'caution'/'safe'
    zone_shape      VARCHAR(20) NOT NULL,         -- 'rectangle'/'circle'/'polygon'
    zone_data       JSONB NOT NULL,                -- 区域坐标数据
    zone_name       VARCHAR(100),
    is_enabled      BOOLEAN DEFAULT TRUE,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.6 具身决策日志表 (embodied_decision_logs)

```sql
CREATE TABLE embodied_decision_logs (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    decision_type   VARCHAR(50) NOT NULL,
    context         JSONB,                          -- 决策上下文
    chosen_action   VARCHAR(100),
    action_params   JSONB,
    confidence      DECIMAL(5,4),
    reasoning       TEXT,
    execution_result VARCHAR(20),
    latency_ms      INT,
    decided_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_embodied_decision_logs_device ON embodied_decision_logs(device_id, decided_at DESC);
```

### 4.7 安全日志表 (safety_logs)

```sql
CREATE TABLE safety_logs (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    event_type      VARCHAR(50) NOT NULL,         -- 'collision'/'emergency_stop'/'zone_violation'
    severity        VARCHAR(20) NOT NULL,
    details         JSONB,
    location        JSONB,
    resolved        BOOLEAN DEFAULT FALSE,
    resolved_at     TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_safety_logs_device ON safety_logs(device_id, created_at DESC);
```

---

## 五、前端页面清单

### 5.1 环境感知

| 页面 | 路由 | 说明 |
|------|------|------|
| 环境感知面板 | /embodied/:device_id/perception | 实时感知数据 |
| 视觉感知 | /embodied/:device_id/perception/visual | 视觉识别结果 |

### 5.2 地图管理

| 页面 | 路由 | 说明 |
|------|------|------|
| 地图查看 | /embodied/:device_id/map | 查看/编辑地图 |
| 定位校准 | /embodied/:device_id/localization | 定位校准工具 |

### 5.3 导航控制

| 页面 | 路由 | 说明 |
|------|------|------|
| 导航控制 | /embodied/:device_id/navigate | 导航控制面板 |
| 跟随模式 | /embodied/:device_id/follow | 跟随模式设置 |
| 自主探索 | /embodied/:device_id/explore | 探索状态和控制 |

### 5.4 动作库

| 页面 | 路由 | 说明 |
|------|------|------|
| 动作库 | /embodied/action-library | 动作库列表 |
| 动作录制 | /embodied/action-library/record | 录制新动作 |
| 动作详情 | /embodied/action-library/:id | 动作详情 |

### 5.5 安全设置

| 页面 | 路由 | 说明 |
|------|------|------|
| 禁区设置 | /embodied/:device_id/safety/zones | 禁区管理 |
| 安全日志 | /embodied/:device_id/safety/logs | 安全事件记录 |

### 5.6 决策日志

| 页面 | 路由 | 说明 |
|------|------|------|
| 决策日志 | /embodied/:device_id/decision/logs | AI决策记录 |
| 决策策略 | /embodied/:device_id/decision/strategy | 决策策略配置 |

---

## 六、验收标准

### 6.1 环境感知

| 验收点 | 标准 |
|--------|------|
| 物体识别准确率 | >90%（常见20类物体） |
| 深度估计误差 | <5%（相对距离） |
| 感知延迟 | <100ms |
| 感知融合 | 多传感器数据一致 |

### 6.2 空间认知

| 验收点 | 标准 |
|--------|------|
| 定位精度 | <10cm |
| 地图精度 | 障碍物位置偏差<15cm |
| 地图更新延迟 | 环境变化后<5分钟更新 |

### 6.3 自主探索

| 验收点 | 标准 |
|--------|------|
| 导航成功率 | >90% |
| 避障成功率 | >95% |
| 自主回充成功率 | >95% |
| 跟随稳定性 | 跟随目标稳定不丢失 |

### 6.4 动作模仿

| 验收点 | 标准 |
|--------|------|
| 动作学习成功率 | >80% |
| 动作执行准确率 | >85% |
| 动作过渡平滑 | 无明显卡顿 |
| 动作分享可用 | 动作可被其他用户使用 |

### 6.5 安全

| 验收点 | 标准 |
|--------|------|
| 碰撞次数 | 0次（严重碰撞） |
| 紧急停止响应 | <100ms |
| 禁区闯入 | 0次 |
| 安全日志完整性 | 100%记录 |


---

## 七、页面布局规范

### 7.1 环境感知面板（/embodied/:device_id/perception）

**布局结构：**
1. 面包屑 → 页面标题 + 设备信息
2. 实时感知数据展示区（白色卡片）：视觉/距离/触觉
3. 感知事件记录表格

**分页：** 右下角，10/20/50/100 条

### 7.2 地图管理页面（/embodied/:device_id/map）

**布局结构：**
1. 面包屑 → 页面标题 + 设备选择
2. 地图可视化区域（白色，全屏卡片）
3. 地图工具栏（缩放/定位/编辑模式切换）
4. 物体标注列表（右侧边栏）

**按钮规范：**
- [定位校准] [保存地图] — 左对齐
- [编辑] [删除] — 物体标注内右对齐

### 7.3 导航控制页面（/embodied/:device_id/navigate）

**布局结构：**
1. 面包屑 → 页面标题
2. 地图区域（白色卡片）+ 导航控制面板
3. 导航历史记录表格

**按钮规范：**
- [开始导航] [停止] [跟随模式] — 左对齐

**分页：** 右下角，10/20/50/100 条

### 7.4 动作库页面（/embodied/action-library）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片 #F2F3F5）：动作类型 / 难度 / 来源
3. 操作栏（录制动作靠左）
4. 动作库列表表格

**按钮规范：**
- [录制动作] — 左对齐
- [预览] [编辑] [分享] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 动作名称 | 150px | - |
| 类型 | 100px | 内置/学习/自定义 |
| 难度 | 100px | easy/medium/hard |
| 时长 | 80px | - |
| 评分 | 80px | - |
| 下载量 | 100px | - |
| 操作 | 120px | 预览/编辑/分享/删除 |

**分页：** 右下角，10/20/50/100 条

### 7.5 禁区设置页面（/embodied/:device_id/safety/zones）

**布局结构：**
1. 面包屑 → 页面标题 + 设备信息
2. 地图区域（白色卡片，显示禁区）
3. 禁区列表表格

**按钮规范：**
- [添加禁区] — 左对齐
- [编辑] [删除] — 行内右对齐

**分页：** 右下角，10/20/50/100 条

### 7.6 决策日志页面（/embodied/:device_id/decision/logs）

**布局结构：**
1. 面包屑 → 页面标题 + 设备信息
2. 筛选区（浅灰卡片）：决策类型 / 日期范围
3. 决策日志列表表格

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 决策类型 | 120px | - |
| 决策上下文 | 200px | - |
| 选择动作 | 150px | - |
| 置信度 | 100px | - |
| 延迟ms | 100px | - |
| 决策时间 | 150px | - |

**分页：** 右下角，10/20/50/100 条

### 7.7 弹窗规范

| 类型 | 使用场景 |
|------|----------|
| Drawer 抽屉 | 动作详情、禁区编辑、决策详情 |
| Dialog 对话框 | 确认删除禁区、紧急停止确认 |
| 全屏模态 | 暂无复杂表单场景 |
