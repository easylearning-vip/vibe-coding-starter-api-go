-- Rollback initial data for vibe-coding-starter (MySQL)

-- Delete sample data in reverse order
DELETE FROM comments WHERE id IN (1, 2);
DELETE FROM article_tags WHERE article_id IN (1, 2);
DELETE FROM articles WHERE id IN (1, 2);
DELETE FROM tags WHERE slug IN ('go', 'gin', 'mysql', 'redis', 'docker', 'javascript', 'react', 'vue', 'nodejs', 'python');
DELETE FROM categories WHERE slug IN ('frontend', 'backend', 'ios', 'android', 'technology', 'web-development', 'mobile-development', 'devops', 'database');
DELETE FROM users WHERE username IN ('admin', 'demo');
