-- 确保使用 im 数据库
CREATE DATABASE IF NOT EXISTS im;
USE im;

-- 如果表不存在，可以先根据模型创建（可选）
CREATE TABLE IF NOT EXISTS `user` (
    `id`         BIGINT AUTO_INCREMENT PRIMARY KEY,
    `uid`        BIGINT UNIQUE NOT NULL,
    `name`       VARCHAR(64) NOT NULL,
    `password`   VARCHAR(255) NOT NULL,
    `email`      VARCHAR(128) UNIQUE NOT NULL,
    `phone`      VARCHAR(20) UNIQUE,
    `status`     TINYINT DEFAULT 1,
    `created_at` DATETIME(3) NULL,
    `updated_at` DATETIME(3) NULL,
    `deleted_at` DATETIME(3) NULL,
    INDEX `idx_user_name` (`name`),
    INDEX `idx_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 批量插入 Mock 数据
-- 注意：Password 通常存放 Bcrypt 加密后的哈希值，这里用简单的占位符
INSERT INTO `user` (`uid`, `name`, `password`, `email`, `phone`, `status`, `created_at`, `updated_at`)
VALUES 
(10001, '张三', '123', 'zhangsan@example.com', '13800000001', 1, NOW(), NOW()),
(10002, '李四', '123', 'lisi@example.com', '13800000002', 1, NOW(), NOW()),
(10003, '王五', '123', 'wangwu@example.com', '13800000003', 1, NOW(), NOW()),
(10004, '赵六', '123', 'zhaoliu@example.com', '13800000004', 2, NOW(), NOW()),
(10005, '小明', '123', 'xiaoming@example.com', '13800000005', 1, NOW(), NOW()),
(10006, '小红', '123', 'xiaohong@example.com', '13800000006', 1, NOW(), NOW()),
(10007, '陈七', '123', 'chenqi@example.com', '13800000007', 1, NOW(), NOW()),
(10008, '老王', '123', 'laowang@example.com', '13800000008', 1, NOW(), NOW()),
(10009, '阿强', '123', 'aqiang@example.com', '13800000009', 1, NOW(), NOW()),
(10010, '露西', '123', 'lucy@example.com', '13800000010', 1, NOW(), NOW());