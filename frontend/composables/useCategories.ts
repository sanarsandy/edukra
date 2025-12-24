// Categories composable for admin category management
export interface Category {
    id: string
    tenant_id?: string
    name: string
    slug: string
    description: string
    course_count?: number
}

export interface CategoriesResponse {
    categories: Category[]
    total: number
}

export const useCategories = () => {
    const api = useApi()
    const loading = ref(false)
    const error = ref<string | null>(null)
    const categories = ref<Category[]>([])
    const total = ref(0)

    // Fetch all categories
    const fetchCategories = async () => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<CategoriesResponse>('/api/categories')
            categories.value = response.categories || []
            total.value = response.total || 0
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch categories'
        } finally {
            loading.value = false
        }
    }

    // Create category
    const createCategory = async (data: Partial<Category>) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Category>('/api/admin/categories', {
                method: 'POST',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to create category'
            return null
        } finally {
            loading.value = false
        }
    }

    // Update category
    const updateCategory = async (id: string, data: Partial<Category>) => {
        loading.value = true
        error.value = null

        try {
            const response = await api.fetch<Category>(`/api/admin/categories/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data),
            })
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update category'
            return null
        } finally {
            loading.value = false
        }
    }

    // Delete category
    const deleteCategory = async (id: string) => {
        loading.value = true
        error.value = null

        try {
            await api.fetch(`/api/admin/categories/${id}`, { method: 'DELETE' })
            return true
        } catch (err: any) {
            error.value = err.message || 'Failed to delete category'
            return false
        } finally {
            loading.value = false
        }
    }

    return {
        loading,
        error,
        categories,
        total,
        fetchCategories,
        createCategory,
        updateCategory,
        deleteCategory,
    }
}
