<template>
  <div class="max-w-4xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <NuxtLink to="/instructor/courses" class="inline-flex items-center gap-2 text-slate-500 hover:text-slate-700 text-sm font-medium mb-4 transition-colors">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 19l-7-7 7-7"/>
        </svg>
        Kembali
      </NuxtLink>
      <div class="flex items-start justify-between gap-4">
        <div>
          <h1 class="text-2xl font-display font-bold text-slate-900">{{ course?.title || 'Loading...' }}</h1>
          <div class="flex items-center gap-3 mt-2">
            <span :class="getStatusClass(course?.status)" class="px-2.5 py-1 text-xs font-semibold rounded-full">{{ getStatusLabel(course?.status) }}</span>
            <span class="text-sm text-slate-500">{{ course?.lessons_count || 0 }} materi</span>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <NuxtLink :to="`/instructor/courses/${courseId}/materials`" class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium rounded-lg transition-colors text-sm flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
            </svg>
            Kelola Materi
          </NuxtLink>
          <button 
            v-if="course?.status === 'draft' || course?.status === 'rejected'"
            @click="submitForReview"
            :disabled="submitting"
            class="px-4 py-2 bg-gradient-to-r from-teal-500 to-cyan-500 text-white font-medium rounded-lg shadow-md hover:shadow-lg transition-all text-sm disabled:opacity-50 flex items-center gap-2"
          >
            <svg v-if="submitting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
            </svg>
            Submit Review
          </button>
        </div>
      </div>
    </div>

    <!-- Success/Error Alerts -->
    <div v-if="successMessage" class="mb-6 p-4 bg-emerald-50 border border-emerald-200 rounded-xl">
      <div class="flex items-center gap-3">
        <svg class="w-5 h-5 text-emerald-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <p class="text-emerald-700 text-sm">{{ successMessage }}</p>
      </div>
    </div>

    <div v-if="error" class="mb-6 p-4 bg-rose-50 border border-rose-200 rounded-xl">
      <div class="flex items-center gap-3">
        <svg class="w-5 h-5 text-rose-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <p class="text-rose-700 text-sm">{{ error }}</p>
      </div>
    </div>

    <!-- Review Notes (if rejected) -->
    <div v-if="course?.status === 'rejected' && course?.review_notes" class="mb-6 p-4 bg-amber-50 border border-amber-200 rounded-xl">
      <div class="flex gap-3">
        <svg class="w-5 h-5 text-amber-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
        </svg>
        <div>
          <p class="font-medium text-amber-800 mb-1">Catatan dari Admin</p>
          <p class="text-amber-700 text-sm">{{ course.review_notes }}</p>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="bg-white rounded-2xl border border-slate-200/50 p-12 text-center">
      <div class="inline-flex items-center gap-2 text-slate-500">
        <svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
        </svg>
        Memuat...
      </div>
    </div>

    <!-- Edit Form -->
    <form v-else @submit.prevent="handleSave" class="bg-white rounded-2xl border border-slate-200/50 shadow-sm overflow-hidden">
      <div class="p-6 space-y-6">
        <!-- Title -->
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Judul Kursus <span class="text-rose-500">*</span></label>
          <input
            v-model="form.title"
            type="text"
            required
            :disabled="course?.status === 'published'"
            class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all disabled:bg-slate-50 disabled:cursor-not-allowed"
          />
        </div>

        <!-- Description -->
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Deskripsi</label>
          <textarea
            v-model="form.description"
            rows="4"
            :disabled="course?.status === 'published'"
            class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all resize-none disabled:bg-slate-50 disabled:cursor-not-allowed"
          ></textarea>
        </div>

        <!-- Category & Price -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">Kategori</label>
            <select
              v-model="form.category_id"
              :disabled="course?.status === 'published'"
              class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all bg-white disabled:bg-slate-50 disabled:cursor-not-allowed"
            >
              <option value="">Pilih kategori</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">Harga (IDR)</label>
            <div class="relative">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-500">Rp</span>
              <input
                v-model.number="form.price"
                type="number"
                min="0"
                :disabled="course?.status === 'published'"
                class="w-full pl-12 pr-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all disabled:bg-slate-50 disabled:cursor-not-allowed"
              />
            </div>
          </div>
        </div>

        <!-- Published Notice -->
        <div v-if="course?.status === 'published'" class="p-4 bg-slate-50 rounded-xl">
          <div class="flex items-center gap-3">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <p class="text-sm text-slate-600">Kursus yang sudah dipublish tidak dapat diedit. Hubungi admin untuk unpublish.</p>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div v-if="course?.status !== 'published'" class="px-6 py-4 bg-slate-50 border-t border-slate-100 flex items-center justify-between">
        <button
          v-if="course?.status === 'draft' || course?.status === 'rejected'"
          type="button"
          @click="deleteCourse"
          :disabled="deleting"
          class="text-rose-600 hover:text-rose-700 font-medium text-sm transition-colors flex items-center gap-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
          </svg>
          {{ deleting ? 'Menghapus...' : 'Hapus Kursus' }}
        </button>
        <div class="flex items-center gap-3">
          <button
            type="submit"
            :disabled="saving"
            class="px-6 py-2.5 bg-gradient-to-r from-teal-500 to-cyan-500 text-white font-semibold rounded-xl shadow-lg shadow-teal-500/25 hover:shadow-xl transition-all disabled:opacity-50 flex items-center gap-2"
          >
            <svg v-if="saving" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
            </svg>
            {{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'instructor',
  middleware: ['instructor']
})

const route = useRoute()
const courseId = route.params.id as string
const instructorPanel = useInstructorPanel()

const loading = ref(true)
const saving = ref(false)
const deleting = ref(false)
const submitting = ref(false)
const error = ref<string | null>(null)
const successMessage = ref<string | null>(null)

const course = ref<any>(null)
const categories = ref<any[]>([])

const form = ref({
  title: '',
  description: '',
  category_id: '',
  price: 0
})

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    draft: 'bg-slate-100 text-slate-700',
    pending_review: 'bg-amber-100 text-amber-700',
    approved: 'bg-blue-100 text-blue-700',
    rejected: 'bg-rose-100 text-rose-700',
    published: 'bg-emerald-100 text-emerald-700'
  }
  return classes[status] || 'bg-slate-100 text-slate-700'
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    draft: 'Draft',
    pending_review: 'Menunggu Review',
    approved: 'Disetujui',
    rejected: 'Ditolak',
    published: 'Dipublish'
  }
  return labels[status] || status
}

const fetchCourse = async () => {
  loading.value = true
  const data = await instructorPanel.fetchCourse(courseId)
  if (data) {
    course.value = data
    form.value = {
      title: data.title || '',
      description: data.description || '',
      category_id: data.category_id || '',
      price: data.price || 0
    }
  }
  loading.value = false
}

const handleSave = async () => {
  saving.value = true
  error.value = null
  successMessage.value = null

  const result = await instructorPanel.updateCourse(courseId, {
    title: form.value.title,
    description: form.value.description,
    category_id: form.value.category_id || undefined,
    price: form.value.price
  })

  if (result) {
    successMessage.value = 'Kursus berhasil disimpan'
    await fetchCourse()
  } else {
    error.value = instructorPanel.error.value || 'Gagal menyimpan'
  }

  saving.value = false
}

const submitForReview = async () => {
  submitting.value = true
  error.value = null

  const result = await instructorPanel.submitCourseForReview(courseId)
  if (result) {
    successMessage.value = result.message
    await fetchCourse()
  } else {
    error.value = instructorPanel.error.value || 'Gagal submit'
  }

  submitting.value = false
}

const deleteCourse = async () => {
  if (!confirm('Apakah Anda yakin ingin menghapus kursus ini?')) return

  deleting.value = true
  const result = await instructorPanel.deleteCourse(courseId)
  if (result) {
    navigateTo('/instructor/courses')
  } else {
    error.value = instructorPanel.error.value || 'Gagal menghapus'
  }
  deleting.value = false
}

onMounted(async () => {
  const [_, catData] = await Promise.all([
    fetchCourse(),
    instructorPanel.fetchCategories()
  ])
  if (catData?.categories) {
    categories.value = catData.categories
  }
})
</script>


