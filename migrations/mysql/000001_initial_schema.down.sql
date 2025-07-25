-- Rollback initial database schema for vibe-coding-starter (MySQL)

-- Drop tables in reverse order to avoid foreign key constraints
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS article_tags;
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
