// Rating composable for course ratings and reviews
import { ref } from 'vue'

export interface CourseRating {
    id: string
    user_id: string
    course_id: string
    rating: number
    review?: string
    created_at: string
    updated_at: string
    user_name: string
}

export interface CourseRatingStats {
    average_rating: number
    total_ratings: number
}

export const useRatings = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)

    const ratings = ref<CourseRating[]>([])
    const stats = ref<CourseRatingStats | null>(null)
    const myRating = ref<CourseRating | null>(null)
    const hasRated = ref(false)

    // Fetch ratings for a course
    const fetchRatings = async (courseId: string, limit = 10) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<{ ratings: CourseRating[] }>(`/api/courses/${courseId}/ratings?limit=${limit}`)
            ratings.value = response.ratings || []
            return ratings.value
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch ratings'
            return []
        } finally {
            loading.value = false
        }
    }

    // Fetch rating stats for a course
    const fetchStats = async (courseId: string) => {
        try {
            const response = await api.fetch<CourseRatingStats>(`/api/courses/${courseId}/ratings/stats`)
            stats.value = response
            return response
        } catch (err: any) {
            console.error('Failed to fetch rating stats:', err)
            return null
        }
    }

    // Fetch user's own rating
    const fetchMyRating = async (courseId: string) => {
        try {
            const response = await api.fetch<{ has_rated: boolean; rating: CourseRating | null }>(`/api/courses/${courseId}/my-rating`)
            hasRated.value = response.has_rated
            myRating.value = response.rating
            return response
        } catch (err: any) {
            console.error('Failed to fetch my rating:', err)
            return null
        }
    }

    // Submit a new rating
    const submitRating = async (courseId: string, rating: number, review?: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<CourseRating>(`/api/courses/${courseId}/ratings`, {
                method: 'POST',
                body: JSON.stringify({ rating, review })
            })
            myRating.value = response
            hasRated.value = true
            // Refresh ratings and stats
            await Promise.all([
                fetchRatings(courseId),
                fetchStats(courseId)
            ])
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to submit rating'
            throw err
        } finally {
            loading.value = false
        }
    }

    // Update existing rating
    const updateRating = async (courseId: string, rating: number, review?: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<CourseRating>(`/api/courses/${courseId}/ratings`, {
                method: 'PUT',
                body: JSON.stringify({ rating, review })
            })
            myRating.value = response
            // Refresh ratings and stats
            await Promise.all([
                fetchRatings(courseId),
                fetchStats(courseId)
            ])
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update rating'
            throw err
        } finally {
            loading.value = false
        }
    }

    // Delete rating
    const deleteRating = async (courseId: string) => {
        loading.value = true
        error.value = null

        try {
            await api.fetch(`/api/courses/${courseId}/ratings`, {
                method: 'DELETE'
            })
            myRating.value = null
            hasRated.value = false
            // Refresh ratings and stats
            await Promise.all([
                fetchRatings(courseId),
                fetchStats(courseId)
            ])
        } catch (err: any) {
            error.value = err.message || 'Failed to delete rating'
            throw err
        } finally {
            loading.value = false
        }
    }

    // Format time ago
    const formatTimeAgo = (dateString: string): string => {
        const date = new Date(dateString)
        const now = new Date()
        const diffMs = now.getTime() - date.getTime()
        const diffMins = Math.floor(diffMs / 60000)
        const diffHours = Math.floor(diffMs / 3600000)
        const diffDays = Math.floor(diffMs / 86400000)

        if (diffMins < 60) return `${diffMins} menit lalu`
        if (diffHours < 24) return `${diffHours} jam lalu`
        if (diffDays < 7) return `${diffDays} hari lalu`
        return date.toLocaleDateString('id-ID')
    }

    return {
        loading,
        error,
        ratings,
        stats,
        myRating,
        hasRated,
        fetchRatings,
        fetchStats,
        fetchMyRating,
        submitRating,
        updateRating,
        deleteRating,
        formatTimeAgo,
    }
}
