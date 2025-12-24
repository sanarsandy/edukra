// Admin middleware - checks if user is admin
export default defineNuxtRouteMiddleware((to, from) => {
    // Get token and user from cookies
    const token = useCookie('token')
    const userCookie = useCookie<{ role?: string } | null>('user')

    // Check if token exists
    if (!token.value) {
        return navigateTo('/admin/login')
    }

    // Check token and role
    try {
        const payload = JSON.parse(atob(token.value.split('.')[1]))
        const exp = payload.exp * 1000

        // Check expiry
        if (Date.now() > exp) {
            token.value = null
            userCookie.value = null
            return navigateTo('/admin/login')
        }

        // Check role from JWT payload (more secure than cookie)
        if (payload.role !== 'admin') {
            return navigateTo('/admin/login')
        }
    } catch (e) {
        token.value = null
        userCookie.value = null
        return navigateTo('/admin/login')
    }
})
