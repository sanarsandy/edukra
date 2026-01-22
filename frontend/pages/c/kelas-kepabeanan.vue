<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'

// SEO & Head
definePageMeta({
  layout: false
})

useHead({
  title: 'Belajar Kepabeanan & Cukai - Webinar Gratis',
  bodyAttrs: {
    class: '!bg-[#050505] !m-0 !p-0 !overflow-x-hidden'
  },
  meta: [
    { name: 'description', content: 'Webinar paling lit buat kamu yang pengen paham dunia ekspor-impor. No boring lecture, 100% real talk!' }
  ],
  link: [
    { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700;800;900&display=swap' }
  ]
})

// Runtime Config for API
const config = useRuntimeConfig()
const apiBase = config.public.apiBase || 'http://localhost:8080'

// Fetch Data Logic
// We hardcode the slug 'kelas-kepabeanan' here because this is a specialized page for that specific campaign.
// Make sure a campaign with slug 'kelas-kepabeanan' exists in the Admin Panel.
const { data: campaign, error: fetchError } = await useFetch('/api/c/kelas-kepabeanan', {
  baseURL: apiBase,
  key: 'custom-landing-kepabeanan'
})

// Form Data
const form = ref({
  fullName: '',
  email: '',
  phone: '',
})
const submitting = ref(false)
const submitError = ref('')
const submitSuccess = ref('')

// Countdown Logic
const targetDate = computed(() => {
    return campaign.value?.end_date ? new Date(campaign.value.end_date).getTime() : new Date('2026-01-27T19:00:00+07:00').getTime()
})

const countdown = ref({ days: '00', hours: '00', minutes: '00', seconds: '00' })
let timer: ReturnType<typeof setInterval>

const updateCountdown = () => {
  const now = new Date().getTime()
  const distance = targetDate.value - now

  if (distance < 0) {
    countdown.value = { days: '00', hours: '00', minutes: '00', seconds: '00' }
    clearInterval(timer)
    return
  }

  const d = Math.floor(distance / (1000 * 60 * 60 * 24))
  const h = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const m = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60))
  const s = Math.floor((distance % (1000 * 60)) / 1000)

  countdown.value = {
    days: d < 10 ? `0${d}` : `${d}`,
    hours: h < 10 ? `0${h}` : `${h}`,
    minutes: m < 10 ? `0${m}` : `${m}`,
    seconds: s < 10 ? `0${s}` : `${s}`
  }
}

onMounted(() => {
  updateCountdown()
  timer = setInterval(updateCountdown, 1000)
})

onUnmounted(() => {
  clearInterval(timer)
})

const handleRegister = async () => {
  submitError.value = ''
  submitSuccess.value = ''

  // Fallback if API hasn't loaded yet?
  if (!campaign.value || !campaign.value.id) {
    submitError.value = 'Data campaign tidak ditemukan. Pastikan campaign dengan slug "kelas-kepabeanan" sudah dibuat di Admin.'
    return
  }

  if (!form.value.email || !form.value.phone) {
    submitError.value = 'Email dan Nomor WhatsApp wajib diisi.'
    return
  }

  submitting.value = true
  try {
    const response = await $fetch<{
      payment_url?: string
      is_free?: boolean
      message?: string
    }>('/api/campaign-checkout', {
      method: 'POST',
      baseURL: apiBase,
      body: {
        campaign_id: campaign.value.id,
        email: form.value.email,
        phone: form.value.phone,
        full_name: form.value.fullName,
        payment_method: '', 
      }
    })

    if (response.is_free || response.message) {
      submitSuccess.value = response.message || 'Pendaftaran berhasil! Cek WhatsApp Anda.'
      form.value = { fullName: '', email: '', phone: '' } // Reset form
      
      // Optional: Redirect to WhatsApp or Thank You page?
      // window.location.href = `https://wa.me/...`
    } else if (response.payment_url) {
      window.location.href = response.payment_url
    }
  } catch (err: any) {
    console.error(err)
    submitError.value = err.data?.error || 'Gagal memproses pendaftaran.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#050505] text-white font-['Outfit'] overflow-x-hidden selection:bg-[#ff3399] selection:text-white">
    
    <!-- Hero Section -->
    <header class="relative min-h-screen flex flex-col justify-center items-center px-4 pt-10 text-center">
      
      <!-- Top Badge -->
      <div class="inline-flex items-center gap-2 bg-neutral-900 border border-neutral-800 rounded-full pl-2 pr-4 py-1.5 mb-12 hover:scale-105 transition-transform cursor-pointer shadow-[0_0_20px_rgba(255,51,153,0.15)]">
         <span class="text-xl">✨</span>
         <span class="text-[#84cc16] text-sm font-bold uppercase tracking-wider">100% GRATIS!</span>
         <span class="text-xl">🔥</span>
      </div>
      
      <!-- Floating Emojis (3D Style) -->
      <div class="absolute top-[20%] left-[10%] text-5xl md:text-7xl animate-float-slow opacity-90 hidden md:block">🚀</div>
      <div class="absolute top-[15%] right-[15%] text-5xl md:text-7xl animate-float-medium opacity-90 hidden md:block">🌍</div>
      <div class="absolute bottom-[30%] left-[15%] text-5xl md:text-7xl animate-float-fast opacity-90 hidden md:block">💡</div>
      <div class="absolute bottom-[35%] right-[10%] text-5xl md:text-7xl animate-float-slow opacity-90 hidden md:block">💎</div>

      <!-- Main Title -->
      <h1 class="relative z-10 text-6xl sm:text-7xl md:text-9xl font-black leading-tight tracking-tight mb-8">
        BELAJAR <br />
        <span class="text-transparent bg-clip-text bg-gradient-to-r from-[#e879f9] via-[#d946ef] to-[#f97316]">
          KEPABEANAN <br /> & CUKAI
        </span>
      </h1>

      <!-- Subtitle -->
      <p class="relative z-10 text-lg sm:text-xl text-neutral-400 max-w-2xl mx-auto mb-12 font-medium">
        Webinar paling lit buat kamu yang pengen paham dunia ekspor-impor. 
        <span class="text-white font-bold">No boring lecture, 100% real talk! 🎤</span>
      </p>

      <!-- Info Pills -->
      <div class="flex flex-wrap justify-center gap-4 mb-20 relative z-10">
        <div class="flex items-center gap-2 bg-neutral-900/80 border border-neutral-800 px-5 py-2.5 rounded-full backdrop-blur-md hover:border-[#d946ef]/50 transition-colors">
          <span class="text-[#a855f7]">📅</span>
          <span class="font-medium text-sm text-neutral-300">27 Jan 2026</span>
        </div>
        <div class="flex items-center gap-2 bg-neutral-900/80 border border-neutral-800 px-5 py-2.5 rounded-full backdrop-blur-md hover:border-[#d946ef]/50 transition-colors">
          <span class="text-[#ec4899]">⏰</span>
          <span class="font-medium text-sm text-neutral-300">19:00 WIB</span>
        </div>
        <div class="flex items-center gap-2 bg-neutral-900/80 border border-neutral-800 px-5 py-2.5 rounded-full backdrop-blur-md hover:border-[#d946ef]/50 transition-colors">
          <img src="/gmeets_logo.png" alt="Google Meet" class="w-5 h-5 object-contain" />
          <span class="font-medium text-sm text-neutral-300">Google Meet</span>
        </div>
      </div>

       <!-- Box Icon -->
       <div class="text-7xl mb-4 animate-bounce relative z-10">📦</div>

      <!-- Countdown Title -->
      <p class="text-xs font-bold text-neutral-500 uppercase tracking-[0.2em] mb-8 relative z-10">WEBINAR DIMULAI DALAM</p>

      <!-- Countdown -->
      <div class="flex flex-wrap justify-center gap-4 relative z-10 mb-20">
        <!-- Days -->
        <div class="text-center">
            <div class="w-24 h-24 bg-neutral-900/50 backdrop-blur-md rounded-3xl border border-white/5 flex items-center justify-center shadow-[0_0_30px_-5px_rgba(168,85,247,0.3)] hover:scale-105 transition-transform">
                <span class="text-5xl font-black bg-gradient-to-b from-[#e879f9] to-[#a855f7] bg-clip-text text-transparent">{{ countdown.days }}</span>
            </div>
            <div class="text-[10px] font-bold text-neutral-500 uppercase mt-3 tracking-wider">HARI</div>
        </div>
        <!-- Hours -->
        <div class="text-center">
            <div class="w-24 h-24 bg-neutral-900/50 backdrop-blur-md rounded-3xl border border-white/5 flex items-center justify-center shadow-[0_0_30px_-5px_rgba(236,72,153,0.3)] hover:scale-105 transition-transform">
                <span class="text-5xl font-black bg-gradient-to-b from-[#f472b6] to-[#ec4899] bg-clip-text text-transparent">{{ countdown.hours }}</span>
            </div>
            <div class="text-[10px] font-bold text-neutral-500 uppercase mt-3 tracking-wider">JAM</div>
        </div>
        <!-- Minutes -->
        <div class="text-center">
            <div class="w-24 h-24 bg-neutral-900/50 backdrop-blur-md rounded-3xl border border-white/5 flex items-center justify-center shadow-[0_0_30px_-5px_rgba(59,130,246,0.3)] hover:scale-105 transition-transform">
                <span class="text-5xl font-black bg-gradient-to-b from-[#60a5fa] to-[#3b82f6] bg-clip-text text-transparent">{{ countdown.minutes }}</span>
            </div>
            <div class="text-[10px] font-bold text-neutral-500 uppercase mt-3 tracking-wider">MENIT</div>
        </div>
        <!-- Seconds -->
        <div class="text-center">
            <div class="w-24 h-24 bg-neutral-900/50 backdrop-blur-md rounded-3xl border border-white/5 flex items-center justify-center shadow-[0_0_30px_-5px_rgba(132,204,22,0.3)] hover:scale-105 transition-transform">
                <span class="text-5xl font-black bg-gradient-to-b from-[#a3e635] to-[#84cc16] bg-clip-text text-transparent">{{ countdown.seconds }}</span>
            </div>
            <div class="text-[10px] font-bold text-neutral-500 uppercase mt-3 tracking-wider">DETIK</div>
        </div>
      </div>
    </header>

    <!-- Features Section -->
    <section class="max-w-7xl mx-auto px-4 py-20">
      <div class="text-center mb-16">
          <div class="text-6xl mb-4">🎯</div>
          <h2 class="text-4xl md:text-5xl font-black mb-4">
              Yang Bakal Kamu <span class="text-transparent bg-clip-text bg-gradient-to-r from-[#e879f9] to-[#f97316]">DAPETIN</span>
          </h2>
          <p class="text-neutral-400 max-w-2xl mx-auto">Materi super praktis yang langsung bisa dipake. Cocok buat pemula sampai yang udah pro! 💪</p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          
        <!-- Card 1 -->
        <div class="bg-gradient-to-br from-purple-900/20 to-neutral-900 border border-white/5 p-10 rounded-[32px] relative overflow-hidden group hover:border-purple-500/30 transition-all">
            <div class="absolute top-0 right-0 w-[300px] h-[300px] bg-purple-600/10 rounded-full blur-[80px] -translate-y-1/2 translate-x-1/2 group-hover:bg-purple-600/20 transition-colors"></div>
            
            <div class="w-16 h-16 bg-gradient-to-br from-[#e879f9] to-[#a855f7] rounded-3xl flex items-center justify-center text-3xl mb-6 shadow-lg shadow-purple-500/20 group-hover:scale-110 transition-transform">📖</div>
            
            <h3 class="text-2xl font-bold mb-3 relative z-10">Kepabeanan 101</h3>
            <p class="text-neutral-400 text-sm leading-relaxed relative z-10">
                Dari nol sampai paham! Konsep dasar yang wajib tau sebelum terjun ke dunia ekspor-impor.
            </p>
            <div class="absolute top-4 right-4 text-2xl">📚</div>
        </div>

        <!-- Card 2 -->
        <div class="bg-gradient-to-br from-pink-900/20 to-neutral-900 border border-white/5 p-10 rounded-[32px] relative overflow-hidden group hover:border-pink-500/30 transition-all">
            <div class="absolute top-0 right-0 w-[300px] h-[300px] bg-pink-600/10 rounded-full blur-[80px] -translate-y-1/2 translate-x-1/2 group-hover:bg-pink-600/20 transition-colors"></div>
            
            <div class="w-16 h-16 bg-gradient-to-br from-[#f472b6] to-[#ec4899] rounded-3xl flex items-center justify-center text-3xl mb-6 shadow-lg shadow-pink-500/20 group-hover:scale-110 transition-transform">📋</div>
            
            <h3 class="text-2xl font-bold mb-3 relative z-10">Anti Ribet Cukai</h3>
            <p class="text-neutral-400 text-sm leading-relaxed relative z-10">
                Step-by-step prosedur cukai yang simple. Goodbye dokumen berantakan!
            </p>
            <div class="absolute top-4 right-4 text-2xl">✨</div>
        </div>

        <!-- Card 3 -->
        <div class="bg-gradient-to-br from-blue-900/20 to-neutral-900 border border-white/5 p-10 rounded-[32px] relative overflow-hidden group hover:border-blue-500/30 transition-all">
            <div class="absolute top-0 right-0 w-[300px] h-[300px] bg-blue-600/10 rounded-full blur-[80px] -translate-y-1/2 translate-x-1/2 group-hover:bg-blue-600/20 transition-colors"></div>
            
            <div class="w-16 h-16 bg-gradient-to-br from-[#60a5fa] to-[#3b82f6] rounded-3xl flex items-center justify-center text-3xl mb-6 shadow-lg shadow-blue-500/20 group-hover:scale-110 transition-transform">🌐</div>
            
            <h3 class="text-2xl font-bold mb-3 relative z-10">Go International!</h3>
            <p class="text-neutral-400 text-sm leading-relaxed relative z-10">
                Tips praktis buat kamu yang mau mulai atau scale up bisnis ke pasar global.
            </p>
            <div class="absolute top-4 right-4 text-2xl">🌍</div>
        </div>

        <!-- Card 4 -->
        <div class="bg-gradient-to-br from-orange-900/20 to-neutral-900 border border-white/5 p-10 rounded-[32px] relative overflow-hidden group hover:border-orange-500/30 transition-all">
            <div class="absolute top-0 right-0 w-[300px] h-[300px] bg-orange-600/10 rounded-full blur-[80px] -translate-y-1/2 translate-x-1/2 group-hover:bg-orange-600/20 transition-colors"></div>
            
            <div class="w-16 h-16 bg-gradient-to-br from-[#fbbf24] to-[#f59e0b] rounded-3xl flex items-center justify-center text-3xl mb-6 shadow-lg shadow-orange-500/20 group-hover:scale-110 transition-transform">📈</div>
            
            <h3 class="text-2xl font-bold mb-3 relative z-10">Karir yang Cuan</h3>
            <p class="text-neutral-400 text-sm leading-relaxed relative z-10">
                Peluang karir di bidang kepabeanan & cukai yang lagi hot banget di 2026!
            </p>
            <div class="absolute top-4 right-4 text-2xl">💰</div>
        </div>
      </div>
    </section>

    <!-- Registration CTA -->
    <section id="register" class="py-24 px-4 text-center">
        <div class="max-w-2xl mx-auto bg-neutral-900/50 backdrop-blur-xl border border-white/10 p-10 rounded-[40px] relative overflow-hidden">
            <!-- Glows -->
            <div class="absolute top-0 left-0 w-full h-full bg-gradient-to-b from-purple-500/5 to-transparent pointer-events-none"></div>

            <h2 class="text-3xl sm:text-4xl font-black mb-4 relative z-10">Siap Jadi Master Kepabeanan? 🚀</h2>
            <p class="text-neutral-400 mb-8 relative z-10">Slot terbatas! Amankan kursimu sekarang sebelum penuh.</p>

            <form @submit.prevent="handleRegister" class="relative z-10 space-y-4 max-w-sm mx-auto">
                <input v-model="form.fullName" type="text" placeholder="Nama Lengkap (Optional)" class="w-full bg-black/40 border border-white/10 rounded-2xl px-6 py-4 text-white placeholder-neutral-500 focus:border-[#d946ef] focus:ring-1 focus:ring-[#d946ef] outline-none transition-all" />
                <input v-model="form.email" type="email" placeholder="Email Address *" required class="w-full bg-black/40 border border-white/10 rounded-2xl px-6 py-4 text-white placeholder-neutral-500 focus:border-[#d946ef] focus:ring-1 focus:ring-[#d946ef] outline-none transition-all" />
                <input v-model="form.phone" type="tel" placeholder="WhatsApp Number *" required class="w-full bg-black/40 border border-white/10 rounded-2xl px-6 py-4 text-white placeholder-neutral-500 focus:border-[#d946ef] focus:ring-1 focus:ring-[#d946ef] outline-none transition-all" />
                
                <div v-if="submitError" class="text-red-500 text-sm font-bold">{{ submitError }}</div>
                <div v-if="submitSuccess" class="text-green-500 text-sm font-bold">{{ submitSuccess }}</div>

                <button :disabled="submitting" class="w-full bg-gradient-to-r from-[#ec4899] to-[#f97316] text-white font-bold text-lg py-4 rounded-2xl shadow-[0_10px_30px_rgba(236,72,153,0.3)] hover:shadow-[0_10px_40px_rgba(236,72,153,0.5)] hover:scale-[1.02] active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed">
                    <span v-if="submitting">Memproses...</span>
                    <span v-else>Daftar Sekarang (Gratis)</span>
                </button>
            </form>
            <p class="text-neutral-600 text-sm mt-6">No credit card required. Instant access.</p>
        </div>
    </section>

    <footer class="text-center py-10 border-t border-white/5 text-neutral-600 text-sm">
        &copy; 2026 Customs Youth Spark. All rights reserved.
    </footer>
    
  </div>
</template>

<style scoped>
@keyframes float-slow {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}
@keyframes float-medium {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-15px); }
}
@keyframes float-fast {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.animate-float-slow { animation: float-slow 4s ease-in-out infinite; }
.animate-float-medium { animation: float-medium 3s ease-in-out infinite; }
.animate-float-fast { animation: float-fast 2.5s ease-in-out infinite; }
</style>
