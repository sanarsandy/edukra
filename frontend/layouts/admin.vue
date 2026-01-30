<template>
  <div class="min-h-screen bg-neutral-100">
    <!-- Sidebar - Admin Theme (Dark) -->
    <aside 
      class="fixed top-0 left-0 z-40 h-screen bg-neutral-900 hidden lg:block transition-all duration-300 overflow-hidden"
      :class="sidebarCollapsed ? 'w-[72px]' : 'w-72'"
    >
      <div class="h-full flex flex-col">
        <!-- Logo -->
        <div class="h-16 flex items-center border-b border-neutral-800" :class="sidebarCollapsed ? 'justify-center' : 'justify-between px-4'">
          <div class="flex items-center gap-3 min-w-0" :class="sidebarCollapsed ? 'justify-center' : ''">
            <img src="/logo.png" alt="EDUKRA Logo" class="w-10 h-10 object-contain" />
            <div v-if="!sidebarCollapsed"><span class="font-display font-bold text-lg text-white">EDUKRA</span><span class="block text-xs text-admin-400 -mt-0.5">Admin Panel</span></div>
          </div>
          <button v-if="!sidebarCollapsed" @click="toggleSidebar" class="p-2 text-neutral-500 hover:text-white hover:bg-neutral-800 rounded-lg transition-colors flex-shrink-0" title="Tutup Sidebar">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 19l-7-7 7-7m8 14l-7-7 7-7"/></svg>
          </button>
        </div>
        
        <!-- Navigation -->
        <nav class="flex-1 py-6 overflow-y-auto overflow-x-hidden" :class="sidebarCollapsed ? 'px-3' : 'px-4'">
          <div class="mb-8">
            <p v-if="!sidebarCollapsed" class="px-3 mb-3 text-xs font-semibold text-neutral-500 uppercase tracking-wider">Overview</p>
            <div class="space-y-1">
              <NuxtLink 
                to="/admin" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin', true) ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Dashboard</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Dashboard
                </div>
              </NuxtLink>
            </div>
          </div>
          
          <div class="mb-8">
            <p v-if="!sidebarCollapsed" class="px-3 mb-3 text-xs font-semibold text-neutral-500 uppercase tracking-wider">Manajemen</p>
            <div class="space-y-1">
              <NuxtLink 
                to="/admin/users" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/users') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Pengguna</span>
                <span v-if="!sidebarCollapsed && stats.total_users > 0" class="ml-auto bg-neutral-700 text-neutral-300 text-xs font-medium px-2 py-0.5 rounded-full">{{ formatCount(stats.total_users) }}</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Pengguna
                </div>
              </NuxtLink>
              
              <NuxtLink 
                to="/admin/courses" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/courses') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
                <span v-if="!sidebarCollapsed">Kursus</span>
                <span v-if="!sidebarCollapsed && stats.total_courses > 0" class="ml-auto bg-neutral-700 text-neutral-300 text-xs font-medium px-2 py-0.5 rounded-full">{{ formatCount(stats.total_courses) }}</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Kursus
                </div>
              </NuxtLink>
              
              <NuxtLink 
                to="/admin/transactions" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/transactions') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Transaksi</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Transaksi
                </div>
              </NuxtLink>
              
              <NuxtLink 
                to="/admin/coupons" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/coupons') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A2 2 0 013 12V7a4 4 0 014-4z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Kupon</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Kupon
                </div>
              </NuxtLink>

              <NuxtLink 
                to="/admin/campaigns" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/campaigns') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Campaigns</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Campaigns
                </div>
              </NuxtLink>

              <NuxtLink 
                to="/admin/blog" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/blog') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Blog</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Blog
                </div>
              </NuxtLink>

              <NuxtLink 
                to="/admin/webinars" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/webinars') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Webinars</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Webinars
                </div>
              </NuxtLink>
              
              <NuxtLink 
                to="/admin/instructors" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/instructors') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Instruktur</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Instruktur
                </div>
              </NuxtLink>
              
              <NuxtLink 
                to="/admin/categories" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/categories') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Kategori</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Kategori
                </div>
              </NuxtLink>
            </div>
          </div>
          
          <div>
            <p v-if="!sidebarCollapsed" class="px-3 mb-3 text-xs font-semibold text-neutral-500 uppercase tracking-wider">Sistem</p>
            <div class="space-y-1">
              <NuxtLink 
                to="/admin/settings" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/settings') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Pengaturan</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Pengaturan
                </div>
              </NuxtLink>
              
              <NuxtLink 
                to="/admin/ai-processing" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/admin/ai-processing') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714c0 .597.237 1.17.659 1.591L19.8 15.3M14.25 3.104c.251.023.501.05.75.082M19.8 15.3l-1.57.393A9.065 9.065 0 0112 15a9.065 9.065 0 00-6.23.693L5 15.3m14.8 0l.21 1.847a1.5 1.5 0 01-1.48 1.71H5.47a1.5 1.5 0 01-1.48-1.71L5 15.3"/>
                </svg>
                <span v-if="!sidebarCollapsed">AI Processing</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  AI Processing
                </div>
              </NuxtLink>
            </div>
          </div>
        </nav>
        
        <!-- Toggle Button (collapsed state) -->
        <div v-if="sidebarCollapsed" class="p-3 border-t border-neutral-800">
          <button @click="toggleSidebar" class="w-full flex items-center justify-center p-2.5 text-neutral-500 hover:text-white hover:bg-neutral-800 rounded-lg transition-colors" title="Buka Sidebar">
            <svg class="w-5 h-5 rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 19l-7-7 7-7m8 14l-7-7 7-7"/></svg>
          </button>
        </div>
      </div>
    </aside>
    
    <!-- Mobile Header -->
    <header class="lg:hidden fixed top-0 left-0 right-0 z-50 h-16 bg-neutral-900 flex items-center justify-between px-4">
      <div class="flex items-center space-x-3">
        <button @click="mobileMenuOpen = true" class="p-2 -ml-2 text-neutral-400 hover:text-white">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 6h16M4 12h16M4 18h16"/>
          </svg>
        </button>
        <img src="/logo.png" alt="EDUKRA Logo" class="w-8 h-8 object-contain" />
        <span class="text-xs text-admin-400 font-medium">Admin</span>
      </div>
    </header>
    
    <!-- Main Content -->
    <div class="transition-all duration-300" :class="sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-72'">
      <!-- Top Bar -->
      <header class="hidden lg:flex h-16 bg-white border-b border-neutral-200 items-center justify-between px-8">
        <div class="flex items-center gap-4">
          <span class="px-3 py-1.5 bg-admin-50 text-admin-700 text-xs font-semibold rounded-full">Admin Panel</span>
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
            <input type="text" placeholder="Cari..." class="w-64 pl-9 pr-4 py-2 bg-neutral-100 border-0 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 focus:bg-white text-sm placeholder-neutral-500 transition-all"/>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <!-- Notification Bell -->
          <div class="relative" ref="notificationRef">
            <button @click="notificationOpen = !notificationOpen" class="relative p-2 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/></svg>
              <span v-if="unreadCount > 0" class="absolute top-1 right-1 w-4 h-4 bg-red-500 rounded-full text-[10px] text-white font-bold flex items-center justify-center">{{ unreadCount > 9 ? '9+' : unreadCount }}</span>
            </button>
            
            <!-- Notification Dropdown -->
            <Transition name="dropdown">
              <div v-if="notificationOpen" class="absolute right-0 mt-2 w-80 bg-white rounded-xl shadow-lg border border-neutral-200 overflow-hidden z-50">
                <div class="p-4 border-b border-neutral-100 flex items-center justify-between">
                  <h3 class="font-semibold text-neutral-900">Notifikasi</h3>
                  <button v-if="unreadCount > 0" @click="markAllRead" class="text-xs text-admin-600 hover:text-admin-700 font-medium">Tandai dibaca</button>
                </div>
                <div class="max-h-80 overflow-y-auto">
                  <div v-if="notifications.length === 0" class="p-6 text-center text-neutral-500 text-sm">
                    Tidak ada notifikasi
                  </div>
                  <div 
                    v-for="notif in notifications" 
                    :key="notif.id" 
                    @click="handleNotificationClick(notif)"
                    class="p-4 hover:bg-neutral-50 border-b border-neutral-100 last:border-0 cursor-pointer transition-colors"
                    :class="!notif.is_read ? 'bg-admin-50/50' : ''"
                  >
                    <div class="flex gap-3">
                      <div 
                        class="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0"
                        :class="getNotificationColor(notif.type)"
                      >
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" :d="getNotificationIcon(notif.type)"/>
                        </svg>
                      </div>
                      <div class="flex-1 min-w-0">
                        <p class="text-sm text-neutral-900" :class="!notif.is_read ? 'font-medium' : ''">{{ notif.title }}</p>
                        <p class="text-xs text-neutral-500 mt-0.5 truncate">{{ notif.message }}</p>
                        <p class="text-xs text-neutral-400 mt-1">{{ notif.time }}</p>
                      </div>
                      <div v-if="!notif.is_read" class="w-2 h-2 bg-admin-500 rounded-full flex-shrink-0 mt-2"></div>
                    </div>
                  </div>
                </div>
                <div class="p-3 border-t border-neutral-100 bg-neutral-50">
                  <button class="w-full text-center text-sm text-admin-600 hover:text-admin-700 font-medium">Lihat Semua Notifikasi</button>
                </div>
              </div>
            </Transition>
          </div>
          <div class="relative" ref="userDropdownRef">
            <button @click="userDropdownOpen = !userDropdownOpen" class="flex items-center gap-3 p-2 hover:bg-neutral-100 rounded-lg transition-colors">
              <div class="w-8 h-8 bg-admin-600 rounded-full flex items-center justify-center text-white font-semibold text-sm">A</div>
              <div class="text-left hidden sm:block"><p class="text-sm font-medium text-neutral-900">Admin</p></div>
              <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 9l-7 7-7-7"/></svg>
            </button>
            <Transition name="dropdown">
              <div v-if="userDropdownOpen" class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-neutral-200 py-2 z-50">
                <NuxtLink to="/admin/settings" class="flex items-center px-4 py-2 text-sm text-neutral-700 hover:bg-neutral-100">Pengaturan</NuxtLink>
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
      <div v-if="mobileMenuOpen" class="fixed inset-0 z-50 bg-black/60 lg:hidden" @click="mobileMenuOpen = false"></div>
    </Transition>
    
    <!-- Mobile Menu -->
    <Transition name="slide">
      <aside v-if="mobileMenuOpen" class="fixed top-0 left-0 z-50 w-72 h-screen bg-neutral-900 lg:hidden">
        <div class="h-full flex flex-col">
          <div class="h-16 flex items-center justify-between px-6 border-b border-neutral-800">
            <div class="flex items-center space-x-3">
              <img src="/logo.png" alt="EDUKRA Logo" class="w-10 h-10 object-contain" />
              <span class="font-display font-bold text-lg text-white">EDUKRA Admin</span>
            </div>
            <button @click="mobileMenuOpen = false" class="p-2 text-neutral-500 hover:text-white">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <nav class="flex-1 px-4 py-6 overflow-y-auto">
            <!-- Overview -->
            <div class="mb-6">
              <p class="px-3 mb-2 text-xs font-semibold text-neutral-500 uppercase tracking-wider">Overview</p>
              <div class="space-y-1">
                <NuxtLink to="/admin" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin', true) ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z"/>
                  </svg>
                  Dashboard
                </NuxtLink>
              </div>
            </div>
            
            <!-- Manajemen -->
            <div class="mb-6">
              <p class="px-3 mb-2 text-xs font-semibold text-neutral-500 uppercase tracking-wider">Manajemen</p>
              <div class="space-y-1">
                <NuxtLink to="/admin/users" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin/users') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"/>
                  </svg>
                  Pengguna
                </NuxtLink>
                <NuxtLink to="/admin/courses" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin/courses') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                  </svg>
                  Kursus
                </NuxtLink>
                <NuxtLink to="/admin/transactions" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin/transactions') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z"/>
                  </svg>
                  Transaksi
                </NuxtLink>
                <NuxtLink to="/admin/coupons" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin/coupons') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A2 2 0 013 12V7a4 4 0 014-4z"/>
                  </svg>
                  Kupon
                </NuxtLink>
                <NuxtLink to="/admin/campaigns" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin/campaigns') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"/>
                  </svg>
                  Campaigns
                </NuxtLink>
                <NuxtLink to="/admin/blog" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin/blog') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"/>
                  </svg>
                  Blog
                </NuxtLink>
                <NuxtLink to="/admin/instructors" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin/instructors') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
                  </svg>
                  Instruktur
                </NuxtLink>
                <NuxtLink to="/admin/categories" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin/categories') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"/>
                  </svg>
                  Kategori
                </NuxtLink>
              </div>
            </div>
            
            <!-- Sistem -->
            <div>
              <p class="px-3 mb-2 text-xs font-semibold text-neutral-500 uppercase tracking-wider">Sistem</p>
              <div class="space-y-1">
                <NuxtLink to="/admin/settings" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin/settings') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                  </svg>
                  Pengaturan
                </NuxtLink>
                <NuxtLink to="/admin/ai-processing" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/admin/ai-processing') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714c0 .597.237 1.17.659 1.591L19.8 15.3M14.25 3.104c.251.023.501.05.75.082M19.8 15.3l-1.57.393A9.065 9.065 0 0112 15a9.065 9.065 0 00-6.23.693L5 15.3m14.8 0l.21 1.847a1.5 1.5 0 01-1.48 1.71H5.47a1.5 1.5 0 01-1.48-1.71L5 15.3"/>
                  </svg>
                  AI Processing
                </NuxtLink>
              </div>
            </div>
          </nav>
        </div>
      </aside>
    </Transition>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const token = useCookie('token')
const userCookie = useCookie('user')
const api = useApi()


const mobileMenuOpen = ref(false)
const sidebarCollapsed = ref(false)
const userDropdownOpen = ref(false)
const userDropdownRef = ref<HTMLElement | null>(null)
const notificationOpen = ref(false)
const notificationRef = ref<HTMLElement | null>(null)

// Notifications - fetched from API
interface Notification {
  id: string
  type: string
  title: string
  message?: string
  reference_id?: string
  reference_type?: string
  is_read: boolean
  created_at: string
  time: string
}
const notifications = ref<Notification[]>([])

const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length)

// Fetch notifications from API
const fetchNotifications = async () => {
  try {
    const response = await api.fetch<{ notifications: Notification[], unread_count: number }>('/api/notifications')
    if (response?.notifications) {
      notifications.value = response.notifications
    }
  } catch (err) {
    console.error('[Admin] Error fetching notifications:', err)
  }
}

const getNotificationColor = (type: string) => {
  const colors: Record<string, string> = {
    user: 'bg-primary-100 text-primary-600',
    payment: 'bg-accent-100 text-accent-600',
    course: 'bg-warm-100 text-warm-600',
    alert: 'bg-admin-100 text-admin-600'
  }
  return colors[type] || 'bg-neutral-100 text-neutral-600'
}

const getNotificationIcon = (type: string) => {
  const icons: Record<string, string> = {
    user: 'M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z',
    payment: 'M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z',
    course: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253',
    alert: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
  }
  return icons[type] || icons.alert
}

const formatTimeAgo = (date: Date) => {
  const seconds = Math.floor((new Date().getTime() - date.getTime()) / 1000)
  if (seconds < 60) return 'Baru saja'
  if (seconds < 3600) return `${Math.floor(seconds / 60)} menit lalu`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)} jam lalu`
  return `${Math.floor(seconds / 86400)} hari lalu`
}

const markAllRead = async () => {
  try {
    await api.fetch('/api/notifications/read-all', { method: 'PUT' })
    notifications.value.forEach(n => n.is_read = true)
  } catch (err) {
    console.error('[Admin] Error marking all read:', err)
  }
}

const handleNotificationClick = async (notif: Notification) => {
  if (!notif.is_read) {
    try {
      await api.fetch(`/api/notifications/${notif.id}/read`, { method: 'PUT' })
      notif.is_read = true
    } catch (err) {
      console.error('[Admin] Error marking notification read:', err)
    }
  }
  notificationOpen.value = false
  // Navigate based on notification type
}

// Dashboard stats for badges
const stats = ref({
  total_users: 0,
  total_courses: 0,
  total_transactions: 0
})

const formatCount = (count: number) => {
  if (count >= 1000) {
    return (count / 1000).toFixed(1).replace(/\.0$/, '') + 'K'
  }
  return count.toString()
}

const isActive = (path: string, exact = false) => {
  if (exact) return route.path === path
  return route.path === path || route.path.startsWith(path + '/')
}

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
  if (process.client) {
    localStorage.setItem('admin-sidebar-collapsed', String(sidebarCollapsed.value))
  }
}

const fetchDashboardStats = async () => {
  try {
    const response = await api.fetch<{ stats: typeof stats.value }>('/api/admin/dashboard')
    if (response.stats) {
      stats.value = response.stats
    }
  } catch (err) {
    // Silently fail - badges will show 0
  }
}

onMounted(() => {
  if (process.client) {
    const saved = localStorage.getItem('admin-sidebar-collapsed')
    if (saved === 'true') {
      sidebarCollapsed.value = true
    }
    document.addEventListener('click', (e) => {
      if (userDropdownRef.value && !userDropdownRef.value.contains(e.target as Node)) {
        userDropdownOpen.value = false
      }
      if (notificationRef.value && !notificationRef.value.contains(e.target as Node)) {
        notificationOpen.value = false
      }
    })
    
    // Fetch stats for sidebar badges
    fetchDashboardStats()
    
    // Fetch notifications
    fetchNotifications()
  }
})

const handleLogout = () => {
  // Clear cookies
  token.value = null
  userCookie.value = null
  
  // Close dropdown
  userDropdownOpen.value = false
  
  // Redirect to admin login
  navigateTo('/admin/login')
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.slide-enter-active, .slide-leave-active { transition: transform 0.3s ease; }
.slide-enter-from, .slide-leave-to { transform: translateX(-100%); }
.dropdown-enter-active, .dropdown-leave-active { transition: all 0.2s ease; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(-10px); }
</style>

