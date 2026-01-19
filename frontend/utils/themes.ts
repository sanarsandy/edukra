export type ButtonStyle = 'solid' | 'outline' | 'gradient' | 'glass' | 'neo-brutalism' | 'soft-shadow'
export type ButtonShape = 'sharp' | 'rounded' | 'pill' | 'leaf'
export type CardStyle = 'flat' | 'card' | 'glass' | 'outline' | 'neumorphism'
export type BackgroundType = 'solid' | 'gradient' | 'image' | 'pattern' | 'animated'

export interface ThemePreset {
    id: string
    name: string
    thumbnail: string
    category: 'minimal' | 'creative' | 'professional' | 'dark' | 'sales' | 'premium'
    styles: {
        // Colors
        primaryColor: string
        accentColor: string
        backgroundColor: string
        buttonColor: string // Text color on button if solid, or border color if outline
        buttonTextColor?: string // Specific text color for buttons

        // Typography
        textPrimaryColor: string
        textSecondaryColor: string
        fontFamilyHeading: string // e.g., 'Playfair Display'
        fontFamilyBody: string    // e.g., 'Inter'

        // Component Styles
        buttonStyle: ButtonStyle
        buttonShape: ButtonShape
        cardStyle: CardStyle

        // Background
        backgroundType: BackgroundType
        backgroundImage?: string // URL of image or pattern
        backgroundGradient?: string // CSS Gradient string
        backgroundOverlay?: string // e.g., 'rgba(0,0,0,0.5)' for image legibility
    }
}

export const THEME_PRESETS: ThemePreset[] = [
    // =============================================
    // 🔥 HIGH-CONVERTING SALES THEMES
    // =============================================

    // Urgency Red - Flash sale, limited time offers
    {
        id: 'sales_urgency',
        name: '🔥 Urgency Red',
        category: 'sales',
        thumbnail: 'linear-gradient(135deg, #dc2626 0%, #991b1b 100%)',
        styles: {
            primaryColor: '#fbbf24', // Yellow for contrast
            accentColor: '#dc2626', // Red
            backgroundColor: '#0f0f0f',
            buttonColor: '#dc2626',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#ffffff',
            textSecondaryColor: '#d4d4d4',
            fontFamilyHeading: 'Poppins',
            fontFamilyBody: 'Inter',
            buttonStyle: 'solid',
            buttonShape: 'rounded',
            cardStyle: 'card',
            backgroundType: 'gradient',
            backgroundGradient: 'linear-gradient(180deg, #0f0f0f 0%, #1a1a1a 100%)',
        }
    },

    // Trust Builder - Course/coaching sales
    {
        id: 'sales_trust',
        name: '💎 Trust Builder',
        category: 'sales',
        thumbnail: 'linear-gradient(135deg, #1e40af 0%, #3b82f6 100%)',
        styles: {
            primaryColor: '#3b82f6', // Blue
            accentColor: '#10b981', // Green for trust
            backgroundColor: '#ffffff',
            buttonColor: '#3b82f6',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#1e293b',
            textSecondaryColor: '#64748b',
            fontFamilyHeading: 'Inter',
            fontFamilyBody: 'Inter',
            buttonStyle: 'solid',
            buttonShape: 'rounded',
            cardStyle: 'card',
            backgroundType: 'solid',
        }
    },

    // Webinar Pro - Live events, webinars
    {
        id: 'sales_webinar',
        name: '📺 Webinar Pro',
        category: 'sales',
        thumbnail: 'linear-gradient(135deg, #7c3aed 0%, #c026d3 100%)',
        styles: {
            primaryColor: '#c026d3', // Fuchsia
            accentColor: '#fbbf24', // Yellow highlight
            backgroundColor: '#0c0a1d',
            buttonColor: '#c026d3',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#f8fafc',
            textSecondaryColor: '#a5b4fc',
            fontFamilyHeading: 'Poppins',
            fontFamilyBody: 'Inter',
            buttonStyle: 'gradient',
            buttonShape: 'pill',
            cardStyle: 'glass',
            backgroundType: 'gradient',
            backgroundGradient: 'linear-gradient(135deg, #0c0a1d 0%, #1e1b4b 50%, #312e81 100%)',
        }
    },

    // Money Maker - High ticket, finance
    {
        id: 'sales_money',
        name: '💰 Money Maker',
        category: 'sales',
        thumbnail: 'linear-gradient(135deg, #15803d 0%, #22c55e 100%)',
        styles: {
            primaryColor: '#22c55e', // Green
            accentColor: '#fbbf24', // Gold
            backgroundColor: '#052e16',
            buttonColor: '#22c55e',
            buttonTextColor: '#052e16',
            textPrimaryColor: '#f0fdf4',
            textSecondaryColor: '#86efac',
            fontFamilyHeading: 'Montserrat',
            fontFamilyBody: 'Inter',
            buttonStyle: 'solid',
            buttonShape: 'rounded',
            cardStyle: 'outline',
            backgroundType: 'gradient',
            backgroundGradient: 'linear-gradient(180deg, #052e16 0%, #14532d 100%)',
        }
    },

    // Launch Mode - Product launches
    {
        id: 'sales_launch',
        name: '🚀 Launch Mode',
        category: 'sales',
        thumbnail: 'linear-gradient(135deg, #f97316 0%, #ea580c 100%)',
        styles: {
            primaryColor: '#f97316', // Orange
            accentColor: '#0ea5e9', // Sky blue
            backgroundColor: '#18181b',
            buttonColor: '#f97316',
            buttonTextColor: '#000000',
            textPrimaryColor: '#fafafa',
            textSecondaryColor: '#a1a1aa',
            fontFamilyHeading: 'Space Grotesk',
            fontFamilyBody: 'Inter',
            buttonStyle: 'solid',
            buttonShape: 'sharp',
            cardStyle: 'card',
            backgroundType: 'solid',
        }
    },

    // =============================================
    // 👑 PREMIUM/LUXURY THEMES
    // =============================================

    // Black Gold - Premium courses, VIP
    {
        id: 'premium_gold',
        name: '👑 Black Gold',
        category: 'premium',
        thumbnail: 'linear-gradient(135deg, #0f0f0f 0%, #1a1a1a 100%)',
        styles: {
            primaryColor: '#d4af37', // Gold
            accentColor: '#f5e6a3', // Light gold
            backgroundColor: '#0a0a0a',
            buttonColor: '#d4af37',
            buttonTextColor: '#0a0a0a',
            textPrimaryColor: '#ffffff',
            textSecondaryColor: '#a3a3a3',
            fontFamilyHeading: 'Playfair Display',
            fontFamilyBody: 'Lato',
            buttonStyle: 'solid',
            buttonShape: 'sharp',
            cardStyle: 'outline',
            backgroundType: 'gradient',
            backgroundGradient: 'linear-gradient(180deg, #0a0a0a 0%, #171717 100%)',
        }
    },

    // Elegant Rose - Women's products, beauty
    {
        id: 'premium_rose',
        name: '🌹 Elegant Rose',
        category: 'premium',
        thumbnail: 'linear-gradient(135deg, #fdf2f8 0%, #fce7f3 100%)',
        styles: {
            primaryColor: '#be185d', // Rose
            accentColor: '#d4af37', // Gold accent
            backgroundColor: '#fdf2f8',
            buttonColor: '#be185d',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#1f2937',
            textSecondaryColor: '#6b7280',
            fontFamilyHeading: 'Playfair Display',
            fontFamilyBody: 'Quicksand',
            buttonStyle: 'soft-shadow',
            buttonShape: 'pill',
            cardStyle: 'card',
            backgroundType: 'solid',
        }
    },

    // Midnight Luxury - High-end coaching
    {
        id: 'premium_midnight',
        name: '🌙 Midnight Luxury',
        category: 'premium',
        thumbnail: 'linear-gradient(135deg, #1e1b4b 0%, #312e81 100%)',
        styles: {
            primaryColor: '#a78bfa', // Violet
            accentColor: '#c4b5fd',
            backgroundColor: '#0f0d1a',
            buttonColor: '#a78bfa',
            buttonTextColor: '#0f0d1a',
            textPrimaryColor: '#f5f3ff',
            textSecondaryColor: '#c4b5fd',
            fontFamilyHeading: 'DM Serif Display',
            fontFamilyBody: 'DM Sans',
            buttonStyle: 'solid',
            buttonShape: 'rounded',
            cardStyle: 'glass',
            backgroundType: 'gradient',
            backgroundGradient: 'linear-gradient(135deg, #0f0d1a 0%, #1e1b4b 100%)',
        }
    },

    // =============================================
    // 🎨 CREATIVE SERIES
    // =============================================

    // Unicorn Dream
    {
        id: 'creative_gradient',
        name: '🦄 Unicorn Dream',
        category: 'creative',
        thumbnail: 'linear-gradient(135deg, #FF9A9E 0%, #FECFEF 100%)',
        styles: {
            primaryColor: '#ffffff',
            accentColor: '#ff6b6b',
            backgroundColor: '#fff0f5',
            buttonColor: '#ec4899',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#1f2937',
            textSecondaryColor: '#6b7280',
            fontFamilyHeading: 'Poppins',
            fontFamilyBody: 'Quicksand',
            buttonStyle: 'gradient',
            buttonShape: 'pill',
            cardStyle: 'glass',
            backgroundType: 'gradient',
            backgroundGradient: 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 99%, #fecfef 100%)',
        }
    },

    // Neo Pop - Bold brutalism
    {
        id: 'creative_neo',
        name: '⚡ Neo Pop',
        category: 'creative',
        thumbnail: '#fbbf24',
        styles: {
            primaryColor: '#000000',
            accentColor: '#a855f7',
            backgroundColor: '#fbbf24',
            buttonColor: '#000000',
            buttonTextColor: '#fbbf24',
            textPrimaryColor: '#000000',
            textSecondaryColor: '#333333',
            fontFamilyHeading: 'Space Grotesk',
            fontFamilyBody: 'Inter',
            buttonStyle: 'neo-brutalism',
            buttonShape: 'sharp',
            cardStyle: 'card',
            backgroundType: 'solid',
        }
    },

    // Sunset Vibes
    {
        id: 'creative_sunset',
        name: '🌅 Sunset Vibes',
        category: 'creative',
        thumbnail: 'linear-gradient(135deg, #f97316 0%, #ec4899 100%)',
        styles: {
            primaryColor: '#ffffff',
            accentColor: '#fbbf24',
            backgroundColor: '#1a1a2e',
            buttonColor: '#f97316',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#ffffff',
            textSecondaryColor: '#d4d4d8',
            fontFamilyHeading: 'Poppins',
            fontFamilyBody: 'Inter',
            buttonStyle: 'gradient',
            buttonShape: 'rounded',
            cardStyle: 'glass',
            backgroundType: 'gradient',
            backgroundGradient: 'linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%)',
        }
    },

    // =============================================
    // 🌑 DARK THEMES
    // =============================================

    // Cyber Glitch - Tech/gaming
    {
        id: 'dark_cyber',
        name: '🎮 Cyber Glitch',
        category: 'dark',
        thumbnail: '#0f172a',
        styles: {
            primaryColor: '#00f2ff',
            accentColor: '#bd00ff',
            backgroundColor: '#050505',
            buttonColor: '#00f2ff',
            buttonTextColor: '#000000',
            textPrimaryColor: '#ffffff',
            textSecondaryColor: '#94a3b8',
            fontFamilyHeading: 'Orbitron',
            fontFamilyBody: 'Inter',
            buttonStyle: 'outline',
            buttonShape: 'sharp',
            cardStyle: 'outline',
            backgroundType: 'solid',
        }
    },

    // Sleek Dark - SaaS, modern apps
    {
        id: 'dark_sleek',
        name: '🖤 Sleek Dark',
        category: 'dark',
        thumbnail: 'linear-gradient(135deg, #18181b 0%, #27272a 100%)',
        styles: {
            primaryColor: '#6366f1', // Indigo
            accentColor: '#818cf8',
            backgroundColor: '#09090b',
            buttonColor: '#6366f1',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#fafafa',
            textSecondaryColor: '#a1a1aa',
            fontFamilyHeading: 'Inter',
            fontFamilyBody: 'Inter',
            buttonStyle: 'solid',
            buttonShape: 'rounded',
            cardStyle: 'card',
            backgroundType: 'solid',
        }
    },

    // =============================================
    // 💼 PROFESSIONAL THEMES
    // =============================================

    // Executive Blue
    {
        id: 'pro_executive',
        name: '💼 Executive Blue',
        category: 'professional',
        thumbnail: '#1e3a8a',
        styles: {
            primaryColor: '#1e40af',
            accentColor: '#3b82f6',
            backgroundColor: '#f8fafc',
            buttonColor: '#1e40af',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#0f172a',
            textSecondaryColor: '#475569',
            fontFamilyHeading: 'Montserrat',
            fontFamilyBody: 'Inter',
            buttonStyle: 'solid',
            buttonShape: 'rounded',
            cardStyle: 'card',
            backgroundType: 'solid',
        }
    },

    // Deep Forest - Nature, wellness
    {
        id: 'pro_forest',
        name: '🌲 Deep Forest',
        category: 'professional',
        thumbnail: '#064e3b',
        styles: {
            primaryColor: '#059669',
            accentColor: '#d97706',
            backgroundColor: '#f0fdf4',
            buttonColor: '#059669',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#14532d',
            textSecondaryColor: '#166534',
            fontFamilyHeading: 'DM Serif Display',
            fontFamilyBody: 'DM Sans',
            buttonStyle: 'soft-shadow',
            buttonShape: 'leaf',
            cardStyle: 'flat',
            backgroundType: 'solid',
        }
    },

    // Corporate Clean
    {
        id: 'pro_corporate',
        name: '🏢 Corporate Clean',
        category: 'professional',
        thumbnail: '#374151',
        styles: {
            primaryColor: '#374151', // Gray 700
            accentColor: '#0ea5e9', // Sky blue
            backgroundColor: '#ffffff',
            buttonColor: '#374151',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#111827',
            textSecondaryColor: '#6b7280',
            fontFamilyHeading: 'Inter',
            fontFamilyBody: 'Inter',
            buttonStyle: 'solid',
            buttonShape: 'rounded',
            cardStyle: 'card',
            backgroundType: 'solid',
        }
    },

    // =============================================
    // ✨ MINIMALIST THEMES
    // =============================================

    // Monochrome Clean
    {
        id: 'minimal_mono',
        name: '⬛ Monochrome',
        category: 'minimal',
        thumbnail: 'linear-gradient(to bottom right, #ffffff, #f3f4f6)',
        styles: {
            primaryColor: '#000000',
            accentColor: '#525252',
            backgroundColor: '#ffffff',
            buttonColor: '#000000',
            buttonTextColor: '#ffffff',
            textPrimaryColor: '#171717',
            textSecondaryColor: '#525252',
            fontFamilyHeading: 'Inter',
            fontFamilyBody: 'Inter',
            buttonStyle: 'solid',
            buttonShape: 'sharp',
            cardStyle: 'flat',
            backgroundType: 'solid',
        }
    },

    // Organic Beige
    {
        id: 'minimal_beige',
        name: '🍂 Organic Beige',
        category: 'minimal',
        thumbnail: '#f5f5dc',
        styles: {
            primaryColor: '#78716c',
            accentColor: '#a8a29e',
            backgroundColor: '#fafaf9',
            buttonColor: '#57534e',
            buttonTextColor: '#fafaf9',
            textPrimaryColor: '#44403c',
            textSecondaryColor: '#78716c',
            fontFamilyHeading: 'Playfair Display',
            fontFamilyBody: 'Lato',
            buttonStyle: 'soft-shadow',
            buttonShape: 'rounded',
            cardStyle: 'card',
            backgroundType: 'solid',
        }
    },

    // Paper White - Ultra clean
    {
        id: 'minimal_paper',
        name: '📄 Paper White',
        category: 'minimal',
        thumbnail: '#fafafa',
        styles: {
            primaryColor: '#18181b',
            accentColor: '#3b82f6',
            backgroundColor: '#fafafa',
            buttonColor: '#18181b',
            buttonTextColor: '#fafafa',
            textPrimaryColor: '#18181b',
            textSecondaryColor: '#71717a',
            fontFamilyHeading: 'Poppins',
            fontFamilyBody: 'Inter',
            buttonStyle: 'outline',
            buttonShape: 'pill',
            cardStyle: 'flat',
            backgroundType: 'solid',
        }
    },
]

