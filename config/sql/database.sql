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
                          `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '图片ID，主键，自增',
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

-- 社区帖子表，用户发布的帖子
CREATE TABLE posts (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '帖子ID，主键，自增',
                       user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                       title VARCHAR(255) NOT NULL COMMENT '帖子标题，最多 255 个字符',
                       content TEXT NOT NULL COMMENT '帖子内容，不能为空',
                       image_urls JSON COMMENT '帖子图片 URL 数组，可为空',
                       visibility ENUM('public', 'friends', 'private') DEFAULT 'public' COMMENT '可见性：公开、好友可见、私密',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                       deleted_at TIMESTAMP NULL COMMENT '删除时间',
                       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户帖子表';

-- 社区评论表，存储用户对帖子或评论的评论
CREATE TABLE comments (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评论ID，主键，自增',
                          post_id BIGINT COMMENT '帖子ID，关联 posts 表，可为空',
                          parent_comment_id BIGINT COMMENT '父评论ID，关联 comments 表，可为空',
                          user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                          content TEXT NOT NULL COMMENT '评论内容，不能为空',
                          status ENUM('normal', 'hidden') DEFAULT 'normal' COMMENT '评论状态',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          deleted_at TIMESTAMP NULL COMMENT '删除时间',
                          FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
                          FOREIGN KEY (parent_comment_id) REFERENCES comments(id) ON DELETE SET NULL,
                          FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子和评论的评论表';

-- 反应表（取代 likes，支持不同的点赞类型）
CREATE TABLE reactions (
                           id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '反应ID，主键，自增',
                           user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                           post_id BIGINT COMMENT '帖子ID，关联 posts 表，可为空',
                           comment_id BIGINT COMMENT '评论ID，关联 comments 表，可为空',
                           reaction_type ENUM('like', 'love', 'haha', 'wow', 'sad', 'angry') NOT NULL COMMENT '反应类型',
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           deleted_at TIMESTAMP NULL COMMENT '删除时间',
                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                           FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
                           FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
                           CHECK (post_id IS NOT NULL OR comment_id IS NOT NULL)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户反应表';

-- 私信消息表
CREATE TABLE messages (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '消息ID，主键，自增',
                          sender_id BIGINT NOT NULL COMMENT '发送者ID',
                          receiver_id BIGINT NOT NULL COMMENT '接收者ID',
                          content TEXT NOT NULL COMMENT '消息内容',
                          status ENUM('unread', 'read', 'retracted') DEFAULT 'unread' COMMENT '消息状态：未读、已读、撤回',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
                          FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
                          FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='私信消息表';

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
CREATE INDEX idx_posts_user ON posts(user_id, created_at);
CREATE INDEX idx_comments_post ON comments(post_id);
CREATE INDEX idx_comments_parent ON comments(parent_comment_id);
CREATE INDEX idx_reactions_user ON reactions(user_id);
CREATE INDEX idx_reactions_post ON reactions(post_id);
CREATE INDEX idx_reactions_comment ON reactions(comment_id);
CREATE INDEX idx_logs_user ON user_activity_logs(user_id);
CREATE INDEX idx_relationships_user ON relationships(user_id, target_id);
CREATE INDEX idx_messages_user ON messages(sender_id, receiver_id);
CREATE INDEX idx_memberships_user ON memberships(user_id);
CREATE INDEX idx_points_user ON points(user_id);
CREATE INDEX idx_teams_creator ON teams(creator_id);
CREATE INDEX idx_team_members_team ON team_members(team_id);
CREATE INDEX idx_team_members_user ON team_members(user_id);
CREATE INDEX idx_generated_interfaces_user ON generated_interfaces(user_id);