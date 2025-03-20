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