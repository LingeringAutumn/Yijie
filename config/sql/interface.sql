-- AI 生成的前端代码表
CREATE TABLE generated_interfaces (
                                      id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '生成界面ID，主键，自增',
                                      user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                                      interface_name VARCHAR(100) NOT NULL COMMENT '生成的界面名称，最长 100 个字符',
                                      preset_id VARCHAR(10) COMMENT '预设场景ID，例如 0001 表示音乐播放器，0002 表示视频播放器，可为空',
                                      description TEXT COMMENT '界面描述信息，存储用户对该界面的备注或用途说明，可为空',
                                      code TEXT NOT NULL COMMENT '生成的前端代码，不可为空',
                                      config_json JSON COMMENT '前端可视化配置信息，例如用户修改的 UI 布局参数，可为空',
                                      is_public BOOLEAN DEFAULT FALSE COMMENT '是否公开该界面，默认不公开',
                                      views INT DEFAULT 0 COMMENT '浏览次数，默认 0，每次被查看时递增',
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI 生成的前端代码存储表';
