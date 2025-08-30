-- +migrate Up
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED DEFAULT NULL COMMENT '操作用户ID',
    username VARCHAR(100) NOT NULL DEFAULT '' COMMENT '用户名',
    action VARCHAR(100) NOT NULL DEFAULT '' COMMENT '操作类型',
    resource VARCHAR(100) NOT NULL DEFAULT '' COMMENT '资源类型',
    method VARCHAR(20) NOT NULL DEFAULT '' COMMENT 'HTTP方法',
    path VARCHAR(500) NOT NULL DEFAULT '' COMMENT '请求路径',
    ip_address VARCHAR(45) NOT NULL DEFAULT '' COMMENT '客户端IP',
    user_agent VARCHAR(500) NOT NULL DEFAULT '' COMMENT '用户代理',
    request_body TEXT COMMENT '请求体',
    response_body TEXT COMMENT '响应体',
    status_code INT NOT NULL DEFAULT 0 COMMENT '响应状态码',
    duration BIGINT NOT NULL DEFAULT 0 COMMENT '请求耗时(纳秒)',
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_resource (resource),
    INDEX idx_ip_address (ip_address),
    INDEX idx_status_code (status_code),
    INDEX idx_timestamp (timestamp),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';

-- +migrate Down
DROP TABLE IF EXISTS audit_logs;