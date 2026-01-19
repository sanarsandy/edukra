<script setup lang="ts">
const props = defineProps<{
  block: any
  styles?: any
}>()
</script>

<template>
  <section class="py-20 bg-[#FAFAFA]" :style="{ fontFamily: block.data.font_family || 'Inter, sans-serif' }">
      <div class="max-w-6xl mx-auto px-4">
          <h2 class="text-2xl md:text-3xl font-bold text-center text-neutral-900 mb-12">
              {{ block.data.title }}
          </h2>

          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
              <div v-for="(item, idx) in block.data.items" :key="idx" class="bg-white p-6 rounded-xl shadow-sm border border-neutral-100 hover:shadow-md transition-shadow">
                   <div class="w-10 h-10 bg-[#FCE7F3] rounded-lg flex items-center justify-center text-xl mb-4 text-[#BE185D]">
                       {{ item.icon }}
                   </div>
                   <h3 class="font-bold text-neutral-900 mb-2">{{ item.title }}</h3>
                   <p class="text-sm text-neutral-500 leading-relaxed">{{ item.text }}</p>
              </div>
          </div>
      </div>
  </section>
</template>

<style scoped>
.font-inter {
    font-family: 'Inter', sans-serif;
}
</style>
