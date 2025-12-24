<template>
  <div class="flex flex-col items-center justify-center py-12 px-4 text-center">
    <!-- Icon -->
    <div class="w-20 h-20 mb-6 rounded-full flex items-center justify-center" :class="iconBgClass">
      <slot name="icon">
        <svg class="w-10 h-10" :class="iconClass" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/>
        </svg>
      </slot>
    </div>
    
    <!-- Title -->
    <h3 class="text-lg font-semibold text-neutral-900 mb-2">{{ title }}</h3>
    
    <!-- Description -->
    <p class="text-neutral-500 max-w-sm mb-6">{{ description }}</p>
    
    <!-- Action Button -->
    <slot name="action">
      <button v-if="actionText" @click="$emit('action')" class="btn-primary">
        {{ actionText }}
      </button>
    </slot>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description: string
  actionText?: string
  variant?: 'primary' | 'accent' | 'neutral'
}>()

defineEmits(['action'])

const props = withDefaults(defineProps<{
  variant?: 'primary' | 'accent' | 'neutral'
}>(), {
  variant: 'primary'
})

const iconBgClass = computed(() => {
  switch (props.variant) {
    case 'accent': return 'bg-accent-100'
    case 'neutral': return 'bg-neutral-100'
    default: return 'bg-primary-100'
  }
})

const iconClass = computed(() => {
  switch (props.variant) {
    case 'accent': return 'text-accent-500'
    case 'neutral': return 'text-neutral-400'
    default: return 'text-primary-500'
  }
})
</script>
