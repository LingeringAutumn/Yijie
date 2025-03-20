-- 团队表，存储团队基本信息
CREATE TABLE teams (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '团队ID，主键，自增',
                       name VARCHAR(100) UNIQUE NOT NULL COMMENT '团队名称，唯一，最多 100 个字符',
                       description TEXT COMMENT '团队简介，可为空',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '团队创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队基本信息表';

-- 团队成员表
CREATE TABLE team_members (
                              id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '团队成员ID，主键，自增',
                              team_id BIGINT NOT NULL COMMENT '团队ID，关联 teams 表',
                              user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                              role ENUM('admin', 'member') NOT NULL DEFAULT 'member' COMMENT '团队角色，admin 为管理员，member 为普通成员',
                              joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
                              FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
                              FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队成员表';
