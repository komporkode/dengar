<script lang="ts">
  import type { NavigationTab } from '../types'
  import { userDataStore } from '../stores/userDataStore'
  import { isSettingsOpen, isSidebarCollapsed, toggleSidebar } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'

  export let activeTab: NavigationTab = 'search'

  $: favCount = $userDataStore.favorites.length
  $: histCount = $userDataStore.history.length
</script>

<aside
  class="bg-surface border-r border-white/[0.06] flex flex-col shrink-0 select-none justify-between h-full transition-all duration-300 ease-in-out {$isSidebarCollapsed ? 'w-20' : 'w-56 xl:w-64'}"
>
  <!-- Brand & Primary Nav Section -->
  <div class="space-y-6 overflow-y-auto {$isSidebarCollapsed ? 'p-2' : 'p-3.5 xl:p-4'}">
    <!-- Logo and Collapse Toggle Header -->
    <div class="flex items-center justify-between {$isSidebarCollapsed ? 'flex-col gap-2 pt-1' : 'px-1.5 pt-1'}">
      <div class="flex items-center gap-2.5 min-w-0">
        <!-- Logo Box (Click to toggle collapse) -->
        <button
          class="relative group w-9 h-9 rounded-xl bg-gradient-to-tr from-crimson-600 via-crimson-500 to-rose-400 flex items-center justify-center shadow-lg shadow-crimson-500/25 shrink-0 cursor-pointer hover:scale-105 active:scale-95 transition-all"
          on:click={toggleSidebar}
          title={$isSidebarCollapsed ? 'Perluas Sidebar' : 'Ciutkan Sidebar'}
        >
          <svg class="w-4.5 h-4.5 text-white" viewBox="0 0 269.22 267.55" fill="none">
            <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
          </svg>
        </button>

        {#if !$isSidebarCollapsed}
          <div class="min-w-0 animate-in fade-in duration-200">
            <div class="flex items-center gap-1.5 leading-none">
              <span class="font-extrabold text-base tracking-tight text-white truncate">Dengar</span>
              <span class="text-[9px] font-bold px-1.5 py-0.5 rounded bg-crimson-500/20 border border-crimson-500/30 text-crimson-400 uppercase tracking-widest shrink-0">Hi-Fi</span>
            </div>
            <span class="text-[10px] text-zinc-400 font-medium leading-tight truncate block mt-0.5">Master Audio Engine</span>
          </div>
        {/if}
      </div>

      <!-- Modern Sidebar Collapse/Expand Toggle Button -->
      <button
        class="p-1.5 rounded-xl text-zinc-400 hover:text-white hover:bg-white/[0.08] transition active:scale-90 cursor-pointer shrink-0"
        on:click={toggleSidebar}
        title={$isSidebarCollapsed ? 'Perluas Sidebar' : 'Ciutkan Sidebar'}
      >
        {#if $isSidebarCollapsed}
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect width="18" height="18" x="3" y="3" rx="2"/>
            <path d="M9 3v18"/>
            <path d="m14 9 3 3-3 3"/>
          </svg>
        {:else}
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect width="18" height="18" x="3" y="3" rx="2"/>
            <path d="M9 3v18"/>
            <path d="m16 15-3-3 3-3"/>
          </svg>
        {/if}
      </button>
    </div>

    <!-- Section: Menjelajah -->
    <div class="space-y-1">
      {#if !$isSidebarCollapsed}
        <p class="px-2.5 text-[11px] font-semibold text-zinc-400 uppercase tracking-wider animate-in fade-in duration-150">
          {$t('section_browse')}
        </p>
      {/if}
      <button
        class="w-full flex items-center rounded-xl text-sm font-medium transition-all duration-150 cursor-pointer {$isSidebarCollapsed ? 'justify-center p-3' : 'gap-3 px-3 py-2.5'} {activeTab === 'search' ? 'bg-white/[0.08] text-white shadow-sm border border-white/[0.04]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04]'}"
        on:click={() => (activeTab = 'search')}
        title={$t('nav_explore')}
      >
        <svg class="w-4 h-4 shrink-0 {activeTab === 'search' ? 'text-crimson-400' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
        </svg>
        {#if !$isSidebarCollapsed}
          <span class="truncate">{$t('nav_explore')}</span>
          {#if activeTab === 'search'}
            <span class="ml-auto text-[10px] bg-crimson-500/20 text-crimson-400 font-semibold px-1.5 py-0.5 rounded shrink-0">{$t('nav_active')}</span>
          {/if}
        {/if}
      </button>
    </div>

    <!-- Section: Koleksi Saya -->
    <div class="space-y-1">
      {#if !$isSidebarCollapsed}
        <p class="px-2.5 text-[11px] font-semibold text-zinc-400 uppercase tracking-wider animate-in fade-in duration-150">
          {$t('section_collection')}
        </p>
      {/if}
      <!-- Lagu Favorit -->
      <button
        class="w-full flex items-center rounded-xl text-sm font-medium transition-all duration-150 cursor-pointer group {$isSidebarCollapsed ? 'justify-center p-3 relative' : 'justify-between px-3 py-2.5'} {activeTab === 'favorites' ? 'bg-white/[0.08] text-white shadow-sm border border-white/[0.04]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04]'}"
        on:click={() => (activeTab = 'favorites')}
        title={$t('nav_favorites')}
      >
        <div class="flex items-center gap-3 min-w-0">
          <svg class="w-4 h-4 shrink-0 {activeTab === 'favorites' ? 'fill-crimson-500 text-crimson-500' : 'text-zinc-400 group-hover:text-crimson-400'} transition" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
          </svg>
          {#if !$isSidebarCollapsed}
            <span class="truncate">{$t('nav_favorites')}</span>
          {/if}
        </div>
        {#if favCount > 0}
          {#if $isSidebarCollapsed}
            <span class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-crimson-500 shadow-xs"></span>
          {:else}
            <span class="text-xs bg-white/[0.06] text-zinc-300 px-2 py-0.5 rounded-full font-mono shrink-0 ml-1">{favCount}</span>
          {/if}
        {/if}
      </button>

      <!-- Riwayat Putar -->
      <button
        class="w-full flex items-center rounded-xl text-sm font-medium transition-all duration-150 cursor-pointer group {$isSidebarCollapsed ? 'justify-center p-3 relative' : 'justify-between px-3 py-2.5'} {activeTab === 'history' ? 'bg-white/[0.08] text-white shadow-sm border border-white/[0.04]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04]'}"
        on:click={() => (activeTab = 'history')}
        title={$t('nav_history')}
      >
        <div class="flex items-center gap-3 min-w-0">
          <svg class="w-4 h-4 shrink-0 {activeTab === 'history' ? 'text-crimson-400' : 'text-zinc-400 group-hover:text-zinc-200'} transition" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="10"></circle>
            <polyline points="12 6 12 12 16 14"></polyline>
          </svg>
          {#if !$isSidebarCollapsed}
            <span class="truncate">{$t('nav_history')}</span>
          {/if}
        </div>
        {#if histCount > 0}
          {#if $isSidebarCollapsed}
            <span class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-crimson-500 shadow-xs"></span>
          {:else}
            <span class="text-xs bg-white/[0.06] text-zinc-300 px-2 py-0.5 rounded-full font-mono shrink-0 ml-1">{histCount}</span>
          {/if}
        {/if}
      </button>
    </div>
  </div>

  <!-- Bottom Sidebar: Clean Settings & About trigger button -->
  <div class="border-t border-white/[0.06] bg-[#0E1017] {$isSidebarCollapsed ? 'p-2 flex justify-center' : 'p-3'}">
    <button
      class="flex items-center rounded-xl text-sm font-medium text-zinc-400 hover:text-zinc-100 hover:bg-white/[0.05] transition group border border-transparent hover:border-white/[0.06] cursor-pointer {$isSidebarCollapsed ? 'p-3 justify-center' : 'w-full gap-3 px-3 py-2.5'}"
      on:click={() => isSettingsOpen.set(true)}
      title={$t('settings_title')}
    >
      <svg class="w-4.5 h-4.5 text-zinc-400 group-hover:text-crimson-400 group-hover:rotate-45 transition duration-200 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="3"></circle>
        <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
      </svg>
      {#if !$isSidebarCollapsed}
        <span class="truncate">{$t('settings_title')}</span>
      {/if}
    </button>
  </div>
</aside>
