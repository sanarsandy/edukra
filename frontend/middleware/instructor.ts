// Instructor middleware - checks if user is instructor
export default defineNuxtRouteMiddleware((to, from) => {
    const token = useCookie('token')
    const userCookie = useCookie<{ role?: string } | null>('user')

    if (!token.value) {
        return navigateTo('/instructor/login')
    }

    try {
        const payload = JSON.parse(atob(token.value.split('.')[1]))
        const exp = payload.exp * 1000

        // Check expiry
        if (Date.now() > exp) {
            token.value = null
            userCookie.value = null
            return navigateTo('/instructor/login')
        }

        // Check role - only instructor allowed
        if (payload.role !== 'instructor') {
            return navigateTo('/instructor/login')
        }
    } catch (e) {
        token.value = null
        userCookie.value = null
        return navigateTo('/instructor/login')
    }
})


