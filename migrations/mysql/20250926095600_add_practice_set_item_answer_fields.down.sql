-- Revert answer fields from practice set item tables

-- Part2
ALTER TABLE part2_practice_set_items
  DROP INDEX idx_part2_practice_set_items_is_correct,
  DROP COLUMN is_correct,
  DROP COLUMN selected_answer;

-- Part3
ALTER TABLE part3_practice_set_items
  DROP INDEX idx_part3_practice_set_items_is_correct,
  DROP COLUMN is_correct,
  DROP COLUMN selected_answer;

-- Part4
ALTER TABLE part4_practice_set_items
  DROP INDEX idx_part4_practice_set_items_is_correct,
  DROP COLUMN is_correct,
  DROP COLUMN selected_answer;

