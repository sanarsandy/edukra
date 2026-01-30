<template>
  <div class="min-h-screen bg-neutral-50">
    <!-- Navigation -->
    <nav class="bg-white border-b border-neutral-200 sticky top-0 z-50">
      <div class="container-custom">
        <div class="flex items-center justify-between h-16">
          <NuxtLink to="/" class="flex items-center space-x-2">
            <img src="/logo.png" alt="EDUKRA" class="w-10 h-10 object-contain" />
            <span class="font-display font-bold text-xl text-neutral-900">EDUKRA</span>
          </NuxtLink>
          <div class="flex items-center gap-4">
            <NuxtLink to="/dashboard/explore" class="text-neutral-600 hover:text-primary-600 transition-colors font-medium">
              Kursus
            </NuxtLink>
            <NuxtLink to="/login" class="px-4 py-2 text-sm font-medium text-primary-600 hover:bg-primary-50 rounded-lg transition-colors">
              Masuk
            </NuxtLink>
            <NuxtLink to="/register" class="px-4 py-2 text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 rounded-lg transition-colors">
              Daftar
            </NuxtLink>
          </div>
        </div>
      </div>
    </nav>

    <!-- Hero -->
    <section class="bg-white border-b border-neutral-200">
      <div class="container-custom py-16 text-center">
        <h1 class="text-4xl md:text-5xl font-display font-bold text-neutral-900 mb-4">
          Blog & Artikel
        </h1>
        <p class="text-lg text-neutral-600 max-w-2xl mx-auto">
          Tips, panduan, dan insight terbaru seputar pembelajaran dan pengembangan karir.
        </p>
      </div>
    </section>

    <!-- Content -->
    <section class="py-12">
      <div class="container-custom">
        <!-- Loading -->
        <div v-if="pending" class="flex items-center justify-center py-24">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        </div>

        <!-- Posts Grid -->
        <div v-else-if="posts && posts.length > 0" class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <article 
            v-for="post in posts" 
            :key="post.id"
            class="bg-white rounded-2xl border border-neutral-200 overflow-hidden hover:shadow-lg transition-all duration-300 group"
          >
            <NuxtLink :to="`/blog/${post.slug}`">
              <div class="h-48 bg-neutral-100 overflow-hidden">
                <img 
                  v-if="post.thumbnail_url"
                  :src="getThumbnailUrl(post.thumbnail_url)"
                  :alt="post.title"
                  class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                />
                <div v-else class="w-full h-full flex items-center justify-center text-neutral-400">
                  <svg class="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"/>
                  </svg>
                </div>
              </div>
              <div class="p-6">
                <div class="flex items-center gap-3 mb-3">
                  <span class="text-xs text-neutral-500">{{ formatDate(post.published_at || post.created_at) }}</span>
                  <span class="text-xs text-neutral-400">•</span>
                  <span class="text-xs text-neutral-500">{{ post.view_count || 0 }} views</span>
                </div>
                <h2 class="text-xl font-bold text-neutral-900 mb-2 group-hover:text-primary-600 transition-colors line-clamp-2">
                  {{ post.title }}
                </h2>
                <p class="text-neutral-600 text-sm line-clamp-3">
                  {{ post.excerpt || stripHtml(post.content).substring(0, 150) }}...
                </p>
                <div class="mt-4 flex items-center text-primary-600 text-sm font-medium">
                  Baca Selengkapnya
                  <svg class="w-4 h-4 ml-1 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                  </svg>
                </div>
              </div>
            </NuxtLink>
          </article>
        </div>

        <!-- Empty State -->
        <div v-else class="text-center py-24">
          <svg class="w-16 h-16 mx-auto text-neutral-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"/>
          </svg>
          <p class="text-neutral-500">Belum ada artikel.</p>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex justify-center mt-12">
          <div class="flex items-center gap-2">
            <button 
              @click="page > 1 && (page = page - 1)"
              :disabled="page <= 1"
              class="px-4 py-2 text-sm font-medium text-neutral-600 bg-white border border-neutral-200 rounded-lg hover:bg-neutral-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Sebelumnya
            </button>
            <span class="px-4 py-2 text-sm text-neutral-600">
              Halaman {{ page }} dari {{ totalPages }}
            </span>
            <button 
              @click="page < totalPages && (page = page + 1)"
              :disabled="page >= totalPages"
              class="px-4 py-2 text-sm font-medium text-neutral-600 bg-white border border-neutral-200 rounded-lg hover:bg-neutral-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Selanjutnya
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-neutral-900 text-neutral-400 py-12 mt-12">
      <div class="container-custom text-center">
        <p>&copy; 2024 EDUKRA. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
useHead({
  title: 'Blog - EDUKRA',
  meta: [
    { name: 'description', content: 'Baca artikel, tips, dan panduan terbaru seputar pembelajaran dan pengembangan karir di EDUKRA.' }
  ]
})

const config = useRuntimeConfig()
const apiBase = config.public.apiBase

const page = ref(1)
const perPage = 9

const { data, pending } = await useFetch(`${apiBase}/api/blog`, {
  query: { page, per_page: perPage },
  watch: [page]
})

const posts = computed(() => data.value?.posts || [])
const totalPages = computed(() => data.value?.total_pages || 1)

const getThumbnailUrl = (url) => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads')) return `${apiBase}${url}`
  return `${apiBase}/api/images/${url}`
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

const stripHtml = (html) => {
  if (!html) return ''
  return html.replace(/<[^>]*>/g, '')
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
