<template>
  <div>
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Kelola Blog</h1>
        <p class="text-neutral-500 mt-1">Buat dan kelola artikel blog untuk SEO.</p>
      </div>
      <NuxtLink 
        to="/admin/blog/create"
        class="inline-flex items-center px-4 py-2.5 bg-admin-600 hover:bg-admin-700 text-white font-medium rounded-lg transition-colors"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        Tulis Artikel
      </NuxtLink>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-neutral-900">{{ stats.total }}</p>
        <p class="text-sm text-neutral-500">Total Artikel</p>
      </div>
      
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-success-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-neutral-900">{{ stats.published }}</p>
        <p class="text-sm text-neutral-500">Dipublikasi</p>
      </div>
      
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-warm-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-warm-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-neutral-900">{{ stats.draft }}</p>
        <p class="text-sm text-neutral-500">Draft</p>
      </div>
      
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-accent-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-neutral-900">{{ stats.totalViews }}</p>
        <p class="text-sm text-neutral-500">Total Views</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-xl border border-neutral-200 p-4 mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <div class="flex-1 min-w-[200px]">
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Cari artikel..." 
              class="w-full pl-10 pr-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
            />
          </div>
        </div>
        <select 
          v-model="filterStatus"
          class="px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
        >
          <option value="all">Semua Status</option>
          <option value="published">Dipublikasi</option>
          <option value="draft">Draft</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-admin-600"></div>
    </div>

    <!-- Posts Table -->
    <div v-else class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-neutral-50 border-b border-neutral-200">
            <tr>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Artikel</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Penulis</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Views</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Status</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Tanggal</th>
              <th class="px-6 py-4 text-right text-xs font-semibold text-neutral-600 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-neutral-100">
            <tr v-for="post in filteredPosts" :key="post.id" class="hover:bg-neutral-50 transition-colors">
              <td class="px-6 py-4">
                <div class="flex items-center gap-4">
                  <div class="w-16 h-12 bg-neutral-100 rounded-lg overflow-hidden flex-shrink-0">
                    <img 
                      v-if="post.thumbnail_url" 
                      :src="getThumbnailUrl(post.thumbnail_url)" 
                      :alt="post.title"
                      class="w-full h-full object-cover"
                    />
                    <div v-else class="w-full h-full flex items-center justify-center text-neutral-400">
                      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                      </svg>
                    </div>
                  </div>
                  <div class="min-w-0">
                    <h3 class="font-semibold text-neutral-900 truncate max-w-xs">{{ post.title }}</h3>
                    <p class="text-sm text-neutral-500 truncate max-w-xs">/blog/{{ post.slug }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <span class="text-neutral-700">{{ post.author?.full_name || 'Admin' }}</span>
              </td>
              <td class="px-6 py-4">
                <span class="text-neutral-900 font-medium">{{ post.view_count || 0 }}</span>
              </td>
              <td class="px-6 py-4">
                <span 
                  :class="[
                    'px-2.5 py-1 rounded-full text-xs font-medium',
                    post.status === 'published' ? 'bg-success-100 text-success-700' : 'bg-neutral-100 text-neutral-600'
                  ]"
                >
                  {{ post.status === 'published' ? 'Dipublikasi' : 'Draft' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-neutral-600">
                <div v-if="post.published_at">
                  {{ formatDate(post.published_at) }}
                </div>
                <div v-else>
                  {{ formatDate(post.created_at) }}
                </div>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <a 
                    v-if="post.status === 'published'"
                    :href="`/blog/${post.slug}`"
                    target="_blank"
                    class="p-2 text-neutral-500 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors"
                    title="Lihat"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
                    </svg>
                  </a>
                  <NuxtLink 
                    :to="`/admin/blog/${post.id}`"
                    class="p-2 text-neutral-500 hover:text-admin-600 hover:bg-admin-50 rounded-lg transition-colors"
                    title="Edit"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                    </svg>
                  </NuxtLink>
                  <button 
                    @click="confirmDelete(post)"
                    class="p-2 text-neutral-500 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    title="Hapus"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredPosts.length === 0">
              <td colspan="6" class="px-6 py-12 text-center text-neutral-500">
                <svg class="w-12 h-12 mx-auto text-neutral-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"/>
                </svg>
                <p>Belum ada artikel. Tulis artikel pertama Anda!</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50" @click="showDeleteModal = false"></div>
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-md p-6">
        <div class="text-center">
          <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
            </svg>
          </div>
          <h3 class="text-lg font-bold text-neutral-900 mb-2">Hapus Artikel?</h3>
          <p class="text-neutral-600 mb-6">
            Artikel <strong>{{ postToDelete?.title }}</strong> akan dihapus permanen. Tindakan ini tidak dapat dibatalkan.
          </p>
          <div class="flex items-center justify-center gap-3">
            <button 
              @click="showDeleteModal = false"
              class="px-5 py-2.5 text-neutral-600 hover:text-neutral-900 font-medium rounded-lg hover:bg-neutral-100 transition-colors"
            >
              Batal
            </button>
            <button 
              @click="deletePost"
              :disabled="deleting"
              class="px-5 py-2.5 bg-red-600 hover:bg-red-700 text-white font-medium rounded-lg transition-colors disabled:opacity-50"
            >
              {{ deleting ? 'Menghapus...' : 'Ya, Hapus' }}
            </button>
          </div>
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
const apiBase = config.public.apiBase

const loading = ref(true)
const posts = ref([])
const searchQuery = ref('')
const filterStatus = ref('all')

// Delete modal state
const showDeleteModal = ref(false)
const postToDelete = ref(null)
const deleting = ref(false)

// Stats
const stats = computed(() => {
  const total = posts.value.length
  const published = posts.value.filter(p => p.status === 'published').length
  const draft = posts.value.filter(p => p.status === 'draft').length
  const totalViews = posts.value.reduce((sum, p) => sum + (p.view_count || 0), 0)
  return { total, published, draft, totalViews }
})

const filteredPosts = computed(() => {
  let result = posts.value
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(p => 
      p.title.toLowerCase().includes(query) || 
      p.slug.toLowerCase().includes(query)
    )
  }
  
  if (filterStatus.value !== 'all') {
    result = result.filter(p => p.status === filterStatus.value)
  }
  
  return result
})

const getThumbnailUrl = (url) => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads')) return `${apiBase}${url}`
  return `${apiBase}/api/images/${url}`
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const loadPosts = async () => {
  loading.value = true
  try {
    const token = useCookie('token')
    const response = await $fetch(`${apiBase}/api/admin/blog`, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    posts.value = response.posts || []
  } catch (error) {
    console.error('Failed to load posts:', error)
  } finally {
    loading.value = false
  }
}

const confirmDelete = (post) => {
  postToDelete.value = post
  showDeleteModal.value = true
}

const deletePost = async () => {
  deleting.value = true
  try {
    const token = useCookie('token')
    await $fetch(`${apiBase}/api/admin/blog/${postToDelete.value.id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token.value}` }
    })
    showDeleteModal.value = false
    await loadPosts()
  } catch (error) {
    console.error('Failed to delete post:', error)
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  loadPosts()
})
</script>
