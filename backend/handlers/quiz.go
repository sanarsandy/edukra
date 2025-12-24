package handlers

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

// getQuizUserID extracts user ID from JWT token for quiz handlers
func getQuizUserID(c echo.Context) string {
	user := c.Get("user")
	if user == nil {
		return ""
	}
	token, ok := user.(*jwt.Token)
	if !ok {
		return ""
	}
	
	// Try pointer first (echojwt uses new(jwt.MapClaims) which returns pointer)
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		// Fallback: try as value type
		if mc, ok := token.Claims.(jwt.MapClaims); ok {
			claims = &mc
		} else {
			return ""
		}
	}
	
	userID, ok := (*claims)["user_id"].(string)
	if !ok {
		return ""
	}
	return userID
}

var quizRepo *postgres.QuizRepository

func initQuizRepos() {
	if quizRepo == nil && db.DB != nil {
		quizRepo = postgres.NewQuizRepository(db.DB)
	}
}

// ============ Admin Handlers ============

// CreateQuiz creates a new quiz for a lesson
func CreateQuiz(c echo.Context) error {
	initQuizRepos()
	
	lessonID := c.Param("lessonId")
	
	var req domain.CreateQuizRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	// Check if quiz already exists for this lesson
	existing, _ := quizRepo.GetByLessonID(lessonID)
	if existing != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Quiz already exists for this lesson"})
	}
	
	quiz := &domain.Quiz{
		LessonID:           lessonID,
		Title:              req.Title,
		Description:        req.Description,
		TimeLimit:          req.TimeLimit,
		PassingScore:       req.PassingScore,
		MaxAttempts:        req.MaxAttempts,
		ShuffleQuestions:   req.ShuffleQuestions,
		ShuffleOptions:     req.ShuffleOptions,
		ShowCorrectAnswers: req.ShowCorrectAnswers,
	}
	
	if quiz.PassingScore == 0 {
		quiz.PassingScore = 70
	}
	
	err := quizRepo.Create(quiz)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create quiz: " + err.Error()})
	}
	
	return c.JSON(http.StatusCreated, quiz)
}

// GetQuiz retrieves a quiz with all questions (admin view with correct answers)
func GetQuiz(c echo.Context) error {
	initQuizRepos()
	
	id := c.Param("id")
	if id == "" {
		id = c.Param("lessonId")
	}
	
	var quiz *domain.Quiz
	var err error
	
	// Try to get by quiz ID first, then by lesson ID
	quiz, err = quizRepo.GetByID(id)
	if quiz == nil {
		quiz, err = quizRepo.GetByLessonID(id)
	}
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch quiz"})
	}
	if quiz == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Quiz not found"})
	}
	
	// Load questions
	quiz.Questions, err = quizRepo.GetQuestionsByQuizID(quiz.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch questions"})
	}
	
	return c.JSON(http.StatusOK, quiz)
}

// UpdateQuiz updates quiz settings
func UpdateQuiz(c echo.Context) error {
	initQuizRepos()
	
	id := c.Param("id")
	
	quiz, err := quizRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch quiz"})
	}
	if quiz == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Quiz not found"})
	}
	
	var req domain.UpdateQuizRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if req.Title != nil {
		quiz.Title = *req.Title
	}
	if req.Description != nil {
		quiz.Description = *req.Description
	}
	if req.TimeLimit != nil {
		quiz.TimeLimit = *req.TimeLimit
	}
	if req.PassingScore != nil {
		quiz.PassingScore = *req.PassingScore
	}
	if req.MaxAttempts != nil {
		quiz.MaxAttempts = *req.MaxAttempts
	}
	if req.ShuffleQuestions != nil {
		quiz.ShuffleQuestions = *req.ShuffleQuestions
	}
	if req.ShuffleOptions != nil {
		quiz.ShuffleOptions = *req.ShuffleOptions
	}
	if req.ShowCorrectAnswers != nil {
		quiz.ShowCorrectAnswers = *req.ShowCorrectAnswers
	}
	
	err = quizRepo.Update(quiz)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update quiz"})
	}
	
	return c.JSON(http.StatusOK, quiz)
}

// DeleteQuiz deletes a quiz
func DeleteQuiz(c echo.Context) error {
	initQuizRepos()
	
	id := c.Param("id")
	
	err := quizRepo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete quiz"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Quiz deleted successfully"})
}

// CreateQuestion adds a question to a quiz
func CreateQuestion(c echo.Context) error {
	initQuizRepos()
	
	quizID := c.Param("quizId")
	
	var req domain.CreateQuestionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if req.QuestionText == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Question text is required"})
	}
	
	question := &domain.Question{
		QuizID:       quizID,
		Type:         req.Type,
		QuestionText: req.QuestionText,
		Explanation:  req.Explanation,
		Points:       req.Points,
		Required:     req.Required,
	}
	
	if question.Points == 0 {
		question.Points = 1
	}
	if question.Type == "" {
		question.Type = domain.QuestionMultipleChoice
	}
	
	err := quizRepo.CreateQuestion(question)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create question: " + err.Error()})
	}
	
	// Create options if provided
	if len(req.Options) > 0 {
		options := make([]domain.Option, len(req.Options))
		for i, opt := range req.Options {
			options[i] = domain.Option{
				QuestionID: question.ID,
				OptionText: opt.OptionText,
				IsCorrect:  opt.IsCorrect,
			}
		}
		err = quizRepo.CreateOptions(options)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create options: " + err.Error()})
		}
		question.Options = options
	}
	
	return c.JSON(http.StatusCreated, question)
}

// UpdateQuestion updates a question
func UpdateQuestion(c echo.Context) error {
	initQuizRepos()
	
	id := c.Param("id")
	
	question, err := quizRepo.GetQuestionByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch question"})
	}
	if question == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Question not found"})
	}
	
	var req domain.UpdateQuestionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if req.Type != nil {
		question.Type = *req.Type
	}
	if req.QuestionText != nil {
		question.QuestionText = *req.QuestionText
	}
	if req.Explanation != nil {
		question.Explanation = *req.Explanation
	}
	if req.Points != nil {
		question.Points = *req.Points
	}
	if req.Required != nil {
		question.Required = *req.Required
	}
	
	err = quizRepo.UpdateQuestion(question)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update question"})
	}
	
	// Replace options if provided
	if req.Options != nil {
		quizRepo.DeleteOptionsByQuestionID(id)
		
		options := make([]domain.Option, len(req.Options))
		for i, opt := range req.Options {
			options[i] = domain.Option{
				QuestionID: question.ID,
				OptionText: opt.OptionText,
				IsCorrect:  opt.IsCorrect,
			}
		}
		err = quizRepo.CreateOptions(options)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update options"})
		}
		question.Options = options
	}
	
	return c.JSON(http.StatusOK, question)
}

// DeleteQuestion deletes a question
func DeleteQuestion(c echo.Context) error {
	initQuizRepos()
	
	id := c.Param("id")
	
	err := quizRepo.DeleteQuestion(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete question"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Question deleted successfully"})
}

// ReorderQuestions updates question order
func ReorderQuestions(c echo.Context) error {
	initQuizRepos()
	
	quizID := c.Param("quizId")
	
	var req struct {
		QuestionIDs []string `json:"question_ids"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	err := quizRepo.ReorderQuestions(quizID, req.QuestionIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reorder questions"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Questions reordered successfully"})
}

// ============ Student Handlers ============

// GetQuizForStudent retrieves a quiz without correct answers
func GetQuizForStudent(c echo.Context) error {
	initQuizRepos()
	
	lessonID := c.Param("lessonId")
	
	quiz, err := quizRepo.GetByLessonID(lessonID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch quiz"})
	}
	if quiz == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Quiz not found"})
	}
	
	// Load questions
	quiz.Questions, err = quizRepo.GetQuestionsByQuizID(quiz.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch questions"})
	}
	
	// Remove correct answer info for students
	for i := range quiz.Questions {
		for j := range quiz.Questions[i].Options {
			quiz.Questions[i].Options[j].IsCorrect = false
		}
	}
	
	// Shuffle if enabled
	if quiz.ShuffleQuestions {
		sort.Slice(quiz.Questions, func(i, j int) bool {
			return time.Now().UnixNano()%2 == 0
		})
	}
	
	return c.JSON(http.StatusOK, quiz)
}

// GetQuizStatus returns quiz status for a user including attempt history and remaining attempts
func GetQuizStatus(c echo.Context) error {
	initQuizRepos()
	
	quizID := c.Param("quizId")
	userID := getQuizUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}
	
	quiz, err := quizRepo.GetByID(quizID)
	if err != nil || quiz == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Quiz not found"})
	}
	
	// Get user's attempts
	attempts, _ := quizRepo.GetAttemptsByQuizAndUser(quizID, userID)
	attemptCount := len(attempts)
	
	// Calculate remaining attempts
	remainingAttempts := -1 // -1 means unlimited
	if quiz.MaxAttempts > 0 {
		remainingAttempts = quiz.MaxAttempts - attemptCount
		if remainingAttempts < 0 {
			remainingAttempts = 0
		}
	}
	
	// Find best score
	var bestScore *int
	var hasPassed bool
	for _, attempt := range attempts {
		if attempt.Score != nil {
			if bestScore == nil || *attempt.Score > *bestScore {
				bestScore = attempt.Score
			}
			if attempt.Passed != nil && *attempt.Passed {
				hasPassed = true
			}
		}
	}
	
	// Check for in-progress attempt
	var inProgressAttempt *domain.QuizAttempt
	for i := range attempts {
		if attempts[i].CompletedAt == nil {
			inProgressAttempt = &attempts[i]
			break
		}
	}
	
	canAttempt := remainingAttempts != 0
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"quiz_id":             quizID,
		"title":               quiz.Title,
		"description":         quiz.Description,
		"time_limit":          quiz.TimeLimit,
		"passing_score":       quiz.PassingScore,
		"max_attempts":        quiz.MaxAttempts,
		"attempt_count":       attemptCount,
		"remaining_attempts":  remainingAttempts,
		"best_score":          bestScore,
		"has_passed":          hasPassed,
		"can_attempt":         canAttempt,
		"in_progress_attempt": inProgressAttempt,
		"attempts":            attempts,
	})
}

// StartQuizAttempt starts a new quiz attempt
func StartQuizAttempt(c echo.Context) error {
	initQuizRepos()
	
	quizID := c.Param("quizId")
	userID := getQuizUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}
	
	quiz, err := quizRepo.GetByID(quizID)
	if err != nil || quiz == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Quiz not found"})
	}
	
	// Check attempt limit
	if quiz.MaxAttempts > 0 {
		count, _ := quizRepo.CountAttemptsByQuizAndUser(quizID, userID)
		if count >= quiz.MaxAttempts {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Maximum attempts reached"})
		}
	}
	
	attempt := &domain.QuizAttempt{
		QuizID: quizID,
		UserID: userID,
	}
	
	err = quizRepo.CreateAttempt(attempt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start attempt"})
	}
	
	return c.JSON(http.StatusCreated, attempt)
}

// SubmitQuizAttempt submits answers and calculates score
func SubmitQuizAttempt(c echo.Context) error {
	initQuizRepos()
	
	attemptID := c.Param("attemptId")
	
	attempt, err := quizRepo.GetAttemptByID(attemptID)
	if err != nil || attempt == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Attempt not found"})
	}
	
	if attempt.CompletedAt != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Attempt already submitted"})
	}
	
	var req domain.SubmitQuizRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	quiz, _ := quizRepo.GetByID(attempt.QuizID)
	questions, _ := quizRepo.GetQuestionsByQuizID(attempt.QuizID)
	
	// Build question map for lookup
	questionMap := make(map[string]domain.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}
	
	// Calculate scores
	totalPoints := 0
	earnedPoints := 0
	correctCount := 0
	
	var answers []domain.Answer
	
	for _, ans := range req.Answers {
		q, exists := questionMap[ans.QuestionID]
		if !exists {
			continue
		}
		
		totalPoints += q.Points
		isCorrect := false
		pointsEarned := 0
		
		switch q.Type {
		case domain.QuestionMultipleChoice, domain.QuestionTrueFalse:
			// Check if selected option is correct
			if len(ans.SelectedOptionIDs) == 1 {
				for _, opt := range q.Options {
					if opt.ID == ans.SelectedOptionIDs[0] && opt.IsCorrect {
						isCorrect = true
						break
					}
				}
			}
			
		case domain.QuestionMultipleAnswer:
			// Check if all correct options are selected and no incorrect ones
			correctOpts := make(map[string]bool)
			for _, opt := range q.Options {
				if opt.IsCorrect {
					correctOpts[opt.ID] = true
				}
			}
			
			selectedSet := make(map[string]bool)
			for _, id := range ans.SelectedOptionIDs {
				selectedSet[id] = true
			}
			
			if len(selectedSet) == len(correctOpts) {
				allCorrect := true
				for id := range selectedSet {
					if !correctOpts[id] {
						allCorrect = false
						break
					}
				}
				isCorrect = allCorrect
			}
			
		case domain.QuestionShortAnswer:
			// Simple exact match (case insensitive, trimmed)
			for _, opt := range q.Options {
				if strings.TrimSpace(strings.ToLower(ans.TextAnswer)) == 
				   strings.TrimSpace(strings.ToLower(opt.OptionText)) {
					isCorrect = true
					break
				}
			}
		}
		
		if isCorrect {
			pointsEarned = q.Points
			earnedPoints += pointsEarned
			correctCount++
		}
		
		answers = append(answers, domain.Answer{
			AttemptID:         attemptID,
			QuestionID:        ans.QuestionID,
			SelectedOptionIDs: ans.SelectedOptionIDs,
			TextAnswer:        ans.TextAnswer,
			IsCorrect:         &isCorrect,
			PointsEarned:      pointsEarned,
		})
	}
	
	// Save answers
	err = quizRepo.CreateAnswers(answers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save answers"})
	}
	
	// Calculate percentage and update attempt
	score := 0
	if totalPoints > 0 {
		score = (earnedPoints * 100) / totalPoints
	}
	
	now := time.Now()
	timeSpent := int(now.Sub(attempt.StartedAt).Seconds())
	passed := score >= quiz.PassingScore
	
	attempt.CompletedAt = &now
	attempt.Score = &score
	attempt.Passed = &passed
	attempt.TimeSpent = &timeSpent
	
	err = quizRepo.UpdateAttempt(attempt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update attempt"})
	}
	
	result := domain.QuizResult{
		AttemptID:      attemptID,
		QuizID:         attempt.QuizID,
		Score:          score,
		Passed:         passed,
		TimeSpent:      timeSpent,
		TotalPoints:    totalPoints,
		EarnedPoints:   earnedPoints,
		CorrectCount:   correctCount,
		TotalQuestions: len(questions),
	}
	
	// Add answer details if quiz allows showing correct answers
	if quiz.ShowCorrectAnswers {
		for i, ans := range answers {
			q := questionMap[ans.QuestionID]
			
			correctAnswers := []string{}
			for _, opt := range q.Options {
				if opt.IsCorrect {
					correctAnswers = append(correctAnswers, opt.ID)
				}
			}
			
			result.Answers = append(result.Answers, domain.AnswerResult{
				QuestionID:    ans.QuestionID,
				QuestionText:  q.QuestionText,
				UserAnswer:    ans.SelectedOptionIDs,
				CorrectAnswer: correctAnswers,
				IsCorrect:     *answers[i].IsCorrect,
				PointsEarned:  ans.PointsEarned,
				MaxPoints:     q.Points,
				Explanation:   q.Explanation,
			})
		}
	}
	
	return c.JSON(http.StatusOK, result)
}

// GetQuizAttemptResult retrieves the result of a completed attempt
func GetQuizAttemptResult(c echo.Context) error {
	initQuizRepos()
	
	attemptID := c.Param("attemptId")
	
	attempt, err := quizRepo.GetAttemptByID(attemptID)
	if err != nil || attempt == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Attempt not found"})
	}
	
	if attempt.CompletedAt == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Attempt not yet completed"})
	}
	
	quiz, _ := quizRepo.GetByID(attempt.QuizID)
	answers, _ := quizRepo.GetAnswersByAttemptID(attemptID)
	questions, _ := quizRepo.GetQuestionsByQuizID(attempt.QuizID)
	
	questionMap := make(map[string]domain.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}
	
	result := domain.QuizResult{
		AttemptID:      attemptID,
		QuizID:         attempt.QuizID,
		Score:          *attempt.Score,
		Passed:         *attempt.Passed,
		TimeSpent:      *attempt.TimeSpent,
		TotalQuestions: len(questions),
	}
	
	if quiz.ShowCorrectAnswers {
		for _, ans := range answers {
			q := questionMap[ans.QuestionID]
			
			correctAnswers := []string{}
			for _, opt := range q.Options {
				if opt.IsCorrect {
					correctAnswers = append(correctAnswers, opt.ID)
				}
			}
			
			result.Answers = append(result.Answers, domain.AnswerResult{
				QuestionID:    ans.QuestionID,
				QuestionText:  q.QuestionText,
				UserAnswer:    ans.SelectedOptionIDs,
				CorrectAnswer: correctAnswers,
				IsCorrect:     *ans.IsCorrect,
				PointsEarned:  ans.PointsEarned,
				MaxPoints:     q.Points,
				Explanation:   q.Explanation,
			})
		}
	}
	
	return c.JSON(http.StatusOK, result)
}

// GetUserQuizAttempts retrieves all attempts by the current user for a quiz
func GetUserQuizAttempts(c echo.Context) error {
	initQuizRepos()
	
	quizID := c.Param("quizId")
	userID := getQuizUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}
	
	attempts, err := quizRepo.GetAttemptsByQuizAndUser(quizID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch attempts"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"attempts": attempts,
		"total":    len(attempts),
	})
}
