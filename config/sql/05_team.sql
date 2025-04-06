-- 团队表，存储团队基本信息
CREATE TABLE teams (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '团队ID，主键，自增',
                       name VARCHAR(255) UNIQUE NOT NULL COMMENT '团队名称，唯一',
                       description TEXT COMMENT '团队简介，可为空',
                       creator_id BIGINT NOT NULL COMMENT '创建者用户ID，关联 users 表',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '团队创建时间',
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                       deleted_at TIMESTAMP DEFAULT NULL COMMENT '删除时间',
                       FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队基本信息表';

-- 团队成员表
CREATE TABLE team_members (
                              id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '团队成员ID，主键，自增',
                              team_id BIGINT NOT NULL COMMENT '团队ID，关联 teams 表',
                              user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                              role ENUM('admin', 'member') NOT NULL DEFAULT 'member' COMMENT '团队角色',
                              status ENUM('active', 'pending', 'removed') DEFAULT 'active' COMMENT '成员状态，pending 代表待审核',
                              joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
                              FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
                              FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队成员表';

-- 索引优化
CREATE INDEX idx_teams_creator ON teams(creator_id);
CREATE INDEX idx_team_members_team ON team_members(team_id);
CREATE INDEX idx_team_members_user ON team_members(user_id);