export const useApi = () => {
    const config = useRuntimeConfig()
    // useCookie MUST be called at composable top level (synchronous)
    const tokenCookie = useCookie('token')

    // Determine API URL based on environment
    // Server: Use internal Docker network URL (apiInternal)
    // Client: Use public URL (apiBase) or proxy
    const apiUrl = process.server
        ? (config.apiInternal as string || config.public.apiBase || 'http://localhost:8080')
        : (config.public.apiBase || 'http://localhost:8080')

    const apiFetch = async <T>(
        endpoint: string,
        options: RequestInit = {}
    ): Promise<T> => {
        const url = `${apiUrl}${endpoint}`

        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        }

        // Read token value from the ref (tokenCookie is reactive)
        const token = tokenCookie.value
        if (token) {
            headers['Authorization'] = `Bearer ${token}`
        }

        try {
            const response = await $fetch<T>(url, {
                ...options,
                method: options.method as 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH',
                headers,
                credentials: 'include',
            })

            return response
        } catch (error: any) {
            // Handle 401 Unauthorized - clear token and redirect to login
            if (error.statusCode === 401 || error.status === 401) {
                // Only redirect if not already on login page and not an auth endpoint
                const isAuthEndpoint = endpoint.includes('/auth/')
                const currentPath = typeof window !== 'undefined' ? window.location.pathname : ''
                const isLoginPage = currentPath.includes('/login')

                if (!isAuthEndpoint && !isLoginPage && process.client) {
                    tokenCookie.value = null
                    navigateTo('/login')
                }
            }
            throw error
        }
    }

    return {
        apiUrl,
        fetch: apiFetch,
    }
}
