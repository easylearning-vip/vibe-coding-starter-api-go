-- Migration: create_product_categorys_table
-- Created: 20250815072718
-- Description: Create product_categories table


CREATE TABLE IF NOT EXISTS product_categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    parent_id INT UNSIGNED NOT NULL,
    sort_order INT NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    -- Indexes
    UNIQUE KEY uk_product_categories_name (name),
    KEY idx_product_categories_created_at (created_at),
    KEY idx_product_categories_updated_at (updated_at),
    KEY idx_product_categories_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

