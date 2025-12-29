<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-neutral-900">Webinar Saya</h1>
      <p class="text-neutral-500 mt-1">Daftar webinar yang telah Anda daftar</p>
    </div>

    <!-- Loading -->
    <div v-if="pending" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-primary-200 border-t-primary-600 rounded-full animate-spin"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="bg-red-50 text-red-600 p-4 rounded-lg">
      Gagal memuat data webinar
    </div>

    <!-- Empty State -->
    <div v-else-if="webinars.length === 0" class="bg-white rounded-xl border border-neutral-200 p-12 text-center">
      <div class="w-20 h-20 bg-gradient-to-br from-purple-100 to-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-10 h-10 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/>
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-neutral-900 mb-2">Belum ada Webinar</h3>
      <p class="text-neutral-500 max-w-sm mx-auto">Anda belum terdaftar di webinar manapun. Daftar kursus untuk ikuti webinar live.</p>
    </div>

    <!-- Webinar List -->
    <div v-else class="space-y-4">
      <!-- Upcoming Section -->
      <div v-if="upcomingWebinars.length > 0">
        <h2 class="text-lg font-semibold text-neutral-900 mb-3">Webinar Mendatang</h2>
        <div class="grid gap-4">
          <div 
            v-for="webinar in upcomingWebinars" 
            :key="webinar.id"
            class="bg-white border border-neutral-200 rounded-xl p-5 hover:shadow-md transition-all"
          >
            <div class="flex flex-col md:flex-row md:items-center gap-4">
              <!-- Date Card -->
              <div class="flex-shrink-0 w-20 h-20 bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl flex flex-col items-center justify-center text-white">
                <span class="text-2xl font-bold">{{ getDay(webinar.scheduled_at) }}</span>
                <span class="text-xs uppercase">{{ getMonth(webinar.scheduled_at) }}</span>
              </div>

              <!-- Info -->
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <h3 class="text-lg font-semibold text-neutral-900">{{ webinar.title }}</h3>
                  <span v-if="webinar.status === 'live'" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-700">
                    🔴 LIVE
                  </span>
                  <span v-else-if="isToday(webinar.scheduled_at)" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-700">
                    Hari Ini
                  </span>
                </div>
                
                <p v-if="webinar.course" class="text-sm text-neutral-500 mb-2">
                  {{ webinar.course.title }}
                </p>

                <div class="flex flex-wrap items-center gap-4 text-sm text-neutral-600">
                  <div class="flex items-center gap-1.5">
                    <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                    <span>{{ formatTime(webinar.scheduled_at) }} WIB</span>
                  </div>
                  <div class="flex items-center gap-1.5">
                    <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4h4"/>
                    </svg>
                    <span>{{ webinar.duration_minutes }} menit</span>
                  </div>
                  <div class="text-purple-600 font-medium">
                    {{ getCountdown(webinar.scheduled_at) }}
                  </div>
                </div>
              </div>

              <!-- Action -->
              <div class="flex-shrink-0">
                <a 
                  v-if="webinar.meeting_url && canJoin(webinar)"
                  :href="webinar.meeting_url"
                  target="_blank"
                  class="inline-flex items-center px-4 py-2.5 bg-purple-600 hover:bg-purple-700 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/>
                  </svg>
                  Join Webinar
                </a>
                <div v-else class="text-sm text-neutral-500 text-center">
                  <div class="px-4 py-2 bg-neutral-100 rounded-lg">
                    Belum dimulai
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Past Webinars Section -->
      <div v-if="pastWebinars.length > 0" class="mt-8">
        <h2 class="text-lg font-semibold text-neutral-900 mb-3">Webinar Sebelumnya</h2>
        <div class="grid gap-3">
          <div 
            v-for="webinar in pastWebinars" 
            :key="webinar.id"
            class="bg-white border border-neutral-200 rounded-lg p-4 opacity-75"
          >
            <div class="flex items-center justify-between">
              <div>
                <h3 class="font-medium text-neutral-900">{{ webinar.title }}</h3>
                <p class="text-sm text-neutral-500">{{ formatFullDate(webinar.scheduled_at) }}</p>
              </div>
              <a 
                v-if="webinar.recording_url"
                :href="webinar.recording_url"
                target="_blank"
                class="inline-flex items-center px-3 py-1.5 bg-neutral-100 hover:bg-neutral-200 text-neutral-700 text-sm font-medium rounded-lg"
              >
                <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                Lihat Rekaman
              </a>
              <span v-else class="text-xs text-neutral-400">Selesai</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'dashboard'
})

const config = useRuntimeConfig()
const token = useCookie('token')

// Fetch user's webinars
const { data, pending, error } = await useFetch('/api/my/webinars?include_completed=true', {
  baseURL: config.public.apiBase,
  headers: {
    Authorization: `Bearer ${token.value}`
  }
})

const webinars = computed(() => data.value || [])

const upcomingWebinars = computed(() => 
  webinars.value.filter(w => w.status === 'upcoming' || w.status === 'live')
)

const pastWebinars = computed(() => 
  webinars.value.filter(w => w.status === 'completed')
)

// Helper functions
function getDay(dateStr) {
  return new Date(dateStr).getDate()
}

function getMonth(dateStr) {
  return new Date(dateStr).toLocaleDateString('id-ID', { month: 'short' })
}

function formatTime(dateStr) {
  return new Date(dateStr).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function formatFullDate(dateStr) {
  return new Date(dateStr).toLocaleDateString('id-ID', { 
    weekday: 'long', 
    day: 'numeric', 
    month: 'long', 
    year: 'numeric' 
  })
}

function isToday(dateStr) {
  const webinarDate = new Date(dateStr).toDateString()
  const today = new Date().toDateString()
  return webinarDate === today
}

function canJoin(webinar) {
  // Can join 15 minutes before and during the webinar
  const now = new Date()
  const scheduled = new Date(webinar.scheduled_at)
  const joinWindow = new Date(scheduled.getTime() - 15 * 60 * 1000) // 15 min before
  const endWindow = new Date(scheduled.getTime() + webinar.duration_minutes * 60 * 1000)
  
  return now >= joinWindow && now <= endWindow
}

function getCountdown(dateStr) {
  const now = new Date()
  const scheduled = new Date(dateStr)
  const diff = scheduled - now
  
  if (diff <= 0) return 'Sedang berlangsung'
  
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  
  if (days > 0) return `${days} hari ${hours} jam lagi`
  if (hours > 0) return `${hours} jam ${minutes} menit lagi`
  return `${minutes} menit lagi`
}
</script>
