<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Kursus Saya</h1>
        <p class="text-neutral-500 mt-1">Kelola dan lanjutkan kursus yang Anda ikuti</p>
      </div>
      <NuxtLink to="/dashboard/explore" class="btn-primary w-full sm:w-auto">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
        </svg>
        Cari Kursus Baru
      </NuxtLink>
    </div>

    <!-- Filter Tabs -->
    <div class="flex space-x-2 mb-6 overflow-x-auto pb-2">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        @click="activeTab = tab.id"
        :class="[
          'px-4 py-2 text-sm font-medium rounded-lg whitespace-nowrap transition-all',
          activeTab === tab.id 
            ? 'bg-primary-600 text-white' 
            : 'bg-white text-neutral-600 border border-neutral-200 hover:bg-neutral-50'
        ]"
      >
        {{ tab.name }}
        <span v-if="tab.count" class="ml-2 px-2 py-0.5 text-xs rounded-full" :class="activeTab === tab.id ? 'bg-primary-500 text-white' : 'bg-neutral-100 text-neutral-600'">
          {{ tab.count }}
        </span>
      </button>
    </div>

    <!-- Courses Grid -->
    <div class="grid md:grid-cols-2 xl:grid-cols-3 gap-5">
      <div v-for="course in filteredCourses" :key="course.id" class="bg-white rounded-xl border border-neutral-200 overflow-hidden hover:shadow-lg transition-all group">
        <div class="h-36 relative overflow-hidden flex items-center justify-center bg-neutral-200">
          <img v-if="course.thumbnail" :src="course.thumbnail" class="w-full h-full object-cover" />
          <svg v-else class="w-14 h-14 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
          </svg>
          <div class="absolute top-3 left-3" v-if="course.category">
            <span class="px-2.5 py-1 bg-black/50 text-white text-xs font-medium rounded-lg backdrop-blur-sm">{{ course.category }}</span>
          </div>
          <div v-if="course.progress >= 100" class="absolute top-3 right-3">
            <span class="px-2.5 py-1 bg-accent-500 text-white text-xs font-medium rounded-lg">Selesai</span>
          </div>
        </div>
        <div class="p-5">
          <h3 class="font-semibold text-neutral-900 mb-1 group-hover:text-primary-600 transition-colors">{{ course.title }}</h3>
          <p class="text-sm text-neutral-500 mb-3">{{ course.instructor }}</p>
          
          <div class="flex items-center text-xs text-neutral-400 mb-4 space-x-4">
            <span class="flex items-center">
              <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
              </svg>
              {{ course.lessons }} Modul
            </span>
            <span class="flex items-center">
              <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              {{ course.duration }}
            </span>
          </div>
          
          <div class="mb-4">
            <div class="flex justify-between text-xs mb-2">
              <span class="text-neutral-500">Progress</span>
              <span class="font-semibold" :class="course.progress === 100 ? 'text-accent-600' : 'text-primary-600'">{{ course.progress }}%</span>
            </div>
            <div class="h-2 bg-neutral-100 rounded-full overflow-hidden">
              <div class="h-full rounded-full transition-all" :class="course.progress === 100 ? 'bg-accent-500' : 'bg-primary-500'" :style="{ width: course.progress + '%' }"></div>
            </div>
          </div>
          
          <NuxtLink :to="`/dashboard/courses/${course.id}`" class="w-full py-2.5 text-sm font-medium text-primary-600 bg-primary-50 rounded-lg hover:bg-primary-100 transition-colors flex items-center justify-center">
            {{ course.progress === 100 ? 'Lihat Kursus' : 'Lanjutkan Belajar' }}
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="filteredCourses.length === 0" class="text-center py-16">
      <div class="w-20 h-20 bg-neutral-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-10 h-10 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-neutral-900 mb-2">Belum ada kursus</h3>
      <p class="text-neutral-500 mb-6">Mulai perjalanan belajar Anda dengan mengikuti kursus</p>
      <NuxtLink to="/dashboard/explore" class="btn-primary">Jelajahi Kursus</NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
const { enrollments, fetchEnrollments, loading, total } = useEnrollments()

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

useHead({
  title: 'Kursus Saya - LearnHub'
})

onMounted(async () => {
  await fetchEnrollments({ limit: 100 })
})

const activeTab = ref('all')

const tabs = computed(() => {
  const all = enrollments.value.length
  const inProgress = enrollments.value.filter(e => (e.progress_percentage || 0) > 0 && (e.progress_percentage || 0) < 100).length
  const completed = enrollments.value.filter(e => (e.progress_percentage || 0) >= 100).length
  return [
    { id: 'all', name: 'Semua', count: all },
    { id: 'in-progress', name: 'Sedang Dipelajari', count: inProgress },
    { id: 'completed', name: 'Selesai', count: completed }
  ]
})

const courses = computed(() => {
  return enrollments.value.map(e => ({
    id: e.course_id,
    title: e.course?.title || 'Unknown Course',
    instructor: e.course?.instructor?.full_name || 'Instructor',
    category: e.course?.category?.name || '',
    color: 'bg-primary-600',
    thumbnail: e.course?.thumbnail_url,
    lessons: e.course?.lessons?.length || 0,
    duration: '-',
    progress: e.progress_percentage || 0
  }))
})

const filteredCourses = computed(() => {
  if (activeTab.value === 'all') return courses.value
  if (activeTab.value === 'in-progress') return courses.value.filter(c => c.progress > 0 && c.progress < 100)
  if (activeTab.value === 'completed') return courses.value.filter(c => c.progress >= 100)
  return courses.value
})
</script>

