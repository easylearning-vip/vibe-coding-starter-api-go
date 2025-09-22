-- Migration: create_part3_answer_options_table
-- Created: 20250922074504
-- Description: Create part3_answer_options table


CREATE TABLE IF NOT EXISTS part3_answer_options (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    conversation_id INT NOT NULL,
    question_number INT NOT NULL,
    option_a VARCHAR(255) NOT NULL,
    option_b VARCHAR(255) NOT NULL,
    option_c VARCHAR(255) NOT NULL,
    option_d VARCHAR(255) NOT NULL,
    correct_answer VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    -- Indexes
    UNIQUE KEY uk_part3_answer_options_name (name),
    KEY idx_part3_answer_options_created_at (created_at),
    KEY idx_part3_answer_options_updated_at (updated_at),
    KEY idx_part3_answer_options_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

