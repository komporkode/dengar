<script lang="ts">
  import { isSettingsOpen, audioQualityStore, showToast } from '../stores/settingsStore'
  import { currentLocale, t } from '../stores/i18nStore'

  let activeTab: 'general' | 'audio' | 'about' = 'general'

  function handleClearCache() {
    // Clear localStorage search/stream cache if any
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
    class="fixed inset-0 bg-black/70 backdrop-blur-md z-50 flex items-center justify-center p-4 select-none"
    on:click={() => isSettingsOpen.set(false)}
    role="presentation"
  >
    <!-- Modal Card -->
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class="w-full max-w-lg bg-surface border border-white/10 rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-in fade-in zoom-in-95 duration-150"
      on:click|stopPropagation
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <!-- Modal Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-white/[0.08] bg-[#0E1017]">
        <div class="flex items-center gap-2.5">
          <svg class="w-5 h-5 text-crimson-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="3"></circle>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
          </svg>
          <h2 class="text-base font-bold text-white">{$t('settings_title')}</h2>
        </div>
        <button
          class="p-1 rounded-lg text-zinc-400 hover:text-white hover:bg-white/[0.08] transition cursor-pointer"
          on:click={() => isSettingsOpen.set(false)}
          title={$t('close')}
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <!-- Navigation Tabs -->
      <div class="flex items-center px-6 pt-3 gap-2 border-b border-white/[0.06] bg-[#0E1017]">
        <button
          class="pb-2.5 px-3 text-xs font-semibold border-b-2 transition cursor-pointer {activeTab === 'general' ? 'border-crimson-500 text-white' : 'border-transparent text-zinc-400 hover:text-zinc-200'}"
          on:click={() => (activeTab = 'general')}
        >
          {$t('tab_general')}
        </button>
        <button
          class="pb-2.5 px-3 text-xs font-semibold border-b-2 transition cursor-pointer {activeTab === 'audio' ? 'border-crimson-500 text-white' : 'border-transparent text-zinc-400 hover:text-zinc-200'}"
          on:click={() => (activeTab = 'audio')}
        >
          {$t('tab_audio')}
        </button>
        <button
          class="pb-2.5 px-3 text-xs font-semibold border-b-2 transition cursor-pointer {activeTab === 'about' ? 'border-crimson-500 text-white' : 'border-transparent text-zinc-400 hover:text-zinc-200'}"
          on:click={() => (activeTab = 'about')}
        >
          {$t('tab_about')}
        </button>
      </div>

      <!-- Tab Contents -->
      <div class="p-6 space-y-5 text-sm">
        {#if activeTab === 'general'}
          <div>
            <div class="block text-xs font-semibold text-zinc-300 uppercase tracking-wider mb-2">
              {$t('lang_select')}
            </div>
            <div class="grid grid-cols-2 gap-3">
              <button
                class="flex items-center justify-between p-3.5 rounded-xl border transition cursor-pointer {$currentLocale === 'id' ? 'bg-crimson-500/15 border-crimson-500 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                on:click={() => currentLocale.set('id')}
              >
                <div class="flex items-center gap-2.5 font-medium">
                  <span class="text-xs font-bold px-1.5 py-0.5 rounded bg-white/10 font-mono">ID</span>
                  <span>{$t('lang_id')}</span>
                </div>
                {#if $currentLocale === 'id'}
                  <svg class="w-4 h-4 text-crimson-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                {/if}
              </button>

              <button
                class="flex items-center justify-between p-3.5 rounded-xl border transition cursor-pointer {$currentLocale === 'en' ? 'bg-crimson-500/15 border-crimson-500 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                on:click={() => currentLocale.set('en')}
              >
                <div class="flex items-center gap-2.5 font-medium">
                  <span class="text-xs font-bold px-1.5 py-0.5 rounded bg-white/10 font-mono">EN</span>
                  <span>{$t('lang_en')}</span>
                </div>
                {#if $currentLocale === 'en'}
                  <svg class="w-4 h-4 text-crimson-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                {/if}
              </button>
            </div>
          </div>

        {:else if activeTab === 'audio'}
          <div class="space-y-4">
            <div>
              <div class="block text-xs font-semibold text-zinc-300 uppercase tracking-wider mb-2">
                {$t('quality_setting')}
              </div>
              <div class="space-y-2">
                <button
                  class="w-full flex items-center justify-between p-3.5 rounded-xl border text-left transition cursor-pointer {$audioQualityStore === 'high' ? 'bg-crimson-500/15 border-crimson-500 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                  on:click={() => audioQualityStore.set('high')}
                >
                  <div>
                    <div class="font-medium text-sm">{$t('quality_high')}</div>
                    <div class="text-xs text-zinc-400 mt-0.5">Format webm/opus dengan bitrate tertinggi yang disediakan server</div>
                  </div>
                  {#if $audioQualityStore === 'high'}
                    <svg class="w-4 h-4 text-crimson-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                      <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                  {/if}
                </button>

                <button
                  class="w-full flex items-center justify-between p-3.5 rounded-xl border text-left transition cursor-pointer {$audioQualityStore === 'saver' ? 'bg-crimson-500/15 border-crimson-500 text-white' : 'bg-white/[0.03] border-white/10 text-zinc-300 hover:bg-white/[0.06]'}"
                  on:click={() => audioQualityStore.set('saver')}
                >
                  <div>
                    <div class="font-medium text-sm">{$t('quality_saver')}</div>
                    <div class="text-xs text-zinc-400 mt-0.5">Mengutamakan kecepatan buffer & hemat bandwidth jaringan</div>
                  </div>
                  {#if $audioQualityStore === 'saver'}
                    <svg class="w-4 h-4 text-crimson-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                      <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                  {/if}
                </button>
              </div>
            </div>

            <div class="p-3 rounded-xl bg-amber-500/10 border border-amber-500/20 text-xs text-amber-300 leading-relaxed">
              {$t('quality_notice')}
            </div>

            <div class="pt-2">
              <button
                class="px-4 py-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.1] border border-white/10 text-xs font-semibold text-zinc-200 transition cursor-pointer"
                on:click={handleClearCache}
              >
                {$t('clear_cache_btn')}
              </button>
            </div>
          </div>

        {:else if activeTab === 'about'}
          <div class="flex flex-col items-center text-center space-y-4 py-2">
            <div class="w-16 h-16 rounded-2xl bg-gradient-to-tr from-crimson-600 via-crimson-500 to-rose-400 flex items-center justify-center shadow-lg shadow-crimson-500/30">
              <svg class="w-9 h-9 text-white" viewBox="0 0 269.22 267.55" fill="none">
                <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
              </svg>
            </div>
            <div>
              <h3 class="text-xl font-extrabold text-white">Dengar</h3>
              <p class="text-xs font-mono text-crimson-400 font-semibold mt-0.5">v1.0.0 — Production</p>
            </div>
            <p class="text-xs text-zinc-300 max-w-sm leading-relaxed">
              {$t('about_desc')}
            </p>
            <div class="w-full pt-3 border-t border-white/[0.08] text-xs text-zinc-400 space-y-1">
              <p>{$t('built_with')}</p>
              <p class="font-medium text-zinc-300">{$t('creator')}: <span class="text-crimson-400 font-bold">Komporkode</span></p>
            </div>
          </div>
        {/if}
      </div>

      <!-- Modal Footer -->
      <div class="px-6 py-3 bg-[#0E1017] border-t border-white/[0.08] flex justify-end">
        <button
          class="px-5 py-2 rounded-xl bg-crimson-500 hover:bg-crimson-600 text-white font-semibold text-xs transition cursor-pointer"
          on:click={() => isSettingsOpen.set(false)}
        >
          {$t('close')}
        </button>
      </div>
    </div>
  </div>
{/if}
