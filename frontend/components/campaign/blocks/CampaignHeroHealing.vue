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
  return dateStr ? new Date(dateStr).getTime() : new Date().getTime() + 86400000
})

const countdown = ref({ days: '00', hours: '00', minutes: '00', seconds: '00' })
const displayDate = computed(() => {
  const dateStr = props.globalEndDate || props.block.data.end_date
  if (!dateStr) return '28 Januari 2025'
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' }).format(date)
})
const displayTime = computed(() => {
  const dateStr = props.globalEndDate || props.block.data.end_date
  if (!dateStr) return '19:00 - 20:00 WIB'
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

const scrollToForm = () => {
    const formEl = document.querySelector('#cta_healing') || document.querySelector('#cta_clean') || document.querySelector('form')
    if (formEl) formEl.scrollIntoView({ behavior: 'smooth' })
}
</script>

<template>
  <section class="healing-hero relative overflow-hidden">
    <!-- Watercolor Background -->
    <div class="absolute inset-0 bg-gradient-to-b from-[#f5f0e8] via-[#e8e4dc] to-[#d4cfc4]"></div>
    
    <!-- Decorative Elements -->
    <div class="absolute inset-0 opacity-30">
      <div class="absolute top-0 left-0 w-96 h-96 bg-gradient-radial from-amber-200/30 to-transparent rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 right-0 w-80 h-80 bg-gradient-radial from-sky-200/30 to-transparent rounded-full blur-3xl"></div>
    </div>

    <div class="relative max-w-6xl mx-auto px-4 py-16 md:py-24">
      <!-- Title Section -->
      <div class="text-center mb-12">
        <!-- Main Title with Elegant Typography -->
        <h1 class="font-serif text-5xl md:text-7xl font-light text-[#4a4540] mb-4 tracking-wide">
          {{ block.data.headline || 'P.U.L.I.H' }}
        </h1>
        <p class="text-xl md:text-2xl text-[#6b635a] font-light italic mb-2">
          {{ block.data.subheadline || 'Seni Berdamai dengan Luka Batin' }}
        </p>
      </div>

      <!-- Hero Image - Meditation Figure -->
      <div class="relative mx-auto max-w-md mb-12">
        <div class="aspect-square rounded-full overflow-hidden shadow-2xl relative">
          <img 
            :src="block.data.image_url || 'https://images.unsplash.com/photo-1506126613408-eca07ce68773?auto=format&fit=crop&q=80&w=800'" 
            :alt="block.data.headline || 'P.U.L.I.H Webinar'"
            class="object-cover w-full h-full"
          />
          <!-- Golden Light Overlay -->
          <div class="absolute inset-0 bg-gradient-to-t from-amber-500/10 via-transparent to-transparent"></div>
        </div>
        <!-- Floating Glow Effect -->
        <div class="absolute -inset-4 bg-gradient-radial from-amber-300/20 to-transparent rounded-full blur-2xl -z-10"></div>
      </div>

      <!-- Event Badge -->
      <div class="text-center mb-8">
        <span class="inline-block bg-[#4a4540]/10 text-[#4a4540] font-semibold text-sm px-6 py-2 rounded-full border border-[#4a4540]/20">
          {{ block.data.badge || 'FREE WEBINAR' }}
        </span>
      </div>

      <!-- Event Details -->
      <div class="flex flex-wrap justify-center gap-6 mb-10 text-[#5a534a]">
        <div class="flex items-center gap-2 bg-white/50 backdrop-blur-sm px-4 py-2 rounded-full">
          <span class="text-lg">📅</span>
          <span class="font-medium">{{ displayDate }}</span>
        </div>
        <div class="flex items-center gap-2 bg-white/50 backdrop-blur-sm px-4 py-2 rounded-full">
          <span class="text-lg">⏰</span>
          <span class="font-medium">{{ displayTime }}</span>
        </div>
        <div class="flex items-center gap-2 bg-white/50 backdrop-blur-sm px-4 py-2 rounded-full">
          <span class="text-lg">💻</span>
          <span class="font-medium">{{ block.data.platform || 'Online via Zoom' }}</span>
        </div>
      </div>

      <!-- Countdown Box -->
      <div class="max-w-lg mx-auto mb-10">
        <div class="bg-white/60 backdrop-blur-md rounded-2xl p-6 shadow-lg border border-white/50">
          <p class="text-xs font-semibold text-[#8a8279] uppercase mb-4 tracking-widest text-center">
            Pendaftaran Ditutup Dalam
          </p>
          <div class="flex justify-center items-center gap-4">
            <div class="text-center">
              <div class="w-16 h-16 bg-gradient-to-br from-[#4a4540] to-[#6b635a] rounded-xl flex items-center justify-center shadow-md">
                <span class="text-2xl font-bold text-white">{{ countdown.days }}</span>
              </div>
              <span class="text-[10px] text-[#8a8279] uppercase font-medium mt-2 block">Hari</span>
            </div>
            <span class="text-2xl text-[#c4bcb0] font-light">:</span>
            <div class="text-center">
              <div class="w-16 h-16 bg-gradient-to-br from-[#4a4540] to-[#6b635a] rounded-xl flex items-center justify-center shadow-md">
                <span class="text-2xl font-bold text-white">{{ countdown.hours }}</span>
              </div>
              <span class="text-[10px] text-[#8a8279] uppercase font-medium mt-2 block">Jam</span>
            </div>
            <span class="text-2xl text-[#c4bcb0] font-light">:</span>
            <div class="text-center">
              <div class="w-16 h-16 bg-gradient-to-br from-[#4a4540] to-[#6b635a] rounded-xl flex items-center justify-center shadow-md">
                <span class="text-2xl font-bold text-white">{{ countdown.minutes }}</span>
              </div>
              <span class="text-[10px] text-[#8a8279] uppercase font-medium mt-2 block">Menit</span>
            </div>
            <span class="text-2xl text-[#c4bcb0] font-light">:</span>
            <div class="text-center">
              <div class="w-16 h-16 bg-gradient-to-br from-[#4a4540] to-[#6b635a] rounded-xl flex items-center justify-center shadow-md">
                <span class="text-2xl font-bold text-white">{{ countdown.seconds }}</span>
              </div>
              <span class="text-[10px] text-[#8a8279] uppercase font-medium mt-2 block">Detik</span>
            </div>
          </div>
        </div>
      </div>

      <!-- CTA Button -->
      <div class="text-center">
        <button 
          @click="scrollToForm"
          class="group relative inline-flex items-center gap-3 bg-gradient-to-r from-[#4a4540] to-[#6b635a] text-white font-semibold py-4 px-10 rounded-full shadow-xl hover:shadow-2xl transition-all duration-300 hover:scale-105"
          :style="block.data.button_color ? { background: block.data.button_color } : {}"
        >
          <span class="text-lg">{{ block.data.button_text || 'DAFTAR GRATIS SEKARANG' }}</span>
          <svg class="w-5 h-5 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
          </svg>
        </button>
        <p class="text-sm text-[#8a8279] mt-4 flex items-center justify-center gap-2">
          <span class="text-amber-600">✨</span>
          {{ block.data.cta_note || 'Kuota terbatas, amankan tempatmu sekarang!' }}
        </p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.healing-hero {
  font-family: 'Inter', 'Segoe UI', sans-serif;
}

.font-serif {
  font-family: 'Playfair Display', 'Georgia', serif;
}

.bg-gradient-radial {
  background: radial-gradient(circle, var(--tw-gradient-from) 0%, var(--tw-gradient-to) 70%);
}
</style>
