-- Add answer fields to practice set item tables

-- Part2
ALTER TABLE part2_practice_set_items
  ADD COLUMN selected_answer CHAR(1) NULL,
  ADD COLUMN is_correct BOOLEAN NOT NULL DEFAULT FALSE,
  ADD INDEX idx_part2_practice_set_items_is_correct (is_correct);

-- Part3
ALTER TABLE part3_practice_set_items
  ADD COLUMN selected_answer CHAR(1) NULL,
  ADD COLUMN is_correct BOOLEAN NOT NULL DEFAULT FALSE,
  ADD INDEX idx_part3_practice_set_items_is_correct (is_correct);

-- Part4
ALTER TABLE part4_practice_set_items
  ADD COLUMN selected_answer CHAR(1) NULL,
  ADD COLUMN is_correct BOOLEAN NOT NULL DEFAULT FALSE,
  ADD INDEX idx_part4_practice_set_items_is_correct (is_correct);

