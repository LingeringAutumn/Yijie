-- 用户表，存储用户基本信息
CREATE TABLE users (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID，主键，自增',
                       username VARCHAR(30) UNIQUE NOT NULL COMMENT '用户名，唯一，最多 10 个中文字符或 30 个英文字符',
                       email VARCHAR(100) UNIQUE NOT NULL COMMENT '用户邮箱，唯一，最大长度 100',
                       phone VARCHAR(11) UNIQUE NOT NULL COMMENT '手机号，唯一，11 位数字',
                       password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希值，存储加密后的密码，原密码应为 8-16 位字母+数字组合',
                       avatar_url TEXT COMMENT '用户头像 URL，可为空',
                       bio TEXT COMMENT '个人简介，可为空',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '账号创建时间',
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '账号更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户基本信息表';

-- 团队表，存储团队的基本信息
CREATE TABLE teams (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '团队ID，主键，自增',
                       name VARCHAR(50) UNIQUE NOT NULL COMMENT '团队名称，唯一，最多 50 个字符',
                       description TEXT COMMENT '团队描述，可为空',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '团队创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队基本信息表';

-- 团队成员表，存储用户与团队的关系
CREATE TABLE team_members (
                              id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '团队成员ID，主键，自增',
                              team_id BIGINT NOT NULL COMMENT '团队ID，关联 teams 表',
                              user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                              role ENUM('admin', 'member') DEFAULT 'member' COMMENT '团队角色，默认为普通成员',
                              joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
                              FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
                              FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队成员表';

-- 会员表，存储用户的会员信息
CREATE TABLE memberships (
                             id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '会员ID，主键，自增',
                             user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                             membership_level ENUM('free', 'silver', 'gold', 'platinum') DEFAULT 'free' COMMENT '会员等级，默认为 free',
                             expire_at TIMESTAMP COMMENT '会员到期时间，可为空',
                             FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户会员信息表';

-- 积分表，存储用户的积分信息
CREATE TABLE points (
                        id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID，主键，自增',
                        user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                        points INT DEFAULT 0 COMMENT '用户积分，默认为 0，不可为负数',
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
                        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分表';

-- 团队项目表，存储团队创建的 AI 生成前端项目
CREATE TABLE team_projects (
                               id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '团队项目ID，主键，自增',
                               team_id BIGINT NOT NULL COMMENT '团队ID，关联 teams 表',
                               project_name VARCHAR(100) NOT NULL COMMENT '项目名称，最多 100 个字符',
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                               FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队项目表';

-- 团队项目权限表，存储团队成员对项目的权限
CREATE TABLE team_permissions (
                                  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '权限ID，主键，自增',
                                  team_project_id BIGINT NOT NULL COMMENT '团队项目ID，关联 team_projects 表',
                                  user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                                  permission ENUM('read', 'write', 'admin') DEFAULT 'read' COMMENT '权限级别，默认为只读',
                                  FOREIGN KEY (team_project_id) REFERENCES team_projects(id) ON DELETE CASCADE,
                                  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队项目权限表';

-- AI 生成的前端代码表
CREATE TABLE generated_interfaces (
                                      id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '生成界面ID，主键，自增',
                                      user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                                      team_project_id BIGINT COMMENT '团队项目ID，关联 team_projects 表，可为空',
                                      interface_name VARCHAR(100) NOT NULL COMMENT '生成的界面名称，最长 100 个字符',
                                      code TEXT NOT NULL COMMENT '生成的前端代码，不可为空',
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                      FOREIGN KEY (team_project_id) REFERENCES team_projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI 生成的前端代码存储表';
