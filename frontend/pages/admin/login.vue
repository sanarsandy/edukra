<template>
  <div class="min-h-screen flex items-center justify-center bg-neutral-900 px-4">
    <div class="w-full max-w-md">
      <!-- Logo and Title -->
      <div class="text-center mb-8">
        <div class="w-16 h-16 bg-gradient-to-br from-red-500 to-red-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-white">Admin Panel</h1>
        <p class="text-neutral-400 mt-1">Masuk ke dashboard administrator</p>
      </div>

      <!-- Login Form -->
      <form @submit.prevent="handleLogin" class="bg-neutral-800 rounded-2xl p-8 shadow-xl border border-neutral-700">
        <!-- Error Message -->
        <div v-if="error" class="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
          <p class="text-red-400 text-sm">{{ error }}</p>
        </div>

        <!-- Email Input -->
        <div class="mb-5">
          <label for="email" class="block text-sm font-medium text-neutral-300 mb-2">Email</label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            autocomplete="email"
            placeholder="admin@lms.local"
            class="w-full px-4 py-3 bg-neutral-700 border border-neutral-600 rounded-lg text-white placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all"
          />
        </div>

        <!-- Password Input -->
        <div class="mb-6">
          <label for="password" class="block text-sm font-medium text-neutral-300 mb-2">Password</label>
          <div class="relative">
            <input
              id="password"
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              required
              autocomplete="current-password"
              placeholder="••••••••"
              class="w-full px-4 py-3 bg-neutral-700 border border-neutral-600 rounded-lg text-white placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all pr-12"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-white transition-colors"
            >
              <svg v-if="showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
              </svg>
              <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
              </svg>
            </button>
          </div>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="loading"
          class="w-full py-3 px-4 bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white font-semibold rounded-lg transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg v-if="loading" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
          </svg>
          <span>{{ loading ? 'Memproses...' : 'Masuk sebagai Admin' }}</span>
        </button>
      </form>

      <!-- Back Link -->
      <div class="text-center mt-6">
        <NuxtLink to="/login" class="text-neutral-400 hover:text-white text-sm transition-colors">
          ← Kembali ke Login User
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

useHead({
  title: 'Admin Login - LMS',
})

const config = useRuntimeConfig()
const token = useCookie('token')
const userCookie = useCookie('user')

const form = ref({
  email: '',
  password: '',
})

const loading = ref(false)
const error = ref('')
const showPassword = ref(false)

// Redirect if already logged in as admin
onMounted(() => {
  if (token.value && userCookie.value) {
    try {
      const user = typeof userCookie.value === 'string' 
        ? JSON.parse(userCookie.value) 
        : userCookie.value
      if (user.role === 'admin') {
        navigateTo('/admin')
      }
    } catch (e) {
      // Invalid cookie, ignore
    }
  }
})

async function handleLogin() {
  loading.value = true
  error.value = ''

  try {
    const baseURL = config.public.apiBase || 'http://localhost:8080'
    
    const response = await fetch(`${baseURL}/api/auth/admin/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(form.value),
    })

    const data = await response.json()

    if (!response.ok) {
      error.value = data.error || 'Login gagal'
      return
    }

    // Save token and user to cookies
    token.value = data.token
    userCookie.value = data.user

    // Redirect to admin dashboard
    await navigateTo('/admin')
  } catch (err: any) {
    error.value = 'Terjadi kesalahan jaringan'
    console.error('Admin login error:', err)
  } finally {
    loading.value = false
  }
}
</script>
