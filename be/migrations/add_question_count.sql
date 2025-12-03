-- Migration: Add question_count column to quizzes table for faster reads
-- This denormalizes the question count to avoid JOIN+COUNT on every list query

-- Add the column with default 0
ALTER TABLE quizzes ADD COLUMN IF NOT EXISTS question_count INTEGER DEFAULT 0;

-- Update existing quizzes with their actual question counts
UPDATE quizzes q
SET question_count = (
    SELECT COUNT(*) 
    FROM questions qu 
    WHERE qu.quiz_id = q.id
);

-- Create index on question_count for potential filtering/sorting
CREATE INDEX IF NOT EXISTS idx_quizzes_question_count ON quizzes(question_count);

-- Optional: Create a trigger to automatically update question_count when questions are added/deleted
-- This keeps the denormalized count in sync

CREATE OR REPLACE FUNCTION update_quiz_question_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE quizzes SET question_count = question_count + 1 WHERE id = NEW.quiz_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE quizzes SET question_count = question_count - 1 WHERE id = OLD.quiz_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Drop trigger if exists and recreate
DROP TRIGGER IF EXISTS trigger_update_question_count ON questions;

CREATE TRIGGER trigger_update_question_count
AFTER INSERT OR DELETE ON questions
FOR EACH ROW
EXECUTE FUNCTION update_quiz_question_count();
