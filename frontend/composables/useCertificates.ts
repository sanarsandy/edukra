// Certificates composable for managing user certificates
import { ref } from 'vue'

export interface Certificate {
    id: string
    user_id: string
    course_id: string
    certificate_number: string
    issued_at: string
    pdf_url?: string
    course_name: string
    user_name: string
}

export const useCertificates = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)

    const certificates = ref<Certificate[]>([])
    const currentCertificate = ref<Certificate | null>(null)

    // Fetch all certificates for current user
    const fetchCertificates = async () => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<{ certificates: Certificate[] }>('/api/certificates')
            certificates.value = response.certificates || []
            return certificates.value
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch certificates'
            return []
        } finally {
            loading.value = false
        }
    }

    // Get a single certificate
    const getCertificate = async (id: string) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Certificate>(`/api/certificates/${id}`)
            currentCertificate.value = response
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch certificate'
            return null
        } finally {
            loading.value = false
        }
    }

    // Download certificate as PDF (opens in new tab for printing)
    const downloadPDF = async (id: string) => {
        const downloadUrl = `/api/certificates/${id}/download`
        // Open in new tab - browser will handle the HTML which can be printed to PDF
        window.open(downloadUrl, '_blank')
    }

    // View certificate in modal (fetch HTML)
    const viewCertificate = async (id: string) => {
        const viewUrl = `/api/certificates/${id}/download`
        return viewUrl
    }

    // Format date for display
    const formatDate = (dateString: string): string => {
        const date = new Date(dateString)
        return date.toLocaleDateString('id-ID', {
            day: 'numeric',
            month: 'long',
            year: 'numeric'
        })
    }

    return {
        loading,
        error,
        certificates,
        currentCertificate,
        fetchCertificates,
        getCertificate,
        downloadPDF,
        viewCertificate,
        formatDate,
    }
}
