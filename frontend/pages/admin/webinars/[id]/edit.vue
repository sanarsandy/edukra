<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="pending" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-admin-200 border-t-admin-600 rounded-full animate-spin"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error || !webinar" class="bg-red-50 text-red-600 p-4 rounded-lg">
      Webinar tidak ditemukan
    </div>

    <!-- Form -->
    <template v-else>
      <!-- Header -->
      <div class="flex items-center gap-4">
        <NuxtLink to="/admin/webinars" class="p-2 hover:bg-neutral-100 rounded-lg transition-colors">
          <svg class="w-5 h-5 text-neutral-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"/>
          </svg>
        </NuxtLink>
        <div>
          <h1 class="text-2xl font-bold text-neutral-900">Edit Webinar</h1>
          <p class="text-neutral-500 mt-1">{{ webinar.title }}</p>
        </div>
      </div>

      <form @submit.prevent="handleSubmit" class="bg-white rounded-xl border border-neutral-200 p-6 space-y-6">
        <!-- Course Info (readonly) -->
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-2">Kursus</label>
          <div class="px-4 py-2.5 bg-neutral-50 border border-neutral-200 rounded-lg text-neutral-600">
            {{ webinar.course?.title || 'N/A' }}
          </div>
        </div>

        <!-- Title -->
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-2">Judul Webinar <span class="text-red-500">*</span></label>
          <input 
            v-model="form.title"
            type="text" 
            class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
            required
          />
        </div>

        <!-- Description -->
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-2">Deskripsi</label>
          <textarea 
            v-model="form.description"
            rows="3"
            class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
          ></textarea>
        </div>

        <!-- Date & Time -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-neutral-700 mb-2">Tanggal & Waktu <span class="text-red-500">*</span></label>
            <input 
              v-model="form.scheduled_at"
              type="datetime-local"
              class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
              required
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-neutral-700 mb-2">Durasi (menit)</label>
            <input 
              v-model.number="form.duration_minutes"
              type="number"
              min="15"
              max="480"
              class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
            />
          </div>
        </div>

        <!-- Meeting URL & Password -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-neutral-700 mb-2">Link Meeting</label>
            <input 
              v-model="form.meeting_url"
              type="url"
              class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-neutral-700 mb-2">Password Meeting</label>
            <input 
              v-model="form.meeting_password"
              type="text"
              class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
            />
          </div>
        </div>

        <!-- Max Participants -->
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-2">Maksimal Peserta</label>
          <input 
            v-model.number="form.max_participants"
            type="number"
            min="1"
            class="w-full md:w-1/3 px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
          />
        </div>

        <!-- Recording URL (for completed webinars) -->
        <div v-if="webinar.status === 'completed'">
          <label class="block text-sm font-medium text-neutral-700 mb-2">Link Rekaman</label>
          <input 
            v-model="form.recording_url"
            type="url"
            class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
            placeholder="https://..."
          />
          <p class="text-xs text-neutral-500 mt-1">Link rekaman webinar untuk peserta yang tidak bisa hadir</p>
        </div>

        <!-- Status -->
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-2">Status</label>
          <div class="flex flex-wrap gap-4">
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="form.status" value="draft" class="text-admin-600" />
              <span class="text-sm">Draft</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="form.status" value="upcoming" class="text-admin-600" />
              <span class="text-sm">Upcoming</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="form.status" value="live" class="text-admin-600" />
              <span class="text-sm">🔴 Live</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="form.status" value="completed" class="text-admin-600" />
              <span class="text-sm">Selesai</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="form.status" value="cancelled" class="text-admin-600" />
              <span class="text-sm">Dibatalkan</span>
            </label>
          </div>
        </div>

        <!-- Submit -->
        <div class="flex items-center justify-end gap-4 pt-4 border-t border-neutral-100">
          <NuxtLink 
            to="/admin/webinars"
            class="px-4 py-2.5 text-sm font-medium text-neutral-700 bg-neutral-100 hover:bg-neutral-200 rounded-lg"
          >
            Batal
          </NuxtLink>
          <button 
            type="submit"
            :disabled="submitting"
            class="px-6 py-2.5 text-sm font-medium text-white bg-admin-600 hover:bg-admin-700 rounded-lg disabled:opacity-50"
          >
            {{ submitting ? 'Menyimpan...' : 'Simpan Perubahan' }}
          </button>
        </div>
      </form>
    </template>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const config = useRuntimeConfig()
const router = useRouter()
const route = useRoute()
const token = useCookie('token')

// Fetch webinar data
const { data: webinar, pending, error } = await useFetch(`/api/admin/webinars/${route.params.id}`, {
  baseURL: config.public.apiBase,
  headers: {
    Authorization: `Bearer ${token.value}`
  }
})

// Form state
const form = ref({
  title: '',
  description: '',
  scheduled_at: '',
  duration_minutes: 60,
  meeting_url: '',
  meeting_password: '',
  max_participants: null,
  recording_url: '',
  status: 'upcoming'
})

const submitting = ref(false)

// Populate form when data loads
watchEffect(() => {
  if (webinar.value) {
    const w = webinar.value
    form.value = {
      title: w.title || '',
      description: w.description || '',
      scheduled_at: w.scheduled_at ? formatDateTimeLocal(w.scheduled_at) : '',
      duration_minutes: w.duration_minutes || 60,
      meeting_url: w.meeting_url || '',
      meeting_password: w.meeting_password || '',
      max_participants: w.max_participants || null,
      recording_url: w.recording_url || '',
      status: w.status || 'upcoming'
    }
  }
})

function formatDateTimeLocal(isoString) {
  const date = new Date(isoString)
  const offset = date.getTimezoneOffset()
  const local = new Date(date.getTime() - offset * 60000)
  return local.toISOString().slice(0, 16)
}

async function handleSubmit() {
  submitting.value = true
  
  try {
    const scheduledAt = new Date(form.value.scheduled_at).toISOString()
    
    const payload = {
      title: form.value.title,
      description: form.value.description || null,
      scheduled_at: scheduledAt,
      duration_minutes: form.value.duration_minutes,
      meeting_url: form.value.meeting_url || null,
      meeting_password: form.value.meeting_password || null,
      max_participants: form.value.max_participants || null,
      recording_url: form.value.recording_url || null,
      status: form.value.status
    }
    
    await $fetch(`/api/admin/webinars/${route.params.id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token.value}`
      },
      body: payload
    })
    
    router.push('/admin/webinars')
  } catch (err) {
    console.error('Failed to update webinar:', err)
    alert('Gagal menyimpan perubahan: ' + (err.data?.error || err.message))
  } finally {
    submitting.value = false
  }
}
</script>
