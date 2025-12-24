// Instructors composable for admin instructor management
export interface Instructor {
    id: string
    email: string
    full_name: string
    specialty?: string
    bio?: string
    avatar_url?: string
    course_count?: number
}

export interface InstructorsResponse {
    instructors: Instructor[]
    total: number
}

export const useInstructors = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)
    const instructors = ref<Instructor[]>([])
    const total = ref(0)

    // Fetch all instructors
    const fetchInstructors = async () => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<InstructorsResponse>('/api/admin/instructors')
            instructors.value = response.instructors || []
            total.value = response.total || 0
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch instructors'
        } finally {
            loading.value = false
        }
    }

    // Create instructor
    const createInstructor = async (data: Partial<Instructor>) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Instructor>('/api/admin/instructors', {
                method: 'POST',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to create instructor'
            return null
        } finally {
            loading.value = false
        }
    }

    // Update instructor
    const updateInstructor = async (id: string, data: Partial<Instructor>) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Instructor>(`/api/admin/instructors/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update instructor'
            return null
        } finally {
            loading.value = false
        }
    }

    // Delete instructor
    const deleteInstructor = async (id: string) => {
        loading.value = true
        error.value = null

        try {
            await api.fetch(`/api/admin/instructors/${id}`, { method: 'DELETE' })
            return true
        } catch (err: any) {
            error.value = err.message || 'Failed to delete instructor'
            return false
        } finally {
            loading.value = false
        }
    }

    return {
        loading,
        error,
        instructors,
        total,
        fetchInstructors,
        createInstructor,
        updateInstructor,
        deleteInstructor,
    }
}
