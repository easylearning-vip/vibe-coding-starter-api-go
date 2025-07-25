-- Remove initial dictionary data for vibe-coding-starter (PostgreSQL)
-- Created: 2025-01-23

-- Delete dictionary items (in reverse order)
DELETE FROM dict_items WHERE category_code IN ('article_status', 'comment_status', 'user_role', 'user_status', 'storage_type');

-- Delete dictionary categories
DELETE FROM dict_categories WHERE code IN ('article_status', 'comment_status', 'user_role', 'user_status', 'storage_type');
