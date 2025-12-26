import { ref, computed, reactive } from 'vue'

interface ChatMessage {
    id: string
    role: 'user' | 'assistant' | 'system'
    content: string
    sources?: Source[]
    tokensUsed?: number
    createdAt: string
}

interface Source {
    lessonId: string
    lessonTitle: string
    contentType: string
    preview: string
    relevance: number
}

interface ChatSession {
    id: string
    userId: string
    courseId: string
    title: string
    isActive: boolean
    messageCount: number
    totalTokensUsed: number
    createdAt: string
    updatedAt: string
}

interface QuotaInfo {
    used: number
    limit: number
    remaining: number
}

interface AIStatus {
    ai_enabled: boolean
    embedding_count: number
    processing_status: any[]
    provider: string
    model: string
}

export function useAITutor(courseId: string) {
    const config = useRuntimeConfig()
    const token = useCookie('token')

    // State
    const messages = ref<ChatMessage[]>([])
    const session = ref<ChatSession | null>(null)
    const quota = ref<QuotaInfo>({ used: 0, limit: 50, remaining: 50 })
    const aiStatus = ref<AIStatus | null>(null)
    const isLoading = ref(false)
    const isSending = ref(false)
    const error = ref<string | null>(null)
    const isOpen = ref(false)

    // Computed
    const isAIEnabled = computed(() => aiStatus.value?.ai_enabled ?? false)
    const hasEmbeddings = computed(() => (aiStatus.value?.embedding_count ?? 0) > 0)
    const remainingQuota = computed(() => quota.value.remaining)
    const canChat = computed(() => isAIEnabled.value && remainingQuota.value > 0)

    // API helpers
    const apiBase = computed(() => config.public.apiBase || 'http://localhost:8080')

    async function fetchWithAuth(url: string, options: RequestInit = {}) {
        const headers = new Headers(options.headers)
        if (token.value) {
            headers.set('Authorization', `Bearer ${token.value}`)
        }
        headers.set('Content-Type', 'application/json')

        const response = await fetch(`${apiBase.value}${url}`, {
            ...options,
            headers
        })

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}))
            throw new Error(errorData.error || `HTTP ${response.status}`)
        }

        return response.json()
    }

    // Actions
    async function loadSession() {
        try {
            isLoading.value = true
            error.value = null

            const data = await fetchWithAuth(`/api/courses/${courseId}/chat/session`)

            session.value = data.session
            messages.value = data.messages || []
            quota.value = data.quota

            if (messages.value.length === 0) {
                // Add welcome message
                messages.value.push({
                    id: 'welcome',
                    role: 'assistant',
                    content: getWelcomeMessage(),
                    createdAt: new Date().toISOString()
                })
            }
        } catch (err: any) {
            error.value = err.message
            console.error('Failed to load AI session:', err)
        } finally {
            isLoading.value = false
        }
    }

    async function loadAIStatus() {
        try {
            const data = await fetchWithAuth(`/api/courses/${courseId}/ai-status`)
            aiStatus.value = data
        } catch (err: any) {
            console.error('Failed to load AI status:', err)
        }
    }

    async function sendMessage(content: string) {
        if (!content.trim() || isSending.value) return

        try {
            isSending.value = true
            error.value = null

            // Add user message immediately
            const tempUserMsg: ChatMessage = {
                id: `temp-${Date.now()}`,
                role: 'user',
                content: content.trim(),
                createdAt: new Date().toISOString()
            }
            messages.value.push(tempUserMsg)

            // Send to API
            const data = await fetchWithAuth(`/api/courses/${courseId}/chat`, {
                method: 'POST',
                body: JSON.stringify({ message: content.trim() })
            })

            // Add assistant response
            const assistantMsg: ChatMessage = {
                id: data.session_id + '-' + Date.now(),
                role: 'assistant',
                content: data.message,
                sources: data.sources,
                tokensUsed: data.tokens_used,
                createdAt: new Date().toISOString()
            }
            messages.value.push(assistantMsg)

            // Update quota
            quota.value = data.quota

            // Update session if provided
            if (data.session_id) {
                session.value = { ...session.value!, id: data.session_id }
            }

        } catch (err: any) {
            error.value = err.message
            // Remove temp message on error
            messages.value = messages.value.filter(m => !m.id.startsWith('temp-'))
            throw err
        } finally {
            isSending.value = false
        }
    }

    async function clearSession() {
        try {
            await fetchWithAuth(`/api/courses/${courseId}/chat/session`, {
                method: 'DELETE'
            })

            messages.value = [{
                id: 'welcome',
                role: 'assistant',
                content: getWelcomeMessage(),
                createdAt: new Date().toISOString()
            }]
            session.value = null

            // Reload session to get new one
            await loadSession()
        } catch (err: any) {
            error.value = err.message
            throw err
        }
    }

    async function loadQuota() {
        try {
            const data = await fetchWithAuth(`/api/courses/${courseId}/chat/quota`)
            quota.value = data
        } catch (err: any) {
            console.error('Failed to load quota:', err)
        }
    }

    function toggleChat() {
        isOpen.value = !isOpen.value
        if (isOpen.value && messages.value.length <= 1) {
            loadSession()
        }
    }

    function openChat() {
        isOpen.value = true
        if (messages.value.length <= 1) {
            loadSession()
        }
    }

    function closeChat() {
        isOpen.value = false
    }

    function getWelcomeMessage(): string {
        return `Halo! 👋 Saya AI Tutor untuk kursus ini.

Saya sudah mempelajari semua materi di kursus ini dan siap membantu Anda memahami kontennya. Silakan tanyakan apa saja tentang materi kursus!

**Beberapa contoh pertanyaan:**
- Apa itu [konsep dari materi]?
- Jelaskan perbedaan antara [A] dan [B]
- Bagaimana cara [melakukan sesuatu]?
- Berikan contoh [topik tertentu]

Saya akan menjawab berdasarkan materi yang tersedia di kursus ini. 🎓`
    }

    // Initialize
    loadAIStatus()

    return {
        // State
        messages,
        session,
        quota,
        aiStatus,
        isLoading,
        isSending,
        error,
        isOpen,

        // Computed
        isAIEnabled,
        hasEmbeddings,
        remainingQuota,
        canChat,

        // Actions
        loadSession,
        loadAIStatus,
        sendMessage,
        clearSession,
        loadQuota,
        toggleChat,
        openChat,
        closeChat
    }
}
