<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { NavigationTab } from '../types'
  import { userDataStore } from '../stores/userDataStore'
  import { isSettingsOpen, isSidebarCollapsed, toggleSidebar } from '../stores/settingsStore'
  import { activePlaylistId, openCreatePlaylistModal, openImportPlaylistModal } from '../stores/playlistStore'
  import { t } from '../stores/i18nStore'

  export let activeTab: NavigationTab = 'search'

  const dispatch = createEventDispatcher<{ tabChange: NavigationTab }>()

  function setTab(tab: NavigationTab) {
    activeTab = tab
    if (tab !== 'settings') {
      isSettingsOpen.set(false)
    }
    dispatch('tabChange', tab)
  }

  function selectPlaylist(id: string) {
    activePlaylistId.set(id)
    setTab('playlist')
  }

  $: favCount = $userDataStore.favorites.length
  $: histCount = $userDataStore.history.length

  // Global high-contrast tooltip for collapsed mode (escapes overflow clipping via fixed position)
  let activeTooltip: { text: string; top: number; left: number; badge?: string } | null = null

  function showTooltip(e: MouseEvent, text: string, badge?: string) {
    if (!$isSidebarCollapsed) return
    const target = e.currentTarget as HTMLElement
    if (!target) return
    const rect = target.getBoundingClientRect()
    activeTooltip = {
      text,
      top: rect.top + rect.height / 2,
      left: rect.right + 10,
      badge,
    }
  }

  function hideTooltip() {
    activeTooltip = null
  }

  $: if (!$isSidebarCollapsed) {
    activeTooltip = null
  }
</script>

<aside
  class="bg-surface border-r border-white/[0.06] flex flex-col shrink-0 select-none justify-between h-full sidebar-animated relative z-50 overflow-visible {$isSidebarCollapsed ? 'w-16' : 'w-56 xl:w-60'}"
>
  <!-- Brand & Primary Nav Section -->
  <div
    class="space-y-3.5 {$isSidebarCollapsed ? 'px-2 py-3' : 'px-3 py-3'} overflow-y-auto overflow-x-hidden flex-1"
    on:scroll={hideTooltip}
  >
    <!-- Logo and Header -->
    <div class="flex items-center {$isSidebarCollapsed ? 'justify-center' : 'justify-between'} h-10 overflow-hidden">
      <div class="flex items-center min-w-0 {$isSidebarCollapsed ? 'justify-center w-full' : ''}">
        <!-- Logo Box (Click to toggle collapse) -->
        <button
          class="w-10 h-10 rounded-xl bg-gradient-to-tr from-crimson-600 via-crimson-500 to-rose-400 flex items-center justify-center shadow-lg shadow-crimson-500/25 shrink-0 cursor-pointer hover:scale-105 active:scale-95 transition-all mx-auto"
          on:click={toggleSidebar}
          on:mouseenter={(e) => showTooltip(e, $isSidebarCollapsed ? $t('sidebar_expand') : $t('sidebar_collapse'))}
          on:mouseleave={hideTooltip}
          aria-label={$isSidebarCollapsed ? $t('sidebar_expand') : $t('sidebar_collapse')}
        >
          <svg class="w-5 h-5 text-white" viewBox="0 0 269.22 267.55" fill="none">
            <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
          </svg>
        </button>

        <!-- Animated Brand Text (Expanded only) -->
        <div class="min-w-0 sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-w-0 opacity-0 -translate-x-2 pointer-events-none' : 'max-w-[140px] opacity-100 translate-x-0 ml-3'}">
          <div class="flex items-center gap-1.5 leading-none">
            <span class="font-extrabold text-base tracking-tight text-white truncate">Dengerin</span>
            <span class="text-[9px] font-bold px-1.5 py-0.5 rounded bg-crimson-500/20 border border-crimson-500/30 text-crimson-400 uppercase tracking-widest shrink-0">Hi-Fi</span>
          </div>
          <span class="text-[10px] text-zinc-400 font-medium leading-tight truncate block mt-0.5">Master Audio Engine</span>
        </div>
      </div>

      <!-- Animated Collapse Toggle Button (Expanded only) -->
      <div class="sidebar-text-transition overflow-hidden shrink-0 {$isSidebarCollapsed ? 'max-w-0 opacity-0 pointer-events-none' : 'max-w-8 opacity-100'}">
        <button
          class="p-1.5 rounded-xl text-zinc-400 hover:text-white hover:bg-white/[0.08] transition active:scale-90 cursor-pointer shrink-0"
          on:click={toggleSidebar}
          title={$t('sidebar_collapse')}
          aria-label={$t('sidebar_collapse')}
        >
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect width="18" height="18" x="3" y="3" rx="2"/>
            <path d="M9 3v18"/>
            <path d="m16 15-3-3 3-3"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Section: Menjelajah -->
    <div class="space-y-1">
      <div class="sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-h-0 opacity-0 -translate-x-2 pointer-events-none my-0' : 'max-h-6 opacity-100 translate-x-0 px-2.5 mb-1'}">
        <p class="text-[10px] font-bold text-zinc-400 uppercase tracking-wider">
          {$t('section_browse')}
        </p>
      </div>
      
      <!-- Tab: Eksplorasi & Cari -->
      <div class="flex {$isSidebarCollapsed ? 'justify-center w-full' : 'w-full'}">
        <button
          class="{$isSidebarCollapsed ? 'w-10 h-10 p-0 justify-center' : 'w-full h-10 px-2.5 justify-start'} flex items-center rounded-xl text-sm font-medium transition-colors duration-200 cursor-pointer relative {activeTab === 'search' ? 'bg-white/[0.1] text-white shadow-sm border border-white/[0.08]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04] border border-transparent'}"
          on:click={() => setTab('search')}
          on:mouseenter={(e) => showTooltip(e, $t('nav_explore'))}
          on:mouseleave={hideTooltip}
          aria-label={$t('nav_explore')}
        >
          <div class="w-5 h-5 flex items-center justify-center shrink-0">
            <svg class="w-4 h-4 {activeTab === 'search' ? 'text-crimson-400' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <circle cx="11" cy="11" r="8"></circle>
              <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
            </svg>
          </div>

          <div class="flex items-center flex-1 min-w-0 sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-w-0 opacity-0 -translate-x-2 pointer-events-none' : 'max-w-[160px] opacity-100 translate-x-0 ml-3'}">
            <span class="truncate">{$t('nav_explore')}</span>
            {#if activeTab === 'search'}
              <span class="ml-auto text-[10px] bg-crimson-500/20 text-crimson-400 font-semibold px-1.5 py-0.5 rounded shrink-0">{$t('nav_active')}</span>
            {/if}
          </div>
        </button>
      </div>

      <!-- Tab: Untuk Anda (For You) -->
      <div class="flex {$isSidebarCollapsed ? 'justify-center w-full' : 'w-full'}">
        <button
          class="{$isSidebarCollapsed ? 'w-10 h-10 p-0 justify-center' : 'w-full h-10 px-2.5 justify-start'} flex items-center rounded-xl text-sm font-medium transition-colors duration-200 cursor-pointer relative {activeTab === 'foryou' ? 'bg-white/[0.1] text-white shadow-sm border border-white/[0.08]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04] border border-transparent'}"
          on:click={() => setTab('foryou')}
          on:mouseenter={(e) => showTooltip(e, $t('nav_foryou'))}
          on:mouseleave={hideTooltip}
          aria-label={$t('nav_foryou')}
        >
          <div class="w-5 h-5 flex items-center justify-center shrink-0">
            <svg class="w-4 h-4 {activeTab === 'foryou' ? 'text-crimson-400' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"></path>
            </svg>
          </div>

          <div class="flex items-center flex-1 min-w-0 sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-w-0 opacity-0 -translate-x-2 pointer-events-none' : 'max-w-[160px] opacity-100 translate-x-0 ml-3'}">
            <span class="truncate">{$t('nav_foryou')}</span>
            {#if activeTab === 'foryou'}
              <span class="ml-auto text-[10px] bg-crimson-500/20 text-crimson-400 font-semibold px-1.5 py-0.5 rounded shrink-0">{$t('nav_active')}</span>
            {/if}
          </div>
        </button>
      </div>
    </div>

    <!-- Section: Koleksi Saya -->
    <div class="space-y-1">
      <div class="sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-h-0 opacity-0 -translate-x-2 pointer-events-none my-0' : 'max-h-6 opacity-100 translate-x-0 px-2.5 mb-1'}">
        <p class="text-[10px] font-bold text-zinc-400 uppercase tracking-wider">
          {$t('section_collection')}
        </p>
      </div>

      <!-- Lagu Favorit -->
      <div class="flex {$isSidebarCollapsed ? 'justify-center w-full' : 'w-full'}">
        <button
          class="{$isSidebarCollapsed ? 'w-10 h-10 p-0 justify-center' : 'w-full h-10 px-2.5 justify-start'} flex items-center rounded-xl text-sm font-medium transition-colors duration-200 cursor-pointer relative {activeTab === 'favorites' ? 'bg-white/[0.1] text-white shadow-sm border border-white/[0.08]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04] border border-transparent'}"
          on:click={() => setTab('favorites')}
          on:mouseenter={(e) => showTooltip(e, $t('nav_favorites'), favCount > 0 ? String(favCount) : undefined)}
          on:mouseleave={hideTooltip}
          aria-label={$t('nav_favorites')}
        >
          <div class="w-5 h-5 flex items-center justify-center shrink-0">
            <svg class="w-4 h-4 {activeTab === 'favorites' ? 'fill-crimson-500 text-crimson-500' : 'text-zinc-400 group-hover:text-crimson-400'} transition" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
            </svg>
          </div>

          {#if favCount > 0 && $isSidebarCollapsed}
            <span class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-crimson-500 ring-2 ring-[#12141C] shadow-[0_0_6px_#E61E38] pointer-events-none"></span>
          {/if}

          <div class="flex items-center flex-1 min-w-0 sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-w-0 opacity-0 -translate-x-2 pointer-events-none' : 'max-w-[160px] opacity-100 translate-x-0 ml-3'}">
            <span class="truncate">{$t('nav_favorites')}</span>
            {#if favCount > 0}
              <span class="ml-auto text-xs bg-white/[0.06] text-zinc-300 px-2 py-0.5 rounded-full font-mono shrink-0">{favCount}</span>
            {:else if activeTab === 'favorites'}
              <span class="ml-auto text-[10px] bg-crimson-500/20 text-crimson-400 font-semibold px-1.5 py-0.5 rounded shrink-0">{$t('nav_active')}</span>
            {/if}
          </div>
        </button>
      </div>

      <!-- Riwayat Putar -->
      <div class="flex {$isSidebarCollapsed ? 'justify-center w-full' : 'w-full'}">
        <button
          class="{$isSidebarCollapsed ? 'w-10 h-10 p-0 justify-center' : 'w-full h-10 px-2.5 justify-start'} flex items-center rounded-xl text-sm font-medium transition-colors duration-200 cursor-pointer relative {activeTab === 'history' ? 'bg-white/[0.1] text-white shadow-sm border border-white/[0.08]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04] border border-transparent'}"
          on:click={() => setTab('history')}
          on:mouseenter={(e) => showTooltip(e, $t('nav_history'), histCount > 0 ? String(histCount) : undefined)}
          on:mouseleave={hideTooltip}
          aria-label={$t('nav_history')}
        >
          <div class="w-5 h-5 flex items-center justify-center shrink-0">
            <svg class="w-4 h-4 {activeTab === 'history' ? 'text-crimson-400' : 'text-zinc-400 group-hover:text-zinc-200'} transition" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
          </div>

          {#if histCount > 0 && $isSidebarCollapsed}
            <span class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-crimson-500 ring-2 ring-[#12141C] shadow-[0_0_6px_#E61E38] pointer-events-none"></span>
          {/if}

          <div class="flex items-center flex-1 min-w-0 sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-w-0 opacity-0 -translate-x-2 pointer-events-none' : 'max-w-[160px] opacity-100 translate-x-0 ml-3'}">
            <span class="truncate">{$t('nav_history')}</span>
            {#if histCount > 0}
              <span class="ml-auto text-xs bg-white/[0.06] text-zinc-300 px-2 py-0.5 rounded-full font-mono shrink-0">{histCount}</span>
            {:else if activeTab === 'history'}
              <span class="ml-auto text-[10px] bg-crimson-500/20 text-crimson-400 font-semibold px-1.5 py-0.5 rounded shrink-0">{$t('nav_active')}</span>
            {/if}
          </div>
        </button>
      </div>

      <!-- Tersimpan / Offline -->
      <div class="flex {$isSidebarCollapsed ? 'justify-center w-full' : 'w-full'}">
        <button
          class="{$isSidebarCollapsed ? 'w-10 h-10 p-0 justify-center' : 'w-full h-10 px-2.5 justify-start'} flex items-center rounded-xl text-sm font-medium transition-colors duration-200 cursor-pointer relative {activeTab === 'offline' ? 'bg-white/[0.1] text-white shadow-sm border border-white/[0.08]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04] border border-transparent'}"
          on:click={() => setTab('offline')}
          on:mouseenter={(e) => showTooltip(e, $t('nav_offline'))}
          on:mouseleave={hideTooltip}
          aria-label={$t('nav_offline')}
        >
          <div class="relative w-5 h-5 flex items-center justify-center shrink-0">
            <svg class="w-4 h-4 {activeTab === 'offline' ? 'text-crimson-400' : 'text-zinc-400 group-hover:text-zinc-200'} transition" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
              <polyline points="7 10 12 15 17 10"></polyline>
              <line x1="12" y1="15" x2="12" y2="3"></line>
            </svg>
          </div>

          <div class="flex items-center flex-1 min-w-0 sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-w-0 opacity-0 -translate-x-2 pointer-events-none' : 'max-w-[160px] opacity-100 translate-x-0 ml-3'}">
            <span class="truncate">{$t('nav_offline')}</span>
            {#if activeTab === 'offline'}
              <span class="ml-auto text-[10px] bg-crimson-500/20 text-crimson-400 font-semibold px-1.5 py-0.5 rounded shrink-0">{$t('nav_active')}</span>
            {/if}
          </div>
        </button>
      </div>

      <!-- Playlists Section -->
      <div class="pt-2">
        <div class="flex items-center justify-between sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-h-0 opacity-0 -translate-x-2 pointer-events-none my-0' : 'max-h-6 opacity-100 translate-x-0 px-2.5 mb-1.5'}">
          <p class="text-[10px] font-bold text-zinc-400 uppercase tracking-wider">
            {$t('section_playlists')}
          </p>
          <div class="flex items-center gap-1">
            <button
              class="p-1 rounded-lg text-zinc-400 hover:text-white hover:bg-white/[0.08] transition cursor-pointer"
              on:click={() => openImportPlaylistModal()}
              title={$t('import_playlist')}
              aria-label={$t('import_playlist')}
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="7 10 12 15 17 10"></polyline>
                <line x1="12" y1="3" x2="12" y2="15"></line>
              </svg>
            </button>
            <button
              class="p-1 rounded-lg text-zinc-400 hover:text-white hover:bg-white/[0.08] transition cursor-pointer"
              on:click={() => openCreatePlaylistModal()}
              title={$t('create_playlist')}
              aria-label={$t('create_playlist')}
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
              </svg>
            </button>
          </div>
        </div>

        {#if $isSidebarCollapsed}
          <div class="flex justify-center w-full py-1">
            <button
              class="w-10 h-10 rounded-xl bg-white/[0.03] hover:bg-white/[0.08] text-zinc-400 hover:text-white border border-dashed border-white/10 hover:border-crimson-500/50 flex items-center justify-center transition-all cursor-pointer shadow-xs active:scale-95"
              on:click={() => openCreatePlaylistModal()}
              on:mouseenter={(e) => showTooltip(e, $t('create_playlist'))}
              on:mouseleave={hideTooltip}
              aria-label={$t('create_playlist')}
            >
              <svg class="w-4 h-4 text-crimson-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
              </svg>
            </button>
          </div>
        {/if}

        <div class="space-y-0.5">
          {#each $userDataStore.playlists as pl (pl.id)}
            <div class="flex {$isSidebarCollapsed ? 'justify-center w-full' : 'w-full'}">
              <button
                class="{$isSidebarCollapsed ? 'w-10 h-10 p-0 justify-center' : 'w-full h-10 px-2.5 justify-start'} flex items-center rounded-xl text-xs font-medium transition-colors duration-200 cursor-pointer relative {activeTab === 'playlist' && $activePlaylistId === pl.id ? 'bg-white/[0.1] text-white shadow-sm border border-white/[0.08]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04] border border-transparent'}"
                on:click={() => selectPlaylist(pl.id)}
                on:mouseenter={(e) => showTooltip(e, pl.name, (pl.tracks || []).length > 0 ? String((pl.tracks || []).length) : undefined)}
                on:mouseleave={hideTooltip}
                aria-label={pl.name}
              >
                <div class="w-5 h-5 flex items-center justify-center shrink-0">
                  <svg class="w-4 h-4 {activeTab === 'playlist' && $activePlaylistId === pl.id ? 'text-crimson-400' : 'text-zinc-500'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <path d="M9 18V5l12-2v13"></path>
                    <circle cx="6" cy="18" r="3"></circle>
                    <circle cx="18" cy="16" r="3"></circle>
                  </svg>
                </div>

                <div class="flex items-center flex-1 min-w-0 sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-w-0 opacity-0 -translate-x-2 pointer-events-none' : 'max-w-[160px] opacity-100 translate-x-0 ml-3'}">
                  <span class="truncate">{pl.name}</span>
                  {#if (pl.tracks || []).length > 0}
                    <span class="ml-auto text-[10px] text-zinc-500 font-mono shrink-0">{(pl.tracks || []).length}</span>
                  {/if}
                </div>
              </button>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </div>

  <!-- Bottom Sidebar: Clean Settings trigger button -->
  <div class="border-t border-white/[0.06] bg-[#0E1017] {$isSidebarCollapsed ? 'p-2 flex justify-center' : 'p-3'} relative overflow-visible">
    <div class="flex {$isSidebarCollapsed ? 'justify-center w-full' : 'w-full'}">
      <button
        class="{$isSidebarCollapsed ? 'w-10 h-10 p-0 justify-center' : 'w-full h-10 px-2.5 justify-start'} flex items-center rounded-xl text-sm font-medium transition-colors cursor-pointer {activeTab === 'settings'
          ? 'bg-crimson-500/15 border border-crimson-500/30 text-white shadow-xs'
          : 'text-zinc-400 hover:text-zinc-100 hover:bg-white/[0.05] border border-transparent hover:border-white/[0.06]'}"
        on:click={() => setTab('settings')}
        on:mouseenter={(e) => showTooltip(e, $t('settings_title'))}
        on:mouseleave={hideTooltip}
        aria-label={$t('settings_title')}
      >
        <div class="w-5 h-5 flex items-center justify-center shrink-0">
          <svg class="w-4.5 h-4.5 {activeTab === 'settings' ? 'text-crimson-400 rotate-45' : 'text-zinc-400 group-hover:text-crimson-400 group-hover:rotate-45'} transition duration-200 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="3"></circle>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
          </svg>
        </div>

        <div class="flex items-center flex-1 min-w-0 sidebar-text-transition overflow-hidden whitespace-nowrap {$isSidebarCollapsed ? 'max-w-0 opacity-0 -translate-x-2 pointer-events-none' : 'max-w-[160px] opacity-100 translate-x-0 ml-3'}">
          <span class="truncate">{$t('settings_title')}</span>
          {#if activeTab === 'settings'}
            <span class="ml-auto text-[10px] bg-crimson-500/20 text-crimson-400 font-semibold px-1.5 py-0.5 rounded shrink-0">{$t('nav_active')}</span>
          {/if}
        </div>
      </button>
    </div>
  </div>

  <!-- High-contrast fixed tooltip overlay for collapsed mode (never clipped by overflow) -->
  {#if $isSidebarCollapsed && activeTooltip}
    <div
      class="fixed z-[9999] pointer-events-none px-2.5 py-1.5 rounded-xl bg-[#161824]/95 border border-white/10 text-xs font-semibold text-white whitespace-nowrap shadow-2xl backdrop-blur-md flex items-center gap-1.5 -translate-y-1/2 animate-scale-in"
      style="top: {activeTooltip.top}px; left: {activeTooltip.left}px;"
    >
      <span>{activeTooltip.text}</span>
      {#if activeTooltip.badge}
        <span class="text-[10px] bg-crimson-500/20 text-crimson-300 px-1.5 py-0.5 rounded-full font-mono">{activeTooltip.badge}</span>
      {/if}
      <span class="absolute right-full top-1/2 -translate-y-1/2 border-4 border-transparent border-r-[#161824]/95"></span>
    </div>
  {/if}
</aside>

<style>
  aside.sidebar-animated {
    transition: width 280ms cubic-bezier(0.4, 0, 0.2, 1);
    will-change: width;
  }
  .sidebar-text-transition {
    transition: opacity 200ms ease, max-width 280ms cubic-bezier(0.4, 0, 0.2, 1), transform 200ms ease, max-height 280ms cubic-bezier(0.4, 0, 0.2, 1);
    will-change: opacity, max-width, transform;
  }
</style>
