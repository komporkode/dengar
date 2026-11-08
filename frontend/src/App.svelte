<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import type {
    NavigationTab,
    SearchResultItem,
    SearchFilter,
  } from "./lib/types";
  import { Search, SearchMore, ReloadInstances } from "../wailsjs/go/main/App";
  import {
    Quit,
    WindowHide,
    WindowMinimise,
    WindowToggleMaximise,
    EventsOn,
  } from "../wailsjs/runtime/runtime.js";
  import { playerStore } from "./lib/stores/playerStore";
  import { userDataStore } from "./lib/stores/userDataStore";
  import { offlineStore } from "./lib/stores/offlineStore";
  import {
    isSettingsOpen,
    isMiniPlayer,
    toastData,
    toastMessage,
    togglePinMiniPlayer,
    exitMiniPlayer,
    enterMiniPlayer,
    isTransitioning,
    isWhatsNewOpen,
    closeActionStore,
    isCloseConfirmOpen,
    showToast,
  } from "./lib/stores/settingsStore";
  import {
    updateState,
    showUpdateBanner,
    isDownloadingUpdate,
    downloadProgressData,
    downloadedFilePath,
    checkForUpdatesInBackground,
    startDownloadUpdate,
    revealDownloadedFile,
    installAndRestartUpdate,
    formatBytes,
  } from "./lib/stores/updateStore";
  import { t } from "./lib/stores/i18nStore";
  import { CURRENT_APP_VERSION } from "./lib/data/changelogs";

  import Sidebar from "./lib/components/Sidebar.svelte";
  import SearchBar from "./lib/components/SearchBar.svelte";
  import TrackList from "./lib/components/TrackList.svelte";
  import ForYouView from "./lib/components/ForYouView.svelte";
  import OfflineView from "./lib/components/OfflineView.svelte";
  import PlayerBar from "./lib/components/PlayerBar.svelte";
  import QueueDrawer from "./lib/components/QueueDrawer.svelte";
  import SettingsView from "./lib/components/SettingsView.svelte";
  import MiniPlayer from "./lib/components/MiniPlayer.svelte";
  import HomeSections from "./lib/components/HomeSections.svelte";
  import PlaylistView from "./lib/components/PlaylistView.svelte";
  import PlaylistModal from "./lib/components/PlaylistModal.svelte";
  import AddToPlaylistModal from "./lib/components/AddToPlaylistModal.svelte";
  import PlaylistShareModal from "./lib/components/PlaylistShareModal.svelte";
  import WhatsNewModal from "./lib/components/WhatsNewModal.svelte";
  import CloseConfirmModal from "./lib/components/CloseConfirmModal.svelte";
  import SleepTimerModal from "./lib/components/SleepTimerModal.svelte";
  import NetworkBanner from "./lib/components/NetworkBanner.svelte";
  import { networkStore } from "./lib/stores/networkStore";
  import { activePlaylistId } from "./lib/stores/playlistStore";
  import { mediaSessionService } from "./lib/services/mediaSessionService";

  let activeTab: NavigationTab = "search";
  let previousTab: NavigationTab = "search";

  function handleTabChange(e: CustomEvent<NavigationTab>) {
    const nextTab = e.detail;
    if (activeTab === "search" && nextTab === "search" && searchQuery) {
      searchQuery = "";
      if (popularTracks.length === 0) {
        loadPopularTracks();
      }
    }
    if (activeTab !== "settings") {
      previousTab = activeTab;
    }
    activeTab = nextTab;
    if (nextTab !== "settings") {
      isSettingsOpen.set(false);
    }
  }

  function handleCloseSettings() {
    isSettingsOpen.set(false);
    activeTab =
      previousTab && previousTab !== "settings" ? previousTab : "search";
  }

  $: if ($isSettingsOpen && activeTab !== "settings") {
    previousTab = activeTab;
    activeTab = "settings";
  }
  let searchQuery = "";
  let currentFilter: SearchFilter = "music_songs";

  // Dedicated popular / trending tracks (protected from search pollution)
  let popularTracks: SearchResultItem[] = [];
  let popularContinuationToken: string = "";
  let isPopularLoading = false;
  let popularError: string | null = null;

  // Search results
  let searchResults: SearchResultItem[] = [];
  let searchContinuationToken: string = "";
  let isSearching = false;
  let searchError: string | null = null;

  let isLoadingMore = false;

  let scrollTimer: any;
  let handleScrollActivity: (e: Event) => void;

  onMount(async () => {
    // 0. Initialize hardware media keys and OS SMTC
    mediaSessionService.init();

    // 0b. Listen to system tray action events
    EventsOn("tray:play_pause", () => {
      playerStore.togglePlay();
    });
    EventsOn("tray:next", () => {
      playerStore.playNext();
    });
    EventsOn("tray:prev", () => {
      playerStore.playPrev();
    });

    // 1. Load user favorites, history & offline tracks from backend JSON storage
    await Promise.all([userDataStore.load(), offlineStore.load()]);

    // 2. Load popular / trending tracks so app is instantly ready with tracks
    loadPopularTracks();

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
    const CURRENT_VERSION = CURRENT_APP_VERSION;
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

    // 7. Global listener for artist search clicks
    window.addEventListener("dengar-search-artist", handleSearchArtistEvent);
  });

  function handleSearchArtistEvent(e: any) {
    const artist = e.detail?.artist;
    if (artist) {
      activeTab = "search";
      searchQuery = artist;
      performSearch(artist, "music_songs");
    }
  }

  function handleCloseWhatsNew() {
    isWhatsNewOpen.set(false);
    try {
      localStorage.setItem("dengerin_seen_version", CURRENT_APP_VERSION);
    } catch (e) {
      console.warn("Save version warning:", e);
    }
  }

  onDestroy(() => {
    window.removeEventListener("keydown", handleGlobalKeydown);
    window.removeEventListener("dengar-search-artist", handleSearchArtistEvent);
    if (handleScrollActivity) {
      window.removeEventListener("scroll", handleScrollActivity, {
        capture: true,
      });
    }
    clearTimeout(scrollTimer);
  });

  async function loadPopularTracks(force = false) {
    if (popularTracks.length > 0 && !force) return;
    if (!$networkStore.isOnline) {
      popularError =
        "Perangkat sedang offline. Sambungkan internet untuk memuat lagu populer.";
      return;
    }
    isPopularLoading = true;
    popularError = null;
    try {
      const fallbackQuery =
        currentFilter === "videos" ? "Trending Music Video" : "Trending Music";
      const res = await Search(fallbackQuery, currentFilter);
      popularTracks = res?.items || [];
      popularContinuationToken = res?.continuationToken || "";
    } catch (err: any) {
      console.error("Load popular tracks failed:", err);
      popularError =
        err?.message || "Gagal memuat lagu populer. Periksa koneksi internet.";
    } finally {
      isPopularLoading = false;
    }
  }

  async function performSearch(query: string, filter: SearchFilter) {
    const q = query.trim();
    if (!q) {
      searchResults = [];
      searchContinuationToken = "";
      return;
    }
    if (!$networkStore.isOnline) {
      searchError =
        "Perangkat sedang offline. Sambungkan internet untuk mencari lagu.";
      searchResults = [];
      searchContinuationToken = "";
      return;
    }
    isSearching = true;
    searchError = null;
    try {
      const res = await Search(q, filter);
      searchResults = res?.items || [];
      searchContinuationToken = res?.continuationToken || "";
    } catch (err: any) {
      console.error("Search failed:", err);
      searchError =
        err?.message || "Gagal mencari lagu. Periksa koneksi internet.";
      searchResults = [];
      searchContinuationToken = "";
    } finally {
      isSearching = false;
    }
  }

  async function handleLoadMore() {
    if (isLoadingMore) return;
    if (searchQuery.trim()) {
      if (!searchContinuationToken || isSearching) return;
      isLoadingMore = true;
      try {
        const res = await SearchMore(searchContinuationToken);
        if (res && res.items && res.items.length > 0) {
          const existingIds = new Set(searchResults.map((i) => i.id));
          const fresh = res.items.filter((i) => !existingIds.has(i.id));
          searchResults = [...searchResults, ...fresh];
          searchContinuationToken = res.continuationToken || "";
        } else {
          searchContinuationToken = "";
        }
      } catch (err: any) {
        console.warn("SearchMore failed:", err);
      } finally {
        isLoadingMore = false;
      }
    } else {
      if (!popularContinuationToken || isPopularLoading) return;
      isLoadingMore = true;
      try {
        const res = await SearchMore(popularContinuationToken);
        if (res && res.items && res.items.length > 0) {
          const existingIds = new Set(popularTracks.map((i) => i.id));
          const fresh = res.items.filter((i) => !existingIds.has(i.id));
          popularTracks = [...popularTracks, ...fresh];
          popularContinuationToken = res.continuationToken || "";
        } else {
          popularContinuationToken = "";
        }
      } catch (err: any) {
        console.warn("SearchMore popular failed:", err);
      } finally {
        isLoadingMore = false;
      }
    }
  }

  let isRetrying = false;
  let homeSectionsRef: any = null;

  async function handleRetryConnection() {
    if (isRetrying) return;
    isRetrying = true;
    try {
      await ReloadInstances();
      searchError = null;
      popularError = null;
      if (searchQuery.trim()) {
        await performSearch(searchQuery.trim(), currentFilter);
      } else {
        await loadPopularTracks(true);
      }
      if (homeSectionsRef?.reload) {
        homeSectionsRef.reload();
      }
      showToast($t("toast_instances_reloaded"));
    } catch (e: any) {
      console.error("Retry failed:", e);
      if (searchQuery.trim()) {
        searchError = e?.message || $t("instance_connect_error");
      } else {
        popularError = e?.message || $t("instance_connect_error");
      }
    } finally {
      isRetrying = false;
    }
  }

  function handleSearchEvent(
    e: CustomEvent<{ query: string; filter: SearchFilter }>,
  ) {
    searchQuery = e.detail.query;
    currentFilter = e.detail.filter;
    if (!searchQuery.trim()) {
      searchResults = [];
      searchContinuationToken = "";
      if (popularTracks.length === 0) {
        loadPopularTracks();
      }
      return;
    }
    performSearch(searchQuery.trim(), currentFilter);
  }

  function handleWindowClose() {
    try {
      const action = $closeActionStore;
      if (action === "ask") {
        isCloseConfirmOpen.set(true);
      } else if (action === "tray") {
        WindowHide();
        showToast("Dengerin tetap berputar di background.");
      } else {
        Quit();
      }
    } catch (e) {
      console.warn("Wails Window close error:", e);
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
  <!-- Seamless Custom Frameless Titlebar -->
  <header
    class="h-9 px-4 flex items-center justify-between bg-transparent text-zinc-400 text-xs shrink-0 select-none z-50 transition-colors"
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
          >{$playerStore.currentTrack.title} • {$playerStore.currentTrack
            .uploaderName}</span
        >
      {:else}
        <span class="text-zinc-500">{$t("master_audio_engine")}</span>
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
        title={$t("window_minimize")}
        aria-label={$t("window_minimize")}
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
        title={$t("window_maximize")}
        aria-label={$t("window_maximize")}
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
        title={$t("window_close")}
        aria-label={$t("window_close")}
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
    <Sidebar bind:activeTab on:tabChange={handleTabChange} />

    <!-- Main Content Viewport -->
    <main
      class="flex-1 flex flex-col overflow-y-auto [scrollbar-gutter:stable] bg-gradient-to-b from-[#12141F] via-obsidian to-obsidian pb-32"
    >
      {#if activeTab === "search"}
        <SearchBar
          bind:query={searchQuery}
          isLoading={isSearching}
          on:search={handleSearchEvent}
        />

        {#if (searchQuery.trim() ? searchError : popularError)}
          <div
            class="mx-4 sm:mx-8 mt-4 flex items-center justify-between gap-3 px-4 sm:px-5 py-3 rounded-2xl bg-crimson-500/10 border border-crimson-500/30 text-crimson-300 text-xs shadow-lg"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div
                class="w-8 h-8 rounded-full bg-crimson-500/20 text-crimson-400 flex items-center justify-center shrink-0"
              >
                <svg
                  class="w-4 h-4"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="8" x2="12"></line>
                  <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
              </div>
              <div class="min-w-0">
                <p class="font-semibold text-white text-xs sm:text-sm">
                  {$t("instance_connect_error")}
                </p>
                <p class="text-[11px] text-zinc-300 truncate">
                  {searchQuery.trim() ? searchError : popularError}
                </p>
              </div>
            </div>

            <button
              type="button"
              class="h-9 px-4 rounded-xl bg-crimson-600 hover:bg-crimson-500 text-white font-bold text-xs inline-flex items-center gap-2 shadow-md shadow-crimson-600/30 transition-all cursor-pointer whitespace-nowrap shrink-0 disabled:opacity-50"
              on:click={handleRetryConnection}
              disabled={isRetrying}
            >
              <svg
                class="w-3.5 h-3.5 {isRetrying ? 'animate-spin' : ''}"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                viewBox="0 0 24 24"
              >
                <path
                  d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"
                />
              </svg>
              <span
                >{isRetrying
                  ? $t("refreshing_instances")
                  : $t("retry_action")}</span
              >
            </button>
          </div>
        {/if}

        {#if !searchQuery.trim()}
          <HomeSections bind:this={homeSectionsRef} />
        {/if}

        <TrackList
          title={searchQuery.trim()
            ? `${$t("search_results_prefix")}“${searchQuery.trim()}”`
            : $t("popular_songs")}
          tracks={searchQuery.trim() ? searchResults : popularTracks}
          isLoading={searchQuery.trim() ? isSearching : isPopularLoading}
          isSearchContext={!!searchQuery.trim()}
          hasMore={searchQuery.trim()
            ? !!searchContinuationToken
            : !!popularContinuationToken}
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
          showBanner={true}
          bannerTag="KOLEKSI PRIBADI • FAVORITES"
          bannerTitle="Lagu Favorit Kamu"
          bannerSubtitle="Semua alunan musik terbaik yang kamu sukai tersimpan rapi di sini"
          emptyMessage="Belum ada lagu yang ditandai favorit. Klik ikon hati pada lagu untuk menyimpannya!"
        />
      {:else if activeTab === "history"}
        <TrackList
          title={$t("nav_history")}
          tracks={$userDataStore.history}
          isLoading={$userDataStore.isLoading}
          emptyMessage={$t("history_empty_desc")}
          isHistory={true}
          onRemoveItem={(track) => {
            userDataStore.removeHistory(track.id);
            showToast($t("toast_history_removed"));
          }}
          onClearAll={() => {
            userDataStore.clearHistory();
            showToast($t("toast_history_cleared"));
          }}
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
      {:else if activeTab === "offline"}
        <OfflineView />
      {:else if activeTab === "settings"}
        <SettingsView on:close={handleCloseSettings} />
      {/if}
    </main>
  </div>

  <!-- Ultra-Premium Queue Drawer -->
  <QueueDrawer />

  <!-- Bottom Floating Island Player -->
  <PlayerBar />

  <!-- Playlist Modals -->
  <PlaylistModal />
  <AddToPlaylistModal />
  <PlaylistShareModal />
  <SleepTimerModal />

  <!-- Toast Notification with optional Undo -->
  {#if $toastData}
    <div
      class="fixed top-12 right-5 z-50 pointer-events-auto animate-slide-up duration-200"
    >
      <div
        class="px-4 py-2.5 rounded-xl bg-[#181A24]/95 border border-white/10 shadow-2xl text-xs font-medium text-white flex items-center gap-3 backdrop-blur-md"
      >
        <span
          class="w-2 h-2 rounded-full bg-crimson-500 shadow-[0_0_8px_#E61E38] shrink-0"
        ></span>
        <span>{$toastData.message}</span>
        {#if $toastData.undoAction}
          <button
            type="button"
            class="ml-1 px-2 py-0.5 rounded-md bg-white/10 hover:bg-white/20 text-crimson-400 font-bold text-[11px] transition cursor-pointer"
            on:click={() => {
              if ($toastData?.undoAction) $toastData.undoAction();
              toastData.set(null);
            }}
          >
            Urungkan
          </button>
        {/if}
      </div>
    </div>
  {:else if $toastMessage}
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
              <div
                class="font-bold text-zinc-100 flex items-center gap-1.5 flex-wrap"
              >
                <span class="text-xs">Versi Baru Tersedia!</span>
                <span
                  class="text-[10px] px-1.5 py-0.5 rounded bg-crimson-500/20 text-crimson-300 font-mono font-semibold"
                  >{$updateState.latestVersion}</span
                >
                <span
                  class="text-[9px] px-1.5 py-0.5 rounded {$updateState.isInstalled
                    ? 'bg-sky-500/20 text-sky-300 border border-sky-500/30'
                    : 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30'} font-mono"
                >
                  {$updateState.isInstalled ? "Installer" : "Portabel"}
                </span>
              </div>
              <p class="text-[11px] text-zinc-400 mt-1 leading-snug">
                {#if $downloadedFilePath}
                  Pembaruan siap dipasang. Klik tombol di bawah untuk mulai
                  ulang.
                {:else if $isDownloadingUpdate}
                  Sedang mengunduh berkas pembaruan...
                {:else}
                  Dengerin {$updateState.latestVersion} ({$updateState.isInstalled
                    ? "Setup .exe"
                    : "Portable .zip"}) siap diunduh otomatis.
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
        <div
          class="flex items-center justify-end gap-2 pt-2 border-t border-white/[0.08]"
        >
          {#if $downloadedFilePath}
            <button
              class="flex-1 py-1.5 px-3 rounded-xl bg-crimson-600 hover:bg-crimson-500 text-white font-semibold transition cursor-pointer text-xs flex items-center justify-center gap-1.5 shadow-md shadow-crimson-600/30"
              on:click={installAndRestartUpdate}
              title={$t("btn_install_restart")}
            >
              <svg
                class="w-3.5 h-3.5 shrink-0"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
              </svg>
              <span>{$t("btn_install_restart")}</span>
            </button>
            <button
              class="py-1.5 px-3 rounded-xl bg-white/10 hover:bg-white/20 text-zinc-200 font-semibold transition cursor-pointer text-xs flex items-center justify-center gap-1.5 shadow-sm shrink-0"
              on:click={revealDownloadedFile}
              title={$t("btn_open_folder")}
            >
              <svg
                class="w-3.5 h-3.5 shrink-0"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                />
              </svg>
              <span>{$t("btn_open_folder")}</span>
            </button>
          {:else if $isDownloadingUpdate}
            <div class="space-y-1.5 w-full py-1">
              <div class="flex items-center justify-between text-[11px] text-zinc-300">
                <span class="flex items-center gap-1.5 font-medium">
                  <span class="w-2 h-2 rounded-full bg-crimson-500 animate-pulse"></span>
                  Mengunduh pembaruan...
                </span>
                <span class="font-mono text-crimson-400 font-bold">
                  {$downloadProgressData ? `${$downloadProgressData.percent}%` : '0%'}
                </span>
              </div>
              <div class="h-2 w-full bg-white/10 rounded-full overflow-hidden relative">
                <div
                  class="h-full bg-gradient-to-r from-crimson-600 to-rose-500 rounded-full transition-all duration-150"
                  style="width: {$downloadProgressData ? $downloadProgressData.percent : 0}%"
                ></div>
              </div>
              <div class="flex items-center justify-between text-[10px] text-zinc-400 font-mono">
                <span>
                  {$downloadProgressData && $downloadProgressData.total > 0
                    ? `${formatBytes($downloadProgressData.downloaded)} / ${formatBytes($downloadProgressData.total)}`
                    : 'Menghubungkan...'}
                </span>
                <span>{$downloadProgressData?.speedText || ''}</span>
              </div>
            </div>
          {:else}
            <button
              class="w-full py-1.5 px-3 rounded-xl bg-crimson-600 hover:bg-crimson-500 text-white font-semibold transition cursor-pointer text-xs shadow-md shadow-crimson-600/30 text-center"
              on:click={() => startDownloadUpdate()}
              title="Unduh paket pembaruan resmi"
            >
              Unduh Pembaruan ({$updateState.isInstalled
                ? "Installer .exe"
                : "Portable .zip"})
            </button>
          {/if}
        </div>
      </div>
    </aside>
  {/if}

  <!-- What's New Release Popup -->
  <WhatsNewModal isOpen={$isWhatsNewOpen} onClose={handleCloseWhatsNew} />

  <!-- Close Confirmation Dialog -->
  <CloseConfirmModal />
</div>
