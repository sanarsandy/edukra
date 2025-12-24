<template>
  <button 
    :type="type"
    class="inline-flex items-center justify-center transition-all"
    :class="[btnClass, { 'opacity-70 cursor-not-allowed': loading || disabled }]"
    :disabled="loading || disabled"
    @click="$emit('click', $event)"
  >
    <!-- Loading Spinner -->
    <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    
    <!-- Button Content -->
    <slot>{{ text }}</slot>
  </button>
</template>

<script setup lang="ts">
defineProps<{
  text?: string
  loading?: boolean
  disabled?: boolean
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'secondary' | 'danger' | 'outline'
}>()

defineEmits(['click'])

const props = withDefaults(defineProps<{
  variant?: 'primary' | 'secondary' | 'danger' | 'outline'
  type?: 'button' | 'submit' | 'reset'
}>(), {
  variant: 'primary',
  type: 'button'
})

const btnClass = computed(() => {
  const base = 'px-4 py-2.5 text-sm font-medium rounded-lg'
  switch (props.variant) {
    case 'secondary':
      return `${base} bg-neutral-100 text-neutral-700 hover:bg-neutral-200`
    case 'danger':
      return `${base} bg-red-600 text-white hover:bg-red-700`
    case 'outline':
      return `${base} border border-neutral-300 text-neutral-700 hover:bg-neutral-50`
    default:
      return `${base} bg-primary-600 text-white hover:bg-primary-700`
  }
})
</script>
