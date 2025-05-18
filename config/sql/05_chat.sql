CREATE TABLE chat_rooms (
                            id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '聊天室 ID',
                            is_group BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否为群聊',
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天房间表';
CREATE TABLE chat_room_members (
                                   room_id BIGINT NOT NULL COMMENT '聊天室 ID',
                                   user_id BIGINT NOT NULL COMMENT '成员用户 ID',
                                   PRIMARY KEY (room_id, user_id),
                                   FOREIGN KEY (room_id) REFERENCES chat_rooms(id) ON DELETE CASCADE,
                                   FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                   INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天室成员表';
CREATE TABLE chat_messages (
                               id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '消息 ID',
                               room_id BIGINT NOT NULL COMMENT '聊天室 ID',
                               sender_id BIGINT NOT NULL COMMENT '发送者用户 ID',
                               content TEXT NOT NULL COMMENT '消息内容',
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
                               FOREIGN KEY (room_id) REFERENCES chat_rooms(id) ON DELETE CASCADE,
                               FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
                               INDEX idx_room_time (room_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天消息表';
