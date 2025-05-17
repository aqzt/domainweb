-- 创建数据库
CREATE DATABASE IF NOT EXISTS domainweb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE domainweb;

-- 创建域名表
CREATE TABLE IF NOT EXISTS domains (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL COMMENT '完整域名',
    tld VARCHAR(50) NOT NULL COMMENT '顶级域名',
    length INT NOT NULL COMMENT '域名长度（不含TLD）',
    structure VARCHAR(50) NOT NULL COMMENT '域名结构',
    register_date DATETIME COMMENT '注册日期',
    expire_date DATETIME COMMENT '到期日期',
    created_at DATETIME NOT NULL COMMENT '记录创建时间',
    updated_at DATETIME NOT NULL COMMENT '记录更新时间',
    UNIQUE KEY idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='域名基本信息表';

-- 创建域名属性表
CREATE TABLE IF NOT EXISTS domain_attributes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    attribute_name VARCHAR(100) NOT NULL COMMENT '属性名称',
    attribute_type VARCHAR(50) NOT NULL COMMENT '属性类型',
    price_factor DECIMAL(10, 2) NOT NULL COMMENT '估价倍数',
    grade_factor DECIMAL(10, 2) NOT NULL COMMENT '等级增量',
    attribute_value VARCHAR(255) NOT NULL COMMENT '属性值',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
    INDEX idx_attribute_type (attribute_type),
    INDEX idx_attribute_name (attribute_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='域名属性表';

-- 创建查询历史表
CREATE TABLE IF NOT EXISTS history_records (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    domain VARCHAR(255) NOT NULL COMMENT '查询的域名',
    grade DECIMAL(10, 2) NOT NULL COMMENT '品相等级',
    price DECIMAL(10, 2) NOT NULL COMMENT '估价结果',
    estimation_date DATETIME NOT NULL COMMENT '查询时间',
    INDEX idx_domain (domain),
    INDEX idx_estimation_date (estimation_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='查询历史记录表';

-- 插入基础属性数据
INSERT INTO domain_attributes (attribute_name, attribute_type, price_factor, grade_factor, attribute_value, created_at, updated_at) VALUES
-- TLD属性
('com后缀', '基础属性', 9.55, 0.5, 'com', NOW(), NOW()),
('net后缀', '基础属性', 2.38, 0.2, 'net', NOW(), NOW()),
('org后缀', '基础属性', 1.90, 0.1, 'org', NOW(), NOW()),
('cn后缀', '基础属性', 1.45, 0.0, 'cn', NOW(), NOW()),
('com.cn后缀', '基础属性', 1.20, -0.1, 'com.cn', NOW(), NOW()),
('cc后缀', '基础属性', 1.15, -0.2, 'cc', NOW(), NOW()),
('co后缀', '基础属性', 1.10, -0.2, 'co', NOW(), NOW()),
('io后缀', '基础属性', 1.80, 0.1, 'io', NOW(), NOW()),
('ai后缀', '基础属性', 2.50, 0.3, 'ai', NOW(), NOW()),

-- 长度属性
('2位长度', '基础属性', 8.50, 1.0, '2', NOW(), NOW()),
('3位长度', '基础属性', 5.20, 0.8, '3', NOW(), NOW()),
('4位长度', '基础属性', 3.60, 0.62, '4', NOW(), NOW()),
('5位长度', '基础属性', 2.10, 0.4, '5', NOW(), NOW()),
('6位长度', '基础属性', 1.50, 0.2, '6', NOW(), NOW()),
('7位长度', '基础属性', 1.20, 0.1, '7', NOW(), NOW()),
('8位长度', '基础属性', 1.10, 0.05, '8', NOW(), NOW()),
('9位及以上长度', '基础属性', 1.00, 0.0, '9', NOW(), NOW()),

-- 结构属性
('纯数字结构', '基础属性', 1.80, 0.5, '纯数字', NOW(), NOW()),
('纯字母结构', '基础属性', 1.26, 0.31, '纯字母', NOW(), NOW()),
('数字字母混合结构', '基础属性', 1.15, 0.2, '数字字母混合', NOW(), NOW()),
('含连字符结构', '基础属性', 0.85, -0.1, '含连字符', NOW(), NOW()),
('其他结构', '基础属性', 0.75, -0.2, '其他', NOW(), NOW()),

-- 其他属性（示例）
('声母属性', '其他属性', 0.85, 0.0, '声母', NOW(), NOW()),
('Alexa排名', '其他属性', 1.00, 0.0, 'Alexa', NOW(), NOW()),
('相关域名未注册', '其他属性', 0.65, -0.1, '未注册', NOW(), NOW()),
('搜索量', '其他属性', 1.80, 0.6, '搜索量', NOW(), NOW()),
('贴吧数量', '其他属性', 2.25, 0.6, '贴吧', NOW(), NOW()),
('百科系数', '其他属性', 1.30, 0.3, '百科', NOW(), NOW()),
('词典记录', '其他属性', 1.35, 0.3, '词典', NOW(), NOW()),
('360搜索指数', '其他属性', 1.00, 0.0, '360搜索', NOW(), NOW()),
('传媒系数', '其他属性', 3.70, 0.9, '传媒', NOW(), NOW()),
('社交系数', '其他属性', 1.00, 0.0, '社交', NOW(), NOW()),
('淘宝商品数量', '其他属性', 1.18, 0.1, '淘宝', NOW(), NOW());
