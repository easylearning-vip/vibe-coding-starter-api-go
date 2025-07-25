-- Rollback initial database schema for vibe-coding-starter (PostgreSQL)

-- Drop triggers
DROP TRIGGER IF EXISTS update_files_updated_at ON files;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP TRIGGER IF EXISTS update_articles_updated_at ON articles;
DROP TRIGGER IF EXISTS update_tags_updated_at ON tags;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order to avoid foreign key constraints
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS article_tags;
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

-- Drop enum types
DROP TYPE IF EXISTS comment_status;
DROP TYPE IF EXISTS article_status;
