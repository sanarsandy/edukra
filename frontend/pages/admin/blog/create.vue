<template>
  <div>
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div class="flex items-center gap-4">
        <NuxtLink 
          to="/admin/blog"
          class="p-2 text-neutral-500 hover:text-neutral-900 hover:bg-neutral-100 rounded-lg transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
          </svg>
        </NuxtLink>
        <div>
          <h1 class="text-2xl font-bold text-neutral-900">Tulis Artikel Baru</h1>
          <p class="text-neutral-500 mt-1">Buat artikel blog untuk meningkatkan SEO.</p>
        </div>
      </div>
    </div>

    <form @submit.prevent="submitForm" class="grid lg:grid-cols-3 gap-8">
      <!-- Main Content -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Title -->
        <div class="bg-white rounded-xl border border-neutral-200 p-6">
          <label class="block text-sm font-medium text-neutral-700 mb-2">Judul Artikel *</label>
          <input 
            v-model="form.title"
            @input="generateSlug"
            type="text" 
            class="w-full px-4 py-3 text-lg border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
            placeholder="Masukkan judul artikel..."
            required
          />
        </div>

        <!-- Slug -->
        <div class="bg-white rounded-xl border border-neutral-200 p-6">
          <label class="block text-sm font-medium text-neutral-700 mb-2">Slug URL *</label>
          <div class="flex items-center gap-2">
            <span class="text-neutral-500">/blog/</span>
            <input 
              v-model="form.slug"
              type="text" 
              class="flex-1 px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
              placeholder="judul-artikel"
              required
            />
          </div>
        </div>

        <!-- Excerpt -->
        <div class="bg-white rounded-xl border border-neutral-200 p-6">
          <label class="block text-sm font-medium text-neutral-700 mb-2">Ringkasan (Excerpt)</label>
          <textarea 
            v-model="form.excerpt"
            rows="3"
            class="w-full px-4 py-3 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
            placeholder="Ringkasan singkat artikel untuk preview..."
          ></textarea>
          <p class="text-xs text-neutral-500 mt-1">Akan ditampilkan di daftar blog dan meta description.</p>
        </div>

        <!-- Content -->
        <div class="bg-white rounded-xl border border-neutral-200 p-6">
          <label class="block text-sm font-medium text-neutral-700 mb-2">Konten Artikel *</label>
          <ClientOnly>
            <BlogEditor v-model="form.content" placeholder="Tulis konten artikel Anda di sini..." />
            <template #fallback>
              <div class="border border-neutral-200 rounded-lg p-6 min-h-[400px] flex items-center justify-center bg-neutral-50">
                <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-admin-600"></div>
              </div>
            </template>
          </ClientOnly>
        </div>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <!-- Publish Card -->
        <div class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-4">Publikasi</h3>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Status</label>
              <select 
                v-model="form.status"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
              >
                <option value="draft">Draft</option>
                <option value="published">Dipublikasi</option>
              </select>
            </div>
            
            <div class="flex gap-3 pt-4 border-t border-neutral-200">
              <button 
                type="submit"
                :disabled="submitting"
                class="flex-1 px-4 py-2.5 bg-admin-600 hover:bg-admin-700 text-white font-medium rounded-lg transition-colors disabled:opacity-50"
              >
                {{ submitting ? 'Menyimpan...' : 'Simpan' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Thumbnail -->
        <div class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-4">Thumbnail</h3>
          
          <div 
            v-if="form.thumbnail_url" 
            class="relative mb-4"
          >
            <img 
              :src="getThumbnailUrl(form.thumbnail_url)" 
              class="w-full h-40 object-cover rounded-lg"
            />
            <button 
              type="button"
              @click="form.thumbnail_url = ''"
              class="absolute top-2 right-2 p-1.5 bg-red-500 text-white rounded-lg hover:bg-red-600"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-neutral-700 mb-2">Upload Gambar</label>
            <div class="relative">
              <input 
                type="file"
                accept="image/*"
                @change="uploadThumbnail"
                :disabled="uploadingThumbnail"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-medium file:bg-admin-50 file:text-admin-700 hover:file:bg-admin-100"
              />
              <div v-if="uploadingThumbnail" class="absolute right-3 top-1/2 -translate-y-1/2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-admin-600"></div>
              </div>
            </div>
            <p class="text-xs text-neutral-500 mt-1">Format: JPG, PNG, WebP. Max 10MB.</p>
          </div>
        </div>

        <!-- SEO -->
        <div class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-4">SEO</h3>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Meta Title</label>
              <input 
                v-model="form.meta_title"
                type="text"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
                placeholder="Judul untuk search engine"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Meta Description</label>
              <textarea 
                v-model="form.meta_description"
                rows="3"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
                placeholder="Deskripsi untuk search engine"
              ></textarea>
            </div>
          </div>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const config = useRuntimeConfig()
const apiBase = config.public.apiBase
const router = useRouter()

const submitting = ref(false)
const uploadingThumbnail = ref(false)
const form = ref({
  title: '',
  slug: '',
  excerpt: '',
  content: '',
  thumbnail_url: '',
  status: 'draft',
  meta_title: '',
  meta_description: ''
})

const generateSlug = () => {
  form.value.slug = form.value.title
    .toLowerCase()
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .trim()
}

const getThumbnailUrl = (url) => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads')) return `${apiBase}${url}`
  // MinIO object key - use public images endpoint
  return `${apiBase}/api/images/${encodeURIComponent(url)}`
}

const uploadThumbnail = async (event) => {
  const file = event.target?.files?.[0]
  if (!file) return
  
  uploadingThumbnail.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    
    const token = useCookie('token')
    const response = await $fetch(`${apiBase}/api/admin/upload`, {
      method: 'POST',
      body: formData,
      headers: {
        'Authorization': `Bearer ${token.value}`
      }
    })
    
    if (response?.url || response?.object_key) {
      form.value.thumbnail_url = response.object_key || response.url
    }
  } catch (error) {
    console.error('Upload error:', error)
    alert('Gagal mengupload gambar')
  } finally {
    uploadingThumbnail.value = false
    event.target.value = ''
  }
}

const submitForm = async () => {
  submitting.value = true
  try {
    const token = useCookie('token')
    const payload = {
      title: form.value.title,
      slug: form.value.slug,
      excerpt: form.value.excerpt || null,
      content: form.value.content,
      thumbnail_url: form.value.thumbnail_url || null,
      status: form.value.status,
      meta_title: form.value.meta_title || null,
      meta_description: form.value.meta_description || null
    }
    
    await $fetch(`${apiBase}/api/admin/blog`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token.value}` },
      body: payload
    })
    
    router.push('/admin/blog')
  } catch (error) {
    console.error('Failed to create post:', error)
    alert(error.data?.error || 'Gagal membuat artikel')
  } finally {
    submitting.value = false
  }
}
</script>
