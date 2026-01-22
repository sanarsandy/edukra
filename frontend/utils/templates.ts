

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
    isRigid?: boolean // If true, user cannot add/remove blocks or change detailed styles
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
    // 0. NEW: Gen Z Webinar (Custom Lovable Design)
    {
        id: 'webinar_gen_z',
        name: 'Webinar Gen Z (Lovable Style)',
        description: 'Style modern & playful dengan elemen 3D floating, gradient, dan glassmorphism.',
        category: 'webinar',
        thumbnail: '🚀',
        isRigid: true,
        styles: {
            primaryColor: '#D946EF', // Fuchsia 500
            accentColor: '#F97316', // Orange 500
            backgroundColor: '#050505', // Deep Dark
            buttonColor: '#D946EF',
            textPrimaryColor: '#FFFFFF',
            textSecondaryColor: '#A3A3A3',
            buttonStyle: 'gradient',
            borderRadius: 'rounded',
            hasGradient: true,
            fontFamily: 'Outfit' as any
        },
        blocks: [
            {
                id: 'hero_gz',
                type: 'hero_gen_z',
                enabled: true,
                order: 1,
                data: {
                    headline: 'BELAJAR KEPABEANAN & CUKAI',
                    subheadline: 'Kupas Tuntas Aturan, Prosedur, dan Strategi Ekspor-Impor Langsung dari Ahlinya.',
                    promo_text: 'Promo Berakhir Dalam:',
                    badge: 'WEBINAR EXCLUSIVE 2024'
                }
            },
            {
                id: 'features_gz',
                type: 'features_gen_z',
                enabled: true,
                order: 2,
                data: {
                    title: 'Yang Bakal Kamu DAPETIN'
                }
            },
            {
                id: 'cta_gz',
                type: 'cta_gen_z',
                enabled: true,
                order: 3,
                data: {
                    headline: 'AMANKAN KUOTA',
                    subheadline: 'SEBELUM HABIS!',
                    cta_text: 'DAFTAR SEKARANG'
                }
            }
        ]
    },
    // NEW: TikTok Shop Style (Clean Minimalist)
    {
        id: 'webinar_tiktok',
        name: 'Webinar TikTok Style (Clean)',
        description: 'Tampilan bersih & minimalis ala TikTok Shop. Fokus konversi tinggi.',
        category: 'webinar',
        thumbnail: '🎵',
        isRigid: true,
        styles: {
            primaryColor: '#F472B6', // Pink 400
            accentColor: '#38BDF8', // Light Blue
            backgroundColor: '#FFFFFF', // White
            buttonColor: '#F472B6',
            textPrimaryColor: '#000000',
            textSecondaryColor: '#6B7280',
            buttonStyle: 'gradient',
            borderRadius: 'rounded',
            hasGradient: true,
            fontFamily: 'Inter' as any
        },
        blocks: [
            {
                id: 'hero_clean',
                type: 'hero_clean',
                enabled: true,
                order: 1,
                data: {
                    headline: 'Rahasia Laris di TikTok Shop!',
                    subheadline: 'Pelajari strategi jualan yang terbukti meningkatkan penjualan 10x lipat di platform TikTok Shop',
                    badge: 'WEBINAR GRATIS'
                }
            },
            {
                id: 'speaker_clean',
                type: 'speaker_clean',
                enabled: true,
                order: 2,
                data: {
                    name: 'Budi Santoso',
                    title: 'TikTok Shop Specialist & Digital Marketing Expert',
                    bio: 'Budi telah membantu lebih dari 200+ UMKM meningkatkan penjualan mereka melalui platform TikTok Shop. Dengan pengalaman 5+ tahun di bidang e-commerce dan sosial media marketing.',
                    stats: [
                        { icon: '👥', label: '200+ UMKM' },
                        { icon: '🏆', label: '5+ Tahun Pengalaman' }
                    ]
                }
            },
            {
                id: 'features_clean',
                type: 'features_clean',
                enabled: true,
                order: 3,
                data: {
                    title: 'Apa Yang Akan Anda Dapatkan?',
                    items: [
                        { title: 'Strategi Algoritma TikTok', text: 'Pelajari cara memanfaatkan algoritma TikTok untuk meningkatkan visibilitas.', icon: '📈' },
                        { title: 'Optimasi Toko Online', text: 'Tips mengoptimalkan tampilan toko dan deskripsi produk.', icon: '🛍️' },
                        { title: 'Teknik Konten Viral', text: 'Formula membuat video pendek yang menarik dan berpotensi viral.', icon: '🎬' },
                        { title: 'Strategi Pricing', text: 'Cara menentukan harga yang kompetitif namun tetap menguntungkan.', icon: '💲' }
                    ]
                }
            },
            {
                id: 'cta_clean',
                type: 'cta_clean',
                enabled: true,
                order: 4,
                data: {
                    headline: 'Daftar Sekarang',
                    subheadline: 'Isi formulir di bawah ini untuk mendapatkan akses ke webinar eksklusif ini'
                    // Form fields are intrinsic to the component
                }
            },
            {
                id: 'faq_clean',
                type: 'faq', // Reuse generic FAQ/Clean variant if possible, or new component. Let's reuse FAQ with 'clean' variant style if possible, but user wants 'persis'. For now, reuse standard FAQ but maybe we need a clean variant.
                // Actually the screenshots show a specific FAQ style. Let's stick to standard FAQ but styled via theme for now to save time, unless it's very specific.
                // Screenshot FAQ looks standard clear text.
                enabled: true,
                order: 5,
                data: {
                    title: 'Pertanyaan Umum',
                    items: [
                        { question: 'Apakah webinar ini benar-benar gratis?', answer: 'Ya, webinar ini sepenuhnya gratis.' },
                        { question: 'Apakah saya akan mendapatkan rekaman?', answer: 'Ya, semua peserta akan mendapatkan akses rekaman.' },
                        { question: 'Apakah saya bisa mengajukan pertanyaan?', answer: 'Tentu saja! Akan ada sesi tanya jawab.' }
                    ]
                }
            }
        ]

    },
    // NEW: Gen Z Customs Webinar (User Request)
    {
        id: 'webinar_gen_z_customs',
        name: 'Webinar Bea Cukai (Gen Z)',
        description: 'Template modern & fun khusus anak muda dengan gaya visual Gen Z.',
        category: 'webinar',
        thumbnail: '🚀',
        isRigid: true,
        styles: {
            primaryColor: '#d946ef', // Fuchsia 500
            accentColor: '#f97316', // Orange 500
            backgroundColor: '#050505',
            buttonColor: 'linear-gradient(to right, #ec4899, #f97316)',
            textPrimaryColor: '#FFFFFF',
            textSecondaryColor: '#a3a3a3',
            buttonStyle: 'gradient',
            borderRadius: 'rounded',
            hasGradient: true,
            fontFamily: 'Outfit' as any
        },
        blocks: [
            {
                id: 'hero_gz',
                type: 'hero_gen_z',
                enabled: true,
                order: 1,
                data: {
                    headline: 'Bea & Cukai Itu Gampang!',
                    subheadline: 'Webinar Gratis: Ngomongin impor–ekspor, pajak barang, dan belanja dari luar negeri tanpa ribet.',
                    badge: 'WEBINAR GRATIS',
                    promo_text: 'JANGAN SAMPAI KETINGGALAN',
                    end_date: '2026-02-15T19:00:00+07:00'
                }
            },
            {
                id: 'features_gz',
                type: 'features_gen_z',
                enabled: true,
                order: 2,
                data: {
                    title: 'Apa Aja Yang Bakal',
                    items: [
                        { title: 'Bea Cukai 101', desc: 'Apa itu Bea dan Cukai (versi paling gampang) & kenapa barang impor kena pajak.', icon: '🧠', color: 'text-[#e879f9]', bg_gradient: 'from-purple-900/20 to-neutral-900', icon_gradient: 'from-[#e879f9] to-[#a855f7]' },
                        { title: 'Studi Kasus Nyata', desc: 'Contoh kasus belanja online, kirim barang, dan kesalahan umum biar kamu nggak boncos.', icon: '📦', color: 'text-[#f472b6]', bg_gradient: 'from-pink-900/20 to-neutral-900', icon_gradient: 'from-[#f472b6] to-[#ec4899]' },
                        { title: 'Anti Salah Langkah', desc: 'Tips praktis & trik biar nggak kaget saat kena bea masuk. Wajib tau!', icon: '💡', color: 'text-[#60a5fa]', bg_gradient: 'from-blue-900/20 to-neutral-900', icon_gradient: 'from-[#60a5fa] to-[#3b82f6]' },
                        { title: 'Tanya Sepuasnya', desc: 'Sesi Q&A langsung bareng narasumber ahli. Tanyain apa aja seputar custom!', icon: '💬', color: 'text-[#fbbf24]', bg_gradient: 'from-orange-900/20 to-neutral-900', icon_gradient: 'from-[#fbbf24] to-[#f59e0b]' }
                    ]
                }
            },
            {
                id: 'speaker_gz',
                type: 'speaker_gen_z',
                enabled: true,
                order: 3,
                data: {
                    name: 'Joni',
                    title: 'Ex Auditor Bea & Cukai',
                    bio: 'Bongkar fakta, mitos, dan kesalahan umum soal Bea & Cukai langsung dari orang dalam. Disampaikan dengan bahasa awam dan santai.',
                    image_url: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=300',
                    stats: [
                        { icon: '🕵️‍♂️', label: '10th+ Pengalaman' },
                        { icon: '📚', label: 'Praktisi Cukai' }
                    ]
                }
            },
            {
                id: 'faq_c',
                type: 'faq',
                enabled: true,
                order: 4,
                data: {
                    title: 'FAQ Singkat',
                    items: [
                        { question: 'Apakah webinar ini benar-benar gratis?', answer: 'Ya, 100% gratis tanpa biaya apa pun.' },
                        { question: 'Apakah pemula boleh ikut?', answer: 'Sangat boleh! Ini memang dibuat untuk pemula.' },
                        { question: 'Apakah dapat sertifikat?', answer: 'Ya, peserta akan mendapatkan e-sertifikat.' }
                    ]
                }
            },
            {
                id: 'cta_gz',
                type: 'cta_gen_z',
                enabled: true,
                order: 5,
                data: {
                    headline: 'Siap Ikutan? 🚀',
                    subheadline: 'Daftar sekarang, kuota terbatas! Amankan tempatmu sekarang.',
                    button_text: 'Daftar Webinar Gratis'
                }
            }
        ]
    },
    // Professional Minimalist Customs Webinar
    {
        id: 'webinar_pro_customs',
        name: 'Webinar Bea Cukai (Professional)',
        description: 'Template profesional & minimalis dengan warna navy-amber yang elegan.',
        category: 'webinar',
        thumbnail: '🛃',
        isRigid: true,
        styles: {
            primaryColor: '#1e3a8a', // Blue 900
            accentColor: '#f59e0b', // Amber 500
            backgroundColor: '#FFFFFF',
            buttonColor: '#1e3a8a',
            textPrimaryColor: '#1e293b',
            textSecondaryColor: '#64748b',
            buttonStyle: 'solid',
            borderRadius: 'rounded',
            hasGradient: false,
            fontFamily: 'Inter' as any
        },
        blocks: [
            {
                id: 'hero_pro',
                type: 'hero_pro',
                enabled: true,
                order: 1,
                data: {
                    headline: 'Webinar Gratis: Bea & Cukai Itu Gampang!',
                    subheadline: 'Pahami regulasi impor-ekspor dengan bahasa sederhana. Cocok untuk pemula yang ingin belajar kepabeanan tanpa pusing!',
                    badge: 'WEBINAR GRATIS',
                    button_text: 'Daftar Sekarang — Gratis',
                    end_date: '2026-01-28T19:05:00+07:00',
                    speaker_image: '/images/Foto Narasumber-Indra Pratama.jpeg',
                    speaker_name: 'Indra Pratama',
                    speaker_title: 'Ex Auditor Bea & Cukai'
                }
            },
            {
                id: 'features_pro',
                type: 'features_pro',
                enabled: true,
                order: 2,
                data: {
                    title: 'Apa yang Akan Kamu Pelajari?',
                    subtitle: 'Materi disusun khusus untuk pemula dengan bahasa yang mudah dipahami',
                    items: [
                        { title: 'Dasar Kepabeanan', text: 'Memahami Bea dan Cukai dengan bahasa sederhana dan contoh nyata.', icon: '📘', image: 'https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?auto=format&fit=crop&q=80&w=400' },
                        { title: 'Prosedur Impor', text: 'Langkah-langkah impor barang dan dokumen yang diperlukan.', icon: '📦', image: 'https://images.unsplash.com/photo-1578575437130-527eed3abbec?auto=format&fit=crop&q=80&w=400' },
                        { title: 'Perhitungan Pajak', text: 'Cara menghitung bea masuk, PPN, dan PPh impor dengan benar.', icon: '💰', image: 'https://images.unsplash.com/photo-1554224155-6726b3ff858f?auto=format&fit=crop&q=80&w=400' },
                        { title: 'Studi Kasus Nyata', text: 'Contoh kasus belanja online dan tips agar tidak kena masalah.', icon: '✅', image: 'https://images.unsplash.com/photo-1553413077-190dd305871c?auto=format&fit=crop&q=80&w=400' }
                    ]
                }
            },
            {
                id: 'speaker_pro',
                type: 'speaker_pro',
                enabled: true,
                order: 3,
                data: {
                    name: 'Indra Pratama',
                    title: 'Ex Auditor Bea & Cukai',
                    bio: 'Berpengalaman lebih dari 10 tahun di bidang kepabeanan. Akan membagikan insight praktis dan tips agar kamu tidak salah langkah saat berurusan dengan Bea & Cukai.',
                    image_url: '/images/Foto Narasumber-Indra Pratama.jpeg',
                    stats: [
                        { icon: '🕵️', label: '10+ Tahun Pengalaman' },
                        { icon: '📚', label: 'Praktisi Kepabeanan' }
                    ]
                }
            },
            {
                id: 'faq_pro',
                type: 'faq',
                enabled: true,
                order: 4,
                data: {
                    title: 'Pertanyaan Umum',
                    items: [
                        { question: 'Apakah webinar ini benar-benar gratis?', answer: 'Ya, 100% gratis tanpa biaya apa pun.' },
                        { question: 'Apakah pemula boleh ikut?', answer: 'Sangat boleh! Webinar ini memang dirancang untuk pemula.' },
                        { question: 'Apakah dapat sertifikat?', answer: 'Ya, semua peserta akan mendapatkan e-sertifikat.' },
                        { question: 'Bagaimana cara bergabung?', answer: 'Setelah mendaftar, link Google Meet akan dikirim via WhatsApp.' }
                    ]
                }
            },
            {
                id: 'cta_pro',
                type: 'cta_pro',
                enabled: true,
                order: 5,
                data: {
                    headline: 'Amankan Tempatmu Sekarang!',
                    subheadline: 'Kuota terbatas! Daftar gratis dan dapatkan akses penuh ke webinar beserta e-sertifikat.',
                    button_text: 'Daftar Webinar Gratis'
                }
            }
        ]
    },
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
                        { emoji: '💬', title: 'Konsultasi 1-on-1', value: 0, description: 'Diskusi privat via Google Meet' }
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
                        { feature: 'Sesi Google Meet Live', val_a: true, val_b: true },
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
