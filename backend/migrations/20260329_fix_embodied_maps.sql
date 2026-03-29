-- 修复 embodied_maps 表缺失的 map_type 列
-- 问题: GORM 模型定义 MapType 为 NOT NULL，但数据库表中该列不存在
-- 解决: 添加列并设置默认值

-- 检查列是否存在，如果不存在则添加
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'embodied_maps' AND column_name = 'map_type'
    ) THEN
        ALTER TABLE embodied_maps ADD COLUMN map_type varchar(20) DEFAULT 'grid';
        RAISE NOTICE 'Column map_type added to embodied_maps';
    ELSE
        RAISE NOTICE 'Column map_type already exists in embodied_maps';
    END IF;
END $$;

-- 同样检查 spatial_positions 表的 map_id 列类型
-- 如果 map_id 是 varchar但模型期望 bigint，需要确保兼容性
-- 注意: spatial_positions.map_id 存储的是字符串(如 'map_001')，不是数字ID
