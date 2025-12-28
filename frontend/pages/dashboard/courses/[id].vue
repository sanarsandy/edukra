<template>
  <div>
    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center min-h-[400px]">
      <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full"></div>
    </div>

    <div v-else-if="course">
      <!-- Breadcrumb -->
      <Breadcrumb :items="[{label: 'Kursus Saya', path: '/dashboard/courses'}, {label: course.title}]" class="mb-6" />
      
      <!-- Course Header -->
      <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden mb-8">
        <div class="h-48 md:h-56 relative bg-primary-600">
          <img v-if="course.thumbnail_url" :src="getThumbnailUrl(course.thumbnail_url)" class="w-full h-full object-cover" />
          <div v-else class="absolute inset-0 flex items-center justify-center">
            <svg class="w-24 h-24 text-white/20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
            </svg>
          </div>
          <div class="absolute top-4 left-4">
            <NuxtLink to="/dashboard/courses" class="flex items-center text-white/80 hover:text-white transition-colors bg-black/30 px-3 py-1.5 rounded-lg">
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 19l-7-7 7-7"/>
              </svg>
              Kembali
            </NuxtLink>
          </div>
        </div>
        <div class="p-6">
          <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
            <div>
              <span v-if="course.category" class="badge-primary mb-3">{{ course.category.name }}</span>
              <h1 class="text-2xl font-bold text-neutral-900 mb-2">{{ course.title }}</h1>
              <p class="text-neutral-500 mb-4">{{ course.description }}</p>
              <div class="flex flex-wrap items-center gap-4 text-sm text-neutral-500">
                <span class="flex items-center">
                  <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                  </svg>
                  {{ course.instructor?.full_name || 'Instructor' }}
                </span>
                <span class="flex items-center">
                  <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                  </svg>
                  {{ actualLessons.length }} Materi
                </span>
              </div>
            </div>
            <div class="flex-shrink-0">
              <div v-if="isEnrolled" class="text-center p-4 bg-neutral-50 rounded-xl">
                <div class="text-3xl font-bold text-primary-600 mb-1">{{ progressPercentage }}%</div>
                <div class="text-sm text-neutral-500">Selesai</div>
              </div>
              <div v-else class="text-center">
                <!-- Discounted Price -->
                <div v-if="course.discount_price && isDiscountActive(course)" class="mb-2">
                  <span class="text-sm text-neutral-400 line-through mr-1">
                    {{ formatPrice(course.price) }}
                  </span>
                  <span class="text-lg font-bold text-red-600">
                    {{ formatPrice(course.discount_price) }}
                  </span>
                </div>
                <!-- Normal Price -->
                <div v-else class="text-lg font-bold text-neutral-900 mb-2">{{ formatPrice(course.price) }}</div>
                <button @click="handleEnroll" class="btn-primary">{{ getEffectivePrice() > 0 ? 'Beli Sekarang' : 'Daftar Gratis' }}</button>
              </div>
            </div>
          </div>
        </div>
      </div>

    <!-- Content Grid -->
    <div class="grid lg:grid-cols-3 gap-6">
      <!-- Lessons List -->
      <div class="lg:col-span-2">
        <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
          <div class="p-5 border-b border-neutral-100 flex items-center justify-between">
            <h2 class="font-semibold text-neutral-900">Daftar Materi</h2>
            <div class="flex gap-2">
              <span 
                v-for="type in contentTypeFilters" 
                :key="type.value"
                @click="activeFilter = activeFilter === type.value ? '' : type.value"
                class="px-2 py-1 text-xs font-medium rounded-full cursor-pointer transition-colors"
                :class="activeFilter === type.value ? 'bg-primary-100 text-primary-700' : 'bg-neutral-100 text-neutral-600 hover:bg-neutral-200'"
              >
                {{ type.label }}
              </span>
            </div>
          </div>
          <div class="p-2">
            <template v-if="lessonsTree.length > 0">
              <StudentLessonTree
                v-for="(lesson, idx) in lessonsTree"
                :key="lesson.id"
                :lesson="lesson"
                :lesson-number="idx + 1"
                :selected-id="selectedLesson?.id"
                :completed-ids="completedLessonIds"
                @select="handleLessonSelect"
              />
            </template>
            <div v-else class="p-8 text-center text-neutral-400">
              Belum ada materi
            </div>
          </div>
        </div>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <!-- Content Preview -->
        <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
          <!-- Video Preview -->
          <div v-if="selectedLesson && selectedLesson.type === 'video'" class="aspect-video flex items-center justify-center relative bg-neutral-900">
            <!-- Embed Video (YouTube, Vimeo, Google Drive, etc) -->
            <template v-if="isEmbedVideo(selectedLesson.videoUrl)">
              <div v-if="!showVideoPlayer" class="text-center text-white cursor-pointer" @click="playVideo">
                <div class="w-20 h-20 mx-auto mb-3 rounded-full flex items-center justify-center transition-colors" :class="getVideoProviderColor(selectedLesson.videoUrl)">
                  <svg class="w-10 h-10 ml-1" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8 5v14l11-7z"/>
                  </svg>
                </div>
                <p class="text-sm opacity-80">Klik untuk memutar {{ getVideoProviderName(selectedLesson.videoUrl) }}</p>
              </div>
              <iframe 
                v-else
                :src="getEmbedUrl(selectedLesson.videoUrl)"
                class="w-full h-full"
                frameborder="0"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; fullscreen"
                allowfullscreen
              ></iframe>
            </template>
            <!-- Direct Video File -->
            <template v-else>
              <div v-if="!showVideoPlayer" class="text-center text-white cursor-pointer" @click="playVideo">
                <div class="w-20 h-20 mx-auto mb-3 bg-white/20 rounded-full flex items-center justify-center hover:bg-white/30 transition-colors">
                  <svg class="w-10 h-10 ml-1" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8 5v14l11-7z"/>
                  </svg>
                </div>
                <p class="text-sm opacity-80">Klik untuk memutar video</p>
              </div>
              <video 
                v-else
                ref="videoPlayer"
                :src="secureVideoUrl || getFileUrl(selectedLesson.videoUrl)"
                controls
                class="w-full h-full"
                @ended="onVideoEnded"
              ></video>
            </template>
          </div>

          <!-- PDF/Document Preview -->
          <div v-else-if="selectedLesson && (selectedLesson.type === 'pdf' || selectedLesson.type === 'document')" class="aspect-video flex items-center justify-center bg-neutral-100">
            <div class="text-center">
              <div class="w-20 h-24 mx-auto mb-4 bg-white rounded-lg shadow-md flex items-center justify-center relative">
                <svg class="w-10 h-10 text-red-500" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm-1 2l5 5h-5V4z"/>
                </svg>
                <span class="absolute -bottom-1 left-1/2 -translate-x-1/2 px-2 py-0.5 bg-red-500 text-white text-xs font-bold rounded">PDF</span>
              </div>
              <p class="text-sm text-neutral-600 font-medium">{{ selectedLesson.title }}</p>
            </div>
          </div>

          <!-- Quiz Preview -->
          <div v-else-if="selectedLesson && selectedLesson.type === 'quiz'" class="aspect-video flex items-center justify-center bg-gradient-to-br from-accent-50 to-accent-100">
            <div class="text-center">
              <div class="w-20 h-20 mx-auto mb-4 bg-accent-500 rounded-full flex items-center justify-center shadow-lg">
                <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                </svg>
              </div>
              <p class="text-lg text-accent-800 font-semibold">{{ selectedLesson.title }}</p>
              <p class="text-sm text-accent-600 mt-1">Uji pemahaman Anda</p>
            </div>
          </div>

          <!-- Text Preview -->
          <div v-else-if="selectedLesson && selectedLesson.type === 'text'" class="aspect-video flex items-center justify-center bg-neutral-50">
            <div class="text-center">
              <div class="w-16 h-16 mx-auto mb-4 bg-blue-100 rounded-full flex items-center justify-center">
                <svg class="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                </svg>
              </div>
              <p class="text-sm text-neutral-600 font-medium">{{ selectedLesson.title }}</p>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="aspect-video bg-neutral-100 flex items-center justify-center">
            <p class="text-neutral-400 text-sm">Pilih materi untuk memulai</p>
          </div>

          <!-- Content Info & Actions -->
          <div v-if="selectedLesson" class="p-4 border-t border-neutral-100">
            <h3 class="font-semibold text-neutral-900 mb-1">{{ selectedLesson.title }}</h3>
            <p class="text-xs text-neutral-500 mb-4">{{ selectedLesson.description || 'Tidak ada deskripsi' }}</p>
            
            <!-- Action Buttons based on type -->
            <div class="flex gap-2">
              <button v-if="selectedLesson.type === 'video'" @click="playVideo" class="flex-1 btn-primary text-sm py-2">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                </svg>
                Putar Video
              </button>
              <button v-else-if="selectedLesson.type === 'pdf' || selectedLesson.type === 'document'" @click="openDocument" class="flex-1 btn-primary text-sm py-2">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                </svg>
                Buka Dokumen
              </button>
              <button v-else-if="selectedLesson.type === 'quiz'" @click="openQuizModal" class="flex-1 btn-primary text-sm py-2">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                </svg>
                Mulai Kuis
              </button>
              <button v-else-if="selectedLesson.type === 'text'" @click="openTextContent" class="flex-1 btn-primary text-sm py-2">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                </svg>
                Baca Materi
              </button>
              
              <!-- Mark Complete Button -->
              <button 
                v-if="isEnrolled && !completedLessonIds.includes(selectedLesson.id)"
                @click="markLessonComplete"
                class="p-2 text-accent-600 bg-accent-50 rounded-lg hover:bg-accent-100 transition-colors"
                title="Tandai Selesai"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M5 13l4 4L19 7"/>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Progress Card -->
        <div class="bg-white rounded-xl border border-neutral-200 p-5">
          <h3 class="font-semibold text-neutral-900 mb-4">Progress Kursus</h3>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-neutral-500">Materi Selesai</span>
              <span class="font-medium text-neutral-900">{{ completedLessonIds.length }}/{{ actualLessons.length }}</span>
            </div>
            <div class="h-2 bg-neutral-100 rounded-full overflow-hidden">
              <div class="h-full bg-primary-500 rounded-full transition-all" :style="{ width: progressPercentage + '%' }"></div>
            </div>
          </div>
          
          <!-- Content Type Summary -->
          <div class="mt-4 pt-4 border-t border-neutral-100">
            <p class="text-xs text-neutral-500 mb-2">Jenis Materi</p>
            <div class="flex flex-wrap gap-2">
              <span v-for="(count, type) in contentTypeCounts" :key="type" class="px-2 py-1 text-xs rounded-full" :class="getContentTypeBadgeClass(type as string)">
                {{ count }} {{ getContentTypeLabel(type as string) }}
              </span>
            </div>
          </div>

          <button 
            v-if="isEnrolled && progressPercentage < 100 && nextUncompletedLesson"
            @click="selectLesson(nextUncompletedLesson)"
            class="w-full btn-primary mt-4"
          >
            Lanjutkan Belajar
          </button>
          <button 
            v-else-if="isEnrolled && progressPercentage >= 100"
            class="w-full py-2.5 text-sm font-medium text-accent-600 bg-accent-50 rounded-lg hover:bg-accent-100 transition-colors mt-4"
          >
            🎉 Kursus Selesai!
          </button>
          <button 
            v-else-if="!isEnrolled"
            @click="handleEnroll"
            class="w-full btn-primary mt-4"
          >
            {{ course.price > 0 ? 'Beli Kursus' : 'Daftar Gratis' }}
          </button>
        </div>

        <!-- Rating Card -->
        <div class="bg-white rounded-xl border border-neutral-200 p-5">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-neutral-900">Rating & Ulasan</h3>
            <div v-if="ratingStats" class="flex items-center gap-1">
              <svg class="w-5 h-5 text-yellow-400 fill-current" viewBox="0 0 24 24">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
              <span class="font-bold text-neutral-900">{{ ratingStats.average_rating.toFixed(1) }}</span>
              <span class="text-neutral-500 text-sm">({{ ratingStats.total_ratings }})</span>
            </div>
          </div>

          <!-- Submit/Edit Rating Form (only if enrolled) -->
          <div v-if="isEnrolled" class="mb-4">
            <p class="text-sm text-neutral-600 mb-2">{{ hasRated ? 'Rating Anda:' : 'Beri rating:' }}</p>
            <div class="flex items-center gap-1 mb-3">
              <button 
                v-for="star in 5" 
                :key="star"
                @click="selectedRating = star"
                class="p-1 transition-colors"
              >
                <svg 
                  class="w-7 h-7 transition-colors" 
                  :class="star <= selectedRating ? 'text-yellow-400 fill-current' : 'text-neutral-300 hover:text-yellow-300'"
                  viewBox="0 0 24 24"
                >
                  <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              </button>
            </div>
            <textarea 
              v-model="reviewText" 
              placeholder="Tulis ulasan (opsional)..."
              class="w-full px-3 py-2 text-sm border border-neutral-200 rounded-lg focus:ring-2 focus:ring-primary-500 resize-none"
              rows="2"
            ></textarea>
            <button 
              @click="handleSubmitRating"
              :disabled="selectedRating === 0 || ratingLoading"
              class="w-full btn-primary mt-2 disabled:opacity-50"
            >
              {{ hasRated ? 'Update Rating' : 'Kirim Rating' }}
            </button>
          </div>

          <!-- Reviews List -->
          <div v-if="courseRatings.length > 0" class="space-y-3 border-t border-neutral-100 pt-4">
            <div v-for="review in courseRatings.slice(0, 3)" :key="review.id" class="text-sm">
              <div class="flex items-center justify-between mb-1">
                <span class="font-medium text-neutral-900">{{ review.user_name }}</span>
                <div class="flex items-center">
                  <svg v-for="s in review.rating" :key="s" class="w-3 h-3 text-yellow-400 fill-current" viewBox="0 0 24 24">
                    <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                  </svg>
                </div>
              </div>
              <p v-if="review.review" class="text-neutral-600 line-clamp-2">{{ review.review }}</p>
              <p class="text-xs text-neutral-400 mt-1">{{ formatRatingTime(review.created_at) }}</p>
            </div>
          </div>
          <p v-else class="text-sm text-neutral-500 text-center py-3">Belum ada ulasan</p>
        </div>
      </div>
    </div>
    </div>

    <!-- Quiz Modal -->
    <Transition name="fade">
      <div v-if="showQuizModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4" @click="closeQuizModal">
        <div class="bg-white rounded-2xl w-full max-w-2xl max-h-[90vh] overflow-hidden flex flex-col" @click.stop>
          <!-- Quiz Header -->
          <div class="p-6 border-b border-neutral-100 flex items-center justify-between">
            <div>
              <h2 class="text-xl font-bold text-neutral-900">{{ currentQuiz?.title || 'Kuis' }}</h2>
              <p class="text-sm text-neutral-500 mt-1">{{ currentQuiz?.description }}</p>
            </div>
            <button @click="closeQuizModal" class="p-2 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- Quiz Content -->
          <div class="flex-1 overflow-y-auto p-6">
            <!-- Loading Quiz -->
            <div v-if="loadingQuiz" class="flex items-center justify-center py-12">
              <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full"></div>
            </div>

            <!-- Quiz Not Started -->
            <div v-else-if="!quizAttempt && currentQuiz" class="text-center py-8">
              <div class="w-20 h-20 mx-auto mb-6 bg-accent-100 rounded-full flex items-center justify-center">
                <svg class="w-10 h-10 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-neutral-900 mb-2">Siap Memulai Kuis?</h3>
              <div class="text-sm text-neutral-500 space-y-1 mb-6">
                <p>{{ currentQuiz.questions?.length || 0 }} Pertanyaan</p>
                <p v-if="currentQuiz.time_limit">Waktu: {{ currentQuiz.time_limit }} menit</p>
                <p>Nilai minimum: {{ currentQuiz.passing_score }}%</p>
              </div>
              <button @click="startQuiz" class="btn-primary px-8">
                Mulai Kuis
              </button>
            </div>

            <!-- Quiz Questions -->
            <div v-else-if="quizAttempt && currentQuiz?.questions" class="space-y-6">
              <div 
                v-for="(question, qIndex) in currentQuiz.questions" 
                :key="question.id"
                class="p-4 border border-neutral-200 rounded-xl"
              >
                <div class="flex items-start gap-3 mb-4">
                  <span class="flex-shrink-0 w-8 h-8 bg-primary-100 text-primary-700 rounded-full flex items-center justify-center text-sm font-semibold">
                    {{ qIndex + 1 }}
                  </span>
                  <p class="text-neutral-900 font-medium">{{ question.question_text }}</p>
                </div>
                
                <!-- Multiple Choice Options -->
                <div v-if="question.question_type === 'multiple_choice' || question.question_type === 'true_false'" class="space-y-2 ml-11">
                  <label 
                    v-for="option in question.options" 
                    :key="option.id"
                    class="flex items-center gap-3 p-3 border border-neutral-200 rounded-lg cursor-pointer hover:bg-neutral-50 transition-colors"
                    :class="{ 'border-primary-500 bg-primary-50': userAnswers[question.id]?.includes(option.id!) }"
                  >
                    <input 
                      type="radio" 
                      :name="`question-${question.id}`"
                      :value="option.id"
                      @change="selectAnswer(question.id, option.id!)"
                      class="w-4 h-4 text-primary-600"
                    />
                    <span class="text-sm text-neutral-700">{{ option.option_text }}</span>
                  </label>
                </div>

                <!-- Multiple Answer Options -->
                <div v-else-if="question.question_type === 'multiple_answer'" class="space-y-2 ml-11">
                  <label 
                    v-for="option in question.options" 
                    :key="option.id"
                    class="flex items-center gap-3 p-3 border border-neutral-200 rounded-lg cursor-pointer hover:bg-neutral-50 transition-colors"
                    :class="{ 'border-primary-500 bg-primary-50': userAnswers[question.id]?.includes(option.id!) }"
                  >
                    <input 
                      type="checkbox"
                      :value="option.id"
                      @change="toggleAnswer(question.id, option.id!)"
                      class="w-4 h-4 text-primary-600 rounded"
                    />
                    <span class="text-sm text-neutral-700">{{ option.option_text }}</span>
                  </label>
                </div>

                <!-- Short Answer -->
                <div v-else-if="question.question_type === 'short_answer'" class="ml-11">
                  <input 
                    type="text"
                    v-model="textAnswers[question.id]"
                    placeholder="Ketik jawaban Anda..."
                    class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
              </div>
            </div>

            <!-- Quiz Result -->
            <div v-else-if="quizResult" class="text-center py-8">
              <div 
                class="w-24 h-24 mx-auto mb-6 rounded-full flex items-center justify-center"
                :class="quizResult.passed ? 'bg-accent-100' : 'bg-red-100'"
              >
                <span class="text-3xl font-bold" :class="quizResult.passed ? 'text-accent-600' : 'text-red-600'">
                  {{ quizResult.score }}%
                </span>
              </div>
              <h3 class="text-xl font-bold mb-2" :class="quizResult.passed ? 'text-accent-600' : 'text-red-600'">
                {{ quizResult.passed ? '🎉 Selamat, Anda Lulus!' : 'Belum Lulus' }}
              </h3>
              <p class="text-neutral-500 mb-6">
                {{ quizResult.correct_count }}/{{ quizResult.total_questions }} jawaban benar
              </p>
              <div class="flex gap-3 justify-center">
                <button @click="closeQuizModal" class="px-6 py-2 text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200">
                  Tutup
                </button>
                <button v-if="!quizResult.passed" @click="retryQuiz" class="btn-primary px-6">
                  Coba Lagi
                </button>
              </div>
            </div>
          </div>

          <!-- Quiz Footer -->
          <div v-if="quizAttempt && !quizResult" class="p-4 border-t border-neutral-100 flex justify-end">
            <button @click="submitQuiz" :disabled="submittingQuiz" class="btn-primary px-8">
              <span v-if="submittingQuiz">Mengirim...</span>
              <span v-else>Kirim Jawaban</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Text Content Modal -->
    <Transition name="fade">
      <div v-if="showTextModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4" @click="showTextModal = false">
        <div class="bg-white rounded-2xl w-full max-w-3xl max-h-[90vh] overflow-hidden flex flex-col" @click.stop>
          <div class="p-6 border-b border-neutral-100 flex items-center justify-between">
            <h2 class="text-xl font-bold text-neutral-900">{{ selectedLesson?.title }}</h2>
            <button @click="showTextModal = false" class="p-2 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <div class="flex-1 overflow-y-auto p-6">
            <div class="prose max-w-none" v-html="selectedLesson?.content || 'Tidak ada konten'"></div>
          </div>
          <div class="p-4 border-t border-neutral-100 flex justify-end gap-2">
            <button 
              v-if="isEnrolled && !completedLessonIds.includes(selectedLesson?.id || '')"
              @click="markLessonComplete(); showTextModal = false"
              class="btn-primary"
            >
              Tandai Selesai
            </button>
            <button @click="showTextModal = false" class="px-4 py-2 text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200">
              Tutup
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Toast -->
    <Transition name="slide-up">
      <div v-if="toast.show" class="fixed bottom-6 right-6 z-50">
        <div class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3" :class="toast.type === 'success' ? 'bg-accent-600 text-white' : 'bg-red-600 text-white'">
          <svg v-if="toast.type === 'success'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
          <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
          <span class="text-sm font-medium">{{ toast.message }}</span>
        </div>
      </div>
    </Transition>

    <!-- AI Tutor Chat Widget -->
    <AIChatWidget 
      v-if="isEnrolled" 
      :course-id="courseId" 
    />
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const config = useRuntimeConfig()
const courseId = computed(() => route.params.id as string)
const apiBase = config.public.apiBase || 'http://localhost:8080'

// Get proper thumbnail URL - handle MinIO objects
const getThumbnailUrl = (url: string | null | undefined): string => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads')) return `${apiBase}${url}`
  return `${apiBase}/api/images/${url}`
}

// Helper: Check if discount is still active
const isDiscountActive = (courseData: any): boolean => {
  if (!courseData?.discount_price) return false
  if (!courseData.discount_valid_until) return true
  return new Date(courseData.discount_valid_until) > new Date()
}

// Helper: Format price to IDR currency
const formatPrice = (price: number): string => {
  if (price === 0) return 'Gratis'
  return new Intl.NumberFormat('id-ID', { 
    style: 'currency', 
    currency: 'IDR', 
    minimumFractionDigits: 0 
  }).format(price)
}

// Helper: Get the effective price (discounted or original)
const getEffectivePrice = (): number => {
  if (course.value && isDiscountActive(course.value)) {
    return course.value.discount_price
  }
  return course.value?.price || 0
}

const { course, fetchCourse, loading: loadingCourse } = useCourses()
const { checkEnrollment, enrollInCourse, updateProgress, loading: loadingEnrollment } = useEnrollments()
const { 
  quiz: currentQuiz, 
  questions: quizQuestions,
  loading: loadingQuiz, 
  getQuizForStudent, 
  startAttempt, 
  submitAttempt 
} = useQuiz()
const { 
  completedLessonIds, 
  fetchCourseProgress, 
  markLessonComplete: markLessonCompleteApi, 
  migrateFromLocalStorage,
  updateWatchTime 
} = useLessonProgress()

// Secure Content for MinIO
const { 
  getSecureVideoUrl, 
  getSecureDocumentUrl,
  getContentBlobUrl,
  isMinioObject, 
  isLegacyUpload,
  loading: loadingSecureContent,
  error: secureContentError
} = useSecureContent()

// Rating functionality
const { 
  ratings: courseRatings, 
  stats: ratingStats, 
  myRating, 
  hasRated,
  loading: ratingLoading,
  fetchRatings, 
  fetchStats, 
  fetchMyRating, 
  submitRating, 
  updateRating, 
  formatTimeAgo: formatRatingTime 
} = useRatings()

const selectedRating = ref(0)
const reviewText = ref('')

const handleSubmitRating = async () => {
  if (selectedRating.value === 0) return
  
  try {
    if (hasRated.value) {
      await updateRating(courseId.value, selectedRating.value, reviewText.value || undefined)
    } else {
      await submitRating(courseId.value, selectedRating.value, reviewText.value || undefined)
    }
    showToast('Rating berhasil disimpan!', 'success')
  } catch (err) {
    showToast('Gagal menyimpan rating', 'error')
  }
}

// Load ratings when course is loaded
const loadRatings = async () => {
  if (courseId.value) {
    await Promise.all([
      fetchRatings(courseId.value),
      fetchStats(courseId.value),
      fetchMyRating(courseId.value)
    ])
    // Set initial values from existing rating
    if (myRating.value) {
      selectedRating.value = myRating.value.rating
      reviewText.value = myRating.value.review || ''
    }
  }
}

// Helper functions for video URL handling
type VideoType = 'youtube' | 'vimeo' | 'gdrive' | 'zoom' | 'msstream' | 'sharepoint' | 'direct' | 'embed' | 'unknown'

const detectVideoType = (url: string | null | undefined): VideoType => {
  if (!url) return 'unknown'
  
  // YouTube
  if (url.includes('youtube.com') || url.includes('youtu.be')) return 'youtube'
  
  // Vimeo
  if (url.includes('vimeo.com')) return 'vimeo'
  
  // Google Drive
  if (url.includes('drive.google.com')) return 'gdrive'
  
  // Zoom Cloud Recording
  if (url.includes('zoom.us/rec')) return 'zoom'
  
  // Microsoft Stream
  if (url.includes('microsoftstream.com') || url.includes('web.microsoftstream.com')) return 'msstream'
  
  // SharePoint / OneDrive (Teams recordings often stored here)
  if (url.includes('sharepoint.com') || url.includes('onedrive.live.com') || url.includes('1drv.ms')) return 'sharepoint'
  
  // Direct video file (mp4, webm, ogg, mov)
  if (/\.(mp4|webm|ogg|mov|m4v)(\?.*)?$/i.test(url)) return 'direct'
  
  // Generic embed (iframe-able URLs)
  if (url.includes('/embed') || url.includes('/player')) return 'embed'
  
  // Default to direct for other URLs (will try <video> tag)
  return 'direct'
}

const isYouTubeUrl = (url: string | null | undefined): boolean => {
  return detectVideoType(url) === 'youtube'
}

const isEmbedVideo = (url: string | null | undefined): boolean => {
  const type = detectVideoType(url)
  return ['youtube', 'vimeo', 'gdrive', 'zoom', 'msstream', 'sharepoint', 'embed'].includes(type)
}

const getEmbedUrl = (url: string): string => {
  const type = detectVideoType(url)
  
  switch (type) {
    case 'youtube': {
      // Handle youtube.com/watch?v=VIDEO_ID
      const watchMatch = url.match(/youtube\.com\/watch\?v=([a-zA-Z0-9_-]+)/)
      if (watchMatch) return `https://www.youtube.com/embed/${watchMatch[1]}`
      // Handle youtu.be/VIDEO_ID
      const shortMatch = url.match(/youtu\.be\/([a-zA-Z0-9_-]+)/)
      if (shortMatch) return `https://www.youtube.com/embed/${shortMatch[1]}`
      // Already embed format
      if (url.includes('/embed/')) return url
      return url
    }
    
    case 'vimeo': {
      // Handle vimeo.com/VIDEO_ID
      const vimeoMatch = url.match(/vimeo\.com\/(\d+)/)
      if (vimeoMatch) return `https://player.vimeo.com/video/${vimeoMatch[1]}`
      // Already player format
      if (url.includes('player.vimeo.com')) return url
      return url
    }
    
    case 'gdrive': {
      // Handle drive.google.com/file/d/FILE_ID/view
      const driveMatch = url.match(/drive\.google\.com\/file\/d\/([a-zA-Z0-9_-]+)/)
      if (driveMatch) return `https://drive.google.com/file/d/${driveMatch[1]}/preview`
      return url
    }
    
    case 'zoom': {
      // Zoom recordings: zoom.us/rec/play/xxx or zoom.us/rec/share/xxx
      // Convert share link to embed format
      if (url.includes('/rec/share/')) {
        return url.replace('/rec/share/', '/rec/play/')
      }
      return url
    }
    
    case 'msstream': {
      // Microsoft Stream: web.microsoftstream.com/video/VIDEO_ID
      const streamMatch = url.match(/microsoftstream\.com\/video\/([a-zA-Z0-9-]+)/)
      if (streamMatch) return `https://web.microsoftstream.com/embed/video/${streamMatch[1]}`
      // Already embed format
      if (url.includes('/embed/')) return url
      return url
    }
    
    case 'sharepoint': {
      // SharePoint/OneDrive videos - add embed parameter
      // Format: https://xxx.sharepoint.com/:v:/xxx or direct link
      if (url.includes('sharepoint.com') && !url.includes('embed=1')) {
        const separator = url.includes('?') ? '&' : '?'
        return `${url}${separator}embed=1`
      }
      return url
    }
    
    default:
      return url
  }
}

const getVideoProviderName = (url: string | null | undefined): string => {
  const type = detectVideoType(url)
  switch (type) {
    case 'youtube': return 'YouTube'
    case 'vimeo': return 'Vimeo'
    case 'gdrive': return 'Google Drive'
    case 'zoom': return 'Zoom Recording'
    case 'msstream': return 'Microsoft Stream'
    case 'sharepoint': return 'Teams/SharePoint'
    default: return 'Video'
  }
}

const getVideoProviderColor = (url: string | null | undefined): string => {
  const type = detectVideoType(url)
  switch (type) {
    case 'youtube': return 'bg-red-600 hover:bg-red-700'
    case 'vimeo': return 'bg-blue-500 hover:bg-blue-600'
    case 'gdrive': return 'bg-green-600 hover:bg-green-700'
    case 'zoom': return 'bg-blue-600 hover:bg-blue-700'
    case 'msstream': return 'bg-purple-600 hover:bg-purple-700'
    case 'sharepoint': return 'bg-teal-600 hover:bg-teal-700'
    default: return 'bg-white/20 hover:bg-white/30'
  }
}

// Helper to get full URL for uploaded files
// Now supports both legacy /uploads paths and MinIO pre-signed URLs
const getFileUrl = (path: string | null | undefined): string => {
  if (!path) return ''
  // If it's already a full URL (pre-signed or external), return as-is
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path
  }
  // If it's a legacy local path (starts with /uploads), still serve from API
  // This maintains backward compatibility for existing content
  if (path.startsWith('/uploads')) {
    return `${apiBase}${path}`
  }
  // For MinIO objects (just object key), we need to get secure URL
  // This will be handled by getSecureVideoUrl/getSecureDocumentUrl
  // For now, return empty - the secure URL should be fetched separately
  return ''
}

// State for secure video URL (fetched from MinIO)
const secureVideoUrl = ref<string | null>(null)

const isEnrolled = ref(false)
const enrollmentData = ref<any>(null)
const activeFilter = ref('')
const selectedLesson = ref<any>(null)
// completedLessonIds is now from useLessonProgress composable

// Video
const showVideoPlayer = ref(false)
const videoPlayer = ref<HTMLVideoElement | null>(null)

// Quiz
const showQuizModal = ref(false)
const quizAttempt = ref<any>(null)
const quizResult = ref<any>(null)
const userAnswers = ref<Record<string, string[]>>({})
const textAnswers = ref<Record<string, string>>({})
const submittingQuiz = ref(false)

// Text Content
const showTextModal = ref(false)

// Toast
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

onMounted(async () => {
  await fetchCourse(courseId.value)
  const check = await checkEnrollment(courseId.value)
  isEnrolled.value = check.enrolled
  enrollmentData.value = check.enrollment
  
  // Migrate localStorage progress to server (one-time migration)
  await migrateFromLocalStorage(courseId.value)
  
  // Fetch progress from server
  if (isEnrolled.value) {
    await fetchCourseProgress(courseId.value)
  }
  
  // Check for payment status from return URL
  const route = useRoute()
  if (route.query.payment === 'pending') {
    toast.value = {
      show: true,
      message: 'Pembayaran sedang diproses. Mohon tunggu status terupdate.',
      type: 'success'
    }
    
    // Check enrollment again after a short delay
    setTimeout(async () => {
      const check = await checkEnrollment(courseId.value)
      if (check.enrolled) {
        isEnrolled.value = true
        enrollmentData.value = check.enrollment
        toast.value = {
          show: true,
          message: 'Pembayaran berhasil! Anda telah terdaftar.',
          type: 'success'
        }
      }
    }, 2000)
  }

  // Load ratings
  await loadRatings()
  
  // Select first lesson
  if (course.value?.lessons?.length) {
    selectedLesson.value = mapLesson(course.value.lessons[0])
  }
})

useHead({
  title: computed(() => course.value?.title ? `${course.value.title} - LearnHub` : 'Loading...')
})

const loading = computed(() => loadingCourse.value || loadingEnrollment.value)

const contentTypeFilters = [
  { value: 'video', label: 'Video' },
  { value: 'pdf', label: 'PDF' },
  { value: 'quiz', label: 'Kuis' }
]

const mapLesson = (l: any) => ({
  id: l.id,
  title: l.title,
  description: l.description,
  type: l.content_type,
  videoUrl: l.video_url,
  content: l.content,
  duration: l.video_duration ? `${Math.floor(l.video_duration / 60)} menit` : '-',
  order: l.order_index,
  is_container: l.is_container || false
})

// All lessons including containers (for backward compatibility)
const lessons = computed(() => {
  if (!course.value?.lessons) return []
  return course.value.lessons.map(mapLesson).sort((a: any, b: any) => a.order - b.order)
})

// Only actual lessons (excludes containers/modules) for progress tracking
const actualLessons = computed(() => {
  return lessons.value.filter((l: any) => !l.is_container)
})

const filteredLessons = computed(() => {
  if (!activeFilter.value) return lessons.value
  return lessons.value.filter((l: any) => l.type === activeFilter.value)
})

// Build tree structure from flat lessons
const lessonsTree = computed(() => {
  if (!course.value?.lessons) return []
  
  const allLessons = course.value.lessons
  const lessonMap = new Map<string, any>()
  const roots: any[] = []
  
  // First pass: create all nodes
  allLessons.forEach((l: any) => {
    lessonMap.set(l.id, { ...l, children: [] })
  })
  
  // Second pass: build tree
  allLessons.forEach((l: any) => {
    const node = lessonMap.get(l.id)!
    if (l.parent_id && lessonMap.has(l.parent_id)) {
      lessonMap.get(l.parent_id)!.children.push(node)
    } else {
      roots.push(node)
    }
  })
  
  // Sort children by order_index
  const sortChildren = (items: any[]) => {
    items.sort((a, b) => (a.order_index || 0) - (b.order_index || 0))
    items.forEach(item => {
      if (item.children?.length) sortChildren(item.children)
    })
  }
  sortChildren(roots)
  
  return roots
})

// Handle lesson selection from tree component
const handleLessonSelect = (lesson: any) => {
  selectedLesson.value = {
    id: lesson.id,
    title: lesson.title,
    description: lesson.description,
    type: lesson.content_type,
    videoUrl: lesson.video_url,
    content: lesson.content,
    duration: lesson.video_duration ? `${Math.floor(lesson.video_duration / 60)} menit` : '-',
    order: lesson.order_index
  }
  showVideoPlayer.value = false
}

const progressPercentage = computed(() => {
  if (!actualLessons.value.length) return 0
  return Math.round((completedLessonIds.value.length / actualLessons.value.length) * 100)
})

const nextUncompletedLesson = computed(() => {
  return actualLessons.value.find((l: any) => !completedLessonIds.value.includes(l.id))
})

const contentTypeCounts = computed(() => {
  const counts: Record<string, number> = {}
  actualLessons.value.forEach((l: any) => {
    counts[l.type] = (counts[l.type] || 0) + 1
  })
  return counts
})

const selectLesson = (lesson: any) => {
  selectedLesson.value = lesson
  showVideoPlayer.value = false
}

const handleEnroll = async () => {
  // If course is paid, redirect to checkout page
  if (course.value && course.value.price > 0) {
    navigateTo(`/checkout/${courseId.value}`)
    return
  }
  
  // Free course - enroll directly
  const result = await enrollInCourse(courseId.value)
  if (result) {
    isEnrolled.value = true
    enrollmentData.value = result
    showToast('Berhasil mendaftar kursus!')
  }
}

// Video Functions
const playVideo = async () => {
  if (!isEnrolled.value) {
    showToast('Silakan daftar kursus terlebih dahulu', 'error')
    return
  }
  
  // For embed videos (YouTube, Vimeo, GDrive), proceed directly
  if (selectedLesson.value && isEmbedVideo(selectedLesson.value.videoUrl)) {
    showVideoPlayer.value = true
    setTimeout(() => {
      if (selectedLesson.value && !completedLessonIds.value.includes(selectedLesson.value.id)) {
        markLessonComplete()
      }
    }, 30000)
    return
  }
  
  // For MinIO objects, fetch content via backend stream
  if (selectedLesson.value && isMinioObject(selectedLesson.value.videoUrl)) {
    // Use blob URL approach - this fetches content through backend proxy
    const url = await getContentBlobUrl(selectedLesson.value.id)
    if (url) {
      secureVideoUrl.value = url
      showVideoPlayer.value = true
      nextTick(() => {
        videoPlayer.value?.play()
      })
    } else {
      showToast(secureContentError.value || 'Gagal memuat video', 'error')
    }
    return
  }
  
  // For legacy /uploads or direct files, use old behavior
  showVideoPlayer.value = true
  nextTick(() => {
    videoPlayer.value?.play()
  })
}

const onVideoEnded = () => {
  if (selectedLesson.value && !completedLessonIds.value.includes(selectedLesson.value.id)) {
    markLessonComplete()
  }
}

// Document Functions
const openDocument = async () => {
  if (!isEnrolled.value) {
    showToast('Silakan daftar kursus terlebih dahulu', 'error')
    return
  }
  
  if (!selectedLesson.value?.videoUrl) {
    showToast('URL dokumen tidak tersedia', 'error')
    return
  }
  
  let docUrl: string
  
  // For MinIO objects, fetch content via backend stream proxy
  if (isMinioObject(selectedLesson.value.videoUrl)) {
    // Use blob URL approach - streamed through backend to avoid signature issues
    const url = await getContentBlobUrl(selectedLesson.value.id)
    if (!url) {
      showToast(secureContentError.value || 'Gagal memuat dokumen', 'error')
      return
    }
    docUrl = url
  } else {
    // Legacy or direct URL
    docUrl = getFileUrl(selectedLesson.value.videoUrl)
  }
  
  window.open(docUrl, '_blank')
  showToast('Membuka dokumen...')
  
  // Auto mark as completed after opening
  if (!completedLessonIds.value.includes(selectedLesson.value.id)) {
    setTimeout(() => markLessonComplete(), 2000)
  }
}

// Text Content Functions
const openTextContent = () => {
  if (!isEnrolled.value) {
    showToast('Silakan daftar kursus terlebih dahulu', 'error')
    return
  }
  showTextModal.value = true
}

// Quiz Functions
const openQuizModal = async () => {
  if (!isEnrolled.value) {
    showToast('Silakan daftar kursus terlebih dahulu', 'error')
    return
  }
  
  // Navigate to dedicated quiz page
  if (selectedLesson.value) {
    // Save current course ID for "back to course" navigation
    localStorage.setItem('lastCourseId', courseId.value)
    navigateTo(`/dashboard/quiz/${selectedLesson.value.id}?lessonId=${selectedLesson.value.id}`)
  }
}

const closeQuizModal = () => {
  showQuizModal.value = false
  quizAttempt.value = null
  quizResult.value = null
}

const startQuiz = async () => {
  if (!currentQuiz.value) return
  const attempt = await startAttempt(currentQuiz.value.id)
  if (attempt) {
    quizAttempt.value = attempt
  }
}

const selectAnswer = (questionId: string, optionId: string) => {
  userAnswers.value[questionId] = [optionId]
}

const toggleAnswer = (questionId: string, optionId: string) => {
  if (!userAnswers.value[questionId]) {
    userAnswers.value[questionId] = []
  }
  const idx = userAnswers.value[questionId].indexOf(optionId)
  if (idx === -1) {
    userAnswers.value[questionId].push(optionId)
  } else {
    userAnswers.value[questionId].splice(idx, 1)
  }
}

const submitQuiz = async () => {
  if (!quizAttempt.value) return
  
  submittingQuiz.value = true
  
  const answers = Object.entries(userAnswers.value).map(([questionId, optionIds]) => ({
    question_id: questionId,
    selected_option_ids: optionIds
  }))
  
  // Add text answers
  Object.entries(textAnswers.value).forEach(([questionId, text]) => {
    if (text) {
      answers.push({
        question_id: questionId,
        text_answer: text
      } as any)
    }
  })
  
  const result = await submitAttempt(quizAttempt.value.id, answers)
  submittingQuiz.value = false
  
  if (result) {
    quizResult.value = result
    if (result.passed && selectedLesson.value) {
      markLessonComplete()
    }
  }
}

const retryQuiz = () => {
  quizAttempt.value = null
  quizResult.value = null
  userAnswers.value = {}
  textAnswers.value = {}
}

// Progress Functions
const markLessonComplete = async () => {
  if (!selectedLesson.value || completedLessonIds.value.includes(selectedLesson.value.id)) return
  
  // Save to server (no more localStorage)
  await markLessonCompleteApi(selectedLesson.value.id)
  
  // Update enrollment progress
  if (enrollmentData.value?.id) {
    await updateProgress(enrollmentData.value.id, progressPercentage.value)
  }
  
  showToast('Materi ditandai selesai!')
  
  // Auto-select next lesson
  if (nextUncompletedLesson.value) {
    setTimeout(() => {
      selectLesson(nextUncompletedLesson.value)
    }, 500)
  }
}

// Helpers
const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const getContentTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    video: 'Video',
    pdf: 'PDF',
    text: 'Teks',
    quiz: 'Kuis'
  }
  return labels[type] || type
}

const getContentTypeBadgeClass = (type: string) => {
  const classes: Record<string, string> = {
    video: 'bg-primary-100 text-primary-700',
    pdf: 'bg-red-100 text-red-700',
    text: 'bg-blue-100 text-blue-700',
    quiz: 'bg-accent-100 text-accent-700'
  }
  return classes[type] || 'bg-neutral-100 text-neutral-600'
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { opacity: 0; transform: translateY(20px); }
</style>
