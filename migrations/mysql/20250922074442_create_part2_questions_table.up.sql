-- Migration: create_part2_questions_table
-- Created: 20250922074442
-- Description: Create part2_questions table


CREATE TABLE IF NOT EXISTS part2_questions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    test_id INT NOT NULL,
    question_number INT NOT NULL,
    question_text VARCHAR(255) NOT NULL,
    option_a VARCHAR(255) NOT NULL,
    option_b VARCHAR(255) NOT NULL,
    option_c VARCHAR(255) NOT NULL,
    correct_answer VARCHAR(255) NOT NULL,
    scenario_id INT NOT NULL,
    difficulty_level_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    -- Indexes
    UNIQUE KEY uk_part2_questions_name (name),
    KEY idx_part2_questions_created_at (created_at),
    KEY idx_part2_questions_updated_at (updated_at),
    KEY idx_part2_questions_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

