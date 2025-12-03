-- Add indexes to improve query performance

-- Index for quizzes table
CREATE INDEX IF NOT EXISTS idx_quizzes_user_id_created_at ON quizzes(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_quizzes_user_id ON quizzes(user_id);

-- Index for questions table
CREATE INDEX IF NOT EXISTS idx_questions_quiz_id ON questions(quiz_id);

-- Index for quiz_attempts table (for history queries)
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_user_id_created_at ON quiz_attempts(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_quiz_id ON quiz_attempts(quiz_id);
