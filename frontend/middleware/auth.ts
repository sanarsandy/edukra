// Auth middleware - checks if user is authenticated as student/instructor (NOT admin)
// Admin users should use /admin routes instead
export default defineNuxtRouteMiddleware((to, from) => {
    // Get token from cookie
    const token = useCookie('token')

    // Check if token exists
    if (!token.value) {
        return navigateTo('/login')
    }

    // Check token validity and role
    try {
        const payload = JSON.parse(atob(token.value.split('.')[1]))
        const exp = payload.exp * 1000

        // Check expiry
        if (Date.now() > exp) {
            const userCookie = useCookie('user')
            token.value = null
            userCookie.value = null
            return navigateTo('/login')
        }

        // If user is admin, redirect them to admin area
        if (payload.role === 'admin') {
            return navigateTo('/admin')
        }
    } catch (e) {
        const userCookie = useCookie('user')
        token.value = null
        userCookie.value = null
        return navigateTo('/login')
    }
})
