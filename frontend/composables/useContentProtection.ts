/**
 * Content Protection Composable
 * Provides utilities to prevent easy downloading of course materials
 */

export const useContentProtection = () => {
    /**
     * Disable right-click context menu on an element
     */
    const disableContextMenu = (element: HTMLElement | null) => {
        if (!element) return

        element.addEventListener('contextmenu', (e) => {
            e.preventDefault()
            return false
        })
    }

    /**
     * Disable text selection on an element
     */
    const disableTextSelection = (element: HTMLElement | null) => {
        if (!element) return

        element.style.userSelect = 'none'
        element.style.webkitUserSelect = 'none'
            ; (element.style as any).msUserSelect = 'none'
        element.style.pointerEvents = 'auto' // Keep clickable
    }

    /**
     * Disable keyboard shortcuts for copying/printing
     */
    const disableCopyShortcuts = (element: HTMLElement | null) => {
        if (!element) return

        element.addEventListener('keydown', (e) => {
            // Ctrl+C, Ctrl+P, Ctrl+S, Ctrl+U (view source)
            if (e.ctrlKey || e.metaKey) {
                if (['c', 'p', 's', 'u', 'a'].includes(e.key.toLowerCase())) {
                    e.preventDefault()
                    return false
                }
            }
            // F12 (DevTools) - optional, can be annoying
            if (e.key === 'F12') {
                e.preventDefault()
                return false
            }
        })
    }

    /**
     * Disable drag on an element (prevents dragging images/videos)
     */
    const disableDrag = (element: HTMLElement | null) => {
        if (!element) return

        element.addEventListener('dragstart', (e) => {
            e.preventDefault()
            return false
        })
    }

    /**
     * Apply all protections to a content container
     */
    const protectContent = (element: HTMLElement | null) => {
        if (!element) return

        disableContextMenu(element)
        disableTextSelection(element)
        disableCopyShortcuts(element)
        disableDrag(element)

        // Add CSS class for additional protection
        element.classList.add('protected-content')
    }

    /**
     * Setup protections on mount (call in onMounted)
     */
    const setupProtection = (containerSelector: string) => {
        if (typeof document === 'undefined') return

        const container = document.querySelector(containerSelector) as HTMLElement
        if (container) {
            protectContent(container)
        }
    }

    /**
     * Add protection styles to document head
     */
    const injectProtectionStyles = () => {
        if (typeof document === 'undefined') return

        const styleId = 'content-protection-styles'
        if (document.getElementById(styleId)) return

        const style = document.createElement('style')
        style.id = styleId
        style.textContent = `
      .protected-content {
        -webkit-user-select: none !important;
        -moz-user-select: none !important;
        -ms-user-select: none !important;
        user-select: none !important;
        -webkit-touch-callout: none !important;
      }
      
      .protected-content img,
      .protected-content video {
        pointer-events: none;
        -webkit-user-drag: none;
        user-drag: none;
      }
      
      .protected-content video::-webkit-media-controls-download-button,
      .protected-content video::-webkit-media-controls-enclosure,
      .protected-content video::-webkit-media-controls-panel {
        /* Allow controls but hide download */
      }
      
      /* Hide download button in video controls */
      video::-webkit-media-controls-download-button {
        display: none !important;
      }
      
      /* PDF iframe protection */
      .protected-pdf-container {
        position: relative;
        overflow: hidden;
      }
      
      .protected-pdf-container::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        pointer-events: none;
      }
    `
        document.head.appendChild(style)
    }

    return {
        disableContextMenu,
        disableTextSelection,
        disableCopyShortcuts,
        disableDrag,
        protectContent,
        setupProtection,
        injectProtectionStyles
    }
}
