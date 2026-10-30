<script lang="ts">
  import { onMount } from 'svelte'
  import { networkStore } from '../stores/networkStore'
  import { audioQualityStore, showToast } from '../stores/settingsStore'

  let showReconnectedPill = false
  let reconnectedTimer: any = null
  let dismissedWeak = false

  onMount(() => {
    return networkStore.onReconnect(() => {
      dismissedWeak = false
      showReconnectedPill = true
      clearTimeout(reconnectedTimer)
      reconnectedTimer = setTimeout(() => {
        showReconnectedPill = false
      }, 3500)
    })
  })

  function handleCheckNow() {
    networkStore.checkNow()
    showToast('Memeriksa konektivitas jaringan...')
  }

  function handleEnableSaverMode() {
    audioQualityStore.set('saver')
    dismissedWeak = true
    showToast('Mode Hemat Data diaktifkan (~70-128kbps) untuk koneksi lambat')
  }
</script>

{#if !$networkStore.isOnline}
  <!-- Offline Global Banner -->
  <aside
    class="w-full bg-crimson-950/90 border-b border-crimson-500/40 text-white px-4 py-2 flex items-center justify-between gap-3 text-xs z-50 backdrop-blur-md shadow-md shadow-crimson-950/40 transition-all duration-300"
    role="alert"
  >
    <div class="flex items-center gap-2.5 min-w-0">
      <div class="w-2.5 h-2.5 rounded-full bg-crimson-500 animate-ping shrink-0"></div>
      <svg class="w-4 h-4 text-crimson-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="1" y1="1" x2="23" y2="23" />
        <path d="M16.72 11.06A10.94 10.94 0 0 1 19 12.55" />
        <path d="M5 12.55a10.94 10.94 0 0 1 5.17-2.39" />
        <path d="M10.71 5.05A16 16 0 0 1 22.58 9" />
        <path d="M1.42 9a15.91 15.91 0 0 1 4.7-2.88" />
        <path d="M8.53 16.11a6 6 0 0 1 6.95 0" />
        <line x1="12" y1="20" x2="12.01" y2="20" />
      </svg>
      <div class="truncate">
        <span class="font-bold text-crimson-200">Perangkat Offline:</span>
        <span class="text-zinc-300 ml-1">Koneksi internet terputus. Pemutaran akan dilanjutkan otomatis saat jaringan kembali.</span>
      </div>
    </div>
    <div class="flex items-center gap-2 shrink-0">
      <button
        class="px-2.5 py-1 rounded-lg bg-crimson-600 hover:bg-crimson-500 text-white font-semibold transition active:scale-95 cursor-pointer shadow-sm text-[11px] flex items-center gap-1"
        on:click={handleCheckNow}
        title="Periksa ulang koneksi jaringan"
      >
        <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        Coba Hubungkan
      </button>
    </div>
  </aside>
{:else if showReconnectedPill}
  <!-- Reconnected Brief Pill -->
  <aside
    class="w-full bg-emerald-950/90 border-b border-emerald-500/40 text-emerald-200 px-4 py-1.5 flex items-center justify-between gap-3 text-xs z-50 backdrop-blur-md shadow-md shadow-emerald-950/30 transition-all duration-300"
    role="status"
  >
    <div class="flex items-center gap-2 min-w-0">
      <svg class="w-4 h-4 text-emerald-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M5 13l4 4L19 7" />
      </svg>
      <span class="font-semibold text-emerald-300">Koneksi Pulih:</span>
      <span class="text-emerald-100/90 truncate">Perangkat kembali terhubung ke internet {$networkStore.latencyMs ? `(~${$networkStore.latencyMs}ms)` : ''}.</span>
    </div>
    <button
      class="text-emerald-400 hover:text-white transition p-1"
      on:click={() => showReconnectedPill = false}
      title="Tutup"
    >
      <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>
    </button>
  </aside>
{:else if $networkStore.isWeak && !dismissedWeak}
  <!-- Weak Network Advisory Banner -->
  <aside
    class="w-full bg-amber-950/80 border-b border-amber-500/30 text-amber-200 px-4 py-1.5 flex items-center justify-between gap-3 text-xs z-40 backdrop-blur-md transition-all duration-300"
    role="status"
  >
    <div class="flex items-center gap-2 min-w-0">
      <svg class="w-4 h-4 text-amber-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <span class="font-semibold text-amber-300">Jaringan Lemah:</span>
      <span class="text-amber-100/80 truncate">Latensi jaringan tinggi {$networkStore.latencyMs ? `(~${$networkStore.latencyMs}ms)` : ''}. Buffering mungkin melambat.</span>
    </div>
    <div class="flex items-center gap-2 shrink-0">
      {#if $audioQualityStore !== 'saver'}
        <button
          class="px-2 py-0.5 rounded bg-amber-600/80 hover:bg-amber-500 text-white font-medium text-[11px] transition cursor-pointer"
          on:click={handleEnableSaverMode}
          title="Beralih ke format audio yang lebih hemat kuota"
        >
          Mode Hemat
        </button>
      {/if}
      <button
        class="text-amber-400 hover:text-white transition p-1"
        on:click={() => dismissedWeak = true}
        title="Tutup peringatan"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>
      </button>
    </div>
  </aside>
{/if}
