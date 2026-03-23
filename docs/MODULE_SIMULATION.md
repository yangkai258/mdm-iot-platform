# 模块 PRD：仿真与测试（MODULE_SIMULATION）

**版本：** V1.1
**所属Phase：** Phase 3（Sprint 23-24）
**优先级：** P1
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

仿真与测试模块为AI电子宠物提供完整的虚拟仿真环境和自动化测试能力，支持虚拟宠物运行、自动化测试用例管理、设备行为回放、场景仿真、数据集管理，以及CI/CD集成，大幅降低真机测试成本，加速产品迭代。

### 1.2 核心价值

- **成本优化**：减少真机测试，缩短迭代周期
- **质量保障**：自动化回归测试，覆盖率持续提升
- **问题定位**：回放系统快速复现问题
- **持续集成**：CI/CD流水线无缝对接

### 1.3 范围边界

**包含：**
- 虚拟宠物仿真
- 自动化测试框架
- 回放系统
- 仿真场景管理
- 压力测试
- 仿真数据集管理
- CI/CD 集成
- A/B实验仿真

**不包含：**
- 具身智能（MODULE_EMBODIED_AI）
- 情感计算（MODULE_AFFECTIVE_COMPUTING）

---

## 二、系统架构

### 2.1 技术架构图

```
┌────────────────────────────────────────────────────────────────────────┐
│                        前端 (Vue 3 + Arco Design)                      │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│  │ 虚拟宠物  │ │ 测试用例  │ │ 回放系统  │ │ 场景管理  │ │ 压力测试  │ │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘ │
└───────┼────────────┼────────────┼────────────┼────────────┼─────────┘
        │            │            │            │            │
        └────────────┴────────────┼────────────┴────────────┘
                                 │ REST API / WebSocket
┌────────────────────────────────┼─────────────────────────────────────┐
│                          后端服务 (Go + Gin)                            │
│  ┌────────────────────────────────────────────────────────────────┐  │
│  │                   Simulation Gateway Service                     │  │
│  ├───┬──────────┬──────────┬──────────┬──────────┬──────────┐    │  │
│  │Pet│ TestCase │ Playback │ Scenario │StressTest│ Dataset │    │  │
│  │API│   API    │   API    │   API    │   API    │   API   │    │  │
│  └───┴──────────┴──────────┴──────────┴──────────┴──────────┘    │  │
│  ┌────────────────────────────────────────────────────────────────┐  │  │
│  │                  Test Execution Engine                          │  │  │
│  │         (Goroutine Pool / WebSocket Push / Report Gen)         │  │  │
│  └────────────────────────────────────────────────────────────────┘  │
│  ┌────────────────────────────────────────────────────────────────┐  │
│  │              Stress Test Engine (自研 k6-style)                 │  │
│  └────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────┘
        │                    │                    │
        ▼                    ▼                    ▼
┌───────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  PostgreSQL   │  │      Redis      │  │   MQTT Broker   │
│ (主数据存储)   │  │ (缓存/队列/会话) │  │   (EMQX 5.0)    │
└───────────────┘  └─────────────────┘  └─────────────────┘
```

### 2.2 业务流程图

#### 2.2.1 自动化测试执行流程

```
用户创建用例
     │
     ▼
┌─────────────┐    否     ┌─────────────────┐
│ 用例校验    │──────────▶│ 返回校验错误     │
└──────┬──────┘           └─────────────────┘
       │ 是
       ▼
┌─────────────┐
│ 用例存储    │
└──────┬──────┘
       ▼
┌──────────────────────────────────────────────────────┐
│                    测试执行引擎                       │
│  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐ │
│  │Goroutine│  │Goroutine│  │Goroutine│  │  ...   │ │
│  │ Pool 1  │  │ Pool 2  │  │ Pool 3  │  │        │ │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬───┘ │
│       └────────────┴──────┬──────┴─────────────┘      │
│                           ▼                             │
│                    ┌────────────┐                       │
│                    │  结果收集   │                       │
│                    └──────┬─────┘                       │
└───────────────────────────┼──────────────────────────────┘
                            ▼
                     ┌────────────┐
                     │  报告生成  │
                     └──────┬─────┘
                            ▼
                     ┌────────────┐
                     │  结果推送  │
                     │(WebSocket) │
                     └────────────┘
```

#### 2.2.2 压力测试流程

```
创建压力测试配置
     │
     ▼
┌─────────────┐    否     ┌─────────────────┐
│ 配置校验    │──────────▶│ 返回校验错误     │
└──────┬──────┘           └─────────────────┘
       │ 是
       ▼
┌─────────────┐
│ 配置存储    │
└──────┬──────┘
       ▼
┌──────────────────────────────────────────────────────┐
│                    压力测试引擎                        │
│  ┌────────────────────────────────────────────────┐  │
│  │            Worker Pool (N goroutines)           │  │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ │  │
│  │  │Worker 1│ │Worker 2│ │Worker 3│ │  ...   │ │  │
│  │  └────┬───┘ └────┬───┘ └────┬───┘ └────┬──┘ │  │
│  └───────┼──────────┼──────────┼──────────┼─────┘  │
│          └──────────┴────┬─────┴──────────┘         │
│                          ▼                           │
│                   ┌──────────────┐                   │
│                   │ Metrics Agg  │                   │
│                   │(Prometheus)  │                   │
│                   └──────────────┘                   │
└──────────────────────────────────────────────────────┘
          │
          ▼
   ┌──────────────┐
   │  报告生成   │
   └──────────────┘
```

#### 2.2.3 回放系统流程

```
设备行为录制
     │
     ▼
┌─────────────┐
│  数据采集   │ ◀── 传感器数据 / MQTT消息 / 用户操作
└──────┬──────┘
       ▼
┌─────────────┐
│  数据压缩   │
└──────┬──────┘
       ▼
┌─────────────┐
│  存储回放   │
└──────┬──────┘
       ▼
┌─────────────┐
│  回放请求   │
└──────┬──────┘
       ▼
┌─────────────┐
│  精准回放   │ ◀── 时间同步 / 变量替换 / 断点控制
└──────┬──────┘
       ▼
┌─────────────┐
│  差异对比   │
└─────────────┘
```

#### 2.2.4 CI/CD 集成流程

```
代码提交
     │
     ▼
┌─────────────┐
│  代码构建   │
└──────┬──────┘
       ▼
   单元测试 ──失败──▶ 通知并阻止合并
       │
      成功
       │
       ▼
┌─────────────┐
│ 自动化测试 │ ──失败──▶ 通知并阻止发布
└──────┬──────┘
       │
      成功
       │
       ▼
┌─────────────┐
│  集成测试   │
└──────┬──────┘
       ▼
┌─────────────┐
│预发布环境   │
└──────┬──────┘
       ▼
┌─────────────┐
│  正式发布   │
└─────────────┘
```

---

## 三、功能详情

### 2.1 虚拟宠物仿真

#### 2.1.1 虚拟宠物环境

| 功能 | 说明 |
|------|------|
| 虚拟宠物渲染 | 3D虚拟宠物形象展示 |
| 环境模拟 | 模拟不同家居环境 |
| 物理模拟 | 虚拟宠物的物理运动 |
| 多宠物支持 | 同时运行多个虚拟宠物 |
| 设备模拟 | 模拟设备传感器数据 |

#### 2.1.2 虚拟交互

| 功能 | 说明 |
|------|------|
| 触摸交互 | 模拟触摸虚拟宠物 |
| 语音交互 | 模拟语音指令和对话 |
| 手势交互 | 模拟手势指令 |
| 环境交互 | 虚拟宠物与虚拟环境交互 |

#### 2.1.3 仿真配置

| 功能 | 说明 |
|------|------|
| 宠物类型 | 支持多种宠物类型（猫/狗/兔子等） |
| 性格配置 | 配置虚拟宠物的性格 |
| 能力配置 | 配置虚拟宠物的AI能力 |
| 环境配置 | 配置虚拟环境参数 |

### 2.2 自动化测试框架

#### 2.2.1 测试用例管理

| 功能 | 说明 |
|------|------|
| 用例创建 | 创建测试用例 |
| 用例分类 | 按模块/功能/优先级分类 |
| 用例版本 | 用例版本管理 |
| 用例依赖 | 定义用例间依赖关系 |
| 用例标签 | 用例标签管理 |

#### 2.2.2 测试执行

| 功能 | 说明 |
|------|------|
| 单用例执行 | 执行单个测试用例 |
| 批量执行 | 批量执行多个用例 |
| 定时执行 | 配置定时执行任务 |
| 并发执行 | 支持多并发执行 |
| 失败重试 | 失败用例自动重试 |

#### 2.2.3 测试报告

| 功能 | 说明 |
|------|------|
| 实时报告 | 执行过程中实时展示结果 |
| 详细报告 | 包含步骤截图/日志 |
| 趋势分析 | 通过率趋势图 |
| 失败分析 | 失败原因分析 |
| 覆盖率统计 | 代码/功能覆盖率统计 |

#### 2.2.4 测试用例类型

| 类型 | 说明 | 示例 |
|------|------|------|
| 功能测试 | 验证功能正确性 | 情绪识别准确率测试 |
| 性能测试 | 验证性能指标 | 响应延迟测试 |
| 压力测试 | 验证极限能力 | 高并发测试 |
| 回归测试 | 验证修改未破坏功能 | 全量回归 |
| 冒烟测试 | 验证核心功能可用 | 快速冒烟 |

### 2.3 回放系统

#### 2.3.1 行为录制

| 功能 | 说明 |
|------|------|
| 自动录制 | 设备行为自动录制 |
| 手动标记 | 手动标记关键事件 |
| 传感器数据 | 录制完整传感器数据 |
| 用户操作 | 录制用户操作序列 |

#### 2.3.2 行为回放

| 功能 | 说明 |
|------|------|
| 精准回放 | 完整复现历史行为 |
| 变量替换 | 回放时替换变量 |
| 断点设置 | 回放中设置断点 |
| 速度控制 | 调整回放速度 |
| 部分回放 | 从指定时间点回放 |

#### 2.3.3 回放分析

| 功能 | 说明 |
|------|------|
| 差异对比 | 回放结果与录制对比 |
| 性能分析 | 回放过程性能分析 |
| 问题定位 | 快速定位问题原因 |
| 分享功能 | 分享回放给开发人员 |

### 2.4 仿真场景管理

#### 2.4.1 场景库

| 功能 | 说明 |
|------|------|
| 预置场景 | 预置常见测试场景 |
| 自定义场景 | 用户创建自定义场景 |
| 场景导入导出 | 场景的导入导出 |
| 场景评分 | 用户对场景评分 |

#### 2.4.2 场景编辑器

| 功能 | 说明 |
|------|------|
| 环境编辑 | 编辑虚拟环境 |
| 物体摆放 | 摆放物体和障碍物 |
| 灯光配置 | 配置光照条件 |
| 声音配置 | 配置背景声音 |
| 事件配置 | 配置触发事件 |

#### 2.4.3 预置场景

| 场景名称 | 说明 |
|----------|------|
| 客厅日常 | 模拟客厅日常活动 |
| 餐厅用餐 | 模拟用餐场景 |
| 夜间睡眠 | 模拟夜间睡眠 |
| 户外散步 | 模拟户外活动 |
| 多人互动 | 模拟多人家庭成员 |
| 紧急情况 | 模拟紧急事件 |

### 2.5 压力测试

#### 2.5.1 并发测试

| 功能 | 说明 |
|------|------|
| 设备并发 | 模拟多设备并发 |
| 用户并发 | 模拟多用户并发 |
| 请求并发 | 模拟API高并发 |
| 逐步加压 | 逐步增加并发量 |

#### 2.5.2 性能测试

| 功能 | 说明 |
|------|------|
| 响应时间 | 测试各接口响应时间 |
| 吞吐量 | 测试系统吞吐量 |
| 资源使用 | 监控CPU/内存/网络 |
| 性能瓶颈 | 定位性能瓶颈 |

#### 2.5.3 稳定性测试

| 功能 | 说明 |
|------|------|
| 长时间运行 | 7x24小时连续运行 |
| 内存泄漏 | 检测内存泄漏 |
| 连接泄漏 | 检测连接泄漏 |
| 故障恢复 | 测试故障自动恢复 |

### 2.6 A/B实验仿真

#### 2.6.1 实验预验证

| 功能 | 说明 |
|------|------|
| 参数仿真 | 仿真不同参数效果 |
| 流量模拟 | 模拟不同流量场景 |
| 效果预估 | 预估实验效果 |
| 风险评估 | 评估实验风险 |

#### 2.6.2 实验对比

| 功能 | 说明 |
|------|------|
| 多方案对比 | 同时对比多个方案 |
| 指标对比 | 对比各项指标 |
| 成本对比 | 对比资源消耗 |
| 推荐建议 | 推荐最优方案 |

---

## 三、API接口定义

### 3.1 虚拟宠物

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/simulation/pets | 创建虚拟宠物 |
| GET | /api/v1/simulation/pets/:id | 获取虚拟宠物 |
| PUT | /api/v1/simulation/pets/:id | 更新虚拟宠物配置 |
| DELETE | /api/v1/simulation/pets/:id | 删除虚拟宠物 |
| POST | /api/v1/simulation/pets/:id/interact | 虚拟交互 |

### 3.2 自动化测试

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/simulation/testcases | 测试用例列表 |
| POST | /api/v1/simulation/testcases | 创建测试用例 |
| GET | /api/v1/simulation/testcases/:id | 用例详情 |
| PUT | /api/v1/simulation/testcases/:id | 更新用例 |
| DELETE | /api/v1/simulation/testcases/:id | 删除用例 |
| POST | /api/v1/simulation/testcases/:id/execute | 执行用例 |
| POST | /api/v1/simulation/testcases/batch-execute | 批量执行 |
| GET | /api/v1/simulation/reports/:id | 测试报告 |
| GET | /api/v1/simulation/reports | 报告列表 |

### 3.3 回放系统

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/simulation/playbacks | 回放列表 |
| POST | /api/v1/simulation/playbacks | 创建回放（录制） |
| GET | /api/v1/simulation/playbacks/:id | 回放详情 |
| POST | /api/v1/simulation/playbacks/:id/play | 开始回放 |
| POST | /api/v1/simulation/playbacks/:id/stop | 停止回放 |
| POST | /api/v1/simulation/playbacks/:id/compare | 差异对比 |
| DELETE | /api/v1/simulation/playbacks/:id | 删除回放 |

### 3.4 仿真场景

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/simulation/scenarios | 场景列表 |
| POST | /api/v1/simulation/scenarios | 创建场景 |
| GET | /api/v1/simulation/scenarios/:id | 场景详情 |
| PUT | /api/v1/simulation/scenarios/:id | 更新场景 |
| DELETE | /api/v1/simulation/scenarios/:id | 删除场景 |
| POST | /api/v1/simulation/scenarios/:id/run | 运行场景 |
| POST | /api/v1/simulation/scenarios/import | 导入场景 |
| POST | /api/v1/simulation/scenarios/export/:id | 导出场景 |

### 3.5 仿真场景管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/simulation/scenes | 场景列表（轻量） |
| POST | /api/v1/simulation/scenes | 创建场景 |
| GET | /api/v1/simulation/scenes/:id | 场景详情 |
| PUT | /api/v1/simulation/scenes/:id | 更新场景 |
| DELETE | /api/v1/simulation/scenes/:id | 删除场景 |
| POST | /api/v1/simulation/scenes/:id/clone | 克隆场景 |
| GET | /api/v1/simulation/scenes/:id/preview | 场景预览 |
| POST | /api/v1/simulation/scenes/batch-delete | 批量删除 |
| POST | /api/v1/simulation/scenes/import | 导入场景 |
| GET | /api/v1/simulation/scenes/export/:id | 导出场景 |

### 3.6 A/B实验仿真

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/simulation/experiments | 实验列表 |
| POST | /api/v1/simulation/experiments | 创建实验 |
| GET | /api/v1/simulation/experiments/:id | 实验详情 |
| PUT | /api/v1/simulation/experiments/:id | 更新实验 |
| DELETE | /api/v1/simulation/experiments/:id | 删除实验 |
| POST | /api/v1/simulation/experiments/:id/start | 开始实验仿真 |
| POST | /api/v1/simulation/experiments/:id/stop | 停止实验 |
| GET | /api/v1/simulation/experiments/:id/compare | 多方案对比 |
| GET | /api/v1/simulation/experiments/:id/recommend | 推荐最优方案 |

### 3.7 仿真数据集管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/simulation/datasets | 数据集列表 |
| POST | /api/v1/simulation/datasets | 创建数据集 |
| GET | /api/v1/simulation/datasets/:id | 数据集详情 |
| PUT | /api/v1/simulation/datasets/:id | 更新数据集 |
| DELETE | /api/v1/simulation/datasets/:id | 删除数据集 |
| POST | /api/v1/simulation/datasets/:id/upload | 上传数据 |
| GET | /api/v1/simulation/datasets/:id/download | 下载数据集 |
| POST | /api/v1/simulation/datasets/:id/split | 数据集划分 |
| GET | /api/v1/simulation/datasets/:id/samples | 数据样本列表 |

### 3.8 压力测试

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/simulation/stress-tests | 创建压力测试 |
| GET | /api/v1/simulation/stress-tests | 压力测试列表 |
| GET | /api/v1/simulation/stress-tests/:id | 测试详情 |
| PUT | /api/v1/simulation/stress-tests/:id | 更新测试配置 |
| DELETE | /api/v1/simulation/stress-tests/:id | 删除测试 |
| POST | /api/v1/simulation/stress-tests/:id/start | 开始测试 |
| POST | /api/v1/simulation/stress-tests/:id/stop | 停止测试 |
| GET | /api/v1/simulation/stress-tests/:id/status | 测试状态 |
| GET | /api/v1/simulation/stress-tests/:id/report | 测试报告 |
| WS | /api/v1/simulation/stress-tests/:id/metrics | 实时指标流 |

### 3.6 仿真数据集

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/simulation/datasets | 创建数据集 |
| GET | /api/v1/simulation/datasets | 数据集列表 |
| GET | /api/v1/simulation/datasets/:id | 数据集详情 |
| PUT | /api/v1/simulation/datasets/:id | 更新数据集 |
| DELETE | /api/v1/simulation/datasets/:id | 删除数据集 |
| POST | /api/v1/simulation/datasets/:id/versions | 创建版本 |
| GET | /api/v1/simulation/datasets/:id/versions | 版本列表 |
| GET | /api/v1/simulation/datasets/:id/versions/:version | 版本详情 |
| POST | /api/v1/simulation/datasets/:id/versions/:version/publish | 发布版本 |
| POST | /api/v1/simulation/datasets/:id/versions/:version/rollback | 回滚版本 |
| POST | /api/v1/simulation/datasets/:id/versions/compare | 对比版本 |
| POST | /api/v1/simulation/datasets/import | 导入数据集 |
| POST | /api/v1/simulation/datasets/export/:id | 导出数据集 |
| GET | /api/v1/simulation/datasets/jobs/:job_id | 导入导出状态 |

### 3.7 CI/CD 集成

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/simulation/webhooks/trigger | Webhook触发 |
| POST | /api/v1/simulation/webhooks/callback/:job_id | CI/CD回调 |
| GET | /api/v1/simulation/integrations | 集成列表 |
| POST | /api/v1/simulation/integrations | 创建集成 |
| PUT | /api/v1/simulation/integrations/:id | 更新集成 |
| DELETE | /api/v1/simulation/integrations/:id | 删除集成 |

---

### 3.8 API详细规格

#### 3.8.1 虚拟宠物仿真 API

##### 3.8.1.1 创建虚拟宠物

**请求**
```http
POST /api/v1/simulation/pets
Content-Type: application/json
```

```json
{
  "pet_name": "小白",
  "pet_type": "cat",
  "personality": {
    "friendliness": 0.8,
    "playfulness": 0.6,
    "energy_level": 0.7
  },
  "capabilities": {
    "speech_recognition": true,
    "emotion_recognition": true,
    "path_planning": true
  },
  "environment_id": 1
}
```

**响应 (201 Created)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "pet_name": "小白",
    "pet_type": "cat",
    "personality": {"friendliness": 0.8, "playfulness": 0.6, "energy_level": 0.7},
    "capabilities": {"speech_recognition": true, "emotion_recognition": true, "path_planning": true},
    "environment_id": 1,
    "status": "idle",
    "current_emotion": null,
    "current_position": null,
    "battery_level": 1.0,
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:00:00Z"
  }
}
```

##### 3.8.1.2 获取虚拟宠物

**请求**
```http
GET /api/v1/simulation/pets/:id
```

**响应 (200 OK)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "pet_name": "小白",
    "pet_type": "cat",
    "personality": {"friendliness": 0.8, "playfulness": 0.6, "energy_level": 0.7},
    "capabilities": {"speech_recognition": true, "emotion_recognition": true, "path_planning": true},
    "environment_id": 1,
    "status": "running",
    "current_emotion": "happy",
    "current_position": {"x": 100, "y": 200, "z": 0},
    "battery_level": 0.85,
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:30:00Z"
  }
}
```

##### 3.8.1.3 虚拟宠物列表

**请求**
```http
GET /api/v1/simulation/pets?page=1&page_size=20&pet_type=cat&status=running
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页条数，默认20，最大100 |
| pet_type | string | 否 | 宠物类型：cat/dog/rabbit |
| status | string | 否 | 状态：idle/running/paused |

**响应 (200 OK)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "pet_name": "小白",
        "pet_type": "cat",
        "status": "running",
        "current_emotion": "happy",
        "created_at": "2026-03-23T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

##### 3.8.1.4 虚拟交互

**请求**
```http
POST /api/v1/simulation/pets/:id/interact
Content-Type: application/json
```

```json
{
  "interaction_type": "touch",
  "position": {"x": 100, "y": 200},
  "parameters": {
    "force": 0.5,
    "duration_ms": 500
  }
}
```

**交互类型**

| 类型 | 说明 |
|------|------|
| touch | 触摸交互 |
| voice | 语音指令 |
| gesture | 手势指令 |
| environmental | 环境交互 |

**响应 (200 OK)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "pet_response": {
      "emotion": "happy",
      "action": "purr",
      "animation": "playful_jump",
      "sound": "meow_soft"
    },
    "response_time_ms": 150
  }
}
```

##### 3.8.1.5 虚拟宠物状态流 (WebSocket)

**连接**
```http
WS /api/v1/simulation/pets/:id/stream
```

**服务端推送**
```json
{
  "event": "pet_status",
  "data": {
    "id": 1,
    "status": "running",
    "emotion": "happy",
    "position": {"x": 105, "y": 202, "z": 0},
    "battery_level": 0.84,
    "timestamp": 1711172400000
  }
}
```

**错误事件**
```json
{
  "event": "error",
  "data": {
    "code": "PET_NOT_RUNNING",
    "message": "宠物当前未运行"
  }
}
```

#### 3.8.2 自动化测试框架 API

##### 3.8.2.1 创建测试用例

**请求**
```http
POST /api/v1/simulation/testcases
Content-Type: application/json
```

```json
{
  "case_name": "情绪识别功能测试",
  "case_type": "functional",
  "module": "emotion_recognition",
  "priority": "high",
  "description": "测试宠物情绪识别功能的准确性",
  "preconditions": "虚拟宠物已创建并处于运行状态",
  "test_steps": [
    {"step": 1, "action": "触发高兴情绪", "expected": "情绪识别结果为happy", "timeout_ms": 5000},
    {"step": 2, "action": "触发悲伤情绪", "expected": "情绪识别结果为sad", "timeout_ms": 5000}
  ],
  "expected_result": "所有测试步骤通过",
  "tags": ["情绪", "功能测试"],
  "dependencies": []
}
```

**用例类型**

| 类型 | 说明 |
|------|------|
| functional | 功能测试 |
| performance | 性能测试 |
| stress | 压力测试 |
| regression | 回归测试 |
| smoke | 冒烟测试 |

**优先级**

| 优先级 | 说明 |
|--------|------|
| high | 高优先级 |
| medium | 中优先级 |
| low | 低优先级 |

**响应 (201 Created)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "case_name": "情绪识别功能测试",
    "case_type": "functional",
    "module": "emotion_recognition",
    "priority": "high",
    "status": "draft",
    "version": 1,
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:00:00Z"
  }
}
```

##### 3.8.2.2 测试用例列表

**请求**
```http
GET /api/v1/simulation/testcases?page=1&page_size=20&case_type=functional&module=emotion_recognition&priority=high
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页条数，默认20 |
| case_type | string | 否 | 用例类型 |
| module | string | 否 | 模块名 |
| priority | string | 否 | 优先级 |
| tags | string | 否 | 标签，逗号分隔 |
| status | string | 否 | 状态 |

**响应 (200 OK)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "case_name": "情绪识别功能测试",
        "case_type": "functional",
        "module": "emotion_recognition",
        "priority": "high",
        "status": "draft",
        "version": 1,
        "created_at": "2026-03-23T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

##### 3.8.2.3 执行单个测试用例

**请求**
```http
POST /api/v1/simulation/testcases/:id/execute
Content-Type: application/json
```

```json
{
  "environment": {
    "pet_id": 1,
    "env_id": 1
  },
  "parameters": {
    "timeout_ms": 30000,
    "retry_count": 2
  }
}
```

**响应 (202 Accepted)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "execution_id": "exec_20260323100001",
    "status": "pending"
  }
}
```

##### 3.8.2.4 批量执行测试用例

**请求**
```http
POST /api/v1/simulation/testcases/batch-execute
Content-Type: application/json
```

```json
{
  "case_ids": [1, 2, 3, 4, 5],
  "environment": {"pet_id": 1},
  "parameters": {
    "concurrency": 5,
    "timeout_ms": 60000,
    "stop_on_failure": false
  }
}
```

**响应 (202 Accepted)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "batch_id": "batch_20260323100001",
    "total_count": 5,
    "pending_count": 5,
    "status": "pending"
  }
}
```

##### 3.8.2.5 获取执行详情

**请求**
```http
GET /api/v1/simulation/executions/:id
```

**响应 (200 OK)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "exec_20260323100001",
    "testcase_id": 1,
    "status": "running",
    "progress": 60,
    "current_step": 2,
    "start_time": "2026-03-23T10:00:00Z",
    "result_details": {
      "steps_completed": 2,
      "steps_passed": 2,
      "steps_failed": 0,
      "logs": ["step 1 passed", "step 2 running"]
    }
  }
}
```

**执行状态**

| 状态 | 说明 |
|------|------|
| pending | 等待执行 |
| running | 执行中 |
| passed | 全部通过 |
| failed | 失败 |
| skipped | 跳过 |
| cancelled | 取消 |

##### 3.8.2.6 执行结果流 (WebSocket)

**连接**
```http
WS /api/v1/simulation/executions/:id/stream
```

**服务端推送（执行中）**
```json
{
  "event": "execution_update",
  "data": {
    "execution_id": "exec_20260323100001",
    "status": "running",
    "progress": 80,
    "step_result": {
      "step": 3,
      "action": "触发悲伤情绪",
      "expected": "情绪识别结果为sad",
      "actual": "情绪识别结果为sad",
      "status": "passed",
      "duration_ms": 120
    }
  }
}
```

**服务端推送（完成）**
```json
{
  "event": "execution_complete",
  "data": {
    "execution_id": "exec_20260323100001",
    "status": "passed",
    "duration_ms": 5200,
    "steps_passed": 3,
    "steps_failed": 0
  }
}
```

#### 3.8.3 回放系统 API

##### 3.8.3.1 创建回放（开始录制）

**请求**
```http
POST /api/v1/simulation/playbacks
Content-Type: application/json
```

```json
{
  "device_id": "device_001",
  "pet_id": 1,
  "record_type": "auto",
  "metadata": {
    "scenario": "客厅日常",
    "notes": "测试宠物自主行为"
  }
}
```

**录制类型**

| 类型 | 说明 |
|------|------|
| auto | 自动录制 |
| manual | 手动录制 |

**响应 (201 Created)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "status": "recording",
    "start_time": "2026-03-23T10:00:00Z"
  }
}
```

##### 3.8.3.2 开始回放

**请求**
```http
POST /api/v1/simulation/playbacks/:id/play
Content-Type: application/json
```

```json
{
  "start_position_ms": 0,
  "speed": 1.0,
  "variables": {
    "device_id": "device_002"
  }
}
```

**响应 (200 OK)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "playback_id": 1,
    "status": "playing",
    "current_position_ms": 0,
    "speed": 1.0
  }
}
```

##### 3.8.3.3 差异对比

**请求**
```http
POST /api/v1/simulation/playbacks/:id/compare
Content-Type: application/json
```

```json
{
  "compare_playback_id": 2,
  "metrics": ["response_time", "accuracy"]
}
```

**响应 (200 OK)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "metrics": {
      "response_time": {
        "playback_1_avg": 120,
        "playback_2_avg": 115,
        "difference_pct": -4.2
      },
      "accuracy": {
        "playback_1": 0.95,
        "playback_2": 0.97,
        "difference_pct": 2.1
      }
    },
    "significant_differences": []
  }
}
```

#### 3.8.4 仿真场景管理 API

##### 3.8.4.1 创建场景

**请求**
```http
POST /api/v1/simulation/scenarios
Content-Type: application/json
```

```json
{
  "scenario_name": "客厅日常场景",
  "scenario_type": "custom",
  "environment": {
    "env_id": 1,
    "objects": [
      {"type": "sofa", "position": {"x": 0, "y": 0}}
    ],
    "lighting": {"ambient": 0.5}
  },
  "events": [
    {
      "type": "pet_action",
      "trigger": "time",
      "params": {"delay_seconds": 60},
      "action": "walk_to",
      "target": {"x": 100, "y": 50}
    }
  ],
  "config": {
    "duration_minutes": 30,
    "ai_enabled": true
  },
  "is_public": false,
  "tags": ["日常", "客厅"]
}
```

**响应 (201 Created)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "scenario_name": "客厅日常场景",
    "scenario_type": "custom",
    "is_public": false,
    "created_at": "2026-03-23T10:00:00Z"
  }
}
```

##### 3.8.4.2 场景列表

**请求**
```http
GET /api/v1/simulation/scenarios?page=1&page_size=20&scenario_type=preset&is_public=true
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页条数，默认20 |
| scenario_type | string | 否 | 场景类型：preset/custom |
| is_public | bool | 否 | 是否公开 |
| tags | string | 否 | 标签 |

##### 3.8.4.3 运行场景

**请求**
```http
POST /api/v1/simulation/scenarios/:id/run
Content-Type: application/json
```

```json
{
  "pet_id": 1,
  "parameters": {
    "speed": 1.0,
    "record_enabled": true
  }
}
```

**响应 (202 Accepted)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "run_id": "run_20260323100001",
    "status": "running",
    "started_at": "2026-03-23T10:00:00Z"
  }
}
```

##### 3.8.4.4 导入场景

**请求**
```http
POST /api/v1/simulation/scenarios/import
Content-Type: multipart/form-data
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 场景文件（JSON格式） |
| overwrite | bool | 否 | 是否覆盖同名场景 |

**响应 (200 OK)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "imported_count": 3,
    "skipped_count": 0,
    "errors": []
  }
}
```

#### 3.8.5 压力测试 API

##### 3.8.5.1 创建压力测试

**请求**
```http
POST /api/v1/simulation/stress-tests
Content-Type: application/json
```

```json
{
  "test_name": "API并发压力测试",
  "test_type": "concurrent",
  "config": {
    "target": {
      "endpoint": "/api/v1/devices/status",
      "method": "POST",
      "headers": {"Content-Type": "application/json"},
      "body_template": {"device_id": "{{device_id}}", "status": "online"}
    },
    "load_pattern": {
      "type": "ramp",
      "initial_vus": 10,
      "max_vus": 100,
      "ramp_up_duration": "2m",
      "hold_duration": "5m"
    },
    "thresholds": {
      "http_req_duration": {"p50": 100, "p95": 200, "p99": 500},
      "http_req_failed": {"rate": 0.01}
    }
  }
}
```

**测试类型**

| 类型 | 说明 |
|------|------|
| concurrent | 并发测试 |
| performance | 性能测试 |
| stability | 稳定性测试 |

**负载模式**

| 模式 | 说明 |
|------|------|
| constant | 恒定并发 |
| ramp | 逐步加压 |
| spike | 峰值测试 |
| soak | 长期稳压 |

**响应 (201 Created)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "test_name": "API并发压力测试",
    "test_type": "concurrent",
    "status": "draft",
    "created_at": "2026-03-23T10:00:00Z"
  }
}
```

##### 3.8.5.2 开始压力测试

**请求**
```http
POST /api/v1/simulation/stress-tests/:id/start
Content-Type: application/json
```

```json
{
  "parameters": {
    "overrides": {
      "max_vus": 150
    }
  }
}
```

**响应 (202 Accepted)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "test_id": 1,
    "status": "running",
    "started_at": "2026-03-23T10:00:00Z"
  }
}
```

##### 3.8.5.3 获取压力测试报告

**请求**
```http
GET /api/v1/simulation/stress-tests/:id/report
```

**响应 (200 OK)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "test_name": "API并发压力测试",
    "test_type": "concurrent",
    "status": "completed",
    "duration_seconds": 420,
    "summary": {
      "total_requests": 50000,
      "failed_requests": 50,
      "success_rate": 99.9,
      "avg_response_time_ms": 85,
      "p50_response_time_ms": 72,
      "p95_response_time_ms": 150,
      "p99_response_time_ms": 200,
      "max_response_time_ms": 500,
      "requests_per_second": 120
    },
    "metrics": {
      "cpu_usage": [{"timestamp": 1711172400, "value": 45.2}],
      "memory_usage": [{"timestamp": 1711172400, "value": 60.5}],
      "network_io": [{"timestamp": 1711172400, "rx_mbps": 100, "tx_mbps": 50}]
    },
    "thresholds_passed": true,
    "generated_at": "2026-03-23T10:30:00Z"
  }
}
```

##### 3.8.5.4 实时指标流 (WebSocket)

**连接**
```http
WS /api/v1/simulation/stress-tests/:id/metrics
```

**服务端推送（实时指标）**
```json
{
  "event": "metrics",
  "data": {
    "vus": 85,
    "request_count": 12500,
    "failed_count": 12,
    "response_time_ms": {
      "avg": 88,
      "p50": 75,
      "p95": 155,
      "p99": 210
    },
    "timestamp": 1711172400000
  }
}
```

#### 3.8.6 仿真数据集管理 API

##### 3.8.6.1 创建数据集

**请求**
```http
POST /api/v1/simulation/datasets
Content-Type: application/json
```

```json
{
  "dataset_name": "传感器仿真数据集v1",
  "dataset_type": "sensor",
  "description": "用于模拟各种传感器数据",
  "tags": ["传感器", "仿真"],
  "schema": {
    "fields": [
      {"name": "temperature", "type": "float", "unit": "celsius", "min": -20, "max": 50},
      {"name": "humidity", "type": "float", "unit": "percent", "min": 0, "max": 100},
      {"name": "timestamp", "type": "datetime"}
    ]
  },
  "is_public": false
}
```

**数据集类型**

| 类型 | 说明 |
|------|------|
| sensor | 传感器数据 |
| voice | 语音数据 |
| image | 图像数据 |
| behavior | 行为数据 |
| environment | 环境数据 |

**响应 (201 Created)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "dataset_name": "传感器仿真数据集v1",
    "dataset_type": "sensor",
    "current_version": "1.0",
    "is_public": false,
    "created_at": "2026-03-23T10:00:00Z"
  }
}
```

##### 3.8.6.2 创建数据集版本

**请求**
```http
POST /api/v1/simulation/datasets/:id/versions
Content-Type: application/json
```

```json
{
  "version": "1.1",
  "description": "新增夜间场景数据",
  "data_file": "dataset_v1.1.csv",
  "is_published": false
}
```

**响应 (201 Created)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 2,
    "dataset_id": 1,
    "version": "1.1",
    "description": "新增夜间场景数据",
    "record_count": 10000,
    "is_published": false,
    "created_at": "2026-03-23T11:00:00Z"
  }
}
```

##### 3.8.6.3 对比数据集版本

**请求**
```http
POST /api/v1/simulation/datasets/:id/versions/compare
Content-Type: application/json
```

```json
{
  "version_a": "1.0",
  "version_b": "1.1"
}
```

**响应 (200 OK)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "version_a": "1.0",
    "version_b": "1.1",
    "record_count_diff": 5000,
    "schema_changes": [
      {"field": "temperature", "change": "unchanged"},
      {"field": "humidity", "change": "unchanged"},
      {"field": "light_level", "change": "added"}
    ],
    "statistical_diff": {
      "temperature": {"mean_diff": 0.5, "std_diff": 0.1},
      "humidity": {"mean_diff": -2.0, "std_diff": 0.3}
    }
  }
}
```

##### 3.8.6.4 导入数据集

**请求**
```http
POST /api/v1/simulation/datasets/import
Content-Type: multipart/form-data
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 数据文件（CSV/JSON格式） |
| dataset_id | int | 否 | 已有数据集ID |
| skip_duplicates | bool | 否 | 跳过重复数据 |

**响应 (202 Accepted)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "job_id": "job_20260323110001",
    "status": "processing",
    "estimated_completion": "2026-03-23T11:05:00Z"
  }
}
```

#### 3.8.7 CI/CD 集成 API

##### 3.8.7.1 Webhook 触发

**请求**
```http
POST /api/v1/simulation/webhooks/trigger
Content-Type: application/json
```

```json
{
  "event": "commit",
  "repository": "mdm-iot-platform",
  "branch": "main",
  "commit_sha": "abc123def456",
  "author": "developer@example.com",
  "test_suite": "regression",
  "callback_url": "https://ci.example.com/callback/token123"
}
```

**触发事件**

| 事件 | 说明 |
|------|------|
| commit | 代码提交 |
| tag | 标签创建 |
| pr_merge | PR合并 |
| manual | 手动触发 |

**测试套件**

| 套件 | 说明 |
|------|------|
| smoke | 冒烟测试 |
| regression | 回归测试 |
| full | 全量测试 |
| custom | 自定义 |

**响应 (202 Accepted)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "job_id": "job_abc123",
    "status": "queued",
    "queue_position": 3,
    "estimated_start": "2026-03-23T10:03:00Z"
  }
}
```

##### 3.8.7.2 创建集成

**请求**
```http
POST /api/v1/simulation/integrations
Content-Type: application/json
```

```json
{
  "integration_name": "GitHub Actions",
  "integration_type": "github_actions",
  "config": {
    "repository": "yangkai258/mdm-iot-platform",
    "workflow_file": ".github/workflows/simulation-test.yml",
    "trigger_on": ["push", "pull_request"]
  },
  "secret": "github_token_xxx",
  "is_active": true
}
```

**响应 (201 Created)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "integration_name": "GitHub Actions",
    "integration_type": "github_actions",
    "is_active": true,
    "created_at": "2026-03-23T10:00:00Z"
  }
}
```

#### 3.8.8 通用错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 资源不存在 |
| 1002 | 参数无效 |
| 1003 | 权限不足 |
| 2001 | 宠物未运行 |
| 2002 | 执行超时 |
| 2003 | 执行失败 |
| 3001 | 场景不存在 |
| 3002 | 场景运行中 |
| 4001 | 压力测试运行中 |
| 4002 | 压力测试未运行 |
| 5001 | 数据集不存在 |
| 5002 | 版本不存在 |
| 5003 | 导入失败 |
| 6001 | Webhook触发失败 |
| 6002 | 集成配置无效 |

---


## 五、数据库设计

### 5.1 虚拟宠物表 (simulation_pets)

``sql
CREATE TABLE simulation_pets (
    id                 BIGSERIAL PRIMARY KEY,
    pet_name           VARCHAR(100) NOT NULL,
    pet_type           VARCHAR(50) NOT NULL,
    personality        JSONB,
    capabilities       JSONB,
    environment_id     BIGINT,
    status             VARCHAR(20) DEFAULT 'idle',
    current_emotion    VARCHAR(50),
    current_position   JSONB,
    battery_level      DECIMAL(5,2) DEFAULT 1.0,
    created_by         BIGINT REFERENCES users(id),
    created_at         TIMESTAMP DEFAULT NOW(),
    updated_at         TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_simulation_pets_type ON simulation_pets(pet_type);
CREATE INDEX idx_simulation_pets_status ON simulation_pets(status);
CREATE INDEX idx_simulation_pets_created_by ON simulation_pets(created_by);
``

### 5.2 虚拟环境表 (simulation_environments)

``sql
CREATE TABLE simulation_environments (
    id              BIGSERIAL PRIMARY KEY,
    env_name        VARCHAR(255) NOT NULL,
    env_type        VARCHAR(50) NOT NULL,
    background_url  VARCHAR(500),
    objects         JSONB[],
    lighting        JSONB,
    config          JSONB,
    is_public       BOOLEAN DEFAULT FALSE,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_simulation_environments_type ON simulation_environments(env_type);
CREATE INDEX idx_simulation_environments_public ON simulation_environments(is_public);
``

### 5.3 仿真场景表 (simulation_scenarios)

``sql
CREATE TABLE simulation_scenarios (
    id              BIGSERIAL PRIMARY KEY,
    scenario_name   VARCHAR(255) NOT NULL,
    scenario_type   VARCHAR(50) NOT NULL,
    environment     JSONB,
    objects         JSONB[],
    events          JSONB[],
    config          JSONB,
    is_public       BOOLEAN DEFAULT FALSE,
    score           DECIMAL(5,2),
    downloads       INT DEFAULT 0,
    tags            VARCHAR(100)[],
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_simulation_scenarios_type ON simulation_scenarios(scenario_type);
CREATE INDEX idx_simulation_scenarios_public ON simulation_scenarios(is_public);
CREATE INDEX idx_simulation_scenarios_tags ON simulation_scenarios USING GIN(tags);
``

### 5.4 测试用例表 (simulation_testcases)

``sql
CREATE TABLE simulation_testcases (
    id              BIGSERIAL PRIMARY KEY,
    case_name       VARCHAR(255) NOT NULL,
    case_type       VARCHAR(30) NOT NULL,
    module          VARCHAR(50),
    priority        VARCHAR(20) DEFAULT 'medium',
    status          VARCHAR(20) DEFAULT 'draft',
    description     TEXT,
    preconditions   TEXT,
    test_steps      JSONB,
    expected_result TEXT,
    tags            VARCHAR(100)[],
    dependencies    BIGINT[],
    version         INT DEFAULT 1,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_testcases_type ON simulation_testcases(case_type);
CREATE INDEX idx_testcases_module ON simulation_testcases(module);
CREATE INDEX idx_testcases_priority ON simulation_testcases(priority);
CREATE INDEX idx_testcases_status ON simulation_testcases(status);
CREATE INDEX idx_testcases_tags ON simulation_testcases USING GIN(tags);
``

### 5.5 测试执行记录表 (test_executions)

``sql
CREATE TABLE test_executions (
    id                 VARCHAR(50) PRIMARY KEY,
    testcase_id        BIGINT REFERENCES simulation_testcases(id),
    batch_id           VARCHAR(50),
    execution_type     VARCHAR(20),
    trigger_params     JSONB,
    status             VARCHAR(20) DEFAULT 'pending',
    progress           INT DEFAULT 0,
    current_step       INT DEFAULT 0,
    start_time         TIMESTAMP,
    end_time           TIMESTAMP,
    duration_ms        INT,
    environment        JSONB,
    result_details     JSONB,
    screenshots        VARCHAR(500)[],
    logs               TEXT,
    error_message      TEXT,
    retry_count        INT DEFAULT 0,
    created_at         TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_test_executions_case ON test_executions(testcase_id, created_at DESC);
CREATE INDEX idx_test_executions_batch ON test_executions(batch_id);
CREATE INDEX idx_test_executions_status ON test_executions(status);
CREATE INDEX idx_test_executions_type ON test_executions(execution_type);
``

### 5.6 测试报告表 (test_reports)

``sql
CREATE TABLE test_reports (
    id              BIGSERIAL PRIMARY KEY,
    report_name     VARCHAR(255),
    report_type     VARCHAR(20) DEFAULT 'execution',
    execution_ids   VARCHAR(50)[],
    batch_id        VARCHAR(50),
    summary         JSONB,
    pass_count      INT DEFAULT 0,
    fail_count      INT DEFAULT 0,
    skip_count      INT DEFAULT 0,
    total_count     INT DEFAULT 0,
    pass_rate       DECIMAL(5,2),
    avg_duration_ms INT,
    coverage        JSONB,
    trend_data      JSONB,
    failed_cases    JSONB,
    created_by      BIGINT REFERENCES users(id),
    generated_at    TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_test_reports_type ON test_reports(report_type);
CREATE INDEX idx_test_reports_generated ON test_reports(generated_at DESC);
``

### 5.7 回放记录表 (playback_records)

``sql
CREATE TABLE playback_records (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100),
    pet_id          BIGINT REFERENCES simulation_pets(id),
    record_type     VARCHAR(20) NOT NULL,
    start_time      TIMESTAMP NOT NULL,
    end_time        TIMESTAMP,
    duration_ms     INT,
    sensor_data     JSONB,
    user_actions    JSONB[],
    events          JSONB[],
    playback_url    VARCHAR(500),
    file_path       VARCHAR(500),
    file_size       BIGINT,
    status          VARCHAR(20) DEFAULT 'recording',
    metadata        JSONB,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_playback_records_device ON playback_records(device_id);
CREATE INDEX idx_playback_records_pet ON playback_records(pet_id);
CREATE INDEX idx_playback_records_status ON playback_records(status);
CREATE INDEX idx_playback_records_time ON playback_records(start_time DESC);
``

### 5.8 压力测试表 (stress_tests)

``sql
CREATE TABLE stress_tests (
    id                 BIGSERIAL PRIMARY KEY,
    test_name          VARCHAR(255) NOT NULL,
    test_type          VARCHAR(20) NOT NULL,
    config             JSONB NOT NULL,
    status             VARCHAR(20) DEFAULT 'draft',
    start_time         TIMESTAMP,
    end_time           TIMESTAMP,
    duration_seconds   INT,
    metrics            JSONB,
    summary            JSONB,
    thresholds_passed  BOOLEAN,
    report_url         VARCHAR(500),
    created_by         BIGINT REFERENCES users(id),
    created_at         TIMESTAMP DEFAULT NOW(),
    updated_at         TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_stress_tests_type ON stress_tests(test_type);
CREATE INDEX idx_stress_tests_status ON stress_tests(status);
CREATE INDEX idx_stress_tests_created ON stress_tests(created_at DESC);
``

### 5.9 仿真数据集表 (simulation_datasets)

``sql
CREATE TABLE simulation_datasets (
    id                 BIGSERIAL PRIMARY KEY,
    dataset_name       VARCHAR(255) NOT NULL,
    dataset_type       VARCHAR(50) NOT NULL,
    description        TEXT,
    schema             JSONB,
    current_version    VARCHAR(20),
    record_count       INT DEFAULT 0,
    file_size          BIGINT DEFAULT 0,
    tags               VARCHAR(100)[],
    is_public          BOOLEAN DEFAULT FALSE,
    created_by         BIGINT REFERENCES users(id),
    created_at         TIMESTAMP DEFAULT NOW(),
    updated_at         TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_simulation_datasets_type ON simulation_datasets(dataset_type);
CREATE INDEX idx_simulation_datasets_public ON simulation_datasets(is_public);
CREATE INDEX idx_simulation_datasets_tags ON simulation_datasets USING GIN(tags);
``

### 5.10 数据集版本表 (simulation_dataset_versions)

``sql
CREATE TABLE simulation_dataset_versions (
    id              BIGSERIAL PRIMARY KEY,
    dataset_id      BIGINT REFERENCES simulation_datasets(id) ON DELETE CASCADE,
    version         VARCHAR(20) NOT NULL,
    description     TEXT,
    record_count    INT DEFAULT 0,
    file_size       BIGINT DEFAULT 0,
    file_path       VARCHAR(500),
    schema_hash     VARCHAR(64),
    stats           JSONB,
    is_published    BOOLEAN DEFAULT FALSE,
    published_at    TIMESTAMP,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_dataset_version_unique ON simulation_dataset_versions(dataset_id, version);
CREATE INDEX idx_dataset_version_published ON simulation_dataset_versions(dataset_id, is_published);
``

### 5.11 CI/CD 集成表 (simulation_integrations)

``sql
CREATE TABLE simulation_integrations (
    id                 BIGSERIAL PRIMARY KEY,
    integration_name   VARCHAR(255) NOT NULL,
    integration_type   VARCHAR(50) NOT NULL,
    config             JSONB,
    secret             VARCHAR(500),
    is_active          BOOLEAN DEFAULT TRUE,
    last_triggered_at  TIMESTAMP,
    last_status        VARCHAR(20),
    created_by         BIGINT REFERENCES users(id),
    created_at         TIMESTAMP DEFAULT NOW(),
    updated_at         TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_simulation_integrations_type ON simulation_integrations(integration_type);
CREATE INDEX idx_simulation_integrations_active ON simulation_integrations(is_active);
``

### 5.12 CI/CD 任务表 (simulation_cicd_jobs)

``sql
CREATE TABLE simulation_cicd_jobs (
    id              VARCHAR(50) PRIMARY KEY,
    integration_id BIGINT REFERENCES simulation_integrations(id),
    event           VARCHAR(20) NOT NULL,
    repository      VARCHAR(255),
    branch          VARCHAR(100),
    commit_sha      VARCHAR(64),
    test_suite      VARCHAR(50),
    callback_url    VARCHAR(500),
    status          VARCHAR(20) DEFAULT 'queued',
    queue_position  INT DEFAULT 0,
    execution_ids   VARCHAR(50)[],
    report_id       BIGINT,
    result_summary  JSONB,
    error_message   TEXT,
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_simulation_cicd_jobs_status ON simulation_cicd_jobs(status);
CREATE INDEX idx_simulation_cicd_jobs_repository ON simulation_cicd_jobs(repository, branch);
CREATE INDEX idx_simulation_cicd_jobs_created ON simulation_cicd_jobs(created_at DESC);
``

### 5.13 数据导入导出任务表 (dataset_jobs)

``sql
CREATE TABLE dataset_jobs (
    id                  VARCHAR(50) PRIMARY KEY,
    job_type            VARCHAR(20) NOT NULL,
    dataset_id          BIGINT REFERENCES simulation_datasets(id),
    version             VARCHAR(20),
    file_path           VARCHAR(500),
    file_size           BIGINT,
    status              VARCHAR(20) DEFAULT 'processing',
    progress            INT DEFAULT 0,
    records_total       INT DEFAULT 0,
    records_processed   INT DEFAULT 0,
    records_failed      INT DEFAULT 0,
    error_message       TEXT,
    estimated_completion TIMESTAMP,
    created_by          BIGINT REFERENCES users(id),
    created_at          TIMESTAMP DEFAULT NOW(),
    completed_at        TIMESTAMP
);

CREATE INDEX idx_dataset_jobs_status ON dataset_jobs(status);
CREATE INDEX idx_dataset_jobs_dataset ON dataset_jobs(dataset_id);
``

## 六、前端页面清单

### 6.1 虚拟宠物

| 页面 | 路由 | 说明 |
|------|------|------|
| 虚拟宠物列表 | /simulation/pets | 虚拟宠物列表 |
| 创建虚拟宠物 | /simulation/pets/create | 创建虚拟宠物 |
| 虚拟宠物详情 | /simulation/pets/:id | 虚拟宠物详情 |
| 虚拟交互 | /simulation/pets/:id/interact | 虚拟交互界面 |

### 6.2 测试用例

| 页面 | 路由 | 说明 |
|------|------|------|
| 用例列表 | /simulation/testcases | 测试用例列表 |
| 创建用例 | /simulation/testcases/create | 创建测试用例 |
| 用例详情 | /simulation/testcases/:id | 用例详情/编辑 |
| 执行结果 | /simulation/testcases/:id/results | 用例执行结果 |

### 6.3 测试报告

| 页面 | 路由 | 说明 |
|------|------|------|
| 报告列表 | /simulation/reports | 报告列表 |
| 报告详情 | /simulation/reports/:id | 报告详情 |
| 趋势分析 | /simulation/reports/trends | 通过率趋势 |

### 6.4 回放系统

| 页面 | 路由 | 说明 |
|------|------|------|
| 回放列表 | /simulation/playbacks | 回放列表 |
| 回放详情 | /simulation/playbacks/:id | 回放详情 |
| 回放播放器 | /simulation/playbacks/:id/player | 回放播放界面 |

### 6.5 仿真场景

| 页面 | 路由 | 说明 |
|------|------|------|
| 场景列表 | /simulation/scenarios | 场景库列表 |
| 创建场景 | /simulation/scenarios/create | 场景编辑器 |
| 场景详情 | /simulation/scenarios/:id | 场景详情 |
| 运行场景 | /simulation/scenarios/:id/run | 场景运行 |
| 场景导入 | /simulation/scenarios/import | 导入场景 |
| 场景导出 | /simulation/scenarios/export/:id | 导出场景 |

### 6.6 压力测试

| 页面 | 路由 | 说明 |
|------|------|------|
| 压力测试列表 | /simulation/stress-tests | 测试列表 |
| 创建测试 | /simulation/stress-tests/create | 创建压力测试 |
| 测试详情 | /simulation/stress-tests/:id | 测试详情/监控 |
| 测试报告 | /simulation/stress-tests/:id/report | 性能报告 |

---


### 6.7 仿真数据集

| 页面 | 路由 | 说明 |
|------|------|------|
| 数据集列表 | /simulation/datasets | 数据集列表 |
| 创建数据集 | /simulation/datasets/create | 创建数据集 |
| 数据集详情 | /simulation/datasets/:id | 数据集详情/版本管理 |
| 数据集导入 | /simulation/datasets/import | 导入数据集 |
| 版本对比 | /simulation/datasets/:id/compare | 版本对比 |

### 6.8 CI/CD 集成

| 页面 | 路由 | 说明 |
|------|------|------|
| 集成列表 | /simulation/integrations | 集成配置列表 |
| 创建集成 | /simulation/integrations/create | 创建集成 |
| 集成详情 | /simulation/integrations/:id | 集成配置详情 |
| 任务列表 | /simulation/integrations/:id/jobs | 任务执行历史 |

## 七、验收标准

### 7.1 虚拟宠物仿真

| 验收点 | 标准 |
|--------|------|
| 虚拟宠物运行 | 虚拟宠物稳定运行不掉帧 |
| 交互响应 | 交互响应<200ms |
| 多宠物支持 | 同时运行5个虚拟宠物流畅 |

### 7.2 自动化测试

| 验收点 | 标准 |
|--------|------|
| 用例执行 | 用例按预期执行，结果准确 |
| 批量执行 | 100个用例并发执行正常 |
| 报告生成 | 执行完成5分钟内生成报告 |
| 覆盖率 | 核心功能覆盖率>80% |

### 7.3 回放系统

| 验收点 | 标准 |
|--------|------|
| 回放精度 | 回放与录制偏差<1% |
| 差异对比 | 差异对比准确 |
| 分享可用 | 回放可被其他用户播放 |

### 7.4 仿真场景

| 验收点 | 标准 |
|--------|------|
| 场景运行 | 预置场景100%可运行 |
| 场景编辑 | 自定义场景编辑保存成功 |
| 场景导入导出 | 导入导出成功率>99% |

### 7.5 压力测试

| 验收点 | 标准 |
|--------|------|
| 性能数据准确性 | 性能数据与实际一致 |
| 并发支持 | 支持100+并发连接 |
| 稳定性测试 | 7x24小时运行无内存泄漏 |

### 7.6 仿真数据集

| 验收点 | 标准 |
|--------|------|
| 数据集创建 | 数据集创建成功，schema验证通过 |
| 数据集版本 | 版本创建、发布、回滚正常 |
| 数据集导入 | 1万条数据导入时间<30秒 |
| 数据集导出 | 导出文件完整性校验通过 |
| 版本对比 | 对比结果准确显示差异 |

### 7.7 CI/CD 集成

| 验收点 | 标准 |
|--------|------|
| Webhook触发 | 代码提交后5分钟内自动触发测试 |
| 集成稳定性 | 集成可用性>99% |
| 结果回调 | 测试结果正确回调至CI系统 |
| 并行任务 | 支持5个以上并行CI任务 |
| 通知准确 | 成功/失败通知100%送达 |

---

## 八、页面布局规范

### 8.1 虚拟宠物列表页面（/simulation/pets）

**布局结构：**
1. 面包屑 → 页面标题
2. 操作栏（创建虚拟宠物靠左）
3. 虚拟宠物卡片网格（展示宠物形象/类型/状态）

**按钮规范：**
- [创建虚拟宠物] — 左对齐
- [交互] [详情] [删除] — 卡片内右对齐

### 8.2 测试用例列表页面（/simulation/testcases）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片 #F2F3F5）：用例类型 / 模块 / 优先级
3. 操作栏（创建用例靠左，批量执行靠右）
4. 测试用例列表表格

**按钮规范：**
- [创建用例] — 左对齐
- [批量执行] — 操作栏中间
- [详情] [编辑] [执行] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 用例名称 | 200px | - |
| 类型 | 100px | 功能/性能/压力/回归/冒烟 |
| 优先级 | 100px | high/medium/low |
| 模块 | 120px | - |
| 状态 | 80px | - |
| 操作 | 120px | 详情/编辑/执行/删除 |

**分页：** 右下角，10/20/50/100 条

### 8.3 测试报告页面（/simulation/reports）

**布局结构：**
1. 面包屑 → 页面标题
2. 统计概览卡片（通过率/覆盖率/总用例数/失败数）—— 白色
3. 趋势图表
4. 测试报告列表表格

**按钮规范：**
- [查看详情] — 行内右对齐

**分页：** 右下角，10/20/50/100 条

### 8.4 回放列表页面（/simulation/playbacks）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：设备ID / 日期范围
3. 操作栏（开始录制靠左）
4. 回放列表表格

**按钮规范：**
- [开始录制] — 左对齐
- [播放] [对比] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 回放ID | 150px | - |
| 设备ID | 150px | - |
| 录制类型 | 100px | auto/manual |
| 时长 | 100px | - |
| 状态 | 100px | recording/completed |
| 创建时间 | 150px | - |
| 操作 | 120px | 播放/对比/删除 |

**分页：** 右下角，10/20/50/100 条

### 8.5 仿真场景页面（/simulation/scenarios）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：场景类型 / 是否公开
3. 操作栏（创建场景/导入靠左，导出靠右）
4. 场景列表卡片网格

**按钮规范：**
- [创建场景] [导入] — 左对齐
- [导出] — 操作栏中间
- [编辑] [运行] [删除] — 卡片内右对齐

### 8.6 压力测试页面（/simulation/stress-tests）

**布局结构：**
1. 面包屑 → 页面标题
2. 操作栏（创建测试靠左）
3. 压力测试列表表格

**按钮规范：**
- [创建测试] — 左对齐
- [开始] [停止] [查看报告] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 测试名称 | 200px | - |
| 测试类型 | 100px | 并发/性能/稳定性 |
| 状态 | 100px | draft/running/completed |
| 创建时间 | 150px | - |
| 操作 | 120px | 开始/停止/查看报告 |

**分页：** 右下角，10/20/50/100 条

### 8.7 弹窗规范

| 类型 | 使用场景 |
|------|----------|
| Drawer 抽屉 | 创建/编辑测试用例、创建/编辑场景、测试报告详情 |
| Dialog 对话框 | 确认删除、确认开始/停止测试 |
| 全屏模态 | 暂无复杂表单场景 |
