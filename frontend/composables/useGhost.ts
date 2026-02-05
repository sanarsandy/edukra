// Ghost API response types
interface GhostPost {
    id: string
    slug: string
    title: string
    html: string
    excerpt?: string
    custom_excerpt?: string
    feature_image?: string
    published_at: string
    created_at: string
    url: string
    primary_author?: {
        name: string
        profile_image?: string
    }
    tags?: Array<{ name: string; slug: string }>
    meta_title?: string
    meta_description?: string
    og_title?: string
    og_description?: string
    og_image?: string
    twitter_title?: string
    twitter_description?: string
    twitter_image?: string
}

interface GhostPagination {
    page: number
    limit: number
    pages: number
    total: number
    next: number | null
    prev: number | null
}

interface GhostPostsResponse {
    posts: GhostPost[]
    meta: {
        pagination: GhostPagination
    }
}

export const useGhost = () => {
    const config = useRuntimeConfig()

    const fetchPosts = (limit = 10, page = 1) => {
        const ghostUrl = config.public.ghostUrl as string
        const ghostKey = config.public.ghostKey as string

        if (!ghostUrl || !ghostKey) {
            console.warn('Ghost URL or Key is missing')
            return { data: ref<GhostPostsResponse | null>(null), pending: ref(false), error: ref('Configuration missing') }
        }

        // Ensure URL doesn't end with slash
        const baseUrl = ghostUrl.replace(/\/$/, '')
        const url = `${baseUrl}/ghost/api/content/posts/?key=${ghostKey}&limit=${limit}&page=${page}&include=tags,authors`

        return useFetch<GhostPostsResponse>(url, {
            key: `ghost-posts-${page}-${limit}`
        })
    }

    const fetchPostBySlug = (slug: string) => {
        const ghostUrl = config.public.ghostUrl as string
        const ghostKey = config.public.ghostKey as string

        // Ensure URL doesn't end with slash
        const baseUrl = ghostUrl.replace(/\/$/, '')
        const url = `${baseUrl}/ghost/api/content/posts/slug/${slug}/?key=${ghostKey}&include=tags,authors`

        return useFetch<GhostPostsResponse>(url, {
            transform: (response) => {
                return response?.posts?.[0] || null
            },
            key: `ghost-post-${slug}`
        })
    }

    return {
        fetchPosts,
        fetchPostBySlug
    }
}
