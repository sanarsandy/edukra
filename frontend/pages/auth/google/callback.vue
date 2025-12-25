<template>
  <div class="min-h-screen flex items-center justify-center bg-neutral-50">
    <div class="text-center">
      <div class="w-16 h-16 mx-auto mb-4">
        <svg class="animate-spin w-16 h-16 text-primary-600" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
      <h2 class="text-xl font-semibold text-neutral-900 mb-2">Memproses Login...</h2>
      <p class="text-neutral-500">Mohon tunggu sebentar</p>
      <p v-if="error" class="mt-4 text-red-600">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

const route = useRoute()
const token = useCookie('token')
const userCookie = useCookie('user')
const error = ref('')

onMounted(async () => {
  try {
    // Get token and user from query params
    const tokenParam = route.query.token as string
    const userParam = route.query.user as string

    if (!tokenParam || !userParam) {
      error.value = 'Data login tidak valid'
      setTimeout(() => navigateTo('/login'), 2000)
      return
    }

    // Parse user data
    const userData = JSON.parse(decodeURIComponent(userParam))

    // Save to cookies
    token.value = tokenParam
    userCookie.value = userData

    // Redirect based on role
    if (userData.role === 'admin') {
      await navigateTo('/admin')
    } else if (userData.role === 'instructor') {
      await navigateTo('/instructor')
    } else {
      await navigateTo('/dashboard')
    }
  } catch (err: any) {
    console.error('Google callback error:', err)
    error.value = 'Terjadi kesalahan saat login'
    setTimeout(() => navigateTo('/login'), 2000)
  }
})
</script>
