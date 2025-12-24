// Lessons composable for admin course materials management
export interface Lesson {
    id: string
    course_id: string
    parent_id?: string | null
    is_container: boolean
    title: string
    description: string
    order_index: number
    content_type: 'video' | 'pdf' | 'quiz' | 'text' | ''
    video_url?: string
    content?: string // Rich text content for text type
    video_duration: number
    security_level: 'public' | 'signed_url' | 'aes_128' | 'full_drm'
    is_preview: boolean
    created_at: string
    updated_at: string
    children?: Lesson[] // Nested children for tree structure
}

export interface LessonsResponse {
    lessons: Lesson[]
    total?: number
}

export interface CreateLessonRequest {
    parent_id?: string | null
    is_container?: boolean
    title: string
    description?: string
    content_type?: string
    video_url?: string
    content?: string
    security_level?: string
    is_preview?: boolean
}

export interface UpdateLessonRequest {
    parent_id?: string | null
    is_container?: boolean
    title?: string
    description?: string
    content_type?: string
    video_url?: string
    content?: string
    security_level?: string
    is_preview?: boolean
    order_index?: number
}

export interface MoveLessonRequest {
    parent_id?: string | null
    order_index: number
}

export const useLessons = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)
    const lessons = ref<Lesson[]>([])
    const lessonsTree = ref<Lesson[]>([])
    const lesson = ref<Lesson | null>(null)
    const total = ref(0)

    // Fetch all lessons for a course (flat list)
    const fetchLessons = async (courseId: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<LessonsResponse>(`/api/courses/${courseId}/lessons`)
            lessons.value = response.lessons || []
            total.value = response.total || lessons.value.length
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch lessons'
            lessons.value = []
        } finally {
            loading.value = false
        }
    }

    // Fetch lessons as a tree structure
    const fetchLessonsTree = async (courseId: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<LessonsResponse>(`/api/admin/courses/${courseId}/lessons/tree`)
            lessonsTree.value = response.lessons || []
            total.value = response.total || 0
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch lessons tree'
            lessonsTree.value = []
        } finally {
            loading.value = false
        }
    }

    // Fetch single lesson by ID
    const fetchLesson = async (lessonId: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Lesson>(`/api/lessons/${lessonId}`)
            lesson.value = response
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch lesson'
            lesson.value = null
            return null
        } finally {
            loading.value = false
        }
    }

    // Create a new lesson
    const createLesson = async (courseId: string, data: CreateLessonRequest) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Lesson>(`/api/admin/courses/${courseId}/lessons`, {
                method: 'POST',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to create lesson'
            return null
        } finally {
            loading.value = false
        }
    }

    // Update a lesson
    const updateLesson = async (lessonId: string, data: UpdateLessonRequest) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Lesson>(`/api/admin/lessons/${lessonId}`, {
                method: 'PUT',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update lesson'
            return null
        } finally {
            loading.value = false
        }
    }

    // Move a lesson to a new parent
    const moveLesson = async (lessonId: string, data: MoveLessonRequest) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Lesson>(`/api/admin/lessons/${lessonId}/move`, {
                method: 'PUT',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to move lesson'
            return null
        } finally {
            loading.value = false
        }
    }

    // Delete a lesson
    const deleteLesson = async (lessonId: string) => {
        loading.value = true
        error.value = null

        try {
            await api.fetch(`/api/admin/lessons/${lessonId}`, { method: 'DELETE' })
            return true
        } catch (err: any) {
            error.value = err.message || 'Failed to delete lesson'
            return false
        } finally {
            loading.value = false
        }
    }

    // Reorder lessons
    const reorderLessons = async (courseId: string, lessonIds: string[]) => {
        loading.value = true
        error.value = null

        try {
            await api.fetch(`/api/admin/courses/${courseId}/lessons/reorder`, {
                method: 'PUT',
                body: JSON.stringify({ lesson_ids: lessonIds }),
            })
            return true
        } catch (err: any) {
            error.value = err.message || 'Failed to reorder lessons'
            return false
        } finally {
            loading.value = false
        }
    }

    return {
        loading,
        error,
        lessons,
        lessonsTree,
        lesson,
        total,
        fetchLessons,
        fetchLessonsTree,
        fetchLesson,
        createLesson,
        updateLesson,
        moveLesson,
        deleteLesson,
        reorderLessons,
    }
}
