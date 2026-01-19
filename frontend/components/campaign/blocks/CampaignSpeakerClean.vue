<script setup lang="ts">
const props = defineProps<{
  block: any
  styles?: any
}>()
</script>

<template>
  <section class="max-w-6xl mx-auto px-4 py-16 bg-white" :style="{ fontFamily: block.data.font_family || 'Inter, sans-serif' }">
      <div class="flex flex-col md:flex-row items-center gap-10 md:gap-16">
          
          <!-- IMAGE LEFT -->
           <div class="w-full md:w-1/3">
                <div class="aspect-[3/4] rounded-2xl overflow-hidden bg-neutral-100 shadow-lg relative">
                    <img :src="block.data.image_url || 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=crop&q=80&w=800'" alt="Speaker" class="object-cover w-full h-full" />
                </div>
           </div>

           <!-- CONTENT RIGHT -->
           <div class="w-full md:w-2/3 text-center md:text-left">
               <!-- Badge -->
                <div class="inline-block bg-blue-100 text-blue-800 font-bold text-[10px] px-3 py-1 rounded-sm mb-4 uppercase tracking-wider">
                    PEMBICARA
                </div>

                <h2 class="text-3xl font-bold text-neutral-900 mb-2">{{ block.data.name }}</h2>
                <p class="text-blue-600 font-medium text-sm mb-6">{{ block.data.title }}</p>
                
                <p class="text-neutral-600 leading-relaxed mb-8 max-w-2xl">
                    {{ block.data.bio }}
                </p>

                <div class="flex flex-wrap justify-center md:justify-start gap-8 border-t border-neutral-100 pt-8">
                    <div v-for="(stat, idx) in block.data.stats" :key="idx" class="flex items-center gap-3">
                         <div class="text-2xl">{{ stat.icon }}</div>
                         <div class="text-sm font-bold text-neutral-800">{{ stat.label }}</div>
                    </div>
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
