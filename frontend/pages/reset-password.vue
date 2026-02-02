<template>
  <div class="min-h-screen flex">
    <!-- Left Side - Branding -->
    <div class="hidden lg:flex lg:w-1/2 bg-primary-600 relative overflow-hidden">
      <div class="absolute inset-0 z-0 overflow-hidden">
        <div class="absolute top-20 -left-4 w-32 h-32 bg-gradient-to-br from-white/60 to-white/10 backdrop-blur-3xl border border-white/40 rounded-2xl shadow-2xl animate-float-cube opacity-60"></div>
        <div class="absolute bottom-40 right-10 w-48 h-48 bg-gradient-to-br from-white/40 to-white/5 backdrop-blur-2xl border border-white/20 rounded-full shadow-xl animate-float-sphere animation-delay-1000 opacity-50"></div>
        <div class="absolute top-1/2 left-20 w-24 h-24 bg-gradient-to-br from-white/40 to-white/10 backdrop-blur-xl border border-white/20 rounded-xl shadow-lg animate-float-prism animation-delay-2000 opacity-40"></div>
        <div class="absolute top-20 left-20 w-72 h-72 bg-white rounded-full blur-3xl opacity-10"></div>
        <div class="absolute bottom-20 right-20 w-96 h-96 bg-white rounded-full blur-3xl opacity-10"></div>
      </div>
      
      <div class="relative z-10 flex flex-col justify-center px-12 xl:px-20">
        <div class="flex items-center space-x-3 mb-12">
          <img src="/logo.png" alt="EDUKRA Logo" class="w-12 h-12 object-contain" />
          <span class="font-display font-bold text-2xl text-white">EDUKRA</span>
        </div>
        <h1 class="text-4xl xl:text-5xl font-bold text-white leading-tight mb-6">
          Buat Password Baru
        </h1>
        <p class="text-lg text-white/80 mb-12 max-w-md">
          Amankan akun Anda dengan password baru yang kuat.
        </p>
      </div>
    </div>
    
    <!-- Right Side - Form -->
    <div class="w-full lg:w-1/2 flex items-center justify-center px-6 py-12 bg-neutral-50">
      <div class="w-full max-w-md">
        <div class="text-center lg:text-left mb-8">
          <h2 class="text-2xl md:text-3xl font-bold text-neutral-900 mb-2">
            Reset Password
          </h2>
          <p class="text-neutral-500">
            Masukkan password baru Anda di bawah ini.
          </p>
        </div>
        
        <form v-if="!success" @submit.prevent="handleSubmit" class="space-y-5">
          <!-- Password -->
          <div>
            <label for="password" class="block text-sm font-medium text-neutral-700 mb-1.5">Password Baru</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              minlength="8"
              class="w-full px-4 py-3 border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all bg-white text-neutral-900"
              placeholder="Minimal 8 karakter"
            />
          </div>

          <!-- Confirm Password -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-neutral-700 mb-1.5">Konfirmasi Password</label>
            <input
              id="confirmPassword"
              v-model="confirmPassword"
              type="password"
              required
              class="w-full px-4 py-3 border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all bg-white text-neutral-900"
              placeholder="Ulangi password baru"
            />
          </div>
          
          <div v-if="error" class="p-3 bg-red-50 border border-red-200 rounded-xl">
            <p class="text-sm text-red-600">{{ error }}</p>
          </div>
          
          <button type="submit" :disabled="loading" class="w-full inline-flex items-center justify-center px-6 py-3.5 text-base font-semibold text-white bg-primary-600 rounded-xl hover:bg-primary-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="loading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Memproses...
            </span>
            <span v-else>Simpan Password Baru</span>
          </button>
        </form>

        <div v-else class="text-center p-6 bg-green-50 rounded-2xl border border-green-100">
          <div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
          </div>
          <h3 class="text-lg font-bold text-neutral-900 mb-2">Berhasil!</h3>
          <p class="text-neutral-600 mb-6">
            Password Anda berhasil diperbarui. Silakan login menggunakan password baru Anda.
          </p>
          <NuxtLink to="/login" class="inline-block px-6 py-3 bg-primary-600 text-white font-semibold rounded-xl hover:bg-primary-700 transition-colors">
            Masuk Sekarang
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

useHead({
  title: 'Reset Password - EDUKRA'
})

const route = useRoute()
const config = useRuntimeConfig()

const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)

const token = computed(() => route.query.token as string)

onMounted(() => {
  if (!token.value) {
    error.value = 'Token tidak valid atau tidak ditemukan.'
  }
})

const handleSubmit = async () => {
  if (password.value !== confirmPassword.value) {
    error.value = 'Password konfirmasi tidak cocok.'
    return
  }

  loading.value = true
  error.value = ''
  
  try {
    const baseURL = config.public.apiBase || 'http://localhost:8080'
    const response = await fetch(`${baseURL}/api/auth/reset-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        token: token.value,
        password: password.value
      })
    })

    const data = await response.json()

    if (!response.ok) {
        error.value = data.error || 'Gagal mereset password.'
        return
    }

    success.value = true
  } catch (err) {
    error.value = 'Terjadi kesalahan jaringan. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.animate-float-cube { animation: float-rotate-12 8s ease-in-out infinite; }
.animate-float-sphere { animation: float-simple 10s ease-in-out infinite; }
.animate-float-prism { animation: float-rotate-minus-12 9s ease-in-out infinite; }
.animation-delay-1000 { animation-delay: 1s; }
.animation-delay-2000 { animation-delay: 2s; }
@keyframes float-rotate-12 { 0%, 100% { transform: translateY(0) rotate(12deg); } 50% { transform: translateY(-20px) rotate(17deg); } }
@keyframes float-rotate-minus-12 { 0%, 100% { transform: translateY(0) rotate(-12deg); } 50% { transform: translateY(-15px) rotate(-7deg); } }
@keyframes float-simple { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-20px); } }
</style>
