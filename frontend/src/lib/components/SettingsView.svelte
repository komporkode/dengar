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
  import {
    equalizerStore,
    loudnessNormalizationStore,
  } from "../stores/equalizerStore";
  import { currentLocale, t } from "../stores/i18nStore";
  import { CURRENT_APP_VERSION } from "../data/changelogs";
  import { offlineStore } from "../stores/offlineStore";
  import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime.js";
  import {
    GetAudioCacheStats,
    ClearAudioCache,
    SetAudioCacheLimit,
    SelectAudioCacheFolder,
    SetAudioCacheFolder,
    ResetAudioCacheFolder,
    GetAudioCacheDir,
    OpenAudioCacheFolder,
  } from "../../../wailsjs/go/main/App";
  import type { backend } from "../../../wailsjs/go/models";
  import {
    updateState,
    isCheckingUpdate,
    isDownloadingUpdate,
    downloadProgressData,
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
  const savedLimitMb =
    typeof localStorage !== "undefined"
      ? localStorage.getItem("dengar_cache_limit_mb")
      : null;
  let selectedLimitMb: number =
    savedLimitMb !== null && !isNaN(Number(savedLimitMb))
      ? Number(savedLimitMb)
      : 1024;

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
        const mb =
          stats.maxBytes > 0 ? Math.round(stats.maxBytes / (1024 * 1024)) : 0;
        selectedLimitMb = mb;
        if (typeof localStorage !== "undefined") {
          localStorage.setItem("dengar_cache_limit_mb", String(mb));
        }
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
      await offlineStore.load();
      showToast($t("toast_cleared_cache"));
    } catch (e) {
      console.warn("Clear cache error:", e);
      showToast("Gagal membersihkan cache");
    } finally {
      isClearingCache = false;
    }
  }

  async function handleLimitChange(limitMb: number) {
    selectedLimitMb = limitMb;
    if (typeof localStorage !== "undefined") {
      localStorage.setItem("dengar_cache_limit_mb", String(limitMb));
    }
    const bytes = limitMb > 0 ? limitMb * 1024 * 1024 : 0;
    try {
      await SetAudioCacheLimit(bytes);
      await loadCacheStats();
      showToast(
        limitMb > 0
          ? `${$t("toast_cache_limit")} ${limitMb >= 1024 ? limitMb / 1024 + " GB" : limitMb + " MB"}`
          : $t("toast_cache_unlimited"),
      );
    } catch (e) {
      console.warn("Set limit error:", e);
    }
  }

  let isChangingCacheDir = false;

  async function handleChangeCacheDir() {
    if (isChangingCacheDir) return;
    try {
      const selected = await SelectAudioCacheFolder();
      if (!selected) return;
      isChangingCacheDir = true;
      await SetAudioCacheFolder(selected, true);
      await loadCacheStats();
      await offlineStore.load();
      showToast($t("toast_cache_dir_changed"));
    } catch (e: any) {
      console.warn("Change cache dir failed:", e);
      showToast($t("toast_cache_dir_failed"));
    } finally {
      isChangingCacheDir = false;
    }
  }

  async function handleResetCacheDir() {
    if (isChangingCacheDir) return;
    try {
      isChangingCacheDir = true;
      await ResetAudioCacheFolder();
      await loadCacheStats();
      await offlineStore.load();
      showToast($t("toast_cache_dir_reset"));
    } catch (e: any) {
      console.warn("Reset cache dir failed:", e);
    } finally {
      isChangingCacheDir = false;
    }
  }

  async function handleOpenCacheFolder() {
    try {
      await OpenAudioCacheFolder();
    } catch (e) {
      console.warn("Open cache folder failed:", e);
    }
  }

  function openGithub() {
    try {
      BrowserOpenURL("https://github.com/komporkode/dengar");
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

<div
  class="w-full flex-1 px-4 sm:px-8 py-6 space-y-6 max-w-6xl mx-auto select-none"
>
  <!-- Top Banner / Navigation Header (Omnitik Fullscreen Tab Header) -->
  <div class="flex items-center gap-3.5 pb-6 border-b border-white/[0.08]">
    <div
      class="w-11 h-11 rounded-2xl bg-gradient-to-tr from-crimson-600 to-rose-500 flex items-center justify-center shadow-lg shadow-crimson-500/25 shrink-0"
    >
      <svg
        class="w-5 h-5 text-white"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path
          d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"
        />
        <circle cx="12" cy="12" r="3" />
      </svg>
    </div>
    <div>
      <h1
        class="text-xl font-bold text-white tracking-tight flex items-center gap-2"
      >
        <span>{$t("settings_title")}</span>
        <span
          class="text-xs px-2 py-0.5 rounded-full bg-crimson-500/20 text-crimson-300 border border-crimson-500/30 font-mono font-medium"
          >{CURRENT_APP_VERSION}</span
        >
      </h1>
      <p class="text-xs text-zinc-400 mt-0.5">
        {$t("settings_subtitle")}
      </p>
    </div>
  </div>

  <!-- Omnitik-style Menu Tabs Navigation (Segmented pill bar) -->
  <div
    class="flex items-center gap-2 p-1.5 rounded-2xl bg-black/30 border border-white/[0.06] backdrop-blur-md"
  >
    <button
      class="flex-1 flex items-center justify-center gap-2 py-2.5 px-4 rounded-xl text-xs font-semibold transition-all cursor-pointer {activeTab ===
      'general'
        ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25'
        : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
      on:click={() => (activeTab = "general")}
    >
      <svg
        class="w-4 h-4 shrink-0"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="2" y1="12" x2="22" y2="12"></line>
        <path
          d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"
        ></path>
      </svg>
      <span>{$t("tab_general_ui")}</span>
    </button>

    <button
      class="flex-1 flex items-center justify-center gap-2 py-2.5 px-4 rounded-xl text-xs font-semibold transition-all cursor-pointer {activeTab ===
      'audio'
        ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25'
        : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
      on:click={() => (activeTab = "audio")}
    >
      <svg
        class="w-4 h-4 shrink-0"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
        <path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"
        ></path>
      </svg>
      <span>{$t("tab_audio_cache")}</span>
    </button>

    <button
      class="flex-1 flex items-center justify-center gap-2 py-2.5 px-4 rounded-xl text-xs font-semibold transition-all cursor-pointer {activeTab ===
      'about'
        ? 'bg-crimson-500 text-white shadow-md shadow-crimson-500/25'
        : 'text-zinc-400 hover:text-white hover:bg-white/[0.05]'}"
      on:click={() => (activeTab = "about")}
    >
      <svg
        class="w-4 h-4 shrink-0"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="16" x2="12" y2="12"></line>
        <line x1="12" y1="8" x2="12.01" y2="8"></line>
      </svg>
      <span>{$t("tab_about_updates")}</span>
    </button>
  </div>

  <!-- Content Sections -->
  {#if activeTab === "general"}
    <!-- TAB 1: General & Interface -->
    <div class="space-y-6 animate-fade-in">
      <!-- Language Selection Card -->
      <div
        class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4"
      >
        <div>
          <h3 class="text-sm font-bold text-white">{$t("lang_title")}</h3>
          <p class="text-xs text-zinc-400 mt-0.5">{$t("lang_desc")}</p>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
          <button
            class="flex items-center justify-between p-4 rounded-2xl border transition-all cursor-pointer {$currentLocale ===
            'id'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => currentLocale.set("id")}
          >
            <div class="flex items-center gap-3">
              <span
                class="w-8 h-8 rounded-xl bg-white/[0.08] flex items-center justify-center font-bold text-xs font-mono"
                >ID</span
              >
              <div class="text-left">
                <p class="text-xs font-bold">{$t("lang_id")}</p>
                <p class="text-[11px] text-zinc-400">{$t("lang_id_sub")}</p>
              </div>
            </div>
            {#if $currentLocale === "id"}
              <div
                class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0"
              >
                <svg
                  class="w-3 h-3"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="3"
                >
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </div>
            {/if}
          </button>

          <button
            class="flex items-center justify-between p-4 rounded-2xl border transition-all cursor-pointer {$currentLocale ===
            'en'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => currentLocale.set("en")}
          >
            <div class="flex items-center gap-3">
              <span
                class="w-8 h-8 rounded-xl bg-white/[0.08] flex items-center justify-center font-bold text-xs font-mono"
                >EN</span
              >
              <div class="text-left">
                <p class="text-xs font-bold">{$t("lang_en")}</p>
                <p class="text-[11px] text-zinc-400">{$t("lang_en_sub")}</p>
              </div>
            </div>
            {#if $currentLocale === "en"}
              <div
                class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0"
              >
                <svg
                  class="w-3 h-3"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="3"
                >
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </div>
            {/if}
          </button>
        </div>
      </div>

      <!-- Close Action Behavior Card -->
      <div
        class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4"
      >
        <div>
          <h3 class="text-sm font-bold text-white">
            {$t("close_action_title")}
          </h3>
          <p class="text-xs text-zinc-400 mt-0.5">{$t("close_action_desc")}</p>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <button
            class="p-4 rounded-2xl border text-left transition-all cursor-pointer {$closeActionStore ===
            'ask'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => {
              closeActionStore.set("ask");
              showToast($t("close_action_toast_ask"));
            }}
          >
            <div class="text-xs font-bold">{$t("close_action_ask")}</div>
            <p class="text-[11px] text-zinc-400 mt-1">
              {$t("close_action_ask_desc")}
            </p>
          </button>

          <button
            class="p-4 rounded-2xl border text-left transition-all cursor-pointer {$closeActionStore ===
            'tray'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => {
              closeActionStore.set("tray");
              showToast($t("close_action_toast_tray"));
            }}
          >
            <div class="text-xs font-bold">{$t("close_action_tray")}</div>
            <p class="text-[11px] text-zinc-400 mt-1">
              {$t("close_action_tray_desc")}
            </p>
          </button>

          <button
            class="p-4 rounded-2xl border text-left transition-all cursor-pointer {$closeActionStore ===
            'quit'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => {
              closeActionStore.set("quit");
              showToast($t("close_action_toast_quit"));
            }}
          >
            <div class="text-xs font-bold">{$t("close_action_quit")}</div>
            <p class="text-[11px] text-zinc-400 mt-1">
              {$t("close_action_quit_desc")}
            </p>
          </button>
        </div>
      </div>

      <!-- Discord RPC Card -->
      <div
        class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] flex items-center justify-between gap-4"
      >
        <div class="flex items-center gap-3.5">
          <div
            class="w-10 h-10 rounded-2xl bg-[#5865F2]/20 text-[#5865F2] flex items-center justify-center shrink-0"
          >
            <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path
                d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028c.462-.63.874-1.295 1.226-1.994.021-.041.001-.09-.041-.106a13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.929 1.793 8.18 1.793 12.061 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.894.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.028zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z"
              />
            </svg>
          </div>
          <div>
            <h3 class="text-sm font-bold text-white">
              {$t("discord_rpc_title")}
            </h3>
            <p class="text-xs text-zinc-400 mt-0.5">{$t("discord_rpc_desc")}</p>
          </div>
        </div>
        <button
          type="button"
          on:click={() => {
            const next = !$isDiscordRpcEnabled;
            isDiscordRpcEnabled.set(next);
            showToast(
              next ? $t("discord_rpc_enabled") : $t("discord_rpc_disabled"),
            );
          }}
          class="w-12 h-6 rounded-full transition-colors relative cursor-pointer shrink-0 {$isDiscordRpcEnabled
            ? 'bg-crimson-600'
            : 'bg-white/10'}"
          aria-label="Toggle Discord RPC"
        >
          <span
            class="absolute top-1 w-4 h-4 rounded-full bg-white transition-transform {$isDiscordRpcEnabled
              ? 'left-7'
              : 'left-1'}"
          ></span>
        </button>
      </div>

      <!-- Keyboard Shortcuts Reference Card -->
      <div
        class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4"
      >
        <div class="flex items-center gap-2 text-white text-sm font-bold">
          <svg
            class="w-4 h-4 text-crimson-400"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            viewBox="0 0 24 24"
          >
            <rect x="2" y="4" width="20" height="16" rx="2" ry="2"></rect>
            <path
              d="M6 8h.001M10 8h.001M14 8h.001M18 8h.001M6 12h.001M18 12h.001M8 16h8"
            ></path>
          </svg>
          <span>{$t("shortcuts_title")}</span>
        </div>
        <div
          class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3 text-xs"
        >
          <div
            class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]"
          >
            <span class="text-zinc-400">{$t("shortcut_play_pause")}</span>
            <kbd
              class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]"
              >{$t("space_key")}</kbd
            >
          </div>
          <div
            class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]"
          >
            <span class="text-zinc-400">{$t("shortcut_volume")}</span>
            <kbd
              class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]"
              >↑ / ↓</kbd
            >
          </div>
          <div
            class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]"
          >
            <span class="text-zinc-400">{$t("shortcut_seek")}</span>
            <kbd
              class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]"
              >← / →</kbd
            >
          </div>
          <div
            class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]"
          >
            <span class="text-zinc-400">{$t("shortcut_prev_next")}</span>
            <kbd
              class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]"
              >J / L</kbd
            >
          </div>
          <div
            class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]"
          >
            <span class="text-zinc-400">{$t("shortcut_mute")}</span>
            <kbd
              class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]"
              >M</kbd
            >
          </div>
          <div
            class="flex items-center justify-between p-3 rounded-xl bg-white/[0.02] border border-white/[0.05]"
          >
            <span class="text-zinc-400">{$t("shortcut_back")}</span>
            <kbd
              class="px-2 py-0.5 rounded bg-white/10 text-zinc-200 font-mono text-[11px]"
              >Esc</kbd
            >
          </div>
        </div>
      </div>
    </div>
  {:else if activeTab === "audio"}
    <!-- TAB 2: Audio & Disk Cache -->
    <div class="space-y-6 animate-fade-in">
      <!-- Audio Streaming Resolution -->
      <div
        class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4"
      >
        <div>
          <h3 class="text-sm font-bold text-white">{$t("quality_setting")}</h3>
          <p class="text-xs text-zinc-400 mt-0.5">{$t("quality_desc")}</p>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
          <button
            class="p-4 rounded-2xl border text-left transition-all cursor-pointer flex items-center justify-between {$audioQualityStore ===
            'high'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => {
              audioQualityStore.set("high");
              showToast($t("quality_high"));
            }}
          >
            <div class="flex items-center gap-3">
              <span
                class="px-2 py-1 rounded-xl bg-amber-500/20 text-amber-300 font-mono text-xs font-bold border border-amber-500/30"
                >HQ</span
              >
              <div>
                <div class="font-bold text-xs leading-tight">
                  {$t("quality_high")}
                </div>
                <div class="text-[11px] text-zinc-400 mt-0.5">
                  {$t("quality_high_sub")}
                </div>
              </div>
            </div>
            {#if $audioQualityStore === "high"}
              <div
                class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0"
              >
                <svg
                  class="w-3 h-3"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="3"
                >
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </div>
            {/if}
          </button>

          <button
            class="p-4 rounded-2xl border text-left transition-all cursor-pointer flex items-center justify-between {$audioQualityStore ===
            'saver'
              ? 'bg-crimson-500/15 border-crimson-500 text-white shadow-sm'
              : 'bg-white/[0.02] border-white/10 text-zinc-300 hover:bg-white/[0.05]'}"
            on:click={() => {
              audioQualityStore.set("saver");
              showToast($t("quality_saver"));
            }}
          >
            <div class="flex items-center gap-3">
              <span
                class="px-2 py-1 rounded-xl bg-emerald-500/20 text-emerald-300 font-mono text-xs font-bold border border-emerald-500/30"
                >ECO</span
              >
              <div>
                <div class="font-bold text-xs leading-tight">
                  {$t("quality_saver")}
                </div>
                <div class="text-[11px] text-zinc-400 mt-0.5">
                  {$t("quality_saver_sub")}
                </div>
              </div>
            </div>
            {#if $audioQualityStore === "saver"}
              <div
                class="w-5 h-5 rounded-full bg-crimson-500 flex items-center justify-center text-white shrink-0"
              >
                <svg
                  class="w-3 h-3"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="3"
                >
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </div>
            {/if}
          </button>
        </div>
      </div>

      <!-- Web Audio Crossfade Slider -->
      <div
        class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4"
      >
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-sm font-bold text-white">
              {$t("crossfade_title")}
            </h3>
            <p class="text-xs text-zinc-400 mt-0.5">{$t("crossfade_desc")}</p>
          </div>
          <span
            class="text-xs font-mono font-bold text-crimson-400 bg-crimson-500/15 border border-crimson-500/25 px-2.5 py-1 rounded-lg"
          >
            {$crossfadeDurationStore > 0
              ? `${$crossfadeDurationStore} ${$t("crossfade_seconds")}`
              : $t("crossfade_disabled")}
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

      <!-- Loudness Normalization -->
      <div
        class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] flex items-center justify-between gap-4"
      >
        <div>
          <h3 class="text-sm font-bold text-white">
            {$t("loudness_title")}
          </h3>
          <p class="text-xs text-zinc-400 mt-0.5">
            {$t("loudness_desc")}
          </p>
        </div>
        <button
          class="w-12 h-6 rounded-full transition-colors relative cursor-pointer {$loudnessNormalizationStore
            ? 'bg-crimson-500'
            : 'bg-white/10'}"
          on:click={() =>
            equalizerStore.setLoudnessNormalization(
              !$loudnessNormalizationStore,
            )}
          role="switch"
          aria-checked={$loudnessNormalizationStore}
          aria-label="Toggle normalisasi kenyaringan suara"
          title="Toggle normalisasi kenyaringan suara"
        >
          <span
            class="w-4 h-4 rounded-full bg-white absolute top-1 transition-transform {$loudnessNormalizationStore
              ? 'left-7'
              : 'left-1'}"
          ></span>
        </button>
      </div>

      <!-- Smart Disk Cache & Offline Storage Management -->
      <div
        class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4"
      >
        <div class="flex items-center justify-between gap-4">
          <div>
            <div class="flex items-center gap-2.5">
              <h3 class="text-sm font-bold text-white">
                {$t("smart_cache_title")}
              </h3>
              <span
                class="text-xs font-mono px-2 py-0.5 rounded-full bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 font-semibold"
              >
                {formatBytes(cacheStats?.totalBytes || 0)}
              </span>
            </div>
            <p class="text-xs text-zinc-400 mt-1">
              {$t("smart_cache_desc_prefix")}{cacheStats?.totalFiles || 0}
              {$t("smart_cache_desc_suffix")}
            </p>
          </div>

          <button
            class="px-4 py-2 rounded-xl bg-crimson-600/20 hover:bg-crimson-600 border border-crimson-500/30 text-xs font-semibold text-crimson-200 hover:text-white transition cursor-pointer flex items-center gap-2 shrink-0 active:scale-95 disabled:opacity-50"
            on:click={handleClearCache}
            disabled={isClearingCache || !cacheStats?.totalFiles}
            title={$t("btn_clear_cache")}
          >
            {#if isClearingCache}
              <svg
                class="w-3.5 h-3.5 animate-spin"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <circle cx="12" cy="12" r="10" stroke-opacity="0.25" />
                <path d="M12 2a10 10 0 0 1 10 10" />
              </svg>
              <span>{$t("clearing_cache")}</span>
            {:else}
              <svg
                class="w-3.5 h-3.5 text-crimson-400"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <polyline points="3 6 5 6 21 6"></polyline>
                <path
                  d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                ></path>
              </svg>
              <span>{$t("btn_clear_cache")}</span>
            {/if}
          </button>
        </div>

        <div
          class="pt-4 border-t border-white/[0.06] flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs"
        >
          <span class="text-zinc-400">{$t("cache_limit_label")}</span>
          <div class="flex items-center gap-2">
            {#each [500, 1024, 2048, 0] as limit}
              <button
                class="px-3 py-1 rounded-xl text-xs font-mono transition cursor-pointer {selectedLimitMb ===
                limit
                  ? 'bg-crimson-500 text-white font-bold shadow-xs'
                  : 'text-zinc-400 hover:text-zinc-200 bg-white/[0.04]'}"
                on:click={() => handleLimitChange(limit)}
              >
                {limit === 0
                  ? $t("cache_unlimited")
                  : limit >= 1024
                    ? limit / 1024 + " GB"
                    : limit + " MB"}
              </button>
            {/each}
          </div>
        </div>

        <!-- Custom Cache Directory Location & Folder Actions -->
        <div class="pt-4 border-t border-white/[0.06] flex flex-col gap-2">
          <div
            class="flex flex-col sm:flex-row sm:items-center justify-between gap-3"
          >
            <div class="min-w-0 flex-1">
              <span class="text-xs font-semibold text-zinc-300"
                >{$t("cache_location_title")}</span
              >
              <p
                class="text-[11px] font-mono text-zinc-400 truncate mt-0.5 bg-black/30 px-2 py-1 rounded-lg border border-white/[0.04]"
                title={cacheStats?.cacheDir || ""}
              >
                {cacheStats?.cacheDir || "Default AppData"}
              </p>
            </div>
            <div class="flex items-center gap-2 shrink-0 pt-1 sm:pt-0">
              <button
                type="button"
                class="px-3 py-1.5 rounded-xl bg-white/[0.05] hover:bg-white/[0.1] border border-white/10 text-xs font-medium text-zinc-300 hover:text-white transition flex items-center gap-1.5 cursor-pointer disabled:opacity-50"
                on:click={handleOpenCacheFolder}
                title={$t("btn_open_folder")}
              >
                <svg
                  class="w-3.5 h-3.5 text-zinc-400"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"
                  ></path>
                </svg>
                <span>{$t("btn_open_folder")}</span>
              </button>

              <button
                type="button"
                class="px-3 py-1.5 rounded-xl bg-crimson-600/20 hover:bg-crimson-600/35 border border-crimson-500/30 text-xs font-semibold text-crimson-200 hover:text-white transition flex items-center gap-1.5 cursor-pointer disabled:opacity-50"
                on:click={handleChangeCacheDir}
                disabled={isChangingCacheDir}
                title={$t("btn_change_location")}
              >
                {#if isChangingCacheDir}
                  <div
                    class="w-3 h-3 border-2 border-white/20 border-t-white rounded-full animate-spin"
                  ></div>
                {:else}
                  <svg
                    class="w-3.5 h-3.5 text-crimson-400"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path
                      d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
                    ></path>
                    <path
                      d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"
                    ></path>
                  </svg>
                {/if}
                <span>{$t("btn_change_location")}</span>
              </button>

              <button
                type="button"
                class="px-2.5 py-1.5 rounded-xl text-xs text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04] transition cursor-pointer disabled:opacity-50"
                on:click={handleResetCacheDir}
                disabled={isChangingCacheDir}
                title={$t("btn_reset_default")}
              >
                {$t("btn_reset_default")}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  {:else if activeTab === "about"}
    <!-- TAB 3: About & Updates -->
    <div class="space-y-6 animate-fade-in">
      <!-- Brand Card -->
      <div
        class="p-6 rounded-2xl bg-gradient-to-br from-[#161824] to-[#10121C] border border-white/10 flex flex-col md:flex-row items-start md:items-center gap-5 shadow-sm"
      >
        <div
          class="w-16 h-16 rounded-2xl bg-gradient-to-tr from-crimson-600 via-crimson-500 to-rose-400 flex items-center justify-center shadow-lg shadow-crimson-500/25 shrink-0"
        >
          <svg
            class="w-8 h-8 text-white"
            viewBox="0 0 269.22 267.55"
            fill="none"
          >
            <path
              fill="currentColor"
              d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"
            />
          </svg>
        </div>
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2.5 flex-wrap">
            <h2 class="text-base font-extrabold text-white">
              Dengerin Desktop
            </h2>
            <span
              class="text-xs font-bold px-2 py-0.5 rounded-full bg-crimson-500/20 border border-crimson-500/30 text-crimson-400 font-mono"
              >{CURRENT_APP_VERSION}</span
            >
            <span
              class="text-[10px] font-medium px-2 py-0.5 rounded-full {$updateState?.isInstalled
                ? 'bg-sky-500/20 border-sky-500/30 text-sky-400'
                : 'bg-emerald-500/20 border-emerald-500/30 text-emerald-400'} border font-mono"
            >
              {$updateState?.isInstalled
                ? $t("package_installed")
                : $t("package_portable")}
            </span>
          </div>
          <p class="text-xs text-zinc-300 leading-relaxed mt-1.5">
            {$t("about_desc")}
          </p>
          <div
            class="flex items-center gap-2 mt-2 text-[10px] font-medium text-zinc-300 flex-wrap"
          >
            <span
              class="px-2 py-0.5 rounded bg-white/[0.06] border border-white/10"
              >{$t("zero_ads")}</span
            >
            <span
              class="px-2 py-0.5 rounded bg-white/[0.06] border border-white/10"
              >{$t("free_rp0")}</span
            >
            <span
              class="px-2 py-0.5 rounded bg-white/[0.06] border border-white/10"
              >{$t("low_ram")}</span
            >
            <span
              class="px-2 py-0.5 rounded bg-white/[0.06] border border-white/10"
              >GNU GPLv3</span
            >
          </div>
          <div class="mt-3.5">
            <button
              type="button"
              class="px-3 py-1.5 rounded-xl bg-crimson-500/15 hover:bg-crimson-500/25 border border-crimson-500/30 text-xs font-semibold text-crimson-300 hover:text-white transition cursor-pointer flex items-center gap-1.5"
              on:click={() => isWhatsNewOpen.set(true)}
            >
              <svg
                class="w-3.5 h-3.5 text-crimson-400 shrink-0"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"
                />
              </svg>
              <span>{$t("whats_new_btn")}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- GitHub & License Card -->
      <div
        class="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.06] space-y-4"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div
              class="w-9 h-9 rounded-xl bg-white/10 flex items-center justify-center text-white shrink-0"
            >
              <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
                <path
                  fill-rule="evenodd"
                  clip-rule="evenodd"
                  d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.53 1.032 1.53 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
                />
              </svg>
            </div>
            <div>
              <h4 class="text-sm font-bold text-white">
                {$t("github_repo_title")}
              </h4>
              <p class="text-xs text-zinc-400 font-mono">
                {$t("github_repo_sub")}
              </p>
            </div>
          </div>

          <button
            class="px-3.5 py-1.5 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-xs font-semibold text-zinc-200 transition cursor-pointer flex items-center gap-2"
            on:click={openGithub}
            title={$t("btn_visit_repo")}
          >
            <span>{$t("btn_visit_repo")}</span>
            <svg
              class="w-3.5 h-3.5 text-zinc-400"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"
              ></path>
              <polyline points="15 3 21 3 21 9"></polyline>
              <line x1="10" y1="14" x2="21" y2="3"></line>
            </svg>
          </button>
        </div>

        <!-- Update Status Section -->
        <div class="pt-4 border-t border-white/[0.06] space-y-3">
          {#if $isCheckingUpdate}
            <div class="flex items-center gap-2 py-1">
              <div
                class="w-4 h-4 border-2 border-white/20 border-t-crimson-500 rounded-full animate-spin"
              ></div>
              <span class="text-zinc-400 text-xs">{$t("checking_updates")}</span
              >
            </div>
          {:else if $updateState?.hasUpdate}
            <div
              class="p-4 rounded-xl bg-crimson-500/10 border border-crimson-500/25 space-y-3"
            >
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2 flex-wrap">
                  <span
                    class="w-2.5 h-2.5 rounded-full bg-crimson-400 shadow-[0_0_8px_#E61E38] animate-pulse"
                  ></span>
                  <span class="text-crimson-300 font-bold text-xs"
                    >{$t("update_available")}</span
                  >
                  <span
                    class="px-2 py-0.5 rounded bg-crimson-500/25 text-white font-mono font-bold text-xs"
                  >
                    {$updateState.latestVersion}
                  </span>
                  <span
                    class="px-2 py-0.5 rounded {$updateState.isInstalled
                      ? 'bg-sky-500/25 text-sky-200 border border-sky-500/30'
                      : 'bg-emerald-500/25 text-emerald-200 border border-emerald-500/30'} font-mono text-[10px]"
                  >
                    {$updateState.isInstalled
                      ? "Installer .exe"
                      : "Portable .zip"}
                  </span>
                </div>
              </div>

              {#if $updateState?.releaseNotes}
                <details
                  class="text-xs text-zinc-300 bg-black/20 p-3 rounded-xl border border-white/[0.06] cursor-pointer"
                >
                  <summary
                    class="font-semibold text-zinc-200 hover:text-white transition flex items-center justify-between"
                  >
                    <span>Lihat Catatan Rilis {$updateState.latestVersion}</span
                    >
                    <span class="text-[10px] text-zinc-500 font-mono"
                      >Klik untuk buka/tutup</span
                    >
                  </summary>
                  <div
                    class="mt-2.5 text-zinc-300 whitespace-pre-wrap font-sans text-xs max-h-40 overflow-y-auto leading-relaxed border-t border-white/[0.06] pt-2.5 scrollbar-thin"
                  >
                    {$updateState.releaseNotes}
                  </div>
                </details>
              {/if}

              <p class="text-xs text-zinc-300 leading-relaxed">
                {#if $downloadedFilePath}
                  {$t("update_download_ready")}
                {:else if $isDownloadingUpdate}
                  Sedang mengunduh berkas pembaruan...
                {:else}
                  {$updateState.latestVersion} {$t("update_ready_to_download")}
                {/if}
              </p>

              <div class="space-y-2 pt-1">
                {#if $downloadedFilePath}
                  <div class="flex items-center gap-2.5">
                    <button
                      class="flex-1 py-2 px-4 rounded-xl bg-crimson-600 hover:bg-crimson-500 text-white font-bold text-xs transition cursor-pointer shadow-md flex items-center justify-center gap-2"
                      on:click={installAndRestartUpdate}
                    >
                      <svg
                        class="w-4 h-4"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        viewBox="0 0 24 24"
                      >
                        <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"
                        ></polygon>
                      </svg>
                      <span>{$t("btn_install_restart")}</span>
                    </button>
                    <button
                      class="py-2 px-4 rounded-xl bg-white/[0.08] hover:bg-white/[0.14] text-zinc-300 hover:text-white font-semibold text-xs transition cursor-pointer shadow-sm flex items-center gap-1.5 shrink-0"
                      on:click={revealDownloadedFile}
                    >
                      <svg
                        class="w-4 h-4"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        viewBox="0 0 24 24"
                      >
                        <path
                          d="M5 19a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2h4l2 2h7a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5z"
                        ></path>
                      </svg>
                      <span>{$t("btn_open_folder")}</span>
                    </button>
                  </div>
                {:else if $isDownloadingUpdate}
                  <div class="space-y-2 w-full py-1">
                    <div
                      class="flex items-center justify-between text-xs text-zinc-300"
                    >
                      <span class="flex items-center gap-2 font-medium">
                        <span
                          class="w-2 h-2 rounded-full bg-crimson-500 animate-pulse"
                        ></span>
                        Mengunduh berkas pembaruan...
                      </span>
                      <span class="font-mono text-crimson-400 font-bold">
                        {$downloadProgressData
                          ? `${$downloadProgressData.percent}%`
                          : "0%"}
                      </span>
                    </div>
                    <div
                      class="h-2 w-full bg-white/10 rounded-full overflow-hidden relative"
                    >
                      <div
                        class="h-full bg-gradient-to-r from-crimson-600 to-rose-500 rounded-full transition-all duration-150"
                        style="width: {$downloadProgressData
                          ? $downloadProgressData.percent
                          : 0}%"
                      ></div>
                    </div>
                    <div
                      class="flex items-center justify-between text-[11px] text-zinc-400 font-mono"
                    >
                      <span>
                        {$downloadProgressData &&
                        $downloadProgressData.total > 0
                          ? `${formatBytes($downloadProgressData.downloaded)} / ${formatBytes($downloadProgressData.total)}`
                          : "Menghubungkan..."}
                      </span>
                      <span>{$downloadProgressData?.speedText || ""}</span>
                    </div>
                  </div>
                {:else}
                  <button
                    class="w-full py-2.5 px-4 rounded-xl bg-crimson-500 hover:bg-crimson-600 text-white font-bold text-xs transition cursor-pointer shadow-md shadow-crimson-500/30 flex items-center justify-center gap-2"
                    on:click={() => startDownloadUpdate()}
                  >
                    <svg
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      viewBox="0 0 24 24"
                    >
                      <path d="M12 4v12m0 0l-4-4m4 4l4-4M4 20h16" />
                    </svg>
                    <span
                      >{$t("btn_download_update")} ({$updateState.isInstalled
                        ? "Installer .exe"
                        : "Portable .zip"})</span
                    >
                  </button>

                  <!-- Alternate Edition Download Option -->
                  {#if $updateState.isInstalled && $updateState.portableUrl}
                    <div class="text-center pt-1">
                      <button
                        type="button"
                        class="text-[11px] text-zinc-400 hover:text-crimson-400 transition cursor-pointer underline underline-offset-2"
                        on:click={() =>
                          startDownloadUpdate(
                            $updateState?.portableUrl,
                            $updateState?.portableFileName,
                          )}
                      >
                        Atau unduh versi Portabel (.zip)
                      </button>
                    </div>
                  {:else if !$updateState.isInstalled && $updateState.installerUrl}
                    <div class="text-center pt-1">
                      <button
                        type="button"
                        class="text-[11px] text-zinc-400 hover:text-crimson-400 transition cursor-pointer underline underline-offset-2"
                        on:click={() =>
                          startDownloadUpdate(
                            $updateState?.installerUrl,
                            $updateState?.installerFileName,
                          )}
                      >
                        Atau unduh versi Installer (.exe)
                      </button>
                    </div>
                  {/if}
                {/if}
              </div>
            </div>
          {:else}
            <div class="flex items-center justify-between py-1">
              <div class="flex items-center gap-2">
                <span
                  class="w-2.5 h-2.5 rounded-full bg-emerald-400 shadow-[0_0_8px_#34D399]"
                ></span>
                <span class="text-emerald-400 text-xs font-semibold">
                  {$t("app_up_to_date")} ({$updateState?.currentVersion ||
                    CURRENT_APP_VERSION})
                </span>
              </div>
              <button
                class="px-3.5 py-1.5 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-xs font-semibold text-zinc-200 hover:text-white transition cursor-pointer disabled:opacity-50"
                on:click={checkForUpdatesManual}
                disabled={$isCheckingUpdate}
              >
                {$t("btn_check_update")}
              </button>
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>
