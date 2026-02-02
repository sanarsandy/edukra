<template>
  <div class="min-h-screen flex">
    <!-- Left Side - Branding (Same as Login) -->
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
          Lupa Password?
        </h1>
        <p class="text-lg text-white/80 mb-12 max-w-md">
          Jangan khawatir, kami akan membantu Anda mendapatkan kembali akses ke akun Anda.
        </p>
      </div>
    </div>
    
    <!-- Right Side - Form -->
    <div class="w-full lg:w-1/2 flex items-center justify-center px-6 py-12 bg-neutral-50">
      <div class="w-full max-w-md">
        <!-- Mobile Logo -->
        <div class="flex items-center justify-center space-x-2 mb-8 lg:hidden">
          <img src="/logo.png" alt="EDUKRA Logo" class="w-10 h-10 object-contain" />
          <span class="font-display font-bold text-xl text-neutral-900">EDUKRA</span>
        </div>
        
        <div class="text-center lg:text-left mb-8">
          <h2 class="text-2xl md:text-3xl font-bold text-neutral-900 mb-2">
            Reset Password
          </h2>
          <p class="text-neutral-500">
            Masukkan email Anda untuk menerima instruksi reset password.
          </p>
        </div>
        
        <form v-if="!success" @submit.prevent="handleSubmit" class="space-y-5">
          <div>
            <label for="email" class="block text-sm font-medium text-neutral-700 mb-1.5">Email</label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              class="w-full px-4 py-3 border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all bg-white text-neutral-900 placeholder-neutral-400"
              placeholder="nama@email.com"
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
            <span v-else>Kirim Link Reset</span>
          </button>
        </form>

        <div v-else class="text-center p-6 bg-green-50 rounded-2xl border border-green-100">
          <div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
            </svg>
          </div>
          <h3 class="text-lg font-bold text-neutral-900 mb-2">Cek Email Anda</h3>
          <p class="text-neutral-600 mb-6">
            Kami telah mengirimkan instruksi reset password ke <strong>{{ email }}</strong>. Silakan cek inbox atau folder spam Anda.
          </p>
          <button @click="success = false" class="text-primary-600 hover:text-primary-700 font-medium">
            Kirim ulang
          </button>
        </div>
        
        <div class="mt-8 text-center">
          <NuxtLink to="/login" class="text-sm text-neutral-500 hover:text-neutral-700">← Kembali ke Login</NuxtLink>
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
  title: 'Lupa Password - EDUKRA'
})

const config = useRuntimeConfig()
const email = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)

const handleSubmit = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const baseURL = config.public.apiBase || 'http://localhost:8080'
    const response = await fetch(`${baseURL}/api/auth/forgot-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value })
    })

    const data = await response.json()

    if (!response.ok) {
        // Even if error, if security dictates, we might want to show success or generic error.
        // But backend sends 200 for user not found too.
        // If 429 (rate limit), show error.
        if (response.status === 429) {
            error.value = 'Terlalu banyak permintaan. Silakan coba lagi nanti.'
            return
        }
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
