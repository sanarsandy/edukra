import type { Block } from '~/utils/templates'
import type { ButtonStyle, ButtonShape, CardStyle, BackgroundType } from '~/utils/themes'

// =============================================
// TYPE DEFINITIONS
// =============================================

export interface CampaignStyles {
    primaryColor: string
    accentColor: string
    backgroundColor: string
    buttonColor: string
    buttonTextColor?: string
    textPrimaryColor: string
    textSecondaryColor: string
    fontFamily?: string
    fontFamilyHeading?: string
    fontFamilyBody?: string
    borderRadius?: string
    hasGradient?: boolean
    buttonStyle?: ButtonStyle
    buttonShape?: ButtonShape
    cardStyle?: CardStyle
    backgroundType?: BackgroundType
    backgroundImage?: string
    backgroundGradient?: string
    layoutMode?: 'wide' | 'mobile'
}

export interface CampaignData {
    id: number
    title: string
    slug: string
    course_id?: number
    course?: {
        id: number
        title: string
        price: number
        thumbnail_url?: string
    }
    webinar_id?: number
    webinar?: {
        id: number
        title: string
        price: number
        cover_image_url?: string
    }
    styles: CampaignStyles | string
    blocks: Block[] | string
    discount_price?: number
    discount_end_date?: string
    globalCustomCSS?: string
    whatsapp_number?: string
    meta_description?: string
    og_image_url?: string
    gtm_id?: string
    facebook_pixel_id?: string
}

// =============================================
// DEFAULT STYLES
// =============================================

export const defaultCampaignStyles: CampaignStyles = {
    primaryColor: '#6366f1',
    accentColor: '#f59e0b',
    backgroundColor: '#ffffff',
    buttonColor: '#6366f1',
    buttonTextColor: '#ffffff',
    textPrimaryColor: '#111827',
    textSecondaryColor: '#4b5563',
    fontFamilyHeading: 'Inter',
    fontFamilyBody: 'Inter',
    buttonStyle: 'solid',
    buttonShape: 'rounded',
    cardStyle: 'flat',
    backgroundType: 'solid',
    backgroundImage: '',
    layoutMode: 'wide'
}

// =============================================
// COMPOSABLE
// =============================================

export function useCampaignStyles(campaign: Ref<CampaignData | null>) {

    // Parse campaign styles
    const campaignStyles = computed((): CampaignStyles => {
        let styles = campaign.value?.styles
        if (typeof styles === 'string') {
            try { styles = JSON.parse(styles) } catch { styles = {} }
        }
        return { ...defaultCampaignStyles, ...(styles || {}) }
    })

    // Check if mobile layout mode
    const isMobileLayout = computed(() => campaignStyles.value.layoutMode === 'mobile')

    // Get container class based on layout mode
    const getContainerClass = (maxWidth = 'max-w-7xl') => {
        return isMobileLayout.value
            ? 'w-full max-w-[480px] mx-auto px-4'
            : `w-full ${maxWidth} mx-auto px-4`
    }

    // Parse and filter enabled blocks
    const enabledBlocks = computed(() => {
        let blocks = campaign.value?.blocks
        if (typeof blocks === 'string') {
            try { blocks = JSON.parse(blocks) } catch { blocks = [] }
        }
        return (blocks || [])
            .filter((b: Block) => b.enabled)
            .sort((a: Block, b: Block) => (a.order || 0) - (b.order || 0))
    })

    // Adjust hex color brightness
    const adjustColor = (hex: string, percent: number): string => {
        if (!hex || !hex.startsWith('#')) return hex
        const num = parseInt(hex.replace('#', ''), 16)
        const r = Math.min(255, Math.max(0, (num >> 16) + percent))
        const g = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + percent))
        const b = Math.min(255, Math.max(0, (num & 0x0000FF) + percent))
        return `#${(1 << 24 | r << 16 | g << 8 | b).toString(16).slice(1)}`
    }

    // Get per-block styles (minHeight, backgroundColor, accentColor)
    const getBlockStyles = (block: Block) => {
        const data = block.data || {}
        return {
            minHeight: data.minHeight ? `${data.minHeight}px` : undefined,
            backgroundColor: data.backgroundColor || undefined,
            accentColor: data.accentColor || campaignStyles.value.accentColor,
            sectionStyle: {
                ...(data.minHeight && { minHeight: `${data.minHeight}px` }),
                ...(data.backgroundColor && { backgroundColor: data.backgroundColor })
            }
        }
    }

    // Scope CSS to specific block for security
    const getScopedCSS = (block: Block): string => {
        const css = block.data.customCSS || ''
        const blockId = `#block-${block.id}`

        return css.split('}').map((rule: string) => {
            if (!rule.trim()) return ''
            const [selector, ...rest] = rule.split('{')
            if (!selector || !rest.length) return rule + '}'

            const scopedSelector = selector.split(',').map((s: string) => {
                const trimmed = s.trim()
                if (!trimmed || trimmed.startsWith('@')) return trimmed
                return `${blockId} ${trimmed}`
            }).join(', ')

            return `${scopedSelector} {${rest.join('{')}}`
        }).join('\n')
    }

    // Page style with CSS variables
    const pageStyle = computed(() => ({
        fontFamily: (campaignStyles.value.fontFamilyBody || campaignStyles.value.fontFamily || 'Inter') + ', sans-serif',
        backgroundColor: campaignStyles.value.backgroundColor,
        backgroundImage: campaignStyles.value.hasGradient
            ? `linear-gradient(135deg, ${campaignStyles.value.backgroundColor} 0%, ${adjustColor(campaignStyles.value.backgroundColor, -20)} 100%)`
            : 'none',
        '--primary-color': campaignStyles.value.primaryColor,
        '--accent-color': campaignStyles.value.accentColor,
        '--text-primary': campaignStyles.value.textPrimaryColor || '#111827',
        '--text-secondary': campaignStyles.value.textSecondaryColor || '#4b5563',
        '--btn-radius': campaignStyles.value.borderRadius === 'pill' ? '999px' : campaignStyles.value.borderRadius === 'sharp' ? '0px' : '0.75rem',
        '--font-heading': (campaignStyles.value.fontFamilyHeading || 'Inter') + ', sans-serif',
        '--font-body': (campaignStyles.value.fontFamilyBody || 'Inter') + ', sans-serif',
        color: campaignStyles.value.textPrimaryColor || '#111827'
    }))

    // Button style computed
    const buttonStyle = computed(() => {
        const s = campaignStyles.value || defaultCampaignStyles

        const shape = s.buttonShape || 'rounded'
        const style = s.buttonStyle || 'solid'
        const bgColor = s.buttonColor || '#6366f1'
        const primary = s.primaryColor || '#6366f1'
        const accent = s.accentColor || '#f59e0b'

        const base = {
            fontFamily: (s.fontFamilyHeading || 'Inter') + ', sans-serif',
            borderRadius: shape === 'pill' ? '999px' : shape === 'sharp' ? '0px' : shape === 'leaf' ? '1.5rem 0.25rem 1.5rem 0.25rem' : '0.75rem',
            transition: 'all 0.2s ease-in-out',
            fontWeight: '600'
        }

        const computedTextColor = s.buttonTextColor || (style === 'outline' || style === 'glass' ? bgColor : '#ffffff')

        switch (style) {
            case 'outline':
                return {
                    ...base,
                    backgroundColor: 'transparent',
                    border: `2px solid ${bgColor}`,
                    color: bgColor
                }
            case 'gradient':
                return {
                    ...base,
                    backgroundImage: `linear-gradient(to right, ${primary}, ${accent})`,
                    color: '#ffffff',
                    border: 'none',
                    boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)'
                }
            case 'glass':
                return {
                    ...base,
                    backgroundColor: 'rgba(255, 255, 255, 0.15)',
                    backdropFilter: 'blur(10px)',
                    border: `1px solid rgba(255, 255, 255, 0.3)`,
                    color: s.buttonTextColor || '#ffffff',
                    boxShadow: '0 8px 32px 0 rgba(31, 38, 135, 0.15)',
                    textShadow: '0 1px 2px rgba(0,0,0,0.1)'
                }
            case 'neo-brutalism':
                return {
                    ...base,
                    backgroundColor: bgColor,
                    color: computedTextColor,
                    border: '3px solid #000000',
                    boxShadow: '5px 5px 0px #000000',
                    fontWeight: '800',
                    textTransform: 'uppercase' as const
                }
            case 'soft-shadow':
                return {
                    ...base,
                    backgroundColor: bgColor,
                    color: computedTextColor,
                    boxShadow: `0 10px 25px -5px ${bgColor}66, 0 8px 10px -6px ${bgColor}33`,
                    border: 'none'
                }
            case 'solid':
            default:
                return {
                    ...base,
                    backgroundColor: bgColor,
                    color: computedTextColor,
                    border: 'none',
                    boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)'
                }
        }
    })

    // Hero style for background images
    const heroStyle = (block: Block) => {
        const bgImage = block.data.background_image || campaign.value?.course?.thumbnail_url || campaign.value?.webinar?.cover_image_url
        return {
            backgroundImage: bgImage ? `url(${bgImage})` : `linear-gradient(135deg, ${campaignStyles.value.backgroundColor}, ${adjustColor(campaignStyles.value.backgroundColor, -30)})`,
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            minHeight: block.data.minHeight ? `${block.data.minHeight}px` : undefined
        }
    }

    // Format price to Indonesian locale
    const formatPrice = (price: number): string => {
        return new Intl.NumberFormat('id-ID').format(price || 0)
    }

    // Get YouTube embed URL from various URL formats
    const getYoutubeEmbedUrl = (url: string): string | null => {
        if (!url) return null

        // Handle various YouTube URL formats
        const patterns = [
            /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([a-zA-Z0-9_-]+)/,
            /youtube\.com\/shorts\/([a-zA-Z0-9_-]+)/
        ]

        for (const pattern of patterns) {
            const match = url.match(pattern)
            if (match && match[1]) {
                return `https://www.youtube.com/embed/${match[1]}`
            }
        }

        return null
    }

    return {
        campaignStyles,
        isMobileLayout,
        getContainerClass,
        enabledBlocks,
        adjustColor,
        getBlockStyles,
        getScopedCSS,
        pageStyle,
        buttonStyle,
        heroStyle,
        formatPrice,
        getYoutubeEmbedUrl
    }
}

// =============================================
// COUNTDOWN COMPOSABLE
// =============================================

export function useCampaignCountdown(discountEndDate: Ref<string | undefined>) {
    const countdown = ref<{ days: number; hours: number; minutes: number; seconds: number } | null>(null)
    let countdownInterval: ReturnType<typeof setInterval> | null = null

    const startCountdown = () => {
        if (!discountEndDate.value) return

        const endDate = new Date(discountEndDate.value).getTime()

        const updateCountdown = () => {
            const now = new Date().getTime()
            const distance = endDate - now

            if (distance <= 0) {
                countdown.value = null
                if (countdownInterval) clearInterval(countdownInterval)
                return
            }

            countdown.value = {
                days: Math.floor(distance / (1000 * 60 * 60 * 24)),
                hours: Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
                minutes: Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60)),
                seconds: Math.floor((distance % (1000 * 60)) / 1000)
            }
        }

        updateCountdown()
        countdownInterval = setInterval(updateCountdown, 1000)
    }

    onMounted(() => {
        startCountdown()
    })

    onUnmounted(() => {
        if (countdownInterval) clearInterval(countdownInterval)
    })

    return { countdown }
}

// =============================================
// SCROLL ANIMATION COMPOSABLE
// =============================================

export function useScrollAnimation() {
    const observedElements = new Set<Element>()
    let observer: IntersectionObserver | null = null

    const initObserver = () => {
        if (typeof window === 'undefined') return

        observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.classList.add('animate')
                    observer?.unobserve(entry.target)
                    observedElements.delete(entry.target)
                }
            })
        }, {
            threshold: 0.1,
            rootMargin: '0px 0px -50px 0px'
        })
    }

    const observeElement = (el: Element | null) => {
        if (!el || !observer) return
        observer.observe(el)
        observedElements.add(el)
    }

    onMounted(() => {
        initObserver()
    })

    onUnmounted(() => {
        observedElements.forEach(el => observer?.unobserve(el))
        observedElements.clear()
        observer?.disconnect()
    })

    return { observeElement }
}
