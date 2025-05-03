-- 视频表，存储视频基本信息
CREATE TABLE videos (
                        video_id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '视频ID，主键，自增',
                        user_id BIGINT NOT NULL COMMENT '作者用户ID，关联 users 表',
                        title VARCHAR(255) NOT NULL COMMENT '视频标题',
                        description TEXT COMMENT '视频描述，可为空',
                        cover_url VARCHAR(255) COMMENT '封面图地址，可为空',
                        video_url VARCHAR(255) NOT NULL COMMENT '视频文件URL',
                        duration_seconds INT UNSIGNED COMMENT '视频时长（单位：秒）',
                        status ENUM('published', 'deleted', 'draft') DEFAULT 'published' COMMENT '视频状态',
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        deleted_at TIMESTAMP NULL COMMENT '逻辑删除时间',
                        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                        INDEX idx_user_id (user_id),
                        INDEX idx_status (status),
                        INDEX idx_created_at (created_at),
                        FULLTEXT INDEX idx_title_description (title, description)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频内容表';

-- 视频统计表，记录播放量、点赞、评论等
CREATE TABLE video_stats (
                             video_id BIGINT PRIMARY KEY COMMENT '视频ID，关联 videos 表',
                             views BIGINT UNSIGNED DEFAULT 0 COMMENT '播放次数',
                             likes BIGINT UNSIGNED DEFAULT 0 COMMENT '点赞次数',
                             comments BIGINT UNSIGNED DEFAULT 0 COMMENT '评论数',
                             hot_score DOUBLE DEFAULT 0 COMMENT '热度评分，用于热榜',
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
                             FOREIGN KEY (video_id) REFERENCES videos(video_id) ON DELETE CASCADE,
                             INDEX idx_hot_score (hot_score),
                             INDEX idx_views (views),
                             INDEX idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频统计表';

-- 视频标签表
CREATE TABLE video_tags (
                            tag_id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '标签ID',
                            video_id BIGINT NOT NULL COMMENT '视频ID',
                            tag VARCHAR(64) NOT NULL COMMENT '标签内容',
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '标签创建时间',
                            FOREIGN KEY (video_id) REFERENCES videos(video_id) ON DELETE CASCADE,
                            INDEX idx_tag (tag),
                            UNIQUE KEY uniq_video_tag (video_id, tag)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频标签表';

