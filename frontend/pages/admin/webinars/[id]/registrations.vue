<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <NuxtLink 
          to="/admin/webinars" 
          class="inline-flex items-center text-sm text-neutral-500 hover:text-neutral-700 mb-2"
        >
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 19l-7-7 7-7"/>
          </svg>
          Kembali ke Webinars
        </NuxtLink>
        <h1 class="text-2xl font-bold text-neutral-900">Peserta Webinar</h1>
        <p v-if="webinar" class="text-neutral-500 mt-1">{{ webinar.title }}</p>
      </div>
      
      <!-- Stats -->
      <div v-if="registrations" class="flex gap-4">
        <div class="bg-purple-50 px-4 py-2 rounded-lg">
          <span class="text-2xl font-bold text-purple-600">{{ registrations.length }}</span>
          <span class="text-sm text-purple-600 ml-1">Terdaftar</span>
        </div>
        <div class="bg-green-50 px-4 py-2 rounded-lg">
          <span class="text-2xl font-bold text-green-600">{{ attendedCount }}</span>
          <span class="text-sm text-green-600 ml-1">Hadir</span>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="pending" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-admin-200 border-t-admin-600 rounded-full animate-spin"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="bg-red-50 text-red-600 p-4 rounded-lg">
      <p class="font-medium">Gagal memuat data peserta</p>
      <button @click="refresh" class="text-sm underline mt-1">Coba lagi</button>
    </div>

    <!-- Empty -->
    <div v-else-if="registrations.length === 0" class="bg-white rounded-xl border border-neutral-200 p-12 text-center">
      <div class="w-20 h-20 bg-neutral-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-10 h-10 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-neutral-900 mb-2">Belum ada Peserta</h3>
      <p class="text-neutral-500">Belum ada pengguna yang terdaftar untuk webinar ini.</p>
    </div>

    <!-- Registrations List -->
    <div v-else class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
      <table class="w-full">
        <thead class="bg-neutral-50 border-b border-neutral-200">
          <tr>
            <th class="text-left text-xs font-semibold text-neutral-500 uppercase tracking-wider px-6 py-3">Peserta</th>
            <th class="text-left text-xs font-semibold text-neutral-500 uppercase tracking-wider px-6 py-3">Kontak</th>
            <th class="text-left text-xs font-semibold text-neutral-500 uppercase tracking-wider px-6 py-3">Terdaftar</th>
            <th class="text-left text-xs font-semibold text-neutral-500 uppercase tracking-wider px-6 py-3">Sumber</th>
            <th class="text-center text-xs font-semibold text-neutral-500 uppercase tracking-wider px-6 py-3">Kehadiran</th>
            <th class="text-center text-xs font-semibold text-neutral-500 uppercase tracking-wider px-6 py-3">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-neutral-100">
          <tr v-for="reg in registrations" :key="reg.id" class="hover:bg-neutral-50">
            <!-- Peserta -->
            <td class="px-6 py-4">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-purple-100 flex items-center justify-center text-purple-600 font-semibold">
                  {{ reg.user?.full_name?.charAt(0)?.toUpperCase() || '?' }}
                </div>
                <div>
                  <p class="font-medium text-neutral-900">{{ reg.user?.full_name || 'Unknown' }}</p>
                  <p class="text-sm text-neutral-500">{{ reg.user?.email || '-' }}</p>
                </div>
              </div>
            </td>
            
            <!-- Kontak -->
            <td class="px-6 py-4">
              <p class="text-sm text-neutral-600">{{ reg.user?.phone || '-' }}</p>
            </td>
            
            <!-- Terdaftar -->
            <td class="px-6 py-4">
              <p class="text-sm text-neutral-600">{{ formatDate(reg.registered_at) }}</p>
            </td>
            
            <!-- Sumber -->
            <td class="px-6 py-4">
              <span 
                class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                :class="getSourceClass(reg.registration_source)"
              >
                {{ getSourceLabel(reg.registration_source) }}
              </span>
            </td>
            
            <!-- Kehadiran -->
            <td class="px-6 py-4 text-center">
              <span v-if="reg.attended" class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-700">
                ✓ Hadir
              </span>
              <span v-else class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-neutral-100 text-neutral-500">
                Belum
              </span>
            </td>
            
            <!-- Aksi -->
            <td class="px-6 py-4 text-center">
              <button 
                v-if="!reg.attended"
                @click="markAttendance(reg)"
                :disabled="markingAttendance === reg.id"
                class="text-sm text-admin-600 hover:text-admin-700 font-medium disabled:opacity-50"
              >
                {{ markingAttendance === reg.id ? 'Loading...' : 'Tandai Hadir' }}
              </button>
              <span v-else class="text-sm text-neutral-400">
                {{ formatDate(reg.attended_at) }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const config = useRuntimeConfig()
const route = useRoute()
const token = useCookie('token')

const webinarId = route.params.id

// Fetch webinar details
const { data: webinar } = await useFetch(`/api/admin/webinars/${webinarId}`, {
  baseURL: config.public.apiBase,
  headers: {
    Authorization: `Bearer ${token.value}`
  }
})

// Fetch registrations
const { data, pending, error, refresh } = await useFetch(`/api/admin/webinars/${webinarId}/registrations`, {
  baseURL: config.public.apiBase,
  headers: {
    Authorization: `Bearer ${token.value}`
  }
})

const registrations = computed(() => data.value?.registrations || [])
const attendedCount = computed(() => registrations.value.filter(r => r.attended).length)

const markingAttendance = ref(null)

async function markAttendance(reg) {
  markingAttendance.value = reg.id
  try {
    await $fetch(`/api/admin/webinars/${webinarId}/attendance/${reg.user_id}`, {
      baseURL: config.public.apiBase,
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token.value}`
      }
    })
    refresh()
  } catch (err) {
    console.error('Failed to mark attendance:', err)
    alert('Gagal menandai kehadiran')
  } finally {
    markingAttendance.value = null
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getSourceClass(source) {
  const classes = {
    purchase: 'bg-green-100 text-green-700',
    campaign: 'bg-blue-100 text-blue-700',
    direct: 'bg-purple-100 text-purple-700',
    admin: 'bg-orange-100 text-orange-700'
  }
  return classes[source] || 'bg-neutral-100 text-neutral-600'
}

function getSourceLabel(source) {
  const labels = {
    purchase: 'Pembelian',
    campaign: 'Campaign',
    direct: 'Langsung',
    admin: 'Admin'
  }
  return labels[source] || source || 'Unknown'
}
</script>
