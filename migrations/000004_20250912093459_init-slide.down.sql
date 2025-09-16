-- Drop indexes first
DROP INDEX IF EXISTS idx_slides_project_id;
DROP INDEX IF EXISTS idx_slides_path;

-- Then drop the table
DROP TABLE IF EXISTS slides;