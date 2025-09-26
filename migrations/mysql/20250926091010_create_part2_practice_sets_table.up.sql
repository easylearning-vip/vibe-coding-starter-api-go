CREATE TABLE IF NOT EXISTS part2_practice_sets (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  scenario_id INT NULL,
  difficulty_level_id INT NULL,
  total_questions INT NOT NULL DEFAULT 0,
  completed_count INT NOT NULL DEFAULT 0,
  correct_count INT NOT NULL DEFAULT 0,
  accuracy DOUBLE NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  INDEX idx_p2ps_user (user_id),
  INDEX idx_p2ps_scenario (scenario_id),
  INDEX idx_p2ps_diff (difficulty_level_id),
  CONSTRAINT fk_p2ps_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

