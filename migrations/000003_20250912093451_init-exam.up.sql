CREATE TABLE IF NOT EXISTS exams (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  user_id UUID NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  
  -- Foreign key constraints
  CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_exams_project_id ON exams(project_id);
CREATE INDEX IF NOT EXISTS idx_exams_user_id ON exams(user_id);