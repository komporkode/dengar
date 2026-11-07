<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { Playlist, SearchResultItem } from "../types";
  import TrackItem from "./TrackItem.svelte";
  import { playerStore } from "../stores/playerStore";
  import { queueStore } from "../stores/queueStore";
  import { userDataStore } from "../stores/userDataStore";
  import {
    openEditPlaylistModal,
    openSharePlaylistModal,
  } from "../stores/playlistStore";
  import { showToast } from "../stores/settingsStore";
  import { t } from "../stores/i18nStore";
  import HeroBanner from "./HeroBanner.svelte";

  export let playlist: Playlist;

  const dispatch = createEventDispatcher<{
    close: void;
  }>();

  $: tracks = playlist.tracks || [];
  let filterQuery = "";
  let sortBy: "default" | "title" | "artist" | "duration" = "default";

  $: displayTracks = (() => {
    let list = [...tracks];
    if (filterQuery.trim()) {
      const q = filterQuery.toLowerCase().trim();
      list = list.filter(
        (t) =>
          t.title.toLowerCase().includes(q) ||
          (t.uploaderName && t.uploaderName.toLowerCase().includes(q)),
      );
    }
    if (sortBy === "title") {
      list.sort((a, b) => a.title.localeCompare(b.title));
    } else if (sortBy === "artist") {
      list.sort((a, b) =>
        (a.uploaderName || "").localeCompare(b.uploaderName || ""),
      );
    } else if (sortBy === "duration") {
      list.sort((a, b) => (b.duration || 0) - (a.duration || 0));
    }
    return list;
  })();

  $: totalDurationFormatted = formatTotalDuration(tracks);

  $: thumbnails = getCollageThumbnails(tracks);

  function getCollageThumbnails(items: SearchResultItem[]): string[] {
    const list: string[] = [];
    for (const t of items) {
      if (t.thumbnailUrl && !list.includes(t.thumbnailUrl)) {
        list.push(t.thumbnailUrl);
      }
      if (list.length >= 4) break;
    }
    return list;
  }

  function formatTotalDuration(items: SearchResultItem[]): string {
    if (!items || items.length === 0) return "";
    const totalSec = items.reduce((acc, item) => acc + (item.duration || 0), 0);
    if (totalSec <= 0) return "";
    const hrs = Math.floor(totalSec / 3600);
    const mins = Math.floor((totalSec % 3600) / 60);
    if (hrs > 0) {
      return `~${hrs} jam ${mins} mnt`;
    }
    return `~${mins} mnt`;
  }

  $: isCurrentListPlaying =
    tracks.length > 0 &&
    $playerStore.isPlaying &&
    tracks.some((t) => t.id === $playerStore.currentTrack?.id);

  function handlePlayAll() {
    if (tracks.length === 0) return;
    if (isCurrentListPlaying) {
      playerStore.togglePlay();
    } else {
      queueStore.setQueue(tracks, 0);
      playerStore.playTrack(tracks[0]);
    }
  }

  function handleShuffleAll() {
    if (tracks.length === 0) return;
    queueStore.enableShuffle();
    queueStore.setQueue(tracks, 0);
    const current = queueStore.getCurrent();
    if (current) {
      playerStore.playTrack(current);
    }
  }

  async function handleDeletePlaylist() {
    if (window.confirm($t("confirm_delete_playlist"))) {
      await userDataStore.deletePlaylist(playlist.id);
      showToast("Playlist berhasil dihapus");
      dispatch("close");
    }
  }

  async function handleRemoveTrack(trackId: string) {
    await userDataStore.removeTrackFromPlaylist(playlist.id, trackId);
    showToast("Lagu dihapus dari playlist");
  }
</script>

<div class="w-full flex flex-col select-none">
  <!-- Cinematic Full-Bleed Music Hero -->
  <div class="px-4 sm:px-8 pt-6 pb-2">
    <HeroBanner
      tag="DAFTAR PUTAR PRIBADI"
      title={playlist.name}
      subtitle={`${tracks.length} ${$t("items_count")}${totalDurationFormatted ? ` • ${totalDurationFormatted}` : ""}${playlist.description ? ` • ${playlist.description}` : ""}`}
      imageUrl={thumbnails[0] || ""}
      trackId={tracks[0]?.id}
      playLabel={isCurrentListPlaying ? $t("pause_all") : $t("play_all")}
      isPlaying={isCurrentListPlaying}
      showPlay={tracks.length > 0}
      showShuffle={tracks.length > 0}
      shuffleTitle={$t("shuffle")}
      on:play={handlePlayAll}
      on:shuffle={handleShuffleAll}
    >
      <button
        type="button"
        class="w-11 h-11 sm:w-12 sm:h-12 rounded-full bg-white/[0.08] hover:bg-white/[0.14] border border-white/10 hover:border-white/20 text-zinc-200 hover:text-white flex items-center justify-center transition-all hover:scale-105 active:scale-95 cursor-pointer shrink-0"
        on:click={() => openSharePlaylistModal(playlist)}
        title={$t("share_playlist")}
        aria-label={$t("share_playlist")}
      >
        <svg
          class="w-4 h-4 shrink-0"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          viewBox="0 0 24 24"
        >
          <circle cx="18" cy="5" r="3"></circle>
          <circle cx="6" cy="12" r="3"></circle>
          <circle cx="18" cy="19" r="3"></circle>
          <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line>
          <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line>
        </svg>
      </button>

      <button
        type="button"
        class="w-11 h-11 sm:w-12 sm:h-12 rounded-full bg-white/[0.08] hover:bg-white/[0.14] border border-white/10 hover:border-white/20 text-zinc-200 hover:text-white flex items-center justify-center transition-all hover:scale-105 active:scale-95 cursor-pointer shrink-0"
        on:click={() => openEditPlaylistModal(playlist)}
        title={$t("edit_playlist")}
        aria-label={$t("edit_playlist")}
      >
        <svg
          class="w-4 h-4 shrink-0"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          viewBox="0 0 24 24"
        >
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
          ></path>
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"
          ></path>
        </svg>
      </button>

      <button
        type="button"
        class="w-11 h-11 sm:w-12 sm:h-12 rounded-full bg-white/[0.08] hover:bg-red-500/20 border border-white/10 hover:border-red-500/30 text-zinc-300 hover:text-red-400 flex items-center justify-center transition-all hover:scale-105 active:scale-95 cursor-pointer shrink-0"
        on:click={handleDeletePlaylist}
        title={$t("delete_playlist")}
        aria-label={$t("delete_playlist")}
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
      </button>
    </HeroBanner>
  </div>

  <!-- Tracks Section -->
  <section class="px-3 sm:px-6 pt-4 pb-28">
    {#if tracks.length === 0}
      <div class="py-16 flex flex-col items-center justify-center text-center">
        <div
          class="w-14 h-14 rounded-2xl bg-white/[0.04] border border-white/[0.08] flex items-center justify-center mb-3 text-zinc-500"
        >
          <svg
            class="w-7 h-7"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            viewBox="0 0 24 24"
          >
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
        </div>
        <h4 class="text-sm font-bold text-zinc-300">{$t("playlist_empty")}</h4>
        <p class="text-xs text-zinc-500 mt-1 max-w-sm">
          {$t("playlist_empty_tip")}
        </p>
      </div>
    {:else}
      <!-- Filter & Sort Bar -->
      <div
        class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-2.5 mb-3"
      >
        <div class="relative flex-1 max-w-xs">
          <svg
            class="w-3.5 h-3.5 absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            viewBox="0 0 24 24"
          >
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
          </svg>
          <input
            type="text"
            placeholder={$t("search_in_playlist")}
            bind:value={filterQuery}
            class="w-full pl-9 pr-3 py-1.5 rounded-xl bg-white/[0.05] border border-white/10 text-xs text-white placeholder-zinc-500 outline-none focus:border-crimson-500 transition"
          />
        </div>
        <div class="flex items-center gap-2 text-xs">
          <span class="text-zinc-500 font-medium text-[11px]">{$t("sort_by")}</span>
          <select
            bind:value={sortBy}
            class="bg-white/[0.06] border border-white/10 text-zinc-300 text-xs rounded-xl px-2.5 py-1.5 outline-none cursor-pointer focus:border-crimson-500"
          >
            <option value="default" class="bg-[#10121A]">{$t("sort_default")}</option>
            <option value="title" class="bg-[#10121A]">{$t("sort_title")}</option>
            <option value="artist" class="bg-[#10121A]">{$t("sort_artist")}</option>
            <option value="duration" class="bg-[#10121A]"
              >{$t("sort_duration")}</option
            >
          </select>
        </div>
      </div>

      <!-- Header Bar (Ramping & Rata Kiri) -->
      <div
        class="flex items-center gap-2.5 sm:gap-3.5 px-2.5 sm:px-3 py-2 text-[11px] font-bold tracking-wider text-zinc-400 uppercase border-b border-white/[0.06] mb-1"
      >
        <div class="w-8 sm:w-9 shrink-0 text-center">#</div>
        <div class="flex-1 min-w-0">{$t("title_and_artist")}</div>
        <div class="hidden lg:block w-44 xl:w-56 shrink-0 pr-3">
          {$t("channel_or_artist")}
        </div>
        <div class="w-24 sm:w-28 shrink-0 text-right pr-1">
          {$t("quality_and_duration")}
        </div>
        <div class="w-16 sm:w-20 shrink-0 text-right pr-1">{$t("actions")}</div>
      </div>

      <!-- Tracks List -->
      <div class="space-y-1">
        {#each displayTracks as track, idx (track.id)}
          <TrackItem
            {track}
            index={idx}
            tracksContext={displayTracks}
            onRemove={() => handleRemoveTrack(track.id)}
          />
        {/each}
      </div>
    {/if}
  </section>
</div>
