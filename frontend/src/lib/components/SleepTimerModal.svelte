<script lang="ts">
  import { fade, scale } from 'svelte/transition'
  import { cubicOut } from 'svelte/easing'
  import {
    isSleepTimerModalOpen,
    sleepTimerRemaining,
    sleepTimerMode,
    closeSleepTimerModal,
    setTimerMinutes,
    setEndOfTrackTimer,
    cancelSleepTimer,
  } from '../stores/sleepTimerStore'
  import { t } from '../stores/i18nStore'

  let customMinutes = 20

  function formatCountdown(sec: number | null): string {
    if (sec === null || sec < 0) return '00:00'
    const m = Math.floor(sec / 60)
    const s = sec % 60
    return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  }

  function handleSelectPreset(mins: number) {
    setTimerMinutes(mins)
    closeSleepTimerModal()
  }

  function handleSelectEndOfTrack() {
    setEndOfTrackTimer()
    closeSleepTimerModal()
  }

  function handleCustomStart() {
    if (customMinutes > 0) {
      setTimerMinutes(customMinutes)
      closeSleepTimerModal()
    }
  }

  function handleCancelTimer() {
    cancelSleepTimer()
    closeSleepTimerModal()
  }
</script>

{#if $isSleepTimerModalOpen}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4 select-none"
    transition:fade={{ duration: 150 }}
    on:click={closeSleepTimerModal}
    role="presentation"
  >
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class="w-full max-w-sm bg-[#11131E] border border-white/[0.12] rounded-3xl shadow-2xl p-6 flex flex-col gap-5 overflow-hidden"
      transition:scale={{ start: 0.95, duration: 180, easing: cubicOut }}
      on:click|stopPropagation
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-white/[0.08] pb-3.5">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-xl bg-purple-500/20 text-purple-400 flex items-center justify-center">
            <svg class="w-4 h-4 fill-none stroke-current stroke-2" viewBox="0 0 24 24">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
            </svg>
          </div>
          <h3 class="text-base font-bold text-white">{$t('sleep_timer_title')}</h3>
        </div>
        <button
          class="text-zinc-400 hover:text-white p-1 rounded-lg hover:bg-white/[0.06] transition cursor-pointer"
          on:click={closeSleepTimerModal}
          aria-label={$t('close')}
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <!-- Active Timer Status Banner -->
      {#if $sleepTimerRemaining !== null}
        <div class="p-4 rounded-2xl bg-purple-500/10 border border-purple-500/20 flex flex-col items-center gap-2">
          <span class="text-[10px] uppercase font-bold tracking-wider text-purple-400">
            {$sleepTimerMode === 'end_of_track' ? $t('sleep_timer_end_of_track') : $t('sleep_timer_time_left')}
          </span>
          <div class="text-3xl font-mono font-extrabold text-white tracking-wider">
            {formatCountdown($sleepTimerRemaining)}
          </div>
          <button
            class="mt-1 px-4 py-1.5 rounded-xl bg-crimson-500/20 hover:bg-crimson-500/30 text-crimson-400 text-xs font-semibold transition cursor-pointer"
            on:click={handleCancelTimer}
          >
            {$t('sleep_timer_turn_off')}
          </button>
        </div>
      {/if}

      <!-- Presets Grid -->
      <div class="space-y-2">
        <span class="text-xs font-semibold text-zinc-400">{$t('sleep_timer_quick_duration')}</span>
        <div class="grid grid-cols-2 gap-2">
          <button
            class="p-2.5 rounded-xl bg-white/[0.04] hover:bg-white/[0.08] text-xs font-semibold text-zinc-200 transition cursor-pointer border border-white/[0.06] text-center"
            on:click={() => handleSelectPreset(15)}
          >
            15 {$t('sleep_timer_minutes')}
          </button>
          <button
            class="p-2.5 rounded-xl bg-white/[0.04] hover:bg-white/[0.08] text-xs font-semibold text-zinc-200 transition cursor-pointer border border-white/[0.06] text-center"
            on:click={() => handleSelectPreset(30)}
          >
            30 {$t('sleep_timer_minutes')}
          </button>
          <button
            class="p-2.5 rounded-xl bg-white/[0.04] hover:bg-white/[0.08] text-xs font-semibold text-zinc-200 transition cursor-pointer border border-white/[0.06] text-center"
            on:click={() => handleSelectPreset(45)}
          >
            45 {$t('sleep_timer_minutes')}
          </button>
          <button
            class="p-2.5 rounded-xl bg-white/[0.04] hover:bg-white/[0.08] text-xs font-semibold text-zinc-200 transition cursor-pointer border border-white/[0.06] text-center"
            on:click={() => handleSelectPreset(60)}
          >
            60 {$t('sleep_timer_minutes')}
          </button>
        </div>

        <button
          class="w-full p-2.5 rounded-xl bg-purple-500/15 hover:bg-purple-500/25 text-purple-300 text-xs font-semibold transition cursor-pointer border border-purple-500/30 flex items-center justify-center gap-2"
          on:click={handleSelectEndOfTrack}
        >
          <svg class="w-3.5 h-3.5 fill-none stroke-current stroke-2" viewBox="0 0 24 24">
            <polygon points="5 4 15 12 5 20 5 4"></polygon>
            <line x1="19" y1="5" x2="19" y2="19"></line>
          </svg>
          <span>{$t('sleep_timer_end_of_track')}</span>
        </button>
      </div>

      <!-- Custom Minutes Input -->
      <div class="space-y-1.5 border-t border-white/[0.08] pt-3">
        <label for="custom-sleep-timer" class="text-xs font-semibold text-zinc-400">{$t('sleep_timer_custom_label')}</label>
        <div class="flex items-center gap-2">
          <input
            id="custom-sleep-timer"
            type="number"
            min="1"
            max="300"
            bind:value={customMinutes}
            class="w-24 bg-[#181A24] border border-white/[0.1] rounded-xl px-3 py-2 text-xs font-mono text-white text-center focus:outline-none focus:border-purple-500/60"
          />
          <button
            class="flex-1 py-2 px-3 rounded-xl bg-white/[0.08] hover:bg-white/[0.12] text-xs font-semibold text-white transition cursor-pointer"
            on:click={handleCustomStart}
          >
            {$t('sleep_timer_start_btn')}
          </button>
        </div>
      </div>

      <p class="text-[11px] text-zinc-400 text-center leading-relaxed">
        {$t('sleep_timer_fade_note')}
      </p>
    </div>
  </div>
{/if}
