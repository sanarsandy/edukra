package domain

import "time"

// QuestionType represents the type of quiz question
type QuestionType string

const (
	QuestionMultipleChoice QuestionType = "multiple_choice"  // Single correct answer
	QuestionMultipleAnswer QuestionType = "multiple_answer"  // Multiple correct answers
	QuestionTrueFalse      QuestionType = "true_false"       // True or False
	QuestionShortAnswer    QuestionType = "short_answer"     // Text input
)

// Quiz represents a quiz attached to a lesson
type Quiz struct {
	ID                 string     `json:"id"`
	LessonID           string     `json:"lesson_id"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	TimeLimit          int        `json:"time_limit"`           // in minutes, 0 = no limit
	PassingScore       int        `json:"passing_score"`        // percentage
	MaxAttempts        int        `json:"max_attempts"`         // 0 = unlimited
	ShuffleQuestions   bool       `json:"shuffle_questions"`
	ShuffleOptions     bool       `json:"shuffle_options"`
	ShowCorrectAnswers bool       `json:"show_correct_answers"`
	Questions          []Question `json:"questions,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// Question represents a quiz question
type Question struct {
	ID           string       `json:"id"`
	QuizID       string       `json:"quiz_id"`
	Type         QuestionType `json:"question_type"`
	QuestionText string       `json:"question_text"`
	Explanation  string       `json:"explanation,omitempty"`
	Points       int          `json:"points"`
	OrderIndex   int          `json:"order_index"`
	Required     bool         `json:"required"`
	Options      []Option     `json:"options,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
}

// Option represents an answer option for a question
type Option struct {
	ID         string `json:"id"`
	QuestionID string `json:"question_id"`
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct,omitempty"` // Hidden for students during quiz
	OrderIndex int    `json:"order_index"`
}

// QuizAttempt represents a user's attempt at a quiz
type QuizAttempt struct {
	ID          string     `json:"id"`
	QuizID      string     `json:"quiz_id"`
	UserID      string     `json:"user_id"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Score       *int       `json:"score,omitempty"`   // percentage
	Passed      *bool      `json:"passed,omitempty"`
	TimeSpent   *int       `json:"time_spent,omitempty"` // in seconds
	Answers     []Answer   `json:"answers,omitempty"`
}

// Answer represents a user's answer to a question
type Answer struct {
	ID                string   `json:"id"`
	AttemptID         string   `json:"attempt_id"`
	QuestionID        string   `json:"question_id"`
	SelectedOptionIDs []string `json:"selected_option_ids,omitempty"` // for multiple choice
	TextAnswer        string   `json:"text_answer,omitempty"`         // for short answer
	IsCorrect         *bool    `json:"is_correct,omitempty"`
	PointsEarned      int      `json:"points_earned"`
}

// CreateQuizRequest represents a request to create a quiz
type CreateQuizRequest struct {
	LessonID           string `json:"lesson_id"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	TimeLimit          int    `json:"time_limit"`
	PassingScore       int    `json:"passing_score"`
	MaxAttempts        int    `json:"max_attempts"`
	ShuffleQuestions   bool   `json:"shuffle_questions"`
	ShuffleOptions     bool   `json:"shuffle_options"`
	ShowCorrectAnswers bool   `json:"show_correct_answers"`
}

// UpdateQuizRequest represents a request to update quiz settings
type UpdateQuizRequest struct {
	Title              *string `json:"title,omitempty"`
	Description        *string `json:"description,omitempty"`
	TimeLimit          *int    `json:"time_limit,omitempty"`
	PassingScore       *int    `json:"passing_score,omitempty"`
	MaxAttempts        *int    `json:"max_attempts,omitempty"`
	ShuffleQuestions   *bool   `json:"shuffle_questions,omitempty"`
	ShuffleOptions     *bool   `json:"shuffle_options,omitempty"`
	ShowCorrectAnswers *bool   `json:"show_correct_answers,omitempty"`
}

// CreateQuestionRequest represents a request to create a question
type CreateQuestionRequest struct {
	QuizID       string       `json:"quiz_id"`
	Type         QuestionType `json:"question_type"`
	QuestionText string       `json:"question_text"`
	Explanation  string       `json:"explanation"`
	Points       int          `json:"points"`
	Required     bool         `json:"required"`
	Options      []CreateOptionRequest `json:"options,omitempty"`
}

// CreateOptionRequest represents a request to create an option
type CreateOptionRequest struct {
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
}

// UpdateQuestionRequest represents a request to update a question
type UpdateQuestionRequest struct {
	Type         *QuestionType `json:"question_type,omitempty"`
	QuestionText *string       `json:"question_text,omitempty"`
	Explanation  *string       `json:"explanation,omitempty"`
	Points       *int          `json:"points,omitempty"`
	Required     *bool         `json:"required,omitempty"`
	Options      []CreateOptionRequest `json:"options,omitempty"` // Replace all options
}

// SubmitQuizRequest represents a user submitting quiz answers
type SubmitQuizRequest struct {
	Answers []SubmitAnswerRequest `json:"answers"`
}

// SubmitAnswerRequest represents a single answer submission
type SubmitAnswerRequest struct {
	QuestionID        string   `json:"question_id"`
	SelectedOptionIDs []string `json:"selected_option_ids,omitempty"`
	TextAnswer        string   `json:"text_answer,omitempty"`
}

// QuizResult represents the result of a quiz attempt
type QuizResult struct {
	AttemptID      string           `json:"attempt_id"`
	QuizID         string           `json:"quiz_id"`
	Score          int              `json:"score"`
	Passed         bool             `json:"passed"`
	TimeSpent      int              `json:"time_spent"`
	TotalPoints    int              `json:"total_points"`
	EarnedPoints   int              `json:"earned_points"`
	CorrectCount   int              `json:"correct_count"`
	TotalQuestions int              `json:"total_questions"`
	Answers        []AnswerResult   `json:"answers,omitempty"`
}

// AnswerResult represents the result of a single answer
type AnswerResult struct {
	QuestionID     string   `json:"question_id"`
	QuestionText   string   `json:"question_text"`
	UserAnswer     []string `json:"user_answer"` // option IDs or text
	CorrectAnswer  []string `json:"correct_answer,omitempty"` // shown if allowed
	IsCorrect      bool     `json:"is_correct"`
	PointsEarned   int      `json:"points_earned"`
	MaxPoints      int      `json:"max_points"`
	Explanation    string   `json:"explanation,omitempty"`
}

// QuizRepository defines the interface for quiz data access
type QuizRepository interface {
	// Quiz CRUD
	GetByID(id string) (*Quiz, error)
	GetByLessonID(lessonID string) (*Quiz, error)
	Create(quiz *Quiz) error
	Update(quiz *Quiz) error
	Delete(id string) error

	// Questions
	GetQuestionsByQuizID(quizID string) ([]Question, error)
	GetQuestionByID(id string) (*Question, error)
	CreateQuestion(question *Question) error
	UpdateQuestion(question *Question) error
	DeleteQuestion(id string) error
	ReorderQuestions(quizID string, questionIDs []string) error

	// Options
	GetOptionsByQuestionID(questionID string) ([]Option, error)
	CreateOptions(options []Option) error
	DeleteOptionsByQuestionID(questionID string) error

	// Attempts
	GetAttemptByID(id string) (*QuizAttempt, error)
	GetAttemptsByQuizAndUser(quizID, userID string) ([]QuizAttempt, error)
	CountAttemptsByQuizAndUser(quizID, userID string) (int, error)
	CreateAttempt(attempt *QuizAttempt) error
	UpdateAttempt(attempt *QuizAttempt) error

	// Answers
	CreateAnswers(answers []Answer) error
	GetAnswersByAttemptID(attemptID string) ([]Answer, error)
}
