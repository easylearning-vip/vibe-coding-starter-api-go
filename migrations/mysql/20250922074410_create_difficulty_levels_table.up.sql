-- Migration: create_difficulty_levels_table
-- Created: 20250922074410
-- Description: Create difficulty_levels table


CREATE TABLE IF NOT EXISTS difficulty_levels (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    characteristics VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    -- Indexes
    UNIQUE KEY uk_difficulty_levels_name (name),
    KEY idx_difficulty_levels_created_at (created_at),
    KEY idx_difficulty_levels_updated_at (updated_at),
    KEY idx_difficulty_levels_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

