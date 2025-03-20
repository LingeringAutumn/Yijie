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