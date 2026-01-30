<template>
  <section class="py-16 md:py-24 bg-neutral-50">
    <div class="container-custom">
      <!-- Section Header -->
      <div class="text-center mb-12">
        <div class="inline-flex items-center px-4 py-2 bg-primary-50 rounded-full mb-4 border border-primary-100">
          <svg class="w-4 h-4 text-primary-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2"/>
          </svg>
          <span class="text-sm font-medium text-primary-700">Blog & Artikel</span>
        </div>
        <h2 class="text-3xl md:text-4xl font-bold text-neutral-900 font-display mb-4">
          Artikel Terbaru
        </h2>
        <p class="text-neutral-600 max-w-2xl mx-auto">
          Baca artikel menarik seputar pembelajaran dan pengembangan diri
        </p>
      </div>

      <!-- Loading -->
      <div v-if="pending" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>

      <!-- Slider -->
      <div v-else-if="posts.length" class="relative">
        <!-- Slider Container -->
        <div class="overflow-hidden rounded-2xl">
          <div 
            class="flex transition-transform duration-500 ease-in-out"
            :style="{ transform: `translateX(-${currentSlide * 100}%)` }"
          >
            <div 
              v-for="post in posts" 
              :key="post.id"
              class="w-full flex-shrink-0 px-2"
            >
              <NuxtLink :to="`/blog/${post.slug}`" class="block">
                <div class="bg-white rounded-2xl overflow-hidden shadow-soft hover:shadow-soft-lg transition-all duration-300 group">
                  <div class="grid md:grid-cols-2 gap-0">
                    <!-- Image -->
                    <div class="relative h-64 md:h-80 overflow-hidden">
                      <img 
                        v-if="post.thumbnail_url"
                        :src="getThumbnailUrl(post.thumbnail_url)" 
                        :alt="post.title"
                        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                      />
                      <div v-else class="w-full h-full bg-gradient-to-br from-primary-100 to-warm-100 flex items-center justify-center">
                        <svg class="w-16 h-16 text-primary-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2"/>
                        </svg>
                      </div>
                    </div>
                    <!-- Content -->
                    <div class="p-8 flex flex-col justify-center">
                      <div class="text-sm text-primary-600 font-medium mb-3">
                        {{ formatDate(post.published_at) }}
                      </div>
                      <h3 class="text-2xl font-bold text-neutral-900 mb-4 group-hover:text-primary-600 transition-colors line-clamp-2">
                        {{ post.title }}
                      </h3>
                      <p class="text-neutral-600 mb-6 line-clamp-3">
                        {{ post.excerpt || stripHtml(post.content) }}
                      </p>
                      <div class="flex items-center text-primary-600 font-medium">
                        <span>Baca Selengkapnya</span>
                        <svg class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"/>
                        </svg>
                      </div>
                    </div>
                  </div>
                </div>
              </NuxtLink>
            </div>
          </div>
        </div>

        <!-- Navigation Dots -->
        <div v-if="posts.length > 1" class="flex justify-center gap-2 mt-6">
          <button
            v-for="(_, index) in posts"
            :key="index"
            @click="goToSlide(index)"
            class="w-3 h-3 rounded-full transition-all duration-300"
            :class="currentSlide === index ? 'bg-primary-600 w-8' : 'bg-neutral-300 hover:bg-neutral-400'"
          />
        </div>

        <!-- Navigation Arrows -->
        <button
          v-if="posts.length > 1"
          @click="prevSlide"
          class="absolute left-4 top-1/2 -translate-y-1/2 w-12 h-12 bg-white/90 backdrop-blur rounded-full shadow-lg flex items-center justify-center text-neutral-700 hover:bg-white hover:text-primary-600 transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
          </svg>
        </button>
        <button
          v-if="posts.length > 1"
          @click="nextSlide"
          class="absolute right-4 top-1/2 -translate-y-1/2 w-12 h-12 bg-white/90 backdrop-blur rounded-full shadow-lg flex items-center justify-center text-neutral-700 hover:bg-white hover:text-primary-600 transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
          </svg>
        </button>
      </div>

      <!-- No Posts -->
      <div v-else class="text-center py-12 text-neutral-500">
        Belum ada artikel yang dipublikasi.
      </div>

      <!-- View All Button -->
      <div v-if="posts.length" class="text-center mt-10">
        <NuxtLink 
          to="/blog" 
          class="inline-flex items-center px-6 py-3 bg-primary-600 hover:bg-primary-700 text-white font-medium rounded-xl transition-colors"
        >
          Lihat Semua Artikel
          <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"/>
          </svg>
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

<script setup>
const config = useRuntimeConfig()
const apiBase = config.public.apiBase

const currentSlide = ref(0)
let autoSlideInterval = null

const { data, pending } = await useFetch('/api/blog', {
  baseURL: apiBase,
  query: { page: 1, per_page: 5 },
  server: false
})

const posts = computed(() => data.value?.posts || [])

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
  return html.replace(/<[^>]*>/g, '').substring(0, 150) + '...'
}

const nextSlide = () => {
  currentSlide.value = (currentSlide.value + 1) % posts.value.length
}

const prevSlide = () => {
  currentSlide.value = currentSlide.value === 0 ? posts.value.length - 1 : currentSlide.value - 1
}

const goToSlide = (index) => {
  currentSlide.value = index
}

const startAutoSlide = () => {
  autoSlideInterval = setInterval(() => {
    if (posts.value.length > 1) {
      nextSlide()
    }
  }, 5000)
}

const stopAutoSlide = () => {
  if (autoSlideInterval) {
    clearInterval(autoSlideInterval)
  }
}

onMounted(() => {
  startAutoSlide()
})

onBeforeUnmount(() => {
  stopAutoSlide()
})
</script>
