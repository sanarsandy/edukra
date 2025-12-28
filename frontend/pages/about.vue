<template>
  <div class="min-h-screen bg-neutral-50 font-sans text-neutral-900">
    <Navbar />

    <!-- Hero Section -->
    <section class="pt-32 pb-16 bg-white relative overflow-hidden">
      <div class="container-custom relative z-10">
        <div class="max-w-3xl mx-auto text-center">
          <div class="inline-flex items-center px-4 py-2 bg-primary-50 rounded-full mb-6 border border-primary-100">
             <span class="w-2 h-2 bg-warm-500 rounded-full mr-2"></span>
             <span class="text-sm font-medium text-primary-700">Tentang Kami</span>
           </div>
          <h1 class="text-4xl md:text-5xl font-display font-bold mb-6 text-neutral-900">
            Membangun Masa Depan Melalui Pendidikan Digital
          </h1>
          <p class="text-lg text-neutral-600 leading-relaxed">
            EDUKRA hadir untuk menjembatani kesenjangan skill di Indonesia dengan menyediakan akses pendidikan berkualitas tinggi yang terjangkau, fleksibel, dan relevan dengan kebutuhan industri.
          </p>
        </div>
      </div>
      
      <!-- Background Elements -->
      <div class="absolute top-0 left-1/2 -translate-x-1/2 w-full h-full max-w-7xl pointer-events-none">
        <div class="absolute top-1/4 left-0 w-64 h-64 bg-primary-100 rounded-full blur-3xl opacity-40"></div>
        <div class="absolute bottom-0 right-0 w-96 h-96 bg-warm-100 rounded-full blur-3xl opacity-40"></div>
      </div>
    </section>

    <!-- Vision & Mission -->
    <section class="py-16 bg-white border-t border-neutral-100">
      <div class="container-custom">
        <div class="grid md:grid-cols-2 gap-12 items-center">
          <div>
            <img src="/images/mentor.webp" alt="Our Vision" class="rounded-3xl shadow-soft-lg object-cover h-96 w-full" />
          </div>
          <div>
            <h2 class="text-3xl font-display font-bold mb-6">Visi Kami</h2>
            <p class="text-neutral-600 mb-8 leading-relaxed">
              Menjadi platform pendidikan teknologi terkemuka di Asia Tenggara yang mencetak 1 juta talenta digital siap kerja pada tahun 2030.
            </p>
            
            <h2 class="text-3xl font-display font-bold mb-6">Misi Kami</h2>
            <ul class="space-y-4">
              <li class="flex items-start">
                <div class="w-8 h-8 rounded-full bg-success-100 flex items-center justify-center mr-4 flex-shrink-0 mt-1">
                  <svg class="w-4 h-4 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                </div>
                <p class="text-neutral-600">Menyediakan kurikulum standar industri yang terus diperbarui.</p>
              </li>
              <li class="flex items-start">
                <div class="w-8 h-8 rounded-full bg-success-100 flex items-center justify-center mr-4 flex-shrink-0 mt-1">
                 <svg class="w-4 h-4 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                </div>
                <p class="text-neutral-600">Memberikan pendampingan mentor (mentorship) yang personal.</p>
              </li>
              <li class="flex items-start">
                <div class="w-8 h-8 rounded-full bg-success-100 flex items-center justify-center mr-4 flex-shrink-0 mt-1">
                  <svg class="w-4 h-4 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                </div>
                <p class="text-neutral-600">Membangun komunitas belajar yang suportif dan kolaboratif.</p>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- Team / Founders -->
    <section class="py-20 bg-neutral-50">
      <div class="container-custom">
        <div class="text-center mb-16">
          <h2 class="text-3xl font-display font-bold mb-4">Tim Kami</h2>
          <p class="text-neutral-600">Orang-orang di balik EDUKRA.</p>
        </div>
        
        <div class="grid md:grid-cols-3 gap-8">
          <div class="bg-white p-6 rounded-2xl shadow-soft text-center group hover:-translate-y-1 transition-transform">
            <div class="w-24 h-24 mx-auto mb-4 relative">
              <div class="absolute inset-0 bg-primary-100 rounded-full scale-110 group-hover:scale-125 transition-transform duration-300"></div>
              <img src="/images/avatar-budi.webp" alt="Founder" class="w-full h-full object-cover rounded-full relative z-10 border-2 border-white" />
            </div>
            <h3 class="font-bold text-lg text-neutral-900">Indra Pratama</h3>
            <p class="text-primary-600 text-sm font-medium">CEO</p>
            <p class="text-neutral-500 text-sm mt-3">Ex-Data Analyst Kemenkeu dengan passion di pendidikan.</p>
          </div>
          
           <div class="bg-white p-6 rounded-2xl shadow-soft text-center group hover:-translate-y-1 transition-transform">
            <div class="w-24 h-24 mx-auto mb-4 relative">
              <div class="absolute inset-0 bg-accent-100 rounded-full scale-110 group-hover:scale-125 transition-transform duration-300"></div>
              <img src="/images/avatar-dewi.webp" alt="Co-Founder" class="w-full h-full object-cover rounded-full relative z-10 border-2 border-white" />
            </div>
            <h3 class="font-bold text-lg text-neutral-900">Gede Agus Adi Wijaya</h3>
            <p class="text-accent-600 text-sm font-medium">COO</p>
            <p class="text-neutral-500 text-sm mt-3">10+ tahun pengalaman IT.</p>
          </div>

           <div class="bg-white p-6 rounded-2xl shadow-soft text-center group hover:-translate-y-1 transition-transform">
            <div class="w-24 h-24 mx-auto mb-4 relative">
              <div class="absolute inset-0 bg-warm-100 rounded-full scale-110 group-hover:scale-125 transition-transform duration-300"></div>
              <img src="/images/avatar-ahmad.webp" alt="CTO" class="w-full h-full object-cover rounded-full relative z-10 border-2 border-white" />
            </div>
            <h3 class="font-bold text-lg text-neutral-900">Ihsan Arsandy</h3>
            <p class="text-warm-600 text-sm font-medium">CTO</p>
            <p class="text-neutral-500 text-sm mt-3">Tech architect di balik platform handal kami.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer Simple -->
    <footer class="bg-white py-8 border-t border-neutral-200 text-center">
      <div class="container-custom">
        <p class="text-neutral-500 text-sm">© 2024 EDUKRA. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
useHead({
  title: 'Tentang Kami - EDUKRA',
  meta: [
    { name: 'description', content: 'Mengenal lebih dekat visi dan misi EDUKRA dalam memajukan pendidikan teknologi di Indonesia.' }
  ]
})
</script>
