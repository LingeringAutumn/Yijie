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

-- 索引优化
CREATE INDEX idx_posts_user ON posts(user_id, created_at);
CREATE INDEX idx_comments_post ON comments(post_id);
CREATE INDEX idx_comments_parent ON comments(parent_comment_id);
CREATE INDEX idx_reactions_user ON reactions(user_id);
CREATE INDEX idx_reactions_post ON reactions(post_id);
CREATE INDEX idx_reactions_comment ON reactions(comment_id);