-- 社区帖子表，用户发布的帖子
CREATE TABLE posts (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '帖子ID，主键，自增',
                       user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                       content TEXT NOT NULL COMMENT '帖子内容，不能为空',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户帖子表';

-- 社区评论表，存储用户对帖子或评论的评论
CREATE TABLE comments (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评论ID，主键，自增',
                          post_id BIGINT COMMENT '帖子ID，关联 posts 表，可为空（如果是对评论的回复）',
                          parent_comment_id BIGINT COMMENT '父评论ID，关联 comments 表，可为空（如果是对帖子评论）',
                          user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                          content TEXT NOT NULL COMMENT '评论内容，不能为空',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
                          FOREIGN KEY (parent_comment_id) REFERENCES comments(id) ON DELETE CASCADE,
                          FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子和评论的评论表';

-- 点赞表，存储用户对帖子或评论的点赞
CREATE TABLE likes (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '点赞ID，主键，自增',
                       user_id BIGINT NOT NULL COMMENT '用户ID，关联 users 表',
                       post_id BIGINT COMMENT '帖子ID，关联 posts 表，可为空（如果是对评论点赞）',
                       comment_id BIGINT COMMENT '评论ID，关联 comments 表，可为空（如果是对帖子点赞）',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',
                       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                       FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
                       FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子和评论的点赞表';