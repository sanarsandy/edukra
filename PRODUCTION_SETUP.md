# LMS v1.1 - Production Setup Guide

## Prerequisites

- Docker & Docker Compose
- PostgreSQL 15+
- Node.js 22+ (for local development)
- Go 1.21+ (for backend development)

## Quick Start (Docker)

```bash
# 1. Clone repository
git clone <repository-url>
cd "LMS v1.1"

# 2. Copy environment file
cp .env.example .env

# 3. Configure environment variables in .env
# Edit .env with your production values

# 4. Start services
docker compose -f docker-compose.prod.yml up -d

# 5. Initialize database (first time only)
docker exec lms-api sh -c "cat /app/migrations/999_production_ready.sql | psql -h db -U lms_user -d lms_db"
docker exec lms-api sh -c "cat /app/seeds/admin.sql | psql -h db -U lms_user -d lms_db"
```

## Environment Variables

Create a `.env` file with the following variables:

```env
# Database
DB_HOST=db
DB_PORT=5432
DB_USER=lms_user
DB_PASSWORD=<strong-password>
DB_NAME=lms_db

# JWT
JWT_SECRET=<your-jwt-secret-min-32-chars>

# Google OAuth (optional)
GOOGLE_CLIENT_ID=<your-google-client-id>
GOOGLE_CLIENT_SECRET=<your-google-client-secret>
GOOGLE_REDIRECT_URL=https://yourdomain.com/api/auth/google/callback

# Frontend
FRONTEND_URL=https://yourdomain.com
NUXT_PUBLIC_API_BASE=https://api.yourdomain.com

# CORS
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

## Default Credentials

| Role  | Email             | Password   |
|-------|-------------------|------------|
| Admin | admin@lms.local   | Admin123!  |

⚠️ **IMPORTANT**: Change the admin password immediately after first login!

## Database Schema

The production migration (`migrations/999_production_ready.sql`) creates:

### Core Tables
- `tenants` - Multi-tenant support
- `users` - User accounts (admin, instructor, student)
- `categories` - Course categories
- `courses` - Course content
- `lessons` - Course lessons (supports nested structure)
- `enrollments` - User course enrollments
- `transactions` - Payment transactions

### Learning Features
- `lesson_progress` - Track lesson completion
- `activity_logs` - User activity tracking
- `learning_streaks` - Daily learning streaks
- `certificates` - Course completion certificates
- `course_ratings` - Course reviews and ratings

### Quiz System
- `quizzes` - Quiz definitions
- `quiz_questions` - Quiz questions
- `quiz_options` - Multiple choice options
- `quiz_attempts` - User quiz attempts
- `quiz_answers` - User answers

## API Endpoints

### Public Endpoints (No Auth)
- `GET /health` - Health check
- `GET /api/courses` - List published courses
- `GET /api/courses/:id` - Get course details
- `GET /api/categories` - List categories
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `POST /api/auth/admin/login` - Admin login

### Protected Endpoints (Require JWT)
- `GET /api/me` - Get current user
- `PUT /api/me` - Update profile
- `GET /api/dashboard` - User dashboard
- `GET /api/enrollments` - My enrollments
- `POST /api/enrollments` - Enroll in course
- `GET /api/certificates` - My certificates
- `GET /api/lessons/:id` - Get lesson content
- `GET/POST /api/lessons/:id/progress` - Lesson progress

### Admin Endpoints (Require Admin Role)
- `GET /api/admin/dashboard` - Admin dashboard
- `GET/POST/PUT/DELETE /api/admin/users` - User management
- `GET/POST/PUT/DELETE /api/admin/courses` - Course management
- `GET/POST/PUT/DELETE /api/admin/lessons` - Lesson management
- `GET/POST/PUT/DELETE /api/admin/categories` - Category management
- `GET /api/admin/transactions` - Transaction management
- `GET/PUT /api/admin/settings` - Site settings

## Ports

| Service  | Port | Description         |
|----------|------|---------------------|
| Frontend | 3000 | Nuxt.js application |
| API      | 8080 | Go Echo backend     |
| Database | 5432 | PostgreSQL          |

## Troubleshooting

### API returns 500 errors
1. Check database connection: `docker logs lms-api`
2. Verify all migrations ran: `docker exec lms-db psql -U lms_user -d lms_db -c "\dt"`
3. Run production migration again if needed

### "missing or malformed jwt" error
- This is expected for protected endpoints without valid token
- Login first to get a valid JWT token

### Database connection refused
- Ensure PostgreSQL container is running: `docker ps`
- Wait a few seconds for database to initialize
- Restart API container: `docker restart lms-api`

### Frontend can't connect to API
1. Check CORS settings in `.env`
2. Verify `NUXT_PUBLIC_API_BASE` is correct
3. Ensure API container is healthy: `curl http://localhost:8080/health`

## Security Checklist

- [ ] Change default admin password
- [ ] Set strong `JWT_SECRET` (min 32 characters)
- [ ] Set strong `DB_PASSWORD`
- [ ] Configure proper `CORS_ALLOWED_ORIGINS`
- [ ] Enable HTTPS in production
- [ ] Configure firewall rules
- [ ] Set up database backups
- [ ] Configure rate limiting (nginx)

## Backup & Restore

### Backup Database
```bash
docker exec lms-db pg_dump -U lms_user lms_db > backup_$(date +%Y%m%d).sql
```

### Restore Database
```bash
cat backup_YYYYMMDD.sql | docker exec -i lms-db psql -U lms_user -d lms_db
```

## Support

For issues or questions, please open a GitHub issue.


