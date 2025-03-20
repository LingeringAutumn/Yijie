-- AI 生成的前端代码表
CREATE TABLE generated_interfaces (
                                      id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '生成界面ID，主键，自增',
                                      user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                                      interface_name VARCHAR(100) NOT NULL COMMENT '生成的界面名称，最长 100 个字符',
                                      code TEXT NOT NULL COMMENT '生成的前端代码，不可为空',
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI 生成的前端代码存储表';