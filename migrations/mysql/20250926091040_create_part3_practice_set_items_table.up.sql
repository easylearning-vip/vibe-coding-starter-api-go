CREATE TABLE IF NOT EXISTS part3_practice_set_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  set_id BIGINT UNSIGNED NOT NULL,
  conversation_id INT NOT NULL,
  question_index INT NOT NULL,
  order_index INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  PRIMARY KEY (id),
  INDEX idx_p3psi_set (set_id),
  INDEX idx_p3psi_conv (conversation_id),
  INDEX idx_p3psi_qidx (question_index),
  CONSTRAINT fk_p3psi_set FOREIGN KEY (set_id) REFERENCES part3_practice_sets(id) ON DELETE CASCADE,
  CONSTRAINT fk_p3psi_conv FOREIGN KEY (conversation_id) REFERENCES part3_conversations(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

