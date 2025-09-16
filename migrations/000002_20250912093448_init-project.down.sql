-- Drop the index first to avoid dependency issues
DROP INDEX IF EXISTS idx_projects_name;

-- Then drop the table
DROP TABLE IF EXISTS projects;