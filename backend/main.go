package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/handlers"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/scheduler"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/service"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/storage"
	customMiddleware "github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
)

func main() {
	// Set Global Timezone to Asia/Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Warning: Failed to load Asia/Jakarta timezone: %v", err)
	} else {
		time.Local = loc
		log.Println("Global Timezone set to Asia/Jakarta")
	}

	// Initialize Database
	db.Init()

	// Initialize MinIO Storage (optional - will skip if not configured)
	if err := storage.InitStorage(); err != nil {
		log.Printf("Warning: Failed to initialize MinIO storage: %v", err)
		log.Println("File uploads will not work without MinIO configuration")
	}

	// Initialize WhatsApp Service
	service.InitWhatsAppService()

	// Initialize and start reminder scheduler
	scheduler.InitScheduler(db.DB)
	scheduler.StartScheduler()
	defer scheduler.StopScheduler()

	e := EchoServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

func EchoServer() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// CORS Configuration
	corsOrigins := []string{"http://localhost:3000"}
	if corsEnv := os.Getenv("CORS_ALLOWED_ORIGINS"); corsEnv != "" {
		corsOrigins = append(corsOrigins, strings.Split(corsEnv, ",")...)
	}
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to API",
			"status":  "healthy",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "up",
		})
	})

	// Auth Routes (Public)
	auth := e.Group("/api/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/admin/login", handlers.AdminLogin)
	auth.POST("/instructor/login", handlers.InstructorLogin)
	auth.GET("/google", handlers.GetGoogleAuthURL)
	auth.GET("/google/callback", handlers.GoogleAuthCallback)

	// Public Course Routes (no auth required for browsing)
	e.GET("/api/courses", handlers.ListCourses)
	e.GET("/api/courses/:id", handlers.GetCourse)
	e.GET("/api/categories", handlers.ListCategories)
	
	// Public Campaign Routes (landing pages) - with rate limiting
	e.GET("/api/c/:slug", handlers.GetCampaignBySlug, customMiddleware.PublicAPIRateLimiter.Middleware())
	e.POST("/api/c/:id/click", handlers.TrackCampaignClick, customMiddleware.TrackingRateLimiter.Middleware())
	e.POST("/api/c/:id/track", handlers.TrackCampaignClick, customMiddleware.TrackingRateLimiter.Middleware())

	// Public Campaign Checkout (guest checkout - with strict rate limiting)
	e.POST("/api/campaign-checkout", handlers.CampaignCheckout, customMiddleware.CheckoutRateLimiter.Middleware())
	e.GET("/api/transaction-status/:order_id", handlers.GetTransactionStatus)
	e.GET("/api/public/payment-methods", handlers.GetPaymentMethods) // Public for campaign checkout
	e.GET("/api/settings", handlers.GetSettings) // Public settings (banner, site info)

	// Public Blog Routes
	e.GET("/api/blog", handlers.ListBlogPostsPublic)
	e.GET("/api/blog/:slug", handlers.GetBlogPostBySlug)

	// Public Webinar Routes
	e.GET("/api/webinars/:id", handlers.GetPublicWebinar)
	e.GET("/api/courses/:id/webinars", handlers.GetCourseWebinars)

	// Protected Routes
	api := e.Group("/api")
	api.Use(customMiddleware.JWTMiddleware())
	
	// User Profile & Auth
	api.GET("/me", handlers.GetMe)
	api.PUT("/me", handlers.UpdateCurrentUser)
	api.PUT("/me/password", handlers.ChangePassword)
	api.POST("/auth/refresh", handlers.RefreshToken)
	api.POST("/auth/logout", handlers.Logout)
	
	// User Dashboard
	api.GET("/dashboard", handlers.GetUserDashboard)
	
	// Enrollments
	api.GET("/enrollments", handlers.ListMyEnrollments)
	api.POST("/enrollments", handlers.EnrollInCourse)
	api.GET("/enrollments/:id", handlers.GetEnrollment)
	api.PUT("/enrollments/:id/progress", handlers.UpdateEnrollmentProgress)
	api.GET("/courses/:courseId/enrollment", handlers.CheckEnrollment)
	
	// Lessons (protected for enrolled users)
	api.GET("/lessons/:id", handlers.GetLesson)
	api.GET("/courses/:courseId/lessons", handlers.ListLessons)

	// Lesson Progress
	api.GET("/lessons/:lessonId/progress", handlers.GetLessonProgress)
	api.POST("/lessons/:lessonId/progress", handlers.UpdateLessonProgressHandler)
	api.POST("/lessons/:lessonId/watchtime", handlers.UpdateWatchTime)
	api.GET("/courses/:courseId/progress", handlers.GetCourseProgressHandler)
	api.POST("/courses/:courseId/progress/bulk", handlers.BulkUpdateProgress)

	// Secure Content Access (pre-signed URLs)
	api.GET("/content/:lessonId/url", handlers.GetSecureContentURL)
	api.GET("/content/:lessonId/document", handlers.GetSecureDocumentURL)
	// Stream content directly through backend (bypasses pre-signed URL issues)
	api.GET("/content/:lessonId/stream", handlers.StreamContent)
	
	// HLS Encrypted Video Streaming
	api.GET("/content/:lessonId/hls/manifest", handlers.GetHLSManifest)
	api.GET("/content/:lessonId/hls/segment/:filename", handlers.GetHLSSegment)
	api.GET("/content/:lessonId/hls/key", handlers.GetHLSKey)
	api.GET("/content/:lessonId/hls/status", handlers.GetHLSStatus)
	
	// External PDF Proxy (bypass CORS for external PDFs)
	api.POST("/content/proxy-pdf", handlers.ProxyExternalPDF)
	
	// Public images (thumbnails, etc) - no auth required
	e.GET("/api/images/:objectKey", handlers.GetPublicImage)

	// Activities & Stats
	api.GET("/activities", handlers.GetRecentActivities)
	api.GET("/activities/weekly", handlers.GetWeeklyActivities)
	api.GET("/stats", handlers.GetLearningStats)

	// Certificates
	api.GET("/certificates", handlers.ListMyCertificates)
	api.GET("/certificates/:id", handlers.GetCertificate)
	api.GET("/certificates/:id/download", handlers.DownloadCertificatePDF)

	// Course Ratings
	api.GET("/courses/:courseId/ratings", handlers.GetCourseRatings)
	api.GET("/courses/:courseId/ratings/stats", handlers.GetCourseRatingStats)
	api.GET("/courses/:courseId/my-rating", handlers.GetMyRating)
	api.POST("/courses/:courseId/ratings", handlers.CreateCourseRating)
	api.PUT("/courses/:courseId/ratings", handlers.UpdateCourseRating)
	api.DELETE("/courses/:courseId/ratings", handlers.DeleteCourseRating)

	// AI Tutor Chat (for students)
	api.POST("/courses/:id/chat", handlers.SendChatMessage)
	api.GET("/courses/:id/chat/session", handlers.GetChatSession)
	api.DELETE("/courses/:id/chat/session", handlers.ClearChatSession)
	api.GET("/courses/:id/chat/quota", handlers.GetChatQuota)
	api.GET("/courses/:id/ai-status", handlers.GetAIStatus)

	// Notifications (all authenticated users)
	api.GET("/notifications", handlers.GetMyNotifications)
	api.PUT("/notifications/:id/read", handlers.MarkNotificationAsRead)
	api.PUT("/notifications/read-all", handlers.MarkAllNotificationsRead)

	// Student Webinars
	api.GET("/my/webinars", handlers.GetMyWebinars)

	// Payment & Checkout
	api.POST("/checkout", handlers.CreateCheckout)
	api.GET("/checkout/config", handlers.GetCheckoutConfig)
	api.GET("/checkout/payment-methods", handlers.GetPaymentMethods)
	api.POST("/coupons/validate", handlers.ValidateCoupon) // Validate coupon at checkout
	api.GET("/my/transactions", handlers.GetMyTransactions)
	api.GET("/enrollments/check/:id", handlers.CheckEnrollment)
	
	// Payment Webhooks (no auth - verified by signature)
	e.POST("/api/webhooks/midtrans", handlers.MidtransWebhook)
	e.POST("/api/webhooks/duitku", handlers.DuitkuWebhook)
	
	// Test endpoint for simulating payments (development only - protected by admin auth)
	e.POST("/api/test/simulate-payment", handlers.SimulatePaymentSuccess, customMiddleware.JWTMiddleware(), customMiddleware.RequireAdmin())
	
	// Admin Routes (requires admin role)
	admin := e.Group("/api/admin")
	admin.Use(customMiddleware.JWTMiddleware())
	admin.Use(customMiddleware.RequireAdmin())
	
	// Admin Dashboard
	admin.GET("/dashboard", handlers.GetAdminDashboard)
	admin.GET("/dashboard/charts", handlers.GetDashboardChartData)
	admin.GET("/dashboard/activities", handlers.GetAdminRecentActivities)
	
	// Admin User Management
	admin.GET("/users", handlers.ListUsers)
	admin.POST("/users", handlers.CreateUser)
	admin.GET("/users/:id", handlers.GetUser)
	admin.PUT("/users/:id", handlers.UpdateUser)
	admin.DELETE("/users/:id", handlers.DeleteUser)
	
	// Admin Course Management
	admin.GET("/courses", handlers.AdminListCourses)
	admin.POST("/courses", handlers.CreateCourse)
	admin.PUT("/courses/:id", handlers.UpdateCourse)
	admin.DELETE("/courses/:id", handlers.DeleteCourse)
	admin.PUT("/courses/:id/publish", handlers.PublishCourse)
	
	// Admin Lesson Management
	admin.POST("/courses/:courseId/lessons", handlers.CreateLesson)
	admin.GET("/courses/:courseId/lessons/tree", handlers.GetLessonTree)
	admin.PUT("/lessons/:id", handlers.UpdateLesson)
	admin.PUT("/lessons/:id/move", handlers.MoveLesson)
	admin.DELETE("/lessons/:id", handlers.DeleteLesson)
	admin.PUT("/courses/:courseId/lessons/reorder", handlers.ReorderLessons)

	
	// Admin Transactions
	admin.GET("/transactions", handlers.ListTransactions)
	admin.GET("/transactions/:id", handlers.GetTransaction)
	admin.PUT("/transactions/:id/status", handlers.UpdateTransactionStatus)
	admin.DELETE("/transactions/:id", handlers.DeleteTransaction)

	// Admin Payment Settings
	admin.GET("/payment/settings", handlers.GetPaymentSettings)
	admin.PUT("/payment/settings", handlers.UpdatePaymentSettings)

	// Admin Categories
	admin.GET("/categories", handlers.ListCategories)
	admin.GET("/categories/:id", handlers.GetCategory)
	admin.POST("/categories", handlers.CreateCategory)
	admin.PUT("/categories/:id", handlers.UpdateCategory)
	admin.DELETE("/categories/:id", handlers.DeleteCategory)

	// Admin Instructors
	admin.GET("/instructors", handlers.ListInstructors)
	admin.GET("/instructors/:id", handlers.GetInstructor)
	admin.POST("/instructors", handlers.CreateInstructor)
	admin.PUT("/instructors/:id", handlers.UpdateInstructor)
	admin.DELETE("/instructors/:id", handlers.DeleteInstructor)

	// Admin Settings
	admin.GET("/settings", handlers.GetSettings)
	admin.PUT("/settings", handlers.UpdateSettings)

	// Admin AI Settings
	admin.GET("/ai/settings", handlers.GetAISettings)
	admin.PUT("/ai/settings", handlers.UpdateAISettings)
	admin.POST("/ai/validate-key", handlers.ValidateAIKey)
	admin.GET("/ai/providers", handlers.GetAIProviders)
	admin.GET("/ai/models", handlers.FetchProviderModels) // Fetch models from provider API
	admin.DELETE("/ai/key", handlers.ClearAPIKey)

	// Admin AI Content Processing
	admin.POST("/courses/:id/process-ai", handlers.ProcessCourseContent)
	admin.GET("/courses/:id/ai-processing-status", handlers.GetProcessingStatus)
	admin.DELETE("/courses/:id/embeddings", handlers.ClearCourseEmbeddings)

	// Admin File Upload
	admin.POST("/upload", handlers.UploadFile)

	// Admin Rating Management
	admin.GET("/ratings", handlers.AdminGetAllRatings)
	admin.GET("/ratings/stats", handlers.AdminGetRatingStats)
	admin.GET("/courses/:courseId/ratings", handlers.AdminGetCourseRatings)
	admin.DELETE("/ratings/:id", handlers.AdminDeleteRating)


	// Admin Campaign Management
	admin.GET("/campaigns", handlers.ListCampaigns)
	admin.POST("/campaigns", handlers.CreateCampaign)
	admin.GET("/campaigns/:id", handlers.GetCampaign)
	admin.PUT("/campaigns/:id", handlers.UpdateCampaign)
	admin.DELETE("/campaigns/:id", handlers.DeleteCampaign)
	admin.GET("/campaigns/:id/analytics", handlers.GetCampaignAnalytics)

	// Admin Blog Management
	admin.GET("/blog", handlers.ListBlogPostsAdmin)
	admin.POST("/blog", handlers.CreateBlogPost)
	admin.GET("/blog/:id", handlers.GetBlogPostAdmin)
	admin.PUT("/blog/:id", handlers.UpdateBlogPost)
	admin.DELETE("/blog/:id", handlers.DeleteBlogPost)
	admin.GET("/blog-categories", handlers.ListBlogCategories)
	admin.POST("/blog-categories", handlers.CreateBlogCategory)
	admin.DELETE("/blog-categories/:id", handlers.DeleteBlogCategory)

	// Admin Quiz Management
	admin.POST("/lessons/:lessonId/quiz", handlers.CreateQuiz)
	admin.GET("/lessons/:lessonId/quiz", handlers.GetQuiz)
	admin.GET("/quizzes/:id", handlers.GetQuiz)
	admin.PUT("/quizzes/:id", handlers.UpdateQuiz)
	admin.DELETE("/quizzes/:id", handlers.DeleteQuiz)
	admin.POST("/quizzes/:quizId/questions", handlers.CreateQuestion)
	admin.PUT("/questions/:id", handlers.UpdateQuestion)
	admin.DELETE("/questions/:id", handlers.DeleteQuestion)
	admin.PUT("/quizzes/:quizId/questions/reorder", handlers.ReorderQuestions)

	// Admin Course Review Management (for instructor workflow)
	admin.GET("/reviews", handlers.AdminListPendingReviews)
	admin.GET("/reviews/stats", handlers.AdminReviewStats)
	admin.GET("/reviews/:id", handlers.AdminReviewCourseDetail)
	admin.POST("/reviews/:id/approve", handlers.AdminApproveCourse)
	admin.POST("/reviews/:id/reject", handlers.AdminRejectCourse)
	admin.POST("/reviews/:id/publish", handlers.AdminPublishCourse)
	admin.POST("/reviews/:id/unpublish", handlers.AdminUnpublishCourse)

	// Admin Coupon Management
	admin.GET("/coupons", handlers.ListCoupons)
	admin.POST("/coupons", handlers.CreateCoupon)
	admin.GET("/coupons/:id", handlers.GetCoupon)
	admin.PUT("/coupons/:id", handlers.UpdateCoupon)
	admin.DELETE("/coupons/:id", handlers.DeleteCoupon)

	// Admin Webinar Management
	admin.GET("/webinars", handlers.ListWebinars)
	admin.POST("/webinars", handlers.CreateWebinar)
	admin.GET("/webinars/:id", handlers.GetWebinar)
	admin.PUT("/webinars/:id", handlers.UpdateWebinar)
	admin.DELETE("/webinars/:id", handlers.DeleteWebinar)
	admin.GET("/webinars/:id/registrations", handlers.GetWebinarRegistrations)
	admin.POST("/webinars/:id/attendance/:user_id", handlers.MarkWebinarAttendance)
	admin.GET("/courses/:id/webinars", handlers.GetWebinarsByCourse)

	// ========================================
	// INSTRUCTOR ROUTES
	// ========================================
	instructor := e.Group("/api/instructor")
	instructor.Use(customMiddleware.JWTMiddleware())
	instructor.Use(customMiddleware.RequireInstructor())

	// Instructor Dashboard
	instructor.GET("/dashboard", handlers.InstructorDashboard)

	// Instructor Course Management
	instructor.GET("/courses", handlers.InstructorListCourses)
	instructor.POST("/courses", handlers.InstructorCreateCourse)
	instructor.GET("/courses/:id", handlers.InstructorGetCourse)
	instructor.PUT("/courses/:id", handlers.InstructorUpdateCourse)
	instructor.DELETE("/courses/:id", handlers.InstructorDeleteCourse)
	instructor.POST("/courses/:id/submit", handlers.InstructorSubmitCourse)

	// Instructor Lesson Management
	instructor.GET("/courses/:courseId/lessons", handlers.InstructorGetLessonTree)
	instructor.GET("/courses/:courseId/lessons/tree", handlers.InstructorGetLessonTree) // alias for tree endpoint
	instructor.POST("/courses/:courseId/lessons", handlers.InstructorCreateLesson)
	instructor.PUT("/lessons/:id", handlers.InstructorUpdateLesson)
	instructor.DELETE("/lessons/:id", handlers.InstructorDeleteLesson)

	// Instructor Quiz Management
	instructor.POST("/lessons/:lessonId/quiz", handlers.InstructorCreateQuiz)
	instructor.GET("/lessons/:lessonId/quiz", handlers.InstructorGetQuiz)
	instructor.PUT("/quizzes/:id", handlers.InstructorUpdateQuiz)
	instructor.DELETE("/quizzes/:id", handlers.InstructorDeleteQuiz)
	instructor.POST("/quizzes/:quizId/questions", handlers.InstructorAddQuestion)
	instructor.PUT("/questions/:id", handlers.InstructorUpdateQuestion)
	instructor.DELETE("/questions/:id", handlers.InstructorDeleteQuestion)

	// Instructor Analytics
	instructor.GET("/courses/:id/students", handlers.InstructorCourseStudents)
	instructor.GET("/courses/:id/ratings", handlers.InstructorCourseRatings)

	// Instructor Notifications
	instructor.GET("/notifications", handlers.InstructorNotifications)
	instructor.PUT("/notifications/:id/read", handlers.MarkNotificationRead)

	// Instructor Categories (read-only)
	instructor.GET("/categories", handlers.ListCategories)

	// Instructor Coupon Management
	instructor.GET("/coupons", handlers.ListInstructorCoupons)
	instructor.POST("/coupons", handlers.CreateInstructorCoupon)

	// Instructor File Upload
	instructor.POST("/upload", handlers.UploadFile)

	// Student Quiz Routes (protected)
	api.GET("/lessons/:lessonId/quiz", handlers.GetQuizForStudent)
	api.POST("/quizzes/:quizId/start", handlers.StartQuizAttempt)
	api.POST("/attempts/:attemptId/submit", handlers.SubmitQuizAttempt)
	api.GET("/attempts/:attemptId/result", handlers.GetQuizAttemptResult)
	api.GET("/quizzes/:quizId/status", handlers.GetQuizStatus)
	api.GET("/quizzes/:quizId/attempts", handlers.GetUserQuizAttempts)

	// NOTE: Static file serving for /uploads has been removed
	// All content is now served via pre-signed URLs from MinIO
	// Use GET /api/content/:lessonId/url to get secure access URLs

	return e
}

