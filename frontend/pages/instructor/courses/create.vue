<template>
  <div class="max-w-3xl mx-auto">
    <div class="mb-8">
      <NuxtLink to="/instructor/courses" class="inline-flex items-center gap-2 text-slate-500 hover:text-slate-700 text-sm font-medium mb-4 transition-colors">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 19l-7-7 7-7"/>
        </svg>
        Kembali ke Kursus
      </NuxtLink>
      <h1 class="text-2xl font-display font-bold text-slate-900">Buat Kursus Baru</h1>
      <p class="text-slate-600 mt-1">Isi informasi dasar kursus. Anda dapat menambahkan materi setelah kursus dibuat.</p>
    </div>

    <!-- Error Alert -->
    <div v-if="error" class="mb-6 p-4 bg-rose-50 border border-rose-200 rounded-xl">
      <div class="flex items-center gap-3">
        <svg class="w-5 h-5 text-rose-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <p class="text-rose-700 text-sm">{{ error }}</p>
      </div>
    </div>

    <form @submit.prevent="handleSubmit" class="bg-white rounded-2xl border border-slate-200/50 shadow-sm overflow-hidden">
      <div class="p-6 space-y-6">
        <!-- Title -->
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Judul Kursus <span class="text-rose-500">*</span></label>
          <input
            v-model="form.title"
            type="text"
            required
            placeholder="Contoh: Belajar JavaScript dari Nol"
            class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all"
          />
        </div>

        <!-- Description -->
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Deskripsi</label>
          <textarea
            v-model="form.description"
            rows="4"
            placeholder="Jelaskan tentang kursus ini, apa yang akan dipelajari..."
            class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all resize-none"
          ></textarea>
        </div>

        <!-- Category & Price Row -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Category -->
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">Kategori</label>
            <select
              v-model="form.category_id"
              class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all bg-white"
            >
              <option value="">Pilih kategori</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>

          <!-- Price -->
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">Harga (IDR)</label>
            <div class="relative">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-500">Rp</span>
              <input
                v-model.number="form.price"
                type="number"
                min="0"
                step="1000"
                placeholder="0"
                class="w-full pl-12 pr-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all"
              />
            </div>
            <p class="mt-1 text-xs text-slate-500">Kosongkan atau 0 untuk kursus gratis</p>
          </div>
        </div>

        <!-- Thumbnail -->
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Thumbnail</label>
          <div class="border-2 border-dashed border-slate-200 rounded-xl p-6 text-center hover:border-teal-400 transition-colors">
            <div v-if="thumbnailPreview" class="relative inline-block">
              <img :src="thumbnailPreview" class="max-h-40 rounded-lg" alt="Preview" />
              <button @click="clearThumbnail" type="button" class="absolute -top-2 -right-2 w-6 h-6 bg-rose-500 text-white rounded-full flex items-center justify-center shadow-lg hover:bg-rose-600 transition-colors">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
            <div v-else>
              <svg class="w-10 h-10 mx-auto text-slate-300 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
              </svg>
              <label class="cursor-pointer">
                <span class="text-teal-600 hover:text-teal-700 font-medium">Upload gambar</span>
                <span class="text-slate-500"> atau drag & drop</span>
                <input type="file" accept="image/*" @change="handleThumbnail" class="hidden" />
              </label>
              <p class="text-xs text-slate-400 mt-1">PNG, JPG hingga 2MB</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="px-6 py-4 bg-slate-50 border-t border-slate-100 flex items-center justify-end gap-3">
        <NuxtLink to="/instructor/courses" class="px-5 py-2.5 text-slate-600 font-medium hover:text-slate-800 transition-colors">
          Batal
        </NuxtLink>
        <button
          type="submit"
          :disabled="loading || !form.title"
          class="px-6 py-2.5 bg-gradient-to-r from-teal-500 to-cyan-500 text-white font-semibold rounded-xl shadow-lg shadow-teal-500/25 hover:shadow-xl hover:shadow-teal-500/30 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          <svg v-if="loading" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ loading ? 'Menyimpan...' : 'Buat Kursus' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'instructor',
  middleware: ['instructor']
})

const instructorPanel = useInstructorPanel()

const loading = ref(false)
const error = ref<string | null>(null)
const categories = ref<any[]>([])

const form = ref({
  title: '',
  description: '',
  category_id: '',
  price: 0,
  thumbnail_url: ''
})

const thumbnailPreview = ref<string | null>(null)
const thumbnailFile = ref<File | null>(null)

const handleThumbnail = (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    if (file.size > 2 * 1024 * 1024) {
      error.value = 'Ukuran file maksimal 2MB'
      return
    }
    thumbnailFile.value = file
    thumbnailPreview.value = URL.createObjectURL(file)
  }
}

const clearThumbnail = () => {
  thumbnailFile.value = null
  thumbnailPreview.value = null
  form.value.thumbnail_url = ''
}

const handleSubmit = async () => {
  if (!form.value.title.trim()) {
    error.value = 'Judul kursus wajib diisi'
    return
  }

  loading.value = true
  error.value = null

  // Upload thumbnail first if exists
  if (thumbnailFile.value) {
    const uploadResult = await instructorPanel.uploadFile(thumbnailFile.value, 'thumbnail')
    if (uploadResult?.url) {
      form.value.thumbnail_url = uploadResult.url
    }
  }

  const result = await instructorPanel.createCourse({
    title: form.value.title,
    description: form.value.description,
    category_id: form.value.category_id || undefined,
    price: form.value.price || 0,
    thumbnail_url: form.value.thumbnail_url || undefined
  })

  if (result?.id) {
    navigateTo(`/instructor/courses/${result.id}`)
  } else {
    error.value = instructorPanel.error.value || 'Gagal membuat kursus'
  }

  loading.value = false
}

onMounted(async () => {
  const data = await instructorPanel.fetchCategories()
  if (data?.categories) {
    categories.value = data.categories
  }
})
</script>


