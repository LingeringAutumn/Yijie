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

-- 索引优化
CREATE INDEX idx_messages_user ON messages(sender_id, receiver_id);