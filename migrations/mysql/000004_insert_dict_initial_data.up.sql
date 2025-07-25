-- Insert initial dictionary data for vibe-coding-starter (MySQL)
-- Created: 2025-01-23

-- Insert dictionary categories
INSERT INTO dict_categories (code, name, description, sort_order) VALUES
('article_status', '文章状态', '文章的发布状态管理', 1),
('comment_status', '评论状态', '评论的审核状态管理', 2),
('user_role', '用户角色', '用户权限角色管理', 3),
('user_status', '用户状态', '用户账户状态管理', 4),
('storage_type', '存储类型', '文件存储类型管理', 5);

-- Insert dictionary items for article_status
INSERT INTO dict_items (category_code, item_key, item_value, description, sort_order, is_active) VALUES
('article_status', 'draft', '草稿', '文章草稿状态，未发布', 1, TRUE),
('article_status', 'published', '已发布', '文章已发布状态，对外可见', 2, TRUE),
('article_status', 'archived', '已归档', '文章已归档状态，不再显示', 3, TRUE);

-- Insert dictionary items for comment_status
INSERT INTO dict_items (category_code, item_key, item_value, description, sort_order, is_active) VALUES
('comment_status', 'pending', '待审核', '评论待审核状态，暂不显示', 1, TRUE),
('comment_status', 'approved', '已批准', '评论已批准状态，对外可见', 2, TRUE),
('comment_status', 'rejected', '已拒绝', '评论已拒绝状态，不予显示', 3, TRUE);

-- Insert dictionary items for user_role
INSERT INTO dict_items (category_code, item_key, item_value, description, sort_order, is_active) VALUES
('user_role', 'admin', '管理员', '系统管理员，拥有所有权限', 1, TRUE),
('user_role', 'user', '普通用户', '普通用户，拥有基本权限', 2, TRUE);

-- Insert dictionary items for user_status
INSERT INTO dict_items (category_code, item_key, item_value, description, sort_order, is_active) VALUES
('user_status', 'active', '活跃', '用户账户正常活跃状态', 1, TRUE),
('user_status', 'inactive', '非活跃', '用户账户非活跃状态', 2, TRUE),
('user_status', 'banned', '已禁用', '用户账户已被禁用', 3, TRUE);

-- Insert dictionary items for storage_type
INSERT INTO dict_items (category_code, item_key, item_value, description, sort_order, is_active) VALUES
('storage_type', 'local', '本地存储', '文件存储在本地服务器', 1, TRUE),
('storage_type', 's3', 'AWS S3', '文件存储在Amazon S3', 2, TRUE),
('storage_type', 'oss', '阿里云OSS', '文件存储在阿里云对象存储', 3, TRUE);
