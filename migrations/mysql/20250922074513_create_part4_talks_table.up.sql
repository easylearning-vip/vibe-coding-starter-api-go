-- Migration: create_part4_talks_table
-- Created: 20250922074513
-- Description: Create part4_talks table


CREATE TABLE IF NOT EXISTS part4_talks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    test_id INT NOT NULL,
    talk_number INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content VARCHAR(255) NOT NULL,
    question1 VARCHAR(255) NOT NULL,
    question2 VARCHAR(255) NOT NULL,
    question3 VARCHAR(255) NOT NULL,
    scenario_id INT NOT NULL,
    difficulty_level_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    -- Indexes
    UNIQUE KEY uk_part4_talks_name (name),
    KEY idx_part4_talks_created_at (created_at),
    KEY idx_part4_talks_updated_at (updated_at),
    KEY idx_part4_talks_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

