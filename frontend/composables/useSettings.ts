// Settings composable for admin platform settings
export interface Settings {
    site_name: string
    site_description: string
    contact_email: string
    currency: string
    language: string
    logo_url?: string
    theme: string
    banner_enabled?: boolean
    banner_text?: string
    banner_link?: string
    banner_bg_color?: string
    banner_text_color?: string
}

export const useSettings = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)
    const settings = ref<Settings | null>(null)

    // Fetch settings
    const fetchSettings = async () => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Settings>('/api/settings')
            settings.value = response
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch settings'
        } finally {
            loading.value = false
        }
    }

    // Update settings
    const updateSettings = async (data: Partial<Settings>) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Settings>('/api/admin/settings', {
                method: 'PUT',
                body: JSON.stringify(data),
            })
            settings.value = response
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update settings'
            return null
        } finally {
            loading.value = false
        }
    }

    return {
        loading,
        error,
        settings,
        fetchSettings,
        updateSettings,
    }
}
