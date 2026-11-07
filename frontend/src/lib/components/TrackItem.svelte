<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import type { SearchResultItem } from "../types";
  import { playerStore } from "../stores/playerStore";
  import { queueStore } from "../stores/queueStore";
  import { userDataStore } from "../stores/userDataStore";
  import { openAddToPlaylistModal } from "../stores/playlistStore";
  import { showToast, audioQualityStore } from "../stores/settingsStore";
  import { t } from "../stores/i18nStore";
  import { offlineStore, handleCoverError } from "../stores/offlineStore";

  export let track: SearchResultItem;
  export let index: number;
  export let tracksContext: SearchResultItem[] = [];
  export let isSearchContext = false;
  export let onRemove: (() => void) | null = null;
  export let removeTitle = "";
  export let isOffline = false;

  $: isCurrent = $playerStore.currentTrack?.id === track.id;
  $: isPlaying = isCurrent && $playerStore.isPlaying;
  $: isLoading = isCurrent && $playerStore.isLoading;
  $: isFav = $userDataStore.favorites.some((f) => f.id === track.id);
  $: isTrackOffline = isOffline || $offlineStore.trackIds.has(track.id);

  let showMenu = false;
  let menuX = 0;
  let menuY = 0;

  function closeMenu() {
    showMenu = false;
  }

  function handleOpenMenu(x: number, y: number) {
    window.dispatchEvent(new CustomEvent("dengar-close-all-menus"));
    const menuWidth = 220;
    const menuHeight = 250;
    menuX = Math.max(12, Math.min(x, window.innerWidth - menuWidth - 16));
    menuY = Math.max(12, Math.min(y, window.innerHeight - menuHeight - 16));
    showMenu = true;
  }

  function handleMoreClick(e: MouseEvent) {
    e.stopPropagation();
    if (showMenu) {
      showMenu = false;
      return;
    }
    const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
    handleOpenMenu(rect.left - 180, rect.bottom + 6);
  }

  function handleContextMenu(e: MouseEvent) {
    e.preventDefault();
    e.stopPropagation();
    handleOpenMenu(e.clientX, e.clientY);
  }

  function handlePlay() {
    closeMenu();
    if (isCurrent) {
      playerStore.togglePlay();
    } else {
      if (isSearchContext) {
        queueStore.playNow(track);
      } else if (tracksContext && tracksContext.length > 0) {
        queueStore.setQueue(tracksContext, index);
      } else {
        queueStore.playNow(track);
      }
      playerStore.playTrack(track);
    }
  }

  function handlePlayNext(e: MouseEvent) {
    e.stopPropagation();
    closeMenu();
    queueStore.insertNext(track);
    showToast($t("toast_will_play_next"));
  }

  function handleSearchArtist(e: MouseEvent) {
    e.stopPropagation();
    closeMenu();
    if (!track.uploaderName) return;
    window.dispatchEvent(
      new CustomEvent("dengar-search-artist", {
        detail: { artist: track.uploaderName },
      }),
    );
  }

  function handleAddToQueue(e: MouseEvent) {
    e.stopPropagation();
    closeMenu();
    queueStore.addToQueue(track);
    showToast($t("added_to_queue"));
  }

  function handleToggleFav(e: MouseEvent) {
    e.stopPropagation();
    const wasFav = isFav;
    userDataStore.toggleFavorite(track);
    showToast(wasFav ? $t("removed_from_favorites") : $t("added_to_favorites"));
  }

  function handleAddToPlaylist(e: MouseEvent) {
    e.stopPropagation();
    closeMenu();
    openAddToPlaylistModal(track);
  }

  async function handleSaveOffline(e: MouseEvent) {
    e.stopPropagation();
    closeMenu();
    try {
      showToast(`${$t("saving_offline")} "${track.title}"`);
      await offlineStore.saveTrack(track);
    } catch (err: any) {
      showToast(err?.message || "Gagal menyimpan offline");
    }
  }

  async function handleToggleOffline(e: MouseEvent) {
    e.stopPropagation();
    closeMenu();
    if (isTrackOffline) {
      try {
        await offlineStore.removeTrack(track.id);
        showToast(`"${track.title}" ${$t("offline_delete_toast")}`);
      } catch {
        showToast($t("offline_delete_fail"));
      }
    } else {
      await handleSaveOffline(e);
    }
  }

  async function handleCopyLink(e: MouseEvent) {
    e.stopPropagation();
    closeMenu();
    try {
      const url = `https://www.youtube.com/watch?v=${track.id}`;
      await navigator.clipboard.writeText(url);
      showToast($t("toast_link_copied"));
    } catch {
      showToast("Gagal menyalin tautan");
    }
  }

  function onGlobalCloseMenus() {
    showMenu = false;
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === "Escape" && showMenu) {
      showMenu = false;
    }
  }

  onMount(() => {
    window.addEventListener("dengar-close-all-menus", onGlobalCloseMenus);
    window.addEventListener("keydown", handleKeyDown);
  });

  onDestroy(() => {
    window.removeEventListener("dengar-close-all-menus", onGlobalCloseMenus);
    window.removeEventListener("keydown", handleKeyDown);
  });
</script>

<div
  class="track-row group flex items-center gap-2.5 sm:gap-3.5 px-2.5 sm:px-3 py-2 sm:py-2.5 rounded-2xl transition-all duration-150 cursor-pointer active:scale-[0.99] select-none {isCurrent
    ? 'bg-crimson-500/[0.08] border border-crimson-500/25 hover:bg-crimson-500/[0.12]'
    : 'hover:bg-white/[0.04] border border-transparent'}"
  on:click={handlePlay}
  on:contextmenu={handleContextMenu}
  role="button"
  tabindex="0"
  on:keydown={(e) => e.key === "Enter" && handlePlay()}
>
  <!-- Col 1: Ramping & Rata Kiri (#) / Play Indicator (32px) -->
  <div
    class="w-8 sm:w-9 shrink-0 flex items-center justify-center text-center relative"
  >
    {#if isPlaying}
      <span
        class="track-index text-crimson-400 font-mono text-sm font-bold flex items-center gap-0.5"
      >
        <span
          class="inline-block w-0.5 sm:w-1 h-3.5 bg-crimson-500 rounded-sm animate-pulse"
        ></span>
        <span
          class="inline-block w-0.5 sm:w-1 h-4 bg-crimson-400 rounded-sm animate-bounce"
        ></span>
        <span
          class="inline-block w-0.5 sm:w-1 h-2 bg-crimson-500 rounded-sm animate-pulse"
        ></span>
      </span>
      <button
        type="button"
        class="play-hover-btn absolute inset-0 m-auto w-7 h-7 rounded-full bg-crimson-500 text-white flex items-center justify-center opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all shadow-md cursor-pointer active:scale-90"
        on:click|stopPropagation={handlePlay}
        title="Jeda Lagu"
        aria-label="Jeda Lagu"
      >
        <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
          <rect height="16" width="4" x="6" y="4"></rect>
          <rect height="16" width="4" x="14" y="4"></rect>
        </svg>
      </button>
    {:else if isLoading}
      <div
        class="w-4 h-4 border-2 border-crimson-500 border-t-transparent rounded-full animate-spin"
      ></div>
    {:else}
      <span
        class="track-index text-xs font-mono font-medium text-zinc-500 group-hover:invisible transition-all"
      >
        {index + 1}
      </span>
      <button
        type="button"
        class="play-hover-btn absolute inset-0 m-auto w-7 h-7 rounded-full bg-white/[0.12] hover:bg-crimson-500 text-white flex items-center justify-center opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all shadow-md cursor-pointer active:scale-90"
        on:click|stopPropagation={handlePlay}
        title="Putar Lagu"
        aria-label="Putar Lagu"
      >
        <svg class="w-3.5 h-3.5 fill-current ml-0.5" viewBox="0 0 24 24">
          <polygon points="5 3 19 12 5 21 5 3"></polygon>
        </svg>
      </button>
    {/if}
  </div>

  <!-- Col 2: Thumbnail + Judul Lagu & Artis (Menempel langsung rapat ke kiri) -->
  <div class="flex-1 min-w-0 flex items-center gap-2.5 sm:gap-3">
    <!-- Album Artwork Thumbnail -->
    <div
      class="relative w-10 h-10 sm:w-11 sm:h-11 shrink-0 rounded-xl overflow-hidden bg-black/40 shadow-xs group-hover:shadow-md transition-shadow"
    >
      {#if track.thumbnailUrl}
        <img
          src={track.thumbnailUrl}
          alt={track.title}
          class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          loading="lazy"
          on:error={(e) => handleCoverError(e, track.id)}
        />
      {:else}
        <div
          class="w-full h-full flex items-center justify-center bg-zinc-800 text-zinc-500"
        >
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path
              d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"
            />
          </svg>
        </div>
      {/if}
    </div>

    <!-- Title & Artist Metadata -->
    <div class="min-w-0 flex-1">
      <p
        class="text-xs sm:text-sm font-semibold truncate {isCurrent
          ? 'text-crimson-400 font-bold'
          : 'text-zinc-200 group-hover:text-white'}"
        title={track.title}
      >
        {track.title}
      </p>
      <div
        class="flex items-center gap-1.5 text-[11px] text-zinc-400 truncate mt-0.5"
      >
        <button
          type="button"
          class="truncate cursor-pointer hover:underline hover:text-crimson-400 transition text-left p-0 bg-transparent border-0"
          on:click|stopPropagation={handleSearchArtist}
          title={`${$t("search_artist_songs")} ${track.uploaderName}`}
        >
          {track.uploaderName}
        </button>
        {#if isTrackOffline}
          <span
            class="inline-flex items-center text-emerald-400 shrink-0"
            title={$t("track_offline_saved")}
          >
            <svg
              class="w-3 h-3"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              viewBox="0 0 24 24"
            >
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
          </span>
        {/if}
      </div>
    </div>
  </div>

  <!-- Col 3: Channel / Artist (Optional di layar lebar) -->
  <div
    class="hidden lg:flex w-44 xl:w-56 shrink-0 items-center text-xs text-zinc-400 truncate pr-3"
  >
    <button
      type="button"
      class="truncate cursor-pointer hover:underline hover:text-crimson-400 transition text-left p-0 bg-transparent border-0"
      on:click|stopPropagation={handleSearchArtist}
      title={`${$t("search_artist_songs")} ${track.uploaderName}`}
    >
      {track.uploaderName}
    </button>
  </div>

  <!-- Col 4: Kualitas & Durasi -->
  <div
    class="w-24 sm:w-28 shrink-0 flex items-center justify-end gap-1.5 text-right font-mono whitespace-nowrap"
  >
    <span
      class="hidden sm:inline-block text-[10px] font-mono px-1.5 py-0.5 rounded {isCurrent
        ? $audioQualityStore === 'saver'
          ? 'bg-emerald-500/10 text-emerald-300 border border-emerald-500/20'
          : 'bg-amber-500/10 text-amber-300 border border-amber-500/20'
        : 'bg-white/[0.06] text-zinc-300'}"
    >
      {$audioQualityStore === "saver"
        ? "ECO Audio"
        : track.duration > 300
          ? "HQ Opus"
          : "HQ Audio"}
    </span>
    <span
      class="text-[11px] sm:text-xs font-mono {isCurrent
        ? 'text-crimson-400 font-bold'
        : 'text-zinc-400'}"
    >
      {track.durationFormatted || "00:00"}
    </span>
  </div>

  <!-- Col 5: Menu Aksi Minimalis (Hanya Love + Titik Tiga ⋮) -->
  <div
    class="w-16 sm:w-20 shrink-0 flex items-center justify-end gap-1 text-zinc-400"
  >
    <!-- Favorite Love Button -->
    <button
      type="button"
      class="w-8 h-8 rounded-lg hover:bg-white/[0.08] flex items-center justify-center transition cursor-pointer"
      on:click={handleToggleFav}
      title={isFav ? $t("removed_from_favorites") : $t("added_to_favorites")}
      aria-label={isFav ? $t("removed_from_favorites") : $t("added_to_favorites")}
    >
      <svg
        class="w-4 h-4 {isFav
          ? 'fill-crimson-500 text-crimson-500'
          : 'text-zinc-400 hover:text-white'}"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        viewBox="0 0 24 24"
      >
        <path
          d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"
        ></path>
      </svg>
    </button>

    <!-- Context Menu Trigger Button (⋮) -->
    <button
      type="button"
      class="w-8 h-8 rounded-lg hover:bg-white/[0.08] hover:text-white flex items-center justify-center transition cursor-pointer {showMenu
        ? 'bg-white/[0.12] text-white'
        : 'text-zinc-400'}"
      on:click={handleMoreClick}
      title={$t("track_options_menu")}
      aria-label={$t("track_options_menu")}
    >
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
        <circle cx="12" cy="5" r="1.75"></circle>
        <circle cx="12" cy="12" r="1.75"></circle>
        <circle cx="12" cy="19" r="1.75"></circle>
      </svg>
    </button>
  </div>
</div>

<!-- Floating Popover Context Menu -->
{#if showMenu}
  <!-- Backdrop click dismisser -->
  <button
    type="button"
    class="fixed inset-0 z-40 bg-transparent border-0 cursor-default"
    on:click|stopPropagation={closeMenu}
    on:contextmenu|preventDefault|stopPropagation={closeMenu}
    aria-label="Tutup menu"
  ></button>

  <!-- Popover Context Menu Window -->
  <div
    class="fixed z-50 w-56 rounded-2xl bg-[#121422]/95 backdrop-blur-2xl border border-white/10 shadow-2xl p-1.5 select-none animate-in fade-in zoom-in-95 duration-100 flex flex-col gap-0.5 text-xs font-medium text-zinc-200"
    style="left: {menuX}px; top: {menuY}px;"
    on:click|stopPropagation
    on:keydown|stopPropagation
    role="menu"
    tabindex="-1"
  >
    <!-- Track Title in Menu Header -->
    <div class="px-2.5 py-2 border-b border-white/[0.08] mb-1">
      <p class="font-bold text-white text-xs truncate">{track.title}</p>
      <p class="text-[10px] text-zinc-400 truncate mt-0.5">
        {track.uploaderName}
      </p>
    </div>

    <!-- Play Now -->
    <button
      type="button"
      class="w-full flex items-center gap-2.5 px-2.5 py-2 rounded-xl hover:bg-white/[0.08] hover:text-white transition cursor-pointer text-left"
      on:click|stopPropagation={handlePlay}
    >
      <svg
        class="w-4 h-4 text-crimson-400 shrink-0"
        fill="currentColor"
        viewBox="0 0 24 24"
      >
        {#if isPlaying}
          <rect height="16" width="4" x="6" y="4"></rect>
          <rect height="16" width="4" x="14" y="4"></rect>
        {:else}
          <polygon points="5 3 19 12 5 21 5 3"></polygon>
        {/if}
      </svg>
      <span>{isPlaying ? "Jeda Pemutaran" : $t("play_now")}</span>
    </button>

    <!-- Play Next -->
    <button
      type="button"
      class="w-full flex items-center gap-2.5 px-2.5 py-2 rounded-xl hover:bg-white/[0.08] hover:text-white transition cursor-pointer text-left"
      on:click|stopPropagation={handlePlayNext}
    >
      <svg
        class="w-4 h-4 text-zinc-400 shrink-0"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        viewBox="0 0 24 24"
      >
        <polygon points="5 4 15 12 5 20 5 4"></polygon>
        <line x1="19" y1="5" x2="19" y2="19"></line>
      </svg>
      <span>{$t("play_next")}</span>
    </button>

    {#if track.uploaderName}
      <!-- Search Artist -->
      <button
        type="button"
        class="w-full flex items-center gap-2.5 px-2.5 py-2 rounded-xl hover:bg-white/[0.08] hover:text-white transition cursor-pointer text-left"
        on:click|stopPropagation={handleSearchArtist}
      >
        <svg
          class="w-4 h-4 text-zinc-400 shrink-0"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          viewBox="0 0 24 24"
        >
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
        </svg>
        <span>{$t("search_artist_songs")}</span>
      </button>
    {/if}

    <!-- Add to Queue -->
    <button
      type="button"
      class="w-full flex items-center gap-2.5 px-2.5 py-2 rounded-xl hover:bg-white/[0.08] hover:text-white transition cursor-pointer text-left"
      on:click|stopPropagation={handleAddToQueue}
    >
      <svg
        class="w-4 h-4 text-zinc-400 shrink-0"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        viewBox="0 0 24 24"
      >
        <line x1="12" y1="5" x2="12" y2="19"></line>
        <line x1="5" y1="12" x2="19" y2="12"></line>
      </svg>
      <span>{$t("added_to_queue")}</span>
    </button>

    <!-- Add to Playlist -->
    <button
      type="button"
      class="w-full flex items-center gap-2.5 px-2.5 py-2 rounded-xl hover:bg-white/[0.08] hover:text-white transition cursor-pointer text-left"
      on:click|stopPropagation={handleAddToPlaylist}
    >
      <svg
        class="w-4 h-4 text-zinc-400 shrink-0"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        viewBox="0 0 24 24"
      >
        <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"></path>
      </svg>
      <span>{$t("add_to_playlist")}</span>
    </button>

    <!-- Save / Remove Offline -->
    <button
      type="button"
      class="w-full flex items-center gap-2.5 px-2.5 py-2 rounded-xl hover:bg-white/[0.08] hover:text-white transition cursor-pointer text-left"
      on:click|stopPropagation={handleToggleOffline}
    >
      {#if isTrackOffline}
        <svg
          class="w-4 h-4 text-emerald-400 shrink-0"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          viewBox="0 0 24 24"
        >
          <polyline points="20 6 9 17 4 12"></polyline>
        </svg>
        <span class="text-emerald-400">{$t("remove_offline")}</span>
      {:else}
        <svg
          class="w-4 h-4 text-zinc-400 shrink-0"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          viewBox="0 0 24 24"
        >
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
          <polyline points="7 10 12 15 17 10"></polyline>
          <line x1="12" y1="15" x2="12" y2="3"></line>
        </svg>
        <span>{$t("btn_save_offline")}</span>
      {/if}
    </button>

    <!-- Copy Track Link -->
    <button
      type="button"
      class="w-full flex items-center gap-2.5 px-2.5 py-2 rounded-xl hover:bg-white/[0.08] hover:text-white transition cursor-pointer text-left"
      on:click|stopPropagation={handleCopyLink}
    >
      <svg
        class="w-4 h-4 text-zinc-400 shrink-0"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        viewBox="0 0 24 24"
      >
        <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"
        ></path>
        <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"
        ></path>
      </svg>
      <span>{$t("copy_link")}</span>
    </button>

    <!-- Delete / Remove if provided (e.g. from Playlist / History) -->
    {#if onRemove}
      <div class="h-[1px] bg-white/[0.08] my-1"></div>
      <button
        type="button"
        class="w-full flex items-center gap-2.5 px-2.5 py-2 rounded-xl hover:bg-red-500/15 text-red-400 hover:text-red-300 transition cursor-pointer text-left"
        on:click|stopPropagation={() => {
          closeMenu();
          onRemove && onRemove();
        }}
      >
        <svg
          class="w-4 h-4 shrink-0"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          viewBox="0 0 24 24"
        >
          <polyline points="3 6 5 6 21 6"></polyline>
          <path
            d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
          ></path>
        </svg>
        <span>{removeTitle || $t("remove_from_list")}</span>
      </button>
    {/if}
  </div>
{/if}
