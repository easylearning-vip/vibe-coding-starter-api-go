-- Initial data for vibe-coding-starter (MySQL)
-- Created: 2025-01-15

-- Insert default users
-- Password: vibecoding 
INSERT INTO users (username, email, password, nickname, role, status, created_at, updated_at) VALUES
('admin', 'admin@vibe-coding.com', '$2a$10$tSfU6C4n.M7bawa9M2tU9utad/vbkr8ncZudjT6HQnR41.Ty8qkGK', 'Admin User', 'admin', 'active', NOW(), NOW()),
('demo', 'demo@vibe-coding.com', '$2a$10$tSfU6C4n.M7bawa9M2tU9utad/vbkr8ncZudjT6HQnR41.Ty8qkGK', 'Demo User', 'user', 'active', NOW(), NOW());

-- Insert main categories
INSERT INTO categories (name, slug, description, sort_order, is_active, created_at, updated_at) VALUES
('Technology', 'technology', 'Articles about technology and programming', 1, TRUE, NOW(), NOW()),
('Web Development', 'web-development', 'Web development tutorials and tips', 2, TRUE, NOW(), NOW()),
('Mobile Development', 'mobile-development', 'Mobile app development guides', 3, TRUE, NOW(), NOW()),
('DevOps', 'devops', 'DevOps practices and tools', 4, TRUE, NOW(), NOW()),
('Database', 'database', 'Database design and optimization', 5, TRUE, NOW(), NOW());

-- Insert subcategories
INSERT INTO categories (name, slug, description, parent_id, sort_order, is_active, created_at, updated_at) VALUES
('Frontend', 'frontend', 'Frontend development technologies', 2, 1, TRUE, NOW(), NOW()),
('Backend', 'backend', 'Backend development frameworks', 2, 2, TRUE, NOW(), NOW()),
('iOS', 'ios', 'iOS app development', 3, 1, TRUE, NOW(), NOW()),
('Android', 'android', 'Android app development', 3, 2, TRUE, NOW(), NOW());

-- Insert default tags
INSERT INTO tags (name, slug, description, color, created_at, updated_at) VALUES
('Go', 'go', 'Go programming language', '#00ADD8', NOW(), NOW()),
('Gin', 'gin', 'Gin web framework', '#007bff', NOW(), NOW()),
('MySQL', 'mysql', 'MySQL database', '#4479A1', NOW(), NOW()),
('Redis', 'redis', 'Redis cache', '#DC382D', NOW(), NOW()),
('Docker', 'docker', 'Docker containerization', '#2496ED', NOW(), NOW()),
('JavaScript', 'javascript', 'JavaScript programming', '#F7DF1E', NOW(), NOW()),
('React', 'react', 'React framework', '#61DAFB', NOW(), NOW()),
('Vue', 'vue', 'Vue.js framework', '#4FC08D', NOW(), NOW()),
('Node.js', 'nodejs', 'Node.js runtime', '#339933', NOW(), NOW()),
('Python', 'python', 'Python programming language', '#3776AB', NOW(), NOW());

-- Insert sample articles
INSERT INTO articles (title, slug, content, excerpt, status, published_at, author_id, category_id, created_at, updated_at) VALUES
('Welcome to Vibe Coding Starter', 'welcome-to-vibe-coding-starter', 
'# Welcome to Vibe Coding Starter

This is a comprehensive Go web application starter template built with the Gin framework. It provides a solid foundation for building modern web applications with best practices.

## Features

- **RESTful API**: Clean and well-structured REST API endpoints
- **Authentication**: JWT-based authentication system
- **Database**: MySQL with GORM ORM
- **Caching**: Redis integration for performance
- **Testing**: Comprehensive test suite with Docker support
- **Code Generation**: Built-in code generators for rapid development
- **Documentation**: Auto-generated API documentation

## Getting Started

1. Clone the repository
2. Configure your database settings
3. Run migrations
4. Start the server

Happy coding!', 
'A comprehensive Go web application starter template built with the Gin framework.',
'published', NOW(), 1, 1, NOW(), NOW()),

('Building RESTful APIs with Go and Gin', 'building-restful-apis-with-go-and-gin',
'# Building RESTful APIs with Go and Gin

Learn how to build robust RESTful APIs using Go and the Gin web framework.

## What is Gin?

Gin is a HTTP web framework written in Go. It features a Martini-like API with much better performance.

## Key Features

- Zero allocation router
- Middleware support
- JSON validation
- Route grouping
- Error management

## Example

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run()
}
```

This is just the beginning of your journey with Go and Gin!',
'Learn how to build robust RESTful APIs using Go and the Gin web framework.',
'published', NOW(), 1, 2, NOW(), NOW());

-- Link articles with tags
INSERT INTO article_tags (article_id, tag_id, created_at) VALUES
(1, 1, NOW()), -- Welcome article -> Go
(1, 2, NOW()), -- Welcome article -> Gin
(1, 3, NOW()), -- Welcome article -> MySQL
(1, 4, NOW()), -- Welcome article -> Redis
(1, 5, NOW()), -- Welcome article -> Docker
(2, 1, NOW()), -- API article -> Go
(2, 2, NOW()); -- API article -> Gin

-- Insert sample comments
INSERT INTO comments (content, author_id, article_id, status, created_at, updated_at) VALUES
('Great starter template! This will save me a lot of time.', 2, 1, 'approved', NOW(), NOW()),
('Very comprehensive guide. Thanks for sharing!', 2, 2, 'approved', NOW(), NOW());
