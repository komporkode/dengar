<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import type {
    NavigationTab,
    SearchResultItem,
    SearchFilter,
  } from "./lib/types";
  import { Search, SearchMore } from "../wailsjs/go/main/App";
  import {
    Quit,
    WindowMinimise,
    WindowToggleMaximise,
  } from "../wailsjs/runtime/runtime.js";
  import { playerStore } from "./lib/stores/playerStore";
  import { userDataStore } from "./lib/stores/userDataStore";
  import {
    isMiniPlayer,
    toastMessage,
    togglePinMiniPlayer,
    exitMiniPlayer,
    enterMiniPlayer,
    isTransitioning,
    isWhatsNewOpen,
  } from "./lib/stores/settingsStore";
  import {
    updateState,
    showUpdateBanner,
    isDownloadingUpdate,
    downloadedFilePath,
    checkForUpdatesInBackground,
    startDownloadUpdate,
    revealDownloadedFile,
    installAndRestartUpdate,
  } from "./lib/stores/updateStore";
  import { t } from "./lib/stores/i18nStore";

  import Sidebar from "./lib/components/Sidebar.svelte";
  import SearchBar from "./lib/components/SearchBar.svelte";
  import TrackList from "./lib/components/TrackList.svelte";
  import ForYouView from "./lib/components/ForYouView.svelte";
  import PlayerBar from "./lib/components/PlayerBar.svelte";
  import QueueDrawer from "./lib/components/QueueDrawer.svelte";
  import SettingsModal from "./lib/components/SettingsModal.svelte";
  import MiniPlayer from "./lib/components/MiniPlayer.svelte";
  import HomeSections from "./lib/components/HomeSections.svelte";
  import PlaylistView from "./lib/components/PlaylistView.svelte";
  import PlaylistModal from "./lib/components/PlaylistModal.svelte";
  import AddToPlaylistModal from "./lib/components/AddToPlaylistModal.svelte";
  import PlaylistShareModal from "./lib/components/PlaylistShareModal.svelte";
  import WhatsNewModal from "./lib/components/WhatsNewModal.svelte";
  import NetworkBanner from "./lib/components/NetworkBanner.svelte";
  import { networkStore } from "./lib/stores/networkStore";
  import { activePlaylistId } from "./lib/stores/playlistStore";

  let activeTab: NavigationTab = "search";
  let searchQuery = "";
  let currentFilter: SearchFilter = "music_songs";
  let searchResults: SearchResultItem[] = [];
  let continuationToken: string = "";
  let isSearching = false;
  let isLoadingMore = false;
  let searchError: string | null = null;

  let scrollTimer: any;
  let handleScrollActivity: (e: Event) => void;

  onMount(async () => {
    // 1. Load user favorites & history from backend JSON storage
    await userDataStore.load();

    // 2. Perform initial search so app is instantly ready with tracks
    performSearch("Trending Music", "music_songs");

    // 3. Register global in-app keyboard shortcuts
    window.addEventListener("keydown", handleGlobalKeydown);

    // 4. Global autohide scrollbar indicator
    handleScrollActivity = (e: Event) => {
      const target = e.target as HTMLElement;
      if (target?.classList) {
        target.classList.add("is-scrolling");
      }
      document.body.classList.add("is-scrolling");
      clearTimeout(scrollTimer);
      scrollTimer = setTimeout(() => {
        if (target?.classList) {
          target.classList.remove("is-scrolling");
        }
        document.body.classList.remove("is-scrolling");
      }, 900);
    };
    window.addEventListener("scroll", handleScrollActivity, {
      capture: true,
      passive: true,
    });

    // 5. Automatic background update check on app launch
    checkForUpdatesInBackground();

    // 6. Check first run for current version to show What's New
    const CURRENT_VERSION = "v1.1.2";
    try {
      const seen = localStorage.getItem("dengerin_seen_version");
      if (seen !== CURRENT_VERSION) {
        setTimeout(() => {
          isWhatsNewOpen.set(true);
        }, 500);
      }
    } catch (e) {
      console.warn("Version check warning:", e);
    }
  });

  function handleCloseWhatsNew() {
    isWhatsNewOpen.set(false);
    try {
      localStorage.setItem("dengerin_seen_version", "v1.1.2");
    } catch (e) {
      console.warn("Save version warning:", e);
    }
  }

  onDestroy(() => {
    window.removeEventListener("keydown", handleGlobalKeydown);
    if (handleScrollActivity) {
      window.removeEventListener("scroll", handleScrollActivity, {
        capture: true,
      });
    }
    clearTimeout(scrollTimer);
  });

  async function performSearch(query: string, filter: SearchFilter) {
    if (!query) {
      searchResults = [];
      continuationToken = "";
      return;
    }
    if (!$networkStore.isOnline) {
      searchError = "Perangkat sedang offline. Sambungkan internet untuk mencari lagu.";
      searchResults = [];
      continuationToken = "";
      return;
    }
    isSearching = true;
    searchError = null;
    try {
      const res = await Search(query, filter);
      searchResults = res?.items || [];
      continuationToken = res?.continuationToken || "";
    } catch (err: any) {
      console.error("Search failed:", err);
      searchError =
        err?.message || "Gagal mencari lagu. Periksa koneksi internet.";
      searchResults = [];
      continuationToken = "";
    } finally {
      isSearching = false;
    }
  }

  async function handleLoadMore() {
    if (!continuationToken || isLoadingMore || isSearching) return;
    isLoadingMore = true;
    try {
      const res = await SearchMore(continuationToken);
      if (res && res.items && res.items.length > 0) {
        const existingIds = new Set(searchResults.map((i) => i.id));
        const fresh = res.items.filter((i) => !existingIds.has(i.id));
        searchResults = [...searchResults, ...fresh];
        continuationToken = res.continuationToken || "";
      } else {
        continuationToken = "";
      }
    } catch (err: any) {
      console.warn("SearchMore failed:", err);
    } finally {
      isLoadingMore = false;
    }
  }

  function handleSearchEvent(
    e: CustomEvent<{ query: string; filter: SearchFilter }>,
  ) {
    searchQuery = e.detail.query;
    currentFilter = e.detail.filter;
    const fallbackQuery =
      currentFilter === "videos" ? "Trending Music Video" : "Trending Music";
    performSearch(searchQuery.trim() || fallbackQuery, currentFilter);
  }

  function handleWindowClose() {
    try {
      Quit();
    } catch (e) {
      console.warn("Wails Quit error:", e);
    }
  }

  function handleWindowMinimise() {
    try {
      WindowMinimise();
    } catch (e) {
      console.warn("Wails WindowMinimise error:", e);
    }
  }

  function handleWindowToggleMaximise() {
    try {
      WindowToggleMaximise();
    } catch (e) {
      console.warn("Wails WindowToggleMaximise error:", e);
    }
  }

  function handleGlobalKeydown(e: KeyboardEvent) {
    const target = e.target as HTMLElement;
    // Avoid intercepting keystrokes when user is typing in search input
    if (
      target &&
      (target.tagName === "INPUT" || target.tagName === "TEXTAREA")
    ) {
      return;
    }

    switch (e.code) {
      case "Space":
        e.preventDefault();
        playerStore.togglePlay();
        break;
      case "ArrowLeft":
        e.preventDefault();
        playerStore.seek(Math.max(0, $playerStore.currentTime - 5));
        break;
      case "ArrowRight":
        e.preventDefault();
        playerStore.seek(
          Math.min($playerStore.duration, $playerStore.currentTime + 5),
        );
        break;
      case "ArrowUp":
        e.preventDefault();
        playerStore.adjustVolume(0.05);
        break;
      case "ArrowDown":
        e.preventDefault();
        playerStore.adjustVolume(-0.05);
        break;
      case "KeyM":
        if (e.ctrlKey || e.metaKey) {
          e.preventDefault();
          if ($isMiniPlayer) {
            exitMiniPlayer();
          } else {
            enterMiniPlayer();
          }
        } else {
          playerStore.toggleMute();
        }
        break;
      case "KeyJ":
        playerStore.playPrev();
        break;
      case "KeyL":
        playerStore.playNext();
        break;
      case "KeyF":
        if ($playerStore.currentTrack) {
          userDataStore.toggleFavorite($playerStore.currentTrack);
        }
        break;
      case "KeyP":
        if ($isMiniPlayer) {
          togglePinMiniPlayer();
        }
        break;
      case "Escape":
        if ($isMiniPlayer) {
          exitMiniPlayer();
        }
        break;
    }
  }
</script>

{#if $isMiniPlayer}
  <!-- Mini Music Player View -->
  <MiniPlayer />
{/if}

<!-- Standard Desktop Application View with Single Custom Frameless Title Bar -->
<div
  class="w-full h-full flex flex-col bg-obsidian text-zinc-100 font-sans antialiased select-none overflow-hidden relative transition-opacity duration-200 ease-out {$isTransitioning &&
  !$isMiniPlayer
    ? 'opacity-0'
    : 'opacity-100'}"
  style:display={$isMiniPlayer ? "none" : "flex"}
>
  <!-- Unified Dark Custom Titlebar (Frameless Windows-style, NO Mac dots) -->
  <header
    class="h-9 px-4 flex items-center justify-between bg-[#08080C] border-b border-white/[0.06] text-zinc-400 text-xs shrink-0 select-none z-50"
    style="--wails-draggable:drag"
  >
    <div class="flex items-center gap-2" style="--wails-draggable:no-drag">
      <div
        class="w-4 h-4 rounded-md bg-gradient-to-tr from-crimson-600 to-rose-500 flex items-center justify-center shadow-xs"
      >
        <svg
          class="w-2.5 h-2.5 text-white"
          viewBox="0 0 269.22 267.55"
          fill="none"
        >
          <path
            fill="currentColor"
            d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"
          />
        </svg>
      </div>
      <span class="text-zinc-200 font-semibold tracking-wide">Dengerin</span>
    </div>

    <!-- Center playing song / status (Responsive: hidden on narrow screens) -->
    <div
      class="hidden sm:flex text-[11px] font-medium text-zinc-400 items-center gap-2 truncate max-w-xs md:max-w-md pointer-events-none"
    >
      {#if $playerStore.currentTrack}
        <span class="text-zinc-300 truncate"
          >{$playerStore.currentTrack.title} — {$playerStore.currentTrack
            .uploaderName}</span
        >
      {:else}
        <span class="text-zinc-500">Master Audio Streaming Engine</span>
      {/if}
    </div>

    <!-- Windows-style window controls (NO Mac dots) -->
    <div
      class="flex items-center text-zinc-400 h-full -mr-4"
      style="--wails-draggable:no-drag"
    >
      <button
        class="h-full px-3.5 hover:text-zinc-100 hover:bg-white/[0.08] transition-colors cursor-pointer flex items-center justify-center"
        on:click={handleWindowMinimise}
        title="Minimalisir"
      >
        <svg
          class="w-3.5 h-3.5"
          fill="none"
          stroke="currentColor"
          stroke-width="1.8"
          viewBox="0 0 24 24"><path d="M19 12H5"></path></svg
        >
      </button>
      <button
        class="h-full px-3.5 hover:text-zinc-100 hover:bg-white/[0.08] transition-colors cursor-pointer flex items-center justify-center"
        on:click={handleWindowToggleMaximise}
        title="Maksimalkan"
      >
        <svg
          class="w-3.5 h-3.5"
          fill="none"
          stroke="currentColor"
          stroke-width="1.8"
          viewBox="0 0 24 24"
          ><rect height="13" rx="1.5" width="13" x="5.5" y="5.5"></rect></svg
        >
      </button>
      <button
        class="h-full px-4 hover:text-white hover:bg-crimson-600 transition-colors cursor-pointer flex items-center justify-center"
        on:click={handleWindowClose}
        title="Tutup"
      >
        <svg
          class="w-3.5 h-3.5"
          fill="none"
          stroke="currentColor"
          stroke-width="1.8"
          viewBox="0 0 24 24"><path d="M18 6L6 18M6 6l12 12"></path></svg
        >
      </button>
    </div>
  </header>

  <!-- Network Connectivity Banner (Offline / Weak / Restored) -->
  <NetworkBanner />

  <!-- Application Body Layout -->
  <div class="flex flex-1 overflow-hidden relative">
    <!-- Left Sidebar -->
    <Sidebar bind:activeTab />

    <!-- Main Content Viewport -->
    <main
      class="flex-1 flex flex-col overflow-y-auto bg-gradient-to-b from-[#12141F] via-obsidian to-obsidian pb-32"
    >
      {#if activeTab === "search"}
        <SearchBar
          bind:query={searchQuery}
          isLoading={isSearching}
          on:search={handleSearchEvent}
        />

        {#if searchError}
          <div
            class="mx-8 mt-4 flex items-center gap-3 px-4 py-3 rounded-xl bg-crimson-500/10 border border-crimson-500/30 text-crimson-400 text-xs"
          >
            <svg
              class="w-4 h-4 shrink-0"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>{searchError}</span>
          </div>
        {/if}

        {#if !searchQuery.trim()}
          <HomeSections />
        {/if}

        <TrackList
          title={searchQuery
            ? `${$t("search_results_prefix")}“${searchQuery}”`
            : $t("popular_songs")}
          tracks={searchResults}
          isLoading={isSearching}
          hasMore={!!continuationToken}
          {isLoadingMore}
          emptyMessage={$t("no_results")}
          on:loadMore={handleLoadMore}
        />
      {:else if activeTab === "foryou"}
        <ForYouView />
      {:else if activeTab === "favorites"}
        <TrackList
          title={$t("nav_favorites")}
          tracks={$userDataStore.favorites}
          isLoading={$userDataStore.isLoading}
          emptyMessage="Belum ada lagu yang ditandai favorit. Klik ikon hati pada lagu untuk menyimpannya!"
        />
      {:else if activeTab === "history"}
        <TrackList
          title={$t("nav_history")}
          tracks={$userDataStore.history}
          isLoading={$userDataStore.isLoading}
          emptyMessage="Belum ada riwayat putar. Mulai putar lagu untuk mencatat riwayat."
        />
      {:else if activeTab === "playlist"}
        {@const selectedPl = $userDataStore.playlists.find(
          (p) => p.id === $activePlaylistId,
        )}
        {#if selectedPl}
          <PlaylistView
            playlist={selectedPl}
            on:close={() => (activeTab = "search")}
          />
        {:else}
          <div class="p-8 text-center text-zinc-500 text-sm">
            Playlist tidak ditemukan.
          </div>
        {/if}
      {/if}
    </main>
  </div>

  <!-- Ultra-Premium Queue Drawer -->
  <QueueDrawer />

  <!-- Bottom Floating Island Player -->
  <PlayerBar />

  <!-- Settings & About Modal -->
  <SettingsModal />

  <!-- Playlist Modals -->
  <PlaylistModal />
  <AddToPlaylistModal />
  <PlaylistShareModal />

  <!-- Toast Notification -->
  {#if $toastMessage}
    <div
      class="fixed top-12 right-5 z-50 pointer-events-none animate-slide-up duration-200"
    >
      <div
        class="px-4 py-2.5 rounded-xl bg-[#181A24]/95 border border-white/10 shadow-2xl text-xs font-medium text-white flex items-center gap-2.5 backdrop-blur-md"
      >
        <span
          class="w-2 h-2 rounded-full bg-crimson-500 shadow-[0_0_8px_#E61E38]"
        ></span>
        <span>{$toastMessage}</span>
      </div>
    </div>
  {/if}

  <!-- Automatic Background Update Available Banner -->
  {#if $showUpdateBanner && $updateState?.hasUpdate}
    <aside
      aria-label="Notifikasi Pembaruan"
      class="fixed top-14 right-5 z-50 animate-spring-up duration-300"
    >
      <div
        class="w-[390px] p-4 rounded-2xl bg-[#12131D]/95 border border-crimson-500/40 shadow-2xl text-xs text-white backdrop-blur-xl space-y-3"
      >
        <!-- Top Row: Icon, Titles & Close button -->
        <div class="flex items-start justify-between gap-2.5">
          <div class="flex items-start gap-2.5 min-w-0">
            <div
              class="w-8 h-8 rounded-xl bg-crimson-500/20 text-crimson-400 flex items-center justify-center shrink-0 border border-crimson-500/30 mt-0.5"
            >
              <svg
                class="w-4 h-4 animate-bounce"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                viewBox="0 0 24 24"
              >
                <path d="M12 4v12m0 0l-4-4m4 4l4-4M4 20h16" />
              </svg>
            </div>
            <div class="min-w-0">
              <div class="font-bold text-zinc-100 flex items-center gap-1.5 flex-wrap">
                <span class="text-xs">Versi Baru Tersedia!</span>
                <span
                  class="text-[10px] px-1.5 py-0.5 rounded bg-crimson-500/20 text-crimson-300 font-mono font-semibold"
                  >{$updateState.latestVersion}</span
                >
                <span
                  class="text-[9px] px-1.5 py-0.5 rounded {$updateState.isInstalled ? 'bg-sky-500/20 text-sky-300 border border-sky-500/30' : 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30'} font-mono"
                >
                  {$updateState.isInstalled ? 'Installer' : 'Portabel'}
                </span>
              </div>
              <p class="text-[11px] text-zinc-400 mt-1 leading-snug">
                {#if $downloadedFilePath}
                  Pembaruan siap dipasang. Klik tombol di bawah untuk mulai ulang.
                {:else if $isDownloadingUpdate}
                  Sedang mengunduh berkas pembaruan...
                {:else}
                  Dengerin {$updateState.latestVersion} ({$updateState.isInstalled ? 'Setup .exe' : 'Portable .zip'}) siap diunduh otomatis.
                {/if}
              </p>
            </div>
          </div>
          <button
            class="p-1 rounded-lg text-zinc-400 hover:text-white hover:bg-white/10 transition cursor-pointer shrink-0 -mt-0.5 -mr-0.5"
            on:click={() => showUpdateBanner.set(false)}
            title="Tutup notifikasi"
          >
            <svg
              class="w-3.5 h-3.5"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Bottom Row: Action Buttons -->
        <div class="flex items-center justify-end gap-2 pt-2 border-t border-white/[0.08]">
          {#if $downloadedFilePath}
            <button
              class="flex-1 py-1.5 px-3 rounded-xl bg-crimson-600 hover:bg-crimson-500 text-white font-semibold transition cursor-pointer text-xs flex items-center justify-center gap-1.5 shadow-md shadow-crimson-600/30"
              on:click={installAndRestartUpdate}
              title="Pasang pembaruan dan mulai ulang aplikasi"
            >
              <span>⚡ Pasang & Mulai Ulang</span>
            </button>
            <button
              class="py-1.5 px-3 rounded-xl bg-white/10 hover:bg-white/20 text-zinc-200 font-semibold transition cursor-pointer text-xs flex items-center justify-center gap-1 shadow-sm shrink-0"
              on:click={revealDownloadedFile}
              title="Buka folder Downloads"
            >
              <span>Buka Folder</span>
            </button>
          {:else if $isDownloadingUpdate}
            <div
              class="flex items-center gap-2 text-zinc-300 text-xs py-1"
            >
              <svg
                class="w-3.5 h-3.5 animate-spin text-crimson-400"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <circle cx="12" cy="12" r="10" stroke-opacity="0.25" />
                <path d="M12 2a10 10 0 0 1 10 10" />
              </svg>
              <span>Mengunduh biner pembaruan...</span>
            </div>
          {:else}
            <button
              class="w-full py-1.5 px-3 rounded-xl bg-crimson-600 hover:bg-crimson-500 text-white font-semibold transition cursor-pointer text-xs shadow-md shadow-crimson-600/30 text-center"
              on:click={() => startDownloadUpdate()}
              title="Unduh paket pembaruan resmi"
            >
              Unduh Pembaruan ({$updateState.isInstalled ? 'Installer .exe' : 'Portable .zip'})
            </button>
          {/if}
        </div>
      </div>
    </aside>
  {/if}

  <!-- What's New Release Popup -->
  <WhatsNewModal isOpen={$isWhatsNewOpen} onClose={handleCloseWhatsNew} />
</div>
