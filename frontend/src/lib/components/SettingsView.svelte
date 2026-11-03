<script lang="ts">
  import { onMount, onDestroy, createEventDispatcher } from "svelte";
  import {
    isSettingsOpen,
    isWhatsNewOpen,
    audioQualityStore,
    showToast,
    closeActionStore,
    isDiscordRpcEnabled,
    crossfadeDurationStore,
    type CloseAction,
  } from "../stores/settingsStore";
  import { currentLocale, t } from "../stores/i18nStore";
  import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime.js";
  import {
    GetAudioCacheStats,
    ClearAudioCache,
    SetAudioCacheLimit,
  } from "../../../wailsjs/go/main/App";
  import type { backend } from "../../../wailsjs/go/models";
  import {
    updateState,
    isCheckingUpdate,
    isDownloadingUpdate,
    downloadedFilePath,
    checkForUpdatesManual,
    startDownloadUpdate,
    revealDownloadedFile,
    installAndRestartUpdate,
  } from "../stores/updateStore";

  const dispatch = createEventDispatcher<{ close: void }>();

  export let initialTab: "general" | "audio" | "about" = "general";
  let activeTab: "general" | "audio" | "about" = initialTab;

  let cacheStats: backend.AudioCacheStats | null = null;
  let isClearingCache = false;
  let selectedLimitMb: number = 1024;

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") {
      closeSettings();
    }
  }

  function closeSettings() {
    isSettingsOpen.set(false);
    dispatch("close");
  }

  async function loadCacheStats() {
    try {
      const stats = await GetAudioCacheStats();
      if (stats) {
        cacheStats = stats;
        selectedLimitMb =
          stats.maxBytes > 0 ? Math.round(stats.maxBytes / (1024 * 1024)) : 0;
      }
    } catch (e) {
      console.warn("Failed to load cache stats:", e);
    }
  }

  function formatBytes(bytes: number): string {
    if (!bytes || bytes <= 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
  }

  async function handleClearCache() {
    if (isClearingCache) return;
    isClearingCache = true;
    try {
      await ClearAudioCache();
      await loadCacheStats();
      showToast("Seluruh berkas cache audio offline berhasil dibersihkan!");
    } catch (e) {
      console.warn("Clear cache error:", e);
      showToast("Gagal membersihkan cache");
    } finally {
      isClearingCache = false;
    }
  }

  async function handleLimitChange(limitMb: number) {
    selectedLimitMb = limitMb;
    const bytes = limitMb > 0 ? limitMb * 1024 * 1024 : 0;
    try {
      await SetAudioCacheLimit(bytes);
      await loadCacheStats();
      showToast(
        limitMb > 0
          ? `Batas kuota cache diatur ke ${limitMb >= 1024 ? limitMb / 1024 + " GB" : limitMb + " MB"}`
          : "Batas kuota cache dinonaktifkan (Tanpa Batas)",
      );
    } catch (e) {
      console.warn("Set limit error:", e);
    }
  }

  function openGithub() {
    try {
      BrowserOpenURL("https://github.com/komporkode/dengerin");
    } catch (e) {
      console.warn("Failed to open github:", e);
    }
  }

  $: if (activeTab === "audio") {
    loadCacheStats();
  }

  onMount(() => {
    window.addEventListener("keydown", handleKeydown);
    if (activeTab === "audio") {
      loadCacheStats();
    }
  });

  onDestroy(() => {
    window.removeEventListener("keydown", handleKeydown);
  });
</script>

<div class="flex-1 w-full flex flex-col overflow-y-auto px-6 lg:px-12 py-8 max-w-5xl mx-auto space-y-6">
  <!-- Top Banner / Navigation Header (Omnitik Fullscreen Tab Header) -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-white/[0.08]">
    <div class="flex items-center gap-3.5">
      <div class="w-11 h-11 rounded-2xl bg-gradient-to-tr from-crimson-600 to-rose-500 flex items-center justify-center shadow-lg shadow-crimson-500/25 shrink-0">
        <svg class="w-5 h-5 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"></circle>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
        </svg>
      </div>
      <div>
        <h1 class="text-xl font-bold text-white tracking-tight flex items-center gap-2">
          <span>{$t("settings_title")}</span>
          <span class="text-xs px-2 py-0.5 rounded-full bg-crimson-500/20 text-crimson-300 border border-crimson-500/30 font-mono font-medium">v1.2.0</span>
        </h1>
        <p class="text-xs text-zinc-400 mt-0.5">
          Kelola preferensi antarmuka, audio engine, kapasitas cache disk, dan info aplikasi.
        </p>
      </div>
    </div>

    <!-- Return to Music Button -->
    <button
      class="px-4 py-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-xs font-semibold text-zinc-200 hover:text-white transition cursor-pointer flex items-center gap-2 shadow-xs shrink-0 active:scale-95 self-start sm:self-center"
      on:click={closeSettings}
      title="Kembali ke tampilan musik (Esc)"
    >
      <svg class="w-4 h-4 text-zinc-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M19 12H5M12 19l-7-7 7-7"/>
      </svg>
      <span>Kembali ke Musik</span>
      <kbd class="px-1.5 py-0.5 rounded bg-black/40 text-[10px] text-zinc-400 font-mono border border-white/10">Esc</kbd>
    </button>
  </div>

  <!-- Omnitik-style Menu Tabs Navigation (Segmented pill bar) -->
  <div class="flex items-center gap-2 p-1.5 rounded-2xl bg-black/30 border border-white/[0.06] backdrop-blur-md">
    <button
      class="flex-1 flex items-center justify-center gap-2 py-2.5 px-4 rounded-xl text-xs font-semibold transition-all cursor-pointer {activeTab === 'general'
        ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25'
        : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
      on:click={() => (activeTab = "general")}
    >
      <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="2" y1="12" x2="22" y2="12"></line>
        <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
      </svg>
      <span>{$t("tab_general")} & Antarmuka</span>
    </button>

    <button
      class="flex-1 flex items-center justify-center gap-2 py-2.5 px-4 rounded-xl text-xs font-semibold transition-all cursor-pointer {activeTab === 'audio'
        ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25'
        : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
      on:click={() => (activeTab = "audio")}
    >
      <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
        <path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"></path>
      </svg>
      <span>{$t("tab_audio")} & Cache Disk</span>
    </button>

    <button
      class="flex-1 flex items-center justify-center gap-2 py-2.5 px-4 rounded-xl text-xs font-semibold transition-all cursor-pointer {activeTab === 'about'
        ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25'
        : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
      on:click={() => (activeTab = "about")}
    >
      <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="16" x2="12" y2="12"></line>
        <line x1="12" y1="8" x2="12.01" y2="8"></line>
      </svg>
      <span>{$t("tab_about")} & Pembaruan</span>
    </button>
  </div>

  <!-- Content Sections -->
  {#if activeTab === "general"}
    <!-- TAB 1: General & Interface -->
    <div class="space-y-6 animate-fade-in">
      <!-- Language Selection Card -->
      <div class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4">
        <div>
          <h3 class="text-sm font-bold text-white">Bahasa Antarmuka (Language)</h3>
          <p class="text-xs text-zinc-400 mt-0.5">Pilih bahasa pengantar untuk navigasi dan deskripsi fitur aplikasi.</p>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
          <button
            class="flex items-center justify-between p-4 rounded-2xl border transition-all cursor-pointer {$currentLocale === 'id'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => currentLocale.set("id")}
          >
            <div class="flex items-center gap-3">
              <span class="w-8 h-8 rounded-xl bg-white/[0.08] flex items-center justify-center font-bold text-xs font-mono">ID</span>
              <div class="text-left">
                <p class="text-xs font-bold">{$t("lang_id")}</p>
                <p class="text-[11px] text-zinc-400">Bahasa Baku Indonesia</p>
              </div>
            </div>
            {#if $currentLocale === "id"}
              <div class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0">
                <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </div>
            {/if}
          </button>

          <button
            class="flex items-center justify-between p-4 rounded-2xl border transition-all cursor-pointer {$currentLocale === 'en'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => currentLocale.set("en")}
          >
            <div class="flex items-center gap-3">
              <span class="w-8 h-8 rounded-xl bg-white/[0.08] flex items-center justify-center font-bold text-xs font-mono">EN</span>
              <div class="text-left">
                <p class="text-xs font-bold">{$t("lang_en")}</p>
                <p class="text-[11px] text-zinc-400">International English</p>
              </div>
            </div>
            {#if $currentLocale === "en"}
              <div class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0">
                <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </div>
            {/if}
          </button>
        </div>
      </div>

      <!-- Close Action Behavior Card -->
      <div class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4">
        <div>
          <h3 class="text-sm font-bold text-white">Tindakan Tombol Tutup ('X')</h3>
          <p class="text-xs text-zinc-400 mt-0.5">Tentukan perilaku aplikasi saat tombol silang (Close) pada title bar jendela ditekan.</p>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <button
            class="p-4 rounded-2xl border text-left transition-all cursor-pointer {$closeActionStore === 'ask'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => {
              closeActionStore.set('ask');
              showToast('Tombol Tutup: Tanya Setiap Kali');
            }}
          >
            <div class="text-xs font-bold">Tanya Selalu</div>
            <p class="text-[11px] text-zinc-400 mt-1">Tampilkan dialog konfirmasi pilihan saat tombol 'X' ditekan.</p>
          </button>

          <button
            class="p-4 rounded-2xl border text-left transition-all cursor-pointer {$closeActionStore === 'tray'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => {
              closeActionStore.set('tray');
              showToast('Tombol Tutup: Sembunyi ke Background');
            }}
          >
            <div class="text-xs font-bold">System Tray</div>
            <p class="text-[11px] text-zinc-400 mt-1">Sembunyikan jendela ke tray bar dan lagu tetap berputar di latar.</p>
          </button>

          <button
            class="p-4 rounded-2xl border text-left transition-all cursor-pointer {$closeActionStore === 'quit'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => {
              closeActionStore.set('quit');
              showToast('Tombol Tutup: Langsung Keluar');
            }}
          >
            <div class="text-xs font-bold">Keluar Sepenuhnya</div>
            <p class="text-[11px] text-zinc-400 mt-1">Hentikan proses aplikasi Dengerin dan keluar secara tuntas.</p>
          </button>
        </div>
      </div>

      <!-- Discord RPC Card -->
      <div class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] flex items-center justify-between gap-4">
        <div class="flex items-center gap-3.5">
          <div class="w-10 h-10 rounded-2xl bg-[#5865F2]/20 text-[#5865F2] flex items-center justify-center font-bold text-lg shrink-0">
            🎮
          </div>
          <div>
            <h3 class="text-sm font-bold text-white">Discord Rich Presence (RPC)</h3>
            <p class="text-xs text-zinc-400 mt-0.5">Tampilkan status lagu, judul, artis, dan durasi yang sedang diputar di profil akun Discord Anda.</p>
          </div>
        </div>
        <button
          type="button"
          on:click={() => {
            const next = !$isDiscordRpcEnabled;
            isDiscordRpcEnabled.set(next);
            showToast(next ? 'Discord Rich Presence Diaktifkan' : 'Discord Rich Presence Dinonaktifkan');
          }}
          class="w-12 h-6 rounded-full transition-colors relative cursor-pointer shrink-0 {$isDiscordRpcEnabled ? 'bg-crimson-600' : 'bg-white/10'}"
          aria-label="Toggle Discord RPC"
        >
          <span
            class="absolute top-1 w-4 h-4 rounded-full bg-white transition-transform {$isDiscordRpcEnabled ? 'left-7' : 'left-1'}"
          ></span>
        </button>
      </div>

      <!-- Keyboard Shortcuts Reference Card -->
      <div class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4">
        <div class="flex items-center gap-2 text-white text-sm font-bold">
          <svg class="w-4 h-4 text-crimson-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <rect x="2" y="4" width="20" height="16" rx="2" ry="2"></rect>
            <path d="M6 8h.001M10 8h.001M14 8h.001M18 8h.001M6 12h.001M18 12h.001M8 16h8"></path>
          </svg>
          <span>Pintasan Keyboard Global (Desktop Shortcuts)</span>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3 text-xs">
          <div class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]">
            <span class="text-zinc-400">Putar / Jeda Lagu</span>
            <kbd class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]">Spasi</kbd>
          </div>
          <div class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]">
            <span class="text-zinc-400">Naikkan / Turunkan Volume</span>
            <kbd class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]">↑ / ↓</kbd>
          </div>
          <div class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]">
            <span class="text-zinc-400">Maju / Mundur 5 Detik</span>
            <kbd class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]">← / →</kbd>
          </div>
          <div class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]">
            <span class="text-zinc-400">Lagu Sebelum / Sesudah</span>
            <kbd class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]">J / L</kbd>
          </div>
          <div class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]">
            <span class="text-zinc-400">Bisukan Suara (Mute)</span>
            <kbd class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]">M</kbd>
          </div>
          <div class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]">
            <span class="text-zinc-400">Kembali ke Musik</span>
            <kbd class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]">Esc</kbd>
          </div>
        </div>
      </div>
    </div>
  {:else if activeTab === "audio"}
    <!-- TAB 2: Audio & Disk Cache -->
    <div class="space-y-6 animate-fade-in">
      <!-- Audio Streaming Resolution -->
      <div class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4">
        <div>
          <h3 class="text-sm font-bold text-white">{$t("quality_setting")}</h3>
          <p class="text-xs text-zinc-400 mt-0.5">Pilih resolusi format kompresi audio yang diambil dari server YouTube.</p>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
          <button
            class="p-4 rounded-2xl border text-left transition-all cursor-pointer flex items-center justify-between {$audioQualityStore === 'high'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => {
              audioQualityStore.set("high");
              showToast("Kualitas audio diatur ke Kualitas Tertinggi (Opus 160kbps / AAC)");
            }}
          >
            <div class="flex items-center gap-3">
              <span class="px-2 py-1 rounded-xl bg-amber-500/20 text-amber-300 font-mono text-xs font-bold border border-amber-500/30">HQ</span>
              <div>
                <div class="font-bold text-xs leading-tight">{$t("quality_high")}</div>
                <div class="text-[11px] text-zinc-400 mt-0.5">Opus 160kbps / AAC • Resolusi Master Hi-Fi</div>
              </div>
            </div>
            {#if $audioQualityStore === "high"}
              <div class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0">
                <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </div>
            {/if}
          </button>

          <button
            class="p-4 rounded-2xl border text-left transition-all cursor-pointer flex items-center justify-between {$audioQualityStore === 'saver'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => {
              audioQualityStore.set("saver");
              showToast("Kualitas audio diatur ke Hemat Kuota (ECO ~70-128kbps)");
            }}
          >
            <div class="flex items-center gap-3">
              <span class="px-2 py-1 rounded-xl bg-emerald-500/20 text-emerald-300 font-mono text-xs font-bold border border-emerald-500/30">ECO</span>
              <div>
                <div class="font-bold text-xs leading-tight">{$t("quality_saver")}</div>
                <div class="text-[11px] text-zinc-400 mt-0.5">Hemat Kuota & Buffer Cepat (~70-128kbps)</div>
              </div>
            </div>
            {#if $audioQualityStore === "saver"}
              <div class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0">
                <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </div>
            {/if}
          </button>
        </div>
      </div>

      <!-- Web Audio Crossfade Slider -->
      <div class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-sm font-bold text-white">Transisi Lagu (Crossfade Audio)</h3>
            <p class="text-xs text-zinc-400 mt-0.5">Peralihan lembut volume antara lagu aktif dengan lagu berikutnya agar tidak hening mendadak.</p>
          </div>
          <span class="text-xs font-mono font-bold text-crimson-400 bg-crimson-500/15 border border-crimson-500/25 px-2.5 py-1 rounded-lg">
            {$crossfadeDurationStore > 0 ? `${$crossfadeDurationStore} detik` : 'Nonaktif (0s)'}
          </span>
        </div>
        <div class="flex items-center gap-4 pt-1">
          <span class="text-xs text-zinc-500 font-mono">0s</span>
          <input
            type="range"
            min="0"
            max="10"
            step="1"
            value={$crossfadeDurationStore}
            on:input={(e) => {
              const val = Number((e.target as HTMLInputElement).value);
              crossfadeDurationStore.set(val);
            }}
            class="w-full accent-crimson-500 cursor-pointer h-2 bg-white/10 rounded-lg appearance-none"
          />
          <span class="text-xs text-zinc-500 font-mono">10s</span>
        </div>
      </div>

      <!-- Smart Disk Cache & Offline Storage Management -->
      <div class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4">
        <div class="flex items-center justify-between gap-4">
          <div>
            <div class="flex items-center gap-2.5">
              <h3 class="text-sm font-bold text-white">Smart Disk Cache & Mode Offline</h3>
              <span class="text-xs font-mono px-2 py-0.5 rounded-full bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 font-semibold">
                {formatBytes(cacheStats?.totalBytes || 0)}
              </span>
            </div>
            <p class="text-xs text-zinc-400 mt-1">
              Lagu yang Anda putar otomatis tersimpan ke penyimpanan lokal ({cacheStats?.totalFiles || 0} lagu siap putar offline tanpa kuota internet).
            </p>
          </div>

          <button
            class="px-4 py-2 rounded-xl bg-crimson-600/20 hover:bg-crimson-600 border border-crimson-500/30 text-xs font-semibold text-crimson-200 hover:text-white transition cursor-pointer flex items-center gap-2 shrink-0 active:scale-95 disabled:opacity-50"
            on:click={handleClearCache}
            disabled={isClearingCache || !cacheStats?.totalFiles}
            title="Bersihkan semua cache audio dari disk"
          >
            {#if isClearingCache}
              <svg class="w-3.5 h-3.5 animate-spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10" stroke-opacity="0.25" />
                <path d="M12 2a10 10 0 0 1 10 10" />
              </svg>
              <span>Membersihkan...</span>
            {:else}
              <svg class="w-3.5 h-3.5 text-crimson-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
              </svg>
              <span>Bersihkan Cache</span>
            {/if}
          </button>
        </div>

        <div class="pt-4 border-t border-white/[0.06] flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs">
          <span class="text-zinc-400">Batas Kuota Kapasitas Disk (LRU Automatic Eviction):</span>
          <div class="flex items-center gap-2">
            {#each [500, 1024, 2048, 0] as limit}
              <button
                class="px-3 py-1 rounded-xl text-xs font-mono transition cursor-pointer {selectedLimitMb === limit ? 'bg-crimson-500 text-white font-bold shadow-xs' : 'text-zinc-400 hover:text-zinc-200 bg-white/[0.04]'}"
                on:click={() => handleLimitChange(limit)}
              >
                {limit === 0 ? 'Unlimited' : limit >= 1024 ? (limit/1024) + ' GB' : limit + ' MB'}
              </button>
            {/each}
          </div>
        </div>
      </div>
    </div>
  {:else if activeTab === "about"}
    <!-- TAB 3: About & Updates -->
    <div class="space-y-6 animate-fade-in">
      <!-- Brand Card -->
      <div class="p-6 rounded-2xl bg-gradient-to-br from-[#161824] to-[#10121C] border border-white/10 flex flex-col md:flex-row items-start md:items-center gap-5 shadow-sm">
        <div class="w-16 h-16 rounded-2xl bg-gradient-to-tr from-crimson-600 via-crimson-500 to-rose-400 flex items-center justify-center shadow-lg shadow-crimson-500/25 shrink-0">
          <svg class="w-8 h-8 text-white" viewBox="0 0 269.22 267.55" fill="none">
            <path
              fill="currentColor"
              d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"
            />
          </svg>
        </div>
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2.5 flex-wrap">
            <h2 class="text-base font-extrabold text-white">Dengerin Desktop</h2>
            <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-crimson-500/20 border border-crimson-500/30 text-crimson-400 font-mono">v1.2.0</span>
            <span class="text-[10px] font-medium px-2 py-0.5 rounded-full {$updateState?.isInstalled ? 'bg-sky-500/20 border-sky-500/30 text-sky-400' : 'bg-emerald-500/20 border-emerald-500/30 text-emerald-400'} border font-mono">
              {$updateState?.isInstalled ? '📦 Terpasang (Setup.exe)' : '🧰 Portabel (.zip)'}
            </span>
          </div>
          <p class="text-xs text-zinc-300 leading-relaxed mt-1.5">
            {$t("about_desc")}
          </p>
          <div class="flex items-center gap-2 mt-2 text-[10px] font-medium text-zinc-300 flex-wrap">
            <span class="px-2 py-0.5 rounded bg-white/[0.06] border border-white/10">🚫 Zero Ads</span>
            <span class="px-2 py-0.5 rounded bg-white/[0.06] border border-white/10">💰 100% Free / Rp0</span>
            <span class="px-2 py-0.5 rounded bg-white/[0.06] border border-white/10">🪶 Ultra Low RAM</span>
            <span class="px-2 py-0.5 rounded bg-white/[0.06] border border-white/10">⚖️ GNU GPLv3</span>
          </div>
          <div class="mt-3.5">
            <button
              type="button"
              class="px-3 py-1.5 rounded-xl bg-crimson-500/15 hover:bg-crimson-500/25 border border-crimson-500/30 text-xs font-semibold text-crimson-300 hover:text-white transition cursor-pointer flex items-center gap-1.5"
              on:click={() => isWhatsNewOpen.set(true)}
            >
              <span>✨ Apa yang Baru di v1.2.0</span>
            </button>
          </div>
        </div>
      </div>

      <!-- GitHub & License Card -->
      <div class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-xl bg-white/10 flex items-center justify-center text-white shrink-0">
              <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
                <path fill-rule="evenodd" clip-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.53 1.032 1.53 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
              </svg>
            </div>
            <div>
              <h4 class="text-sm font-bold text-white">GitHub & Komporkode</h4>
              <p class="text-xs text-zinc-400 font-mono">GNU GPLv3 Open Source Repository</p>
            </div>
          </div>

          <button
            class="px-3.5 py-1.5 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-xs font-semibold text-zinc-200 transition cursor-pointer flex items-center gap-2"
            on:click={openGithub}
            title="Buka Repositori di GitHub"
          >
            <span>Kunjungi Repositori</span>
            <svg class="w-3.5 h-3.5 text-zinc-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
              <polyline points="15 3 21 3 21 9"></polyline>
              <line x1="10" y1="14" x2="21" y2="3"></line>
            </svg>
          </button>
        </div>

        <!-- Update Status Section -->
        <div class="pt-4 border-t border-white/[0.06] space-y-3">
          {#if $isCheckingUpdate}
            <div class="flex items-center gap-2 py-1">
              <div class="w-4 h-4 border-2 border-white/20 border-t-crimson-500 rounded-full animate-spin"></div>
              <span class="text-zinc-400 text-xs">Memeriksa pembaruan di server GitHub...</span>
            </div>
          {:else if $updateState?.hasUpdate}
            <div class="p-4 rounded-xl bg-crimson-500/10 border border-crimson-500/25 space-y-3">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2 flex-wrap">
                  <span class="w-2.5 h-2.5 rounded-full bg-crimson-400 shadow-[0_0_8px_#E61E38] animate-pulse"></span>
                  <span class="text-crimson-300 font-bold text-xs">Pembaruan Tersedia:</span>
                  <span class="px-2 py-0.5 rounded bg-crimson-500/25 text-white font-mono font-bold text-xs">
                    {$updateState.latestVersion}
                  </span>
                  <span class="px-2 py-0.5 rounded {$updateState.isInstalled ? 'bg-sky-500/25 text-sky-200 border border-sky-500/30' : 'bg-emerald-500/25 text-emerald-200 border border-emerald-500/30'} font-mono text-[10px]">
                    {$updateState.isInstalled ? 'Installer .exe' : 'Portable .zip'}
                  </span>
                </div>
              </div>

              <p class="text-xs text-zinc-300 leading-relaxed">
                {#if $downloadedFilePath}
                  Berkas telah siap di folder Downloads. Klik tombol di bawah untuk memasang dan me-restart aplikasi.
                {:else if $isDownloadingUpdate}
                  Sedang mengunduh berkas biner pembaruan resmi dari GitHub...
                {:else}
                  Versi {$updateState.latestVersion} siap diunduh otomatis untuk perangkat Anda.
                {/if}
              </p>

              <div class="flex items-center gap-2.5 pt-1">
                {#if $downloadedFilePath}
                  <button
                    class="flex-1 py-2 px-4 rounded-xl bg-crimson-600 hover:bg-crimson-500 text-white font-bold text-xs transition cursor-pointer shadow-md flex items-center justify-center gap-2"
                    on:click={installAndRestartUpdate}
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                      <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
                    </svg>
                    <span>⚡ Pasang & Mulai Ulang Aplikasi</span>
                  </button>
                  <button
                    class="py-2 px-4 rounded-xl bg-white/[0.08] hover:bg-white/[0.14] text-zinc-300 hover:text-white font-semibold text-xs transition cursor-pointer shadow-sm flex items-center gap-1.5 shrink-0"
                    on:click={revealDownloadedFile}
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                      <path d="M5 19a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2h4l2 2h7a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5z"></path>
                    </svg>
                    <span>Buka Folder</span>
                  </button>
                {:else if $isDownloadingUpdate}
                  <div class="flex items-center gap-2.5 text-zinc-300 text-xs py-1">
                    <div class="w-4 h-4 border-2 border-white/20 border-t-crimson-500 rounded-full animate-spin"></div>
                    <span>Mengunduh berkas pembaruan...</span>
                  </div>
                {:else}
                  <button
                    class="w-full py-2 px-4 rounded-xl bg-crimson-500 hover:bg-crimson-600 text-white font-bold text-xs transition cursor-pointer shadow-md shadow-crimson-500/30 flex items-center justify-center gap-2"
                    on:click={() => startDownloadUpdate()}
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                      <path d="M12 4v12m0 0l-4-4m4 4l4-4M4 20h16" />
                    </svg>
                    <span>Unduh Pembaruan ({$updateState.isInstalled ? 'Installer' : 'Portabel'})</span>
                  </button>
                {/if}
              </div>
            </div>
          {:else}
            <div class="flex items-center justify-between py-1">
              <div class="flex items-center gap-2">
                <span class="w-2.5 h-2.5 rounded-full bg-emerald-400 shadow-[0_0_8px_#34D399]"></span>
                <span class="text-emerald-400 text-xs font-semibold">
                  Aplikasi sudah versi terbaru ({$updateState?.currentVersion || 'v1.2.0'})
                </span>
              </div>
              <button
                class="px-3.5 py-1.5 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-xs font-semibold text-zinc-200 hover:text-white transition cursor-pointer disabled:opacity-50"
                on:click={checkForUpdatesManual}
                disabled={$isCheckingUpdate}
              >
                Cek Pembaruan
              </button>
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>
