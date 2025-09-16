-- Drop indexes first
DROP INDEX IF EXISTS idx_exams_project_id;
DROP INDEX IF EXISTS idx_exams_user_id;

-- Then drop the table
DROP TABLE IF EXISTS exams;