-- Migration: create_part4_practice_items_table
-- Created: 20250926090030
-- Description: Create part4_practice_items table for user practice records (talk questions)

CREATE TABLE IF NOT EXISTS part4_practice_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    talk_id INT NOT NULL,
    question_index INT NOT NULL,
    scenario_id INT NULL,
    difficulty_level_id INT NULL,
    selected_option_id INT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    KEY idx_part4_practice_items_user_id (user_id),
    KEY idx_part4_practice_items_talk_id (talk_id),
    KEY idx_part4_practice_items_question_index (question_index),
    KEY idx_part4_practice_items_scenario_id (scenario_id),
    KEY idx_part4_practice_items_difficulty_level_id (difficulty_level_id),
    KEY idx_part4_practice_items_is_correct (is_correct),
    KEY idx_part4_practice_items_created_at (created_at),
    KEY idx_part4_practice_items_deleted_at (deleted_at),

    CONSTRAINT fk_part4_practice_items_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_part4_practice_items_talk_id FOREIGN KEY (talk_id) REFERENCES part4_talks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

