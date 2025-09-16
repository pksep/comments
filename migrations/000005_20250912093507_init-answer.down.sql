-- Drop indexes first
DROP INDEX IF EXISTS idx_answers_exam_id;
DROP INDEX IF EXISTS idx_answers_slide_id;

-- Then drop the table
DROP TABLE IF EXISTS answers;