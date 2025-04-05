-- 用户表，存储用户基本信息
CREATE TABLE users (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID，主键，自增',
                       username VARCHAR(255) UNIQUE NOT NULL COMMENT '用户名，最多 10 个中文字符或 30 个英文字符',
                       email VARCHAR(255) UNIQUE NOT NULL COMMENT '用户邮箱，唯一，最大长度 100',
                       phone VARCHAR(255) UNIQUE COMMENT '手机号，唯一，11 位数字，可为空',
                       password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希值，存储加密后的密码',
                       avatar_url TEXT COMMENT '用户头像 URL，可为空',
                       bio TEXT COMMENT '个人简介，可为空',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '账号创建时间',
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '账号更新时间',
                       deleted_at TIMESTAMP NULL COMMENT '账号删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户基本信息表';

-- 图片表，存储用户上传的图片信息
CREATE TABLE `images` (
                          `image_id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '图片ID，主键，自增',
                          `uid` BIGINT NOT NULL COMMENT '用户ID，关联 users 表，标识图片的上传用户',
                          `url` VARCHAR(255) NOT NULL COMMENT '图片地址，存储图片在存储服务上的访问URL',
                          `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '图片上传时间',
                          `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '图片信息更新时间',
                          `deleted_at` TIMESTAMP NULL COMMENT '图片逻辑删除时间，若为NULL表示图片未被删除',
                          FOREIGN KEY (`uid`) REFERENCES `users`(`id`) ON DELETE CASCADE,
                          INDEX idx_images_uid (`uid`),
                          INDEX idx_images_created_at (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户图片信息表';

-- 会员表，存储用户的会员信息
CREATE TABLE memberships (
                             id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '会员ID，主键，自增',
                             user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                             membership_level ENUM('free', 'silver', 'gold', 'platinum') DEFAULT 'free' COMMENT '会员等级',
                             status ENUM('active', 'expired', 'pending') DEFAULT 'active' COMMENT '会员状态',
                             expires_at TIMESTAMP COMMENT '会员到期时间',
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '会员激活时间',
                             FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                             INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户会员信息表';

-- 积分表，存储用户的积分信息
CREATE TABLE points (
                        id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID，主键，自增',
                        user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                        points INT UNSIGNED DEFAULT 0 COMMENT '当前可用积分，不可为负数',
                        total_points INT UNSIGNED DEFAULT 0 COMMENT '累计获得的总积分',
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
                        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                        INDEX idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分表';

-- 关注/拉黑表
CREATE TABLE relationships (
                               id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '关系ID，主键，自增',
                               user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                               target_id BIGINT NOT NULL COMMENT '目标用户ID，关联 users 表',
                               status ENUM('follow', 'block', 'mute') NOT NULL DEFAULT 'follow' COMMENT '关系状态：关注、拉黑、静音',
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                               FOREIGN KEY (target_id) REFERENCES users(id) ON DELETE CASCADE,
                               UNIQUE KEY uniq_user_target (user_id, target_id),
                               INDEX idx_relationships_target (target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户关系表，支持关注和拉黑';

-- 用户操作日志表
CREATE TABLE user_activity_logs (
                                    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID，主键，自增',
                                    user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                                    action VARCHAR(255) NOT NULL COMMENT '操作类型，如登录、发布帖子、点赞',
                                    target_id BIGINT UNSIGNED COMMENT '操作目标ID，如帖子ID、评论ID，可为空',
                                    target_type ENUM('post', 'comment', 'reaction', 'message', 'other') NOT NULL COMMENT '目标类型',
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
                                    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                    INDEX idx_target_id (target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户活动日志表';

-- 索引优化
CREATE INDEX idx_logs_user ON user_activity_logs(user_id);
CREATE INDEX idx_relationships_user ON relationships(user_id, target_id);
CREATE INDEX idx_memberships_user ON memberships(user_id);
CREATE INDEX idx_points_user ON points(user_id);