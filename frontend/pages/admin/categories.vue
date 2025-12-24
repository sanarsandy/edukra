<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Kategori Kursus</h1>
        <p class="text-neutral-500 mt-1">Kelola kategori untuk mengorganisir kursus</p>
      </div>
      <button @click="openAddModal" class="btn-admin w-full sm:w-auto">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
        </svg>
        Tambah Kategori
      </button>
    </div>

    <!-- Categories Grid -->
    <div class="grid md:grid-cols-2 xl:grid-cols-3 gap-6">
      <div v-for="category in categories" :key="category.id" class="bg-white rounded-xl border border-neutral-200 overflow-hidden hover:shadow-lg transition-all group">
        <div :class="['h-32 flex items-center justify-center relative', category.color || getColorClass(categories.indexOf(category))]">
          <!-- If icon is SVG path (starts with M), render as SVG -->
          <svg v-if="isSvgPath(category.icon)" class="w-12 h-12 text-white/50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" :d="category.icon"/>
          </svg>
          <!-- If icon is emoji or text, render as text -->
          <span v-else class="text-5xl">{{ category.icon || '📚' }}</span>
          <div class="absolute top-3 right-3 flex gap-1 transition-opacity">
            <button @click="openEditModal(category)" class="p-1.5 bg-white/20 hover:bg-white/30 rounded-lg text-white transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
              </svg>
            </button>
            <button @click="openDeleteModal(category)" class="p-1.5 bg-white/20 hover:bg-red-500 rounded-lg text-white transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </button>
          </div>
        </div>
        <div class="p-5">
          <h3 class="font-semibold text-neutral-900 mb-1">{{ category.name }}</h3>
          <p class="text-sm text-neutral-500 mb-4">{{ category.description }}</p>
          <div class="flex items-center justify-between">
            <span class="text-sm text-neutral-600">
              <span class="font-semibold text-neutral-900">{{ category.course_count || 0 }}</span> Kursus
            </span>
            <span class="text-xs text-neutral-400">{{ (category.students || 0).toLocaleString() }} Siswa</span>
          </div>
        </div>
      </div>

      <!-- Add Category Card -->
      <div 
        @click="openAddModal"
        class="bg-neutral-50 rounded-xl border-2 border-dashed border-neutral-300 flex flex-col items-center justify-center p-8 cursor-pointer hover:border-admin-500 hover:bg-admin-50 transition-all group"
      >
        <div class="w-16 h-16 bg-neutral-200 group-hover:bg-admin-100 rounded-full flex items-center justify-center mb-4 transition-colors">
          <svg class="w-8 h-8 text-neutral-400 group-hover:text-admin-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
        </div>
        <p class="text-sm font-medium text-neutral-500 group-hover:text-admin-600 transition-colors">Tambah Kategori Baru</p>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <Transition name="fade">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="closeModal">
        <div class="bg-white rounded-xl p-6 w-full max-w-md m-4" @click.stop>
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-neutral-900">{{ isEditing ? 'Edit Kategori' : 'Tambah Kategori' }}</h3>
            <button @click="closeModal" class="p-2 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          
          <form @submit.prevent="saveCategory" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Nama Kategori <span class="text-red-500">*</span></label>
              <input 
                v-model="form.name"
                type="text" 
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="Contoh: Mobile Development"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Deskripsi</label>
              <textarea 
                v-model="form.description"
                rows="3" 
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm resize-none"
                placeholder="Deskripsi singkat kategori"
              ></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Warna</label>
              <div class="flex gap-2 flex-wrap">
                <button 
                  v-for="color in colorOptions" 
                  :key="color.value"
                  type="button"
                  @click="form.color = color.value"
                  class="w-8 h-8 rounded-lg transition-all"
                  :class="[color.class, form.color === color.value ? 'ring-2 ring-offset-2 ring-neutral-900' : '']"
                ></button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Ikon</label>
              <div class="grid grid-cols-6 gap-2">
                <button 
                  v-for="icon in iconOptions" 
                  :key="icon.name"
                  type="button"
                  @click="form.icon = icon.path"
                  class="p-3 rounded-lg border transition-all flex items-center justify-center"
                  :class="form.icon === icon.path ? 'border-admin-500 bg-admin-50 text-admin-600' : 'border-neutral-200 text-neutral-400 hover:border-neutral-300'"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" :d="icon.path"/>
                  </svg>
                </button>
              </div>
            </div>
            
            <div class="flex gap-3 pt-4">
              <button type="button" @click="closeModal" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
                Batal
              </button>
              <button type="submit" class="flex-1 py-2.5 text-sm font-medium text-white bg-admin-600 rounded-lg hover:bg-admin-700 transition-colors">
                {{ isEditing ? 'Simpan Perubahan' : 'Tambah Kategori' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- Delete Modal -->
    <Transition name="fade">
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showDeleteModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-sm m-4" @click.stop>
          <div class="text-center">
            <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-neutral-900 mb-2">Hapus Kategori</h3>
            <p class="text-sm text-neutral-500 mb-6">Apakah Anda yakin ingin menghapus kategori <strong>{{ selectedCategory?.name }}</strong>?</p>
          </div>
          <div class="flex gap-3">
            <button @click="showDeleteModal = false" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">Batal</button>
            <button @click="deleteCategory" class="flex-1 py-2.5 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors">Hapus</button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Toast -->
    <Transition name="slide-up">
      <div v-if="toast.show" class="fixed bottom-6 right-6 z-50">
        <div class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 bg-accent-600 text-white">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
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
  title: 'Kategori Kursus - Admin'
})

// API Composable
const { 
  loading, 
  error, 
  categories, 
  total, 
  fetchCategories, 
  createCategory: createCategoryApi, 
  updateCategory: updateCategoryApi, 
  deleteCategory: deleteCategoryApi 
} = useCategories()

const showModal = ref(false)
const showDeleteModal = ref(false)
const isEditing = ref(false)
const selectedCategory = ref<any>(null)
const saving = ref(false)
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

const colorOptions = [
  { value: 'bg-primary-600', class: 'bg-primary-600' },
  { value: 'bg-accent-600', class: 'bg-accent-600' },
  { value: 'bg-warm-500', class: 'bg-warm-500' },
  { value: 'bg-cyan-600', class: 'bg-cyan-600' },
  { value: 'bg-rose-500', class: 'bg-rose-500' },
  { value: 'bg-indigo-600', class: 'bg-indigo-600' }
]

const iconOptions = [
  { name: 'code', path: 'M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4' },
  { name: 'design', path: 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z' },
  { name: 'chart', path: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z' },
  { name: 'mobile', path: 'M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z' },
  { name: 'business', path: 'M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z' },
  { name: 'marketing', path: 'M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z' },
  { name: 'video', path: 'M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z' },
  { name: 'book', path: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253' },
  { name: 'music', path: 'M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3' },
  { name: 'camera', path: 'M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z M15 13a3 3 0 11-6 0 3 3 0 016 0z' },
  { name: 'globe', path: 'M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9' },
  { name: 'health', path: 'M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z' }
]

const form = ref({
  name: '',
  description: '',
  icon: '',
  color: 'bg-primary-600'
})

// Load categories on mount
onMounted(async () => {
  await fetchCategories()
})

const getColorClass = (index: number) => {
  const colors = ['bg-primary-600', 'bg-accent-600', 'bg-warm-500', 'bg-cyan-600', 'bg-rose-500', 'bg-indigo-600']
  return colors[index % colors.length]
}

// Check if icon string is an SVG path (starts with M or m for moveto command)
const isSvgPath = (icon: string | null | undefined): boolean => {
  if (!icon || typeof icon !== 'string') return false
  return /^[Mm]\s*[\d.-]/.test(icon.trim())
}

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const openAddModal = () => {
  isEditing.value = false
  form.value = { name: '', description: '', icon: iconOptions[0].path, color: 'bg-primary-600' }
  showModal.value = true
}

const openEditModal = (category: any) => {
  isEditing.value = true
  selectedCategory.value = category
  form.value = { 
    name: category.name, 
    description: category.description || '', 
    icon: category.icon || iconOptions[0].path,
    color: category.color || 'bg-primary-600'
  }
  showModal.value = true
}

const openDeleteModal = (category: any) => {
  selectedCategory.value = category
  showDeleteModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveCategory = async () => {
  saving.value = true
  try {
    if (isEditing.value && selectedCategory.value) {
      const result = await updateCategoryApi(selectedCategory.value.id, form.value)
      if (result) {
        showToast('Kategori berhasil diperbarui')
        await fetchCategories()
      } else {
        showToast('Gagal memperbarui kategori', 'error')
      }
    } else {
      const result = await createCategoryApi(form.value)
      if (result) {
        showToast('Kategori berhasil ditambahkan')
        await fetchCategories()
      } else {
        showToast('Gagal menambahkan kategori', 'error')
      }
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    closeModal()
  }
}

const deleteCategory = async () => {
  if (!selectedCategory.value) return
  
  saving.value = true
  try {
    const success = await deleteCategoryApi(selectedCategory.value.id)
    if (success) {
      showToast('Kategori berhasil dihapus')
      await fetchCategories()
    } else {
      showToast('Gagal menghapus kategori', 'error')
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    showDeleteModal.value = false
    selectedCategory.value = null
  }
}
</script>



<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { opacity: 0; transform: translateY(20px); }
</style>
