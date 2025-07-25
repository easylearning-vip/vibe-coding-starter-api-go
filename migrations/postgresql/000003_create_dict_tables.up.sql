-- Create dictionary tables for vibe-coding-starter (PostgreSQL)
-- Created: 2025-01-23

-- Dictionary categories table
CREATE TABLE dict_categories (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uk_dict_categories_code UNIQUE (code)
);

-- Create indexes for dict_categories table
CREATE INDEX idx_dict_categories_sort_order ON dict_categories(sort_order);
CREATE INDEX idx_dict_categories_created_at ON dict_categories(created_at);
CREATE INDEX idx_dict_categories_deleted_at ON dict_categories(deleted_at);

-- Dictionary items table
CREATE TABLE dict_items (
    id BIGSERIAL PRIMARY KEY,
    category_code VARCHAR(50) NOT NULL,
    item_key VARCHAR(50) NOT NULL,
    item_value VARCHAR(200) NOT NULL,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uk_dict_items_category_key UNIQUE (category_code, item_key),
    CONSTRAINT fk_dict_items_category_code FOREIGN KEY (category_code) REFERENCES dict_categories(code) ON DELETE CASCADE
);

-- Create indexes for dict_items table
CREATE INDEX idx_dict_items_category_code ON dict_items(category_code);
CREATE INDEX idx_dict_items_is_active ON dict_items(is_active);
CREATE INDEX idx_dict_items_sort_order ON dict_items(sort_order);
CREATE INDEX idx_dict_items_created_at ON dict_items(created_at);
CREATE INDEX idx_dict_items_deleted_at ON dict_items(deleted_at);
