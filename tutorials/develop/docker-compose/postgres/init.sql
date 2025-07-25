-- PostgreSQL 初始化脚本
-- 用于开发环境的数据库初始化

-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 数据库初始化完成
SELECT 'Database schema initialization completed' as status;
