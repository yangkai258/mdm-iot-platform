-- 组织管理相关表：companies, departments, positions, employees
-- 这些表在模型中定义但数据库中缺失

-- 公司表
CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(50),
    company_code VARCHAR(50) UNIQUE NOT NULL,
    company_name VARCHAR(200) NOT NULL,
    short_name VARCHAR(100),
    logo VARCHAR(500),
    province VARCHAR(50),
    city VARCHAR(50),
    district VARCHAR(50),
    address VARCHAR(500),
    legal_person VARCHAR(50),
    contact VARCHAR(50),
    phone VARCHAR(20),
    email VARCHAR(100),
    status INT DEFAULT 1,
    sort INT DEFAULT 0,
    remark VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_companies_tenant ON companies(tenant_id);
CREATE INDEX IF NOT EXISTS idx_companies_deleted ON companies(deleted_at);

-- 部门表
CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    tenant_id UUID,
    dept_code VARCHAR(50) NOT NULL,
    dept_name VARCHAR(100) NOT NULL,
    parent_id INT,
    level INT DEFAULT 1,
    path VARCHAR(500),
    manager_id INT,
    phone VARCHAR(20),
    email VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    sort_order INT DEFAULT 0,
    company_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_departments_tenant ON departments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_departments_parent ON departments(parent_id);
CREATE INDEX IF NOT EXISTS idx_departments_deleted ON departments(deleted_at);

-- 岗位表
CREATE TABLE IF NOT EXISTS positions (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(50),
    pos_code VARCHAR(50) NOT NULL,
    pos_name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    level INT DEFAULT 1,
    dept_id INT,
    company_id INT,
    description VARCHAR(500),
    status INT DEFAULT 1,
    sort INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_positions_tenant ON positions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_positions_dept ON positions(dept_id);
CREATE INDEX IF NOT EXISTS idx_positions_deleted ON positions(deleted_at);

-- 员工表
CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(50),
    emp_code VARCHAR(50) UNIQUE NOT NULL,
    emp_name VARCHAR(50) NOT NULL,
    gender VARCHAR(10),
    birth_date TIMESTAMP,
    phone VARCHAR(20),
    email VARCHAR(100),
    id_card VARCHAR(20),
    photo VARCHAR(500),
    province VARCHAR(50),
    city VARCHAR(50),
    district VARCHAR(50),
    address VARCHAR(500),
    dept_id INT,
    position_id INT,
    company_id INT,
    entry_date TIMESTAMP,
    emp_status INT DEFAULT 1,
    status INT DEFAULT 1,
    remark VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_employees_tenant ON employees(tenant_id);
CREATE INDEX IF NOT EXISTS idx_employees_dept ON employees(dept_id);
CREATE INDEX IF NOT EXISTS idx_employees_position ON employees(position_id);
CREATE INDEX IF NOT EXISTS idx_employees_deleted ON employees(deleted_at);

-- 插入示例数据
INSERT INTO companies (tenant_id, company_code, company_name, short_name, status) VALUES
('e6cbcb82-9bd6-4803-8bf7-b4b1af8eaec2', 'COMP001', '示例科技有限公司', '示例科技', 1)
ON CONFLICT (company_code) DO NOTHING;

INSERT INTO departments (tenant_id, dept_code, dept_name, level, status) VALUES
('e6cbcb82-9bd6-4803-8bf7-b4b1af8eaec2', 'DEPT001', '技术部', 1, 'active'),
('e6cbcb82-9bd6-4803-8bf7-b4b1af8eaec2', 'DEPT002', '运营部', 1, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO positions (tenant_id, pos_code, pos_name, category, level, status) VALUES
('e6cbcb82-9bd6-4803-8bf7-b4b1af8eaec2', 'POS001', '高级工程师', '技术', 3, 1),
('e6cbcb82-9bd6-4803-8bf7-b4b1af8eaec2', 'POS002', '产品经理', '产品', 2, 1),
('e6cbcb82-9bd6-4803-8bf7-b4b1af8eaec2', 'POS003', '运营专员', '运营', 1, 1)
ON CONFLICT DO NOTHING;
