<template>
  <div>
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-neutral-900">AI Content Processing</h1>
      <p class="text-neutral-500 mt-1">Proses konten kursus untuk AI Tutor</p>
    </div>

    <!-- AI Status Check -->
    <div v-if="!aiEnabled" class="bg-amber-50 border border-amber-200 rounded-xl p-4 mb-6 flex items-center gap-3">
      <svg class="w-5 h-5 text-amber-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
      </svg>
      <div>
        <p class="text-sm font-medium text-amber-800">AI Tutor belum aktif</p>
        <p class="text-xs text-amber-600">Aktifkan AI Tutor di <NuxtLink to="/admin/settings" class="underline">Settings</NuxtLink> terlebih dahulu</p>
      </div>
    </div>

    <!-- Courses List -->
    <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
      <div class="p-4 border-b border-neutral-100 flex items-center justify-between">
        <h3 class="font-semibold text-neutral-900">Daftar Kursus</h3>
        <button 
          @click="processAllCourses" 
          :disabled="processingAll || !aiEnabled"
          class="px-4 py-2 text-sm font-medium bg-admin-600 text-white rounded-lg hover:bg-admin-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          <svg v-if="processingAll" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span>Proses Semua Kursus</span>
        </button>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="p-8 text-center text-neutral-500">
        <div class="animate-spin w-8 h-8 border-4 border-admin-500 border-t-transparent rounded-full mx-auto mb-3"></div>
        <p>Memuat daftar kursus...</p>
      </div>

      <!-- Course List -->
      <div v-else class="divide-y divide-neutral-100">
        <div 
          v-for="course in courses" 
          :key="course.id"
          class="p-4 flex items-center justify-between hover:bg-neutral-50 transition-colors"
        >
          <div class="flex items-center gap-4">
            <!-- Thumbnail -->
            <div class="w-16 h-12 bg-neutral-100 rounded-lg overflow-hidden flex-shrink-0">
              <img v-if="course.thumbnail_url" :src="getFullUrl(course.thumbnail_url)" class="w-full h-full object-cover" />
              <div v-else class="w-full h-full flex items-center justify-center text-neutral-400">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
              </div>
            </div>
            
            <!-- Info -->
            <div>
              <h4 class="font-medium text-neutral-900">{{ course.title }}</h4>
              <p class="text-xs text-neutral-500">{{ course.lessons_count || 0 }} materi</p>
            </div>
          </div>

          <!-- Status & Actions -->
          <div class="flex items-center gap-4">
            <!-- Processing Status -->
            <div class="text-sm">
              <template v-if="processingStatus[course.id]?.status === 'processing'">
                <span class="flex items-center gap-2 text-blue-600">
                  <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Memproses...
                </span>
              </template>
              <template v-else-if="processingStatus[course.id]?.status === 'completed'">
                <span class="flex items-center gap-2 text-accent-600">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                  {{ processingStatus[course.id]?.embedding_count || 0 }} embeddings
                </span>
              </template>
              <template v-else-if="processingStatus[course.id]?.status === 'failed'">
                <span class="flex items-center gap-2 text-red-600">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                  </svg>
                  Gagal
                </span>
              </template>
              <template v-else>
                <span class="text-neutral-400">Belum diproses</span>
              </template>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-2">
              <button 
                @click="processCourse(course.id)"
                :disabled="processingStatus[course.id]?.status === 'processing' || !aiEnabled"
                class="px-3 py-1.5 text-xs font-medium bg-admin-100 text-admin-700 rounded-lg hover:bg-admin-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ processingStatus[course.id]?.status === 'completed' ? 'Proses Ulang' : 'Proses' }}
              </button>
              <button 
                v-if="processingStatus[course.id]?.embedding_count > 0"
                @click="clearEmbeddings(course.id)"
                :disabled="processingStatus[course.id]?.status === 'processing'"
                class="px-3 py-1.5 text-xs font-medium bg-red-100 text-red-700 rounded-lg hover:bg-red-200 transition-colors disabled:opacity-50"
              >
                Hapus
              </button>
            </div>
          </div>
        </div>

        <div v-if="courses.length === 0" class="p-8 text-center text-neutral-500">
          Tidak ada kursus tersedia
        </div>
      </div>
    </div>

    <!-- Info Card -->
    <div class="mt-6 bg-blue-50 rounded-xl p-4 border border-blue-100">
      <h4 class="font-medium text-blue-900 mb-2">ℹ️ Cara Kerja AI Content Processing</h4>
      <ul class="text-sm text-blue-800 space-y-1">
        <li>• Konten text dari materi akan dipecah menjadi chunks</li>
        <li>• Setiap chunk dikonversi menjadi vector embeddings</li>
        <li>• AI Tutor akan mencari konteks relevan dari embeddings saat menjawab</li>
        <li>• Proses ulang jika Anda mengupdate materi kursus</li>
      </ul>
    </div>

    <!-- Toast -->
    <Transition name="slide-up">
      <div v-if="toast.show" class="fixed bottom-6 right-6 z-50">
        <div 
          class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3"
          :class="toast.type === 'success' ? 'bg-accent-600 text-white' : 'bg-red-600 text-white'"
        >
          <svg v-if="toast.type === 'success'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
          <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
          <span class="text-sm font-medium">{{ toast.message }}</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

useHead({
  title: 'AI Content Processing - Admin'
})

const config = useRuntimeConfig()
const apiBase = config.public.apiBase

interface Course {
  id: string
  title: string
  thumbnail_url?: string
  lessons_count?: number
}

interface ProcessingStatus {
  status: 'pending' | 'processing' | 'completed' | 'failed'
  embedding_count?: number
  error?: string
}

const loading = ref(true)
const courses = ref<Course[]>([])
const processingStatus = ref<Record<string, ProcessingStatus>>({})
const aiEnabled = ref(false)
const processingAll = ref(false)
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

const getFullUrl = (url: string) => {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return `${apiBase}${url}`
}

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const fetchCourses = async () => {
  try {
    const token = useCookie('token')
    const data = await $fetch<{ courses: Course[] }>(`${apiBase}/api/admin/courses`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    courses.value = data.courses || []
    
    // Fetch status for each course
    for (const course of courses.value) {
      await fetchProcessingStatus(course.id)
    }
  } catch (err) {
    console.error('Failed to fetch courses:', err)
  } finally {
    loading.value = false
  }
}

const fetchAIStatus = async () => {
  try {
    const token = useCookie('token')
    const data = await $fetch<{ enabled: boolean }>(`${apiBase}/api/admin/ai/settings`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    aiEnabled.value = data.enabled
  } catch (err) {
    console.error('Failed to fetch AI status:', err)
  }
}

const fetchProcessingStatus = async (courseId: string) => {
  try {
    const token = useCookie('token')
    const data = await $fetch<ProcessingStatus>(`${apiBase}/api/admin/courses/${courseId}/ai-processing-status`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    processingStatus.value[courseId] = data
  } catch (err) {
    // No status yet
    processingStatus.value[courseId] = { status: 'pending' }
  }
}

const processCourse = async (courseId: string) => {
  processingStatus.value[courseId] = { status: 'processing' }
  
  try {
    const token = useCookie('token')
    await $fetch(`${apiBase}/api/admin/courses/${courseId}/process-ai`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    
    // Poll for completion
    pollProcessingStatus(courseId)
    showToast('Proses dimulai...')
  } catch (err: any) {
    processingStatus.value[courseId] = { status: 'failed', error: err.message }
    showToast('Gagal memproses: ' + (err.data?.error || err.message), 'error')
  }
}

const pollProcessingStatus = async (courseId: string) => {
  let pollCount = 0
  const maxPollsToIgnoreCompleted = 5 // Ignore 'completed' for first 5 polls (10 seconds)
  
  const interval = setInterval(async () => {
    pollCount++
    const currentStatus = processingStatus.value[courseId]?.status
    
    try {
      const token = useCookie('token')
      const data = await $fetch<ProcessingStatus>(`${apiBase}/api/admin/courses/${courseId}/ai-processing-status`, {
        headers: { 'Authorization': `Bearer ${token.value}` }
      })
      
      // If we just started processing and server returns 'completed', 
      // keep showing 'processing' for first few polls (wait for new processing to take effect)
      if (currentStatus === 'processing' && data.status === 'completed' && pollCount <= maxPollsToIgnoreCompleted) {
        console.log(`[Poll ${pollCount}] Ignoring 'completed', waiting for processing to start...`)
        return // Keep current 'processing' state
      }
      
      // Update status from server
      processingStatus.value[courseId] = data
      
      // Stop polling if not processing
      if (data.status !== 'processing') {
        clearInterval(interval)
        
        if (data.status === 'completed') {
          showToast(`Berhasil! ${data.embedding_count || 0} embeddings dibuat`)
        }
      }
    } catch (err) {
      console.error('Poll error:', err)
    }
  }, 2000)
  
  // Stop after 5 minutes
  setTimeout(() => clearInterval(interval), 300000)
}

const clearEmbeddings = async (courseId: string) => {
  if (!confirm('Yakin ingin menghapus embeddings untuk kursus ini?')) return
  
  try {
    const token = useCookie('token')
    await $fetch(`${apiBase}/api/admin/courses/${courseId}/embeddings`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    
    processingStatus.value[courseId] = { status: 'pending' }
    showToast('Embeddings berhasil dihapus')
  } catch (err) {
    showToast('Gagal menghapus embeddings', 'error')
  }
}

const processAllCourses = async () => {
  if (!confirm('Proses semua kursus? Ini mungkin memakan waktu lama.')) return
  
  processingAll.value = true
  
  for (const course of courses.value) {
    await processCourse(course.id)
    // Wait a bit between courses
    await new Promise(resolve => setTimeout(resolve, 1000))
  }
  
  processingAll.value = false
}

onMounted(async () => {
  await fetchAIStatus()
  await fetchCourses()
})
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
