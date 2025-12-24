// Courses composable for API operations
export interface Course {
    id: string
    tenant_id: string
    instructor_id?: string
    category_id?: string
    title: string
    slug: string
    description: string
    thumbnail_url?: string
    price: number
    currency: string
    is_published: boolean
    is_featured: boolean
    created_at: string
    updated_at: string
    lessons?: Lesson[]
    instructor?: User
    category?: Category
}

export interface Category {
    id: string
    name: string
    slug: string
}

export interface Lesson {
    id: string
    course_id: string
    title: string
    description: string
    order_index: number
    content_type: 'video' | 'pdf' | 'quiz' | 'text'
    video_url?: string
    video_duration: number
    security_level: string
    is_preview: boolean
    created_at: string
    updated_at: string
}

export interface User {
    id: string
    email: string
    full_name: string
    role: string
    avatar_url?: string
}

export interface CoursesResponse {
    courses: Course[]
    total: number
    limit: number
    offset: number
}

export const useCourses = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)
    const courses = ref<Course[]>([])
    const course = ref<Course | null>(null)
    const total = ref(0)

    // Fetch all published courses (Public)
    const fetchCourses = async (params: { limit?: number; offset?: number; tenant_id?: string } = {}) => {
        loading.value = true
        error.value = null

        const queryParams = new URLSearchParams()
        if (params.limit) queryParams.append('limit', String(params.limit))
        if (params.offset) queryParams.append('offset', String(params.offset))
        if (params.tenant_id) queryParams.append('tenant_id', params.tenant_id)

        try {
            const response = await api.fetch<CoursesResponse>(`/api/courses?${queryParams.toString()}`)
            courses.value = response.courses || []
            total.value = response.total || 0
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch courses'
        } finally {
            loading.value = false
        }
    }

    // Fetch all courses for admin (Admin)
    const fetchAdminCourses = async (params: { limit?: number; offset?: number; tenant_id?: string } = {}) => {
        loading.value = true
        error.value = null

        const queryParams = new URLSearchParams()
        if (params.limit) queryParams.append('limit', String(params.limit))
        if (params.offset) queryParams.append('offset', String(params.offset))
        if (params.tenant_id) queryParams.append('tenant_id', params.tenant_id)

        try {
            const response = await api.fetch<CoursesResponse>(`/api/admin/courses?${queryParams.toString()}`)
            courses.value = response.courses || []
            total.value = response.total || 0
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch courses'
        } finally {
            loading.value = false
        }
    }

    // Fetch single course with lessons
    const fetchCourse = async (id: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Course>(`/api/courses/${id}`)
            course.value = response
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch course'
        } finally {
            loading.value = false
        }
    }

    // Admin: Create course
    const createCourse = async (data: Partial<Course>) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Course>('/api/admin/courses', {
                method: 'POST',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to create course'
            return null
        } finally {
            loading.value = false
        }
    }

    // Admin: Update course
    const updateCourse = async (id: string, data: Partial<Course>) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Course>(`/api/admin/courses/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update course'
            return null
        } finally {
            loading.value = false
        }
    }

    // Admin: Delete course
    const deleteCourse = async (id: string) => {
        loading.value = true
        error.value = null

        try {
            await api.fetch(`/api/admin/courses/${id}`, { method: 'DELETE' })
            return true
        } catch (err: any) {
            error.value = err.message || 'Failed to delete course'
            return false
        } finally {
            loading.value = false
        }
    }

    // Admin: Toggle publish status
    const togglePublish = async (id: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Course>(`/api/admin/courses/${id}/publish`, {
                method: 'PUT',
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to toggle publish'
            return null
        } finally {
            loading.value = false
        }
    }

    return {
        loading,
        error,
        courses,
        course,
        total,
        fetchCourses,
        fetchAdminCourses,
        fetchCourse,
        createCourse,
        updateCourse,
        deleteCourse,
        togglePublish,
    }
}
