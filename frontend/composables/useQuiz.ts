import { ref } from 'vue'

export interface Quiz {
    id: string
    lesson_id: string
    title: string
    description: string
    time_limit: number
    passing_score: number
    max_attempts: number
    shuffle_questions: boolean
    shuffle_options: boolean
    show_correct_answers: boolean
    questions?: Question[]
    created_at: string
    updated_at: string
}

export interface Question {
    id: string
    quiz_id: string
    question_type: 'multiple_choice' | 'multiple_answer' | 'true_false' | 'short_answer'
    question_text: string
    explanation?: string
    points: number
    order_index: number
    required: boolean
    options?: Option[]
    created_at: string
}

export interface Option {
    id?: string
    question_id?: string
    option_text: string
    is_correct: boolean
    order_index?: number
}

export interface CreateQuizRequest {
    title: string
    description?: string
    time_limit?: number
    passing_score?: number
    max_attempts?: number
    shuffle_questions?: boolean
    shuffle_options?: boolean
    show_correct_answers?: boolean
}

export interface UpdateQuizRequest {
    title?: string
    description?: string
    time_limit?: number
    passing_score?: number
    max_attempts?: number
    shuffle_questions?: boolean
    shuffle_options?: boolean
    show_correct_answers?: boolean
}

export interface CreateQuestionRequest {
    question_type: string
    question_text: string
    explanation?: string
    points?: number
    required?: boolean
    options?: { option_text: string; is_correct: boolean }[]
}

export interface UpdateQuestionRequest {
    question_type?: string
    question_text?: string
    explanation?: string
    points?: number
    required?: boolean
    options?: { option_text: string; is_correct: boolean }[]
}

export interface QuizAttempt {
    id: string
    quiz_id: string
    user_id: string
    started_at: string
    completed_at?: string
    score?: number
    passed?: boolean
    time_spent?: number
}

export interface QuizResult {
    attempt_id: string
    quiz_id: string
    score: number
    passed: boolean
    time_spent: number
    total_points: number
    earned_points: number
    correct_count: number
    total_questions: number
    answers?: AnswerResult[]
}

export interface AnswerResult {
    question_id: string
    question_text: string
    user_answer: string[]
    correct_answer?: string[]
    is_correct: boolean
    points_earned: number
    max_points: number
    explanation?: string
}

export interface QuizStatus {
    quiz_id: string
    title: string
    description: string
    time_limit: number
    passing_score: number
    max_attempts: number
    attempt_count: number
    remaining_attempts: number // -1 means unlimited
    best_score?: number
    has_passed: boolean
    can_attempt: boolean
    in_progress_attempt?: QuizAttempt
    attempts: QuizAttempt[]
}

export function useQuiz() {
    const api = useApi()

    const quiz = ref<Quiz | null>(null)
    const questions = ref<Question[]>([])
    const loading = ref(false)
    const error = ref<string | null>(null)

    // Admin: Get quiz by lesson ID
    const getQuizByLesson = async (lessonId: string): Promise<Quiz | null> => {
        loading.value = true
        error.value = null
        try {
            const data = await api.fetch<Quiz>(`/api/admin/lessons/${lessonId}/quiz`)
            quiz.value = data
            questions.value = data.questions || []
            return data
        } catch (err: any) {
            if (err.status === 404 || err.statusCode === 404) {
                quiz.value = null
                return null
            }
            error.value = err.message || 'Failed to fetch quiz'
            return null
        } finally {
            loading.value = false
        }
    }

    // Admin: Create quiz
    const createQuiz = async (lessonId: string, data: CreateQuizRequest): Promise<Quiz | null> => {
        loading.value = true
        error.value = null
        try {
            const result = await api.fetch<Quiz>(`/api/admin/lessons/${lessonId}/quiz`, {
                method: 'POST',
                body: JSON.stringify(data)
            })
            quiz.value = result
            return result
        } catch (err: any) {
            error.value = err.message || 'Failed to create quiz'
            return null
        } finally {
            loading.value = false
        }
    }

    // Admin: Update quiz
    const updateQuiz = async (quizId: string, data: UpdateQuizRequest): Promise<Quiz | null> => {
        loading.value = true
        error.value = null
        try {
            const result = await api.fetch<Quiz>(`/api/admin/quizzes/${quizId}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            })
            quiz.value = result
            return result
        } catch (err: any) {
            error.value = err.message || 'Failed to update quiz'
            return null
        } finally {
            loading.value = false
        }
    }

    // Admin: Delete quiz
    const deleteQuiz = async (quizId: string): Promise<boolean> => {
        loading.value = true
        error.value = null
        try {
            await api.fetch(`/api/admin/quizzes/${quizId}`, { method: 'DELETE' })
            quiz.value = null
            questions.value = []
            return true
        } catch (err: any) {
            error.value = err.message || 'Failed to delete quiz'
            return false
        } finally {
            loading.value = false
        }
    }

    // Admin: Add question
    const addQuestion = async (quizId: string, data: CreateQuestionRequest): Promise<Question | null> => {
        loading.value = true
        error.value = null
        try {
            const result = await api.fetch<Question>(`/api/admin/quizzes/${quizId}/questions`, {
                method: 'POST',
                body: JSON.stringify(data)
            })
            questions.value.push(result)
            return result
        } catch (err: any) {
            error.value = err.message || 'Failed to add question'
            return null
        } finally {
            loading.value = false
        }
    }

    // Admin: Update question
    const updateQuestion = async (questionId: string, data: UpdateQuestionRequest): Promise<Question | null> => {
        loading.value = true
        error.value = null
        try {
            const result = await api.fetch<Question>(`/api/admin/questions/${questionId}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            })
            const idx = questions.value.findIndex(q => q.id === questionId)
            if (idx !== -1) questions.value[idx] = result
            return result
        } catch (err: any) {
            error.value = err.message || 'Failed to update question'
            return null
        } finally {
            loading.value = false
        }
    }

    // Admin: Delete question
    const deleteQuestion = async (questionId: string): Promise<boolean> => {
        loading.value = true
        error.value = null
        try {
            await api.fetch(`/api/admin/questions/${questionId}`, { method: 'DELETE' })
            questions.value = questions.value.filter(q => q.id !== questionId)
            return true
        } catch (err: any) {
            error.value = err.message || 'Failed to delete question'
            return false
        } finally {
            loading.value = false
        }
    }

    // Admin: Reorder questions
    const reorderQuestions = async (quizId: string, questionIds: string[]): Promise<boolean> => {
        loading.value = true
        error.value = null
        try {
            await api.fetch(`/api/admin/quizzes/${quizId}/questions/reorder`, {
                method: 'PUT',
                body: JSON.stringify({ question_ids: questionIds })
            })
            return true
        } catch (err: any) {
            error.value = err.message || 'Failed to reorder questions'
            return false
        } finally {
            loading.value = false
        }
    }

    // Student: Get quiz for taking
    const getQuizForStudent = async (lessonId: string): Promise<Quiz | null> => {
        loading.value = true
        error.value = null
        try {
            const data = await api.fetch<Quiz>(`/api/lessons/${lessonId}/quiz`)
            quiz.value = data
            questions.value = data.questions || []
            return data
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch quiz'
            return null
        } finally {
            loading.value = false
        }
    }

    // Student: Start attempt
    const startAttempt = async (quizId: string): Promise<QuizAttempt | null> => {
        loading.value = true
        error.value = null
        try {
            const result = await api.fetch<QuizAttempt>(`/api/quizzes/${quizId}/start`, {
                method: 'POST'
            })
            return result
        } catch (err: any) {
            error.value = err.message || 'Failed to start quiz'
            return null
        } finally {
            loading.value = false
        }
    }

    // Student: Submit answers
    const submitAttempt = async (attemptId: string, answers: { question_id: string; selected_option_ids?: string[]; text_answer?: string }[]): Promise<QuizResult | null> => {
        loading.value = true
        error.value = null
        try {
            const result = await api.fetch<QuizResult>(`/api/attempts/${attemptId}/submit`, {
                method: 'POST',
                body: JSON.stringify({ answers })
            })
            return result
        } catch (err: any) {
            error.value = err.message || 'Failed to submit quiz'
            return null
        } finally {
            loading.value = false
        }
    }

    // Student: Get attempt result
    const getAttemptResult = async (attemptId: string): Promise<QuizResult | null> => {
        loading.value = true
        error.value = null
        try {
            const result = await api.fetch<QuizResult>(`/api/attempts/${attemptId}/result`)
            return result
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch result'
            return null
        } finally {
            loading.value = false
        }
    }

    // Student: Get user's attempts
    const getUserAttempts = async (quizId: string): Promise<QuizAttempt[]> => {
        loading.value = true
        error.value = null
        try {
            const data = await api.fetch<{ attempts: QuizAttempt[] }>(`/api/quizzes/${quizId}/attempts`)
            return data.attempts || []
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch attempts'
            return []
        } finally {
            loading.value = false
        }
    }

    // Student: Get quiz status (attempt history, remaining attempts, etc.)
    const getQuizStatus = async (quizId: string): Promise<QuizStatus | null> => {
        loading.value = true
        error.value = null
        try {
            const result = await api.fetch<QuizStatus>(`/api/quizzes/${quizId}/status`)
            return result
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch quiz status'
            return null
        } finally {
            loading.value = false
        }
    }

    return {
        quiz,
        questions,
        loading,
        error,
        // Admin
        getQuizByLesson,
        createQuiz,
        updateQuiz,
        deleteQuiz,
        addQuestion,
        updateQuestion,
        deleteQuestion,
        reorderQuestions,
        // Student
        getQuizForStudent,
        startAttempt,
        submitAttempt,
        getAttemptResult,
        getUserAttempts,
        getQuizStatus
    }
}
