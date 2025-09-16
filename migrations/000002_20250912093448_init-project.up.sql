CREATE TABLE IF NOT EXISTS projects (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  exam_duration INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- Add index for better query performance on project names
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
