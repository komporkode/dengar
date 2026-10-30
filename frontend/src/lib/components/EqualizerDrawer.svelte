<script lang="ts">
  import { equalizerStore, EQ_BANDS, EQ_PRESETS } from '../stores/equalizerStore'
  import { fade, scale } from 'svelte/transition'

  export let isOpen = false
  export let onClose: () => void

  $: eq = $equalizerStore

  function handleBackdrop(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      onClose()
    }
  }

  function formatGain(val: number): string {
    if (val > 0) return `+${val}dB`
    if (val < 0) return `${val}dB`
    return '0dB'
  }
</script>

{#if isOpen}
  <!-- Backdrop click catcher -->
  <div
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-3 sm:p-4 bg-black/60 backdrop-blur-xs select-none"
    on:click={handleBackdrop}
    on:keydown={(e) => e.key === 'Escape' && onClose()}
    transition:fade={{ duration: 150 }}
    role="presentation"
    tabindex="-1"
  >
    <div
      class="w-full max-w-md bg-[#11131E]/95 border border-white/[0.12] rounded-3xl shadow-2xl p-5 sm:p-6 flex flex-col gap-5 text-white backdrop-blur-2xl relative"
      transition:scale={{ duration: 180, start: 0.95 }}
    >
      <!-- Top edge accent -->
      <div class="absolute inset-x-0 top-0 h-[1px] bg-gradient-to-r from-transparent via-crimson-400/40 to-transparent"></div>

      <!-- Header -->
      <div class="flex items-center justify-between border-b border-white/[0.08] pb-3">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-xl bg-crimson-500/15 border border-crimson-500/30 flex items-center justify-center text-crimson-400">
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
          </div>
          <div>
            <h3 class="font-bold text-sm text-white">Equalizer Audio</h3>
            <p class="text-[11px] text-zinc-400 font-mono">5-Band Studio Filter</p>
          </div>
        </div>

        <!-- Toggle Switch ON/OFF -->
        <div class="flex items-center gap-3">
          <button
            type="button"
            class="flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-semibold transition-all cursor-pointer {eq.isEnabled ? 'bg-crimson-500 text-white shadow-lg shadow-crimson-500/30' : 'bg-white/[0.08] text-zinc-400 hover:text-white'}"
            on:click={() => equalizerStore.setEnabled(!eq.isEnabled)}
          >
            <span class="w-2 h-2 rounded-full {eq.isEnabled ? 'bg-white animate-pulse' : 'bg-zinc-500'}"></span>
            <span>{eq.isEnabled ? 'Aktif' : 'Bypass'}</span>
          </button>

          <button
            class="p-1.5 rounded-xl text-zinc-400 hover:text-white hover:bg-white/[0.08] transition cursor-pointer"
            on:click={onClose}
            aria-label="Tutup"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
      </div>

      <!-- Presets Carousel / Pills -->
      <div class="space-y-1.5">
        <div class="flex items-center justify-between text-[11px] font-semibold text-zinc-400 uppercase tracking-wider">
          <span>Preset Suara</span>
          <button
            class="text-crimson-400 hover:text-crimson-300 text-[10px] cursor-pointer"
            on:click={() => equalizerStore.resetToFlat()}
          >
            Reset ke Flat
          </button>
        </div>
        <div class="flex items-center gap-1.5 overflow-x-auto pb-1 scrollbar-none">
          {#each Object.keys(EQ_PRESETS) as presetName}
            <button
              class="px-2.5 py-1 rounded-xl text-xs font-medium whitespace-nowrap transition cursor-pointer {eq.isEnabled && eq.preset === presetName ? 'bg-white/20 text-white border border-white/20 shadow-xs' : 'bg-white/[0.04] text-zinc-400 hover:text-white hover:bg-white/[0.08] border border-transparent'}"
              on:click={() => equalizerStore.setPreset(presetName)}
            >
              {presetName}
            </button>
          {/each}
        </div>
      </div>

      <!-- 5-Band Sliders -->
      <div class="grid grid-cols-5 gap-2 sm:gap-3 py-2 px-1 bg-black/30 rounded-2xl border border-white/[0.05]">
        {#each EQ_BANDS as band, idx}
          {@const gain = eq.gains[idx] || 0}
          <div class="flex flex-col items-center gap-2">
            <!-- Current Gain Indicator -->
            <span class="text-[10px] font-mono font-bold {eq.isEnabled ? (gain > 0 ? 'text-crimson-400' : gain < 0 ? 'text-indigo-400' : 'text-zinc-400') : 'text-zinc-600'}">
              {formatGain(gain)}
            </span>

            <!-- Vertical Slider Shell -->
            <div class="relative h-32 sm:h-36 w-8 flex items-center justify-center">
              <!-- Center Zero Tick Mark -->
              <div class="absolute left-1 right-1 top-1/2 -translate-y-1/2 h-[1px] bg-white/20 pointer-events-none"></div>

              <!-- Slider Track -->
              <div class="h-full w-1.5 bg-white/10 rounded-full overflow-hidden relative pointer-events-none">
                <!-- Center to fill visualization -->
                {#if gain > 0}
                  <div
                    class="absolute left-0 right-0 bottom-1/2 bg-gradient-to-t from-crimson-500 to-rose-400 rounded-full transition-all"
                    style="height: {(gain / 12) * 50}%;"
                  ></div>
                {:else if gain < 0}
                  <div
                    class="absolute left-0 right-0 top-1/2 bg-gradient-to-b from-indigo-500 to-sky-400 rounded-full transition-all"
                    style="height: {(Math.abs(gain) / 12) * 50}%;"
                  ></div>
                {/if}
              </div>

              <!-- Native Vertical Input (Rotated) -->
              <input
                type="range"
                min="-12"
                max="12"
                step="0.5"
                value={gain}
                class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                style="writing-mode: vertical-lr; direction: rtl;"
                disabled={!eq.isEnabled}
                on:input={(e) => equalizerStore.setBandGain(idx, parseFloat(e.currentTarget.value))}
              />
            </div>

            <!-- Frequency Label -->
            <span class="text-[10px] font-mono font-semibold text-zinc-400 text-center whitespace-nowrap">
              {band.label}
            </span>
          </div>
        {/each}
      </div>

      <!-- Footer Info -->
      <div class="flex items-center justify-between text-[11px] text-zinc-400 pt-1">
        <span class="flex items-center gap-1.5">
          <span class="w-1.5 h-1.5 rounded-full {eq.isEnabled ? 'bg-emerald-400' : 'bg-zinc-600'}"></span>
          <span>{eq.isEnabled ? 'Pemrosesan audio aktif' : 'Equalizer nonaktif'}</span>
        </span>
        <button
          class="text-xs font-semibold px-4 py-1.5 rounded-xl bg-white/[0.08] hover:bg-white/[0.14] text-white transition cursor-pointer"
          on:click={onClose}
        >
          Selesai
        </button>
      </div>
    </div>
  </div>
{/if}
