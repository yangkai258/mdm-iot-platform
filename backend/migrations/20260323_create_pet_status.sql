-- ============================================================
-- Sprint 9: 宠物状态表
-- 创建时间: 2026-03-23
-- 描述: 存储宠物的实时状态信息
-- ============================================================

CREATE TABLE IF NOT EXISTS pet_status (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(64) NOT NULL UNIQUE,
    pet_name VARCHAR(32) NOT NULL DEFAULT '小爪',
    pet_type VARCHAR(16) DEFAULT 'cat',
    personality JSONB DEFAULT '{}',
    appearance JSONB DEFAULT '{}',
    mood INT DEFAULT 50,
    energy INT DEFAULT 100,
    hunger INT DEFAULT 0,
    position_x FLOAT DEFAULT 0,
    position_y FLOAT DEFAULT 0,
    current_expression VARCHAR(32) DEFAULT 'happy',
    current_action VARCHAR(32),
    is_online BOOLEAN DEFAULT FALSE,
    last_seen_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pet_status_device_id ON pet_status(device_id);

COMMENT ON TABLE pet_status IS '宠物状态表';
COMMENT ON COLUMN pet_status.device_id IS '设备ID';
COMMENT ON COLUMN pet_status.pet_name IS '宠物名称';
COMMENT ON COLUMN pet_status.pet_type IS '宠物类型: cat, dog, rabbit等';
COMMENT ON COLUMN pet_status.personality IS '性格特征JSON';
COMMENT ON COLUMN pet_status.appearance IS '外观特征JSON';
COMMENT ON COLUMN pet_status.mood IS '心情值 0-100';
COMMENT ON COLUMN pet_status.energy IS '能量值 0-100';
COMMENT ON COLUMN pet_status.hunger IS '饥饿度 0-100';
COMMENT ON COLUMN pet_status.position_x IS 'X坐标位置';
COMMENT ON COLUMN pet_status.position_y IS 'Y坐标位置';
COMMENT ON COLUMN pet_status.current_expression IS '当前表情';
COMMENT ON COLUMN pet_status.current_action IS '当前动作';
COMMENT ON COLUMN pet_status.is_online IS '是否在线';
