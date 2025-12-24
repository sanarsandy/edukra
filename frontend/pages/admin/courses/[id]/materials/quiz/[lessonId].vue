<template>
  <div class="min-h-screen bg-neutral-50">
    <!-- Header -->
    <div class="bg-white border-b border-neutral-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <NuxtLink :to="`/admin/courses/${courseId}/materials`" class="p-2 hover:bg-neutral-100 rounded-lg transition-colors">
              <svg class="w-5 h-5 text-neutral-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
              </svg>
            </NuxtLink>
            <div>
              <h1 class="text-xl font-bold text-neutral-900">Editor Kuis</h1>
              <p class="text-sm text-neutral-500">{{ lesson?.title }}</p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <button v-if="quiz" @click="showSettingsModal = true" class="px-4 py-2 text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors flex items-center gap-2">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
              </svg>
              Pengaturan
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
      </div>

      <!-- No Quiz Yet - Create -->
      <div v-else-if="!quiz" class="bg-white rounded-xl border border-neutral-200 p-8 text-center">
        <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">Belum ada kuis</h3>
        <p class="text-neutral-500 mb-6">Buat kuis baru untuk materi ini</p>
        <button @click="showCreateModal = true" class="px-6 py-3 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors font-medium">
          Buat Kuis
        </button>
      </div>

      <!-- Quiz Editor -->
      <div v-else class="space-y-6">
        <!-- Quiz Info Card -->
        <div class="bg-white rounded-xl border border-neutral-200 p-6">
          <div class="flex items-start justify-between">
            <div>
              <h2 class="text-xl font-bold text-neutral-900">{{ quiz.title }}</h2>
              <p class="text-neutral-500 mt-1">{{ quiz.description || 'Tidak ada deskripsi' }}</p>
            </div>
            <div class="flex items-center gap-4 text-sm text-neutral-600">
              <div class="flex items-center gap-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                {{ quiz.time_limit > 0 ? `${quiz.time_limit} menit` : 'Tidak terbatas' }}
              </div>
              <div class="flex items-center gap-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                Lulus: {{ quiz.passing_score }}%
              </div>
              <div class="flex items-center gap-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
                </svg>
                {{ quiz.max_attempts > 0 ? `${quiz.max_attempts}x percobaan` : 'Unlimited' }}
              </div>
            </div>
          </div>
        </div>

        <!-- Questions List -->
        <div class="bg-white rounded-xl border border-neutral-200">
          <div class="p-4 border-b border-neutral-200 flex items-center justify-between">
            <h3 class="font-semibold text-neutral-900">Pertanyaan ({{ questions.length }})</h3>
            <button @click="openAddQuestion" class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors text-sm font-medium flex items-center gap-2">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
              </svg>
              Tambah Pertanyaan
            </button>
          </div>

          <!-- Empty State -->
          <div v-if="questions.length === 0" class="p-8 text-center text-neutral-500">
            <p>Belum ada pertanyaan. Klik tombol di atas untuk menambahkan.</p>
          </div>

          <!-- Question Items -->
          <div v-else class="divide-y divide-neutral-200">
            <div v-for="(question, index) in questions" :key="question.id" class="p-4 hover:bg-neutral-50 transition-colors">
              <div class="flex items-start gap-4">
                <div class="flex-shrink-0 w-8 h-8 bg-neutral-100 rounded-full flex items-center justify-center text-sm font-medium text-neutral-600">
                  {{ index + 1 }}
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="px-2 py-0.5 rounded text-xs font-medium" :class="getQuestionTypeBadge(question.question_type)">
                      {{ getQuestionTypeLabel(question.question_type) }}
                    </span>
                    <span class="text-xs text-neutral-500">{{ question.points }} poin</span>
                  </div>
                  <p class="text-neutral-900 font-medium">{{ question.question_text }}</p>
                  <div v-if="question.options && question.options.length > 0" class="mt-2 space-y-1">
                    <div v-for="option in question.options" :key="option.id" class="flex items-center gap-2 text-sm">
                      <span :class="option.is_correct ? 'text-green-500' : 'text-neutral-400'">
                        {{ option.is_correct ? '✓' : '○' }}
                      </span>
                      <span :class="option.is_correct ? 'text-green-700 font-medium' : 'text-neutral-600'">
                        {{ option.option_text }}
                      </span>
                    </div>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <button @click="openEditQuestion(question)" class="p-2 hover:bg-neutral-200 rounded-lg transition-colors">
                    <svg class="w-4 h-4 text-neutral-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                    </svg>
                  </button>
                  <button @click="confirmDeleteQuestion(question)" class="p-2 hover:bg-red-100 rounded-lg transition-colors">
                    <svg class="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Quiz Modal -->
    <Transition name="fade">
      <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showCreateModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-md m-4" @click.stop>
          <h2 class="text-xl font-bold text-neutral-900 mb-6">Buat Kuis Baru</h2>
          <form @submit.prevent="handleCreateQuiz" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Judul Kuis *</label>
              <input v-model="quizForm.title" type="text" required class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" placeholder="Contoh: Kuis Bab 1" />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Deskripsi</label>
              <textarea v-model="quizForm.description" rows="2" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" placeholder="Deskripsi singkat kuis"></textarea>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Batas Waktu (menit)</label>
                <input v-model.number="quizForm.time_limit" type="number" min="0" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" placeholder="0 = tidak terbatas" />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Nilai Kelulusan (%)</label>
                <input v-model.number="quizForm.passing_score" type="number" min="0" max="100" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Maksimal Percobaan</label>
              <input v-model.number="quizForm.max_attempts" type="number" min="0" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" placeholder="0 = tidak terbatas" />
            </div>
            <div class="flex items-center gap-4">
              <label class="flex items-center gap-2 cursor-pointer">
                <input v-model="quizForm.show_correct_answers" type="checkbox" class="w-4 h-4 text-red-500 rounded" />
                <span class="text-sm text-neutral-700">Tampilkan jawaban benar setelah submit</span>
              </label>
            </div>
            <div class="flex justify-end gap-3 pt-4">
              <button type="button" @click="showCreateModal = false" class="px-4 py-2 text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors">Batal</button>
              <button type="submit" :disabled="saving" class="px-6 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors disabled:opacity-50">
                {{ saving ? 'Menyimpan...' : 'Buat Kuis' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- Settings Modal -->
    <Transition name="fade">
      <div v-if="showSettingsModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showSettingsModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-md m-4" @click.stop>
          <h2 class="text-xl font-bold text-neutral-900 mb-6">Pengaturan Kuis</h2>
          <form @submit.prevent="handleUpdateSettings" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Judul Kuis *</label>
              <input v-model="settingsForm.title" type="text" required class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Deskripsi</label>
              <textarea v-model="settingsForm.description" rows="2" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent"></textarea>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Batas Waktu (menit)</label>
                <input v-model.number="settingsForm.time_limit" type="number" min="0" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Nilai Kelulusan (%)</label>
                <input v-model.number="settingsForm.passing_score" type="number" min="0" max="100" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Maksimal Percobaan</label>
              <input v-model.number="settingsForm.max_attempts" type="number" min="0" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" />
            </div>
            <div class="space-y-2">
              <label class="flex items-center gap-2 cursor-pointer">
                <input v-model="settingsForm.shuffle_questions" type="checkbox" class="w-4 h-4 text-red-500 rounded" />
                <span class="text-sm text-neutral-700">Acak urutan pertanyaan</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input v-model="settingsForm.shuffle_options" type="checkbox" class="w-4 h-4 text-red-500 rounded" />
                <span class="text-sm text-neutral-700">Acak urutan pilihan jawaban</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input v-model="settingsForm.show_correct_answers" type="checkbox" class="w-4 h-4 text-red-500 rounded" />
                <span class="text-sm text-neutral-700">Tampilkan jawaban benar setelah submit</span>
              </label>
            </div>
            <div class="flex justify-between pt-4">
              <button type="button" @click="confirmDeleteQuiz" class="px-4 py-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors">Hapus Kuis</button>
              <div class="flex gap-3">
                <button type="button" @click="showSettingsModal = false" class="px-4 py-2 text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors">Batal</button>
                <button type="submit" :disabled="saving" class="px-6 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors disabled:opacity-50">
                  {{ saving ? 'Menyimpan...' : 'Simpan' }}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- Question Modal -->
    <Transition name="fade">
      <div v-if="showQuestionModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showQuestionModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-2xl m-4 max-h-[90vh] overflow-y-auto" @click.stop>
          <h2 class="text-xl font-bold text-neutral-900 mb-6">{{ editingQuestion ? 'Edit Pertanyaan' : 'Tambah Pertanyaan' }}</h2>
          <form @submit.prevent="handleSaveQuestion" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Tipe Pertanyaan *</label>
              <select v-model="questionForm.question_type" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent">
                <option value="multiple_choice">Pilihan Ganda (Satu Jawaban)</option>
                <option value="multiple_answer">Pilihan Ganda (Multi Jawaban)</option>
                <option value="true_false">Benar/Salah</option>
                <option value="short_answer">Jawaban Singkat</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Pertanyaan *</label>
              <textarea v-model="questionForm.question_text" rows="3" required class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" placeholder="Tulis pertanyaan di sini..."></textarea>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Poin</label>
                <input v-model.number="questionForm.points" type="number" min="1" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" />
              </div>
              <div class="flex items-end">
                <label class="flex items-center gap-2 cursor-pointer pb-2">
                  <input v-model="questionForm.required" type="checkbox" class="w-4 h-4 text-red-500 rounded" />
                  <span class="text-sm text-neutral-700">Wajib dijawab</span>
                </label>
              </div>
            </div>

            <!-- Options for Multiple Choice / True False -->
            <div v-if="questionForm.question_type !== 'short_answer'">
              <label class="block text-sm font-medium text-neutral-700 mb-2">Pilihan Jawaban</label>
              <div class="space-y-2">
                <div v-for="(option, index) in questionForm.options" :key="index" class="flex items-center gap-2">
                  <input 
                    v-if="questionForm.question_type === 'multiple_choice' || questionForm.question_type === 'true_false'"
                    type="radio" 
                    :name="'correct-option'" 
                    :checked="option.is_correct"
                    @change="setCorrectOption(index)"
                    class="w-4 h-4 text-green-500"
                  />
                  <input 
                    v-else
                    type="checkbox" 
                    v-model="option.is_correct"
                    class="w-4 h-4 text-green-500 rounded"
                  />
                  <input v-model="option.option_text" type="text" class="flex-1 px-3 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" placeholder="Teks pilihan" />
                  <button v-if="questionForm.options.length > 2" type="button" @click="removeOption(index)" class="p-2 text-red-500 hover:bg-red-50 rounded-lg">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                    </svg>
                  </button>
                </div>
              </div>
              <button v-if="questionForm.question_type !== 'true_false'" type="button" @click="addOption" class="mt-2 text-sm text-red-500 hover:text-red-600">
                + Tambah pilihan
              </button>
            </div>

            <!-- Correct Answer for Short Answer -->
            <div v-else>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Jawaban Benar</label>
              <input v-model="questionForm.options[0].option_text" type="text" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" placeholder="Masukkan jawaban yang benar" />
              <p class="text-xs text-neutral-500 mt-1">Jawaban akan dicocokkan secara case-insensitive</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Penjelasan (opsional)</label>
              <textarea v-model="questionForm.explanation" rows="2" class="w-full px-4 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent" placeholder="Penjelasan yang ditampilkan setelah menjawab"></textarea>
            </div>

            <div class="flex justify-end gap-3 pt-4">
              <button type="button" @click="showQuestionModal = false" class="px-4 py-2 text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors">Batal</button>
              <button type="submit" :disabled="saving" class="px-6 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors disabled:opacity-50">
                {{ saving ? 'Menyimpan...' : 'Simpan' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- Toast -->
    <Transition name="slide-up">
      <div v-if="toast.show" class="fixed bottom-4 right-4 z-50">
        <div :class="['px-4 py-3 rounded-lg shadow-lg', toast.type === 'success' ? 'bg-green-500' : 'bg-red-500', 'text-white']">
          {{ toast.message }}
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { useQuiz, type Question } from '~/composables/useQuiz'

definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const route = useRoute()
const lessonId = computed(() => route.params.lessonId as string)
const courseId = computed(() => route.params.id as string)

useHead({
  title: 'Editor Kuis - Admin'
})

const { quiz, questions, loading, getQuizByLesson, createQuiz, updateQuiz, deleteQuiz, addQuestion, updateQuestion, deleteQuestion } = useQuiz()
const { fetchLesson, lesson } = useLessons()

const saving = ref(false)
const showCreateModal = ref(false)
const showSettingsModal = ref(false)
const showQuestionModal = ref(false)
const editingQuestion = ref<Question | null>(null)

const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

const quizForm = ref({
  title: '',
  description: '',
  time_limit: 0,
  passing_score: 70,
  max_attempts: 0,
  show_correct_answers: true
})

const settingsForm = ref({
  title: '',
  description: '',
  time_limit: 0,
  passing_score: 70,
  max_attempts: 0,
  shuffle_questions: false,
  shuffle_options: false,
  show_correct_answers: true
})

const questionForm = ref({
  question_type: 'multiple_choice' as string,
  question_text: '',
  explanation: '',
  points: 10,
  required: true,
  options: [
    { option_text: '', is_correct: true },
    { option_text: '', is_correct: false },
    { option_text: '', is_correct: false },
    { option_text: '', is_correct: false }
  ]
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const getQuestionTypeBadge = (type: string) => {
  const badges: Record<string, string> = {
    'multiple_choice': 'bg-blue-100 text-blue-700',
    'multiple_answer': 'bg-purple-100 text-purple-700',
    'true_false': 'bg-green-100 text-green-700',
    'short_answer': 'bg-orange-100 text-orange-700'
  }
  return badges[type] || 'bg-neutral-100 text-neutral-700'
}

const getQuestionTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    'multiple_choice': 'Pilihan Ganda',
    'multiple_answer': 'Multi Jawaban',
    'true_false': 'Benar/Salah',
    'short_answer': 'Jawaban Singkat'
  }
  return labels[type] || type
}

const handleCreateQuiz = async () => {
  saving.value = true
  const result = await createQuiz(lessonId.value, quizForm.value)
  saving.value = false
  if (result) {
    showCreateModal.value = false
    showToast('Kuis berhasil dibuat')
    // Copy to settings form
    settingsForm.value = { ...result } as any
  } else {
    showToast('Gagal membuat kuis', 'error')
  }
}

const handleUpdateSettings = async () => {
  if (!quiz.value) return
  saving.value = true
  const result = await updateQuiz(quiz.value.id, settingsForm.value)
  saving.value = false
  if (result) {
    showSettingsModal.value = false
    showToast('Pengaturan berhasil disimpan')
  } else {
    showToast('Gagal menyimpan pengaturan', 'error')
  }
}

const confirmDeleteQuiz = async () => {
  if (!quiz.value) return
  if (!confirm('Apakah Anda yakin ingin menghapus kuis ini? Semua pertanyaan juga akan dihapus.')) return
  
  saving.value = true
  const result = await deleteQuiz(quiz.value.id)
  saving.value = false
  if (result) {
    showSettingsModal.value = false
    showToast('Kuis berhasil dihapus')
  } else {
    showToast('Gagal menghapus kuis', 'error')
  }
}

const openAddQuestion = () => {
  editingQuestion.value = null
  questionForm.value = {
    question_type: 'multiple_choice',
    question_text: '',
    explanation: '',
    points: 10,
    required: true,
    options: [
      { option_text: '', is_correct: true },
      { option_text: '', is_correct: false },
      { option_text: '', is_correct: false },
      { option_text: '', is_correct: false }
    ]
  }
  showQuestionModal.value = true
}

const openEditQuestion = (question: Question) => {
  editingQuestion.value = question
  questionForm.value = {
    question_type: question.question_type,
    question_text: question.question_text,
    explanation: question.explanation || '',
    points: question.points,
    required: question.required,
    options: question.options?.map(o => ({ option_text: o.option_text, is_correct: o.is_correct })) || []
  }
  showQuestionModal.value = true
}

const handleSaveQuestion = async () => {
  if (!quiz.value) return
  
  saving.value = true
  
  // Prepare options based on type
  let options = questionForm.value.options
  if (questionForm.value.question_type === 'true_false') {
    options = [
      { option_text: 'Benar', is_correct: questionForm.value.options[0]?.is_correct || false },
      { option_text: 'Salah', is_correct: !questionForm.value.options[0]?.is_correct }
    ]
  } else if (questionForm.value.question_type === 'short_answer') {
    options = [{ option_text: questionForm.value.options[0]?.option_text || '', is_correct: true }]
  }
  
  const data = {
    question_type: questionForm.value.question_type,
    question_text: questionForm.value.question_text,
    explanation: questionForm.value.explanation,
    points: questionForm.value.points,
    required: questionForm.value.required,
    options
  }
  
  let result
  if (editingQuestion.value) {
    result = await updateQuestion(editingQuestion.value.id, data)
  } else {
    result = await addQuestion(quiz.value.id, data)
  }
  
  saving.value = false
  
  if (result) {
    showQuestionModal.value = false
    showToast(editingQuestion.value ? 'Pertanyaan berhasil diperbarui' : 'Pertanyaan berhasil ditambahkan')
  } else {
    showToast('Gagal menyimpan pertanyaan', 'error')
  }
}

const confirmDeleteQuestion = async (question: Question) => {
  if (!confirm('Hapus pertanyaan ini?')) return
  
  const result = await deleteQuestion(question.id)
  if (result) {
    showToast('Pertanyaan berhasil dihapus')
  } else {
    showToast('Gagal menghapus pertanyaan', 'error')
  }
}

const setCorrectOption = (index: number) => {
  questionForm.value.options.forEach((opt, i) => {
    opt.is_correct = i === index
  })
}

const addOption = () => {
  questionForm.value.options.push({ option_text: '', is_correct: false })
}

const removeOption = (index: number) => {
  questionForm.value.options.splice(index, 1)
}

// Watch for question type change
watch(() => questionForm.value.question_type, (newType) => {
  if (newType === 'true_false') {
    questionForm.value.options = [
      { option_text: 'Benar', is_correct: true },
      { option_text: 'Salah', is_correct: false }
    ]
  } else if (newType === 'short_answer') {
    questionForm.value.options = [{ option_text: '', is_correct: true }]
  } else if (questionForm.value.options.length < 2) {
    questionForm.value.options = [
      { option_text: '', is_correct: true },
      { option_text: '', is_correct: false },
      { option_text: '', is_correct: false },
      { option_text: '', is_correct: false }
    ]
  }
})

onMounted(async () => {
  await fetchLesson(lessonId.value)
  const quizData = await getQuizByLesson(lessonId.value)
  if (quizData) {
    settingsForm.value = {
      title: quizData.title,
      description: quizData.description,
      time_limit: quizData.time_limit,
      passing_score: quizData.passing_score,
      max_attempts: quizData.max_attempts,
      shuffle_questions: quizData.shuffle_questions,
      shuffle_options: quizData.shuffle_options,
      show_correct_answers: quizData.show_correct_answers
    }
  }
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
.slide-up-enter-active, .slide-up-leave-active {
  transition: all 0.3s ease;
}
.slide-up-enter-from, .slide-up-leave-to {
  transform: translateY(20px);
  opacity: 0;
}
</style>
