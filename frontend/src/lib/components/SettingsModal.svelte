<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { fade, scale } from 'svelte/transition'
  import { cubicOut } from 'svelte/easing'
  import { isSettingsOpen, audioQualityStore, showToast } from '../stores/settingsStore'
  import { currentLocale, t } from '../stores/i18nStore'
  import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime.js'
  import { CheckForUpdates } from '../../../wailsjs/go/main/App'

  let activeTab: 'general' | 'audio' | 'about' = 'general'
  let updateStatus: 'idle' | 'checking' | 'latest' | 'update_available' | 'error' = 'idle'
  let latestVersion = ''
  let releaseUrl = ''

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && $isSettingsOpen) {
      isSettingsOpen.set(false)
    }
  }

  onMount(() => {
    window.addEventListener('keydown', handleKeydown)
  })

  onDestroy(() => {
    window.removeEventListener('keydown', handleKeydown)
  })

  function handleClearCache() {
    try {
      showToast($t('toast_cleared_cache'))
    } catch (e) {
      console.warn('Clear cache error:', e)
    }
  }

  async function handleCheckUpdate() {
    updateStatus = 'checking'
    try {
      const res = await CheckForUpdates()
      if (res && res.hasUpdate) {
        updateStatus = 'update_available'
        latestVersion = res.latestVersion
        releaseUrl = res.releaseUrl || 'https://github.com/komporkode/dengar/releases'
      } else {
        updateStatus = 'latest'
      }
    } catch (e) {
      console.warn('Update check failed:', e)
      updateStatus = 'error'
    }
  }

  function openGithub() {
    try {
      BrowserOpenURL('https://github.com/komporkode/dengar')
    } catch (e) {
      window.open('https://github.com/komporkode/dengar', '_blank')
    }
  }

  function openReleasePage() {
    const url = releaseUrl || 'https://github.com/komporkode/dengar/releases'
    try {
      BrowserOpenURL(url)
    } catch (e) {
      window.open(url, '_blank')
    }
  }
</script>

{#if $isSettingsOpen}
  <!-- Backdrop with fade transition -->
  <div
    class="fixed inset-0 bg-black/65 backdrop-blur-md z-50 flex items-center justify-center p-4 select-none"
    transition:fade={{ duration: 180 }}
    on:click={() => isSettingsOpen.set(false)}
    role="presentation"
  >
    <!-- Modal Card (Optimized 700px width, 490px height, spacious & zero scroll) -->
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class="w-full max-w-2xl md:max-w-[700px] h-[490px] bg-[#11131E]/95 border border-white/[0.12] rounded-3xl shadow-2xl flex flex-col md:flex-row overflow-hidden relative"
      transition:scale={{ start: 0.96, duration: 180, easing: cubicOut }}
      on:click|stopPropagation
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <!-- Left Tab Rail -->
      <div class="w-full md:w-52 bg-[#0A0C14]/90 border-b md:border-b-0 md:border-r border-white/[0.08] p-4 flex flex-col justify-between shrink-0 h-auto md:h-full">
        <div>
          <!-- Logo & Brand Header in Settings -->
          <div class="flex items-center gap-2.5 px-2 py-2 mb-4 border-b border-white/[0.06]">
            <div class="w-7 h-7 rounded-lg bg-gradient-to-tr from-crimson-600 to-rose-500 flex items-center justify-center shadow-xs shrink-0">
              <svg class="w-3.5 h-3.5 text-white" viewBox="0 0 269.22 267.55" fill="none">
                <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
              </svg>
            </div>
            <div>
              <h2 class="text-xs font-bold text-white tracking-wide">{$t('settings_title')}</h2>
              <p class="text-[10px] text-zinc-500 font-mono">Dengar Desktop</p>
            </div>
          </div>

          <!-- Tab Selection Navigation -->
          <nav class="space-y-1">
            <button
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer text-left {activeTab === 'general' ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25' : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
              on:click={() => (activeTab = 'general')}
            >
              <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="2" y1="12" x2="22" y2="12"></line>
                <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
              </svg>
              <span>{$t('tab_general')}</span>
            </button>

            <button
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer text-left {activeTab === 'audio' ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25' : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
              on:click={() => (activeTab = 'audio')}
            >
              <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
                <path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"></path>
              </svg>
              <span>{$t('tab_audio')}</span>
            </button>

            <button
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer text-left {activeTab === 'about' ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25' : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
              on:click={() => (activeTab = 'about')}
            >
              <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="16" x2="12" y2="12"></line>
                <line x1="12" y1="8" x2="12.01" y2="8"></line>
              </svg>
              <span>{$t('tab_about')}</span>
            </button>
          </nav>
        </div>

        <div class="pt-3 border-t border-white/[0.06] hidden md:block">
          <p class="text-[10px] text-zinc-400 font-medium">Dengar Desktop</p>
          <p class="text-[10px] text-zinc-500 font-mono">v1.1.0 • Open Source</p>
        </div>
      </div>

      <!-- Right Content Area (Clean, Airy & Zero-Scroll Layout) -->
      <div class="flex-1 flex flex-col justify-between p-5 md:p-6 overflow-hidden h-full">
        <!-- Top Section Header with Clean Close 'X' Button -->
        <div class="flex items-start justify-between pb-3 border-b border-white/[0.06] shrink-0">
          <div>
            <h3 class="text-sm font-bold text-white">
              {activeTab === 'general' ? $t('tab_general') : activeTab === 'audio' ? $t('tab_audio') : $t('tab_about')}
            </h3>
            <p class="text-[11px] text-zinc-400 mt-0.5">
              {activeTab === 'general' ? 'Atur preferensi bahasa antarmuka dan kendali aplikasi.' : activeTab === 'audio' ? 'Pengaturan resolusi streaming dan manajemen cache audio.' : 'Informasi aplikasi, repositori sumber terbuka, dan pembaruan.'}
            </p>
          </div>
          <button
            class="w-7 h-7 rounded-lg bg-white/[0.05] hover:bg-white/[0.12] text-zinc-400 hover:text-white transition flex items-center justify-center cursor-pointer shrink-0 -mt-1 -mr-1"
            on:click={() => isSettingsOpen.set(false)}
            title={$t('close') + ' (Esc)'}
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <!-- Middle Content Area (No scrollbars, perfectly sized) -->
        <div class="flex-1 py-3 flex flex-col justify-center">
          {#if activeTab === 'general'}
            <!-- TAB 1: General & Language -->
            <div class="space-y-3.5">
              <div>
                <span class="text-[10px] font-semibold text-zinc-400 uppercase tracking-wider">
                  {$t('lang_select')}
                </span>
                <div class="grid grid-cols-2 gap-3 mt-1.5">
                  <button
                    class="flex items-center justify-between p-3 rounded-2xl border transition-all cursor-pointer {$currentLocale === 'id' ? 'bg-crimson-500/15 border-crimson-500 shadow-sm shadow-crimson-500/15 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                    on:click={() => currentLocale.set('id')}
                  >
                    <div class="flex items-center gap-2.5">
                      <span class="w-7 h-7 rounded-xl bg-white/[0.08] flex items-center justify-center font-bold text-[11px] font-mono">ID</span>
                      <div class="text-left">
                        <p class="text-xs font-bold leading-tight">{$t('lang_id')}</p>
                        <p class="text-[10px] text-zinc-400">Bahasa Baku</p>
                      </div>
                    </div>
                    {#if $currentLocale === 'id'}
                      <div class="w-4 h-4 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0">
                        <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"></polyline>
                        </svg>
                      </div>
                    {/if}
                  </button>

                  <button
                    class="flex items-center justify-between p-3 rounded-2xl border transition-all cursor-pointer {$currentLocale === 'en' ? 'bg-crimson-500/15 border-crimson-500 shadow-sm shadow-crimson-500/15 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                    on:click={() => currentLocale.set('en')}
                  >
                    <div class="flex items-center gap-2.5">
                      <span class="w-7 h-7 rounded-xl bg-white/[0.08] flex items-center justify-center font-bold text-[11px] font-mono">EN</span>
                      <div class="text-left">
                        <p class="text-xs font-bold leading-tight">{$t('lang_en')}</p>
                        <p class="text-[10px] text-zinc-400">International</p>
                      </div>
                    </div>
                    {#if $currentLocale === 'en'}
                      <div class="w-4 h-4 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0">
                        <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"></polyline>
                        </svg>
                      </div>
                    {/if}
                  </button>
                </div>
              </div>

              <!-- Shortcut Reference Card -->
              <div class="p-3 rounded-2xl bg-white/[0.02] border border-white/[0.06]">
                <div class="flex items-center gap-1.5 text-zinc-300 text-xs font-semibold mb-2">
                  <svg class="w-3.5 h-3.5 text-crimson-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <rect x="2" y="4" width="20" height="16" rx="2" ry="2"></rect>
                    <path d="M6 8h.001M10 8h.001M14 8h.001M18 8h.001M6 12h.001M18 12h.001M8 16h8"></path>
                  </svg>
                  <span>Pintasan Keyboard Global</span>
                </div>
                <div class="grid grid-cols-2 gap-x-4 gap-y-1.5 text-[11px]">
                  <div class="flex items-center justify-between text-zinc-400">
                    <span>Play / Pause</span>
                    <kbd class="px-1.5 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[10px]">Spasi</kbd>
                  </div>
                  <div class="flex items-center justify-between text-zinc-400">
                    <span>Atur Volume</span>
                    <kbd class="px-1.5 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[10px]">↑ / ↓</kbd>
                  </div>
                  <div class="flex items-center justify-between text-zinc-400">
                    <span>Maju / Mundur 5d</span>
                    <kbd class="px-1.5 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[10px]">← / →</kbd>
                  </div>
                  <div class="flex items-center justify-between text-zinc-400">
                    <span>Lagu Sblm / Ssdh</span>
                    <kbd class="px-1.5 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[10px]">J / L</kbd>
                  </div>
                  <div class="flex items-center justify-between text-zinc-400">
                    <span>Bisukan (Mute)</span>
                    <kbd class="px-1.5 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[10px]">M</kbd>
                  </div>
                  <div class="flex items-center justify-between text-zinc-400">
                    <span>Kembali / Tutup</span>
                    <kbd class="px-1.5 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[10px]">Esc</kbd>
                  </div>
                </div>
              </div>
            </div>

          {:else if activeTab === 'audio'}
            <!-- TAB 2: Audio Quality & Cache -->
            <div class="space-y-3">
              <div>
                <span class="text-[10px] font-semibold text-zinc-400 uppercase tracking-wider">
                  {$t('quality_setting')}
                </span>
                <div class="grid grid-cols-2 gap-3 mt-1.5">
                  <button
                    class="p-3 rounded-2xl border text-left transition-all cursor-pointer flex items-center justify-between {$audioQualityStore === 'high' ? 'bg-crimson-500/15 border-crimson-500 shadow-sm shadow-crimson-500/15 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                    on:click={() => audioQualityStore.set('high')}
                  >
                    <div class="flex items-center gap-2.5">
                      <span class="px-1.5 py-0.5 rounded-lg bg-amber-500/20 text-amber-300 font-mono text-[10px] font-bold border border-amber-500/30">HQ</span>
                      <div>
                        <div class="font-bold text-xs text-white leading-tight">{$t('quality_high')}</div>
                        <div class="text-[10px] text-zinc-400 mt-0.5">Opus 160kbps / AAC • Resolusi Master</div>
                      </div>
                    </div>
                    {#if $audioQualityStore === 'high'}
                      <div class="w-4 h-4 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0">
                        <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"></polyline>
                        </svg>
                      </div>
                    {/if}
                  </button>

                  <button
                    class="p-3 rounded-2xl border text-left transition-all cursor-pointer flex items-center justify-between {$audioQualityStore === 'saver' ? 'bg-crimson-500/15 border-crimson-500 shadow-sm shadow-crimson-500/15 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                    on:click={() => audioQualityStore.set('saver')}
                  >
                    <div class="flex items-center gap-2.5">
                      <span class="px-1.5 py-0.5 rounded-lg bg-emerald-500/20 text-emerald-300 font-mono text-[10px] font-bold border border-emerald-500/30">ECO</span>
                      <div>
                        <div class="font-bold text-xs text-white leading-tight">{$t('quality_saver')}</div>
                        <div class="text-[10px] text-zinc-400 mt-0.5">Buffer cepat & hemat bandwidth</div>
                      </div>
                    </div>
                    {#if $audioQualityStore === 'saver'}
                      <div class="w-4 h-4 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0">
                        <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"></polyline>
                        </svg>
                      </div>
                    {/if}
                  </button>
                </div>
              </div>

              <!-- Compact Notice Box -->
              <div class="p-2.5 rounded-xl bg-amber-500/10 border border-amber-500/20 text-[11px] text-amber-300/90 leading-relaxed flex items-center gap-2">
                <svg class="w-3.5 h-3.5 text-amber-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="8" x2="12" y2="12"></line>
                  <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                <span>{$t('quality_notice')}</span>
              </div>

              <!-- Cache Clean Bar -->
              <div class="p-3 rounded-2xl bg-white/[0.02] border border-white/[0.06] flex items-center justify-between">
                <div>
                  <h4 class="text-xs font-bold text-white">Cache & Stream Resolver</h4>
                  <p class="text-[10px] text-zinc-400 mt-0.5">Bersihkan data cache streaming audio sementara di perangkat</p>
                </div>
                <button
                  class="px-3.5 py-1.5 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-xs font-semibold text-zinc-200 hover:text-white transition cursor-pointer flex items-center gap-1.5 shrink-0"
                  on:click={handleClearCache}
                >
                  <svg class="w-3 h-3 text-zinc-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                  <span>{$t('clear_cache_btn')}</span>
                </button>
              </div>
            </div>

          {:else if activeTab === 'about'}
            <!-- TAB 3: About & Updates -->
            <div class="space-y-3">
              <!-- Compact Unified Brand Card -->
              <div class="p-3.5 rounded-2xl bg-gradient-to-br from-[#161824] to-[#10121C] border border-white/10 flex items-center gap-3.5 shadow-sm">
                <div class="w-12 h-12 rounded-2xl bg-gradient-to-tr from-crimson-600 via-crimson-500 to-rose-400 flex items-center justify-center shadow-lg shadow-crimson-500/25 shrink-0">
                  <svg class="w-6 h-6 text-white" viewBox="0 0 269.22 267.55" fill="none">
                    <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
                  </svg>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <h4 class="text-sm font-extrabold text-white">Dengar Desktop</h4>
                    <span class="text-[9px] font-bold px-1.5 py-0.2 rounded bg-crimson-500/20 border border-crimson-500/30 text-crimson-400 font-mono">v1.1.0</span>
                  </div>
                  <p class="text-[11px] text-zinc-300 leading-snug mt-0.5 truncate">
                    {$t('about_desc')}
                  </p>
                  <div class="flex items-center gap-1.5 mt-1.5 text-[9px] font-medium text-zinc-300">
                    <span class="px-1.5 py-0.2 rounded bg-white/[0.06] border border-white/10">🚫 Zero Ads</span>
                    <span class="px-1.5 py-0.2 rounded bg-white/[0.06] border border-white/10">💰 100% Free / Rp0</span>
                    <span class="px-1.5 py-0.2 rounded bg-white/[0.06] border border-white/10">🪶 Low RAM</span>
                  </div>
                </div>
              </div>

              <!-- Combined GitHub & Update Card -->
              <div class="p-3 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-2.5">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2.5">
                    <div class="w-7 h-7 rounded-lg bg-white/10 flex items-center justify-center text-white shrink-0">
                      <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
                        <path fill-rule="evenodd" clip-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.53 1.032 1.53 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
                      </svg>
                    </div>
                    <div>
                      <h4 class="text-xs font-bold text-white leading-tight">GitHub & Komporkode</h4>
                      <p class="text-[10px] text-zinc-400 font-mono">GNU GPLv3 Open Source License</p>
                    </div>
                  </div>

                  <button
                    class="px-2.5 py-1 rounded-lg bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-xs font-semibold text-zinc-200 transition cursor-pointer flex items-center gap-1.5"
                    on:click={openGithub}
                    title="Buka Repositori di GitHub"
                  >
                    <span>Kunjungi</span>
                    <svg class="w-3 h-3 text-zinc-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                      <polyline points="15 3 21 3 21 9"></polyline>
                      <line x1="10" y1="14" x2="21" y2="3"></line>
                    </svg>
                  </button>
                </div>

                <!-- Update Status & Action Row -->
                <div class="pt-2 border-t border-white/[0.06] flex items-center justify-between text-xs">
                  <div class="flex items-center gap-2">
                    {#if updateStatus === 'checking'}
                      <div class="w-3 h-3 border-2 border-white/20 border-t-crimson-500 rounded-full animate-spin"></div>
                      <span class="text-zinc-400 text-[11px]">Memeriksa pembaruan di GitHub...</span>
                    {:else if updateStatus === 'latest'}
                      <span class="w-2 h-2 rounded-full bg-emerald-400 shadow-[0_0_6px_#34D399]"></span>
                      <span class="text-emerald-400 text-[11px] font-medium">Aplikasi sudah versi terbaru (v1.1.0)</span>
                    {:else if updateStatus === 'update_available'}
                      <span class="w-2 h-2 rounded-full bg-crimson-400 shadow-[0_0_6px_#E61E38] animate-pulse"></span>
                      <span class="text-crimson-300 text-[11px] font-bold">Versi baru tersedia: {latestVersion}</span>
                    {:else if updateStatus === 'error'}
                      <span class="w-2 h-2 rounded-full bg-amber-400"></span>
                      <span class="text-zinc-400 text-[11px]">Belum ada rilis baru di GitHub</span>
                    {:else}
                      <span class="text-zinc-400 text-[11px]">Status rilis: <strong class="text-zinc-200 font-mono">v1.1.0 Stable</strong></span>
                    {/if}
                  </div>

                  {#if updateStatus === 'update_available'}
                    <button
                      class="px-2.5 py-1 rounded-lg bg-crimson-500 hover:bg-crimson-600 text-white font-bold text-xs transition cursor-pointer shadow-sm shadow-crimson-500/30"
                      on:click={openReleasePage}
                    >
                      Unduh Update
                    </button>
                  {:else}
                    <button
                      class="px-2.5 py-1 rounded-lg bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-[11px] font-medium text-zinc-300 hover:text-white transition cursor-pointer disabled:opacity-50"
                      on:click={handleCheckUpdate}
                      disabled={updateStatus === 'checking'}
                    >
                      {updateStatus === 'checking' ? 'Memeriksa...' : 'Cek Pembaruan'}
                    </button>
                  {/if}
                </div>
              </div>
            </div>
          {/if}
        </div>

        <!-- Dialog Footer (Clean & Compact) -->
        <div class="pt-3 border-t border-white/[0.06] flex items-center justify-between shrink-0">
          <span class="text-[10px] text-zinc-500 font-mono">Dengar Desktop • Tekan Esc untuk menutup</span>
          <button
            class="px-4 py-1.5 rounded-xl bg-crimson-500 hover:bg-crimson-600 text-white font-semibold text-xs transition cursor-pointer shadow-md shadow-crimson-500/25"
            on:click={() => isSettingsOpen.set(false)}
          >
            {$t('close')}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
