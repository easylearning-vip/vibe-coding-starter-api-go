CREATE TABLE IF NOT EXISTS part2_practice_set_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  set_id BIGINT UNSIGNED NOT NULL,
  question_id INT NOT NULL,
  order_index INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  INDEX idx_p2psi_set (set_id),
  INDEX idx_p2psi_question (question_id),
  CONSTRAINT fk_p2psi_set FOREIGN KEY (set_id) REFERENCES part2_practice_sets(id) ON DELETE CASCADE,
  CONSTRAINT fk_p2psi_q FOREIGN KEY (question_id) REFERENCES part2_questions(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

