<template>
  <div>
    <!-- Loading -->
    <div v-if="pending" class="min-h-screen flex items-center justify-center bg-neutral-900">
      <div class="animate-spin rounded-full h-12 w-12 border-4 border-primary-500 border-t-transparent"></div>
    </div>

    <!-- Not Found -->
    <div v-else-if="!campaign" class="min-h-screen flex items-center justify-center bg-neutral-900 text-white">
      <div class="text-center">
        <h1 class="text-4xl font-bold mb-4">Halaman Tidak Ditemukan</h1>
        <p class="text-neutral-400 mb-6">Campaign yang Anda cari tidak tersedia</p>
        <NuxtLink to="/" class="px-6 py-3 bg-primary-600 rounded-lg font-medium hover:bg-primary-700">Kembali ke Home</NuxtLink>
      </div>
    </div>

    <!-- Campaign Content -->
    <div v-else class="min-h-screen" :style="pageStyle">
      <template v-for="block in enabledBlocks" :key="block.id">
        
        <!-- ========== HERO SECTION (Improved) ========== -->
        <section v-if="block.type === 'hero'" class="relative min-h-[100svh] flex items-center justify-center text-white overflow-hidden" :style="heroStyle(block)">
          <!-- Gradient Overlay -->
          <div class="absolute inset-0 bg-gradient-to-b from-black/70 via-black/50 to-black/80"></div>
          
          <!-- Urgency Badge -->
          <div v-if="countdown" class="absolute top-4 left-1/2 -translate-x-1/2 z-20">
            <div class="bg-red-500 text-white px-4 py-2 rounded-full text-sm font-bold animate-pulse flex items-center gap-2">
              <span>🔥</span> Promo Berakhir dalam {{ countdown.days }}h {{ countdown.hours }}m
            </div>
          </div>
          
          <!-- Content -->
          <div class="relative z-10 text-center max-w-4xl mx-auto px-4 py-16 sm:py-24">
            <!-- Badge -->
            <div v-if="block.data.badge" class="inline-block bg-white/20 backdrop-blur px-4 py-1.5 rounded-full text-sm font-medium mb-6">
              {{ block.data.badge }}
            </div>
            
            <!-- Headline -->
            <h1 class="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-extrabold mb-4 sm:mb-6 leading-tight tracking-tight">
              {{ block.data.headline }}
            </h1>
            
            <!-- Subheadline -->
            <p class="text-base sm:text-lg md:text-xl text-white/90 mb-6 sm:mb-8 max-w-2xl mx-auto leading-relaxed">
              {{ block.data.subheadline }}
            </p>
            
            <!-- Price Preview -->
            <div v-if="displayPrice" class="mb-6">
              <span v-if="coursePrice > displayPrice" class="text-white/50 line-through text-lg mr-2">Rp {{ formatPrice(coursePrice) }}</span>
              <span class="text-2xl sm:text-3xl font-bold" :style="{color: campaignStyles.primaryColor}">Rp {{ formatPrice(displayPrice) }}</span>
            </div>
            
            <!-- CTA Button -->
            <button @click="handleBuy" class="w-full sm:w-auto px-8 sm:px-12 py-4 sm:py-5 text-white text-lg sm:text-xl font-bold rounded-2xl transition-all transform hover:scale-105 active:scale-95 shadow-2xl" :style="buttonStyle">
              {{ block.data.cta_text || '🚀 Daftar Sekarang' }}
            </button>
            
            <!-- Trust Indicator -->
            <p class="mt-4 text-white/60 text-sm">✓ Akses Seumur Hidup • ✓ Garansi 30 Hari</p>
          </div>
          
          <!-- Scroll Indicator -->
          <div class="absolute bottom-6 left-1/2 -translate-x-1/2 animate-bounce hidden sm:block">
            <svg class="w-6 h-6 text-white/50" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"/></svg>
          </div>
        </section>

        <!-- ========== COUNTDOWN SECTION ========== -->
        <section v-else-if="block.type === 'countdown' && countdown" class="py-8 sm:py-12 text-white" :style="{backgroundColor: campaignStyles.backgroundColor}">
          <div class="max-w-4xl mx-auto px-4 text-center">
            <p class="text-base sm:text-lg text-white/80 mb-4 sm:mb-6">{{ block.data.label || '⏰ Penawaran berakhir dalam:' }}</p>
            <div class="flex justify-center gap-2 sm:gap-4">
              <div v-for="(value, label) in {Hari: countdown.days, Jam: countdown.hours, Menit: countdown.minutes, Detik: countdown.seconds}" :key="label" class="bg-black/30 backdrop-blur rounded-xl p-3 sm:p-5 min-w-[60px] sm:min-w-[80px]">
                <div class="text-2xl sm:text-4xl md:text-5xl font-bold tabular-nums" :style="{color: campaignStyles.primaryColor}">{{ String(value).padStart(2, '0') }}</div>
                <div class="text-xs sm:text-sm text-white/60 mt-1">{{ label }}</div>
              </div>
            </div>
          </div>
        </section>

        <!-- ========== BENEFITS SECTION (Improved) ========== -->
        <section v-else-if="block.type === 'benefits'" class="bg-white py-12 sm:py-20">
          <div class="max-w-5xl mx-auto px-4">
            <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-center mb-3 text-neutral-900">{{ block.data.title || 'Yang Akan Anda Dapatkan' }}</h2>
            <p class="text-center text-neutral-500 mb-8 sm:mb-12 max-w-2xl mx-auto">{{ block.data.subtitle }}</p>
            
            <!-- Benefits Grid - Stack on mobile, 2 cols on tablet+ -->
            <div class="grid sm:grid-cols-2 gap-4 sm:gap-6">
              <div v-for="(item, idx) in block.data.items" :key="idx" 
                   class="flex items-start gap-4 p-4 sm:p-5 bg-gradient-to-br from-neutral-50 to-white rounded-2xl border border-neutral-100 hover:border-neutral-200 transition-all hover:shadow-lg group">
                <!-- Icon -->
                <div class="w-12 h-12 sm:w-14 sm:h-14 rounded-2xl flex items-center justify-center flex-shrink-0 transition-transform group-hover:scale-110" 
                     :style="{backgroundColor: campaignStyles.primaryColor + '15'}">
                  <svg class="w-6 h-6 sm:w-7 sm:h-7" :style="{color: campaignStyles.primaryColor}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
                  </svg>
                </div>
                <!-- Text -->
                <div class="flex-1 min-w-0">
                  <h4 v-if="item.title" class="font-bold text-neutral-900 mb-1">{{ item.title }}</h4>
                  <p class="text-neutral-600 text-sm sm:text-base leading-relaxed">{{ item.text }}</p>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- ========== PRICING SECTION (Enhanced) ========== -->
        <section v-else-if="block.type === 'pricing'" id="pricing" class="py-12 sm:py-20 text-white" :style="{background: `linear-gradient(135deg, ${campaignStyles.backgroundColor}, ${adjustColor(campaignStyles.backgroundColor, -30)})`}">
          <div class="max-w-lg mx-auto px-4">
            <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-center mb-2">🎁 Penawaran Spesial</h2>
            <p class="text-center text-white/70 mb-8">Khusus untuk Anda yang serius ingin berkembang</p>
            
            <!-- Pricing Card -->
            <div class="bg-white/10 backdrop-blur-xl rounded-3xl p-6 sm:p-8 border border-white/20 relative overflow-hidden">
              <!-- Popular Badge -->
              <div class="absolute -top-1 -right-8 bg-yellow-400 text-yellow-900 text-xs font-bold px-8 py-1 rotate-45 transform origin-center">
                BEST SELLER
              </div>
              
              <!-- Price Display -->
              <div class="text-center mb-6">
                <div v-if="coursePrice > displayPrice" class="flex items-center justify-center gap-2 mb-2">
                  <span class="text-white/50 line-through text-lg">Rp {{ formatPrice(coursePrice) }}</span>
                  <span class="bg-red-500 text-white text-xs font-bold px-2 py-1 rounded-full">
                    HEMAT {{ Math.round((1 - displayPrice/coursePrice) * 100) }}%
                  </span>
                </div>
                <div class="text-4xl sm:text-5xl md:text-6xl font-black" :style="{color: campaignStyles.primaryColor}">
                  Rp {{ formatPrice(displayPrice) }}
                </div>
                <p class="text-white/60 mt-2 text-sm">Sekali bayar, akses selamanya</p>
              </div>
              
              <!-- What's Included -->
              <div class="space-y-3 mb-8">
                <div v-for="feature in ['Akses semua materi kursus', 'Update materi gratis', 'Sertifikat kelulusan', 'Konsultasi via grup', 'Garansi 30 hari']" 
                     :key="feature" class="flex items-center gap-3">
                  <div class="w-5 h-5 rounded-full bg-green-500/20 flex items-center justify-center flex-shrink-0">
                    <svg class="w-3 h-3 text-green-400" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/></svg>
                  </div>
                  <span class="text-white/90 text-sm sm:text-base">{{ feature }}</span>
                </div>
              </div>
              
              <!-- CTA Button -->
              <button @click="handleBuy" class="w-full py-4 sm:py-5 text-white text-lg sm:text-xl font-bold rounded-2xl transition-all transform hover:scale-[1.02] active:scale-[0.98] shadow-xl" :style="buttonStyle">
                {{ block.data.cta_text || '🔥 Daftar Sekarang' }}
              </button>
              
              <!-- Trust badges -->
              <div class="flex items-center justify-center gap-4 mt-6 text-white/50 text-xs">
                <span>🔒 Pembayaran Aman</span>
                <span>💳 Berbagai Metode</span>
              </div>
            </div>
          </div>
        </section>

        <!-- ========== TESTIMONIALS SECTION (Carousel) ========== -->
        <section v-else-if="block.type === 'testimonials' && block.data.items?.length" class="bg-neutral-50 py-12 sm:py-20">
          <div class="max-w-6xl mx-auto px-4">
            <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-center mb-3 text-neutral-900">{{ block.data.title || 'Apa Kata Mereka?' }}</h2>
            <p class="text-center text-neutral-500 mb-8 sm:mb-12">Kisah sukses alumni kami</p>
            
            <!-- Testimonial Cards - Horizontal scroll on mobile -->
            <div class="flex gap-4 sm:gap-6 overflow-x-auto pb-4 snap-x snap-mandatory scrollbar-hide sm:grid sm:grid-cols-2 lg:grid-cols-3 sm:overflow-visible">
              <div v-for="(item, idx) in block.data.items" :key="idx" 
                   class="flex-shrink-0 w-[85vw] sm:w-auto snap-center bg-white p-5 sm:p-6 rounded-2xl shadow-sm border border-neutral-100 hover:shadow-lg transition-all">
                <!-- Stars -->
                <div class="flex gap-1 mb-4">
                  <span v-for="i in 5" :key="i" class="text-yellow-400">⭐</span>
                </div>
                
                <!-- Quote -->
                <p class="text-neutral-700 mb-6 leading-relaxed text-sm sm:text-base">"{{ item.text }}"</p>
                
                <!-- Author -->
                <div class="flex items-center gap-3">
                  <div class="w-11 h-11 rounded-full overflow-hidden flex-shrink-0" :style="{backgroundColor: campaignStyles.primaryColor + '20'}">
                    <img v-if="item.avatar" :src="item.avatar" class="w-full h-full object-cover" :alt="item.name"/>
                    <div v-else class="w-full h-full flex items-center justify-center">
                      <span :style="{color: campaignStyles.primaryColor}" class="font-bold text-lg">{{ item.name?.charAt(0) || '?' }}</span>
                    </div>
                  </div>
                  <div>
                    <div class="font-semibold text-neutral-900">{{ item.name }}</div>
                    <div v-if="item.role" class="text-sm text-neutral-500">{{ item.role }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- ========== VIDEO SECTION (NEW) ========== -->
        <section v-else-if="block.type === 'video'" class="bg-neutral-900 py-12 sm:py-20">
          <div class="max-w-4xl mx-auto px-4">
            <h2 v-if="block.data.title" class="text-2xl sm:text-3xl font-bold text-center mb-3 text-white">{{ block.data.title }}</h2>
            <p v-if="block.data.subtitle" class="text-center text-neutral-400 mb-8">{{ block.data.subtitle }}</p>
            
            <!-- Video Embed -->
            <div class="relative aspect-video rounded-2xl overflow-hidden shadow-2xl bg-neutral-800">
              <iframe v-if="block.data.youtube_url" 
                      :src="getYoutubeEmbedUrl(block.data.youtube_url)" 
                      class="absolute inset-0 w-full h-full"
                      frameborder="0" 
                      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" 
                      allowfullscreen>
              </iframe>
              <div v-else class="absolute inset-0 flex items-center justify-center">
                <div class="text-center text-neutral-500">
                  <span class="text-4xl mb-2">🎬</span>
                  <p>Video tidak tersedia</p>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- ========== TRUST BADGES (NEW) ========== -->
        <section v-else-if="block.type === 'trust'" class="bg-white py-8 sm:py-12 border-y border-neutral-100">
          <div class="max-w-5xl mx-auto px-4">
            <div class="flex flex-wrap items-center justify-center gap-6 sm:gap-10">
              <div class="flex items-center gap-2 text-neutral-600">
                <span class="text-2xl">🔒</span>
                <span class="text-sm font-medium">Pembayaran Aman</span>
              </div>
              <div class="flex items-center gap-2 text-neutral-600">
                <span class="text-2xl">💯</span>
                <span class="text-sm font-medium">Garansi 30 Hari</span>
              </div>
              <div class="flex items-center gap-2 text-neutral-600">
                <span class="text-2xl">🎓</span>
                <span class="text-sm font-medium">Sertifikat</span>
              </div>
              <div class="flex items-center gap-2 text-neutral-600">
                <span class="text-2xl">♾️</span>
                <span class="text-sm font-medium">Akses Selamanya</span>
              </div>
            </div>
          </div>
        </section>

        <!-- ========== INSTRUCTOR SECTION (Improved) ========== -->
        <section v-else-if="block.type === 'instructor'" class="bg-white py-12 sm:py-20">
          <div class="max-w-4xl mx-auto px-4">
            <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-center mb-8 sm:mb-12 text-neutral-900">👨‍🏫 Tentang Instruktur</h2>
            
            <div class="flex flex-col sm:flex-row items-center sm:items-start gap-6 sm:gap-8 bg-gradient-to-br from-neutral-50 to-white p-6 sm:p-8 rounded-3xl border border-neutral-100">
              <!-- Photo -->
              <div class="w-28 h-28 sm:w-36 sm:h-36 rounded-2xl overflow-hidden flex-shrink-0 shadow-lg">
                <img v-if="instructorData(block).avatar" :src="instructorData(block).avatar" class="w-full h-full object-cover" alt="Instructor"/>
                <div v-else class="w-full h-full bg-neutral-200 flex items-center justify-center text-5xl">👨‍🏫</div>
              </div>
              
              <!-- Info -->
              <div class="flex-1 text-center sm:text-left">
                <h3 class="text-xl sm:text-2xl font-bold text-neutral-900 mb-2">{{ instructorData(block).name }}</h3>
                <p class="text-neutral-600 leading-relaxed mb-4">{{ instructorData(block).bio }}</p>
                
                <!-- Stats -->
                <div class="flex flex-wrap justify-center sm:justify-start gap-4 text-sm">
                  <div class="flex items-center gap-1 text-neutral-500">
                    <span>👥</span> 1,000+ Siswa
                  </div>
                  <div class="flex items-center gap-1 text-neutral-500">
                    <span>⭐</span> 4.9 Rating
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- ========== FAQ SECTION (Accordion) ========== -->
        <section v-else-if="block.type === 'faq' && block.data.items?.length" class="bg-neutral-50 py-12 sm:py-20">
          <div class="max-w-3xl mx-auto px-4">
            <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-center mb-3 text-neutral-900">{{ block.data.title || '❓ Pertanyaan Umum' }}</h2>
            <p class="text-center text-neutral-500 mb-8 sm:mb-12">Jawaban untuk pertanyaan yang sering diajukan</p>
            
            <div class="space-y-3">
              <details v-for="(item, idx) in block.data.items" :key="idx" 
                       class="bg-white rounded-2xl overflow-hidden group shadow-sm border border-neutral-100">
                <summary class="font-semibold text-neutral-900 cursor-pointer p-5 flex justify-between items-center hover:bg-neutral-50 transition-colors">
                  <span class="pr-4">{{ item.question }}</span>
                  <svg class="w-5 h-5 text-neutral-400 group-open:rotate-180 transition-transform flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
                  </svg>
                </summary>
                <div class="px-5 pb-5 text-neutral-600 leading-relaxed border-t border-neutral-100 pt-4">
                  {{ item.answer }}
                </div>
              </details>
            </div>
          </div>
        </section>

        <!-- ========== CTA BANNER ========== -->
        <section v-else-if="block.type === 'cta_banner'" class="py-10 sm:py-16" :style="buttonStyle">
          <div class="max-w-4xl mx-auto px-4 text-center text-white">
            <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold mb-4">{{ block.data.headline || '🚀 Siap Untuk Memulai?' }}</h2>
            <p class="text-white/80 mb-6 max-w-2xl mx-auto">{{ block.data.subheadline || 'Jangan lewatkan kesempatan untuk meningkatkan skill Anda' }}</p>
            <button @click="handleBuy" class="px-8 sm:px-12 py-4 bg-white text-lg font-bold rounded-2xl hover:bg-neutral-100 transition-all transform hover:scale-105 active:scale-95 shadow-xl" :style="{color: campaignStyles.buttonColor}">
              {{ block.data.cta_text || 'Daftar Sekarang' }}
            </button>
          </div>
        </section>

        <!-- ========== STATISTICS BLOCK (NEW) ========== -->
        <section v-else-if="block.type === 'statistics'" class="py-12 sm:py-16 bg-white">
          <div class="max-w-5xl mx-auto px-4">
            <h2 v-if="block.data.title" class="text-2xl sm:text-3xl font-bold text-center mb-8">{{ block.data.title }}</h2>
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 sm:gap-6">
              <div v-for="(stat, idx) in block.data.items" :key="idx" class="text-center p-4 sm:p-6 bg-neutral-50 rounded-2xl">
                <div class="text-3xl sm:text-4xl md:text-5xl font-black mb-2" :style="{color: campaignStyles.primaryColor}">
                  {{ stat.value }}{{ stat.suffix || '' }}
                </div>
                <div class="text-sm sm:text-base font-medium" :style="{color: campaignStyles.textSecondaryColor}">{{ stat.label }}</div>
              </div>
            </div>
          </div>
        </section>

        <!-- ========== BONUS BLOCK (NEW) ========== -->
        <section v-else-if="block.type === 'bonus'" class="py-12 sm:py-20 bg-gradient-to-br from-yellow-50 to-orange-50">
          <div class="max-w-3xl mx-auto px-4">
            <div class="text-center mb-8 sm:mb-12">
              <span class="inline-block bg-yellow-400 text-yellow-900 text-xs font-bold px-3 py-1 rounded-full mb-4">🎁 BONUS SPESIAL</span>
              <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-neutral-900">{{ block.data.title || 'Bonus Eksklusif untuk Anda!' }}</h2>
            </div>
            
            <div class="space-y-4">
              <div v-for="(bonus, idx) in block.data.items" :key="idx" 
                   class="flex items-start gap-4 p-4 sm:p-5 bg-white rounded-2xl shadow-sm border border-orange-100">
                <div class="w-12 h-12 sm:w-14 sm:h-14 rounded-xl flex items-center justify-center flex-shrink-0 text-2xl" 
                     :style="{backgroundColor: campaignStyles.primaryColor + '15'}">
                  {{ bonus.emoji || '🎁' }}
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 flex-wrap">
                    <h4 class="font-bold">{{ bonus.title }}</h4>
                    <span v-if="bonus.value" class="text-sm line-through text-neutral-400">Rp {{ formatPrice(bonus.value) }}</span>
                  </div>
                  <p v-if="bonus.description" class="text-sm mt-1" :style="{color: campaignStyles.textSecondaryColor}">{{ bonus.description }}</p>
                </div>
                <div class="text-green-600 font-bold text-sm flex-shrink-0">GRATIS</div>
              </div>
            </div>
            
            <div v-if="block.data.total_value" class="mt-8 p-4 sm:p-6 bg-gradient-to-r from-green-500 to-emerald-600 rounded-2xl text-white text-center">
              <p class="text-sm opacity-80 mb-1">Total Nilai Bonus</p>
              <p class="text-2xl sm:text-3xl font-black">Rp {{ formatPrice(block.data.total_value) }}</p>
              <p class="text-sm font-medium mt-1">GRATIS untuk Anda!</p>
            </div>
          </div>
        </section>

        <!-- ========== CURRICULUM BLOCK (NEW) ========== -->
        <section v-else-if="block.type === 'curriculum'" class="py-12 sm:py-20 bg-white">
          <div class="max-w-3xl mx-auto px-4">
            <div class="text-center mb-8 sm:mb-12">
              <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-neutral-900">{{ block.data.title || '📚 Materi yang Akan Anda Pelajari' }}</h2>
              <p v-if="block.data.subtitle" class="text-neutral-500 mt-2">{{ block.data.subtitle }}</p>
            </div>
            
            <div class="space-y-3">
              <details v-for="(module, idx) in block.data.modules" :key="idx" 
                       class="bg-neutral-50 rounded-2xl overflow-hidden group border border-neutral-100">
                <summary class="p-4 sm:p-5 cursor-pointer flex items-center gap-4 hover:bg-neutral-100 transition-colors">
                  <div class="w-10 h-10 rounded-xl flex items-center justify-center font-bold text-white flex-shrink-0" 
                       :style="{backgroundColor: campaignStyles.primaryColor}">
                    {{ idx + 1 }}
                  </div>
                  <div class="flex-1 min-w-0">
                    <h4 class="font-semibold">{{ module.title }}</h4>
                    <p v-if="module.lessons_count" class="text-sm" :style="{color: campaignStyles.textSecondaryColor}">{{ module.lessons_count }} pelajaran</p>
                  </div>
                  <svg class="w-5 h-5 text-neutral-400 group-open:rotate-180 transition-transform flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
                  </svg>
                </summary>
                <div class="px-4 sm:px-5 pb-4 sm:pb-5 pt-0">
                  <ul v-if="module.lessons?.length" class="space-y-2 ml-14">
                    <li v-for="(lesson, lidx) in module.lessons" :key="lidx" 
                        class="flex items-center gap-2 text-sm" :style="{color: campaignStyles.textSecondaryColor}">
                      <svg class="w-4 h-4 text-green-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                      </svg>
                      {{ lesson }}
                    </li>
                  </ul>
                  <p v-else class="text-sm ml-14" :style="{color: campaignStyles.textSecondaryColor}">{{ module.description || 'Materi akan segera diupdate' }}</p>
                </div>
              </details>
            </div>
            
            <!-- Total Info -->
            <div v-if="block.data.total_lessons || block.data.total_hours" 
                 class="mt-8 flex flex-wrap justify-center gap-4 sm:gap-8 text-center">
              <div v-if="block.data.total_lessons" class="px-6 py-3 bg-neutral-100 rounded-xl">
                <span class="font-bold text-lg" :style="{color: campaignStyles.primaryColor}">{{ block.data.total_lessons }}</span>
                <span class="text-neutral-600 ml-1">Pelajaran</span>
              </div>
              <div v-if="block.data.total_hours" class="px-6 py-3 bg-neutral-100 rounded-xl">
                <span class="font-bold text-lg" :style="{color: campaignStyles.primaryColor}">{{ block.data.total_hours }}</span>
                <span class="text-neutral-600 ml-1">Jam Video</span>
              </div>
            </div>
          </div>
        </section>

        <!-- ========== GALLERY BLOCK (NEW) ========== -->
        <section v-else-if="block.type === 'gallery'" class="py-12 sm:py-20 bg-neutral-900 text-white overflow-hidden">
          <div class="max-w-7xl mx-auto px-4">
            <h2 v-if="block.data.title" class="text-2xl sm:text-3xl font-bold text-center mb-10">{{ block.data.title }}</h2>
            
            <!-- Scrollable Gallery -->
            <div class="flex gap-4 sm:gap-6 overflow-x-auto pb-8 snap-x snap-mandatory scrollbar-hide px-4 sm:px-0 -mx-4 sm:mx-0">
              <div v-for="(item, idx) in block.data.items" :key="idx" 
                   class="relative w-[80vw] sm:w-[500px] flex-shrink-0 aspect-video rounded-2xl overflow-hidden snap-center group">
                <img :src="item.image" :alt="item.caption || 'Gallery Image'" 
                     class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"/>
                <div v-if="item.caption" class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent flex items-end p-6">
                  <p class="font-medium text-lg">{{ item.caption }}</p>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- ========== ACHIEVEMENT/BADGES BLOCK (NEW) ========== -->
        <section v-else-if="block.type === 'achievement'" class="py-10 border-y border-neutral-100 bg-neutral-50" :style="pageStyle">
          <div class="max-w-5xl mx-auto px-4">
             <div class="flex flex-wrap justify-center gap-4 sm:gap-8">
               <div v-for="(badge, idx) in block.data.items" :key="idx" 
                    class="flex items-center gap-3 bg-white px-5 py-3 rounded-full shadow-sm border border-neutral-100 transform hover:-translate-y-1 transition-transform">
                 <span class="text-2xl">{{ badge.emoji }}</span>
                 <span class="font-bold text-neutral-800" :style="{color: campaignStyles.textPrimaryColor}">{{ badge.text }}</span>
               </div>
             </div>
          </div>
        </section>

        <!-- ========== COMPARISON TABLE BLOCK (NEW) ========== -->
        <section v-else-if="block.type === 'comparison'" class="py-16 sm:py-24 bg-white" :style="pageStyle">
          <div class="max-w-4xl mx-auto px-4">
            <h2 v-if="block.data.title" class="text-3xl font-bold text-center mb-12" :style="{color: campaignStyles.textPrimaryColor}">{{ block.data.title }}</h2>
            
            <div class="overflow-x-auto rounded-3xl shadow-xl border border-neutral-100">
              <table class="w-full text-left border-collapse">
                <thead>
                  <tr>
                    <th class="p-6 bg-neutral-50 text-neutral-500 font-medium w-1/3">{{ block.data.headers?.[0] || 'Fitur' }}</th>
                    <th class="p-6 bg-white text-center font-bold text-xl w-1/3 border-l border-neutral-100">{{ block.data.headers?.[1] || 'Basic' }}</th>
                    <th class="p-6 text-center font-bold text-xl w-1/3 border-l text-white relative overflow-hidden" :style="{backgroundColor: campaignStyles.primaryColor}">
                      <span class="relative z-10">{{ block.data.headers?.[2] || 'Pro' }}</span>
                      <div class="absolute inset-0 bg-white/10"></div>
                    </th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-neutral-100">
                  <tr v-for="(row, idx) in block.data.rows" :key="idx" class="hover:bg-neutral-50/50 transition-colors">
                    <td class="p-5 font-medium text-neutral-700">{{ row.feature }}</td>
                    
                    <!-- Column A -->
                    <td class="p-5 text-center border-l border-neutral-100">
                      <span v-if="row.val_a === true" class="inline-flex items-center justify-center w-8 h-8 rounded-full bg-green-100 text-green-600">✓</span>
                      <span v-else-if="row.val_a === false" class="inline-flex items-center justify-center w-8 h-8 rounded-full bg-neutral-100 text-neutral-400">✕</span>
                      <span v-else class="font-medium text-neutral-600">{{ row.text_a }}</span>
                    </td>
                    
                    <!-- Column B (Highlighted) -->
                    <td class="p-5 text-center border-l border-neutral-100 bg-primary-50/5">
                      <span v-if="row.val_b === true" class="inline-flex items-center justify-center w-8 h-8 rounded-full text-white shadow-sm" :style="{backgroundColor: campaignStyles.primaryColor}">✓</span>
                      <span v-else-if="row.val_b === false" class="inline-flex items-center justify-center w-8 h-8 rounded-full bg-neutral-100 text-neutral-400">✕</span>
                      <span v-else class="font-bold" :style="{color: campaignStyles.primaryColor}">{{ row.text_b }}</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </section>

        <!-- ========== FLOATING CHAT BLOCK (NEW) ========== -->
        <a v-else-if="block.type === 'floating_chat'"
           :href="`https://wa.me/${block.data.phone}?text=${encodeURIComponent(block.data.message || '')}`"
           target="_blank"
           class="fixed bottom-24 sm:bottom-8 right-4 sm:right-8 z-50 flex items-center gap-2 px-5 py-3 rounded-full text-white font-bold shadow-2xl hover:scale-105 active:scale-95 transition-all animate-bounce-slow"
           :class="{'!bottom-24': enabledBlocks.some(b => b.type === 'cta_banner')}" 
           style="background-color: #25D366">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12.031 6.172c-5.616 0-10.187 4.516-10.187 10.063 0 2.054.618 3.96 1.706 5.571l-1.551 5.059 5.378-1.401c1.517.781 3.22 1.258 5.053 1.258 5.618 0 10.189-4.515 10.189-10.063s-4.571-10.488-10.188-10.488z"/></svg>
          <span v-if="block.data.label">{{ block.data.label }}</span>
        </a>

        <!-- ========== SOCIAL PROOF BLOCK (NEW) ========== -->
        <div v-else-if="block.type === 'social_proof'" 
             class="fixed bottom-24 sm:bottom-8 left-4 sm:left-8 z-40 bg-white/95 backdrop-blur shadow-xl border border-neutral-200 rounded-2xl px-4 py-3 flex items-center gap-3 transition-all duration-500 transform hover:scale-105"
             :class="{'translate-y-2 opacity-0': !showSocialProof, '!bottom-24': enabledBlocks.some(b => b.type === 'cta_banner')}">
          <div class="relative flex-shrink-0">
             <span class="absolute -top-1 -right-1 flex h-3 w-3">
               <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
               <span class="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
             </span>
             <span class="text-2xl">👀</span>
          </div>
          <div class="text-sm">
             <div class="font-black text-neutral-900 leading-tight">
               <span class="tabular-nums">{{ liveViewerCount }}</span> orang
             </div>
             <div class="text-neutral-500 text-xs truncate max-w-[150px] sm:max-w-xs">{{ block.data.text || 'sedang melihat ini' }}</div>
          </div>
        </div>

      </template>

      <!-- ========== DEFAULT CTA FOOTER (only if no cta_banner block enabled) ========== -->
      <section v-if="!enabledBlocks.some(b => b.type === 'cta_banner')" class="py-12 text-white text-center" :style="buttonStyle">
        <div class="max-w-2xl mx-auto px-4">
          <h2 class="text-2xl sm:text-3xl font-bold mb-4">🎯 Siap Untuk Memulai?</h2>
          <p class="opacity-90 mb-6">Daftar sekarang dan tingkatkan skill Anda!</p>
          <button @click="handleBuy" class="px-10 py-4 bg-white font-bold rounded-2xl hover:bg-neutral-100 transition-all transform hover:scale-105 active:scale-95 shadow-lg" :style="{color: campaignStyles.buttonColor}">
            🚀 Daftar Sekarang
          </button>
        </div>
      </section>
    </div>
    <!-- ========== MOBILE STICKY CTA ========== -->
    <ClientOnly>
      <div v-if="campaign" class="fixed bottom-0 left-0 right-0 bg-white/95 backdrop-blur-xl border-t border-neutral-200 p-3 sm:p-4 md:hidden z-50 safe-area-pb">
        <div class="flex items-center gap-3">
          <div class="flex-1 min-w-0">
            <div class="text-sm text-neutral-500 truncate">{{ campaign.course?.title }}</div>
            <div class="font-bold text-lg" :style="{color: campaignStyles.primaryColor}">Rp {{ formatPrice(displayPrice) }}</div>
          </div>
          <button @click="handleBuy" class="px-6 py-3 text-white font-bold rounded-xl flex-shrink-0 min-h-[48px]" :style="buttonStyle">
            Beli Sekarang
          </button>
        </div>
      </div>
    </ClientOnly>

    <!-- Checkout Modal -->
    <CampaignCheckoutModal
      :show="showCheckoutModal"
      :campaign-id="campaign?.id || ''"
      :course-name="campaign?.course?.title || ''"
      :course-price="campaign?.course?.price || 0"
      :discount-price="campaign?.course?.discount_price"
      @close="showCheckoutModal = false"
      @success="handleCheckoutSuccess"
    />
  </div>
</template>

<script setup lang="ts">
interface Campaign {
  id: string
  slug: string
  title: string
  meta_description?: string
  course_id?: string
  end_date?: string
  blocks?: Block[]
  styles?: CampaignStyles
  course?: {
    id: string
    slug: string
    title: string
    price: number
    discount_price?: number
    instructor?: {
      full_name: string
      bio?: string
      avatar_url?: string
    }
  }
}

interface Block {
  id: string
  type: string
  enabled: boolean
  order: number
  data: any
}

interface CampaignStyles {
  primaryColor: string
  accentColor: string
  backgroundColor: string
  buttonColor: string
  fontFamily: string
  textPrimaryColor?: string
  textSecondaryColor?: string
  buttonStyle?: 'solid' | 'outline' | 'gradient'
  borderRadius?: 'sharp' | 'rounded' | 'pill'
  hasGradient?: boolean
}

definePageMeta({ layout: false })

const route = useRoute()
const router = useRouter()
const slug = computed(() => route.params.slug as string)

const config = useRuntimeConfig()
// Use internal API for SSR, public API for client
const apiBase = import.meta.server 
  ? (config.apiInternal || 'http://api:8080')
  : (config.public.apiBase || 'http://localhost:8080')

const { data: campaign, pending, refresh } = await useFetch<Campaign>(`/api/c/${slug.value}`, {
  baseURL: apiBase,
  key: `campaign-${slug.value}`,
  watch: [slug]
})

// Force refresh on mount if client to ensure fresh data
onMounted(() => {
  refresh()
})

useHead({
  title: computed(() => campaign.value?.title || 'Campaign'),
  meta: [{ name: 'description', content: campaign.value?.meta_description || '' }]
})

const defaultStyles: CampaignStyles = {
  primaryColor: '#6366f1',
  accentColor: '#f59e0b',
  backgroundColor: '#1f2937',
  buttonColor: '#6366f1',
  fontFamily: 'Inter',
  textPrimaryColor: '#111827',
  textSecondaryColor: '#4b5563',
  buttonStyle: 'solid',
  borderRadius: 'rounded',
  hasGradient: false
}

const campaignStyles = computed(() => {
  let styles = campaign.value?.styles
  if (typeof styles === 'string') {
    try { styles = JSON.parse(styles) } catch { styles = {} }
  }
  return { ...defaultStyles, ...(styles || {}) }
})

const enabledBlocks = computed(() => {
  let blocks = campaign.value?.blocks
  if (typeof blocks === 'string') {
    try { blocks = JSON.parse(blocks) } catch { blocks = [] }
  }
  return (blocks || []).filter((b: Block) => b.enabled).sort((a: Block, b: Block) => (a.order || 0) - (b.order || 0))
})

const adjustColor = (hex: string, percent: number) => {
  const num = parseInt(hex.replace('#', ''), 16)
  const r = Math.min(255, Math.max(0, (num >> 16) + percent))
  const g = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + percent))
  const b = Math.min(255, Math.max(0, (num & 0x0000FF) + percent))
  return `#${(1 << 24 | r << 16 | g << 8 | b).toString(16).slice(1)}`
}

const pageStyle = computed(() => ({
  fontFamily: campaignStyles.value.fontFamily + ', sans-serif',
  backgroundColor: campaignStyles.value.backgroundColor,
  backgroundImage: campaignStyles.value.hasGradient 
    ? `linear-gradient(135deg, ${campaignStyles.value.backgroundColor} 0%, ${adjustColor(campaignStyles.value.backgroundColor, -20)} 100%)` 
    : 'none',
  '--primary-color': campaignStyles.value.primaryColor,
  '--accent-color': campaignStyles.value.accentColor,
  '--text-primary': campaignStyles.value.textPrimaryColor || '#111827',
  '--text-secondary': campaignStyles.value.textSecondaryColor || '#4b5563',
  '--btn-radius': campaignStyles.value.borderRadius === 'pill' ? '999px' : campaignStyles.value.borderRadius === 'sharp' ? '0px' : '0.75rem',
  color: campaignStyles.value.textPrimaryColor || '#111827'
}))

const buttonStyle = computed(() => {
  const styles = campaignStyles.value
  const base = {
    borderRadius: styles.borderRadius === 'pill' ? '9999px' : styles.borderRadius === 'sharp' ? '0px' : '0.75rem',
    fontWeight: '600',
    transition: 'all 0.3s'
  }
  
  if (styles.buttonStyle === 'outline') {
    return {
      ...base,
      border: `2px solid ${styles.buttonColor}`,
      color: styles.buttonColor,
      backgroundColor: 'transparent'
    }
  } else if (styles.buttonStyle === 'gradient') {
    return {
       ...base,
       background: `linear-gradient(to right, ${styles.buttonColor}, ${styles.accentColor || styles.buttonColor})`,
       color: '#ffffff',
       border: 'none'
    }
  }
  
  // Solid default
  return {
    ...base,
    backgroundColor: styles.buttonColor,
    color: '#ffffff',
    border: 'none'
  }
})

// Course price from linked course
const coursePrice = computed(() => campaign.value?.course?.price || 0)

// Display price - use discount_price if available, otherwise course price
const displayPrice = computed(() => {
  const discountPrice = campaign.value?.course?.discount_price
  if (discountPrice && discountPrice > 0) {
    return discountPrice
  }
  return coursePrice.value
})

const heroStyle = (block: Block) => {
  const bg = block.data.background_image
  return bg 
    ? { backgroundImage: `url(${bg})`, backgroundSize: 'cover', backgroundPosition: 'center' }
    : { backgroundColor: campaignStyles.value.backgroundColor }
}


const instructorData = (block: Block) => {
  if (block.data.auto_fill && campaign.value?.course?.instructor) {
    const inst = campaign.value.course.instructor
    return { name: inst.full_name, bio: inst.bio || '', avatar: inst.avatar_url }
  }
  return { name: block.data.name || 'Instruktur', bio: block.data.bio || '', avatar: block.data.avatar }
}

const formatPrice = (price: number) => new Intl.NumberFormat('id-ID').format(price || 0)

// Helper to extract YouTube video ID and create embed URL
const getYoutubeEmbedUrl = (url: string) => {
  if (!url) return ''
  // Match various YouTube URL formats
  const regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|&v=)([^#&?]*).*/
  const match = url.match(regExp)
  const videoId = (match && match[2].length === 11) ? match[2] : null
  return videoId ? `https://www.youtube.com/embed/${videoId}` : ''
}

// Countdown
const countdown = ref<{days: number, hours: number, minutes: number, seconds: number} | null>(null)
let countdownInterval: ReturnType<typeof setInterval> | null = null

const startCountdown = () => {
  if (!campaign.value?.end_date) return
  
  const update = () => {
    const end = new Date(campaign.value!.end_date!).getTime()
    const now = Date.now()
    const diff = end - now

    if (diff <= 0) {
      countdown.value = { days: 0, hours: 0, minutes: 0, seconds: 0 }
      if (countdownInterval) clearInterval(countdownInterval)
      return
    }

    countdown.value = {
      days: Math.floor(diff / (1000 * 60 * 60 * 24)),
      hours: Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
      minutes: Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60)),
      seconds: Math.floor((diff % (1000 * 60)) / 1000)
    }
  }

  update()
  countdownInterval = setInterval(update, 1000)
}

// Social Proof Logic
const liveViewerCount = ref(0)
const showSocialProof = ref(false)
let socialProofInterval: ReturnType<typeof setInterval> | null = null

const initSocialProof = () => {
  const block = enabledBlocks.value.find(b => b.type === 'social_proof')
  if (!block) return

  const min = block.data.min_count || 10
  const max = block.data.max_count || 50
  
  // Initial value
  liveViewerCount.value = Math.floor(Math.random() * (max - min + 1)) + min
  showSocialProof.value = true

  // Simulation
  socialProofInterval = setInterval(() => {
    const change = Math.floor(Math.random() * 5) - 2 // -2 to +2
    let newVal = liveViewerCount.value + change
    if (newVal < min) newVal = min
    if (newVal > max) newVal = max
    liveViewerCount.value = newVal
  }, (block.data.interval || 5) * 1000)
}

onMounted(() => {
  startCountdown()
  initSocialProof()
  trackView()
})

onUnmounted(() => {
  if (countdownInterval) clearInterval(countdownInterval)
  if (socialProofInterval) clearInterval(socialProofInterval)
})

const trackView = async () => {
  if (!campaign.value?.id) return
  try {
    await $fetch(`/api/c/${campaign.value.id}/track`, {
      method: 'POST',
      baseURL: useRuntimeConfig().public.apiBase || 'http://localhost:8080',
      body: { event_type: 'view' }
    })
  } catch (e) {}
}

// Checkout modal state
const showCheckoutModal = ref(false)

const handleBuy = async () => {
  if (!campaign.value?.course) {
    alert('Kursus belum terhubung')
    return
  }
  
  // Track click
  try {
    await $fetch(`/api/c/${campaign.value.id}/track`, {
      method: 'POST',
      baseURL: config.public.apiBase || 'http://localhost:8080',
      body: { event_type: 'click' }
    })
  } catch (e) {}

  // Show checkout modal
  showCheckoutModal.value = true
}

const handleCheckoutSuccess = (result: { isFree: boolean, message?: string }) => {
  if (result.isFree) {
    alert(result.message || 'Berhasil terdaftar!')
    // Redirect to success page or course
    router.push(`/courses/${campaign.value?.course?.slug}`)
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&family=Poppins:wght@400;500;600;700&family=Playfair+Display:wght@400;600;700&family=Roboto:wght@400;500;700&display=swap');

/* Hide scrollbar for testimonials carousel */
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

/* Safe area padding for mobile bottom sticky */
.safe-area-pb {
  padding-bottom: max(env(safe-area-inset-bottom), 12px);
}

/* Smooth scroll snap */
.snap-x {
  scroll-snap-type: x mandatory;
}
.snap-center {
  scroll-snap-align: center;
}

/* Tabular nums for countdown */
.tabular-nums {
  font-variant-numeric: tabular-nums;
}
</style>
