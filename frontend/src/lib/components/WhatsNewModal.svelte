<script lang="ts">
  import { CHANGELOG_DATA } from '../data/changelogs'
  import { t } from '../stores/i18nStore'
  import { fade, scale } from 'svelte/transition'

  export let isOpen = false
  export let onClose: () => void

  const currentRelease = CHANGELOG_DATA[0]

  function handleBackdrop(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      onClose()
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      onClose()
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if isOpen && currentRelease}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-md select-none"
    on:click={handleBackdrop}
    on:keydown={(e) => e.key === 'Escape' && onClose()}
    transition:fade={{ duration: 150 }}
    role="presentation"
    tabindex="-1"
  >
    <div
      class="w-full max-w-lg bg-[#11131E]/95 border border-white/[0.12] rounded-3xl shadow-2xl p-6 sm:p-7 flex flex-col gap-5 text-white backdrop-blur-2xl relative max-h-[90vh] overflow-hidden"
      transition:scale={{ duration: 180, start: 0.95 }}
    >
      <!-- Top edge accent -->
      <div class="absolute inset-x-0 top-0 h-[1.5px] bg-gradient-to-r from-transparent via-crimson-400/60 to-transparent"></div>

      <!-- Radiant glow behind header -->
      <div class="absolute -right-12 -top-12 w-48 h-48 bg-crimson-600/20 rounded-full blur-3xl pointer-events-none"></div>

      <!-- Header -->
      <div class="flex items-start justify-between relative z-10">
        <div class="space-y-1">
          <div class="flex items-center gap-2">
            <span class="px-2 py-0.5 rounded-full text-[10px] font-bold font-mono bg-crimson-500/20 text-crimson-400 border border-crimson-500/30">
              {currentRelease.version}
            </span>
            <span class="text-[10px] font-semibold text-zinc-400 uppercase tracking-wider">
              {currentRelease.badge}
            </span>
            <span class="text-zinc-600 text-xs">•</span>
            <span class="text-[10px] text-zinc-400">
              {currentRelease.date}
            </span>
          </div>
          <h2 class="text-lg sm:text-xl font-extrabold text-white tracking-tight">
            {currentRelease.title}
          </h2>
          <p class="text-xs text-zinc-300 leading-relaxed pt-0.5">
            {currentRelease.summary}
          </p>
        </div>

        <button
          class="p-1.5 rounded-xl text-zinc-400 hover:text-white hover:bg-white/[0.08] transition cursor-pointer shrink-0"
          on:click={onClose}
          aria-label={$t('close')}
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <!-- Changelog Sections (Scrollable content) -->
      <div class="space-y-4 overflow-y-auto pr-1 py-1 max-h-[46vh] scrollbar-thin">
        {#each currentRelease.sections as section}
          <div class="p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] space-y-2">
            <div class="flex items-center gap-2 text-xs font-bold text-white">
              <span class="w-4 h-4 shrink-0 {section.color} flex items-center justify-center">
                {#if section.icon === 'desktop'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect>
                    <line x1="8" y1="21" x2="16" y2="21"></line>
                    <line x1="12" y1="17" x2="12" y2="21"></line>
                  </svg>
                {:else if section.icon === 'database'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <ellipse cx="12" cy="5" rx="9" ry="3"></ellipse>
                    <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"></path>
                    <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"></path>
                  </svg>
                {:else if section.icon === 'sliders'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <line x1="4" y1="21" x2="4" y2="14"></line>
                    <line x1="4" y1="10" x2="4" y2="3"></line>
                    <line x1="12" y1="21" x2="12" y2="12"></line>
                    <line x1="12" y1="8" x2="12" y2="3"></line>
                    <line x1="20" y1="21" x2="20" y2="16"></line>
                    <line x1="20" y1="12" x2="20" y2="3"></line>
                    <line x1="1" y1="14" x2="7" y2="14"></line>
                    <line x1="9" y1="8" x2="15" y2="8"></line>
                    <line x1="17" y1="16" x2="23" y2="16"></line>
                  </svg>
                {:else if section.icon === 'shield'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
                  </svg>
                {:else if section.icon === 'rocket'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.59 14.37a6 6 0 01-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 006.16-12.12A14.98 14.98 0 009.63 8.41m5.96 5.96a14.926 14.926 0 01-5.841 2.58m-.119-8.54a6 6 0 00-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 00-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 01-2.448-2.448 14.9 14.9 0 01.06-.312m-2.24 2.39a4.493 4.493 0 00-1.757 4.306 4.493 4.493 0 004.306-1.758M16 8a2 2 0 11-4 0 2 2 0 014 0z"/>
                  </svg>
                {:else}
                  <!-- sparkles default -->
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"/>
                  </svg>
                {/if}
              </span>
              <span class={section.color}>{section.title}</span>
            </div>
            <ul class="space-y-1.5 pl-6 list-disc list-outside text-[11px] text-zinc-300 leading-relaxed">
              {#each section.items as item}
                <li>{item}</li>
              {/each}
            </ul>
          </div>
        {/each}
      </div>

      <!-- Footer CTA -->
      <div class="pt-2 border-t border-white/[0.08] flex items-center justify-between">
        <span class="text-[10px] text-zinc-400">
          Dengerin Desktop • {$t('license_sub')}
        </span>
        <button
          class="px-5 py-2.5 rounded-xl bg-gradient-to-r from-crimson-600 to-crimson-500 hover:from-crimson-500 hover:to-crimson-400 text-white font-semibold text-xs shadow-lg shadow-crimson-500/30 hover:scale-[1.02] active:scale-95 transition-all cursor-pointer"
          on:click={onClose}
        >
          {$t('close')}
        </button>
      </div>
    </div>
  </div>
{/if}
