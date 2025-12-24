// Lesson Progress composable for tracking learning progress
import { ref } from 'vue'

export interface LessonProgress {
    id: string
    user_id: string
    lesson_id: string
    is_completed: boolean
    watch_time: number
    completed_at?: string
    created_at: string
    updated_at: string
}

export interface CourseProgress {
    course_id: string
    total_lessons: number
    completed_lessons: number
    progress_percentage: number
    total_watch_time: number
    lesson_progress: LessonProgress[]
}

export interface Activity {
    id: string
    user_id: string
    activity_type: string
    reference_id?: string
    reference_type?: string
    description: string
    metadata?: Record<string, any>
    created_at: string
}

export interface LearningStats {
    total_watch_time: number
    quizzes_passed: number
    lessons_completed: number
    courses_completed: number
    current_streak: number
}

export const useLessonProgress = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)

    const courseProgress = ref<CourseProgress | null>(null)
    const completedLessonIds = ref<string[]>([])

    // Fetch course progress from server
    const fetchCourseProgress = async (courseId: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<CourseProgress>(`/api/courses/${courseId}/progress`)
            courseProgress.value = response

            // Extract completed lesson IDs
            completedLessonIds.value = response.lesson_progress
                ?.filter(p => p.is_completed)
                ?.map(p => p.lesson_id) || []

            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch progress'
            return null
        } finally {
            loading.value = false
        }
    }

    // Mark a single lesson as complete
    const markLessonComplete = async (lessonId: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<LessonProgress>(`/api/lessons/${lessonId}/progress`, {
                method: 'POST',
                body: JSON.stringify({ is_completed: true }),
            })

            // Update local state
            if (!completedLessonIds.value.includes(lessonId)) {
                completedLessonIds.value.push(lessonId)
            }

            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to mark lesson complete'
            return null
        } finally {
            loading.value = false
        }
    }

    // Update watch time for a lesson
    const updateWatchTime = async (lessonId: string, seconds: number) => {
        if (seconds <= 0) return

        try {
            await api.fetch(`/api/lessons/${lessonId}/watchtime`, {
                method: 'POST',
                body: JSON.stringify({ seconds }),
            })
        } catch (err: any) {
            console.error('Failed to update watch time:', err)
        }
    }

    // Bulk sync lessons from localStorage to server (migration)
    const migrateFromLocalStorage = async (courseId: string) => {
        const storageKey = `course-progress-${courseId}`
        const savedProgress = localStorage.getItem(storageKey)

        if (!savedProgress) return false

        try {
            const lessonIds = JSON.parse(savedProgress) as string[]
            if (!Array.isArray(lessonIds) || lessonIds.length === 0) {
                localStorage.removeItem(storageKey)
                return false
            }

            // Sync to server
            await api.fetch<CourseProgress>(`/api/courses/${courseId}/progress/bulk`, {
                method: 'POST',
                body: JSON.stringify({ lesson_ids: lessonIds }),
            })

            // Remove from localStorage after successful sync
            localStorage.removeItem(storageKey)
            console.log(`Migrated ${lessonIds.length} lesson progress to server`)

            return true
        } catch (err: any) {
            console.error('Failed to migrate progress:', err)
            return false
        }
    }

    // Get progress for a specific lesson
    const getLessonProgress = async (lessonId: string) => {
        try {
            const response = await api.fetch<LessonProgress>(`/api/lessons/${lessonId}/progress`)
            return response
        } catch (err: any) {
            return null
        }
    }

    // Check if a lesson is completed (from local state)
    const isLessonCompleted = (lessonId: string): boolean => {
        return completedLessonIds.value.includes(lessonId)
    }

    return {
        loading,
        error,
        courseProgress,
        completedLessonIds,
        fetchCourseProgress,
        markLessonComplete,
        updateWatchTime,
        migrateFromLocalStorage,
        getLessonProgress,
        isLessonCompleted,
    }
}

// Separate composable for activities and stats
export const useLearningActivities = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)

    const activities = ref<Activity[]>([])
    const stats = ref<LearningStats | null>(null)

    // Fetch recent activities
    const fetchActivities = async (limit: number = 10) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<{ activities: Activity[] }>(`/api/activities?limit=${limit}`)
            activities.value = response.activities || []
            return activities.value
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch activities'
            return []
        } finally {
            loading.value = false
        }
    }

    // Fetch learning stats
    const fetchStats = async () => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<LearningStats>('/api/stats')
            stats.value = response
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch stats'
            return null
        } finally {
            loading.value = false
        }
    }

    // Format watch time to human readable
    const formatWatchTime = (seconds: number): string => {
        if (seconds < 60) return `${seconds} detik`
        const minutes = Math.floor(seconds / 60)
        if (minutes < 60) return `${minutes} menit`
        const hours = Math.floor(minutes / 60)
        const remainingMinutes = minutes % 60
        if (remainingMinutes === 0) return `${hours} jam`
        return `${hours} jam ${remainingMinutes} menit`
    }

    // Format time ago
    const formatTimeAgo = (dateString: string): string => {
        const date = new Date(dateString)
        const now = new Date()
        const diff = now.getTime() - date.getTime()

        const minutes = Math.floor(diff / 60000)
        if (minutes < 1) return 'Baru saja'
        if (minutes < 60) return `${minutes} menit lalu`

        const hours = Math.floor(minutes / 60)
        if (hours < 24) return `${hours} jam lalu`

        const days = Math.floor(hours / 24)
        if (days === 1) return 'Kemarin'
        if (days < 7) return `${days} hari lalu`

        return date.toLocaleDateString('id-ID')
    }

    // Get activity icon based on type
    const getActivityIcon = (type: string): { icon: string; bgClass: string; colorClass: string } => {
        switch (type) {
            case 'lesson_complete':
                return {
                    icon: 'M5 13l4 4L19 7',
                    bgClass: 'bg-accent-100',
                    colorClass: 'text-accent-600'
                }
            case 'quiz_pass':
                return {
                    icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
                    bgClass: 'bg-primary-100',
                    colorClass: 'text-primary-600'
                }
            case 'quiz_fail':
                return {
                    icon: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
                    bgClass: 'bg-warm-100',
                    colorClass: 'text-warm-600'
                }
            case 'course_complete':
                return {
                    icon: 'M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z',
                    bgClass: 'bg-rose-100',
                    colorClass: 'text-rose-600'
                }
            case 'enroll':
                return {
                    icon: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253',
                    bgClass: 'bg-blue-100',
                    colorClass: 'text-blue-600'
                }
            default:
                return {
                    icon: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
                    bgClass: 'bg-neutral-100',
                    colorClass: 'text-neutral-600'
                }
        }
    }

    return {
        loading,
        error,
        activities,
        stats,
        fetchActivities,
        fetchStats,
        formatWatchTime,
        formatTimeAgo,
        getActivityIcon,
    }
}
