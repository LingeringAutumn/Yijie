-- 视频点赞表（用户行为）
-- 如果系统后期添加更多用户行为（如收藏、转发等），建议拆成 user_behavior.sql 单独维护
CREATE TABLE video_likes (
                             user_id BIGINT NOT NULL COMMENT '用户ID',
                             video_id BIGINT NOT NULL COMMENT '视频ID',
                             is_liked BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否点赞，支持取消点赞逻辑',
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             PRIMARY KEY (user_id, video_id),
                             FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                             FOREIGN KEY (video_id) REFERENCES videos(video_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频点赞行为表';