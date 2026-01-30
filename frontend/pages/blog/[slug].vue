<template>
  <div class="min-h-screen bg-white">
    <!-- Navigation -->
    <nav class="bg-white border-b border-neutral-200 sticky top-0 z-50">
      <div class="container-custom">
        <div class="flex items-center justify-between h-16">
          <NuxtLink to="/" class="flex items-center space-x-2">
            <img src="/logo.png" alt="EDUKRA" class="w-10 h-10 object-contain" />
            <span class="font-display font-bold text-xl text-neutral-900">EDUKRA</span>
          </NuxtLink>
          <div class="flex items-center gap-4">
            <NuxtLink to="/blog" class="text-neutral-600 hover:text-primary-600 transition-colors font-medium">
              Blog
            </NuxtLink>
            <NuxtLink to="/dashboard/explore" class="text-neutral-600 hover:text-primary-600 transition-colors font-medium">
              Kursus
            </NuxtLink>
            <NuxtLink to="/login" class="px-4 py-2 text-sm font-medium text-primary-600 hover:bg-primary-50 rounded-lg transition-colors">
              Masuk
            </NuxtLink>
          </div>
        </div>
      </div>
    </nav>

    <!-- Loading -->
    <div v-if="pending" class="flex items-center justify-center py-32">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
    </div>

    <!-- Not Found -->
    <div v-else-if="!post" class="text-center py-32">
      <h1 class="text-2xl font-bold text-neutral-900 mb-4">Artikel Tidak Ditemukan</h1>
      <NuxtLink to="/blog" class="text-primary-600 hover:underline">
        ← Kembali ke Blog
      </NuxtLink>
    </div>

    <!-- Article -->
    <article v-else>
      <!-- Hero -->
      <header class="bg-neutral-50 border-b border-neutral-200">
        <div class="container-custom max-w-4xl py-12">
          <NuxtLink to="/blog" class="inline-flex items-center text-sm text-primary-600 hover:text-primary-700 mb-6">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
            </svg>
            Kembali ke Blog
          </NuxtLink>
          
          <h1 class="text-3xl md:text-4xl lg:text-5xl font-display font-bold text-neutral-900 mb-6 leading-tight">
            {{ post.title }}
          </h1>
          
          <div class="flex items-center gap-4 text-sm text-neutral-500">
            <div class="flex items-center gap-2">
              <div class="w-10 h-10 bg-primary-100 rounded-full flex items-center justify-center text-primary-600 font-bold">
                {{ getInitials(post.author?.full_name || 'Admin') }}
              </div>
              <div>
                <p class="font-medium text-neutral-900">{{ post.author?.full_name || 'Admin' }}</p>
                <p>{{ formatDate(post.published_at || post.created_at) }}</p>
              </div>
            </div>
            <span class="text-neutral-300">|</span>
            <span>{{ post.view_count || 0 }} views</span>
          </div>
        </div>
      </header>

      <!-- Featured Image -->
      <div v-if="post.thumbnail_url" class="container-custom max-w-4xl py-8">
        <img 
          :src="getThumbnailUrl(post.thumbnail_url)"
          :alt="post.title"
          class="w-full h-auto rounded-2xl shadow-lg"
        />
      </div>

      <!-- Content -->
      <div class="container-custom max-w-4xl py-8">
        <div 
          class="prose prose-lg max-w-none prose-headings:font-display prose-headings:text-neutral-900 prose-p:text-neutral-700 prose-a:text-primary-600 prose-strong:text-neutral-900"
          v-html="post.content"
        ></div>
      </div>

      <!-- Share & CTA -->
      <div class="container-custom max-w-4xl py-12 border-t border-neutral-200">
        <div class="bg-primary-50 rounded-2xl p-8 text-center">
          <h3 class="text-2xl font-bold text-neutral-900 mb-3">Tertarik Belajar Lebih?</h3>
          <p class="text-neutral-600 mb-6">Jelajahi kursus berkualitas dari instruktur terbaik kami.</p>
          <NuxtLink 
            to="/dashboard/explore"
            class="inline-flex items-center px-6 py-3 bg-primary-600 hover:bg-primary-700 text-white font-medium rounded-xl transition-colors"
          >
            Lihat Kursus
            <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"/>
            </svg>
          </NuxtLink>
        </div>
      </div>
    </article>

    <!-- Footer -->
    <footer class="bg-neutral-900 text-neutral-400 py-12">
      <div class="container-custom text-center">
        <p>&copy; 2024 EDUKRA. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
const route = useRoute()
const config = useRuntimeConfig()
const apiBase = config.public.apiBase

const { data: post, pending, error } = await useFetch(`/api/blog/${route.params.slug}`, {
  baseURL: apiBase,
  server: false
})

// SEO
useHead(() => ({
  title: post.value?.meta_title || post.value?.title ? `${post.value.meta_title || post.value.title} - EDUKRA Blog` : 'Blog - EDUKRA',
  meta: [
    { name: 'description', content: post.value?.meta_description || post.value?.excerpt || '' },
    { property: 'og:title', content: post.value?.meta_title || post.value?.title || '' },
    { property: 'og:description', content: post.value?.meta_description || post.value?.excerpt || '' },
    { property: 'og:image', content: post.value?.thumbnail_url || '' },
    { property: 'og:type', content: 'article' }
  ]
}))

const getThumbnailUrl = (url) => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads')) return `${apiBase}${url}`
  // MinIO object key - use public images endpoint
  return `${apiBase}/api/images/${encodeURIComponent(url)}`
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

const getInitials = (name) => {
  if (!name) return 'A'
  return name.split(' ').map(n => n[0]).join('').toUpperCase().substring(0, 2)
}
</script>

<style scoped>
/* Blog content styles for special elements */
:deep(.prose table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1.5rem 0;
}

:deep(.prose th),
:deep(.prose td) {
  border: 1px solid #e5e7eb;
  padding: 0.75rem 1rem;
  text-align: left;
}

:deep(.prose th) {
  background-color: #f3f4f6;
  font-weight: 600;
}

:deep(.prose tr:nth-child(even)) {
  background-color: #f9fafb;
}

/* Task list */
:deep(.prose ul[data-type="taskList"]) {
  list-style: none;
  padding-left: 0;
}

:deep(.prose ul[data-type="taskList"] li) {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

:deep(.prose ul[data-type="taskList"] li input[type="checkbox"]) {
  margin-top: 0.25rem;
}

/* YouTube embed */
:deep(.prose iframe) {
  width: 100%;
  aspect-ratio: 16 / 9;
  border-radius: 0.75rem;
  margin: 1.5rem 0;
}

/* Highlight */
:deep(.prose mark) {
  padding: 0.125rem 0.25rem;
  border-radius: 0.25rem;
}

/* Subscript & Superscript */
:deep(.prose sub) {
  font-size: 0.75em;
  vertical-align: sub;
}

:deep(.prose sup) {
  font-size: 0.75em;
  vertical-align: super;
}

/* Blockquote */
:deep(.prose blockquote) {
  border-left: 4px solid #6366f1;
  padding-left: 1rem;
  margin: 1.5rem 0;
  font-style: italic;
  color: #6b7280;
  background-color: #f9fafb;
  padding: 1rem;
  border-radius: 0 0.5rem 0.5rem 0;
}

/* Code blocks */
:deep(.prose pre) {
  background-color: #1f2937;
  color: #f3f4f6;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
  margin: 1.5rem 0;
}

:deep(.prose code) {
  background-color: #f3f4f6;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
  font-size: 0.875em;
  color: #dc2626;
}

:deep(.prose pre code) {
  background-color: transparent;
  padding: 0;
  color: inherit;
}

/* Images */
:deep(.prose img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.75rem;
  margin: 1.5rem auto;
}
</style>
