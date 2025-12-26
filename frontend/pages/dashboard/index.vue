<template>
  <div>
    <!-- Welcome Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-neutral-900">Selamat Datang, {{ user?.full_name || 'Student' }}! 👋</h1>
      <p class="text-neutral-500 mt-1">Lanjutkan perjalanan belajar Anda hari ini.</p>
    </div>
    
    <!-- Stats Grid -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-neutral-500 mb-1">Kursus Aktif</p>
            <p class="text-2xl font-bold text-neutral-900">{{ stats?.in_progress || 0 }}</p>
          </div>
          <div class="w-11 h-11 bg-primary-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
            </svg>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-neutral-500 mb-1">Modul Selesai</p>
            <p class="text-2xl font-bold text-neutral-900">{{ stats?.completed_courses || 0 }}</p>
          </div>
          <div class="w-11 h-11 bg-accent-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-neutral-500 mb-1">Jam Belajar</p>
            <p class="text-2xl font-bold text-neutral-900">-</p>
          </div>
          <div class="w-11 h-11 bg-warm-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-warm-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-neutral-500 mb-1">Sertifikat</p>
            <p class="text-2xl font-bold text-neutral-900">-</p>
          </div>
          <div class="w-11 h-11 bg-rose-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-rose-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/>
            </svg>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Continue Learning -->
    <div class="mb-8" v-if="loading || (recentCourses && recentCourses.length > 0)">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-lg font-semibold text-neutral-900">Lanjutkan Belajar</h2>
        <NuxtLink to="/dashboard/courses" class="text-sm text-primary-600 hover:text-primary-700 font-medium">Lihat Semua →</NuxtLink>
      </div>
      
      <!-- Loading State -->
      <div v-if="loading" class="grid md:grid-cols-2 xl:grid-cols-3 gap-5">
        <div class="bg-white rounded-xl h-64 animate-pulse bg-neutral-100"></div>
        <div class="bg-white rounded-xl h-64 animate-pulse bg-neutral-100"></div>
      </div>

      <!-- Real Course Cards -->
      <div v-else class="grid md:grid-cols-2 xl:grid-cols-3 gap-5">
        <div v-for="enrollment in recentCourses" :key="enrollment.id" class="bg-white rounded-xl border border-neutral-200 overflow-hidden group hover:shadow-lg transition-all">
          <div class="h-36 bg-primary-600 relative overflow-hidden flex items-center justify-center">
            <img v-if="enrollment.course?.thumbnail_url" :src="getThumbnailUrl(enrollment.course.thumbnail_url)" class="w-full h-full object-cover" />
            <svg v-else class="w-14 h-14 text-white/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
            </svg>
            <div v-if="enrollment.course?.category" class="absolute top-3 left-3">
              <span class="px-2.5 py-1 bg-white/20 text-white text-xs font-medium rounded-lg">{{ enrollment.course.category.name }}</span>
            </div>
          </div>
          <div class="p-5">
            <h3 class="font-semibold text-neutral-900 mb-1 line-clamp-1">{{ enrollment.course?.title || 'Untitled Course' }}</h3>
            <p class="text-sm text-neutral-500 mb-4">{{ enrollment.course?.instructor?.full_name || 'Instructor' }}</p>
            <div class="mb-4">
              <div class="flex justify-between text-xs mb-2">
                <span class="text-neutral-500">Progress</span>
                <span class="font-semibold text-primary-600">{{ enrollment.progress_percentage || 0 }}%</span>
              </div>
              <div class="h-2 bg-neutral-100 rounded-full overflow-hidden">
                <div class="h-full bg-primary-500 rounded-full" :style="{ width: (enrollment.progress_percentage || 0) + '%' }"></div>
              </div>
            </div>
            <NuxtLink :to="`/dashboard/courses/${enrollment.course_id}`" class="block w-full text-center py-2.5 text-sm font-medium text-primary-600 bg-primary-50 rounded-lg hover:bg-primary-100 transition-colors">
              Lanjutkan Belajar
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Empty State -->
    <div v-else-if="!loading && (!recentCourses || recentCourses.length === 0)" class="text-center py-12 bg-white rounded-xl border border-neutral-200 mb-8">
      <h3 class="text-lg font-semibold text-neutral-900">Belum ada kursus aktif</h3>
      <p class="text-neutral-500 mt-2 mb-4">Mulai belajar dengan mendaftar kursus baru</p>
      <NuxtLink to="/dashboard/explore" class="inline-block px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700">Jelajahi Kursus</NuxtLink>
    </div>
    
    <!-- Recent Activity -->
    <div>
      <h2 class="text-lg font-semibold text-neutral-900 mb-5">Aktivitas Terakhir</h2>
      <div class="bg-white rounded-xl border border-neutral-200 divide-y divide-neutral-100">
        <!-- Loading State -->
        <div v-if="loadingActivities" class="p-4 text-center">
          <div class="animate-spin w-6 h-6 border-2 border-primary-500 border-t-transparent rounded-full mx-auto"></div>
        </div>
        
        <!-- Real Activities -->
        <div v-else-if="activities.length > 0" v-for="activity in activities" :key="activity.id" class="flex items-center gap-4 p-4 hover:bg-neutral-50 transition-colors">
          <div 
            class="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0"
            :class="getActivityIcon(activity.activity_type).bgClass"
          >
            <svg class="w-5 h-5" :class="getActivityIcon(activity.activity_type).colorClass" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" :d="getActivityIcon(activity.activity_type).icon"/>
            </svg>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-neutral-900 truncate">{{ activity.description }}</p>
            <p class="text-xs text-neutral-500">{{ activity.metadata?.course_name || '' }}</p>
          </div>
          <span class="text-xs text-neutral-400 flex-shrink-0">{{ formatTimeAgo(activity.created_at) }}</span>
        </div>
        
        <!-- Empty State -->
        <div v-else class="p-6 text-center">
          <p class="text-sm text-neutral-500">Belum ada aktivitas</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { user } = useAuth()
const { stats, recentCourses, fetchUserDashboard, loading } = useDashboard()
const { activities, fetchActivities, formatTimeAgo, getActivityIcon, loading: loadingActivities } = useLearningActivities()

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

useHead({
  title: 'Dashboard - LearnHub'
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

onMounted(async () => {
  await Promise.all([
    fetchUserDashboard(),
    fetchActivities(5)
  ])
})
</script>

