<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuiz, type Quiz, type Question, type QuizAttempt, type QuizResult, type QuizStatus } from '~/composables/useQuiz'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const route = useRoute()
const router = useRouter()
const lessonId = computed(() => route.params.id as string)

const { 
  quiz, 
  questions, 
  loading, 
  error,
  getQuizForStudent,
  startAttempt,
  submitAttempt,
  getQuizStatus
} = useQuiz()

// State
const currentQuiz = ref<Quiz | null>(null)
const quizStatusData = ref<QuizStatus | null>(null)
const quizAttempt = ref<QuizAttempt | null>(null)
const quizResult = ref<QuizResult | null>(null)
const previousAttempts = ref<QuizAttempt[]>([])
const currentQuestionIndex = ref(0)
const userAnswers = ref<Record<string, string[]>>({})
const textAnswers = ref<Record<string, string>>({})
const submitting = ref(false)
const timeRemaining = ref(0)
const timerInterval = ref<ReturnType<typeof setInterval> | null>(null)

// Computed for attempt limits
const remainingAttempts = computed(() => quizStatusData.value?.remaining_attempts ?? -1)
const canAttempt = computed(() => quizStatusData.value?.can_attempt ?? true)
const bestScore = computed(() => quizStatusData.value?.best_score)
const hasPassed = computed(() => quizStatusData.value?.has_passed ?? false)

// Quiz status
const quizViewStatus = computed(() => {
  if (quizResult.value) return 'result'
  if (quizAttempt.value) return 'in-progress'
  return 'ready'
})

// Current question
const currentQuestion = computed(() => {
  if (!currentQuiz.value?.questions) return null
  return currentQuiz.value.questions[currentQuestionIndex.value]
})

const totalQuestions = computed(() => currentQuiz.value?.questions?.length || 0)

const isFirstQuestion = computed(() => currentQuestionIndex.value === 0)
const isLastQuestion = computed(() => currentQuestionIndex.value === totalQuestions.value - 1)

// Check if all questions answered
const allQuestionsAnswered = computed(() => {
  if (!currentQuiz.value?.questions) return false
  return currentQuiz.value.questions.every((q: Question) => {
    if (q.question_type === 'short_answer') {
      return textAnswers.value[q.id]?.trim()
    }
    return userAnswers.value[q.id]?.length > 0
  })
})

// Get answered count
const answeredCount = computed(() => {
  if (!currentQuiz.value?.questions) return 0
  return currentQuiz.value.questions.filter((q: Question) => {
    if (q.question_type === 'short_answer') {
      return textAnswers.value[q.id]?.trim()
    }
    return userAnswers.value[q.id]?.length > 0
  }).length
})

// Format time
const formatTime = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

// Timer color
const timerColor = computed(() => {
  if (timeRemaining.value <= 60) return 'text-red-600'
  if (timeRemaining.value <= 300) return 'text-amber-600'
  return 'text-neutral-600'
})

// Timer key for localStorage
const timerKey = computed(() => `quiz_timer_${lessonId.value}`)
const answersKey = computed(() => `quiz_answers_${lessonId.value}`)

// Navigation
const goToQuestion = (index: number) => {
  if (index >= 0 && index < totalQuestions.value) {
    currentQuestionIndex.value = index
  }
}

const nextQuestion = () => {
  if (!isLastQuestion.value) {
    currentQuestionIndex.value++
  }
}

const prevQuestion = () => {
  if (!isFirstQuestion.value) {
    currentQuestionIndex.value--
  }
}

// Answer handling with localStorage save
const selectAnswer = (questionId: string, optionId: string) => {
  userAnswers.value[questionId] = [optionId]
  saveAnswersToLocal()
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
  saveAnswersToLocal()
}

// Save answers to localStorage
const saveAnswersToLocal = () => {
  if (quizAttempt.value) {
    localStorage.setItem(answersKey.value, JSON.stringify({
      attemptId: quizAttempt.value.id,
      userAnswers: userAnswers.value,
      textAnswers: textAnswers.value
    }))
  }
}

// Clear localStorage data
const clearLocalData = () => {
  localStorage.removeItem(timerKey.value)
  localStorage.removeItem(answersKey.value)
}

// Check if question is answered
const isQuestionAnswered = (questionId: string, type: string) => {
  if (type === 'short_answer') {
    return !!textAnswers.value[questionId]?.trim()
  }
  return (userAnswers.value[questionId]?.length || 0) > 0
}

// Start quiz
const handleStartQuiz = async () => {
  if (!currentQuiz.value || !canAttempt.value) return
  
  const attempt = await startAttempt(currentQuiz.value.id)
  if (attempt) {
    quizAttempt.value = attempt
    currentQuestionIndex.value = 0
    
    // Start timer if time limit exists
    if (currentQuiz.value.time_limit > 0) {
      timeRemaining.value = currentQuiz.value.time_limit * 60
      startTimer()
    }
  }
}

// Timer management with localStorage persistence
const startTimer = () => {
  if (timerInterval.value) clearInterval(timerInterval.value)
  
  timerInterval.value = setInterval(() => {
    timeRemaining.value--
    
    // Save remaining time to localStorage every 5 seconds
    if (timeRemaining.value % 5 === 0) {
      localStorage.setItem(timerKey.value, JSON.stringify({
        remaining: timeRemaining.value,
        attemptId: quizAttempt.value?.id,
        savedAt: Date.now()
      }))
    }
    
    if (timeRemaining.value <= 0) {
      handleSubmitQuiz()
    }
  }, 1000)
}

const stopTimer = () => {
  if (timerInterval.value) {
    clearInterval(timerInterval.value)
    timerInterval.value = null
  }
}

// Submit quiz
const handleSubmitQuiz = async () => {
  if (!quizAttempt.value || !currentQuiz.value) return
  
  stopTimer()
  submitting.value = true

  const answers = currentQuiz.value.questions?.map((q: Question) => ({
    question_id: q.id,
    selected_option_ids: userAnswers.value[q.id] || [],
    text_answer: textAnswers.value[q.id] || ''
  })) || []

  const result = await submitAttempt(quizAttempt.value.id, answers)
  if (result) {
    quizResult.value = result
    clearLocalData()
  }
  
  submitting.value = false
}

// Retry quiz
const handleRetryQuiz = async () => {
  quizAttempt.value = null
  quizResult.value = null
  userAnswers.value = {}
  textAnswers.value = {}
  currentQuestionIndex.value = 0
  
  // Refresh status data
  if (currentQuiz.value) {
    const status = await getQuizStatus(currentQuiz.value.id)
    if (status) {
      quizStatusData.value = status
      previousAttempts.value = status.attempts || []
    }
  }
}

// Back to course
const goBackToCourse = () => {
  const courseId = localStorage.getItem('lastCourseId')
  if (courseId) {
    router.push(`/dashboard/courses/${courseId}`)
  } else {
    router.push('/dashboard/courses')
  }
}

// Load quiz data
onMounted(async () => {
  if (lessonId.value) {
    const quizData = await getQuizForStudent(lessonId.value)
    if (quizData) {
      currentQuiz.value = quizData
      
      // Load quiz status
      const status = await getQuizStatus(quizData.id)
      if (status) {
        quizStatusData.value = status
        previousAttempts.value = status.attempts || []
        
        // Check for in-progress attempt
        if (status.in_progress_attempt) {
          quizAttempt.value = status.in_progress_attempt
          
          // Try to restore timer and answers from localStorage
          try {
            const savedTimer = localStorage.getItem(timerKey.value)
            if (savedTimer && quizData.time_limit > 0) {
              const timerData = JSON.parse(savedTimer)
              if (timerData.attemptId === status.in_progress_attempt.id) {
                // Calculate elapsed time since save
                const elapsed = Math.floor((Date.now() - timerData.savedAt) / 1000)
                timeRemaining.value = Math.max(0, timerData.remaining - elapsed)
                if (timeRemaining.value > 0) {
                  startTimer()
                } else {
                  // Time expired while away
                  handleSubmitQuiz()
                }
              }
            } else if (quizData.time_limit > 0) {
              // Start fresh timer
              const startedAt = new Date(status.in_progress_attempt.started_at).getTime()
              const elapsed = Math.floor((Date.now() - startedAt) / 1000)
              timeRemaining.value = Math.max(0, quizData.time_limit * 60 - elapsed)
              if (timeRemaining.value > 0) {
                startTimer()
              } else {
                handleSubmitQuiz()
              }
            }
            
            // Restore saved answers
            const savedAnswers = localStorage.getItem(answersKey.value)
            if (savedAnswers) {
              const answersData = JSON.parse(savedAnswers)
              if (answersData.attemptId === status.in_progress_attempt.id) {
                userAnswers.value = answersData.userAnswers || {}
                textAnswers.value = answersData.textAnswers || {}
              }
            }
          } catch (e) {
            console.error('Error restoring quiz state:', e)
          }
        }
      }
    }
  }
})

onUnmounted(() => {
  stopTimer()
})
</script>

<template>
  <div class="min-h-screen bg-neutral-50">
    <!-- Header -->
    <div class="bg-white border-b border-neutral-200 sticky top-0 z-10">
      <div class="max-w-4xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <!-- Back Button & Title -->
          <div class="flex items-center gap-4">
            <button @click="goBackToCourse" class="p-2 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
              </svg>
            </button>
            <div>
              <h1 class="text-lg font-semibold text-neutral-900">{{ currentQuiz?.title || 'Kuis' }}</h1>
              <p v-if="quizViewStatus === 'in-progress'" class="text-sm text-neutral-500">
                Pertanyaan {{ currentQuestionIndex + 1 }} dari {{ totalQuestions }}
              </p>
            </div>
          </div>

          <!-- Timer (if in progress) -->
          <div v-if="quizViewStatus === 'in-progress' && currentQuiz?.time_limit" class="flex items-center gap-2 px-4 py-2 bg-neutral-100 rounded-lg">
            <svg class="w-5 h-5" :class="timerColor" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <span class="font-mono font-semibold" :class="timerColor">{{ formatTime(timeRemaining) }}</span>
          </div>

          <!-- Progress Indicator -->
          <div v-if="quizViewStatus === 'in-progress'" class="hidden md:flex items-center gap-2 text-sm text-neutral-500">
            <span>{{ answeredCount }}/{{ totalQuestions }} dijawab</span>
          </div>
        </div>

        <!-- Progress Bar -->
        <div v-if="quizViewStatus === 'in-progress'" class="mt-3 h-1.5 bg-neutral-100 rounded-full overflow-hidden">
          <div 
            class="h-full bg-primary-500 rounded-full transition-all duration-300"
            :style="{ width: `${((currentQuestionIndex + 1) / totalQuestions) * 100}%` }"
          ></div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="max-w-4xl mx-auto px-4 py-8">
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="animate-spin w-10 h-10 border-4 border-primary-500 border-t-transparent rounded-full"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-20">
        <div class="w-16 h-16 mx-auto mb-4 bg-red-100 rounded-full flex items-center justify-center">
          <svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
          </svg>
        </div>
        <p class="text-neutral-600">{{ error }}</p>
        <button @click="goBackToCourse" class="mt-4 btn-primary">Kembali</button>
      </div>

      <!-- Ready State - Show Quiz Info -->
      <div v-else-if="quizViewStatus === 'ready' && currentQuiz" class="max-w-lg mx-auto">
        <div class="bg-white rounded-2xl shadow-sm border border-neutral-200 overflow-hidden">
          <!-- Quiz Icon -->
          <div class="p-8 bg-gradient-to-br from-primary-50 to-accent-50 flex items-center justify-center">
            <div class="w-24 h-24 bg-white rounded-2xl shadow-lg flex items-center justify-center">
              <svg class="w-12 h-12 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/>
              </svg>
            </div>
          </div>

          <!-- Quiz Info -->
          <div class="p-6">
            <h2 class="text-xl font-bold text-neutral-900 mb-2">{{ currentQuiz.title }}</h2>
            <p v-if="currentQuiz.description" class="text-neutral-600 mb-6">{{ currentQuiz.description }}</p>

            <!-- Quiz Stats -->
            <div class="grid grid-cols-2 gap-4 mb-6">
              <div class="p-4 bg-neutral-50 rounded-xl text-center">
                <p class="text-2xl font-bold text-primary-600">{{ currentQuiz.questions?.length || 0 }}</p>
                <p class="text-sm text-neutral-500">Pertanyaan</p>
              </div>
              <div class="p-4 bg-neutral-50 rounded-xl text-center">
                <p class="text-2xl font-bold text-primary-600">{{ currentQuiz.passing_score }}%</p>
                <p class="text-sm text-neutral-500">Nilai Minimum</p>
              </div>
              <div v-if="currentQuiz.time_limit" class="p-4 bg-neutral-50 rounded-xl text-center">
                <p class="text-2xl font-bold text-amber-600">{{ currentQuiz.time_limit }}</p>
                <p class="text-sm text-neutral-500">Menit</p>
              </div>
              <div v-if="currentQuiz.max_attempts" class="p-4 bg-neutral-50 rounded-xl text-center">
                <p class="text-2xl font-bold" :class="remainingAttempts === 0 ? 'text-red-600' : 'text-accent-600'">
                  {{ remainingAttempts === -1 ? '∞' : remainingAttempts }}
                </p>
                <p class="text-sm text-neutral-500">Sisa Percobaan</p>
              </div>
            </div>

            <!-- Best Score Badge -->
            <div v-if="bestScore !== undefined" class="mb-4 p-4 rounded-xl" :class="hasPassed ? 'bg-green-50 border border-green-200' : 'bg-amber-50 border border-amber-200'">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium" :class="hasPassed ? 'text-green-700' : 'text-amber-700'">
                  {{ hasPassed ? '✓ Sudah Lulus' : 'Nilai Terbaik' }}
                </span>
                <span class="text-lg font-bold" :class="hasPassed ? 'text-green-600' : 'text-amber-600'">
                  {{ bestScore }}%
                </span>
              </div>
            </div>

            <!-- Previous Attempts -->
            <div v-if="previousAttempts.length > 0" class="mb-6">
              <h3 class="text-sm font-semibold text-neutral-700 mb-3">Riwayat Percobaan</h3>
              <div class="space-y-2">
                <div 
                  v-for="(attempt, idx) in previousAttempts.slice(0, 3)" 
                  :key="attempt.id"
                  class="flex items-center justify-between p-3 bg-neutral-50 rounded-lg"
                >
                  <span class="text-sm text-neutral-600">Percobaan {{ previousAttempts.length - idx }}</span>
                  <div class="flex items-center gap-2">
                    <span 
                      class="px-2 py-1 rounded text-xs font-medium"
                      :class="attempt.passed ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
                    >
                      {{ attempt.score }}%
                    </span>
                    <span 
                      class="text-xs"
                      :class="attempt.passed ? 'text-green-600' : 'text-red-600'"
                    >
                      {{ attempt.passed ? 'Lulus' : 'Tidak Lulus' }}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Max Attempts Reached Warning -->
            <div v-if="!canAttempt" class="mb-4 p-4 bg-red-50 border border-red-200 rounded-xl">
              <div class="flex items-center gap-2 text-red-700">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                </svg>
                <span class="text-sm font-medium">Anda telah mencapai batas maksimal percobaan</span>
              </div>
            </div>

            <!-- Start Button -->
            <button 
              @click="handleStartQuiz"
              :disabled="loading || !canAttempt"
              class="w-full btn-primary py-3 text-lg"
              :class="{ 'opacity-50 cursor-not-allowed': !canAttempt }"
            >
              {{ !canAttempt ? 'Batas Percobaan Tercapai' : (previousAttempts.length > 0 ? 'Mulai Ulang Kuis' : 'Mulai Kuis') }}
            </button>
          </div>
        </div>
      </div>

      <!-- In Progress State - Show Question -->
      <div v-else-if="quizViewStatus === 'in-progress' && currentQuestion" class="max-w-2xl mx-auto">
        <div class="bg-white rounded-2xl shadow-sm border border-neutral-200 p-6">
          <!-- Question Number -->
          <div class="flex items-center gap-3 mb-6">
            <span class="flex-shrink-0 w-10 h-10 bg-primary-100 text-primary-700 rounded-full flex items-center justify-center text-lg font-bold">
              {{ currentQuestionIndex + 1 }}
            </span>
            <div class="flex-1">
              <span 
                class="px-2 py-1 text-xs font-medium rounded"
                :class="{
                  'bg-blue-100 text-blue-700': currentQuestion.question_type === 'multiple_choice',
                  'bg-purple-100 text-purple-700': currentQuestion.question_type === 'multiple_answer',
                  'bg-amber-100 text-amber-700': currentQuestion.question_type === 'true_false',
                  'bg-green-100 text-green-700': currentQuestion.question_type === 'short_answer'
                }"
              >
                {{ currentQuestion.question_type === 'multiple_choice' ? 'Pilihan Ganda' : 
                   currentQuestion.question_type === 'multiple_answer' ? 'Pilihan Berganda' :
                   currentQuestion.question_type === 'true_false' ? 'Benar/Salah' : 'Jawaban Singkat' }}
              </span>
            </div>
          </div>

          <!-- Question Text -->
          <h2 class="text-lg font-medium text-neutral-900 mb-6">{{ currentQuestion.question_text }}</h2>

          <!-- Multiple Choice / True False Options -->
          <div v-if="currentQuestion.question_type === 'multiple_choice' || currentQuestion.question_type === 'true_false'" class="space-y-3">
            <label 
              v-for="option in currentQuestion.options" 
              :key="option.id"
              class="flex items-center gap-4 p-4 border-2 rounded-xl cursor-pointer transition-all"
              :class="userAnswers[currentQuestion.id]?.includes(option.id!) 
                ? 'border-primary-500 bg-primary-50' 
                : 'border-neutral-200 hover:border-primary-200 hover:bg-neutral-50'"
            >
              <input 
                type="radio" 
                :name="`question-${currentQuestion.id}`"
                :value="option.id"
                :checked="userAnswers[currentQuestion.id]?.includes(option.id!)"
                @change="selectAnswer(currentQuestion.id, option.id!)"
                class="w-5 h-5 text-primary-600"
              />
              <span class="text-neutral-800">{{ option.option_text }}</span>
            </label>
          </div>

          <!-- Multiple Answer Options -->
          <div v-else-if="currentQuestion.question_type === 'multiple_answer'" class="space-y-3">
            <p class="text-sm text-neutral-500 mb-2">Pilih semua jawaban yang benar</p>
            <label 
              v-for="option in currentQuestion.options" 
              :key="option.id"
              class="flex items-center gap-4 p-4 border-2 rounded-xl cursor-pointer transition-all"
              :class="userAnswers[currentQuestion.id]?.includes(option.id!) 
                ? 'border-primary-500 bg-primary-50' 
                : 'border-neutral-200 hover:border-primary-200 hover:bg-neutral-50'"
            >
              <input 
                type="checkbox"
                :value="option.id"
                :checked="userAnswers[currentQuestion.id]?.includes(option.id!)"
                @change="toggleAnswer(currentQuestion.id, option.id!)"
                class="w-5 h-5 text-primary-600 rounded"
              />
              <span class="text-neutral-800">{{ option.option_text }}</span>
            </label>
          </div>

          <!-- Short Answer -->
          <div v-else-if="currentQuestion.question_type === 'short_answer'">
            <textarea 
              v-model="textAnswers[currentQuestion.id]"
              placeholder="Ketik jawaban Anda di sini..."
              rows="4"
              class="w-full px-4 py-3 border-2 border-neutral-200 rounded-xl focus:outline-none focus:border-primary-500 resize-none"
            ></textarea>
          </div>
        </div>

        <!-- Navigation -->
        <div class="flex items-center justify-between mt-6">
          <button 
            @click="prevQuestion"
            :disabled="isFirstQuestion"
            class="flex items-center gap-2 px-6 py-3 bg-white border border-neutral-200 rounded-xl text-neutral-700 hover:bg-neutral-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
            </svg>
            <span>Sebelumnya</span>
          </button>

          <!-- Question Bubbles -->
          <div class="hidden md:flex items-center gap-2">
            <button
              v-for="(q, idx) in currentQuiz?.questions"
              :key="q.id"
              @click="goToQuestion(idx)"
              class="w-8 h-8 rounded-full text-sm font-medium transition-all"
              :class="{
                'bg-primary-500 text-white': idx === currentQuestionIndex,
                'bg-green-100 text-green-700': idx !== currentQuestionIndex && isQuestionAnswered(q.id, q.question_type),
                'bg-neutral-100 text-neutral-500': idx !== currentQuestionIndex && !isQuestionAnswered(q.id, q.question_type)
              }"
            >
              {{ idx + 1 }}
            </button>
          </div>

          <div class="flex gap-3">
            <button 
              v-if="!isLastQuestion"
              @click="nextQuestion"
              class="flex items-center gap-2 px-6 py-3 bg-white border border-neutral-200 rounded-xl text-neutral-700 hover:bg-neutral-50"
            >
              <span>Selanjutnya</span>
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </button>

            <button 
              v-if="isLastQuestion || allQuestionsAnswered"
              @click="handleSubmitQuiz"
              :disabled="submitting"
              class="flex items-center gap-2 px-6 py-3 btn-primary"
            >
              <span v-if="submitting">Mengirim...</span>
              <span v-else>Kirim Jawaban</span>
              <svg v-if="!submitting" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Result State -->
      <div v-else-if="quizViewStatus === 'result' && quizResult" class="max-w-lg mx-auto">
        <div class="bg-white rounded-2xl shadow-sm border border-neutral-200 overflow-hidden">
          <!-- Result Header -->
          <div 
            class="p-8 flex flex-col items-center"
            :class="quizResult.passed ? 'bg-gradient-to-br from-green-50 to-accent-50' : 'bg-gradient-to-br from-red-50 to-amber-50'"
          >
            <div 
              class="w-28 h-28 rounded-full flex items-center justify-center mb-4"
              :class="quizResult.passed ? 'bg-white shadow-lg' : 'bg-white shadow-lg'"
            >
              <span 
                class="text-4xl font-bold"
                :class="quizResult.passed ? 'text-green-600' : 'text-red-600'"
              >
                {{ quizResult.score }}%
              </span>
            </div>
            <h2 
              class="text-2xl font-bold mb-2"
              :class="quizResult.passed ? 'text-green-700' : 'text-red-700'"
            >
              {{ quizResult.passed ? '🎉 Selamat, Anda Lulus!' : 'Belum Lulus' }}
            </h2>
            <p class="text-neutral-600">
              {{ quizResult.correct_count }} dari {{ quizResult.total_questions }} jawaban benar
            </p>
          </div>

          <!-- Result Stats -->
          <div class="p-6">
            <div class="grid grid-cols-2 gap-4 mb-6">
              <div class="p-4 bg-neutral-50 rounded-xl text-center">
                <p class="text-lg font-bold text-neutral-900">{{ quizResult.earned_points }}</p>
                <p class="text-sm text-neutral-500">Poin Diraih</p>
              </div>
              <div class="p-4 bg-neutral-50 rounded-xl text-center">
                <p class="text-lg font-bold text-neutral-900">{{ quizResult.total_points }}</p>
                <p class="text-sm text-neutral-500">Total Poin</p>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex gap-3">
              <button 
                @click="goBackToCourse"
                class="flex-1 py-3 px-4 bg-neutral-100 text-neutral-700 rounded-xl hover:bg-neutral-200 font-medium"
              >
                Kembali ke Kursus
              </button>
              <button 
                v-if="!quizResult.passed"
                @click="handleRetryQuiz"
                class="flex-1 btn-primary py-3"
              >
                Coba Lagi
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.btn-primary {
  @apply bg-primary-500 text-white font-medium rounded-xl hover:bg-primary-600 transition-colors px-4 py-2 disabled:opacity-50 disabled:cursor-not-allowed;
}
</style>
