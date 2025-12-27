// Composable for payment/checkout functionality
export function usePayment() {
    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase

    // Get payment configuration (client key, enabled status)
    const fetchPaymentConfig = async () => {
        try {
            const token = useCookie('token')
            const data = await $fetch<{
                enabled: boolean
                provider: string
                client_key: string
                is_production: boolean
            }>(`${apiBase}/api/checkout/config`, {
                headers: token.value ? { 'Authorization': `Bearer ${token.value}` } : {}
            })
            return data
        } catch (err) {
            console.error('Failed to fetch payment config:', err)
            return null
        }
    }

    // Create checkout and get snap token
    const createCheckout = async (courseId: string, returnUrl?: string, paymentMethod?: string) => {
        const token = useCookie('token')

        if (!token.value) {
            throw new Error('Anda harus login terlebih dahulu')
        }

        const response = await $fetch<{
            transaction_id: string
            order_id: string
            snap_token?: string
            payment_url?: string
            client_key?: string
            expired_at?: string
            is_free: boolean
            message?: string
        }>(`${apiBase}/api/checkout`, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${token.value}` },
            body: {
                course_id: courseId,
                return_url: returnUrl || window.location.origin + '/dashboard/courses/' + courseId,
                payment_method: paymentMethod // For Duitku payment method selection
            }
        })

        return response
    }

    // Get user's transactions
    const fetchMyTransactions = async (limit = 20, offset = 0) => {
        const token = useCookie('token')

        if (!token.value) {
            throw new Error('Unauthorized')
        }

        const data = await $fetch<{
            transactions: any[]
            total: number
            limit: number
            offset: number
        }>(`${apiBase}/api/my/transactions?limit=${limit}&offset=${offset}`, {
            headers: { 'Authorization': `Bearer ${token.value}` }
        })

        return data
    }

    // Load Midtrans Snap.js script
    const loadSnapScript = (clientKey: string, isProduction: boolean): Promise<void> => {
        return new Promise((resolve, reject) => {
            // Check if already loaded
            if ((window as any).snap) {
                resolve()
                return
            }

            const script = document.createElement('script')
            script.src = isProduction
                ? 'https://app.midtrans.com/snap/snap.js'
                : 'https://app.sandbox.midtrans.com/snap/snap.js'
            script.setAttribute('data-client-key', clientKey)
            script.onload = () => resolve()
            script.onerror = () => reject(new Error('Failed to load Midtrans Snap'))
            document.head.appendChild(script)
        })
    }

    // Open Midtrans Snap payment popup
    const openSnapPayment = (
        snapToken: string,
        callbacks?: {
            onSuccess?: (result: any) => void
            onPending?: (result: any) => void
            onError?: (result: any) => void
            onClose?: () => void
        }
    ) => {
        const snap = (window as any).snap
        if (!snap) {
            throw new Error('Midtrans Snap is not loaded')
        }

        snap.pay(snapToken, {
            onSuccess: (result: any) => {
                console.log('[Snap] Payment success:', result)
                callbacks?.onSuccess?.(result)
            },
            onPending: (result: any) => {
                console.log('[Snap] Payment pending:', result)
                callbacks?.onPending?.(result)
            },
            onError: (result: any) => {
                console.error('[Snap] Payment error:', result)
                callbacks?.onError?.(result)
            },
            onClose: () => {
                console.log('[Snap] Popup closed')
                callbacks?.onClose?.()
            }
        })
    }

    // Format currency for display
    const formatPrice = (price: number, currency: string = 'IDR') => {
        if (price === 0) return 'Gratis'
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: currency,
            minimumFractionDigits: 0
        }).format(price)
    }

    // Get transaction status label
    const getStatusLabel = (status: string) => {
        const statusMap: Record<string, { label: string; color: string }> = {
            'pending': { label: 'Menunggu Pembayaran', color: 'text-amber-600 bg-amber-50' },
            'settlement': { label: 'Berhasil', color: 'text-green-600 bg-green-50' },
            'capture': { label: 'Berhasil', color: 'text-green-600 bg-green-50' },
            'deny': { label: 'Ditolak', color: 'text-red-600 bg-red-50' },
            'cancel': { label: 'Dibatalkan', color: 'text-red-600 bg-red-50' },
            'expire': { label: 'Kedaluwarsa', color: 'text-neutral-600 bg-neutral-100' },
            'failure': { label: 'Gagal', color: 'text-red-600 bg-red-50' },
            'refund': { label: 'Refund', color: 'text-blue-600 bg-blue-50' }
        }
        return statusMap[status] || { label: status, color: 'text-neutral-600 bg-neutral-100' }
    }

    return {
        fetchPaymentConfig,
        createCheckout,
        fetchMyTransactions,
        loadSnapScript,
        openSnapPayment,
        formatPrice,
        getStatusLabel
    }
}
