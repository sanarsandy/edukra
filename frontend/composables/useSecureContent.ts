/**
 * Composable for handling secure content access via pre-signed URLs
 * Used for video and document content that requires enrollment verification
 */

interface SecureContentResponse {
    url: string
    expires_at: string
    type: string
}

export const useSecureContent = () => {
    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || 'http://localhost:8080'

    // Cache for pre-signed URLs to avoid unnecessary API calls
    const urlCache = ref<Map<string, { url: string, expiresAt: Date }>>(new Map())

    const loading = ref(false)
    const error = ref<string | null>(null)

    /**
     * Get a secure URL for video content
     * @param lessonId - The lesson ID to get content for
     * @returns Pre-signed URL or null if not available
     */
    const getSecureVideoUrl = async (lessonId: string): Promise<string | null> => {
        // Check cache first
        const cached = urlCache.value.get(`video-${lessonId}`)
        if (cached && cached.expiresAt > new Date()) {
            return cached.url
        }

        loading.value = true
        error.value = null

        try {
            const token = useCookie('token')
            const response = await fetch(`${apiBase}/api/content/${lessonId}/url`, {
                headers: {
                    'Authorization': `Bearer ${token.value}`
                }
            })

            if (!response.ok) {
                if (response.status === 403) {
                    error.value = 'Anda harus terdaftar di kursus ini untuk mengakses konten'
                    return null
                }
                if (response.status === 404) {
                    error.value = 'Konten tidak ditemukan'
                    return null
                }
                throw new Error('Failed to get secure URL')
            }

            const data: SecureContentResponse = await response.json()

            // Cache the URL (subtract 5 minutes from expiry as buffer)
            const expiresAt = new Date(data.expires_at)
            expiresAt.setMinutes(expiresAt.getMinutes() - 5)

            urlCache.value.set(`video-${lessonId}`, {
                url: data.url,
                expiresAt
            })

            return data.url
        } catch (err: any) {
            error.value = err.message || 'Failed to load content'
            console.error('Error getting secure URL:', err)
            return null
        } finally {
            loading.value = false
        }
    }

    /**
     * Get a secure URL for document content (PDF, etc.)
     * @param lessonId - The lesson ID to get content for
     * @returns Pre-signed URL or null if not available
     */
    const getSecureDocumentUrl = async (lessonId: string): Promise<string | null> => {
        // Check cache first
        const cached = urlCache.value.get(`doc-${lessonId}`)
        if (cached && cached.expiresAt > new Date()) {
            return cached.url
        }

        loading.value = true
        error.value = null

        try {
            const token = useCookie('token')
            const response = await fetch(`${apiBase}/api/content/${lessonId}/document`, {
                headers: {
                    'Authorization': `Bearer ${token.value}`
                }
            })

            if (!response.ok) {
                if (response.status === 403) {
                    error.value = 'Anda harus terdaftar di kursus ini untuk mengakses konten'
                    return null
                }
                throw new Error('Failed to get secure URL')
            }

            const data: SecureContentResponse = await response.json()

            // Cache the URL
            const expiresAt = new Date(data.expires_at)
            expiresAt.setMinutes(expiresAt.getMinutes() - 5)

            urlCache.value.set(`doc-${lessonId}`, {
                url: data.url,
                expiresAt
            })

            return data.url
        } catch (err: any) {
            error.value = err.message || 'Failed to load content'
            console.error('Error getting secure URL:', err)
            return null
        } finally {
            loading.value = false
        }
    }

    /**
     * Clear cached URL for a specific lesson
     */
    const clearCache = (lessonId?: string) => {
        if (lessonId) {
            urlCache.value.delete(`video-${lessonId}`)
            urlCache.value.delete(`doc-${lessonId}`)
        } else {
            urlCache.value.clear()
        }
    }

    /**
     * Check if a URL is a MinIO object (not a legacy /uploads path or embed URL)
     */
    const isMinioObject = (url: string | null | undefined): boolean => {
        if (!url) return false
        // MinIO objects don't start with /uploads or http
        // They are stored as just the object key like "1234567890_video.mp4"
        return !url.startsWith('/uploads') &&
            !url.startsWith('http://') &&
            !url.startsWith('https://') &&
            !url.includes('youtube.com') &&
            !url.includes('youtu.be') &&
            !url.includes('vimeo.com') &&
            !url.includes('drive.google.com')
    }

    /**
     * Check if URL is a legacy /uploads path
     */
    const isLegacyUpload = (url: string | null | undefined): boolean => {
        if (!url) return false
        return url.startsWith('/uploads')
    }

    /**
     * Get a stream URL for content (bypasses pre-signed URL issues)
     * This returns a URL that streams content directly through the backend
     * @param lessonId - The lesson ID to get content for
     * @returns Stream URL that includes auth token
     */
    const getStreamUrl = (lessonId: string): string => {
        const token = useCookie('token')
        // Return the stream URL - the token will be added via fetch headers
        // For direct src usage in iframe/video, we need a different approach
        return `${apiBase}/api/content/${lessonId}/stream`
    }

    /**
     * Fetch content as blob URL (for cases where we need direct URL in src)
     * This fetches the content and creates a local blob URL
     */
    const getContentBlobUrl = async (lessonId: string): Promise<string | null> => {
        loading.value = true
        error.value = null

        try {
            const token = useCookie('token')
            const response = await fetch(`${apiBase}/api/content/${lessonId}/stream`, {
                headers: {
                    'Authorization': `Bearer ${token.value}`
                }
            })

            if (!response.ok) {
                if (response.status === 403) {
                    error.value = 'Anda harus terdaftar di kursus ini untuk mengakses konten'
                    return null
                }
                if (response.status === 404) {
                    error.value = 'Konten tidak ditemukan'
                    return null
                }
                throw new Error('Failed to get content')
            }

            const blob = await response.blob()
            const blobUrl = URL.createObjectURL(blob)

            // Cache the blob URL
            urlCache.value.set(`blob-${lessonId}`, {
                url: blobUrl,
                expiresAt: new Date(Date.now() + 3600000) // 1 hour
            })

            return blobUrl
        } catch (err: any) {
            error.value = err.message || 'Failed to load content'
            console.error('Error getting content:', err)
            return null
        } finally {
            loading.value = false
        }
    }

    /**
     * Get thumbnail URL for course thumbnails
     * This uses the public /api/images/:objectKey endpoint for MinIO objects
     * @param thumbnailUrl - The thumbnail URL from database
     * @returns Proper URL for displaying the thumbnail
     */
    const getThumbnailUrl = (thumbnailUrl: string | null | undefined): string => {
        if (!thumbnailUrl) {
            return '/placeholder-course.jpg' // Fallback placeholder
        }

        // If it's a legacy /uploads path, return as-is (handled by nginx or static serve)
        if (thumbnailUrl.startsWith('/uploads')) {
            return thumbnailUrl
        }

        // If it's already a full URL (http/https), return as-is
        if (thumbnailUrl.startsWith('http://') || thumbnailUrl.startsWith('https://')) {
            return thumbnailUrl
        }

        // It's a MinIO object key - use the public images endpoint
        return `${apiBase}/api/images/${thumbnailUrl}`
    }

    return {
        loading: readonly(loading),
        error: readonly(error),
        getSecureVideoUrl,
        getSecureDocumentUrl,
        getStreamUrl,
        getContentBlobUrl,
        getThumbnailUrl,
        clearCache,
        isMinioObject,
        isLegacyUpload
    }
}
