<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps<{
  block: any
  styles?: any
  globalEndDate?: string | null
}>()

// Countdown Logic
const targetDate = computed(() => {
  const dateStr = props.globalEndDate || props.block.data.end_date
  return dateStr ? new Date(dateStr).getTime() : new Date('2026-02-15T19:00:00+07:00').getTime()
})

const countdown = ref({ days: '00', hours: '00', minutes: '00', seconds: '00' })
const displayDate = computed(() => {
  const dateStr = props.globalEndDate || props.block.data.end_date
  if (!dateStr) return '15 Feb 2026'
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }).format(date)
})
const displayTime = computed(() => {
  const dateStr = props.globalEndDate || props.block.data.end_date
  if (!dateStr) return '19:00 WIB'
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('id-ID', { hour: '2-digit', minute: '2-digit', timeZoneName: 'short' }).format(date)
})

let timer: any = null

const updateCountdown = () => {
  const now = new Date().getTime()
  const distance = targetDate.value - now

  if (distance < 0) {
    countdown.value = { days: '00', hours: '00', minutes: '00', seconds: '00' }
    if (timer) clearInterval(timer)
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

watch(() => [props.block.data.end_date, props.globalEndDate], () => {
    updateCountdown()
})

onMounted(() => {
  updateCountdown()
  if (typeof window !== 'undefined') {
    timer = setInterval(updateCountdown, 1000)
  }
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

// Speaker data from block or defaults
const speakerImage = computed(() => props.block.data.speaker_image || 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=400')
const speakerName = computed(() => props.block.data.speaker_name || 'Joni')
const speakerTitle = computed(() => props.block.data.speaker_title || 'Ex Auditor Bea & Cukai')
</script>

<template>
  <header class="relative min-h-screen flex flex-col justify-center px-4 py-16 bg-gradient-to-br from-blue-950 via-blue-900 to-slate-900 font-inter overflow-hidden">
      <!-- Background Pattern -->
      <div class="absolute inset-0 opacity-10">
        <div class="absolute top-0 left-0 w-full h-full" style="background-image: url('data:image/svg+xml,%3Csvg width=\'60\' height=\'60\' viewBox=\'0 0 60 60\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cg fill=\'none\' fill-rule=\'evenodd\'%3E%3Cg fill=\'%23ffffff\' fill-opacity=\'0.1\'%3E%3Cpath d=\'M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z\'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E');"></div>
      </div>
      
      <!-- Decorative Elements -->
      <div class="absolute top-20 right-10 w-72 h-72 bg-amber-500/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-20 left-10 w-96 h-96 bg-blue-500/20 rounded-full blur-3xl"></div>
      
      <!-- Customs Icons Decoration -->
      <div class="absolute top-1/4 left-10 text-6xl opacity-20 hidden lg:block">📦</div>
      <div class="absolute bottom-1/3 right-16 text-5xl opacity-20 hidden lg:block">🛃</div>
      <div class="absolute top-1/2 right-1/4 text-4xl opacity-15 hidden lg:block">🌐</div>
      
      <div class="max-w-7xl mx-auto w-full relative z-10">
          <div class="grid lg:grid-cols-2 gap-12 lg:gap-16 items-center">
              <!-- Left Content -->
              <div class="text-left text-white order-2 lg:order-1">
                  <!-- Badge -->
                  <div v-if="block.data.badge" class="inline-flex items-center gap-2 bg-amber-500/20 border border-amber-400/30 backdrop-blur-sm rounded-full px-4 py-2 mb-6">
                     <span class="w-2 h-2 bg-amber-400 rounded-full animate-pulse"></span>
                     <span class="text-amber-300 text-sm font-semibold uppercase tracking-wide">{{ block.data.badge }}</span>
                  </div>
                  
                  <!-- Headline -->
                  <h1 class="text-4xl sm:text-5xl lg:text-6xl font-bold leading-tight mb-6 text-amber-400">
                    {{ block.data.headline || 'Webinar Gratis: Bea & Cukai Itu Gampang!' }}
                  </h1>
                  
                  <!-- Subheadline -->
                  <p class="text-lg text-blue-100/80 leading-relaxed mb-8 max-w-xl">
                    {{ block.data.subheadline || 'Pahami regulasi kepabeanan dengan cara yang mudah dipahami. Cocok untuk pemula!' }}
                  </p>
                  
                  <!-- Event Info Pills -->
                  <div class="flex flex-wrap gap-3 mb-8">
                    <div class="flex items-center gap-2 bg-white/10 backdrop-blur-sm border border-white/10 rounded-xl px-4 py-3">
                      <span class="text-amber-400">📅</span>
                      <span class="text-white font-medium">{{ displayDate }}</span>
                    </div>
                    <div class="flex items-center gap-2 bg-white/10 backdrop-blur-sm border border-white/10 rounded-xl px-4 py-3">
                      <span class="text-amber-400">⏰</span>
                      <span class="text-white font-medium">{{ displayTime }}</span>
                    </div>
                    <div class="flex items-center gap-2 bg-white/10 backdrop-blur-sm border border-white/10 rounded-xl px-4 py-3">
                      <img src="/gmeets_logo.png" alt="Google Meet" class="w-5 h-5 object-contain" />
                      <span class="text-white font-medium">Google Meet</span>
                    </div>
                  </div>
                  
                  <!-- CTA Button -->
                  <a href="#register" class="inline-flex items-center gap-3 bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-white font-bold px-8 py-4 rounded-xl shadow-lg shadow-amber-500/30 transition-all hover:shadow-xl hover:shadow-amber-500/40 hover:-translate-y-0.5 text-lg">
                    {{ block.data.button_text || 'Daftar Sekarang — Gratis' }}
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                    </svg>
                  </a>
                  
                  <!-- Trust badges -->
                  <div class="flex items-center gap-6 mt-8 text-blue-200/70 text-sm">
                      <div class="flex items-center gap-2">
                          <svg class="w-5 h-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                          </svg>
                          <span>100% Gratis</span>
                      </div>
                      <div class="flex items-center gap-2">
                          <svg class="w-5 h-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                          </svg>
                          <span>E-Sertifikat</span>
                      </div>
                  </div>
              </div>
              
              <!-- Right - Speaker Card + Countdown -->
              <div class="relative order-1 lg:order-2">
                  <!-- Speaker Card -->
                  <div class="relative bg-white/10 backdrop-blur-xl rounded-3xl border border-white/20 p-6 lg:p-8 shadow-2xl">
                      <!-- Speaker Image -->
                      <div class="relative mb-6">
                          <div class="aspect-[4/3] rounded-2xl overflow-hidden shadow-2xl border-4 border-white/20">
                              <img 
                                :src="speakerImage" 
                                :alt="speakerName"
                                class="w-full h-full object-cover object-top"
                              />
                          </div>
                          <!-- Speaker Badge -->
                          <div class="absolute -bottom-4 left-1/2 -translate-x-1/2 bg-gradient-to-r from-amber-500 to-amber-600 text-white px-6 py-2 rounded-full font-semibold shadow-lg text-sm whitespace-nowrap">
                              🎤 Narasumber
                          </div>
                      </div>
                      
                      <!-- Speaker Info -->
                      <div class="text-center mb-6 pt-2">
                          <h3 class="text-2xl font-bold text-white mb-1">{{ speakerName }}</h3>
                          <p class="text-amber-300 font-medium">{{ speakerTitle }}</p>
                      </div>
                      
                      <!-- Countdown -->
                      <div class="bg-white/5 rounded-2xl p-5 border border-white/10">
                          <p class="text-center text-xs font-semibold text-blue-200/70 uppercase tracking-wider mb-4">Webinar Dimulai Dalam</p>
                          <div class="grid grid-cols-4 gap-3">
                              <div class="text-center">
                                  <div class="bg-white/10 backdrop-blur rounded-xl p-3">
                                      <span class="text-2xl lg:text-3xl font-bold text-white">{{ countdown.days }}</span>
                                  </div>
                                  <span class="text-[10px] font-medium text-blue-200/60 uppercase mt-1 block">Hari</span>
                              </div>
                              <div class="text-center">
                                  <div class="bg-white/10 backdrop-blur rounded-xl p-3">
                                      <span class="text-2xl lg:text-3xl font-bold text-white">{{ countdown.hours }}</span>
                                  </div>
                                  <span class="text-[10px] font-medium text-blue-200/60 uppercase mt-1 block">Jam</span>
                              </div>
                              <div class="text-center">
                                  <div class="bg-white/10 backdrop-blur rounded-xl p-3">
                                      <span class="text-2xl lg:text-3xl font-bold text-white">{{ countdown.minutes }}</span>
                                  </div>
                                  <span class="text-[10px] font-medium text-blue-200/60 uppercase mt-1 block">Menit</span>
                              </div>
                              <div class="text-center">
                                  <div class="bg-amber-500/30 backdrop-blur rounded-xl p-3">
                                      <span class="text-2xl lg:text-3xl font-bold text-amber-300">{{ countdown.seconds }}</span>
                                  </div>
                                  <span class="text-[10px] font-medium text-blue-200/60 uppercase mt-1 block">Detik</span>
                              </div>
                          </div>
                      </div>
                  </div>
                  
                  <!-- Floating customs image -->
                  <div class="absolute -bottom-8 -right-8 w-24 h-24 bg-white/10 backdrop-blur-lg rounded-2xl border border-white/10 flex items-center justify-center text-5xl shadow-xl hidden lg:flex">
                      📦
                  </div>
              </div>
          </div>
      </div>
      
      <!-- Bottom wave decoration -->
      <div class="absolute bottom-0 left-0 right-0">
        <svg viewBox="0 0 1440 120" fill="none" xmlns="http://www.w3.org/2000/svg" class="w-full">
          <path d="M0 120L48 108C96 96 192 72 288 66C384 60 480 72 576 78C672 84 768 84 864 78C960 72 1056 60 1152 60C1248 60 1344 72 1392 78L1440 84V120H1392C1344 120 1248 120 1152 120C1056 120 960 120 864 120C768 120 672 120 576 120C480 120 384 120 288 120C192 120 96 120 48 120H0Z" fill="white"/>
        </svg>
      </div>
  </header>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');

.font-inter {
    font-family: 'Inter', sans-serif;
}
</style>
