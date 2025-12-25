<template>
  <div class="min-h-screen bg-neutral-50">
    <!-- Sidebar -->
    <aside 
      class="fixed top-0 left-0 z-40 h-screen bg-white border-r border-neutral-200 hidden lg:block transition-all duration-300 overflow-hidden"
      :class="sidebarCollapsed ? 'w-[72px]' : 'w-72'"
    >
      <div class="h-full flex flex-col">
        <!-- Logo -->
        <div class="h-16 flex items-center border-b border-neutral-100" :class="sidebarCollapsed ? 'justify-center' : 'justify-between px-4'">
          <div class="flex items-center gap-3 min-w-0" :class="sidebarCollapsed ? 'justify-center' : ''">
<img src="/logo.png" alt="EDUKRA Logo" class="w-10 h-10 object-contain" />
            <span v-if="!sidebarCollapsed" class="font-display font-bold text-xl text-neutral-900 whitespace-nowrap">EDUKRA</span>
          </div>
          <button v-if="!sidebarCollapsed" @click="toggleSidebar" class="p-2 text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors flex-shrink-0" title="Tutup Sidebar">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 19l-7-7 7-7m8 14l-7-7 7-7"/></svg>
          </button>
        </div>
        
        <!-- Navigation -->
        <nav class="flex-1 py-6 overflow-y-auto overflow-x-hidden" :class="sidebarCollapsed ? 'px-3' : 'px-4'">
          <div class="mb-8">
            <p v-if="!sidebarCollapsed" class="px-3 mb-3 text-xs font-semibold text-neutral-400 uppercase tracking-wider">Menu Utama</p>
            <div class="space-y-1">
              <NuxtLink 
                to="/dashboard" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/dashboard', true) ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100 hover:text-neutral-900',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
                </svg>
                <span v-if="!sidebarCollapsed">Dashboard</span>
                <!-- Tooltip for collapsed state -->
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-neutral-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Dashboard
                </div>
              </NuxtLink>
              
              <NuxtLink 
                to="/dashboard/courses" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/dashboard/courses') ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100 hover:text-neutral-900',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
                <span v-if="!sidebarCollapsed">Kursus Saya</span>
                <span v-if="!sidebarCollapsed" class="ml-auto bg-primary-100 text-primary-700 text-xs font-medium px-2 py-0.5 rounded-full">5</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-neutral-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Kursus Saya
                </div>
              </NuxtLink>
              
              <NuxtLink 
                to="/dashboard/explore" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/dashboard/explore') ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100 hover:text-neutral-900',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Jelajahi Kursus</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-neutral-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Jelajahi Kursus
                </div>
              </NuxtLink>
            </div>
          </div>
          
          <div class="mb-8">
            <p v-if="!sidebarCollapsed" class="px-3 mb-3 text-xs font-semibold text-neutral-400 uppercase tracking-wider">Pencapaian</p>
            <div class="space-y-1">
              <NuxtLink 
                to="/dashboard/certificates" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/dashboard/certificates') ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100 hover:text-neutral-900',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Sertifikat</span>
                <span v-if="!sidebarCollapsed" class="ml-auto bg-accent-100 text-accent-700 text-xs font-medium px-2 py-0.5 rounded-full">2</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-neutral-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Sertifikat
                </div>
              </NuxtLink>

              <NuxtLink 
                to="/dashboard/progress" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/dashboard/progress') ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100 hover:text-neutral-900',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Progress</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-neutral-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Progress
                </div>
              </NuxtLink>
            </div>
          </div>
          
          <div>
            <p v-if="!sidebarCollapsed" class="px-3 mb-3 text-xs font-semibold text-neutral-400 uppercase tracking-wider">Pengaturan</p>
            <div class="space-y-1">
              <NuxtLink 
                to="/dashboard/settings" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/dashboard/settings') ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100 hover:text-neutral-900',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Pengaturan</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-neutral-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Pengaturan
                </div>
              </NuxtLink>
            </div>
          </div>
        </nav>
        
        <!-- Toggle Button (collapsed state) -->
        <div v-if="sidebarCollapsed" class="p-3 border-t border-neutral-100">
          <button @click="toggleSidebar" class="w-full flex items-center justify-center p-2.5 text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors" title="Buka Sidebar">
            <svg class="w-5 h-5 rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 19l-7-7 7-7m8 14l-7-7 7-7"/></svg>
          </button>
        </div>
      </div>
    </aside>
    
    <!-- Mobile Header -->
    <header class="lg:hidden fixed top-0 left-0 right-0 z-50 h-16 bg-white border-b border-neutral-200 flex items-center justify-between px-4">
      <div class="flex items-center space-x-3">
        <button @click="mobileMenuOpen = true" class="p-2 -ml-2 text-neutral-600 hover:text-neutral-900">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 6h16M4 12h16M4 18h16"/>
          </svg>
        </button>
        <img src="/logo.png" alt="EDUKRA Logo" class="w-8 h-8 object-contain" />
      </div>
      <button class="p-2 text-neutral-500 hover:text-neutral-700">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
        </svg>
      </button>
    </header>
    
    <!-- Main Content -->
    <div class="transition-all duration-300" :class="sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-72'">
      <!-- Top Bar -->
      <header class="hidden lg:flex h-16 bg-white border-b border-neutral-200 items-center justify-between px-8">
        <div class="flex items-center flex-1">
          <div class="relative w-full max-w-md">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
            <input type="text" placeholder="Cari kursus, modul, atau topik..." class="w-full pl-10 pr-4 py-2.5 bg-neutral-100 border-0 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:bg-white text-sm placeholder-neutral-500 transition-all"/>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <!-- Dark Mode Toggle (disabled) -->
          <!-- <DarkModeToggle /> -->
          <!-- Notification Dropdown -->
          <NotificationDropdown />
          <!-- User Profile Dropdown -->
          <div class="relative" ref="userDropdownRef">
            <button @click="userDropdownOpen = !userDropdownOpen" class="flex items-center gap-3 p-2 hover:bg-neutral-100 rounded-lg transition-colors">
              <div class="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center text-white font-semibold text-sm">{{ userInitial }}</div>
              <div class="text-left hidden sm:block"><p class="text-sm font-medium text-neutral-900">{{ userName }}</p></div>
              <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 9l-7 7-7-7"/></svg>
            </button>
            <Transition name="dropdown">
              <div v-if="userDropdownOpen" class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-neutral-200 py-2 z-50">
                <NuxtLink to="/dashboard/settings" class="flex items-center px-4 py-2 text-sm text-neutral-700 hover:bg-neutral-100">Profil Saya</NuxtLink>
                <NuxtLink to="/dashboard/settings" class="flex items-center px-4 py-2 text-sm text-neutral-700 hover:bg-neutral-100">Pengaturan</NuxtLink>
                <div class="border-t border-neutral-100 my-2"></div>
                <button @click="handleLogout" class="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50">Keluar</button>
              </div>
            </Transition>
          </div>
        </div>
      </header>
      
      <!-- Page Content -->
      <main class="p-6 lg:p-8 pt-20 lg:pt-8">
        <slot />
      </main>
    </div>
    
    <!-- Mobile Menu Overlay -->
    <Transition name="fade">
      <div v-if="mobileMenuOpen" class="fixed inset-0 z-50 bg-black/50 lg:hidden" @click="mobileMenuOpen = false"></div>
    </Transition>
    
    <!-- Mobile Menu -->
    <Transition name="slide">
      <aside v-if="mobileMenuOpen" class="fixed top-0 left-0 z-50 w-72 h-screen bg-white lg:hidden">
        <div class="h-full flex flex-col">
          <div class="h-16 flex items-center justify-between px-6 border-b border-neutral-100">
            <div class="flex items-center space-x-3">
              <img src="/logo.png" alt="EDUKRA Logo" class="w-10 h-10 object-contain" />
              <span class="font-display font-bold text-xl text-neutral-900">EDUKRA</span>
            </div>
            <button @click="mobileMenuOpen = false" class="p-2 text-neutral-400 hover:text-neutral-600">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <nav class="flex-1 px-4 py-6 overflow-y-auto">
            <div class="space-y-1">
              <NuxtLink to="/dashboard" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/dashboard', true) ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100'">
                <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
                </svg>
                Dashboard
              </NuxtLink>
              <NuxtLink to="/dashboard/courses" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/dashboard/courses') ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100'">
                <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
                Kursus Saya
              </NuxtLink>
              <NuxtLink to="/dashboard/explore" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/dashboard/explore') ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100'">
                <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
                Jelajahi
              </NuxtLink>
              <NuxtLink to="/dashboard/certificates" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/dashboard/certificates') ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100'">
                <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/></svg>
                Sertifikat
              </NuxtLink>
              <NuxtLink to="/dashboard/progress" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/dashboard/progress') ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100'">
                <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>
                Progress
              </NuxtLink>
              <NuxtLink to="/dashboard/settings" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/dashboard/settings') ? 'bg-primary-50 text-primary-700' : 'text-neutral-600 hover:bg-neutral-100'">
                <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                </svg>
                Pengaturan
              </NuxtLink>
            </div>
          </nav>
          <!-- Mobile User Profile -->
          <div class="p-4 border-t border-neutral-100">
            <div class="flex items-center space-x-3 p-3 rounded-xl bg-neutral-50">
              <div class="w-10 h-10 bg-primary-600 rounded-full flex items-center justify-center text-white font-semibold text-sm">{{ userInitial }}</div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-neutral-900 truncate">{{ userName }}</p>
                <p class="text-xs text-neutral-500 truncate">{{ userEmail }}</p>
              </div>
            </div>
            <button @click="handleLogout" class="w-full mt-3 flex items-center justify-center px-4 py-2.5 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors">
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/></svg>
              Keluar
            </button>
          </div>
        </div>
      </aside>
    </Transition>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const token = useCookie('token')
const userCookie = useCookie('user')

const mobileMenuOpen = ref(false)
const sidebarCollapsed = ref(false)
const userDropdownOpen = ref(false)
const userDropdownRef = ref<HTMLElement | null>(null)

const isActive = (path: string, exact = false) => {
  if (exact) return route.path === path
  return route.path === path || route.path.startsWith(path + '/')
}

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
  if (process.client) {
    localStorage.setItem('sidebar-collapsed', String(sidebarCollapsed.value))
  }
}

onMounted(() => {
  if (process.client) {
    const saved = localStorage.getItem('sidebar-collapsed')
    if (saved === 'true') {
      sidebarCollapsed.value = true
    }
    document.addEventListener('click', (e) => {
      if (userDropdownRef.value && !userDropdownRef.value.contains(e.target as Node)) {
        userDropdownOpen.value = false
      }
    })
  }
})

// Computed user info from cookie
const parsedUser = computed(() => {
  if (!userCookie.value) return null
  if (typeof userCookie.value === 'string') {
    try {
      return JSON.parse(userCookie.value)
    } catch { return null }
  }
  return userCookie.value
})

const userName = computed(() => parsedUser.value?.full_name || 'User')
const userEmail = computed(() => parsedUser.value?.email || '')
const userInitial = computed(() => userName.value.charAt(0).toUpperCase())

const handleLogout = () => {
  // Clear cookies
  token.value = null
  userCookie.value = null
  
  // Close dropdown
  userDropdownOpen.value = false
  
  // Redirect to login
  navigateTo('/login')
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
.slide-enter-active, .slide-leave-active {
  transition: transform 0.3s ease;
}
.slide-enter-from, .slide-leave-to {
  transform: translateX(-100%);
}
.fade-text-enter-active, .fade-text-leave-active {
  transition: opacity 0.15s ease;
}
.fade-text-enter-from, .fade-text-leave-to {
  opacity: 0;
}
.dropdown-enter-active, .dropdown-leave-active {
  transition: all 0.2s ease;
}
.dropdown-enter-from, .dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
