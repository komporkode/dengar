<script lang="ts">
  import { CHANGELOG_DATA } from '../data/changelogs'
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
          aria-label="Tutup"
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
              <span>{section.icon}</span>
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
          Dengerin Desktop • Sumber Terbuka
        </span>
        <button
          class="px-5 py-2.5 rounded-xl bg-gradient-to-r from-crimson-600 to-crimson-500 hover:from-crimson-500 hover:to-crimson-400 text-white font-semibold text-xs shadow-lg shadow-crimson-500/30 hover:scale-[1.02] active:scale-95 transition-all cursor-pointer"
          on:click={onClose}
        >
          Mulai Mendengarkan
        </button>
      </div>
    </div>
  </div>
{/if}
