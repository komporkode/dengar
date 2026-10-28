<script lang="ts">
  import { isSettingsOpen, audioQualityStore, showToast } from '../stores/settingsStore'
  import { currentLocale, t } from '../stores/i18nStore'

  let activeTab: 'general' | 'audio' | 'about' = 'general'

  function handleClearCache() {
    try {
      showToast($t('toast_cleared_cache'))
    } catch (e) {
      console.warn('Clear cache error:', e)
    }
  }
</script>

{#if $isSettingsOpen}
  <!-- Backdrop -->
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4 select-none animate-in fade-in duration-200"
    on:click={() => isSettingsOpen.set(false)}
    role="presentation"
  >
    <!-- Modal Card (Modern Luxury Hi-Fi Glassmorphism) -->
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class="w-full max-w-2xl bg-[#11131E]/95 border border-white/[0.12] rounded-3xl shadow-2xl flex flex-col md:flex-row overflow-hidden animate-scale-in duration-200 max-h-[85vh]"
      on:click|stopPropagation
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <!-- Left Tab Rail -->
      <div class="w-full md:w-56 bg-[#0B0D15]/90 border-b md:border-b-0 md:border-r border-white/[0.08] p-4 flex flex-col justify-between shrink-0">
        <div>
          <!-- Logo & Brand Header in Settings -->
          <div class="flex items-center gap-2.5 px-2 py-2 mb-4 border-b border-white/[0.06]">
            <div class="w-7 h-7 rounded-lg bg-gradient-to-tr from-crimson-600 to-rose-500 flex items-center justify-center shadow-xs">
              <svg class="w-3.5 h-3.5 text-white" viewBox="0 0 269.22 267.55" fill="none">
                <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
              </svg>
            </div>
            <div>
              <h2 class="text-xs font-bold text-white tracking-wide">{$t('settings_title')}</h2>
              <p class="text-[10px] text-zinc-500 font-mono">Dengar Desktop</p>
            </div>
          </div>

          <!-- Tab Selection Buttons -->
          <nav class="space-y-1">
            <button
              class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-xs font-semibold transition-all cursor-pointer text-left {activeTab === 'general' ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25' : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
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
              class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-xs font-semibold transition-all cursor-pointer text-left {activeTab === 'audio' ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25' : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
              on:click={() => (activeTab = 'audio')}
            >
              <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
                <path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"></path>
              </svg>
              <span>{$t('tab_audio')}</span>
            </button>

            <button
              class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-xs font-semibold transition-all cursor-pointer text-left {activeTab === 'about' ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25' : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
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

        <div class="pt-4 border-t border-white/[0.06] hidden md:block">
          <p class="text-[10px] text-zinc-500">Dengar v1.0.0</p>
          <p class="text-[10px] text-zinc-600 font-mono">Hi-Fi Audio Player</p>
        </div>
      </div>

      <!-- Right Content Area -->
      <div class="flex-1 flex flex-col justify-between overflow-y-auto p-6 space-y-6">
        <div>
          {#if activeTab === 'general'}
            <div class="space-y-4">
              <div>
                <h3 class="text-sm font-bold text-white mb-1">{$t('tab_general')}</h3>
                <p class="text-xs text-zinc-400">Atur preferensi bahasa antarmuka aplikasi.</p>
              </div>

              <div class="space-y-2 pt-2">
                <div class="text-[11px] font-semibold text-zinc-400 uppercase tracking-wider">
                  {$t('lang_select')}
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <button
                    class="flex items-center justify-between p-4 rounded-2xl border transition-all cursor-pointer {$currentLocale === 'id' ? 'bg-crimson-500/15 border-crimson-500 shadow-sm shadow-crimson-500/15 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                    on:click={() => currentLocale.set('id')}
                  >
                    <div class="flex items-center gap-3 font-medium">
                      <span class="w-8 h-8 rounded-xl bg-white/[0.08] flex items-center justify-center font-bold text-xs font-mono">ID</span>
                      <div class="text-left">
                        <p class="text-xs font-bold">{$t('lang_id')}</p>
                        <p class="text-[10px] text-zinc-400">Bahasa Baku</p>
                      </div>
                    </div>
                    {#if $currentLocale === 'id'}
                      <div class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white">
                        <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"></polyline>
                        </svg>
                      </div>
                    {/if}
                  </button>

                  <button
                    class="flex items-center justify-between p-4 rounded-2xl border transition-all cursor-pointer {$currentLocale === 'en' ? 'bg-crimson-500/15 border-crimson-500 shadow-sm shadow-crimson-500/15 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                    on:click={() => currentLocale.set('en')}
                  >
                    <div class="flex items-center gap-3 font-medium">
                      <span class="w-8 h-8 rounded-xl bg-white/[0.08] flex items-center justify-center font-bold text-xs font-mono">EN</span>
                      <div class="text-left">
                        <p class="text-xs font-bold">{$t('lang_en')}</p>
                        <p class="text-[10px] text-zinc-400">International</p>
                      </div>
                    </div>
                    {#if $currentLocale === 'en'}
                      <div class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white">
                        <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"></polyline>
                        </svg>
                      </div>
                    {/if}
                  </button>
                </div>
              </div>
            </div>

          {:else if activeTab === 'audio'}
            <div class="space-y-4">
              <div>
                <h3 class="text-sm font-bold text-white mb-1">{$t('tab_audio')}</h3>
                <p class="text-xs text-zinc-400">Pengaturan resolusi streaming dan manajemen penyimpanan sementara audio.</p>
              </div>

              <!-- Stream Quality Options -->
              <div class="space-y-2 pt-2">
                <div class="text-[11px] font-semibold text-zinc-400 uppercase tracking-wider">
                  {$t('quality_setting')}
                </div>
                <div class="space-y-2.5">
                  <button
                    class="w-full flex items-center justify-between p-4 rounded-2xl border text-left transition-all cursor-pointer {$audioQualityStore === 'high' ? 'bg-crimson-500/15 border-crimson-500 shadow-sm shadow-crimson-500/15 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                    on:click={() => audioQualityStore.set('high')}
                  >
                    <div class="flex items-center gap-3">
                      <span class="px-2 py-1 rounded-lg bg-amber-500/20 text-amber-300 font-mono text-[10px] font-bold border border-amber-500/30">HQ</span>
                      <div>
                        <div class="font-semibold text-xs text-white">{$t('quality_high')}</div>
                        <div class="text-[11px] text-zinc-400 mt-0.5">Format audio Opus & AAC resolusi tertinggi dari YouTube</div>
                      </div>
                    </div>
                    {#if $audioQualityStore === 'high'}
                      <div class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0">
                        <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"></polyline>
                        </svg>
                      </div>
                    {/if}
                  </button>

                  <button
                    class="w-full flex items-center justify-between p-4 rounded-2xl border text-left transition-all cursor-pointer {$audioQualityStore === 'saver' ? 'bg-crimson-500/15 border-crimson-500 shadow-sm shadow-crimson-500/15 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                    on:click={() => audioQualityStore.set('saver')}
                  >
                    <div class="flex items-center gap-3">
                      <span class="px-2 py-1 rounded-lg bg-white/10 text-zinc-300 font-mono text-[10px] font-bold">ECO</span>
                      <div>
                        <div class="font-semibold text-xs text-white">{$t('quality_saver')}</div>
                        <div class="text-[11px] text-zinc-400 mt-0.5">Mengutamakan kecepatan buffer & hemat bandwidth jaringan</div>
                      </div>
                    </div>
                    {#if $audioQualityStore === 'saver'}
                      <div class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0">
                        <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"></polyline>
                        </svg>
                      </div>
                    {/if}
                  </button>
                </div>
              </div>

              <!-- Audio Notice Box -->
              <div class="p-3.5 rounded-2xl bg-amber-500/10 border border-amber-500/20 text-xs text-amber-300/90 leading-relaxed flex items-start gap-2.5">
                <svg class="w-4 h-4 text-amber-400 shrink-0 mt-0.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="8" x2="12" y2="12"></line>
                  <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                <span>{$t('quality_notice')}</span>
              </div>

              <!-- Cache clear action -->
              <div class="pt-2">
                <button
                  class="px-4 py-2.5 rounded-xl bg-white/[0.06] hover:bg-white/[0.1] border border-white/10 text-xs font-semibold text-zinc-200 transition cursor-pointer flex items-center gap-2"
                  on:click={handleClearCache}
                >
                  <svg class="w-3.5 h-3.5 text-zinc-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                  <span>{$t('clear_cache_btn')}</span>
                </button>
              </div>
            </div>

          {:else if activeTab === 'about'}
            <div class="space-y-4">
              <!-- Hi-Fi Player Hero Card in About -->
              <div class="p-6 rounded-2xl bg-gradient-to-br from-[#181A28] to-[#12141F] border border-white/10 flex flex-col items-center text-center shadow-lg relative overflow-hidden">
                <div class="w-16 h-16 rounded-2xl bg-gradient-to-tr from-crimson-600 via-crimson-500 to-rose-400 flex items-center justify-center shadow-xl shadow-crimson-500/30 mb-3">
                  <svg class="w-9 h-9 text-white" viewBox="0 0 269.22 267.55" fill="none">
                    <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
                  </svg>
                </div>
                <div class="flex items-center gap-2 mb-1">
                  <h4 class="text-lg font-extrabold text-white">Dengar Desktop</h4>
                  <span class="text-[9px] font-bold px-1.5 py-0.5 rounded bg-crimson-500/20 border border-crimson-500/30 text-crimson-400 uppercase font-mono">v1.0.0</span>
                </div>
                <p class="text-xs text-zinc-300 max-w-sm leading-relaxed mb-4">
                  {$t('about_desc')}
                </p>

                <!-- Feature Pill Badges -->
                <div class="flex flex-wrap items-center justify-center gap-1.5 text-[10px] font-medium text-zinc-300">
                  <span class="px-2 py-0.5 rounded-full bg-white/[0.06] border border-white/10">🚫 Zero Ads</span>
                  <span class="px-2 py-0.5 rounded-full bg-white/[0.06] border border-white/10">💰 100% Free / Rp0</span>
                  <span class="px-2 py-0.5 rounded-full bg-white/[0.06] border border-white/10">⚡ Go + Svelte + Wails</span>
                  <span class="px-2 py-0.5 rounded-full bg-white/[0.06] border border-white/10">🪶 Low RAM Usage</span>
                </div>
              </div>

              <!-- Credits & License -->
              <div class="p-4 rounded-2xl bg-white/[0.02] border border-white/[0.06] flex items-center justify-between text-xs">
                <div>
                  <p class="text-[11px] text-zinc-500">{$t('creator')}</p>
                  <p class="text-xs font-bold text-white mt-0.5">Komporkode</p>
                </div>
                <div class="text-right">
                  <p class="text-[11px] text-zinc-500">Lisensi</p>
                  <p class="text-xs font-semibold text-crimson-400 mt-0.5 font-mono">MIT Open Source</p>
                </div>
              </div>
            </div>
          {/if}
        </div>

        <!-- Dialog Footer -->
        <div class="pt-4 border-t border-white/[0.08] flex items-center justify-end">
          <button
            class="px-5 py-2 rounded-xl bg-crimson-500 hover:bg-crimson-600 text-white font-semibold text-xs transition cursor-pointer shadow-md shadow-crimson-500/25"
            on:click={() => isSettingsOpen.set(false)}
          >
            {$t('close')}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
