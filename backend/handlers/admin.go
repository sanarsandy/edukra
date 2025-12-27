package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

var transactionRepo *postgres.TransactionRepository

func initAdminRepos() {
	if transactionRepo == nil && db.DB != nil {
		transactionRepo = postgres.NewTransactionRepository(db.DB)
	}
	initUserRepos()
	initCourseRepos()
	initEnrollmentRepos()
}

// TransactionWithDetails includes user and course info
type TransactionWithDetails struct {
	*postgres.Transaction
	User   *UserBrief   `json:"user,omitempty"`
	Course *CourseBrief `json:"course,omitempty"`
}

// UserBrief is a simplified user info for transactions
type UserBrief struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

// CourseBrief is a simplified course info for transactions
type CourseBrief struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// enrichTransaction adds user and course details to a transaction
func enrichTransaction(tx *postgres.Transaction) *TransactionWithDetails {
	result := &TransactionWithDetails{Transaction: tx}
	
	// Get user info
	if tx.UserID != "" && userRepo != nil {
		user, err := userRepo.GetByID(tx.UserID)
		if err == nil && user != nil {
			result.User = &UserBrief{
				ID:       user.ID,
				Email:    user.Email,
				FullName: user.FullName,
			}
		}
	}
	
	// Get course info
	if tx.CourseID != nil && *tx.CourseID != "" && courseRepo != nil {
		course, err := courseRepo.GetByID(*tx.CourseID)
		if err == nil && course != nil {
			result.Course = &CourseBrief{
				ID:    course.ID,
				Title: course.Title,
			}
		}
	}
	
	return result
}

// GetAdminDashboard returns dashboard statistics
func GetAdminDashboard(c echo.Context) error {
	initAdminRepos()
	
	tenantID := c.QueryParam("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	
	// Get counts
	userCount, _ := userRepo.CountByTenant(tenantID)
	courseCount, _ := courseRepo.CountByTenant(tenantID)
	transactionCount, _ := transactionRepo.CountByTenant(tenantID)
	totalRevenue, _ := transactionRepo.SumByTenant(tenantID)
	
	// Get recent users
	recentUsers, _ := userRepo.ListByTenant(tenantID, 5, 0)
	
	// Get recent transactions with user/course details
	recentTransactions, _ := transactionRepo.ListByTenant(tenantID, 5, 0)
	var enrichedTransactions []*TransactionWithDetails
	for _, tx := range recentTransactions {
		enrichedTransactions = append(enrichedTransactions, enrichTransaction(tx))
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"stats": map[string]interface{}{
			"total_users":        userCount,
			"total_courses":      courseCount,
			"total_transactions": transactionCount,
			"total_revenue":      totalRevenue,
		},
		"recent_users":        recentUsers,
		"recent_transactions": enrichedTransactions,
	})
}

// ListTransactions returns all transactions (admin only)
func ListTransactions(c echo.Context) error {
	initAdminRepos()
	
	tenantID := c.QueryParam("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	
	status := c.QueryParam("status")
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}
	
	var transactions []*postgres.Transaction
	var err error
	
	if status != "" {
		transactions, err = transactionRepo.ListByStatus(tenantID, status, limit, offset)
	} else {
		transactions, err = transactionRepo.ListByTenant(tenantID, limit, offset)
	}
	
	if err != nil {
		log.Printf("[ListTransactions] Error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch transactions"})
	}
	
	// Enrich transactions with user and course details
	var enrichedTransactions []*TransactionWithDetails
	for _, tx := range transactions {
		enrichedTransactions = append(enrichedTransactions, enrichTransaction(tx))
	}
	
	total, _ := transactionRepo.CountByTenant(tenantID)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactions": enrichedTransactions,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

// GetTransaction returns a single transaction
func GetTransaction(c echo.Context) error {
	initAdminRepos()
	
	id := c.Param("id")
	transaction, err := transactionRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch transaction"})
	}
	
	if transaction == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Transaction not found"})
	}
	
	// Return enriched transaction with user and course details
	return c.JSON(http.StatusOK, enrichTransaction(transaction))
}

// UpdateTransactionStatus updates transaction status (admin/webhook)
func UpdateTransactionStatus(c echo.Context) error {
	initAdminRepos()
	
	id := c.Param("id")
	
	var req struct {
		Status string `json:"status"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	validStatuses := map[string]bool{"pending": true, "success": true, "failed": true, "refunded": true}
	if !validStatuses[req.Status] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid status"})
	}
	
	transaction, err := transactionRepo.GetByID(id)
	if err != nil || transaction == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Transaction not found"})
	}
	
	err = transactionRepo.UpdateStatus(id, req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update transaction"})
	}
	
	// If payment successful, create enrollment
	if req.Status == "success" && transaction.CourseID != nil {
		enrollment := &postgres.Enrollment{
			UserID:        transaction.UserID,
			CourseID:      *transaction.CourseID,
			TransactionID: &id,
		}
		enrollmentRepo.Create(enrollment)
	}
	
	transaction, _ = transactionRepo.GetByID(id)
	return c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction deletes a transaction (only non-success)
func DeleteTransaction(c echo.Context) error {
	initAdminRepos()
	
	id := c.Param("id")
	
	transaction, err := transactionRepo.GetByID(id)
	if err != nil || transaction == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Transaction not found"})
	}
	
	// Prevent deletion of successful transactions
	if transaction.Status == "success" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot delete successful transactions"})
	}
	
	err = transactionRepo.Delete(id)
	if err != nil {
		log.Printf("[DeleteTransaction] Error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete transaction"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Transaction deleted successfully"})
}

// GetDashboardChartData returns monthly stats for charts
func GetDashboardChartData(c echo.Context) error {
	initAdminRepos()
	
	tenantID := c.QueryParam("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	
	// Get revenue data
	revenue, months, err := transactionRepo.GetMonthlyRevenue(tenantID)
	if err != nil {
		log.Printf("[Dashboard] Error fetching revenue: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch revenue data"})
	}
	
	// Get user growth data
	users, _, err := userRepo.GetMonthlyGrowth(tenantID)
	if err != nil {
		log.Printf("[Dashboard] Error fetching user growth: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user growth data"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"labels":  months,
		"revenue": revenue,
		"users":   users,
	})
}

// GetUserDashboard returns user dashboard statistics
func GetUserDashboard(c echo.Context) error {
	initAdminRepos()
	
	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}
	
	// Get user's stats
	enrolledCount, _ := enrollmentRepo.CountByUser(userID)
	completedCount, _ := enrollmentRepo.CountCompletedByUser(userID)
	
	// Get recent enrollments with courses
	enrollments, _ := enrollmentRepo.ListByUser(userID, 5, 0)
	
	type EnrollmentWithCourse struct {
		*postgres.Enrollment
		Course interface{} `json:"course,omitempty"`
	}
	
	var enriched []EnrollmentWithCourse
	for _, e := range enrollments {
		ec := EnrollmentWithCourse{Enrollment: e}
		if courseRepo != nil {
			course, _ := courseRepo.GetByID(e.CourseID)
			ec.Course = course
		}
		enriched = append(enriched, ec)
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"stats": map[string]interface{}{
			"enrolled_courses":  enrolledCount,
			"completed_courses": completedCount,
			"in_progress":       enrolledCount - completedCount,
		},
		"recent_courses": enriched,
	})
}
