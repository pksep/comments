CREATE TABLE IF NOT EXISTS answers (
  id UUID PRIMARY KEY,
  exam_id UUID NOT NULL,
  slide_id UUID NOT NULL,
  is_correct BOOLEAN NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  
  -- Foreign key constraints
  CONSTRAINT fk_exam FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
  CONSTRAINT fk_slide FOREIGN KEY (slide_id) REFERENCES slides(id) ON DELETE CASCADE,
  
  -- Ensure one answer per exam-slide combination
  CONSTRAINT unique_exam_slide UNIQUE (exam_id, slide_id)
);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_answers_exam_id ON answers(exam_id);
CREATE INDEX IF NOT EXISTS idx_answers_slide_id ON answers(slide_id);