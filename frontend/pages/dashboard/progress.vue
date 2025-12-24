<template>
  <div>
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-neutral-900">Progress Belajar</h1>
      <p class="text-neutral-500 mt-1">Pantau perkembangan belajar Anda</p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-neutral-500 mb-1">Total Jam Belajar</p>
            <p class="text-2xl font-bold text-neutral-900">{{ stats.totalHours }}</p>
          </div>
          <div class="w-11 h-11 bg-primary-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-neutral-500 mb-1">Modul Selesai</p>
            <p class="text-2xl font-bold text-neutral-900">{{ stats.completedLessons }}</p>
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
            <p class="text-sm text-neutral-500 mb-1">Kuis Lulus</p>
            <p class="text-2xl font-bold text-neutral-900">{{ stats.quizzesPassed }}</p>
          </div>
          <div class="w-11 h-11 bg-warm-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-warm-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
            </svg>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-neutral-500 mb-1">Streak Belajar</p>
            <p class="text-2xl font-bold text-neutral-900">{{ stats.streak }} hari</p>
          </div>
          <div class="w-11 h-11 bg-rose-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-rose-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z"/>
            </svg>
          </div>
        </div>
      </div>
    </div>

    <div class="grid lg:grid-cols-3 gap-6">
      <!-- Weekly Activity -->
      <div class="lg:col-span-2 bg-white rounded-xl border border-neutral-200 p-6">
        <h3 class="font-semibold text-neutral-900 mb-6">Aktivitas Mingguan</h3>
        <div class="flex items-end justify-between h-40 gap-2">
          <div v-for="day in weeklyActivity" :key="day.name" class="flex-1 flex flex-col items-center">
            <div class="w-full bg-neutral-100 rounded-t-lg relative" style="height: 120px;">
              <div 
                class="absolute bottom-0 w-full bg-primary-500 rounded-t-lg transition-all"
                :style="{ height: getBarHeight(day.hours) + '%' }"
              ></div>
            </div>
            <span class="text-xs text-neutral-500 mt-2">{{ day.name }}</span>
            <span class="text-xs font-medium text-neutral-700">{{ day.hours }}</span>
          </div>
        </div>
      </div>

      <!-- Course Progress -->
      <div class="bg-white rounded-xl border border-neutral-200 p-6">
        <h3 class="font-semibold text-neutral-900 mb-6">Progress Kursus</h3>
        <div class="space-y-4">
          <div v-for="course in courseProgress" :key="course.name">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm text-neutral-700 truncate pr-4">{{ course.name }}</span>
              <span class="text-sm font-medium" :class="course.progress === 100 ? 'text-accent-600' : 'text-primary-600'">{{ course.progress }}%</span>
            </div>
            <div class="h-2 bg-neutral-100 rounded-full overflow-hidden">
              <div 
                class="h-full rounded-full transition-all"
                :class="course.progress === 100 ? 'bg-accent-500' : 'bg-primary-500'"
                :style="{ width: course.progress + '%' }"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="mt-6 bg-white rounded-xl border border-neutral-200">
      <div class="p-5 border-b border-neutral-100">
        <h3 class="font-semibold text-neutral-900">Aktivitas Terakhir</h3>
      </div>
      <div class="divide-y divide-neutral-100">
        <div v-for="activity in recentActivity" :key="activity.id" class="flex items-center gap-4 p-4 hover:bg-neutral-50 transition-colors">
          <div 
            class="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0"
            :class="activity.iconBg"
          >
            <svg class="w-5 h-5" :class="activity.iconColor" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" :d="activity.icon"/>
            </svg>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-neutral-900 truncate">{{ activity.title }}</p>
            <p class="text-xs text-neutral-500">{{ activity.course }}</p>
          </div>
          <span class="text-xs text-neutral-400 flex-shrink-0">{{ activity.time }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { stats: dashboardStats, fetchUserDashboard, loading: loadingDashboard } = useDashboard()
const { enrollments, fetchEnrollments, loading: loadingEnrollments } = useEnrollments()
const { stats: learningStats, activities: recentActivities, fetchStats, fetchActivities, formatTimeAgo, formatWatchTime, getActivityIcon, loading: loadingStats } = useLearningActivities()

// Fetch weekly activities
const api = useApi()
const weeklyData = ref<{ name: string; count: number }[]>([])
const loadingWeekly = ref(false)

const fetchWeeklyActivities = async () => {
  loadingWeekly.value = true
  try {
    const response = await api.fetch<{ weekly: { name: string; count: number }[] }>('/api/activities/weekly')
    weeklyData.value = response.weekly || []
  } catch (err) {
    console.error('Failed to fetch weekly activities:', err)
  } finally {
    loadingWeekly.value = false
  }
}

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

useHead({
  title: 'Progress Belajar - LearnHub'
})

onMounted(async () => {
  await Promise.all([
    fetchUserDashboard(),
    fetchEnrollments({ limit: 50 }),
    fetchStats(),
    fetchActivities(10),
    fetchWeeklyActivities()
  ])
})

const loading = computed(() => loadingDashboard.value || loadingEnrollments.value || loadingStats.value)

// Map real stats from learning stats API
const stats = computed(() => ({
  totalHours: learningStats.value ? formatWatchTime(learningStats.value.total_watch_time) : '-',
  completedLessons: learningStats.value?.lessons_completed || dashboardStats.value?.completed_courses || 0,
  quizzesPassed: learningStats.value?.quizzes_passed || 0,
  streak: learningStats.value?.current_streak || 0
}))

// Course progress from enrollments
const courseProgress = computed(() => {
  return enrollments.value.map(e => ({
    name: e.course?.title || 'Unknown Course',
    progress: e.progress_percentage || 0
  }))
})

// Weekly activity from API
const weeklyActivity = computed(() => {
  if (weeklyData.value.length > 0) {
    return weeklyData.value.map(d => ({
      name: d.name,
      hours: d.count // Using count as activity level
    }))
  }
  // Fallback to empty chart
  return ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min'].map(name => ({
    name,
    hours: 0
  }))
})

// Calculate bar height as percentage relative to max value
const getBarHeight = (count: number): number => {
  const maxVal = Math.max(...weeklyActivity.value.map(d => d.hours), 1)
  const percentage = (count / maxVal) * 100
  return Math.min(percentage, 100) // Cap at 100%
}

// Recent activity from API
const recentActivity = computed(() => {
  if (!recentActivities.value || recentActivities.value.length === 0) {
    return [{
      id: '1',
      title: 'Belum ada aktivitas',
      course: 'Mulai belajar untuk melihat aktivitas',
      time: '-',
      icon: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
      iconBg: 'bg-neutral-100',
      iconColor: 'text-neutral-600'
    }]
  }
  
  return recentActivities.value.map(a => {
    const iconInfo = getActivityIcon(a.activity_type)
    return {
      id: a.id,
      title: a.description,
      course: a.metadata?.course_name || '',
      time: formatTimeAgo(a.created_at),
      icon: iconInfo.icon,
      iconBg: iconInfo.bgClass,
      iconColor: iconInfo.colorClass
    }
  })
})
</script>

