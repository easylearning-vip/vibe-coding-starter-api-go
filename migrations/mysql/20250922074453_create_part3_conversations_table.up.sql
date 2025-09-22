-- Migration: create_part3_conversations_table
-- Created: 20250922074453
-- Description: Create part3_conversations table


CREATE TABLE IF NOT EXISTS part3_conversations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    test_id INT NOT NULL,
    conversation_number INT NOT NULL,
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
    UNIQUE KEY uk_part3_conversations_name (name),
    KEY idx_part3_conversations_created_at (created_at),
    KEY idx_part3_conversations_updated_at (updated_at),
    KEY idx_part3_conversations_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

