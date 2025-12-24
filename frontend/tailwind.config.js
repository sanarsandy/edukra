/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./app.vue",
    "./error.vue",
  ],
  theme: {
    extend: {
      colors: {
        // EDUKRA "Calm Energy" Design System
        // Muted tones untuk kesan "Mahal" dan nyaman untuk durasi layar lama

        // Primary: Dark Blue (Biru Tua - lebih jelas biru, bukan hitam)
        primary: {
          50: '#e8f0f7',
          100: '#c5d9eb',
          200: '#9fc1de',
          300: '#79a9d1',
          400: '#5391c4',
          500: '#3478b0',
          600: '#1E4976',  // Dark Blue - MAIN (jelas biru tua)
          700: '#1a3f67',
          800: '#163558',
          900: '#122b49',
        },
        // Secondary: Soft Teal (Penyeimbang - modern tech, sejuk)
        secondary: {
          50: '#f0f6f9',
          100: '#d9e8ef',
          200: '#b3d1df',
          300: '#8dbacf',
          400: '#5D8AA8',  // Soft Teal - MAIN
          500: '#4a7a9a',
          600: '#3d6580',
          700: '#305066',
          800: '#233b4c',
          900: '#162632',
        },
        // Accent: Muted Clay/Coral (CTA, Notifikasi, Highlight)
        accent: {
          50: '#fdf5f2',
          100: '#fae8e2',
          200: '#f5d1c5',
          300: '#e8b5a3',
          400: '#D98C75',  // Muted Coral - MAIN (Energi Cakra)
          500: '#c47962',
          600: '#a86550',
          700: '#8c513e',
          800: '#703d2c',
          900: '#54291a',
        },
        // Success: Sage Green (Progress bar, Growth, Success state)
        success: {
          50: '#f2f7f5',
          100: '#e0ece7',
          200: '#c1d9cf',
          300: '#a2c6b7',
          400: '#7FA99B',  // Sage Green - MAIN
          500: '#6a9485',
          600: '#557d6f',
          700: '#406659',
          800: '#2b4f43',
          900: '#16382d',
        },
        // Warm: Highlights (softer amber)
        warm: {
          50: '#fefbf5',
          100: '#fdf5e5',
          200: '#fbe8c5',
          300: '#f9dba5',
          400: '#f5c465',
          500: '#f0ad25',
          600: '#d8981e',
          700: '#b07c18',
          800: '#886012',
          900: '#60440c',
        },
        // Neutral: Canvas - Clean & Modern
        neutral: {
          50: '#F9FAFB',   // Background (off-white, tidak silau)
          100: '#f3f4f6',
          200: '#E5E7EB',  // Border/Divider
          300: '#d1d5db',
          400: '#9ca3af',
          500: '#6b7280',
          600: '#4b5563',
          700: '#374151',  // Text Body (lebih nyaman dari hitam legam)
          800: '#1f2937',
          900: '#111827',
        },
        // Admin: Dark Blue theme (sama dengan primary)
        admin: {
          50: '#e8f0f7',
          100: '#c5d9eb',
          200: '#9fc1de',
          300: '#79a9d1',
          400: '#5391c4',
          500: '#3478b0',
          600: '#1E4976',  // Dark Blue
          700: '#1a3f67',
          800: '#163558',
          900: '#122b49',
        },
      },
      fontFamily: {
        // Body Font: Inter - high readability
        sans: ['Inter', 'system-ui', 'sans-serif'],
        // Heading Font: Plus Jakarta Sans - modern, geometric
        display: ['Plus Jakarta Sans', 'Inter', 'system-ui', 'sans-serif'],
      },
      borderRadius: {
        // EDUKRA rounded corners: 8-12px for cards/buttons
        'edukra': '10px',
        'edukra-sm': '8px',
        'edukra-lg': '12px',
      },
      boxShadow: {
        // Soft shadows for "floating" effect
        'soft': '0 1px 3px 0 rgb(0 0 0 / 0.04), 0 1px 2px -1px rgb(0 0 0 / 0.04)',
        'soft-md': '0 4px 6px -1px rgb(0 0 0 / 0.05), 0 2px 4px -2px rgb(0 0 0 / 0.05)',
        'soft-lg': '0 10px 15px -3px rgb(0 0 0 / 0.05), 0 4px 6px -4px rgb(0 0 0 / 0.05)',
        'edukra': '0 4px 6px -1px rgba(0, 0, 0, 0.05)',
        'edukra-hover': '0 8px 12px -2px rgba(0, 0, 0, 0.08)',
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-out',
        'slide-up': 'slideUp 0.5s ease-out',
        'spin-slow': 'spin 2s linear infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
}
