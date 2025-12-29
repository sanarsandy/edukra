<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center gap-4">
      <NuxtLink to="/admin/webinars" class="p-2 hover:bg-neutral-100 rounded-lg transition-colors">
        <svg class="w-5 h-5 text-neutral-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"/>
        </svg>
      </NuxtLink>
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Buat Webinar Baru</h1>
        <p class="text-neutral-500 mt-1">Jadwalkan webinar untuk kursus Anda</p>
      </div>
    </div>

    <!-- Form -->
    <form @submit.prevent="handleSubmit" class="bg-white rounded-xl border border-neutral-200 p-6 space-y-6">
      <!-- Course Selection -->
      <div>
        <label class="block text-sm font-medium text-neutral-700 mb-2">Kursus <span class="text-red-500">*</span></label>
        <select 
          v-model="form.course_id"
          class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
          required
        >
          <option value="">Pilih Kursus</option>
          <option v-for="course in courses" :key="course.id" :value="course.id">
            {{ course.title }}
          </option>
        </select>
      </div>

      <!-- Title -->
      <div>
        <label class="block text-sm font-medium text-neutral-700 mb-2">Judul Webinar <span class="text-red-500">*</span></label>
        <input 
          v-model="form.title"
          type="text" 
          class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
          placeholder="Contoh: Sesi Live Q&A - Belajar React"
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
          placeholder="Deskripsi singkat tentang webinar ini..."
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
            placeholder="60"
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
            placeholder="https://zoom.us/j/xxxxx atau https://meet.google.com/xxx"
          />
          <p class="text-xs text-neutral-500 mt-1">Link Zoom, Google Meet, atau platform lainnya</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-2">Password Meeting</label>
          <input 
            v-model="form.meeting_password"
            type="text"
            class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
            placeholder="Optional"
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
          placeholder="Tidak terbatas"
        />
        <p class="text-xs text-neutral-500 mt-1">Kosongkan jika tidak ada batasan</p>
      </div>

      <!-- Status -->
      <div>
        <label class="block text-sm font-medium text-neutral-700 mb-2">Status</label>
        <div class="flex gap-4">
          <label class="flex items-center gap-2 cursor-pointer">
            <input type="radio" v-model="form.status" value="draft" class="text-admin-600 focus:ring-admin-500" />
            <span class="text-sm text-neutral-700">Draft</span>
          </label>
          <label class="flex items-center gap-2 cursor-pointer">
            <input type="radio" v-model="form.status" value="upcoming" class="text-admin-600 focus:ring-admin-500" />
            <span class="text-sm text-neutral-700">Upcoming (aktif)</span>
          </label>
        </div>
      </div>

      <!-- Submit -->
      <div class="flex items-center justify-end gap-4 pt-4 border-t border-neutral-100">
        <NuxtLink 
          to="/admin/webinars"
          class="px-4 py-2.5 text-sm font-medium text-neutral-700 bg-neutral-100 hover:bg-neutral-200 rounded-lg transition-colors"
        >
          Batal
        </NuxtLink>
        <button 
          type="submit"
          :disabled="submitting"
          class="px-6 py-2.5 text-sm font-medium text-white bg-admin-600 hover:bg-admin-700 rounded-lg transition-colors disabled:opacity-50"
        >
          {{ submitting ? 'Menyimpan...' : 'Simpan Webinar' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const config = useRuntimeConfig()
const router = useRouter()
const token = useCookie('token')

// Form state
const form = ref({
  course_id: '',
  title: '',
  description: '',
  scheduled_at: '',
  duration_minutes: 60,
  meeting_url: '',
  meeting_password: '',
  max_participants: null,
  status: 'upcoming'
})

const submitting = ref(false)

// Fetch courses for dropdown
const { data: coursesData } = await useFetch('/api/admin/courses', {
  baseURL: config.public.apiBase,
  headers: {
    Authorization: `Bearer ${token.value}`
  }
})

const courses = computed(() => coursesData.value?.courses || [])

// Handle form submission
async function handleSubmit() {
  submitting.value = true
  
  try {
    // Convert datetime-local to ISO 8601
    const scheduledAt = new Date(form.value.scheduled_at).toISOString()
    
    const payload = {
      ...form.value,
      scheduled_at: scheduledAt,
      description: form.value.description || null,
      meeting_url: form.value.meeting_url || null,
      meeting_password: form.value.meeting_password || null,
      max_participants: form.value.max_participants || null
    }
    
    await $fetch('/api/admin/webinars', {
      baseURL: config.public.apiBase,
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token.value}`
      },
      body: payload
    })
    
    router.push('/admin/webinars')
  } catch (err) {
    console.error('Failed to create webinar:', err)
    alert('Gagal membuat webinar: ' + (err.data?.error || err.message))
  } finally {
    submitting.value = false
  }
}
</script>
