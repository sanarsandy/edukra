import { defineStore } from 'pinia'

export interface User {
    id: string
    email: string
    full_name: string
    role?: string
    // Add other user fields as needed
}

export const useAuthStore = defineStore('auth', () => {
    const token = useCookie<string | null>('token', {
        maxAge: 60 * 60 * 24 * 3, // 3 days
        sameSite: 'lax',
        secure: process.env.NODE_ENV === 'production'
    })
    
    const user = useCookie<User | null>('user', {
        maxAge: 60 * 60 * 24 * 3, // 3 days
        sameSite: 'lax',
        secure: process.env.NODE_ENV === 'production'
    })

    const isAuthenticated = computed(() => !!token.value)

    function setToken(newToken: string) {
        token.value = newToken
    }

    function setUser(newUser: User) {
        user.value = newUser
    }

    function logout() {
        token.value = null
        user.value = null
        navigateTo('/login')
    }

    return {
        token,
        user,
        isAuthenticated,
        setToken,
        setUser,
        logout
    }
})

