# Sprint 30 规划

**时间**：2026-11-29
**状态**：待开始
**Sprint 周期**：2 周（2026-11-29 ～ 2026-12-12）

---

## 一、Sprint 目标

**目标：** 健康管理完善与宠物社交

完善数字孪生健康管理（睡眠分析、体重追踪、饮食记录），并实现宠物社交功能，构建完整的宠物健康管理和社交生态。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **睡眠分析完善** | 实现深睡/浅睡/REM分析，生成睡眠报告 | sleep_analysis.go | P0 |
| P0-2 | **体重追踪完善** | 实现体重曲线和营养建议 | weight_tracking.go | P0 |
| P1-1 | **饮食记录 API** | 完成 `/api/v1/diet/*` 饮食记录 | diet_controller.go | P1 |
| P1-2 | **营养分析服务** | 基于饮食数据生成营养分析 | nutrition_analyzer.go | P1 |
| P2-1 | **宠物社交完善** | 完善宠物朋友圈和互动功能 | pet_social_enhanced.go | P2 |
| P2-2 | **宠物视频互动** | 实现宠物视频录制和分享 | pet_video互动_enhanced.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **睡眠分析报告页面** | 完成 SleepAnalysisView.vue 睡眠报告 | SleepAnalysisView.vue | P0 |
| PF0-2 | **体重追踪页面** | 完成 WeightTrackingView.vue 体重曲线和营养建议 | WeightTrackingView.vue | P0 |
| PF1-1 | **饮食记录页面** | 完成 DietRecordView.vue 饮食记录和统计 | DietRecordView.vue | P1 |
| PF1-2 | **营养分析页面** | 完成 NutritionAnalysisView.vue 营养报告 | NutritionAnalysisView.vue | P1 |
| PF2-1 | **宠物社交完善页面** | 完成 PetSocialEnhancedView.vue 朋友圈互动 | PetSocialEnhancedView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/health/sleep/:pet_id` | GET | 睡眠数据 |
| `GET /api/v1/health/sleep/:pet_id/report` | GET | 睡眠分析报告 |
| `GET /api/v1/health/weight/:pet_id` | GET | 体重历史 |
| `GET /api/v1/health/weight/:pet_id/nutrition-advice` | GET | 营养建议 |
| `GET /api/v1/diet/meals` | GET | 饮食记录列表 |
| `POST /api/v1/diet/meals` | POST | 记录饮食 |
| `PUT /api/v1/diet/meals/:id` | PUT | 更新饮食记录 |
| `DELETE /api/v1/diet/meals/:id` | DELETE | 删除饮食记录 |
| `GET /api/v1/diet/:pet_id/nutrition-report` | GET | 营养分析报告 |
| `GET /api/v1/social/feed` | GET | 宠物朋友圈 |
| `POST /api/v1/social/feed` | POST | 发布动态 |
| `POST /api/v1/social/feed/:id/like` | POST | 点赞 |
| `POST /api/v1/social/feed/:id/comment` | POST | 评论 |

### 数据库设计

```sql
-- 睡眠分析表
CREATE TABLE sleep_analysis (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    sleep_date      DATE NOT NULL,
    sleep_duration  INT,                            -- 总睡眠时长 分钟
    deep_sleep_min  INT,                            -- 深睡分钟
    light_sleep_min INT,                            -- 浅睡分钟
    rem_min         INT,                            -- REM分钟
    awakenings      INT,                            -- 夜醒次数
    sleep_quality   VARCHAR(20),                    -- 'poor'/'fair'/'good'/'excellent'
    sleep_report    JSONB,                           -- 详细报告
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(pet_id, sleep_date)
);

-- 饮食记录表
CREATE TABLE diet_meals (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    meal_type       VARCHAR(20) NOT NULL,            -- 'breakfast'/'lunch'/'dinner'/'snack'
    food_items      JSONB NOT NULL,                   -- 食物明细
    calories        DECIMAL(8,2),                    -- 总热量 kcal
    recorded_at     TIMESTAMP NOT NULL,
    notes           TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_pet_diet (pet_id, recorded_at DESC)
);

-- 营养分析表
CREATE TABLE nutrition_reports (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    report_date     DATE NOT NULL,
    total_calories  DECIMAL(8,2),
    protein_grams   DECIMAL(8,2),
    fat_grams       DECIMAL(8,2),
    carbs_grams     DECIMAL(8,2),
    fiber_grams     DECIMAL(8,2),
    nutrition_score DECIMAL(5,2),                    -- 营养评分 0-100
    advice          TEXT,                            -- 营养建议
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(pet_id, report_date)
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 睡眠分析 | 深睡/浅睡/REM分析准确率>85% | 对比专业设备 |
| 体重追踪 | 体重曲线展示正常，趋势预测准确>80% | 趋势分析测试 |
| 饮食记录 | 记录准确性>90%，营养统计正确 | 饮食测试 |
| 宠物社交 | 动态发布/点赞/评论正常 | 社交功能测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 睡眠报告生成 | < 5s |
| 营养分析报告 | < 3s |
| 社交feed加载 | < 2s |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| MODULE_DIGITAL_TWIN | 数字孪生健康管理 |
| Sprint 28 宠物社交 | 社交基础能力 |
| 营养数据库 | 食物营养数据 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 睡眠分析不准确 | 健康建议错误 | 专业设备对比校准 |
| 营养数据不完整 | 分析偏差 | 扩充食物数据库 |
| 社交内容违规 | 平台风险 | 内容审核机制 |
