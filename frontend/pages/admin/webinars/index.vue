<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Webinar Management</h1>
        <p class="text-neutral-500 mt-1">Kelola webinar untuk kursus Anda</p>
      </div>
      <NuxtLink 
        to="/admin/webinars/create" 
        class="inline-flex items-center justify-center px-4 py-2.5 bg-admin-600 hover:bg-admin-700 text-white text-sm font-medium rounded-lg transition-colors shadow-sm"
      >
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4.5v15m7.5-7.5h-15"/></svg>
        Tambah Webinar
      </NuxtLink>
    </div>

    <!-- Loading State -->
    <div v-if="pending" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-admin-200 border-t-admin-600 rounded-full animate-spin"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 text-red-600 p-4 rounded-lg">
      <p class="font-medium">Gagal memuat data webinar</p>
      <button @click="refresh" class="text-sm underline mt-1">Coba lagi</button>
    </div>

    <!-- Empty State -->
    <div v-else-if="webinars.length === 0" class="bg-white rounded-xl border border-neutral-200 p-12 text-center">
      <div class="w-20 h-20 bg-gradient-to-br from-purple-100 to-admin-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-10 h-10 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/>
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-neutral-900 mb-2">Belum ada Webinar</h3>
      <p class="text-neutral-500 max-w-sm mx-auto mb-6">Buat webinar untuk kursus Anda dan mulai mengadakan sesi live.</p>
      <NuxtLink to="/admin/webinars/create" class="inline-flex items-center px-4 py-2 bg-admin-600 text-white rounded-lg hover:bg-admin-700 text-sm font-medium">
        Buat Webinar Pertama →
      </NuxtLink>
    </div>

    <!-- Webinar List -->
    <div v-else class="space-y-4">
      <div 
        v-for="webinar in webinars" 
        :key="webinar.id"
        class="bg-white border border-neutral-200 rounded-xl p-5 hover:shadow-md transition-all"
      >
        <div class="flex items-start gap-4">
          <!-- Icon -->
          <div class="w-14 h-14 rounded-lg flex items-center justify-center flex-shrink-0" 
               :class="getStatusBgColor(webinar.status)">
            <svg class="w-7 h-7" :class="getStatusIconColor(webinar.status)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/>
            </svg>
          </div>

          <!-- Info -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <h3 class="text-lg font-semibold text-neutral-900 truncate">{{ webinar.title }}</h3>
              <span 
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                :class="getStatusClass(webinar.status)"
              >
                {{ getStatusLabel(webinar.status) }}
              </span>
            </div>

            <!-- Course -->
            <p v-if="webinar.course" class="text-sm text-neutral-500 mb-2">
              <span class="text-neutral-400">Kursus:</span> {{ webinar.course.title }}
            </p>

            <!-- Schedule Info -->
            <div class="flex flex-wrap items-center gap-4 text-sm mb-3">
              <div class="flex items-center gap-1.5 text-neutral-600">
                <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
                <span>{{ formatDate(webinar.scheduled_at) }}</span>
              </div>
              <div class="flex items-center gap-1.5 text-neutral-600">
                <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <span>{{ formatTime(webinar.scheduled_at) }} WIB</span>
              </div>
              <div class="flex items-center gap-1.5 text-neutral-600">
                <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4h4m-4 0l-.707.707M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <span>{{ webinar.duration_minutes }} menit</span>
              </div>
            </div>

            <!-- Stats -->
            <div class="flex items-center gap-4 text-sm">
              <div class="flex items-center gap-1.5 text-purple-600 font-medium">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
                </svg>
                <span>{{ webinar.registrations_count || 0 }} peserta</span>
              </div>
              <a v-if="webinar.meeting_url" :href="webinar.meeting_url" target="_blank" class="text-admin-600 hover:underline text-xs">
                Buka Link Meeting →
              </a>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-2">
            <NuxtLink 
              :to="`/admin/webinars/${webinar.id}/registrations`"
              class="p-2 text-neutral-400 hover:text-purple-600 hover:bg-purple-50 rounded-lg transition-colors"
              title="Lihat Peserta"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
              </svg>
            </NuxtLink>
            <NuxtLink 
              :to="`/admin/webinars/${webinar.id}/edit`"
              class="p-2 text-neutral-400 hover:text-admin-600 hover:bg-admin-50 rounded-lg transition-colors"
              title="Edit Webinar"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
              </svg>
            </NuxtLink>
            <button 
              @click="confirmDelete(webinar)"
              class="p-2 text-neutral-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
              title="Hapus Webinar"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="deleteModal.show" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="fixed inset-0 bg-black/50" @click="deleteModal.show = false"></div>
      <div class="relative bg-white rounded-xl p-6 max-w-md w-full mx-4 shadow-xl">
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">Hapus Webinar?</h3>
        <p class="text-neutral-600 text-sm mb-4">
          Apakah Anda yakin ingin menghapus webinar <strong>{{ deleteModal.webinar?.title }}</strong>? 
          Tindakan ini tidak dapat dibatalkan.
        </p>
        <div class="flex gap-3 justify-end">
          <button 
            @click="deleteModal.show = false"
            class="px-4 py-2 text-sm font-medium text-neutral-700 bg-neutral-100 hover:bg-neutral-200 rounded-lg"
          >
            Batal
          </button>
          <button 
            @click="deleteWebinar"
            :disabled="deleting"
            class="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-lg disabled:opacity-50"
          >
            {{ deleting ? 'Menghapus...' : 'Hapus Webinar' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const config = useRuntimeConfig()
const token = useCookie('token')

// Fetch webinars
const { data, pending, error, refresh } = await useFetch('/api/admin/webinars', {
  baseURL: config.public.apiBase,
  headers: {
    Authorization: `Bearer ${token.value}`
  }
})

const webinars = computed(() => data.value?.webinars || [])

// Delete modal state
const deleteModal = ref({ show: false, webinar: null })
const deleting = ref(false)

function confirmDelete(webinar) {
  deleteModal.value = { show: true, webinar }
}

async function deleteWebinar() {
  if (!deleteModal.value.webinar) return
  
  deleting.value = true
  try {
    await $fetch(`/api/admin/webinars/${deleteModal.value.webinar.id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token.value}`
      }
    })
    deleteModal.value.show = false
    refresh()
  } catch (err) {
    console.error('Failed to delete webinar:', err)
    alert('Gagal menghapus webinar')
  } finally {
    deleting.value = false
  }
}

// Helper functions
function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

function formatTime(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getStatusClass(status) {
  const classes = {
    draft: 'bg-neutral-100 text-neutral-600',
    upcoming: 'bg-blue-100 text-blue-700',
    live: 'bg-red-100 text-red-700',
    completed: 'bg-green-100 text-green-700',
    cancelled: 'bg-orange-100 text-orange-700'
  }
  return classes[status] || classes.draft
}

function getStatusLabel(status) {
  const labels = {
    draft: 'Draft',
    upcoming: 'Upcoming',
    live: '🔴 LIVE',
    completed: 'Selesai',
    cancelled: 'Dibatalkan'
  }
  return labels[status] || status
}

function getStatusBgColor(status) {
  const colors = {
    draft: 'bg-neutral-100',
    upcoming: 'bg-blue-100',
    live: 'bg-red-100',
    completed: 'bg-green-100',
    cancelled: 'bg-orange-100'
  }
  return colors[status] || colors.draft
}

function getStatusIconColor(status) {
  const colors = {
    draft: 'text-neutral-500',
    upcoming: 'text-blue-600',
    live: 'text-red-600',
    completed: 'text-green-600',
    cancelled: 'text-orange-600'
  }
  return colors[status] || colors.draft
}
</script>
