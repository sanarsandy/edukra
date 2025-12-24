// Users composable for admin user management
import type { AuthUser } from './useAuth'

export interface UsersResponse {
    users: AuthUser[]
    total: number
    limit: number
    offset: number
}

export const useUsers = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)
    const users = ref<AuthUser[]>([])
    const total = ref(0)

    const page = ref(1)
    const limit = ref(10)

    const offset = computed(() => (page.value - 1) * limit.value)
    const totalPages = computed(() => Math.ceil(total.value / limit.value))

    // Fetch all users (admin)
    const fetchUsers = async (params: { limit?: number; offset?: number; tenant_id?: string } = {}) => {
        loading.value = true
        error.value = null

        // Update local state if params provided
        if (params.limit) limit.value = params.limit

        // Calculate offset based on page if not provided
        const currentOffset = params.offset !== undefined ? params.offset : offset.value

        const queryParams = new URLSearchParams()
        queryParams.append('limit', String(limit.value))
        queryParams.append('offset', String(currentOffset))
        if (params.tenant_id) queryParams.append('tenant_id', params.tenant_id)

        try {
            const response = await api.fetch<UsersResponse>(`/api/admin/users?${queryParams.toString()}`)
            users.value = response.users || []
            total.value = response.total || 0
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch users'
        } finally {
            loading.value = false
        }
    }

    // Get single user (admin)
    const getUser = async (id: string): Promise<AuthUser | null> => {
        try {
            const response = await api.fetch<AuthUser>(`/api/admin/users/${id}`)
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch user'
            return null
        }
    }

    // Create new user (admin)
    const createUser = async (data: { email: string; password: string; full_name: string; role: string }) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<AuthUser>('/api/admin/users', {
                method: 'POST',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to create user'
            return null
        } finally {
            loading.value = false
        }
    }

    // Update user (admin)
    const updateUser = async (id: string, data: Partial<AuthUser>) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<AuthUser>(`/api/admin/users/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update user'
            return null
        } finally {
            loading.value = false
        }
    }

    // Delete user (admin)
    const deleteUser = async (id: string) => {
        loading.value = true
        error.value = null

        try {
            await api.fetch(`/api/admin/users/${id}`, { method: 'DELETE' })
            return true
        } catch (err: any) {
            error.value = err.message || 'Failed to delete user'
            return false
        } finally {
            loading.value = false
        }
    }

    const nextPage = () => {
        if (page.value < totalPages.value) {
            page.value++
            fetchUsers()
        }
    }

    const prevPage = () => {
        if (page.value > 1) {
            page.value--
            fetchUsers()
        }
    }

    const setPage = (p: number) => {
        if (p >= 1 && p <= totalPages.value) {
            page.value = p
            fetchUsers()
        }
    }

    return {
        loading,
        error,
        users,
        total,
        page,
        limit,
        totalPages,
        fetchUsers,
        getUser,
        createUser,
        updateUser,
        deleteUser,
        nextPage,
        prevPage,
        setPage
    }
}
