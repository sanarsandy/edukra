<template>
  <div>
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-neutral-900">Jelajahi Kursus</h1>
      <p class="text-neutral-500 mt-1">Temukan kursus yang sesuai dengan minat dan tujuan Anda</p>
    </div>

    <!-- Search & Filter -->
    <div class="flex flex-col sm:flex-row gap-4 mb-8">
      <div class="relative flex-1">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari kursus..." 
          class="w-full pl-10 pr-4 py-3 bg-white border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm placeholder-neutral-500 transition-all"
        />
      </div>
      <select 
        v-model="selectedCategory"
        class="px-4 py-3 bg-white border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 text-sm text-neutral-700"
      >
        <option value="">Semua Kategori</option>
        <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
      </select>
    </div>

    <!-- Categories Pills -->
    <div class="flex flex-wrap gap-2 mb-8" v-if="!loadingCategories">
      <button 
        v-for="cat in ['Semua', ...categories]" 
        :key="cat"
        @click="selectedCategory = cat === 'Semua' ? '' : cat"
        :class="[
          'px-4 py-2 text-sm font-medium rounded-full transition-all',
          (selectedCategory === '' && cat === 'Semua') || selectedCategory === cat
            ? 'bg-primary-600 text-white' 
            : 'bg-white text-neutral-600 border border-neutral-200 hover:bg-neutral-50'
        ]"
      >
        {{ cat }}
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loadingCourses" class="grid md:grid-cols-2 xl:grid-cols-3 gap-6">
      <div class="bg-white rounded-xl h-80 animate-pulse bg-neutral-100"></div>
      <div class="bg-white rounded-xl h-80 animate-pulse bg-neutral-100"></div>
      <div class="bg-white rounded-xl h-80 animate-pulse bg-neutral-100"></div>
    </div>

    <!-- Courses Grid -->
    <div v-else class="grid md:grid-cols-2 xl:grid-cols-3 gap-6">
      <div 
        v-for="course in filteredCourses" 
        :key="course.id" 
        class="bg-white rounded-xl border border-neutral-200 overflow-hidden hover:shadow-lg transition-all group"
      >
        <div class="h-40 relative overflow-hidden flex items-center justify-center bg-gray-200">
          <img v-if="course.thumbnail_url" :src="getThumbnailUrl(course.thumbnail_url)" class="w-full h-full object-cover" />
          <svg v-else class="w-16 h-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
          </svg>
          <div class="absolute top-3 left-3" v-if="course.category">
            <span class="px-2.5 py-1 bg-black/50 text-white text-xs font-medium rounded-lg backdrop-blur-sm">{{ course.category.name }}</span>
          </div>
          <div v-if="course.is_featured" class="absolute top-3 right-3">
            <span class="px-2.5 py-1 bg-warm-500 text-white text-xs font-medium rounded-lg">Featured</span>
          </div>
        </div>
        <div class="p-5">
          <h3 class="font-semibold text-neutral-900 mb-1 group-hover:text-primary-600 transition-colors line-clamp-1">{{ course.title }}</h3>
          <p class="text-sm text-neutral-500 mb-3">{{ course.instructor?.full_name || 'Instructor' }}</p>
          
          <div class="flex items-center mb-3">
            <div class="flex items-center">
              <svg v-for="i in 5" :key="i" class="w-4 h-4" :class="i <= Math.round(getCourseRating(course.id)) ? 'text-yellow-400 fill-current' : 'text-neutral-200 fill-current'" viewBox="0 0 20 20">
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
              </svg>
              <span v-if="getCourseRating(course.id)" class="text-sm text-neutral-700 ml-1 font-medium">{{ getCourseRating(course.id).toFixed(1) }}</span>
              <span v-else class="text-sm text-neutral-400 ml-1">Belum ada rating</span>
            </div>
            <span class="text-sm text-neutral-500 ml-3">{{ course.lessons?.length || 0 }} Modul</span>
          </div>
          
          <div class="flex items-center justify-between pt-4 border-t border-neutral-100">
            <div>
              <span class="text-lg font-bold text-neutral-900">{{ new Intl.NumberFormat('id-ID', { style: 'currency', currency: course.currency || 'IDR' }).format(course.price) }}</span>
            </div>
            <NuxtLink :to="`/dashboard/courses/${course.id}`" class="px-4 py-2 text-sm font-medium text-primary-600 bg-primary-50 rounded-lg hover:bg-primary-100 transition-colors">
              Lihat Detail
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!loadingCourses && filteredCourses.length === 0" class="text-center py-16">
      <div class="w-20 h-20 bg-neutral-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-10 h-10 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-neutral-900 mb-2">Tidak ada kursus ditemukan</h3>
      <p class="text-neutral-500">Coba ubah kata kunci pencarian atau filter</p>
    </div>
  </div>
</template>

<script setup lang="ts">
const { courses, fetchCourses, loading: loadingCourses } = useCourses()
const { categories: apiCategories, fetchCategories, loading: loadingCategories } = useCategories()

// Rating stats for courses
const api = useApi()
const ratingStats = ref<Record<string, number>>({})

const fetchCourseRatings = async () => {
  if (!courses.value || courses.value.length === 0) return
  
  for (const course of courses.value) {
    try {
      const stats = await api.fetch<{ average_rating: number }>(`/api/courses/${course.id}/ratings/stats`)
      if (stats?.average_rating) {
        ratingStats.value[course.id] = stats.average_rating
      }
    } catch (err) {
      // Ignore errors for individual courses
    }
  }
}

const getCourseRating = (courseId: string): number => {
  return ratingStats.value[courseId] || 0
}

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

useHead({
  title: 'Jelajahi Kursus - LearnHub'
})

onMounted(async () => {
  await Promise.all([
    fetchCourses({ limit: 100 }),
    fetchCategories()
  ])
  // Fetch ratings after courses loaded
  await fetchCourseRatings()
})

const searchQuery = ref('')
const selectedCategory = ref('')

const categories = computed(() => apiCategories.value?.map((c: any) => c.name) || [])

const filteredCourses = computed(() => {
  if (!courses.value) return []
  return courses.value.filter(course => {
    const matchesSearch = course.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          course.instructor?.full_name?.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesCategory = !selectedCategory.value || course.category?.name === selectedCategory.value
    return matchesSearch && matchesCategory
  })
})

// Get proper thumbnail URL - handle MinIO objects
const getThumbnailUrl = (url: string | null | undefined): string => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads')) {
    const config = useRuntimeConfig()
    return `${config.public.apiBase}${url}`
  }
  const config = useRuntimeConfig()
  return `${config.public.apiBase}/api/images/${url}`
}
</script>

