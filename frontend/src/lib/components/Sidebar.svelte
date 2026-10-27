<script lang="ts">
  import type { NavigationTab } from '../types'
  import { userDataStore } from '../stores/userDataStore'

  export let activeTab: NavigationTab = 'search'

  $: favCount = $userDataStore.favorites.length
  $: histCount = $userDataStore.history.length
</script>

<aside class="w-64 bg-surface border-r border-white/[0.06] flex flex-col shrink-0 select-none justify-between h-full">
  <!-- Brand & Primary Nav Section -->
  <div class="p-4 space-y-6 overflow-y-auto">
    <!-- Logo and Hi-Fi Indicator -->
    <div class="flex items-center justify-between px-2 pt-1">
      <div class="flex items-center gap-2.5">
        <div class="w-8 h-8 rounded-xl bg-gradient-to-tr from-crimson-600 via-crimson-500 to-rose-400 flex items-center justify-center shadow-lg shadow-crimson-500/20">
          <svg class="w-4 h-4 text-white" viewBox="0 0 269.22 267.55" fill="none">
            <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
          </svg>
        </div>
        <div>
          <div class="flex items-center gap-1.5 leading-none">
            <span class="font-extrabold text-base tracking-tight text-white">Dengar</span>
            <span class="text-[9px] font-bold px-1.5 py-0.5 rounded bg-crimson-500/20 border border-crimson-500/30 text-crimson-400 uppercase tracking-widest">Hi-Fi</span>
          </div>
          <span class="text-[10px] text-zinc-400 font-medium leading-tight">Master Audio Engine</span>
        </div>
      </div>
    </div>

    <!-- Section: Menjelajah -->
    <div class="space-y-1">
      <p class="px-2.5 text-[11px] font-semibold text-zinc-400 uppercase tracking-wider">Menjelajah</p>
      <button
        class="w-full flex items-center gap-3 px-3 py-2 rounded-xl text-sm font-medium transition cursor-pointer text-left {activeTab === 'search' ? 'bg-white/[0.08] text-white shadow-sm border border-white/[0.04]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04]'}"
        on:click={() => (activeTab = 'search')}
      >
        <svg class="w-4 h-4 {activeTab === 'search' ? 'text-crimson-400' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <span>Eksplorasi & Cari</span>
        {#if activeTab === 'search'}
          <span class="ml-auto text-[10px] bg-crimson-500/20 text-crimson-400 font-semibold px-1.5 py-0.5 rounded">Aktif</span>
        {/if}
      </button>
    </div>

    <!-- Section: Koleksi Saya -->
    <div class="space-y-1">
      <p class="px-2.5 text-[11px] font-semibold text-zinc-400 uppercase tracking-wider">Koleksi Saya</p>
      <button
        class="w-full flex items-center justify-between px-3 py-2 rounded-xl text-sm font-medium transition cursor-pointer group text-left {activeTab === 'favorites' ? 'bg-white/[0.08] text-white shadow-sm border border-white/[0.04]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04]'}"
        on:click={() => (activeTab = 'favorites')}
      >
        <div class="flex items-center gap-3">
          <svg class="w-4 h-4 {activeTab === 'favorites' ? 'fill-crimson-500 text-crimson-500' : 'text-zinc-400 group-hover:text-crimson-400'} transition" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
          </svg>
          <span>Lagu Favorit</span>
        </div>
        {#if favCount > 0}
          <span class="text-xs bg-white/[0.06] text-zinc-300 px-2 py-0.5 rounded-full font-mono">{favCount}</span>
        {/if}
      </button>

      <button
        class="w-full flex items-center justify-between px-3 py-2 rounded-xl text-sm font-medium transition cursor-pointer group text-left {activeTab === 'history' ? 'bg-white/[0.08] text-white shadow-sm border border-white/[0.04]' : 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04]'}"
        on:click={() => (activeTab = 'history')}
      >
        <div class="flex items-center gap-3">
          <svg class="w-4 h-4 {activeTab === 'history' ? 'text-crimson-400' : 'text-zinc-400 group-hover:text-zinc-200'} transition" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="10"></circle>
            <polyline points="12 6 12 12 16 14"></polyline>
          </svg>
          <span>Riwayat Putar</span>
        </div>
        {#if histCount > 0}
          <span class="text-xs bg-white/[0.06] text-zinc-300 px-2 py-0.5 rounded-full font-mono">{histCount}</span>
        {/if}
      </button>
    </div>
  </div>

  <!-- Bottom Sidebar: Backend Status Chip -->
  <div class="p-3 border-t border-white/[0.06] bg-[#0E1017] space-y-2">
    <!-- Backend status -->
    <div class="flex items-center justify-between px-2.5 py-2 rounded-xl bg-white/[0.03] border border-white/[0.06]">
      <div class="flex items-center gap-2.5">
        <span class="relative flex h-2 w-2">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
        </span>
        <div>
          <p class="text-[11px] font-semibold text-zinc-200 leading-tight">Server Hi-Fi Publik</p>
          <p class="text-[10px] text-emerald-400 font-mono">Terhubung • 24-bit/96kHz</p>
        </div>
      </div>
      <svg class="w-3.5 h-3.5 text-zinc-500" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <rect height="8" rx="2" ry="2" width="20" x="2" y="2"></rect>
        <rect height="8" rx="2" ry="2" width="20" x="2" y="14"></rect>
        <line x1="6" x2="6.01" y1="6" y2="6"></line>
        <line x1="6" x2="6.01" y1="18" y2="18"></line>
      </svg>
    </div>
  </div>
</aside>
