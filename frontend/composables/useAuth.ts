// Auth composable for authentication
export interface AuthUser {
    id: string
    email: string
    full_name: string
    role: 'admin' | 'instructor' | 'student'
    avatar_url?: string
    auth_provider?: string
    created_at?: string
}

export interface LoginRequest {
    email: string
    password: string
}

export interface RegisterRequest {
    email: string
    password: string
    full_name: string
}

export interface AuthResponse {
    token: string
    user: AuthUser
}

export const useAuth = () => {
    const config = useRuntimeConfig()
    const baseURL = config.public.apiBase || 'http://localhost:8080'

    const loading = ref(false)
    const error = ref<string | null>(null)
    const token = useCookie('token')
    const userCookie = useCookie<AuthUser | null>('user')
    const user = ref<AuthUser | null>(null)

    // Initialize user from cookie
    const initAuth = () => {
        if (userCookie.value) {
            user.value = typeof userCookie.value === 'string'
                ? JSON.parse(userCookie.value)
                : userCookie.value
        }
    }

    // Check if authenticated
    const isAuthenticated = computed(() => !!token.value && !!user.value)

    // Check if admin
    const isAdmin = computed(() => user.value?.role === 'admin')

    // Check if instructor
    const isInstructor = computed(() => user.value?.role === 'instructor' || user.value?.role === 'admin')

    // Get permissions from JWT
    const permissions = computed(() => {
        if (!token.value) return []
        try {
            const payload = JSON.parse(atob(token.value.split('.')[1]))
            return payload.permissions || []
        } catch {
            return []
        }
    })

    // Check if user has specific permission
    const hasPermission = (permission: string): boolean => {
        return permissions.value.includes(permission)
    }

    // Login
    const login = async (credentials: LoginRequest) => {
        loading.value = true
        error.value = null

        try {
            const response = await fetch(`${baseURL}/api/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(credentials),
            })

            const data = await response.json()

            if (!response.ok) {
                error.value = data.error || 'Login gagal'
                return null
            }

            token.value = data.token
            user.value = data.user
            userCookie.value = data.user

            return data as AuthResponse
        } catch (err: any) {
            error.value = 'Terjadi kesalahan jaringan'
            return null
        } finally {
            loading.value = false
        }
    }

    // Admin Login
    const adminLogin = async (credentials: LoginRequest) => {
        loading.value = true
        error.value = null

        try {
            const response = await fetch(`${baseURL}/api/auth/admin/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(credentials),
            })

            const data = await response.json()

            if (!response.ok) {
                error.value = data.error || 'Login gagal'
                return null
            }

            token.value = data.token
            user.value = data.user
            userCookie.value = data.user

            return data as AuthResponse
        } catch (err: any) {
            error.value = 'Terjadi kesalahan jaringan'
            return null
        } finally {
            loading.value = false
        }
    }

    // Register
    const register = async (data: RegisterRequest) => {
        loading.value = true
        error.value = null

        try {
            const response = await fetch(`${baseURL}/api/auth/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data),
            })

            const result = await response.json()

            if (!response.ok) {
                error.value = result.error || 'Registrasi gagal'
                return null
            }

            // Auto login after register if token provided
            if (result.token) {
                token.value = result.token
                user.value = result.user
                userCookie.value = result.user
            }

            return result
        } catch (err: any) {
            error.value = 'Terjadi kesalahan jaringan'
            return null
        } finally {
            loading.value = false
        }
    }

    // Logout
    const logout = () => {
        token.value = null
        user.value = null
        userCookie.value = null

        // Optionally call backend logout
        fetch(`${baseURL}/api/auth/logout`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token.value}`
            },
        }).catch(() => { })

        navigateTo('/login')
    }

    // Admin Logout
    const adminLogout = () => {
        token.value = null
        user.value = null
        userCookie.value = null
        navigateTo('/admin/login')
    }

    // Refresh token
    const refreshToken = async () => {
        if (!token.value) return false

        try {
            const response = await fetch(`${baseURL}/api/auth/refresh`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token.value}`
                },
            })

            if (!response.ok) {
                logout()
                return false
            }

            const data = await response.json()
            token.value = data.token
            return true
        } catch {
            return false
        }
    }

    // Get current user profile
    const fetchProfile = async () => {
        if (!token.value) return null

        loading.value = true
        error.value = null

        try {
            const response = await fetch(`${baseURL}/api/me`, {
                headers: { 'Authorization': `Bearer ${token.value}` },
            })

            if (!response.ok) {
                if (response.status === 401) {
                    logout()
                }
                return null
            }

            const data = await response.json()
            user.value = data
            userCookie.value = data
            return data as AuthUser
        } catch (err: any) {
            error.value = err.message || 'Gagal memuat profil'
            return null
        } finally {
            loading.value = false
        }
    }

    // Initialize on client side
    if (process.client) {
        initAuth()
    }

    return {
        loading,
        error,
        user,
        token,
        isAuthenticated,
        isAdmin,
        isInstructor,
        permissions,
        hasPermission,
        login,
        adminLogin,
        register,
        logout,
        adminLogout,
        refreshToken,
        fetchProfile,
        initAuth,
    }
}
