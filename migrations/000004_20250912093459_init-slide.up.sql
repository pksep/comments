CREATE TABLE IF NOT EXISTS slides (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  path TEXT NOT NULL,
  text TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  
  -- Foreign key constraint
  CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Add index for better query performance
CREATE INDEX IF NOT EXISTS idx_slides_project_id ON slides(project_id);

-- Add index for path lookups if needed
CREATE INDEX IF NOT EXISTS idx_slides_path ON slides(path);