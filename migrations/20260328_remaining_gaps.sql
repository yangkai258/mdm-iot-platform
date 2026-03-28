-- 设备影子快照表
CREATE TABLE IF NOT EXISTS device_shadow_snapshots (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(64) NOT NULL,
    snapshot_id VARCHAR(64) UNIQUE NOT NULL,
    version INTEGER DEFAULT 1,
    snapshot_type VARCHAR(20) DEFAULT 'manual',
    reason VARCHAR(200),
    desired_state TEXT,
    desired_version INTEGER DEFAULT 0,
    reported_state TEXT,
    reported_version INTEGER DEFAULT 0,
    metadata TEXT,
    delta TEXT,
    tags VARCHAR(500),
    state_diff INTEGER DEFAULT 0,
    is_healthy BOOLEAN DEFAULT TRUE,
    file_url VARCHAR(512),
    file_size BIGINT DEFAULT 0,
    checksum VARCHAR(64),
    expires_at TIMESTAMP,
    created_by VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dss_device ON device_shadow_snapshots(device_id);

-- 设备影子快照导出记录表
CREATE TABLE IF NOT EXISTS device_shadow_snapshot_exports (
    id BIGSERIAL PRIMARY KEY,
    snapshot_id VARCHAR(64) NOT NULL,
    device_id VARCHAR(64) NOT NULL,
    format VARCHAR(20),
    file_url VARCHAR(512),
    file_size BIGINT DEFAULT 0,
    download_count INTEGER DEFAULT 0,
    created_by VARCHAR(64),
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dsse_snapshot ON device_shadow_snapshot_exports(snapshot_id);

-- OTA固件兼容性矩阵表
CREATE TABLE IF NOT EXISTS ota_compatibility_matrix (
    id BIGSERIAL PRIMARY KEY,
    matrix_id VARCHAR(64) UNIQUE NOT NULL,
    hardware_model VARCHAR(64) NOT NULL,
    hardware_version VARCHAR(50),
    from_firmware VARCHAR(50) NOT NULL,
    to_firmware VARCHAR(50) NOT NULL,
    compatibility_status VARCHAR(20) DEFAULT 'compatible',
    compatibility_score DECIMAL(5,2) DEFAULT 100,
    min_battery_level INTEGER DEFAULT 30,
    min_storage_kb INTEGER DEFAULT 1024,
    min_memory_kb INTEGER DEFAULT 512,
    network_required BOOLEAN DEFAULT FALSE,
    constraints TEXT,
    warning TEXT,
    breaking_changes TEXT,
    rollback_supported BOOLEAN DEFAULT TRUE,
    rollback_min_version VARCHAR(50),
    success_count INTEGER DEFAULT 0,
    failure_count INTEGER DEFAULT 0,
    success_rate DECIMAL(5,2) DEFAULT 0,
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP,
    verified_by VARCHAR(64),
    is_active BOOLEAN DEFAULT TRUE,
    created_by VARCHAR(64),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ocm_hardware ON ota_compatibility_matrix(hardware_model);
CREATE INDEX IF NOT EXISTS idx_ocm_firmware ON ota_compatibility_matrix(from_firmware, to_firmware);

-- OTA兼容性测试记录表
CREATE TABLE IF NOT EXISTS ota_compatibility_tests (
    id BIGSERIAL PRIMARY KEY,
    matrix_id VARCHAR(64) NOT NULL,
    device_id VARCHAR(64) NOT NULL,
    from_firmware VARCHAR(50) NOT NULL,
    to_firmware VARCHAR(50) NOT NULL,
    test_result VARCHAR(20),
    error_code VARCHAR(50),
    error_message TEXT,
    duration INTEGER DEFAULT 0,
    environment TEXT,
    tester_id VARCHAR(64),
    test_type VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_oct_matrix ON ota_compatibility_tests(matrix_id);

-- SDK包表
CREATE TABLE IF NOT EXISTS sdk_packages (
    id BIGSERIAL PRIMARY KEY,
    sdk_id VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(200),
    description TEXT,
    category VARCHAR(50),
    platform VARCHAR(50),
    language VARCHAR(30),
    current_version VARCHAR(50),
    version_count INTEGER DEFAULT 0,
    download_count INTEGER DEFAULT 0,
    install_count INTEGER DEFAULT 0,
    star_count INTEGER DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    developer_id VARCHAR(64),
    developer_name VARCHAR(100),
    status VARCHAR(20) DEFAULT 'draft',
    is_official BOOLEAN DEFAULT FALSE,
    is_featured BOOLEAN DEFAULT FALSE,
    tags VARCHAR(500),
    category_tags VARCHAR(200),
    doc_url VARCHAR(512),
    github_url VARCHAR(512),
    changelog_url VARCHAR(512),
    icon_url VARCHAR(512),
    screenshots TEXT,
    license VARCHAR(50),
    review_status VARCHAR(20) DEFAULT 'pending',
    review_comment TEXT,
    reviewed_at TIMESTAMP,
    reviewed_by VARCHAR(64),
    created_by VARCHAR(64),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sdk_name ON sdk_packages(name);
CREATE INDEX IF NOT EXISTS idx_sdk_category ON sdk_packages(category);

-- SDK版本表
CREATE TABLE IF NOT EXISTS sdk_versions (
    id BIGSERIAL PRIMARY KEY,
    sdk_id VARCHAR(64) NOT NULL,
    version VARCHAR(50) NOT NULL,
    file_url VARCHAR(512),
    file_size BIGINT DEFAULT 0,
    file_hash VARCHAR(64),
    min_platform_version VARCHAR(50),
    release_date TIMESTAMP,
    release_notes TEXT,
    compatible_platforms VARCHAR(200),
    dependencies TEXT,
    status VARCHAR(20) DEFAULT 'stable',
    is_recommended BOOLEAN DEFAULT FALSE,
    download_count INTEGER DEFAULT 0,
    created_by VARCHAR(64),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sdkv_sdk ON sdk_versions(sdk_id);

-- SDK下载记录表
CREATE TABLE IF NOT EXISTS sdk_downloads (
    id BIGSERIAL PRIMARY KEY,
    sdk_id VARCHAR(64) NOT NULL,
    version_id BIGINT NOT NULL,
    version VARCHAR(50) NOT NULL,
    user_id VARCHAR(64),
    user_name VARCHAR(100),
    project_id VARCHAR(64),
    ip_address VARCHAR(50),
    user_agent VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sdkd_sdk ON sdk_downloads(sdk_id);
