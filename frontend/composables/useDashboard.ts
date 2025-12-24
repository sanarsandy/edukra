// Dashboard composable for dashboard data
import type { AuthUser } from './useAuth'
import type { Enrollment } from './useEnrollments'

export interface DashboardStats {
    enrolled_courses: number
    completed_courses: number
    in_progress: number
}

export interface AdminDashboardStats {
    total_users: number
    total_courses: number
    total_transactions: number
    total_revenue: number
}

export interface Transaction {
    id: string
    tenant_id: string
    user_id: string
    course_id?: string
    amount: number
    currency: string
    status: 'pending' | 'success' | 'failed' | 'refunded'
    payment_gateway: string
    payment_method?: string
    created_at: string
}

export interface UserDashboardResponse {
    stats: DashboardStats
    recent_courses: Enrollment[]
}

export interface AdminDashboardResponse {
    stats: AdminDashboardStats
    recent_users: AuthUser[]
    recent_transactions: Transaction[]
}

export interface TransactionsResponse {
    transactions: Transaction[]
    total: number
    limit: number
    offset: number
}

export const useDashboard = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)

    // User dashboard
    const stats = ref<DashboardStats | null>(null)
    const recentCourses = ref<Enrollment[]>([])

    // Admin dashboard
    const adminStats = ref<AdminDashboardStats | null>(null)
    const recentUsers = ref<AuthUser[]>([])
    const transactions = ref<Transaction[]>([])
    const transactionsTotal = ref(0)

    // Fetch user dashboard
    const fetchUserDashboard = async () => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<UserDashboardResponse>('/api/dashboard')
            stats.value = response.stats
            recentCourses.value = response.recent_courses || []
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch dashboard'
        } finally {
            loading.value = false
        }
    }

    // Fetch admin dashboard
    const fetchAdminDashboard = async (tenantId?: string) => {
        loading.value = true
        error.value = null

        const params = tenantId ? `?tenant_id=${tenantId}` : ''

        try {
            const response = await api.fetch<AdminDashboardResponse>(`/api/admin/dashboard${params}`)
            adminStats.value = response.stats
            recentUsers.value = response.recent_users || []
            transactions.value = response.recent_transactions || []
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch admin dashboard'
        } finally {
            loading.value = false
        }
    }

    // Fetch transactions (admin)
    const fetchTransactions = async (params: { limit?: number; offset?: number; status?: string } = {}) => {
        loading.value = true
        error.value = null

        const queryParams = new URLSearchParams()
        if (params.limit) queryParams.append('limit', String(params.limit))
        if (params.offset) queryParams.append('offset', String(params.offset))
        if (params.status) queryParams.append('status', params.status)

        try {
            const response = await api.fetch<TransactionsResponse>(`/api/admin/transactions?${queryParams.toString()}`)
            transactions.value = response.transactions || []
            transactionsTotal.value = response.total || 0
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch transactions'
        } finally {
            loading.value = false
        }
    }

    // Update transaction status (admin)
    const updateTransactionStatus = async (id: string, status: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Transaction>(`/api/admin/transactions/${id}/status`, {
                method: 'PUT',
                body: JSON.stringify({ status }),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update transaction'
            return null
        } finally {
            loading.value = false
        }
    }

    return {
        loading,
        error,
        stats,
        recentCourses,
        adminStats,
        recentUsers,
        transactions,
        transactionsTotal,
        fetchUserDashboard,
        fetchAdminDashboard,
        fetchTransactions,
        updateTransactionStatus,
    }
}
