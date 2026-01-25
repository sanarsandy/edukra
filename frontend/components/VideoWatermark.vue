<template>
  <div class="video-watermark pointer-events-none select-none" :class="positionClass">
    <div class="watermark-text" :style="textStyle">
      {{ displayText }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface Props {
  userEmail?: string
  userId?: string
  opacity?: number
  rotateInterval?: number // seconds between position changes
}

const props = withDefaults(defineProps<Props>(), {
  opacity: 0.15,
  rotateInterval: 30
})

// Position rotation
const positions = [
  'top-left',
  'top-right', 
  'bottom-left',
  'bottom-right',
  'center'
] as const

const currentPositionIndex = ref(0)
let rotationTimer: ReturnType<typeof setInterval> | null = null

const positionClass = computed(() => {
  return `position-${positions[currentPositionIndex.value]}`
})

const displayText = computed(() => {
  if (props.userEmail) {
    // Partially mask email for privacy but still identifiable
    const [local, domain] = props.userEmail.split('@')
    if (local && domain) {
      const maskedLocal = local.length > 3 
        ? `${local.slice(0, 2)}***${local.slice(-1)}`
        : local
      return `${maskedLocal}@${domain}`
    }
    return props.userEmail
  }
  if (props.userId) {
    return `ID: ${props.userId.slice(0, 8)}...`
  }
  return ''
})

const textStyle = computed(() => ({
  opacity: props.opacity,
  fontSize: '14px',
  fontFamily: 'monospace',
  color: 'white',
  textShadow: '1px 1px 2px rgba(0,0,0,0.5)',
  whiteSpace: 'nowrap'
}))

const rotatePosition = () => {
  currentPositionIndex.value = (currentPositionIndex.value + 1) % positions.length
}

onMounted(() => {
  if (props.rotateInterval > 0) {
    rotationTimer = setInterval(rotatePosition, props.rotateInterval * 1000)
  }
})

onUnmounted(() => {
  if (rotationTimer) {
    clearInterval(rotationTimer)
  }
})
</script>

<style scoped>
.video-watermark {
  position: absolute;
  z-index: 10;
  padding: 8px 12px;
  transition: all 0.5s ease-in-out;
}

.watermark-text {
  transform: rotate(-15deg);
}

/* Position variants */
.position-top-left {
  top: 10%;
  left: 5%;
}

.position-top-right {
  top: 10%;
  right: 5%;
}

.position-bottom-left {
  bottom: 15%;
  left: 5%;
}

.position-bottom-right {
  bottom: 15%;
  right: 5%;
}

.position-center {
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.position-center .watermark-text {
  transform: rotate(-15deg) scale(1.2);
}
</style>
