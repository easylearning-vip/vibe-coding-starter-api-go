-- Create dictionary tables for vibe-coding-starter (MySQL)
-- Created: 2025-01-23

-- Dictionary categories table
CREATE TABLE dict_categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    UNIQUE KEY uk_dict_categories_code (code),
    KEY idx_dict_categories_sort_order (sort_order),
    KEY idx_dict_categories_created_at (created_at),
    KEY idx_dict_categories_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dictionary items table
CREATE TABLE dict_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    category_code VARCHAR(50) NOT NULL,
    item_key VARCHAR(50) NOT NULL,
    item_value VARCHAR(200) NOT NULL,
    description TEXT,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    UNIQUE KEY uk_dict_items_category_key (category_code, item_key),
    KEY idx_dict_items_category_code (category_code),
    KEY idx_dict_items_is_active (is_active),
    KEY idx_dict_items_sort_order (sort_order),
    KEY idx_dict_items_created_at (created_at),
    KEY idx_dict_items_deleted_at (deleted_at),

    CONSTRAINT fk_dict_items_category_code FOREIGN KEY (category_code) REFERENCES dict_categories(code) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
