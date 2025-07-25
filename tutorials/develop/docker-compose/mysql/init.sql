-- MySQL 初始化脚本
-- 用于开发环境的数据库初始化

-- 创建开发数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `vibe_coding_starter` 
    CHARACTER SET utf8mb4 
    COLLATE utf8mb4_unicode_ci;

-- 数据库初始化完成
SELECT 'Database schema initialization completed' as status;
