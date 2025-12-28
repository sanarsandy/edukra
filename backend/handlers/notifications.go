package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	customMiddleware "github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
)

// NotificationResponse represents a notification for API response
type NotificationResponse struct {
	ID            string     `json:"id"`
	Type          string     `json:"type"`
	Title         string     `json:"title"`
	Message       *string    `json:"message,omitempty"`
	ReferenceID   *string    `json:"reference_id,omitempty"`
	ReferenceType *string    `json:"reference_type,omitempty"`
	IsRead        bool       `json:"is_read"`
	CreatedAt     time.Time  `json:"created_at"`
	Time          string     `json:"time"` // Human-readable relative time
}

// GetMyNotifications returns notifications for the authenticated user
func GetMyNotifications(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	rows, err := db.DB.Query(`
		SELECT id, type, title, message, reference_id, reference_type, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 20
	`, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil notifikasi"})
	}
	defer rows.Close()

	var notifications []NotificationResponse
	for rows.Next() {
		var n NotificationResponse
		var message, refID, refType sql.NullString

		err := rows.Scan(&n.ID, &n.Type, &n.Title, &message, &refID, &refType, &n.IsRead, &n.CreatedAt)
		if err != nil {
			continue
		}

		if message.Valid {
			n.Message = &message.String
		}
		if refID.Valid {
			n.ReferenceID = &refID.String
		}
		if refType.Valid {
			n.ReferenceType = &refType.String
		}
		
		// Format relative time
		n.Time = formatNotificationTime(n.CreatedAt)

		notifications = append(notifications, n)
	}

	// Get unread count
	var unreadCount int
	db.DB.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`, userID).Scan(&unreadCount)

	// Return empty array instead of null
	if notifications == nil {
		notifications = []NotificationResponse{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

// MarkNotificationAsRead marks a single notification as read
func MarkNotificationAsRead(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	notifID := c.Param("id")

	result, err := db.DB.Exec(`UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`, notifID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal update notifikasi"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Notifikasi tidak ditemukan"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Notifikasi ditandai sudah dibaca"})
}

// MarkAllNotificationsRead marks all notifications as read for a user
func MarkAllNotificationsRead(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec(`UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false`, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal update notifikasi"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Semua notifikasi ditandai sudah dibaca"})
}

// formatNotificationTime formats a time as relative string
func formatNotificationTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "Baru saja"
	} else if diff < time.Hour {
		mins := int(diff.Minutes())
		return fmt.Sprintf("%d menit lalu", mins)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 jam lalu"
		}
		return fmt.Sprintf("%d jam lalu", hours)
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "Kemarin"
		}
		return fmt.Sprintf("%d hari lalu", days)
	}

	return t.Format("2 Jan 2006")
}
