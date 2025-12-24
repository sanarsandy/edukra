package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// QuizRepository implements domain.QuizRepository using PostgreSQL
type QuizRepository struct {
	db *sqlx.DB
}

// NewQuizRepository creates a new QuizRepository
func NewQuizRepository(db *sqlx.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

// GetByID retrieves a quiz by ID
func (r *QuizRepository) GetByID(id string) (*domain.Quiz, error) {
	query := `
		SELECT id, lesson_id, title, description, time_limit, passing_score, 
		       max_attempts, shuffle_questions, shuffle_options, show_correct_answers,
		       created_at, updated_at
		FROM quizzes WHERE id = $1
	`
	
	var quiz domain.Quiz
	var description sql.NullString
	
	err := r.db.QueryRow(query, id).Scan(
		&quiz.ID, &quiz.LessonID, &quiz.Title, &description, &quiz.TimeLimit,
		&quiz.PassingScore, &quiz.MaxAttempts, &quiz.ShuffleQuestions,
		&quiz.ShuffleOptions, &quiz.ShowCorrectAnswers, &quiz.CreatedAt, &quiz.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if description.Valid {
		quiz.Description = description.String
	}
	
	return &quiz, nil
}

// GetByLessonID retrieves a quiz by lesson ID
func (r *QuizRepository) GetByLessonID(lessonID string) (*domain.Quiz, error) {
	query := `
		SELECT id, lesson_id, title, description, time_limit, passing_score, 
		       max_attempts, shuffle_questions, shuffle_options, show_correct_answers,
		       created_at, updated_at
		FROM quizzes WHERE lesson_id = $1
	`
	
	var quiz domain.Quiz
	var description sql.NullString
	
	err := r.db.QueryRow(query, lessonID).Scan(
		&quiz.ID, &quiz.LessonID, &quiz.Title, &description, &quiz.TimeLimit,
		&quiz.PassingScore, &quiz.MaxAttempts, &quiz.ShuffleQuestions,
		&quiz.ShuffleOptions, &quiz.ShowCorrectAnswers, &quiz.CreatedAt, &quiz.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if description.Valid {
		quiz.Description = description.String
	}
	
	return &quiz, nil
}

// Create inserts a new quiz
func (r *QuizRepository) Create(quiz *domain.Quiz) error {
	query := `
		INSERT INTO quizzes (lesson_id, title, description, time_limit, passing_score,
		                     max_attempts, shuffle_questions, shuffle_options, 
		                     show_correct_answers, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	
	now := time.Now()
	quiz.CreatedAt = now
	quiz.UpdatedAt = now
	
	var desc *string
	if quiz.Description != "" {
		desc = &quiz.Description
	}
	
	return r.db.QueryRow(query,
		quiz.LessonID, quiz.Title, desc, quiz.TimeLimit, quiz.PassingScore,
		quiz.MaxAttempts, quiz.ShuffleQuestions, quiz.ShuffleOptions,
		quiz.ShowCorrectAnswers, quiz.CreatedAt, quiz.UpdatedAt,
	).Scan(&quiz.ID)
}

// Update updates an existing quiz
func (r *QuizRepository) Update(quiz *domain.Quiz) error {
	query := `
		UPDATE quizzes SET
			title = $2, description = $3, time_limit = $4, passing_score = $5,
			max_attempts = $6, shuffle_questions = $7, shuffle_options = $8,
			show_correct_answers = $9, updated_at = $10
		WHERE id = $1
	`
	
	quiz.UpdatedAt = time.Now()
	
	var desc *string
	if quiz.Description != "" {
		desc = &quiz.Description
	}
	
	_, err := r.db.Exec(query,
		quiz.ID, quiz.Title, desc, quiz.TimeLimit, quiz.PassingScore,
		quiz.MaxAttempts, quiz.ShuffleQuestions, quiz.ShuffleOptions,
		quiz.ShowCorrectAnswers, quiz.UpdatedAt,
	)
	return err
}

// Delete removes a quiz by ID
func (r *QuizRepository) Delete(id string) error {
	query := `DELETE FROM quizzes WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetQuestionsByQuizID retrieves all questions for a quiz
func (r *QuizRepository) GetQuestionsByQuizID(quizID string) ([]domain.Question, error) {
	query := `
		SELECT id, quiz_id, question_type, question_text, explanation, points, 
		       order_index, required, created_at
		FROM quiz_questions 
		WHERE quiz_id = $1
		ORDER BY order_index ASC
	`
	
	rows, err := r.db.Query(query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var questions []domain.Question
	for rows.Next() {
		var q domain.Question
		var explanation sql.NullString
		
		err := rows.Scan(
			&q.ID, &q.QuizID, &q.Type, &q.QuestionText, &explanation,
			&q.Points, &q.OrderIndex, &q.Required, &q.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if explanation.Valid {
			q.Explanation = explanation.String
		}
		
		// Load options for this question
		q.Options, err = r.GetOptionsByQuestionID(q.ID)
		if err != nil {
			return nil, err
		}
		
		questions = append(questions, q)
	}
	
	return questions, nil
}

// GetQuestionByID retrieves a question by ID
func (r *QuizRepository) GetQuestionByID(id string) (*domain.Question, error) {
	query := `
		SELECT id, quiz_id, question_type, question_text, explanation, points, 
		       order_index, required, created_at
		FROM quiz_questions WHERE id = $1
	`
	
	var q domain.Question
	var explanation sql.NullString
	
	err := r.db.QueryRow(query, id).Scan(
		&q.ID, &q.QuizID, &q.Type, &q.QuestionText, &explanation,
		&q.Points, &q.OrderIndex, &q.Required, &q.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if explanation.Valid {
		q.Explanation = explanation.String
	}
	
	// Load options
	q.Options, err = r.GetOptionsByQuestionID(q.ID)
	if err != nil {
		return nil, err
	}
	
	return &q, nil
}

// CreateQuestion inserts a new question
func (r *QuizRepository) CreateQuestion(question *domain.Question) error {
	// Get current question count for order_index
	var count int
	r.db.QueryRow(`SELECT COUNT(*) FROM quiz_questions WHERE quiz_id = $1`, question.QuizID).Scan(&count)
	question.OrderIndex = count
	
	query := `
		INSERT INTO quiz_questions (quiz_id, question_type, question_text, explanation,
		                            points, order_index, required, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	
	question.CreatedAt = time.Now()
	
	var explanation *string
	if question.Explanation != "" {
		explanation = &question.Explanation
	}
	
	return r.db.QueryRow(query,
		question.QuizID, question.Type, question.QuestionText, explanation,
		question.Points, question.OrderIndex, question.Required, question.CreatedAt,
	).Scan(&question.ID)
}

// UpdateQuestion updates an existing question
func (r *QuizRepository) UpdateQuestion(question *domain.Question) error {
	query := `
		UPDATE quiz_questions SET
			question_type = $2, question_text = $3, explanation = $4,
			points = $5, required = $6
		WHERE id = $1
	`
	
	var explanation *string
	if question.Explanation != "" {
		explanation = &question.Explanation
	}
	
	_, err := r.db.Exec(query,
		question.ID, question.Type, question.QuestionText, explanation,
		question.Points, question.Required,
	)
	return err
}

// DeleteQuestion removes a question by ID
func (r *QuizRepository) DeleteQuestion(id string) error {
	query := `DELETE FROM quiz_questions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// ReorderQuestions updates the order of questions in a quiz
func (r *QuizRepository) ReorderQuestions(quizID string, questionIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	query := `UPDATE quiz_questions SET order_index = $1 WHERE id = $2 AND quiz_id = $3`
	
	for i, qID := range questionIDs {
		_, err := tx.Exec(query, i, qID, quizID)
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

// GetOptionsByQuestionID retrieves all options for a question
func (r *QuizRepository) GetOptionsByQuestionID(questionID string) ([]domain.Option, error) {
	query := `
		SELECT id, question_id, option_text, is_correct, order_index
		FROM quiz_options 
		WHERE question_id = $1
		ORDER BY order_index ASC
	`
	
	rows, err := r.db.Query(query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var options []domain.Option
	for rows.Next() {
		var o domain.Option
		err := rows.Scan(&o.ID, &o.QuestionID, &o.OptionText, &o.IsCorrect, &o.OrderIndex)
		if err != nil {
			return nil, err
		}
		options = append(options, o)
	}
	
	return options, nil
}

// CreateOptions inserts multiple options
func (r *QuizRepository) CreateOptions(options []domain.Option) error {
	if len(options) == 0 {
		return nil
	}
	
	query := `
		INSERT INTO quiz_options (question_id, option_text, is_correct, order_index)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	
	for i := range options {
		options[i].OrderIndex = i
		err := r.db.QueryRow(query,
			options[i].QuestionID, options[i].OptionText, options[i].IsCorrect, options[i].OrderIndex,
		).Scan(&options[i].ID)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// DeleteOptionsByQuestionID removes all options for a question
func (r *QuizRepository) DeleteOptionsByQuestionID(questionID string) error {
	query := `DELETE FROM quiz_options WHERE question_id = $1`
	_, err := r.db.Exec(query, questionID)
	return err
}

// GetAttemptByID retrieves an attempt by ID
func (r *QuizRepository) GetAttemptByID(id string) (*domain.QuizAttempt, error) {
	query := `
		SELECT id, quiz_id, user_id, started_at, completed_at, score, passed, time_spent
		FROM quiz_attempts WHERE id = $1
	`
	
	var a domain.QuizAttempt
	var completedAt sql.NullTime
	var score sql.NullInt64
	var passed sql.NullBool
	var timeSpent sql.NullInt64
	
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.QuizID, &a.UserID, &a.StartedAt, &completedAt, &score, &passed, &timeSpent,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if completedAt.Valid {
		a.CompletedAt = &completedAt.Time
	}
	if score.Valid {
		s := int(score.Int64)
		a.Score = &s
	}
	if passed.Valid {
		a.Passed = &passed.Bool
	}
	if timeSpent.Valid {
		t := int(timeSpent.Int64)
		a.TimeSpent = &t
	}
	
	return &a, nil
}

// GetAttemptsByQuizAndUser retrieves all attempts by a user for a quiz
func (r *QuizRepository) GetAttemptsByQuizAndUser(quizID, userID string) ([]domain.QuizAttempt, error) {
	query := `
		SELECT id, quiz_id, user_id, started_at, completed_at, score, passed, time_spent
		FROM quiz_attempts 
		WHERE quiz_id = $1 AND user_id = $2
		ORDER BY started_at DESC
	`
	
	rows, err := r.db.Query(query, quizID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var attempts []domain.QuizAttempt
	for rows.Next() {
		var a domain.QuizAttempt
		var completedAt sql.NullTime
		var score sql.NullInt64
		var passed sql.NullBool
		var timeSpent sql.NullInt64
		
		err := rows.Scan(
			&a.ID, &a.QuizID, &a.UserID, &a.StartedAt, &completedAt, &score, &passed, &timeSpent,
		)
		if err != nil {
			return nil, err
		}
		
		if completedAt.Valid {
			a.CompletedAt = &completedAt.Time
		}
		if score.Valid {
			s := int(score.Int64)
			a.Score = &s
		}
		if passed.Valid {
			a.Passed = &passed.Bool
		}
		if timeSpent.Valid {
			t := int(timeSpent.Int64)
			a.TimeSpent = &t
		}
		
		attempts = append(attempts, a)
	}
	
	return attempts, nil
}

// CountAttemptsByQuizAndUser counts the number of attempts
func (r *QuizRepository) CountAttemptsByQuizAndUser(quizID, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM quiz_attempts WHERE quiz_id = $1 AND user_id = $2`
	var count int
	err := r.db.QueryRow(query, quizID, userID).Scan(&count)
	return count, err
}

// CreateAttempt inserts a new attempt
func (r *QuizRepository) CreateAttempt(attempt *domain.QuizAttempt) error {
	query := `
		INSERT INTO quiz_attempts (quiz_id, user_id, started_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	
	attempt.StartedAt = time.Now()
	
	return r.db.QueryRow(query,
		attempt.QuizID, attempt.UserID, attempt.StartedAt,
	).Scan(&attempt.ID)
}

// UpdateAttempt updates an existing attempt
func (r *QuizRepository) UpdateAttempt(attempt *domain.QuizAttempt) error {
	query := `
		UPDATE quiz_attempts SET
			completed_at = $2, score = $3, passed = $4, time_spent = $5
		WHERE id = $1
	`
	
	_, err := r.db.Exec(query,
		attempt.ID, attempt.CompletedAt, attempt.Score, attempt.Passed, attempt.TimeSpent,
	)
	return err
}

// CreateAnswers inserts multiple answers
func (r *QuizRepository) CreateAnswers(answers []domain.Answer) error {
	if len(answers) == 0 {
		return nil
	}
	
	query := `
		INSERT INTO quiz_answers (attempt_id, question_id, selected_option_ids, 
		                          text_answer, is_correct, points_earned)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	
	for i := range answers {
		var textAnswer *string
		if answers[i].TextAnswer != "" {
			textAnswer = &answers[i].TextAnswer
		}
		
		err := r.db.QueryRow(query,
			answers[i].AttemptID, answers[i].QuestionID, 
			pq.Array(answers[i].SelectedOptionIDs), textAnswer,
			answers[i].IsCorrect, answers[i].PointsEarned,
		).Scan(&answers[i].ID)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// GetAnswersByAttemptID retrieves all answers for an attempt
func (r *QuizRepository) GetAnswersByAttemptID(attemptID string) ([]domain.Answer, error) {
	query := `
		SELECT id, attempt_id, question_id, selected_option_ids, text_answer, 
		       is_correct, points_earned
		FROM quiz_answers WHERE attempt_id = $1
	`
	
	rows, err := r.db.Query(query, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var answers []domain.Answer
	for rows.Next() {
		var a domain.Answer
		var textAnswer sql.NullString
		var isCorrect sql.NullBool
		
		err := rows.Scan(
			&a.ID, &a.AttemptID, &a.QuestionID, 
			pq.Array(&a.SelectedOptionIDs), &textAnswer,
			&isCorrect, &a.PointsEarned,
		)
		if err != nil {
			return nil, err
		}
		
		if textAnswer.Valid {
			a.TextAnswer = textAnswer.String
		}
		if isCorrect.Valid {
			a.IsCorrect = &isCorrect.Bool
		}
		
		answers = append(answers, a)
	}
	
	return answers, nil
}
