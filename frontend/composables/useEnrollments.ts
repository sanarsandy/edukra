// Enrollments composable for user course enrollments
import type { Course } from './useCourses'

export interface Enrollment {
    id: string
    user_id: string
    course_id: string
    transaction_id?: string
    progress_percentage: number
    completed_at?: string
    enrolled_at: string
    course?: Course
}

export interface EnrollmentsResponse {
    enrollments: Enrollment[]
    total: number
    limit: number
    offset: number
}

export interface EnrollmentCheck {
    enrolled: boolean
    enrollment?: Enrollment
}

export const useEnrollments = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)
    const enrollments = ref<Enrollment[]>([])
    const total = ref(0)

    // Fetch user's enrollments
    const fetchEnrollments = async (params: { limit?: number; offset?: number } = {}) => {
        loading.value = true
        error.value = null

        const queryParams = new URLSearchParams()
        if (params.limit) queryParams.append('limit', String(params.limit))
        if (params.offset) queryParams.append('offset', String(params.offset))

        try {
            const response = await api.fetch<EnrollmentsResponse>(`/api/enrollments?${queryParams.toString()}`)
            enrollments.value = response.enrollments || []
            total.value = response.total || 0
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch enrollments'
        } finally {
            loading.value = false
        }
    }

    // Enroll in a course
    const enrollInCourse = async (courseId: string, transactionId?: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Enrollment>('/api/enrollments', {
                method: 'POST',
                body: JSON.stringify({
                    course_id: courseId,
                    transaction_id: transactionId
                }),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to enroll'
            return null
        } finally {
            loading.value = false
        }
    }

    // Check if enrolled in a course
    const checkEnrollment = async (courseId: string): Promise<EnrollmentCheck> => {
        try {
            const response = await api.fetch<EnrollmentCheck>(`/api/courses/${courseId}/enrollment`)
            return response
        } catch (err: any) {
            return { enrolled: false }
        }
    }

    // Update progress
    const updateProgress = async (enrollmentId: string, progress: number) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Enrollment>(`/api/enrollments/${enrollmentId}/progress`, {
                method: 'PUT',
                body: JSON.stringify({ progress }),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update progress'
            return null
        } finally {
            loading.value = false
        }
    }

    // Get enrollment by ID
    const getEnrollment = async (id: string): Promise<Enrollment | null> => {
        try {
            const response = await api.fetch<Enrollment>(`/api/enrollments/${id}`)
            return response
        } catch (err: any) {
            return null
        }
    }

    return {
        loading,
        error,
        enrollments,
        total,
        fetchEnrollments,
        enrollInCourse,
        checkEnrollment,
        updateProgress,
        getEnrollment,
    }
}
