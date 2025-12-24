<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Kursus Saya</h1>
        <p class="text-neutral-500 mt-1">Kelola semua kursus yang Anda buat</p>
      </div>
      <div class="flex gap-3">
        <button @click="openAddModal" class="btn-admin">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
          Tambah Kursus
        </button>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Total Kursus</p>
        <p class="text-2xl font-bold text-neutral-900">{{ courses.length }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Dipublikasi</p>
        <p class="text-2xl font-bold text-accent-600">{{ publishedCount }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Menunggu Review</p>
        <p class="text-2xl font-bold text-warm-600">{{ pendingCount }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Draft</p>
        <p class="text-2xl font-bold text-neutral-500">{{ draftCount }}</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-col sm:flex-row gap-4 mb-6">
      <div class="relative flex-1">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari kursus..." 
          class="w-full pl-10 pr-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
        />
      </div>
      <select v-model="categoryFilter" class="px-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm">
        <option value="">Semua Kategori</option>
        <option v-for="cat in categoriesList" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
      </select>
      <select v-model="statusFilter" class="px-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm">
        <option value="">Semua Status</option>
        <option value="published">Published</option>
        <option value="pending_review">Menunggu Review</option>
        <option value="draft">Draft</option>
      </select>
    </div>

    <!-- Courses Table -->
    <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-neutral-50 border-b border-neutral-200">
            <tr>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Kursus</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Kategori</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Harga</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Status</th>
              <th class="px-6 py-4 text-right text-xs font-semibold text-neutral-600 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-neutral-100">
            <tr v-if="loading" class="hover:bg-neutral-50">
              <td colspan="5" class="px-6 py-8 text-center text-neutral-500">
                <div class="flex items-center justify-center gap-2">
                  <svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Memuat...
                </div>
              </td>
            </tr>
            <tr v-else-if="filteredCourses.length === 0" class="hover:bg-neutral-50">
              <td colspan="5" class="px-6 py-8 text-center text-neutral-500">
                Tidak ada kursus ditemukan
              </td>
            </tr>
            <tr v-for="course in filteredCourses" :key="course.id" class="hover:bg-neutral-50 transition-colors">
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div 
                    v-if="course.thumbnail_url" 
                    class="w-12 h-12 rounded-lg overflow-hidden flex-shrink-0"
                  >
                    <img :src="course.thumbnail_url" :alt="course.title" class="w-full h-full object-cover" />
                  </div>
                  <div v-else :class="['w-12 h-12 rounded-lg flex items-center justify-center flex-shrink-0', getCourseColor(course.id)]">
                    <svg class="w-6 h-6 text-white/70" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                    </svg>
                  </div>
                  <div>
                    <p class="text-sm font-medium text-neutral-900">{{ course.title }}</p>
                    <p class="text-xs text-neutral-500">{{ course.description?.substring(0, 50) || 'No description' }}...</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-neutral-600">{{ course.category?.name || '-' }}</td>
              <td class="px-6 py-4 text-sm text-neutral-900 font-medium">{{ formatCurrency(course.price) }}</td>
              <td class="px-6 py-4">
                <span 
                  class="px-2.5 py-1 text-xs font-medium rounded-full"
                  :class="getStatusClass(course.status)"
                >
                  {{ getStatusLabel(course.status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <NuxtLink :to="`/instructor/courses/${course.id}/materials`" class="p-2 text-neutral-400 hover:text-accent-600 hover:bg-accent-50 rounded-lg transition-colors" title="Kelola Materi">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                    </svg>
                  </NuxtLink>
                  <button @click="openEditModal(course)" class="p-2 text-neutral-400 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
                    </svg>
                  </button>
                  <button 
                    v-if="course.status === 'draft'"
                    @click="submitForReview(course)" 
                    class="p-2 text-neutral-400 hover:text-accent-600 hover:bg-accent-50 rounded-lg transition-colors" 
                    title="Ajukan Review"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                  </button>
                  <button @click="openDeleteModal(course)" class="p-2 text-neutral-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add/Edit Course Modal -->
    <Transition name="fade">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="closeModal">
        <div class="bg-white rounded-xl p-6 w-full max-w-lg m-4 max-h-[90vh] overflow-y-auto" @click.stop>
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-neutral-900">{{ isEditing ? 'Edit Kursus' : 'Tambah Kursus' }}</h3>
            <button @click="closeModal" class="p-2 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          
          <form @submit.prevent="saveCourse" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Judul Kursus <span class="text-red-500">*</span></label>
              <input 
                v-model="form.title"
                type="text" 
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="Masukkan judul kursus"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Deskripsi</label>
              <textarea 
                v-model="form.description"
                rows="3"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="Masukkan deskripsi kursus"
              ></textarea>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Kategori <span class="text-red-500">*</span></label>
                <select 
                  v-model="form.category_id"
                  required
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                >
                  <option value="">Pilih kategori</option>
                  <option v-for="cat in categoriesList" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Harga (Rp) <span class="text-red-500">*</span></label>
                <input 
                  v-model.number="form.price"
                  type="number" 
                  required
                  min="0"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                  placeholder="Contoh: 499000"
                />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Thumbnail URL</label>
              <input 
                v-model="form.thumbnail_url"
                type="url"
                placeholder="https://example.com/image.jpg"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
              />
            </div>
            
            <div class="bg-amber-50 border border-amber-200 rounded-lg p-4">
              <p class="text-sm text-amber-800">
                <strong>Catatan:</strong> Kursus akan disimpan sebagai <strong>Draft</strong>. Setelah selesai, ajukan untuk direview oleh admin.
              </p>
            </div>
            
            <div class="flex gap-3 pt-4">
              <button type="button" @click="closeModal" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
                Batal
              </button>
              <button type="submit" :disabled="saving" class="flex-1 py-2.5 text-sm font-medium text-white bg-admin-600 rounded-lg hover:bg-admin-700 transition-colors disabled:opacity-50">
                {{ isEditing ? 'Simpan Perubahan' : 'Tambah Kursus' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- Delete Confirmation Modal -->
    <Transition name="fade">
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showDeleteModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-sm m-4" @click.stop>
          <div class="text-center">
            <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-neutral-900 mb-2">Hapus Kursus</h3>
            <p class="text-sm text-neutral-500 mb-6">Apakah Anda yakin ingin menghapus <strong>{{ selectedCourse?.title }}</strong>? Semua data kursus akan dihapus permanen.</p>
          </div>
          <div class="flex gap-3">
            <button @click="showDeleteModal = false" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
              Batal
            </button>
            <button @click="confirmDeleteCourse" :disabled="saving" class="flex-1 py-2.5 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50">
              Hapus
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Submit for Review Modal -->
    <Transition name="fade">
      <div v-if="showReviewModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showReviewModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-sm m-4" @click.stop>
          <div class="text-center">
            <div class="w-16 h-16 bg-accent-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-neutral-900 mb-2">Ajukan Review</h3>
            <p class="text-sm text-neutral-500 mb-6">Apakah Anda yakin ingin mengajukan <strong>{{ selectedCourse?.title }}</strong> untuk direview? Admin akan mereview kursus Anda.</p>
          </div>
          <div class="flex gap-3">
            <button @click="showReviewModal = false" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
              Batal
            </button>
            <button @click="confirmSubmitReview" :disabled="saving" class="flex-1 py-2.5 text-sm font-medium text-white bg-accent-600 rounded-lg hover:bg-accent-700 transition-colors disabled:opacity-50">
              Ajukan
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Toast -->
    <Transition name="slide-up">
      <div v-if="toast.show" class="fixed bottom-6 right-6 z-50">
        <div class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3" :class="toast.type === 'success' ? 'bg-accent-600 text-white' : 'bg-red-600 text-white'">
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
  layout: 'instructor',
  middleware: 'instructor'
})

useHead({
  title: 'Kursus Saya - Instructor'
})

const instructorPanel = useInstructorPanel()
const { categories: categoriesList, fetchCategories } = useCategories()

const loading = ref(true)
const courses = ref<any[]>([])
const searchQuery = ref('')
const categoryFilter = ref('')
const statusFilter = ref('')
const showModal = ref(false)
const showDeleteModal = ref(false)
const showReviewModal = ref(false)
const isEditing = ref(false)
const selectedCourse = ref<any>(null)
const saving = ref(false)

const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

const form = ref({
  title: '',
  description: '',
  price: 0,
  category_id: '',
  thumbnail_url: ''
})

// Load courses on mount
onMounted(async () => {
  await loadCourses()
})

const loadCourses = async () => {
  loading.value = true
  try {
    await fetchCategories()
    const data = await instructorPanel.fetchCourses()
    if (data?.courses) {
      courses.value = data.courses
    }
  } finally {
    loading.value = false
  }
}

const publishedCount = computed(() => courses.value.filter(c => c.status === 'published').length)
const pendingCount = computed(() => courses.value.filter(c => c.status === 'pending_review').length)
const draftCount = computed(() => courses.value.filter(c => c.status === 'draft').length)

const filteredCourses = computed(() => {
  return courses.value.filter(course => {
    const matchesSearch = course.title.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = !statusFilter.value || course.status === statusFilter.value
    const matchesCategory = !categoryFilter.value || course.category_id === categoryFilter.value
    return matchesSearch && matchesStatus && matchesCategory
  })
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const formatCurrency = (amount: number) => {
  return `Rp ${(amount || 0).toLocaleString('id-ID')}`
}

const getCourseColor = (id: string) => {
  const colors = ['bg-primary-600', 'bg-accent-600', 'bg-warm-500', 'bg-cyan-600', 'bg-rose-500', 'bg-indigo-600']
  const hash = (id || '').split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  return colors[hash % colors.length]
}

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    'draft': 'bg-neutral-100 text-neutral-700',
    'pending_review': 'bg-warm-100 text-warm-700',
    'published': 'bg-accent-100 text-accent-700',
    'rejected': 'bg-red-100 text-red-700'
  }
  return classes[status] || classes.draft
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    'draft': 'Draft',
    'pending_review': 'Menunggu Review',
    'published': 'Published',
    'rejected': 'Ditolak'
  }
  return labels[status] || status
}

const openAddModal = () => {
  isEditing.value = false
  form.value = { 
    title: '', 
    description: '', 
    price: 0, 
    category_id: '',
    thumbnail_url: ''
  }
  showModal.value = true
}

const openEditModal = (course: any) => {
  isEditing.value = true
  selectedCourse.value = course
  form.value = { 
    title: course.title,
    description: course.description || '',
    price: course.price,
    category_id: course.category_id || '',
    thumbnail_url: course.thumbnail_url || ''
  }
  showModal.value = true
}

const openDeleteModal = (course: any) => {
  selectedCourse.value = course
  showDeleteModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const submitForReview = (course: any) => {
  selectedCourse.value = course
  showReviewModal.value = true
}

const saveCourse = async () => {
  saving.value = true
  try {
    const courseData = {
      title: form.value.title,
      description: form.value.description,
      price: form.value.price,
      category_id: form.value.category_id,
      thumbnail_url: form.value.thumbnail_url || null,
    }

    if (isEditing.value && selectedCourse.value) {
      const result = await instructorPanel.updateCourse(selectedCourse.value.id, courseData)
      if (result) {
        showToast('Kursus berhasil diperbarui')
        await loadCourses()
      } else {
        showToast('Gagal memperbarui kursus', 'error')
      }
    } else {
      const result = await instructorPanel.createCourse(courseData)
      if (result) {
        showToast('Kursus berhasil ditambahkan')
        await loadCourses()
      } else {
        showToast('Gagal menambahkan kursus', 'error')
      }
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    closeModal()
  }
}

const confirmDeleteCourse = async () => {
  if (!selectedCourse.value) return
  
  saving.value = true
  try {
    const success = await instructorPanel.deleteCourse(selectedCourse.value.id)
    if (success) {
      showToast('Kursus berhasil dihapus')
      await loadCourses()
    } else {
      showToast('Gagal menghapus kursus', 'error')
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    showDeleteModal.value = false
    selectedCourse.value = null
  }
}

const confirmSubmitReview = async () => {
  if (!selectedCourse.value) return
  
  saving.value = true
  try {
    const success = await instructorPanel.submitForReview(selectedCourse.value.id)
    if (success) {
      showToast('Kursus berhasil diajukan untuk review')
      await loadCourses()
    } else {
      showToast('Gagal mengajukan review', 'error')
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    showReviewModal.value = false
    selectedCourse.value = null
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { opacity: 0; transform: translateY(20px); }
</style>
