

export interface Block {
    id: string
    type: string
    variant?: string // 'default', 'split', 'centered', 'card', 'video_bg', etc.
    enabled: boolean
    order: number
    data: any
}

export interface TemplateStyle {
    primaryColor: string
    accentColor: string
    backgroundColor: string
    buttonColor: string
    textPrimaryColor: string
    textSecondaryColor: string
    buttonStyle: 'solid' | 'outline' | 'gradient'
    borderRadius: 'sharp' | 'rounded' | 'pill'
    hasGradient: boolean
    fontFamily: 'Inter' | 'Poppins' | 'Playfair Display' | 'Roboto'
    layoutMode?: 'wide' | 'mobile' // Optional: default to 'wide'
}

export interface CampaignTemplate {
    id: string
    name: string
    description: string
    category: 'bio' | 'webinar' | 'ecourse' | 'product' | 'service'
    thumbnail: string
    styles: TemplateStyle
    blocks: Block[]
}

// NEW: Pre-defined Design Themes (presets) to quickly switch look & feel
export const CAMPAIGN_THEMES: { id: string, name: string, styles: TemplateStyle }[] = [
    {
        id: 'modern_clean',
        name: 'Modern Clean',
        styles: {
            primaryColor: '#2563eb', // Blue
            accentColor: '#f59e0b',
            backgroundColor: '#ffffff',
            buttonColor: '#2563eb',
            textPrimaryColor: '#111827',
            textSecondaryColor: '#4b5563',
            buttonStyle: 'solid',
            borderRadius: 'rounded',
            hasGradient: false,
            fontFamily: 'Inter'
        }
    },
    {
        id: 'cyber_dark',
        name: 'Cyberpunk Dark',
        styles: {
            primaryColor: '#00f2ff', // Cyan
            accentColor: '#fe00fe', // Neon Pink
            backgroundColor: '#0f172a', // Dark Slate
            buttonColor: '#fe00fe',
            textPrimaryColor: '#f8fafc',
            textSecondaryColor: '#cbd5e1',
            buttonStyle: 'gradient',
            borderRadius: 'sharp',
            hasGradient: true,
            fontFamily: 'Roboto'
        }
    },
    {
        id: 'elegant_serif',
        name: 'Elegant Luxury',
        styles: {
            primaryColor: '#1c1917', // Stone 900
            accentColor: '#d4af37', // Gold
            backgroundColor: '#f5f5f4', // Stone 100
            buttonColor: '#1c1917',
            textPrimaryColor: '#1c1917',
            textSecondaryColor: '#57534e',
            buttonStyle: 'outline',
            borderRadius: 'pill',
            hasGradient: false,
            fontFamily: 'Playfair Display'
        }
    },
    {
        id: 'playful_pop',
        name: 'Playful Pop',
        styles: {
            primaryColor: '#ec4899', // Pink
            accentColor: '#facc15', // Yellow
            backgroundColor: '#fffbeb', // Amber 50
            buttonColor: '#ec4899',
            textPrimaryColor: '#111827',
            textSecondaryColor: '#4b5563',
            buttonStyle: 'solid',
            borderRadius: 'rounded',
            hasGradient: true,
            fontFamily: 'Poppins'
        }
    }
]

export const CAMPAIGN_TEMPLATES: CampaignTemplate[] = [
    // 1. Bio Link (Personal Brand) - Minimalist & Clean
    {
        id: 'bio_minimal',
        name: 'Simple Bio Link',
        description: 'Tampilan bersih untuk personal branding dan kumpulan link.',
        category: 'bio',
        thumbnail: '👤',
        styles: {
            primaryColor: '#18181b', // Zinc 900
            accentColor: '#3b82f6', // Blue 500
            backgroundColor: '#fafafa', // Zinc 50
            buttonColor: '#18181b',
            textPrimaryColor: '#18181b',
            textSecondaryColor: '#52525b',
            buttonStyle: 'outline',
            borderRadius: 'pill',
            hasGradient: false,
            fontFamily: 'Inter'
        },
        blocks: [
            {
                id: 'hero_bio',
                type: 'hero',
                variant: 'centered',
                enabled: true,
                order: 1,
                data: {
                    headline: 'Nama Anda / Brand',
                    subheadline: 'Membantu Anda belajar skill baru setiap hari. Cek link di bawah!',
                    cta_text: '',
                    cta_link: '',
                    background_image: null,
                    badge: 'Digital Creator'
                }
            },

            {
                id: 'bonus_links',
                type: 'bonus', // Repurposing bonus as link stack or we use multiple CTAs. 
                // actually standard blocks in create.vue: hero, countdown, benefits, pricing, testimonials, video, trust, faq, instructor, cta_banner, statistics, bonus, curriculum, gallery, achievement, comparison, floating_chat, social_proof
                // For Bio link, usually it's just buttons. Using 'cta_banner' multiple times or a dedicated 'links' block if exists. 
                // Assuming we don't have a 'links' block yet based on create.vue. 
                // We will use 'benefits' to list key offering links or just the hero + CTA.
                // Let's use 'hero' + 'benefits' (icons as links) + 'social_proof'.
                // Actually, let's stick to available blocks.
                enabled: true,
                order: 2,
                data: {
                    title: 'Produk Unggulan',
                    subtitle: 'Rekomendasi terbaik untuk Anda',
                    items: [
                        { emoji: '📘', title: 'Ebook Gratis', value: 0, description: 'Panduan lengkap untuk pemula' },
                        { emoji: '🚀', title: 'Kelas Online', value: 0, description: 'Belajar intensif selama 4 minggu' },
                        { emoji: '💬', title: 'Konsultasi 1-on-1', value: 0, description: 'Diskusi privat via Zoom' }
                    ],
                    total_value: 0
                }
            },
            {
                id: 'instructor_bio',
                type: 'instructor',
                enabled: true,
                order: 3,
                data: { auto_fill: true }
            }
        ]
    },

    // 2. Webinar Registration (Urgency) - High Contrast & Bold
    {
        id: 'webinar_bold',
        name: 'High Conversion Webinar',
        description: 'Fokus pada urgensi dan konversi pendaftaran webinar.',
        category: 'webinar',
        thumbnail: '🔥',
        styles: {
            primaryColor: '#dc2626', // Red 600
            accentColor: '#f59e0b', // Amber 500
            backgroundColor: '#111827', // Gray 900
            buttonColor: '#dc2626',
            textPrimaryColor: '#ffffff',
            textSecondaryColor: '#d1d5db',
            buttonStyle: 'solid',
            borderRadius: 'rounded',
            hasGradient: true,
            fontFamily: 'Poppins'
        },
        blocks: [
            {
                id: 'hero_web',
                type: 'hero',
                enabled: true,
                order: 1,
                data: {
                    headline: 'Cara Menguasai Skill High-Income di 2024',
                    subheadline: 'Webinar Eksklusif: Strategi rahasia yang belum pernah dibongkar sebelumnya.',
                    cta_text: 'DAFTAR GRATIS SEKARANG',
                    cta_link: '#register',
                    badge: 'LIVE WEBINAR'
                }
            },
            {
                id: 'countdown_web',
                type: 'countdown',
                enabled: true,
                order: 2,
                data: { label: 'Webinar Dimulai Dalam:' }
            },
            {
                id: 'video_web',
                type: 'video',
                enabled: true,
                order: 3,
                data: {
                    title: 'Tonton Cuplikan Materi',
                    subtitle: 'Apa saja yang akan kita bahas?',
                    youtube_url: '' // User to fill
                }
            },
            {
                id: 'benefits_web',
                type: 'benefits',
                enabled: true,
                order: 4,
                data: {
                    title: 'Kenapa Anda Wajib Hadir?',
                    subtitle: 'Ilmu daging yang akan Anda bawa pulang:',
                    items: [
                        { icon: 'lock', text: 'Studi kasus real yang tidak ada di Google' },
                        { icon: 'trending-up', text: 'Framework step-by-step yang terbukti' },
                        { icon: 'gift', text: 'Bonus template senilai Rp 1.5 Juta' }
                    ]
                }
            },
            {
                id: 'instructor_web',
                type: 'instructor',
                enabled: true,
                order: 5,
                data: { auto_fill: true }
            },
            {
                id: 'cta_web',
                type: 'cta_banner',
                enabled: true,
                order: 6,
                data: {
                    headline: 'Kuota Terbatas! Amankan Slot Anda.',
                    subheadline: 'Hanya tersisa untuk 50 pendaftar lagi.',
                    cta_text: 'Saya Mau Daftar!'
                }
            }
        ]
    },

    // 3. E-Course Sales (Professional) - Trustworthy & Detailed
    {
        id: 'ecourse_pro',
        name: 'Professional Course',
        description: 'Layout lengkap dan profesional untuk menjual kelas online.',
        category: 'ecourse',
        thumbnail: '🎓',
        styles: {
            primaryColor: '#2563eb', // Blue 600
            accentColor: '#10b981', // Emerald 500
            backgroundColor: '#ffffff',
            buttonColor: '#2563eb',
            textPrimaryColor: '#1e293b',
            textSecondaryColor: '#64748b',
            buttonStyle: 'solid',
            borderRadius: 'rounded',
            hasGradient: false,
            fontFamily: 'Inter'
        },
        blocks: [
            {
                id: 'hero_course',
                type: 'hero',
                enabled: true,
                order: 1,
                data: {
                    headline: 'Mastering Digital Marketing: Zero to Hero',
                    subheadline: 'Panduan terlengkap menjadi digital marketer handal dalam 8 minggu.',
                    cta_text: 'Mulai Belajar Sekarang',
                    cta_link: '#pricing',
                    badge: 'BEST SELLER'
                }
            },
            {
                id: 'statistics_course',
                type: 'statistics',
                enabled: true,
                order: 2,
                data: {
                    title: '',
                    items: [
                        { value: '2500', suffix: '+', label: 'Alumni Sukses' },
                        { value: '4.9', suffix: '/5', label: 'Rating' },
                        { value: '50', suffix: 'Jam', label: 'Materi Video' }
                    ]
                }
            },
            {
                id: 'curriculum_course',
                type: 'curriculum',
                enabled: true,
                order: 3,
                data: {
                    title: 'Kurikulum Terstruktur',
                    subtitle: 'Dari dasar hingga mahir, semua dipelajari.',
                    modules: [
                        { title: 'Modul 1: Fundamental', lessons_count: 5, lessons: ['Mindset Digital', 'Market Research'] },
                        { title: 'Modul 2: Content Strategy', lessons_count: 8, lessons: ['Copywriting', 'Visual Design'] }
                    ],
                    total_lessons: 45,
                    total_hours: 12
                }
            },
            {
                id: 'testimonials_course',
                type: 'testimonials',
                enabled: true,
                order: 4,
                data: { title: 'Kata Mereka yang Sudah Bergabung', items: [] }
            },
            {
                id: 'pricing_course',
                type: 'pricing',
                enabled: true,
                order: 5,
                data: {
                    original_price: 999000,
                    discount_price: 499000,
                    currency: 'IDR',
                    cta_text: 'Ambil Promo Ini',
                    show_timer: true
                }
            },
            {
                id: 'faq_course',
                type: 'faq',
                enabled: true,
                order: 6,
                data: { title: 'Pertanyaan yang Sering Diajukan', items: [] }
            }
        ]
    },

    // 4. Digital Product / Ebook (Minimalist) - Focus on Value
    {
        id: 'ebook_clean',
        name: 'Clean E-Book Sales',
        description: 'Desain minimalis fokus pada cover produk dan manfaat.',
        category: 'product',
        thumbnail: '📚',
        styles: {
            primaryColor: '#000000',
            accentColor: '#84cc16', // Lime 500
            backgroundColor: '#f3f4f6', // Gray 100
            buttonColor: '#000000',
            textPrimaryColor: '#1f2937',
            textSecondaryColor: '#4b5563',
            buttonStyle: 'solid',
            borderRadius: 'rounded',
            hasGradient: false,
            fontFamily: 'Poppins'
        },
        blocks: [
            {
                id: 'hero_book',
                type: 'hero',
                enabled: true,
                order: 1,
                data: {
                    headline: 'The Ultimate Guide to Productivity',
                    subheadline: 'Tingkatkan produktivitas Anda 10x lipat dengan metode sederhana.',
                    cta_text: 'Download E-Book (Instant Access)',
                    cta_link: '#pricing',
                    badge: 'NEW RELEASE'
                }
            },
            {
                id: 'benefits_book',
                type: 'benefits',
                enabled: true,
                order: 2,
                data: {
                    title: 'Apa yang akan Anda pelajari?',
                    subtitle: 'Bukan sekadar teori, tapi panduan praktis.',
                    items: [
                        { icon: 'check', text: 'Teknik Time Blocking' },
                        { icon: 'check', text: 'Format To-Do List yang Efektif' },
                        { icon: 'check', text: 'Mengelola Energi, Bukan Waktu' }
                    ]
                }
            },
            {
                id: 'bonus_book',
                type: 'bonus',
                enabled: true,
                order: 3,
                data: {
                    title: 'Bonus Spesial Hari Ini',
                    items: [
                        { emoji: '📋', title: 'Planner Template', value: 50000, description: 'Format PDF siap cetak' },
                        { emoji: '🎵', title: 'Focus Playlist', value: 25000, description: 'Musik untuk fokus kerja' }
                    ],
                    total_value: 75000
                }
            },
            {
                id: 'pricing_book',
                type: 'pricing',
                enabled: true,
                order: 4,
                data: {
                    original_price: 150000,
                    discount_price: 49000,
                    currency: 'IDR',
                    cta_text: 'Download Sekarang',
                    show_timer: false
                }
            }
        ]
    },

    // 5. Coaching / Consultation (Vibrant) - Personal & Energetic
    {
        id: 'coaching_vibrant',
        name: 'Vibrant Coaching',
        description: 'Desain energik untuk menawarkan jasa konsultasi atau coaching.',
        category: 'service',
        thumbnail: '🚀',
        styles: {
            primaryColor: '#7c3aed', // Violet 600
            accentColor: '#f472b6', // Pink 400
            backgroundColor: '#fdf4ff', // Fuchsia 50
            buttonColor: '#7c3aed',
            textPrimaryColor: '#4c1d95',
            textSecondaryColor: '#6d28d9',
            buttonStyle: 'outline',
            borderRadius: 'rounded',
            hasGradient: true,
            fontFamily: 'Playfair Display'
        },
        blocks: [
            {
                id: 'hero_coach',
                type: 'hero',
                enabled: true,
                order: 1,
                data: {
                    headline: 'Transformasi Bisnis Anda dalam 30 Hari',
                    subheadline: 'Program mentoring intensif untuk scaling bisnis Anda ke level berikutnya.',
                    cta_text: 'Book Free Call',
                    cta_link: '#chat',
                    badge: 'LIMITED SLOTS'
                }
            },
            {
                id: 'trust_coach',
                type: 'trust',
                enabled: true,
                order: 2,
                data: {}
            },
            {
                id: 'video_coach',
                type: 'video',
                enabled: true,
                order: 3,
                data: {
                    title: 'Dengar Cerita Mereka',
                    subtitle: 'Klien yang telah berhasil scaling bisnisnya.',
                    youtube_url: ''
                }
            },
            {
                id: 'comparison_coach',
                type: 'comparison',
                enabled: true,
                order: 4,
                data: {
                    title: 'Pilih Paket Coaching',
                    headers: ['Manfaat', 'Group', 'Private Vip'],
                    rows: [
                        { feature: 'Sesi Zoom Live', val_a: true, val_b: true },
                        { feature: 'Grup Support WA', val_a: true, val_b: true },
                        { feature: 'Review Bisnis Personal', val_a: false, val_b: true },
                        { feature: 'Akses Direct WA Mentor', val_a: false, val_b: true }
                    ]
                }
            },
            {
                id: 'floating_chat_coach',
                type: 'floating_chat',
                enabled: true,
                order: 5,
                data: {
                    phone: '',
                    message: 'Halo Coach, saya mau tanya program mentoringnya.',
                    label: 'Tanya Dulu Yuk'
                }
            }
        ]
    }
]
