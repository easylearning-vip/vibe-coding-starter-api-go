-- Initial data for vibe-coding-starter (PostgreSQL)
-- Created: 2025-01-15

-- Insert default users
-- Password: vibecoding
INSERT INTO users (username, email, password, first_name, last_name, is_active, is_admin, email_verified_at, created_at, updated_at) VALUES
('admin', 'admin@vibe-coding.com', '$2a$10$tSfU6C4n.M7bawa9M2tU9utad/vbkr8ncZudjT6HQnR41.Ty8qkGK', 'Admin', 'User', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('demo', 'demo@vibe-coding.com', '$2a$10$tSfU6C4n.M7bawa9M2tU9utad/vbkr8ncZudjT6HQnR41.Ty8qkGK', 'Demo', 'User', TRUE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert main categories
INSERT INTO categories (name, slug, description, sort_order, is_active, created_at, updated_at) VALUES
('Technology', 'technology', 'Articles about technology and programming', 1, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Web Development', 'web-development', 'Web development tutorials and tips', 2, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Mobile Development', 'mobile-development', 'Mobile app development guides', 3, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('DevOps', 'devops', 'DevOps practices and tools', 4, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Database', 'database', 'Database design and optimization', 5, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert subcategories
INSERT INTO categories (name, slug, description, parent_id, sort_order, is_active, created_at, updated_at) VALUES
('Frontend', 'frontend', 'Frontend development technologies', 2, 1, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Backend', 'backend', 'Backend development frameworks', 2, 2, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('iOS', 'ios', 'iOS app development', 3, 1, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Android', 'android', 'Android app development', 3, 2, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert default tags
INSERT INTO tags (name, slug, description, color, created_at, updated_at) VALUES
('Go', 'go', 'Go programming language', '#00ADD8', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Gin', 'gin', 'Gin web framework', '#007bff', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('PostgreSQL', 'postgresql', 'PostgreSQL database', '#336791', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Redis', 'redis', 'Redis cache', '#DC382D', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Docker', 'docker', 'Docker containerization', '#2496ED', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('JavaScript', 'javascript', 'JavaScript programming', '#F7DF1E', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('React', 'react', 'React framework', '#61DAFB', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Vue', 'vue', 'Vue.js framework', '#4FC08D', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Node.js', 'nodejs', 'Node.js runtime', '#339933', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Python', 'python', 'Python programming language', '#3776AB', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample articles
INSERT INTO articles (title, slug, content, excerpt, status, published_at, author_id, category_id, created_at, updated_at) VALUES
('Welcome to Vibe Coding Starter', 'welcome-to-vibe-coding-starter', 
'# Welcome to Vibe Coding Starter

This is a comprehensive Go web application starter template built with the Gin framework. It provides a solid foundation for building modern web applications with best practices.

## Features

- **RESTful API**: Clean and well-structured REST API endpoints
- **Authentication**: JWT-based authentication system
- **Database**: PostgreSQL with GORM ORM
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
'published', CURRENT_TIMESTAMP, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

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
'published', CURRENT_TIMESTAMP, 1, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Link articles with tags
INSERT INTO article_tags (article_id, tag_id, created_at) VALUES
(1, 1, CURRENT_TIMESTAMP), -- Welcome article -> Go
(1, 2, CURRENT_TIMESTAMP), -- Welcome article -> Gin
(1, 3, CURRENT_TIMESTAMP), -- Welcome article -> PostgreSQL
(1, 4, CURRENT_TIMESTAMP), -- Welcome article -> Redis
(1, 5, CURRENT_TIMESTAMP), -- Welcome article -> Docker
(2, 1, CURRENT_TIMESTAMP), -- API article -> Go
(2, 2, CURRENT_TIMESTAMP); -- API article -> Gin

-- Insert sample comments
INSERT INTO comments (content, author_id, article_id, status, created_at, updated_at) VALUES
('Great starter template! This will save me a lot of time.', 2, 1, 'approved', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Very comprehensive guide. Thanks for sharing!', 2, 2, 'approved', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
