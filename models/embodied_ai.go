package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ============ 具身智能状态 ============

// EmbodiedAIState 具身AI状态
type EmbodiedAIState struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	StateKey       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"state_key"` // 唯一标识（如 pet_uuid + "_brain"）
	EntityType     string         `gorm:"type:varchar(32);not null" json:"entity_type"`         // pet/robot/device
	EntityID       string         `gorm:"type:varchar(64);index;not null" json:"entity_id"`     // 关联实体ID

	// 位置与空间
	PositionX      float64        `gorm:"type:decimal(10,4)" json:"position_x"`      // X坐标
	PositionY      float64        `gorm:"type:decimal(10,4)" json:"position_y"`      // Y坐标
	PositionZ      float64        `gorm:"type:decimal(10,4)" json:"position_z"`      // Z坐标（高度）
	Orientation    float64        `gorm:"type:decimal(8,4)" json:"orientation"`    // 朝向角度（弧度）
	LocationName   string         `gorm:"type:varchar(128)" json:"location_name"`   // 当前位置名称

	// 感知状态
	PerceptionMode string         `gorm:"type:varchar(32);default:'active'" json:"perception_mode"` // active/sleeping/hibernating
	AlertLevel     string         `gorm:"type:varchar(16);default:'normal'" json:"alert_level"`    // normal/alert/emergency

	// 能量与情绪
	EnergyLevel    float64        `gorm:"type:decimal(5,2);default:100" json:"energy_level"`     // 能量 0-100
	Mood           string         `gorm:"type:varchar(32);default:'happy'" json:"mood"`          // happy/calm/excited/anxious/sad/sleepy
	MoodIntensity  float64        `gorm:"type:decimal(3,2);default:0.5" json:"mood_intensity"`  // 情绪强度 0-1

	// 当前目标
	CurrentGoal    string         `gorm:"type:varchar(256)" json:"current_goal"`    // 当前目标描述
	GoalProgress   float64        `gorm:"type:decimal(5,2)" json:"goal_progress"`   // 目标进度 0-100
	GoalTargetX    float64        `gorm:"type:decimal(10,4)" json:"goal_target_x"`  // 目标位置X
	GoalTargetY    float64        `gorm:"type:decimal(10,4)" json:"goal_target_y"`  // 目标位置Y
	GoalTargetZ    float64        `gorm:"type:decimal(10,4)" json:"goal_target_z"`  // 目标位置Z

	// 环境感知缓存
	LastSeenObjects pq.StringArray `gorm:"type:text[]" json:"last_seen_objects"`    // 最近看到的物体
	NearbyEntities  pq.StringArray `gorm:"type:text[]" json:"nearby_entities"`       // 附近实体
	EnvironmentTags pq.StringArray `gorm:"type:text[]" json:"environment_tags"`      // 环境标签

	// 能力状态
	Capabilities   string         `gorm:"type:jsonb" json:"capabilities"`           // 可用能力JSON
	SensoryStatus  string         `gorm:"type:jsonb" json:"sensory_status"`         // 传感器状态JSON

	// 时间戳
	LastUpdatedAt  time.Time      `json:"last_updated_at"`
	LastActiveAt   time.Time      `json:"last_active_at"`
	TenantID       string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (s *EmbodiedAIState) BeforeCreate(tx *gorm.DB) error {
	if s.StateKey == "" {
		s.StateKey = uuid.New().String()
	}
	return nil
}

// TableName 表名
func (EmbodiedAIState) TableName() string {
	return "embodied_ai_states"
}

// ============ 探索会话 ============

// ExplorationSession 探索会话
type ExplorationSession struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	SessionKey       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"session_key"`
	EntityID         string         `gorm:"type:varchar(64);index;not null" json:"entity_id"`        // 探索主体ID
	EntityType       string         `gorm:"type:varchar(32);not null" json:"entity_type"`          // pet/robot

	// 探索配置
	Strategy         string         `gorm:"type:varchar(32);not null" json:"strategy"`            // random/coverage/frontier/goal_oriented
	ExplorationGoal   string         `gorm:"type:varchar(256)" json:"exploration_goal"`            // 探索目标描述
	MaxDuration       int            `gorm:"type:integer" json:"max_duration"`                     // 最大持续时间（秒）
	MaxDistance       float64        `gorm:"type:decimal(10,2)" json:"max_distance"`               // 最大探索距离（米）

	// 探索范围
	StartX            float64        `gorm:"type:decimal(10,4)" json:"start_x"`
	StartY            float64        `gorm:"type:decimal(10,4)" json:"start_y"`
	BoundaryMinX      float64        `gorm:"type:decimal(10,4)" json:"boundary_min_x"`
	BoundaryMaxX      float64        `gorm:"type:decimal(10,4)" json:"boundary_max_x"`
	BoundaryMinY      float64        `gorm:"type:decimal(10,4)" json:"boundary_min_y"`
	BoundaryMaxY      float64        `gorm:"type:decimal(10,4)" json:"boundary_max_y"`

	// 路径记录
	Waypoints         pq.StringArray `gorm:"type:text[]" json:"waypoints"`                          // [{"x":1,"y":2,"t":"2024-01-01T10:00:00Z"},...]
	VisitedAreas      pq.StringArray `gorm:"type:text[]" json:"visited_areas"`                    // 已访问区域
	DiscoveredObjects pq.StringArray `gorm:"type:text[]" json:"discovered_objects"`                // 发现的物体

	// 探索结果
	CoverageRate      float64        `gorm:"type:decimal(5,2)" json:"coverage_rate"`              // 覆盖率 0-100
	PathLength         float64        `gorm:"type:decimal(10,2)" json:"path_length"`              // 路径总长度（米）
	NewDiscoveryCount  int            `gorm:"type:integer;default:0" json:"new_discovery_count"`  // 新发现数量

	// 状态
	Status             string         `gorm:"type:varchar(20);default:'active'" json:"status"`    // active/paused/completed/aborted
	Progress           float64        `gorm:"type:decimal(5,2)" json:"progress"`                  // 进度 0-100

	// 时间
	StartedAt          time.Time      `json:"started_at"`
	EndedAt           *time.Time      `json:"ended_at"`
	PauseDuration      int            `gorm:"type:integer;default:0" json:"pause_duration"`       // 暂停总时长（秒）

	TenantID           string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (e *ExplorationSession) BeforeCreate(tx *gorm.DB) error {
	if e.SessionKey == "" {
		e.SessionKey = uuid.New().String()
	}
	return nil
}

// TableName 表名
func (ExplorationSession) TableName() string {
	return "exploration_sessions"
}

// ============ 环境地图 ============

// EnvironmentMap 环境地图
type EnvironmentMap struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	MapKey          string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"map_key"`
	EntityID        string         `gorm:"type:varchar(64);index;not null" json:"entity_id"`      // 所属实体ID
	EntityType      string         `gorm:"type:varchar(32);not null" json:"entity_type"`        // pet/robot

	// 地图元数据
	MapName         string         `gorm:"type:varchar(128)" json:"map_name"`
	MapType         string         `gorm:"type:varchar(32);default:'2d_grid'" json:"map_type"` // 2d_grid/3d_pointcloud/topological
	Resolution      float64        `gorm:"type:decimal(8,4)" json:"resolution"`                // 分辨率（米/格）
	Scale            float64        `gorm:"type:decimal(10,2)" json:"scale"`                    // 地图总尺寸（米）

	// 空间范围
	OriginX          float64        `gorm:"type:decimal(10,4)" json:"origin_x"`                // 原点X
	OriginY          float64        `gorm:"type:decimal(10,4)" json:"origin_y"`                // 原点Y
	Width            float64        `gorm:"type:decimal(10,2)" json:"width"`                   // 宽度（米）
	Height           float64        `gorm:"type:decimal(10,2)" json:"height"`                  // 高度（米）

	// 地图数据
	GridData         string         `gorm:"type:jsonb" json:"grid_data"`                        // 栅格数据（占用/空闲/未知）
	Obstacles        string         `gorm:"type:jsonb" json:"obstacles"`                       // 障碍物列表
	Landmarks        string         `gorm:"type:jsonb" json:"landmarks"`                       // 地标点
	Regions          string         `gorm:"type:jsonb" json:"regions"`                        // 区域划分

	// 语义信息
	SemanticLabels   pq.StringArray `gorm:"type:text[]" json:"semantic_labels"`                // 语义标签
	RoomNames        pq.StringArray `gorm:"type:text[]" json:"room_names"`                     // 房间名称

	// 地图质量
	Confidence       float64        `gorm:"type:decimal(5,2)" json:"confidence"`               // 置信度 0-100
	ExploredRatio    float64        `gorm:"type:decimal(5,2)" json:"explored_ratio"`           // 探索比例 0-100

	// 状态
	IsActive         bool           `gorm:"type:boolean;default:true" json:"is_active"`
	Version          int            `gorm:"type:integer;default:1" json:"version"`

	// 时间
	LastUpdatedAt    time.Time      `json:"last_updated_at"`
	TenantID          string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (m *EnvironmentMap) BeforeCreate(tx *gorm.DB) error {
	if m.MapKey == "" {
		m.MapKey = uuid.New().String()
	}
	return nil
}

// TableName 表名
func (EnvironmentMap) TableName() string {
	return "environment_maps"
}

// ============ 请求/响应结构 ============

// PerceiveRequest 环境感知请求
type PerceiveRequest struct {
	EntityID    string  `json:"entity_id" binding:"required"`                           // 实体ID
	EntityType  string  `json:"entity_type" binding:"required"`                         // pet/robot
	SensorData  string  `json:"sensor_data" gorm:"type:jsonb"`                         // 传感器原始数据
	PerceiveMode string `json:"perceive_mode" gorm:"size:32"`                          // visual/audio/tactile/all
	FocusArea   string  `json:"focus_area" gorm:"type:jsonb"`                          // 关注区域
}

// SpatialRequest 空间认知请求
type SpatialRequest struct {
	EntityID    string  `json:"entity_id" binding:"required"`
	QueryType   string  `json:"query_type" binding:"required"`                         // position/direction/distance/map
	TargetX     float64 `json:"target_x"`                                             // 目标位置X
	TargetY     float64 `json:"target_y"`                                             // 目标位置Y
	TargetID    string  `json:"target_id"`                                            // 目标实体ID
	MapKey      string  `json:"map_key"`                                              // 地图Key
	QueryScope  string  `json:"query_scope"`                                          // local/global
}

// ExploreRequest 自主探索请求
type ExploreRequest struct {
	EntityID       string  `json:"entity_id" binding:"required"`
	EntityType     string  `json:"entity_type" binding:"required"`
	Strategy       string  `json:"strategy" binding:"required"`                         // random/coverage/frontier/goal_oriented
	ExplorationGoal string `json:"exploration_goal"`                                   // 探索目标
	MaxDuration    int     `json:"max_duration"`                                      // 最大持续时间（秒）
	MaxDistance    float64 `json:"max_distance"`                                      // 最大探索距离（米）
	Boundary       string  `json:"boundary" gorm:"type:jsonb"`                         // 探索边界
}

// InteractRequest 环境交互请求
type InteractRequest struct {
	EntityID     string `json:"entity_id" binding:"required"`
	ActionType   string `json:"action_type" binding:"required"`                       // grasp/move/place/push/pull/activate
	TargetX      float64 `json:"target_x"`                                            // 目标位置X
	TargetY      float64 `json:"target_y"`                                            // 目标位置Y
	TargetZ      float64 `json:"target_z"`                                            // 目标位置Z
	TargetObject string `json:"target_object"`                                        // 目标物体ID
	Parameters   string `json:"parameters" gorm:"type:jsonb"`                         // 动作参数
}

// EmbodiedAIStateResponse AI状态响应
type EmbodiedAIStateResponse struct {
	StateKey      string   `json:"state_key"`
	EntityID      string   `json:"entity_id"`
	EntityType    string   `json:"entity_type"`
	Position      SpatialPosition `json:"position"`
	Orientation   float64  `json:"orientation"`
	LocationName  string   `json:"location_name"`
	PerceptionMode string  `json:"perception_mode"`
	AlertLevel    string   `json:"alert_level"`
	EnergyLevel   float64  `json:"energy_level"`
	Mood          string   `json:"mood"`
	MoodIntensity float64  `json:"mood_intensity"`
	CurrentGoal   string   `json:"current_goal"`
	GoalProgress  float64  `json:"goal_progress"`
	Capabilities  string   `json:"capabilities"`
	LastActiveAt  string   `json:"last_active_at"`
}

// SpatialPosition 位置信息
type SpatialPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// CapabilitiesResponse 能力列表响应
type CapabilitiesResponse struct {
	EntityID   string       `json:"entity_id"`
	EntityType string      `json:"entity_type"`
	Capabilities []Capability `json:"capabilities"`
}

// Capability 单个能力
type Capability struct {
	Key         string  `json:"key"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Available   bool    `json:"available"`
	Confidence  float64 `json:"confidence"`
}

// PerceiveResponse 感知响应
type PerceiveResponse struct {
	EntityID       string   `json:"entity_id"`
	Timestamp      string   `json:"timestamp"`
	Objects        []string `json:"objects"`        // 检测到的物体
	Events         []string `json:"events"`         // 发生的事件
	Entities       []string `json:"entities"`       // 检测到的实体
	EnvironmentTags []string `json:"environment_tags"`
	SaliencyMap    string   `json:"saliency_map"`   // 显著图
	Confidence      float64  `json:"confidence"`
}

// SpatialResponse 空间认知响应
type SpatialResponse struct {
	EntityID      string   `json:"entity_id"`
	QueryType     string   `json:"query_type"`
	Position      SpatialPosition `json:"position,omitempty"`
	Orientation   float64  `json:"orientation,omitempty"`
	Distance      float64  `json:"distance,omitempty"`
	Direction     float64  `json:"direction,omitempty"`
	NearbyObjects []string `json:"nearby_objects,omitempty"`
	MapData       string   `json:"map_data,omitempty"`
}

// ExploreResponse 探索响应
type ExploreResponse struct {
	SessionKey    string   `json:"session_key"`
	Strategy      string   `json:"strategy"`
	Status        string   `json:"status"`
	Progress      float64  `json:"progress"`
	Waypoints     []string `json:"waypoints,omitempty"`
	DiscoveredObjects []string `json:"discovered_objects,omitempty"`
	EstimatedTime int      `json:"estimated_time_remaining,omitempty"`
	CoverageRate  float64  `json:"coverage_rate,omitempty"`
}

// InteractResponse 交互响应
type InteractResponse struct {
	EntityID    string  `json:"entity_id"`
	ActionType  string  `json:"action_type"`
	Success     bool    `json:"success"`
	NewPosition SpatialPosition `json:"new_position,omitempty"`
	Result      string  `json:"result"`
	Feedback    string  `json:"feedback"`
	EnergyUsed  float64 `json:"energy_used"`
}

// ExploreListQuery 探索会话列表查询
type ExploreListQuery struct {
	EntityID  string `form:"entity_id"`
	Status   string `form:"status"`
	Strategy string `form:"strategy"`
}
