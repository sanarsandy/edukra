// Composable for Instructor Panel
export interface InstructorStats {
    total_courses: number
    published_courses: number
    draft_courses: number
    pending_review: number
    total_students: number
    average_rating: number
    total_revenue: number
}

export interface InstructorCourse {
    id: string
    title: string
    slug: string
    description: string
    thumbnail_url?: string
    price: number
    status: 'draft' | 'pending_review' | 'approved' | 'rejected' | 'published'
    lessons_count: number
    is_published: boolean
    enrolled_count: number
    category_name?: string
    review_notes?: string
    created_at: string
    updated_at: string
}

export interface Notification {
    id: string
    type: string
    title: string
    message?: string
    reference_id?: string
    reference_type?: string
    is_read: boolean
    created_at: string
}

export const useInstructorPanel = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)

    // ========================================
    // DASHBOARD
    // ========================================
    const fetchDashboard = async () => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{
                stats: InstructorStats
                recent_courses: InstructorCourse[]
            }>('/api/instructor/dashboard')
        } catch (err: any) {
            error.value = err.message || 'Gagal memuat dashboard'
            return null
        } finally {
            loading.value = false
        }
    }

    // ========================================
    // COURSES
    // ========================================
    const fetchCourses = async (params?: { status?: string; limit?: number; offset?: number }) => {
        loading.value = true
        error.value = null
        try {
            const query = new URLSearchParams()
            if (params?.status) query.set('status', params.status)
            if (params?.limit) query.set('limit', String(params.limit))
            if (params?.offset) query.set('offset', String(params.offset))

            return await api.fetch<{
                courses: InstructorCourse[]
                total: number
                limit: number
                offset: number
            }>(`/api/instructor/courses?${query}`)
        } catch (err: any) {
            error.value = err.message || 'Gagal memuat kursus'
            return null
        } finally {
            loading.value = false
        }
    }

    const fetchCourse = async (id: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<InstructorCourse>(`/api/instructor/courses/${id}`)
        } catch (err: any) {
            error.value = err.message || 'Gagal memuat kursus'
            return null
        } finally {
            loading.value = false
        }
    }

    const createCourse = async (data: {
        title: string
        description?: string
        price?: number
        category_id?: string
        thumbnail_url?: string
    }) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ id: string; message: string }>('/api/instructor/courses', {
                method: 'POST',
                body: data
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal membuat kursus'
            return null
        } finally {
            loading.value = false
        }
    }

    const updateCourse = async (id: string, data: Partial<{
        title: string
        description: string
        price: number
        category_id: string
        thumbnail_url: string
        duration: string
    }>) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ message: string }>(`/api/instructor/courses/${id}`, {
                method: 'PUT',
                body: data
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal update kursus'
            return null
        } finally {
            loading.value = false
        }
    }

    const deleteCourse = async (id: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ message: string }>(`/api/instructor/courses/${id}`, {
                method: 'DELETE'
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal hapus kursus'
            return null
        } finally {
            loading.value = false
        }
    }

    const submitForReview = async (id: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ message: string; status: string }>(`/api/instructor/courses/${id}/submit`, {
                method: 'POST'
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal submit kursus'
            return null
        } finally {
            loading.value = false
        }
    }

    // ========================================
    // LESSONS
    // ========================================
    const fetchLessons = async (courseId: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ lessons: any[]; total: number }>(`/api/instructor/courses/${courseId}/lessons`)
        } catch (err: any) {
            error.value = err.message || 'Gagal memuat materi'
            return null
        } finally {
            loading.value = false
        }
    }

    const fetchLessonsTree = async (courseId: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ lessons: any[] }>(`/api/instructor/courses/${courseId}/lessons/tree`)
        } catch (err: any) {
            error.value = err.message || 'Gagal memuat materi'
            return null
        } finally {
            loading.value = false
        }
    }

    const createLesson = async (courseId: string, data: {
        title: string
        description?: string
        content_type?: string
        video_url?: string
        content?: string
        is_container?: boolean
        parent_id?: string
        is_preview?: boolean
    }) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<any>(`/api/instructor/courses/${courseId}/lessons`, {
                method: 'POST',
                body: data
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal membuat materi'
            return null
        } finally {
            loading.value = false
        }
    }

    const updateLesson = async (lessonId: string, data: Partial<{
        title: string
        description: string
        content_type: string
        video_url: string
        content: string
        is_preview: boolean
    }>) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<any>(`/api/instructor/lessons/${lessonId}`, {
                method: 'PUT',
                body: data
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal update materi'
            return null
        } finally {
            loading.value = false
        }
    }

    const deleteLesson = async (lessonId: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ message: string }>(`/api/instructor/lessons/${lessonId}`, {
                method: 'DELETE'
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal hapus materi'
            return null
        } finally {
            loading.value = false
        }
    }

    // ========================================
    // QUIZ
    // ========================================
    const createQuiz = async (lessonId: string, data: {
        title: string
        description?: string
        time_limit?: number
        passing_score?: number
        max_attempts?: number
        show_correct_answers?: boolean
    }) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ id: string; message: string }>(`/api/instructor/lessons/${lessonId}/quiz`, {
                method: 'POST',
                body: data
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal membuat kuis'
            return null
        } finally {
            loading.value = false
        }
    }

    const getQuizByLesson = async (lessonId: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{
                quiz: any
                questions: any[]
            }>(`/api/instructor/lessons/${lessonId}/quiz`)
        } catch (err: any) {
            error.value = err.message || 'Gagal memuat kuis'
            return null
        } finally {
            loading.value = false
        }
    }

    const updateQuiz = async (quizId: string, data: {
        title?: string
        description?: string
        time_limit?: number
        passing_score?: number
        max_attempts?: number
        show_correct_answers?: boolean
    }) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ message: string }>(`/api/instructor/quizzes/${quizId}`, {
                method: 'PUT',
                body: data
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal update kuis'
            return null
        } finally {
            loading.value = false
        }
    }

    const deleteQuiz = async (quizId: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ message: string }>(`/api/instructor/quizzes/${quizId}`, {
                method: 'DELETE'
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal hapus kuis'
            return null
        } finally {
            loading.value = false
        }
    }

    const addQuestion = async (quizId: string, data: {
        question_type: string
        question_text: string
        explanation?: string
        points?: number
        required?: boolean
        options?: { option_text: string; is_correct: boolean }[]
    }) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ id: string; message: string }>(`/api/instructor/quizzes/${quizId}/questions`, {
                method: 'POST',
                body: data
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal menambah soal'
            return null
        } finally {
            loading.value = false
        }
    }

    const updateQuestion = async (questionId: string, data: {
        question_type?: string
        question_text?: string
        explanation?: string
        points?: number
        required?: boolean
        options?: { id?: string; option_text: string; is_correct: boolean }[]
    }) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ message: string }>(`/api/instructor/questions/${questionId}`, {
                method: 'PUT',
                body: data
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal update soal'
            return null
        } finally {
            loading.value = false
        }
    }

    const deleteQuestion = async (questionId: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ message: string }>(`/api/instructor/questions/${questionId}`, {
                method: 'DELETE'
            })
        } catch (err: any) {
            error.value = err.message || 'Gagal hapus soal'
            return null
        } finally {
            loading.value = false
        }
    }

    // ========================================
    // ANALYTICS
    // ========================================
    const fetchCourseStudents = async (courseId: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{ students: any[]; total: number }>(`/api/instructor/courses/${courseId}/students`)
        } catch (err: any) {
            error.value = err.message || 'Gagal memuat data siswa'
            return null
        } finally {
            loading.value = false
        }
    }

    const fetchCourseRatings = async (courseId: string) => {
        loading.value = true
        error.value = null
        try {
            return await api.fetch<{
                average_rating: number
                total_ratings: number
                ratings: any[]
            }>(`/api/instructor/courses/${courseId}/ratings`)
        } catch (err: any) {
            error.value = err.message || 'Gagal memuat rating'
            return null
        } finally {
            loading.value = false
        }
    }

    // ========================================
    // NOTIFICATIONS
    // ========================================
    const fetchNotifications = async () => {
        try {
            return await api.fetch<{
                notifications: Notification[]
                unread_count: number
            }>('/api/instructor/notifications')
        } catch {
            return null
        }
    }

    const markNotificationRead = async (id: string) => {
        try {
            return await api.fetch<{ message: string }>(`/api/instructor/notifications/${id}/read`, {
                method: 'PUT'
            })
        } catch {
            return null
        }
    }

    // ========================================
    // CATEGORIES (read-only)
    // ========================================
    const fetchCategories = async () => {
        try {
            return await api.fetch<{ categories: any[] }>('/api/instructor/categories')
        } catch {
            return null
        }
    }

    // ========================================
    // FILE UPLOAD
    // ========================================
    const uploadFile = async (file: File, type: 'thumbnail' | 'video' | 'document' = 'thumbnail') => {
        loading.value = true
        error.value = null
        try {
            const formData = new FormData()
            formData.append('file', file)
            formData.append('type', type)

            const token = useCookie('token')
            const config = useRuntimeConfig()
            const baseURL = config.public.apiBase || 'http://localhost:8080'

            const response = await fetch(`${baseURL}/api/instructor/upload`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token.value}`
                },
                body: formData
            })

            if (!response.ok) {
                const data = await response.json()
                throw new Error(data.error || 'Upload gagal')
            }

            return await response.json()
        } catch (err: any) {
            error.value = err.message || 'Gagal upload file'
            return null
        } finally {
            loading.value = false
        }
    }

    // ========================================
    // AUTH
    // ========================================
    const instructorLogin = async (email: string, password: string) => {
        loading.value = true
        error.value = null
        try {
            const config = useRuntimeConfig()
            const baseURL = config.public.apiBase || 'http://localhost:8080'

            const response = await fetch(`${baseURL}/api/auth/instructor/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            })

            const data = await response.json()

            if (!response.ok) {
                error.value = data.error || 'Login gagal'
                return null
            }

            // Save token and user
            const token = useCookie('token')
            const userCookie = useCookie('user')
            token.value = data.token
            userCookie.value = data.user

            return data
        } catch (err: any) {
            error.value = 'Terjadi kesalahan jaringan'
            return null
        } finally {
            loading.value = false
        }
    }

    const instructorLogout = () => {
        const token = useCookie('token')
        const userCookie = useCookie('user')
        token.value = null
        userCookie.value = null
        navigateTo('/instructor/login')
    }

    return {
        loading,
        error,
        // Dashboard
        fetchDashboard,
        // Courses
        fetchCourses,
        fetchCourse,
        createCourse,
        updateCourse,
        deleteCourse,
        submitForReview,
        // Lessons
        fetchLessons,
        fetchLessonsTree,
        createLesson,
        updateLesson,
        deleteLesson,
        // Quiz
        createQuiz,
        getQuizByLesson,
        updateQuiz,
        deleteQuiz,
        addQuestion,
        updateQuestion,
        deleteQuestion,
        // Analytics
        fetchCourseStudents,
        fetchCourseRatings,
        // Notifications
        fetchNotifications,
        markNotificationRead,
        // Categories
        fetchCategories,
        // Upload
        uploadFile,
        // Auth
        instructorLogin,
        instructorLogout
    }
}


