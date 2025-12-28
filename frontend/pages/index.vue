<template>
  <div class="min-h-screen bg-white">
    <!-- Navigation -->
    <nav class="fixed top-0 left-0 right-0 z-50 bg-white border-b border-neutral-100">
      <div class="container-custom">
        <div class="flex items-center justify-between h-16 md:h-18">
          <!-- Logo -->
          <div class="flex items-center space-x-2">
            <img src="/logo.png" alt="EDUKRA Logo" class="w-10 h-10 object-contain" />
            <span class="font-display font-bold text-xl text-neutral-900">EDUKRA</span>
          </div>
          
          <!-- Desktop Navigation -->
          <div class="hidden md:flex items-center space-x-8">
            <a href="#features" class="text-neutral-600 hover:text-primary-600 transition-colors font-medium">Fitur</a>
            <a href="#courses" class="text-neutral-600 hover:text-primary-600 transition-colors font-medium">Kursus</a>
            <a href="#testimonials" class="text-neutral-600 hover:text-primary-600 transition-colors font-medium">Testimoni</a>
            <!-- <a href="#pricing" class="text-neutral-600 hover:text-primary-600 transition-colors font-medium">Harga</a> -->
          </div>
          
          <!-- Auth Buttons -->
          <div class="flex items-center space-x-3">
            <NuxtLink to="/login" class="btn-outline hidden sm:inline-flex">
              Masuk
            </NuxtLink>
            <NuxtLink to="/register" class="inline-flex items-center justify-center px-6 py-3 text-sm font-semibold text-white bg-primary-600 rounded-xl hover:bg-primary-700 transition-all">
              Daftar Gratis
            </NuxtLink>
          </div>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="pt-32 md:pt-40 pb-16 md:pb-24 bg-neutral-50">
      <div class="container-custom">
        <div class="grid lg:grid-cols-2 gap-12 lg:gap-16 items-center">
            <!-- Hero Content -->
           <div class="animate-fade-in relative z-10 order-2 lg:order-1">
             <div class="inline-flex items-center px-4 py-2 bg-primary-50 rounded-full mb-6 animate-fade-in border border-primary-100">
               <span class="w-2 h-2 bg-warm-500 rounded-full mr-2 animate-pulse"></span>
               <span class="text-sm font-medium text-primary-700">Platform Belajar Terpercaya</span>
             </div>
             
             <h1 class="text-3xl md:text-4xl lg:text-5xl font-display font-bold leading-tight mb-6 animate-fade-in text-neutral-900">
               Tingkatkan Skill, Buka <br/>
               <span class="text-primary-600 relative">
                 Peluang Karir Baru
                 <svg class="absolute -bottom-2 left-0 w-full h-3 text-warm-400 opacity-60" viewBox="0 0 200 9" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M2.00025 6.99997C25.7958 2.00003 103.653 -2.25704 197.904 2.87265" stroke="currentColor" stroke-width="3" stroke-linecap="round"/></svg>
               </span><br/>
               Bersama Mentor Ahli
             </h1>
             
             <p class="text-lg text-neutral-600 mb-8 max-w-xl leading-relaxed">
               Ribuan peserta telah berhasil melalui program Webinar Prakerja kami. 
               Belajar fleksibel, materi berkualitas, dan bersertifikat resmi.
             </p>
             
             <div class="flex flex-col sm:flex-row gap-4 mb-12">
               <NuxtLink to="/register" class="inline-flex items-center justify-center px-8 py-3.5 text-base font-bold text-white bg-warm-500 rounded-xl hover:bg-warm-600 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300">
                 Jelajahi Program
                 <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                   <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
                 </svg>
               </NuxtLink>
               <NuxtLink to="/about" class="inline-flex items-center justify-center px-8 py-3.5 text-base font-semibold text-neutral-700 bg-white border border-neutral-200 rounded-xl hover:bg-neutral-50 hover:border-neutral-300 transition-all">
                 Tentang Kami
               </NuxtLink>
             </div>
             
             <!-- Stats with Modern Icons -->
             <div class="flex items-center space-x-8 pt-8 border-t border-neutral-100">
               <div class="flex items-center gap-3">
                 <div class="w-10 h-10 rounded-full bg-primary-50 flex items-center justify-center text-primary-600">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/></svg>
                 </div>
                 <div>
                   <div class="text-2xl font-bold text-neutral-900 font-display">10K+</div>
                   <div class="text-xs text-neutral-500 font-medium uppercase tracking-wide">Siswa Aktif</div>
                 </div>
               </div>
               <div class="w-px h-10 bg-neutral-200"></div>
               <div class="flex items-center gap-3">
                 <div class="w-10 h-10 rounded-full bg-warm-50 flex items-center justify-center text-warm-600">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                 </div>
                 <div>
                   <div class="text-2xl font-bold text-neutral-900 font-display">95%</div>
                   <div class="text-xs text-neutral-500 font-medium uppercase tracking-wide">Lulus Sertifikasi</div>
                 </div>
               </div>
             </div>
           </div>
          
          <!-- Hero Visual (Image + Floating Card) -->
          <div class="relative animate-slide-up order-1 lg:order-2 flex justify-center lg:justify-end">
            <!-- Hero Image -->
            <div class="relative z-10 w-full max-w-md lg:max-w-full">
              <div class="absolute inset-0 bg-gradient-to-tr from-warm-500/20 to-primary-500/20 rounded-[2rem] transform rotate-3 blur-2xl -z-10"></div>
              <img src="/images/hero.webp" alt="Happy Student" class="w-full h-auto rounded-[2rem] shadow-2xl border-4 border-white relative z-10" />
              
              <!-- Floating Dynamic Card (Course Preview) -->
              <div class="absolute -bottom-12 -left-4 md:-left-12 z-20 w-72 md:w-80">
                <div class="card-hover p-6 relative overflow-hidden bg-white shadow-soft-lg rounded-2xl border border-neutral-100">
                  <div class="absolute top-0 right-0 w-20 h-20 bg-primary-50 rounded-bl-full -mr-10 -mt-10 z-0"></div>
                  
                  <div class="relative z-10">
                    <div class="flex items-center justify-between mb-4">
                      <span class="text-xs font-bold text-primary-600 bg-primary-50 px-2 py-1 rounded-md uppercase tracking-wide">Sedang Populer</span>
                      <div class="flex space-x-1">
                        <span class="w-2 h-2 rounded-full bg-red-400 animate-pulse"></span>
                        <span class="text-[10px] text-neutral-400">Live</span>
                      </div>
                    </div>
                    
                    <div v-if="pendingHero" class="flex justify-center items-center h-32">
                      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
                    </div>
                    <div v-else-if="errorHero" class="text-center text-red-500 text-sm">
                      <p>Gagal memuat.</p>
                    </div>
                    <Transition v-else name="fade" mode="out-in">
                      <div v-if="featuredCourses && featuredCourses.length > 0" :key="currentCourseIndex">
                        <h3 class="font-bold text-neutral-900 text-lg mb-1 leading-snug line-clamp-2">{{ featuredCourses[currentCourseIndex].title }}</h3>
                        <p class="text-xs text-neutral-500 mb-3">{{ featuredCourses[currentCourseIndex].category }}</p>
                        
                        <div class="flex items-center justify-between text-xs text-neutral-600 mb-4">
                          <div class="flex items-center">
                            <svg class="w-4 h-4 text-warning-400 mr-1" fill="currentColor" viewBox="0 0 20 20"><path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/></svg>
                            <span class="font-bold">{{ featuredCourses[currentCourseIndex].rating }}</span>
                          </div>
                          <span>{{ featuredCourses[currentCourseIndex].modules }} Modul</span>
                        </div>
                        
                        <NuxtLink to="/register" class="w-full inline-flex items-center justify-center px-4 py-2 text-sm font-bold text-white bg-primary-600 rounded-lg hover:bg-primary-700 transition-all">
                          Gabung Kelas
                        </NuxtLink>
                      </div>
                    </Transition>
                  </div>
                </div>
              </div>
              
              <!-- Small Badge -->
              <div class="absolute top-8 -right-4 bg-white p-3 rounded-xl shadow-lg z-20 animate-float-delayed flex items-center gap-3">
                 <div class="bg-green-100 p-2 rounded-full text-green-600">
                   <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                 </div>
                 <div>
                   <p class="font-bold text-xs text-gray-800">Terverifikasi</p>
                   <p class="text-[10px] text-gray-500">Resmi Kemnaker</p>
                 </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Popular Courses Section -->
    <section id="courses" class="section-padding bg-white">
      <div class="container-custom">
        <div class="text-center max-w-3xl mx-auto mb-16">
          <span class="badge-primary mb-4">Kursus Populer</span>
          <h2 class="text-3xl md:text-4xl font-bold mb-4">
            Mulai Dari Kursus Terlaris
          </h2>
          <p class="text-lg text-neutral-600">
            Pilih dari ratusan kursus berkualitas tinggi yang dirancang oleh instruktur berpengalaman.
          </p>
        </div>

        <!-- Loading State -->
        <div v-if="loadingCourses" class="grid md:grid-cols-2 xl:grid-cols-3 gap-6">
          <div class="bg-white rounded-xl h-96 animate-pulse bg-neutral-100"></div>
          <div class="bg-white rounded-xl h-96 animate-pulse bg-neutral-100"></div>
          <div class="bg-white rounded-xl h-96 animate-pulse bg-neutral-100"></div>
        </div>

        <!-- Courses Grid -->
        <div v-else-if="popularCourses && popularCourses.length > 0" class="grid md:grid-cols-2 xl:grid-cols-3 gap-6 mb-8">
          <div 
            v-for="course in popularCourses.slice(0, 6)" 
            :key="course.id" 
            class="bg-white rounded-xl border border-neutral-200 overflow-hidden hover:shadow-lg transition-all duration-300 group cursor-pointer"
            @click="$router.push(`/dashboard/courses/${course.id}`)"
          >
            <div class="h-48 relative overflow-hidden bg-gradient-to-br from-primary-50 to-accent-50">
              <img 
                v-if="course.thumbnail_url" 
                :src="getThumbnailUrl(course.thumbnail_url)" 
                :alt="course.title"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
              />
              <div v-else class="w-full h-full flex items-center justify-center">
                <svg class="w-20 h-20 text-neutral-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
              </div>
              
              <div class="absolute top-3 left-3" v-if="course.category">
                <span class="px-3 py-1 bg-black/60 backdrop-blur-sm text-white text-xs font-medium rounded-lg">
                  {{ course.category.name }}
                </span>
              </div>
              <div v-if="course.is_featured" class="absolute top-3 right-3">
                <span class="px-3 py-1 bg-warm-500 text-white text-xs font-medium rounded-lg shadow-lg">
                  ⭐ Featured
                </span>
              </div>
            </div>
            
            <div class="p-5">
              <h3 class="font-semibold text-lg text-neutral-900 mb-2 group-hover:text-primary-600 transition-colors line-clamp-2">
                {{ course.title }}
              </h3>
              <p class="text-sm text-neutral-500 mb-4">
                {{ course.instructor?.full_name || 'Instructor' }}
              </p>
              
              <div class="flex items-center justify-between mb-4">
                <div class="flex items-center space-x-2">
                  <div class="flex items-center">
                    <svg v-for="i in 5" :key="i" class="w-4 h-4" 
                      :class="i <= 4 ? 'text-warm-400 fill-current' : 'text-neutral-200 fill-current'" 
                      viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                    </svg>
                  </div>
                  <span class="text-sm text-neutral-700 font-medium">4.8</span>
                  <span class="text-sm text-neutral-500">(128)</span>
                </div>
                <span class="text-sm text-neutral-500">{{ course.lessons?.length || 0 }} Modul</span>
              </div>
              
              <div class="flex items-center justify-between pt-4 border-t border-neutral-100">
                <div>
                  <span class="text-xl font-bold text-neutral-900">
                    {{ course.price === 0 ? 'Gratis' : new Intl.NumberFormat('id-ID', { style: 'currency', currency: course.currency || 'IDR', minimumFractionDigits: 0 }).format(course.price) }}
                  </span>
                </div>
                <button class="px-4 py-2 text-sm font-semibold text-primary-600 bg-primary-50 rounded-lg hover:bg-primary-100 transition-colors">
                  Lihat Detail
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="text-center py-16">
          <svg class="w-16 h-16 text-neutral-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
          </svg>
          <p class="text-neutral-500">Kursus akan segera tersedia</p>
        </div>

        <!-- View All Button -->
        <div v-if="popularCourses && popularCourses.length > 0" class="text-center mt-8">
          <NuxtLink to="/dashboard/explore" class="inline-flex items-center px-6 py-3 text-sm font-semibold text-primary-600 bg-primary-50 rounded-xl hover:bg-primary-100 transition-colors">
            Lihat Semua Kursus
            <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </NuxtLink>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section class="section-padding bg-gradient-to-b from-neutral-50 to-white">
      <div class="container-custom">
        <div class="text-center max-w-3xl mx-auto mb-16">
          <span class="badge-success mb-4">Cara Kerja</span>
          <h2 class="text-3xl md:text-4xl font-bold mb-4">
            Mulai Belajar dalam 3 Langkah Mudah
          </h2>
          <p class="text-lg text-neutral-600">
            Proses sederhana untuk memulai perjalanan belajar Anda
          </p>
        </div>
        
        <div class="grid md:grid-cols-3 gap-8 lg:gap-12">
          <!-- Step 1 -->
          <div class="text-center relative">
            <div class="inline-flex items-center justify-center w-20 h-20 bg-primary-100 rounded-2xl mb-6 relative z-10">
              <div class="absolute -top-2 -right-2 w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center text-white font-bold text-sm">
                1
              </div>
              <svg class="w-10 h-10 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-3 text-neutral-900">Daftar Akun</h3>
            <p class="text-neutral-600">
              Buat akun gratis dalam hitungan detik. Tidak perlu kartu kredit untuk memulai.
            </p>
          </div>

          <!-- Step 2 -->
          <div class="text-center relative">
            <div class="absolute top-10 left-1/2 w-full h-0.5 bg-gradient-to-r from-success-200 via-success-300 to-transparent hidden md:block"></div>
            <div class="inline-flex items-center justify-center w-20 h-20 bg-accent-100 rounded-2xl mb-6 relative z-10">
              <div class="absolute -top-2 -right-2 w-8 h-8 bg-accent-600 rounded-full flex items-center justify-center text-white font-bold text-sm">
                2
              </div>
              <svg class="w-10 h-10 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-3 text-neutral-900">Pilih Kursus</h3>
            <p class="text-neutral-600">
              Jelajahi katalog kursus kami dan pilih yang sesuai dengan minat dan tujuan Anda.
            </p>
          </div>

          <!-- Step 3 -->
          <div class="text-center relative">
            <div class="absolute top-10 right-1/2 w-full h-0.5 bg-gradient-to-l from-success-200 via-success-300 to-transparent hidden md:block"></div>
            <div class="inline-flex items-center justify-center w-20 h-20 bg-success-100 rounded-2xl mb-6 relative z-10">
              <div class="absolute -top-2 -right-2 w-8 h-8 bg-success-600 rounded-full flex items-center justify-center text-white font-bold text-sm">
                3
              </div>
              <svg class="w-10 h-10 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/>
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-3 text-neutral-900">Mulai Belajar</h3>
            <p class="text-neutral-600">
              Mulai belajar dengan materi berkualitas tinggi dan dapatkan sertifikat setelah selesai.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="section-padding bg-white relative overflow-hidden">
      <!-- Background Texture (Asset 3 Substitution) -->
      <div class="absolute inset-0 bg-gradient-to-b from-white via-primary-50/30 to-white -z-10"></div>
      
      <div class="container-custom">
        <div class="text-center max-w-3xl mx-auto mb-16">
          <span class="px-4 py-1.5 rounded-full bg-accent-50 text-accent-600 font-semibold text-sm mb-6 inline-block tracking-wide uppercase">Fitur Unggulan</span>
          <h2 class="text-3xl md:text-4xl font-display font-bold mb-4 text-neutral-900 leading-tight">
            Semua yang Anda Butuhkan
          </h2>
          <p class="text-lg text-neutral-600 leading-relaxed">
            Platform kami dilengkapi dengan fitur-fitur terbaik untuk memastikan pengalaman belajar Anda efektif dan menyenangkan.
          </p>
        </div>
        
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <!-- Feature Cards with Cleaner Design -->
          <div class="group bg-white rounded-2xl p-8 border border-neutral-100 shadow-soft hover:shadow-soft-lg transition-all duration-300 hover:-translate-y-1">
            <div class="w-14 h-14 bg-primary-50 rounded-xl flex items-center justify-center mb-6 group-hover:bg-primary-100 transition-colors">
              <svg class="w-7 h-7 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-3 text-neutral-900">Video HD Berkualitas</h3>
            <p class="text-neutral-600 mb-6 leading-relaxed">Materi video dengan kualitas tinggi yang dilindungi dan mudah diakses dari perangkat apapun.</p>
            <div class="pt-6 border-t border-neutral-50 group-hover:border-neutral-100 transition-colors">
              <span class="text-primary-600 font-semibold text-sm flex items-center">
                Pelajari Lebih Lanjut
                <svg class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
              </span>
            </div>
          </div>
          
          <div class="group bg-white rounded-2xl p-8 border border-neutral-100 shadow-soft hover:shadow-soft-lg transition-all duration-300 hover:-translate-y-1">
            <div class="w-14 h-14 bg-accent-50 rounded-xl flex items-center justify-center mb-6 group-hover:bg-accent-100 transition-colors">
              <svg class="w-7 h-7 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-3 text-neutral-900">Kuis Interaktif</h3>
            <p class="text-neutral-600 mb-6 leading-relaxed">Uji pemahaman Anda dengan kuis di setiap modul dan dapatkan feedback instan.</p>
            <div class="pt-6 border-t border-neutral-50 group-hover:border-neutral-100 transition-colors">
              <span class="text-accent-600 font-semibold text-sm flex items-center">
                Pelajari Lebih Lanjut
                <svg class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
              </span>
            </div>
          </div>
          
          <div class="group bg-white rounded-2xl p-8 border border-neutral-100 shadow-soft hover:shadow-soft-lg transition-all duration-300 hover:-translate-y-1">
            <div class="w-14 h-14 bg-warm-50 rounded-xl flex items-center justify-center mb-6 group-hover:bg-warm-100 transition-colors">
              <svg class="w-7 h-7 text-warm-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138z"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-3 text-neutral-900">Sertifikat Resmi</h3>
            <p class="text-neutral-600 mb-6 leading-relaxed">Dapatkan sertifikat yang diakui industri setelah menyelesaikan kursus untuk portofolio Anda.</p>
            <div class="pt-6 border-t border-neutral-50 group-hover:border-neutral-100 transition-colors">
              <span class="text-warm-600 font-semibold text-sm flex items-center">
                Pelajari Lebih Lanjut
                <svg class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
              </span>
            </div>
          </div>
          
          <div class="group bg-white rounded-2xl p-8 border border-neutral-100 shadow-soft hover:shadow-soft-lg transition-all duration-300 hover:-translate-y-1">
            <div class="w-14 h-14 bg-rose-50 rounded-xl flex items-center justify-center mb-6 group-hover:bg-rose-100 transition-colors">
              <svg class="w-7 h-7 text-rose-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-3 text-neutral-900">Komunitas Aktif</h3>
            <p class="text-neutral-600 mb-6 leading-relaxed">Bergabung dengan komunitas learner dan diskusi dengan sesama siswa dan mentor.</p>
            <div class="pt-6 border-t border-neutral-50 group-hover:border-neutral-100 transition-colors">
              <span class="text-rose-600 font-semibold text-sm flex items-center">
                Pelajari Lebih Lanjut
                <svg class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
              </span>
            </div>
          </div>
          
          <div class="group bg-white rounded-2xl p-8 border border-neutral-100 shadow-soft hover:shadow-soft-lg transition-all duration-300 hover:-translate-y-1">
            <div class="w-14 h-14 bg-sky-50 rounded-xl flex items-center justify-center mb-6 group-hover:bg-sky-100 transition-colors">
              <svg class="w-7 h-7 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-3 text-neutral-900">Akses Seumur Hidup</h3>
            <p class="text-neutral-600 mb-6 leading-relaxed">Sekali beli, akses selamanya. Belajar kapan saja sesuai jadwal Anda sendiri.</p>
            <div class="pt-6 border-t border-neutral-50 group-hover:border-neutral-100 transition-colors">
              <span class="text-sky-600 font-semibold text-sm flex items-center">
                Pelajari Lebih Lanjut
                <svg class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
              </span>
            </div>
          </div>
          
          <div class="group bg-white rounded-2xl p-8 border border-neutral-100 shadow-soft hover:shadow-soft-lg transition-all duration-300 hover:-translate-y-1">
            <div class="w-14 h-14 bg-cyan-50 rounded-xl flex items-center justify-center mb-6 group-hover:bg-cyan-100 transition-colors">
              <svg class="w-7 h-7 text-cyan-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-3 text-neutral-900">Mobile Friendly</h3>
            <p class="text-neutral-600 mb-6 leading-relaxed">Belajar dari mana saja dengan tampilan optimal di semua perangkat smartphone.</p>
            <div class="pt-6 border-t border-neutral-50 group-hover:border-neutral-100 transition-colors">
              <span class="text-cyan-600 font-semibold text-sm flex items-center">
                Pelajari Lebih Lanjut
                <svg class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Mentor Section (New Humanized Section) -->
    <section class="section-padding bg-neutral-50 overflow-hidden relative">
      <div class="container-custom">
        <div class="grid lg:grid-cols-2 gap-12 lg:gap-20 items-center">
          <!-- Image Content -->
          <div class="relative animate-fade-in order-2 lg:order-1">
            <div class="relative rounded-3xl overflow-hidden shadow-2xl border-4 border-white transform rotate-2 hover:rotate-0 transition-transform duration-500">
               <div class="absolute inset-0 bg-primary-900/10 hover:bg-transparent transition-colors z-10"></div>
               <img src="/images/mentor.webp" alt="Mentoring Session" class="w-full h-auto object-cover transform scale-105 hover:scale-110 transition-transform duration-700" />
            </div>
            
            <!-- Floating Badge -->
            <div class="absolute -bottom-4 -left-4 card p-4 bg-white/90 backdrop-blur shadow-soft rounded-2xl animate-float">
              <div class="flex items-center space-x-2">
                <div class="flex -space-x-3">
                  <img src="/images/avatar-budi.webp" alt="Student" class="w-8 h-8 rounded-full border-2 border-white object-cover" />
                  <img src="/images/avatar-dewi.webp" alt="Student" class="w-8 h-8 rounded-full border-2 border-white object-cover" />
                  <img src="/images/avatar-ahmad.webp" alt="Student" class="w-8 h-8 rounded-full border-2 border-white object-cover" />
                </div>
                <span class="text-sm font-medium text-neutral-700 ml-2">+2.5K bergabung</span>
              </div>
            </div>
            
            <!-- Pattern -->
            <div class="absolute -top-12 -left-12 w-64 h-64 bg-accent-100 rounded-full blur-3xl -z-10 opacity-60"></div>
          </div>
          
          <!-- Text Content -->
          <div class="order-1 lg:order-2">
            <div class="inline-flex items-center px-4 py-2 bg-warm-50 rounded-full mb-6 animate-fade-in border border-warm-100">
               <span class="w-2 h-2 bg-warm-500 rounded-full mr-2"></span>
               <span class="text-sm font-medium text-warm-700">Bimbingan Eksklusif</span>
             </div>
            <h2 class="text-3xl md:text-4xl font-display font-bold mb-6 text-neutral-900 leading-tight">
              Belajar Tidak Sendiri, <br/>
              <span class="text-primary-600">Ada Mentor yang Menemani</span>
            </h2>
            <p class="text-lg text-neutral-600 mb-8 leading-relaxed">
              Kami memadukan teknologi dengan sentuhan manusia. Anda tidak hanya menonton video, tapi juga bisa berdiskusi, bertanya, dan mendapatkan feedback langsung dari praktisi ahli.
            </p>
            
            <ul class="space-y-4 mb-8">
              <li class="flex items-start">
                <div class="flex-shrink-0 w-6 h-6 rounded-full bg-success-100 flex items-center justify-center mt-1 mr-4">
                  <svg class="w-4 h-4 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                </div>
                <div>
                  <h4 class="font-bold text-neutral-900">Sesi Tanya Jawab Langsung</h4>
                  <p class="text-sm text-neutral-500">Jadwal rutin webinar dan diskusi grup.</p>
                </div>
              </li>
              <li class="flex items-start">
                <div class="flex-shrink-0 w-6 h-6 rounded-full bg-success-100 flex items-center justify-center mt-1 mr-4">
                  <svg class="w-4 h-4 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                </div>
                <div>
                  <h4 class="font-bold text-neutral-900">Review Portofolio</h4>
                  <p class="text-sm text-neutral-500">Dapatkan masukan konstruktif untuk karya Anda.</p>
                </div>
              </li>
              <li class="flex items-start">
                <div class="flex-shrink-0 w-6 h-6 rounded-full bg-success-100 flex items-center justify-center mt-1 mr-4">
                  <svg class="w-4 h-4 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                </div>
                <div>
                  <h4 class="font-bold text-neutral-900">Komunitas Suportif</h4>
                  <p class="text-sm text-neutral-500">Jaringan dengan ribuan pembelajar lainnya.</p>
                </div>
              </li>
            </ul>
            
            <div class="inline-flex items-center text-primary-600 font-bold hover:text-primary-700 transition-colors group cursor-pointer">
              Kenalan dengan Mentor
              <svg class="w-5 h-5 ml-2 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"/></svg>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Testimonials Section -->
    <section id="testimonials" class="section-padding bg-neutral-50">
      <div class="container-custom">
        <div class="text-center max-w-3xl mx-auto mb-16">
          <span class="px-4 py-1.5 rounded-full bg-warm-50 text-warm-600 font-semibold text-sm mb-6 inline-block tracking-wide uppercase">Kata Alumni</span>
          <h2 class="text-3xl md:text-4xl font-display font-bold mb-4 text-neutral-900 leading-tight">Apa Kata Mereka?</h2>
          <p class="text-lg text-neutral-600 leading-relaxed">Ribuan siswa telah merasakan manfaat belajar di platform kami.</p>
        </div>
        
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <!-- Testimonial 1 -->
          <div class="bg-white p-8 rounded-2xl shadow-soft hover:shadow-soft-lg transition-all duration-300 border border-neutral-100 flex flex-col">
            <div class="flex items-center space-x-1 mb-6">
              <template v-for="i in 5" :key="i">
                <svg class="w-5 h-5 text-warm-400" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                </svg>
              </template>
            </div>
            <p class="text-neutral-700 mb-8 leading-relaxed italic flex-grow">
              "Platform ini benar-benar mengubah cara saya belajar. Materi yang terstruktur dan kualitas video yang sangat baik membuat proses belajar menjadi lebih menyenangkan!"
            </p>
            <div class="flex items-center space-x-4 pt-6 border-t border-neutral-50">
              <img src="/images/avatar-budi.webp" alt="Budi Santoso" class="w-12 h-12 rounded-full object-cover ring-2 ring-primary-50" />
              <div>
                <p class="font-bold text-neutral-900 text-sm">B Pandega</p>
                <p class="text-xs text-neutral-500 font-medium">Software Engineer</p>
              </div>
            </div>
          </div>
          
          <!-- Testimonial 2 -->
          <div class="bg-white p-8 rounded-2xl shadow-soft hover:shadow-soft-lg transition-all duration-300 border border-neutral-100 flex flex-col">
            <div class="flex items-center space-x-1 mb-6">
              <template v-for="i in 5" :key="i">
                <svg class="w-5 h-5 text-warm-400" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                </svg>
              </template>
            </div>
            <p class="text-neutral-700 mb-8 leading-relaxed italic flex-grow">
              "Sertifikatnya diakui perusahaan tempat saya bekerja. Investasi terbaik untuk karir saya! Dalam 6 bulan saya berhasil naik jabatan."
            </p>
            <div class="flex items-center space-x-4 pt-6 border-t border-neutral-50">
              <img src="/images/avatar-dewi.webp" alt="Dewi Sartika" class="w-12 h-12 rounded-full object-cover ring-2 ring-primary-50" />
              <div>
                <p class="font-bold text-neutral-900 text-sm">Indra Hura</p>
                <p class="text-xs text-neutral-500 font-medium">Data Analyst</p>
              </div>
            </div>
          </div>
          
          <!-- Testimonial 3 -->
          <div class="bg-white p-8 rounded-2xl shadow-soft hover:shadow-soft-lg transition-all duration-300 border border-neutral-100 flex flex-col">
            <div class="flex items-center space-x-1 mb-6">
              <template v-for="i in 5" :key="i">
                <svg class="w-5 h-5 text-warm-400" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                </svg>
              </template>
            </div>
            <p class="text-neutral-700 mb-8 leading-relaxed italic flex-grow">
              "Dari nol sampai bisa landing job sebagai developer. Semua berkat kursus di platform ini. Sekarang saya bekerja di startup teknologi terkemuka!"
            </p>
            <div class="flex items-center space-x-4 pt-6 border-t border-neutral-50">
              <img src="/images/avatar-ahmad.webp" alt="Ahmad Dhani" class="w-12 h-12 rounded-full object-cover ring-2 ring-primary-50" />
              <div>
                <p class="font-bold text-neutral-900 text-sm">Sugeng Widodo</p>
                <p class="text-xs text-neutral-500 font-medium">Full Stack Developer</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="section-padding bg-primary-600 relative overflow-hidden">
      <!-- Background Pattern -->
      <div class="absolute inset-0 opacity-10">
        <div class="absolute top-20 left-20 w-64 h-64 bg-white rounded-full blur-3xl"></div>
        <div class="absolute bottom-20 right-20 w-96 h-96 bg-white rounded-full blur-3xl"></div>
      </div>
      
      <div class="container-custom text-center relative z-10">
        <h2 class="text-3xl md:text-4xl lg:text-5xl font-display font-bold text-white mb-6">
          Siap Memulai Perjalanan Belajar Anda?
        </h2>
        <p class="text-lg md:text-xl text-primary-100 mb-10 max-w-2xl mx-auto leading-relaxed">
          Bergabung dengan <span class="font-semibold text-white">10,000+ siswa</span> yang sudah mengembangkan skill mereka dan mencapai tujuan karir bersama kami.
        </p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
          <NuxtLink to="/register" class="inline-flex items-center justify-center px-8 py-4 text-base font-bold text-white bg-warm-500 rounded-xl hover:bg-warm-400 hover:scale-105 transition-all shadow-xl shadow-warm-900/20">
            Daftar Sekarang - Gratis
            <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </NuxtLink>
          <NuxtLink to="/dashboard/explore" class="inline-flex items-center justify-center px-8 py-4 text-base font-semibold text-white border-2 border-white/30 rounded-xl hover:bg-white/10 hover:border-white/50 transition-all backdrop-blur-sm">
            Lihat Semua Kursus
          </NuxtLink>
        </div>
        
        <!-- Trust Indicators -->
        <div class="flex flex-wrap items-center justify-center gap-8 mt-12 pt-8 border-t border-white/10">
          <div class="flex items-center space-x-2 text-white/90">
            <svg class="w-5 h-5 text-warm-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
            </svg>
            <span class="text-sm font-medium">Tidak Perlu Kartu Kredit</span>
          </div>
          <div class="flex items-center space-x-2 text-white/90">
            <svg class="w-5 h-5 text-warm-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
            </svg>
            <span class="text-sm font-medium">Akses Selamanya</span>
          </div>
          <div class="flex items-center space-x-2 text-white/90">
            <svg class="w-5 h-5 text-warm-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
            </svg>
            <span class="text-sm font-medium">Sertifikat Diakui</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-neutral-900 text-neutral-400 py-16">
      <div class="container-custom">
        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-12 mb-12">
          <div class="lg:col-span-1">
            <div class="flex items-center space-x-2 mb-4">
              <img src="/logo.png" alt="EDUKRA Logo" class="w-10 h-10 object-contain" />
              <span class="font-display font-bold text-xl text-white">EDUKRA</span>
            </div>
            <p class="text-sm leading-relaxed">EDUKRA adalah pusat energi pendidikan - platform LMS modern untuk mengembangkan skill dan mencapai tujuan karir Anda.</p>
          </div>
          
          <div>
            <h4 class="font-semibold text-white mb-4">Quick Links</h4>
            <ul class="space-y-3 text-sm">
              <li><a href="#" class="hover:text-white transition-colors">Tentang Kami</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Kursus</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Instruktur</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Kontak</a></li>
            </ul>
          </div>
          
          <div>
            <h4 class="font-semibold text-white mb-4">Support</h4>
            <ul class="space-y-3 text-sm">
              <li><a href="#" class="hover:text-white transition-colors">FAQ</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Help Center</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Privacy Policy</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Terms of Service</a></li>
            </ul>
          </div>
          
          <div>
            <h4 class="font-semibold text-white mb-4">Hubungi Kami</h4>
            <ul class="space-y-3 text-sm">
              <li>edukra.id@gmail.com</li>
              <li>+62 812 3456 7890</li>
              <li>Jakarta, Indonesia</li>
            </ul>
          </div>
        </div>
        
        <div class="pt-8 border-t border-neutral-800 text-center text-sm">
          <p>&copy; 2024 EDUKRA. All rights reserved.</p>
        </div>
      </div>
    </footer>

    <!-- WhatsApp Floating Button -->
    <a 
      href="https://wa.me/6281381936988?text=Halo%20Admin%20Edukra%2C%20saya%20tertarik%20untuk%20mengetahui%20lebih%20lanjut%20tentang%20kursus%20yang%20tersedia.%20Bisa%20tolong%20dibantu%3F" 
      target="_blank" 
      rel="noopener noreferrer" 
      class="fixed bottom-6 right-6 z-50 flex items-center justify-center w-14 h-14 bg-[#25D366] text-white rounded-full shadow-lg hover:bg-[#20bd5a] hover:scale-110 transition-all duration-300 group"
      aria-label="Chat via WhatsApp"
    >
      <!-- Tooltip -->
      <span class="absolute right-16 top-1/2 -translate-y-1/2 px-3 py-2 bg-white text-neutral-800 text-sm font-medium rounded-lg shadow-lg opacity-0 group-hover:opacity-100 transition-opacity duration-300 whitespace-nowrap pointer-events-none">
        Tanya Admin
      </span>
      <svg class="w-8 h-8" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/>
      </svg>
    </a>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, h, computed, watchEffect } from 'vue'

useHead({
  title: 'EDUKRA - The Engine of Learning | Platform LMS Modern Indonesia',
  meta: [
    { name: 'description', content: 'EDUKRA adalah pusat energi pendidikan - platform LMS modern untuk mengembangkan skill dan mencapai tujuan karir. Akses ribuan kursus berkualitas dari instruktur terbaik.' }
  ]
})

const config = useRuntimeConfig()
const apiBase = config.public.apiBase || 'http://localhost:8080'

// Get proper thumbnail URL - handle MinIO objects
const getThumbnailUrl = (url: string | null | undefined): string => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads')) return `${apiBase}${url}`
  return `${apiBase}/api/images/${url}`
}

// Featured Courses Carousel (for hero section)
const currentCourseIndex = ref(0)
let carouselInterval: any = null

// Icon Styles for cycling
const iconStyles = [
  {
    bg: 'bg-secondary-100',
    color: 'text-secondary-500',
    path: 'M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4'
  },
  {
    bg: 'bg-accent-100',
    color: 'text-accent-500',
    path: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 002 2h2a2 2 0 002-2z'
  },
  {
    bg: 'bg-primary-100',
    color: 'text-primary-500',
    path: 'M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z'
  },
  {
    bg: 'bg-rose-100',
    color: 'text-rose-500',
    path: 'M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01'
  }
]

// Fetch Featured Courses for Hero Section (Client-side only to avoid SSR issues)
const { data: coursesData, pending: pendingHero, error: errorHero } = useFetch('/api/courses', {
  baseURL: config.public.apiBase,
  query: { limit: 5 },
  server: false,
  lazy: true
})

// Fetch Popular Courses for Courses Section
const { data: popularCoursesData, pending: loadingCourses, error: errorCourses } = useFetch('/api/courses', {
  baseURL: config.public.apiBase,
  query: { limit: 6 },
  server: false,
  lazy: true
})

// Popular courses for the courses section
const popularCourses = computed(() => {
  return popularCoursesData.value?.courses || []
})

// Map API data to UI format for hero featured courses
const featuredCourses = computed(() => {
  const courses = coursesData.value?.courses || []
  
  // If no courses found, return empty array (or could return mock fallback)
  if (courses.length === 0) return []

  return courses.map((course: any, index: number) => {
    const style = iconStyles[index % iconStyles.length]
    return {
      title: course.title,
      modules: course.lessons_count || 10, // Default if 0
      duration: course.duration || '10 Jam', // Default if empty
      level: 'Semua Level', // Default
      levelColor: style.color, // Use theme color
      rating: '5.0', // Default
      iconBg: style.bg,
      iconColor: style.color,
      icon: h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: style.path })
      ])
    }
  })
})

onMounted(() => {
  carouselInterval = setInterval(() => {
    if (featuredCourses.value.length > 0) {
      currentCourseIndex.value = (currentCourseIndex.value + 1) % featuredCourses.value.length
    }
  }, 4000)
})

onBeforeUnmount(() => {
  if (carouselInterval) clearInterval(carouselInterval)
})
</script>

<style scoped>
.animate-float-cube {
  animation: float-rotate-12 8s ease-in-out infinite;
}
.animate-float-sphere {
  animation: float-simple 6s ease-in-out infinite;
}
.animate-float-prism {
  animation: float-rotate-minus-12 7s ease-in-out infinite;
}
.animation-delay-1000 {
  animation-delay: 1s;
}
.animation-delay-2000 {
  animation-delay: 2s;
}

@keyframes float-rotate-12 {
  0%, 100% { transform: translateY(0) rotate(12deg); }
  50% { transform: translateY(-20px) rotate(17deg); }
}

@keyframes float-rotate-minus-12 {
  0%, 100% { transform: translateY(0) rotate(-12deg); }
  50% { transform: translateY(-15px) rotate(-7deg); }
}

@keyframes float-simple {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-18px); }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease, transform 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
