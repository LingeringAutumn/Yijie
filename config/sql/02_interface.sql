-- AI 生成界面存储表
CREATE TABLE generated_interfaces (
                                      id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '生成界面ID，主键，自增',
                                      user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                                      interface_name VARCHAR(255) NOT NULL COMMENT '生成的界面名称',
                                      preset_id VARCHAR(255) COMMENT '预设场景ID，可为空',
                                      description TEXT COMMENT '界面描述信息，可为空',
                                      code LONGTEXT NOT NULL COMMENT '生成的前端代码，不可为空',
                                      config_json MEDIUMTEXT COMMENT '前端可视化配置信息，可为空',
                                      visibility ENUM('private', 'public', 'unlisted') DEFAULT 'private' COMMENT '界面可见性',
                                      views INT DEFAULT 0 COMMENT '浏览次数',
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                      deleted_at TIMESTAMP DEFAULT NULL COMMENT '删除时间',
                                      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI 生成的前端代码存储表';

-- 索引优化
CREATE INDEX idx_generated_interfaces_user ON generated_interfaces(user_id);