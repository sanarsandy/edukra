<template>
  <div class="min-h-screen flex">
    <!-- Left Side - Branding -->
    <div class="hidden lg:flex lg:w-1/2 bg-primary-600 relative overflow-hidden">
      <!-- Background Pattern -->
      <!-- Geometric Float Background -->
      <div class="absolute inset-0 z-0 overflow-hidden">
        <!-- Glass Cube -->
        <div class="absolute top-20 -left-4 w-32 h-32 bg-gradient-to-br from-white/60 to-white/10 backdrop-blur-3xl border border-white/40 rounded-2xl shadow-2xl animate-float-cube opacity-60"></div>
        
        <!-- Glass Sphere -->
        <div class="absolute bottom-40 right-10 w-48 h-48 bg-gradient-to-br from-white/40 to-white/5 backdrop-blur-2xl border border-white/20 rounded-full shadow-xl animate-float-sphere animation-delay-1000 opacity-50"></div>
        
        <!-- Glass Rounded Square -->
        <div class="absolute top-1/2 left-20 w-24 h-24 bg-gradient-to-br from-white/40 to-white/10 backdrop-blur-xl border border-white/20 rounded-xl shadow-lg animate-float-prism animation-delay-2000 opacity-40"></div>
        
        <!-- Background Glows -->
        <div class="absolute top-20 left-20 w-72 h-72 bg-white rounded-full blur-3xl opacity-10"></div>
        <div class="absolute bottom-20 right-20 w-96 h-96 bg-white rounded-full blur-3xl opacity-10"></div>
      </div>
      
      <!-- Content -->
      <div class="relative z-10 flex flex-col justify-center px-12 xl:px-20">
        <!-- Logo -->
        <div class="flex items-center space-x-3 mb-12">
          <img src="/logo.png" alt="EDUKRA Logo" class="w-12 h-12 object-contain" />
          <span class="font-display font-bold text-2xl text-white">EDUKRA</span>
        </div>
        
        <h1 class="text-4xl xl:text-5xl font-bold text-white leading-tight mb-6">
          Mulai Perjalanan Belajar Anda!
        </h1>
        <p class="text-lg text-white/80 mb-12 max-w-md">
          Bergabung dengan 10,000+ siswa yang telah mengembangkan skill mereka bersama EDUKRA.
        </p>
        
        <!-- Stats -->
        <div class="grid grid-cols-3 gap-6">
          <div class="bg-white/10 rounded-xl p-4 text-center">
            <div class="text-2xl font-bold text-white">500+</div>
            <div class="text-sm text-white/70">Kursus</div>
          </div>
          <div class="bg-white/10 rounded-xl p-4 text-center">
            <div class="text-2xl font-bold text-white">50+</div>
            <div class="text-sm text-white/70">Mentor</div>
          </div>
          <div class="bg-white/10 rounded-xl p-4 text-center">
            <div class="text-2xl font-bold text-white">98%</div>
            <div class="text-sm text-white/70">Puas</div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Right Side - Register Form -->
    <div class="w-full lg:w-1/2 flex items-center justify-center px-6 py-12 bg-neutral-50">
      <div class="w-full max-w-md">
        <!-- Mobile Logo -->
        <div class="flex items-center justify-center space-x-2 mb-8 lg:hidden">
          <img src="/logo.png" alt="EDUKRA Logo" class="w-10 h-10 object-contain" />
          <span class="font-display font-bold text-xl text-neutral-900">EDUKRA</span>
        </div>
        
        <!-- Header -->
        <div class="text-center lg:text-left mb-8">
          <h2 class="text-2xl md:text-3xl font-bold text-neutral-900 mb-2">Buat Akun Baru</h2>
          <p class="text-neutral-500">
            Sudah punya akun? 
            <NuxtLink to="/login" class="text-success-400 hover:text-success-500 font-medium">Masuk di sini</NuxtLink>
          </p>
        </div>
        
        <!-- Form -->
        <form @submit.prevent="handleRegister" class="space-y-5">
          <div>
            <label for="full_name" class="block text-sm font-medium text-neutral-700 mb-1.5">Nama Lengkap</label>
            <input
              id="full_name"
              v-model="fullName"
              type="text"
              required
              class="w-full px-4 py-3 border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all bg-white text-neutral-900 placeholder-neutral-400"
              placeholder="Masukkan nama lengkap"
            />
          </div>
          
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
          
          <div>
            <label for="password" class="block text-sm font-medium text-neutral-700 mb-1.5">Password</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              minlength="6"
              class="w-full px-4 py-3 border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all bg-white text-neutral-900 placeholder-neutral-400"
              placeholder="Minimal 6 karakter"
            />
          </div>
          
          <div class="flex items-start space-x-3">
            <input
              id="terms"
              v-model="agreeTerms"
              type="checkbox"
              required
              class="mt-1 h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded"
            />
            <label for="terms" class="text-sm text-neutral-600">
              Saya setuju dengan 
              <a href="#" class="text-primary-600 hover:underline">Syarat & Ketentuan</a> dan 
              <a href="#" class="text-primary-600 hover:underline">Kebijakan Privasi</a>
            </label>
          </div>
          
          <div v-if="error" class="p-3 bg-red-50 border border-red-200 rounded-xl">
            <p class="text-sm text-red-600">{{ error }}</p>
          </div>
          
          <div v-if="success" class="p-3 bg-accent-50 border border-accent-200 rounded-xl">
            <p class="text-sm text-accent-700">{{ success }}</p>
          </div>
          
          <button type="submit" :disabled="loading" class="w-full inline-flex items-center justify-center px-6 py-3.5 text-base font-semibold text-white bg-primary-600 rounded-xl hover:bg-primary-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="loading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Memproses...
            </span>
            <span v-else>Daftar Sekarang</span>
          </button>
        </form>
        
        <!-- Divider -->
        <div class="relative my-8">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-neutral-200"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-4 bg-neutral-50 text-neutral-500">atau daftar dengan</span>
          </div>
        </div>
        
        <button type="button" @click="handleGoogleRegister" class="w-full flex items-center justify-center px-4 py-3 border border-neutral-200 rounded-xl bg-white hover:bg-neutral-50 transition-colors">
          <svg class="w-5 h-5 mr-3" viewBox="0 0 24 24">
            <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
            <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
            <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
            <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
          </svg>
          <span class="font-medium text-neutral-700">Google</span>
        </button>
        
        <div class="mt-8 text-center">
          <NuxtLink to="/" class="text-sm text-neutral-500 hover:text-neutral-700">← Kembali ke Beranda</NuxtLink>
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
  title: 'Daftar - EDUKRA | The Engine of Learning'
})

const config = useRuntimeConfig()
const token = useCookie('token')
const userCookie = useCookie('user')

const fullName = ref('')
const email = ref('')
const password = ref('')
const agreeTerms = ref(false)
const loading = ref(false)
const error = ref('')
const success = ref('')

// Redirect if already logged in
onMounted(() => {
  if (token.value && userCookie.value) {
    navigateTo('/dashboard')
  }
})

const handleRegister = async () => {
  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const baseURL = config.public.apiBase || 'http://localhost:8080'
    
    const response = await fetch(`${baseURL}/api/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        full_name: fullName.value,
        email: email.value,
        password: password.value
      })
    })

    const data = await response.json()

    if (!response.ok) {
      error.value = data.error || 'Pendaftaran gagal. Silakan coba lagi.'
      return
    }

    // If token returned, auto login
    if (data.token) {
      token.value = data.token
      userCookie.value = data.user
      success.value = 'Akun berhasil dibuat! Mengalihkan ke dashboard...'
      setTimeout(() => {
        navigateTo('/dashboard')
      }, 1500)
    } else {
      success.value = 'Akun berhasil dibuat! Mengalihkan ke halaman login...'
      setTimeout(() => {
        navigateTo('/login')
      }, 2000)
    }
  } catch (err: any) {
    error.value = 'Terjadi kesalahan jaringan. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}

const handleGoogleRegister = async () => {
  try {
    const baseURL = config.public.apiBase || 'http://localhost:8080'
    const response = await fetch(`${baseURL}/api/auth/google`)
    const data = await response.json()
    
    if (data.auth_url) {
      window.location.href = data.auth_url
    }
  } catch (err: any) {
    error.value = 'Gagal memulai pendaftaran dengan Google.'
  }
}
</script>

<style scoped>
.animate-float-cube {
  animation: float-rotate-12 8s ease-in-out infinite;
}
.animate-float-sphere {
  animation: float-simple 10s ease-in-out infinite;
}
.animate-float-prism {
  animation: float-rotate-minus-12 9s ease-in-out infinite;
}
.animation-delay-1000 {
  animation-delay: 1s;
}
.animation-delay-2000 {
  animation-delay: 2s;
}

@keyframes float-rotate-12 {
  0%, 100% { transform: translateY(0) rotate(12deg); }
  50% { transform: translateY(-20px) rotate(17deg); }
}

@keyframes float-rotate-minus-12 {
  0%, 100% { transform: translateY(0) rotate(-12deg); }
  50% { transform: translateY(-15px) rotate(-7deg); }
}

@keyframes float-simple {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}
</style>
