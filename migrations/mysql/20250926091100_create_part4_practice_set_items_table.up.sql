CREATE TABLE IF NOT EXISTS part4_practice_set_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  set_id BIGINT UNSIGNED NOT NULL,
  talk_id INT NOT NULL,
  question_index INT NOT NULL,
  order_index INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  INDEX idx_p4psi_set (set_id),
  INDEX idx_p4psi_talk (talk_id),
  INDEX idx_p4psi_qidx (question_index),
  CONSTRAINT fk_p4psi_set FOREIGN KEY (set_id) REFERENCES part4_practice_sets(id) ON DELETE CASCADE,
  CONSTRAINT fk_p4psi_talk FOREIGN KEY (talk_id) REFERENCES part4_talks(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

