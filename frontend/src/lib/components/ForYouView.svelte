<script lang="ts">
  import { onMount } from "svelte";
  import type { SearchResultItem } from "../types";
  import { GetHomeData, RefreshHomeData } from "../../../wailsjs/go/main/App";
  import { playerStore } from "../stores/playerStore";
  import { queueStore } from "../stores/queueStore";
  import { userDataStore } from "../stores/userDataStore";
  import { t } from "../stores/i18nStore";
  import TrackList from "./TrackList.svelte";
  import HeroBanner from "./HeroBanner.svelte";

  let topPlayed: SearchResultItem[] = [];
  let recommended: SearchResultItem[] = [];
  let isLoading = true;
  let isRefreshing = false;

  // Dynamic greeting based on user's local hour
  function getGreeting(): string {
    const hour = new Date().getHours();
    if (hour >= 4 && hour < 11) return $t("greeting_morning");
    if (hour >= 11 && hour < 15) return $t("greeting_afternoon");
    if (hour >= 15 && hour < 19) return $t("greeting_evening");
    return $t("greeting_night");
  }

  // Combined deduplicated mix pool for 1-click personal playback
  $: mixTracks = [...topPlayed, ...recommended].filter(
    (item, idx, arr) => arr.findIndex((t) => t.id === item.id) === idx,
  );

  // Top featured tracks for banner and showcase
  $: sampleArtworks = (mixTracks.length > 0 ? mixTracks : topPlayed)
    .filter((t) => !!t.thumbnailUrl)
    .slice(0, 8);

  $: featuredTrack = sampleArtworks[0] || mixTracks[0] || topPlayed[0] || recommended[0];

  $: isMixPlaying =
    $playerStore.isPlaying &&
    mixTracks.some((t) => t.id === $playerStore.currentTrack?.id);

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    isLoading = true;
    try {
      const data = await GetHomeData();
      topPlayed = data?.topPlayed || [];
      recommended = data?.recommended || [];
    } catch (e) {
      console.warn("Failed to load For You data:", e);
    } finally {
      isLoading = false;
    }
  }

  async function handleRefresh() {
    if (isRefreshing) return;
    isRefreshing = true;
    try {
      const data = await RefreshHomeData();
      topPlayed = data?.topPlayed || [];
      recommended = data?.recommended || [];
    } catch (e) {
      console.warn("Failed to refresh For You data:", e);
    } finally {
      isRefreshing = false;
    }
  }

  function handlePlayMix(shuffle = false) {
    const pool =
      mixTracks.length > 0
        ? mixTracks
        : topPlayed.length > 0
          ? topPlayed
          : recommended;
    if (pool.length === 0) return;
    queueStore.setQueue(pool, 0);
    if (shuffle) {
      queueStore.enableShuffle();
    }
    const current = queueStore.getCurrent();
    if (current) {
      playerStore.playTrack(current);
    }
  }

  function handlePlaySection(tracks: SearchResultItem[], shuffle = false) {
    if (!tracks || tracks.length === 0) return;
    queueStore.setQueue(tracks, 0);
    if (shuffle) {
      queueStore.enableShuffle();
    }
    const current = queueStore.getCurrent();
    if (current) {
      playerStore.playTrack(current);
    }
  }
</script>

<div class="px-4 sm:px-8 py-6 space-y-9 select-none">
  <!-- Cinematic Full-Bleed Music Hero: Studio Quality -->
  <HeroBanner
    tag={`PERSONAL MIX • ${getGreeting()}`}
    title={featuredTrack ? featuredTrack.title : $t("nav_foryou")}
    subtitle={featuredTrack ? `${featuredTrack.uploaderName} • ${mixTracks.length > 0 ? `${mixTracks.length} Trek Pilihan Khusus Untukmu` : $t("foryou_hero_desc")}` : $t("foryou_hero_desc")}
    imageUrl={featuredTrack?.thumbnailUrl || ""}
    trackId={featuredTrack?.id}
    playLabel={isMixPlaying ? $t("pause_all") : $t("foryou_play_now")}
    isPlaying={isMixPlaying}
    showPlay={mixTracks.length > 0}
    showShuffle={mixTracks.length > 0}
    on:play={() => handlePlayMix(false)}
    on:shuffle={() => handlePlayMix(true)}
  >
    <button
      type="button"
      class="w-11 h-11 sm:w-12 sm:h-12 rounded-full bg-white/[0.08] hover:bg-white/[0.14] border border-white/10 hover:border-white/20 text-zinc-200 hover:text-white flex items-center justify-center transition-all hover:scale-105 active:scale-95 cursor-pointer disabled:opacity-50 shrink-0"
      on:click={handleRefresh}
      disabled={isRefreshing}
      title={$t("foryou_refresh_mix")}
    >
      <svg
        class="w-4 h-4 {isRefreshing ? 'animate-spin text-crimson-400' : 'text-zinc-200'} shrink-0"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67" />
      </svg>
    </button>
  </HeroBanner>

  <!-- Horizontal Carousel: Koleksi Populer (Inspired by Flume Popular Albums Carousel) -->
  {#if sampleArtworks.length > 0}
    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="w-2 h-2 rounded-full bg-crimson-500"></span>
          <h3
            class="text-sm sm:text-base font-extrabold text-white tracking-tight"
          >
            Koleksi Populer & Pilihan Studio
          </h3>
        </div>
        <span class="text-[11px] text-zinc-500 font-mono"
          >Geser untuk melihat lainnya →</span
        >
      </div>

      <div
        class="flex items-center gap-3.5 overflow-x-auto pb-2 scrollbar-thin"
      >
        {#each sampleArtworks as item, i}
          <div
            class="group w-36 sm:w-40 shrink-0 p-2.5 rounded-2xl bg-white/[0.03] hover:bg-white/[0.08] border border-white/[0.06] hover:border-crimson-500/30 transition-all cursor-pointer space-y-2 active:scale-95"
            on:click={() => {
              queueStore.setQueue(sampleArtworks, i);
              playerStore.playTrack(item);
            }}
            role="button"
            tabindex="0"
            on:keydown={(e) => e.key === "Enter"}
          >
            <div
              class="relative w-full aspect-square rounded-xl overflow-hidden bg-zinc-900 shadow-md"
            >
              <img
                src={item.thumbnailUrl}
                alt={item.title}
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                loading="lazy"
              />
              <div
                class="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
              >
                <div
                  class="w-9 h-9 rounded-full bg-crimson-500 text-white flex items-center justify-center shadow-lg shadow-crimson-500/40 transform group-hover:scale-105 transition-transform"
                >
                  <svg class="w-4 h-4 fill-current ml-0.5" viewBox="0 0 24 24"
                    ><polygon points="5 3 19 12 5 21 5 3" /></svg
                  >
                </div>
              </div>
            </div>

            <div class="space-y-0.5">
              <p
                class="text-xs font-bold text-white truncate group-hover:text-crimson-400 transition-colors"
              >
                {item.title}
              </p>
              <p class="text-[10px] text-zinc-400 truncate">
                {item.uploaderName || "Artis Musik"}
              </p>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  {#if isLoading}
    <div
      class="py-16 text-center text-zinc-500 flex flex-col items-center justify-center gap-3"
    >
      <div
        class="w-6 h-6 border-2 border-white/15 border-t-crimson-500 rounded-full animate-spin"
      ></div>
      <p class="text-xs font-medium">Meracik musik spesial untukmu...</p>
    </div>
  {:else}
    <!-- Section 1: Sering Kamu Putar -->
    {#if topPlayed.length > 0}
      <div>
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2.5">
            <div
              class="w-2.5 h-2.5 rounded-full bg-crimson-500 shadow-[0_0_8px_#E61E38]"
            ></div>
            <h2 class="text-lg font-bold text-white tracking-tight">
              Sering Kamu Putar
            </h2>
          </div>
          <div class="flex items-center gap-2 shrink-0">
            <button
              class="h-9 px-4 rounded-full bg-crimson-600 hover:bg-crimson-500 border border-transparent text-white font-semibold text-xs transition cursor-pointer flex items-center gap-1.5 shadow-sm shadow-crimson-600/30 whitespace-nowrap shrink-0"
              on:click={() => handlePlaySection(topPlayed, false)}
            >
              <svg class="w-3.5 h-3.5 fill-current shrink-0" viewBox="0 0 24 24"
                ><polygon points="5 3 19 12 5 21 5 3" /></svg
              >
              <span class="whitespace-nowrap">Putar Semua</span>
            </button>
            <button
              class="h-9 px-4 rounded-full bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-zinc-300 hover:text-white font-semibold text-xs transition cursor-pointer flex items-center gap-1.5 whitespace-nowrap shrink-0"
              on:click={() => handlePlaySection(topPlayed, true)}
            >
              <svg
                class="w-3.5 h-3.5 shrink-0"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M16 3h5v5M4 20L21 3M21 16v5h-5M15 15l6 6M4 4l5 5" />
              </svg>
              <span class="whitespace-nowrap">Putar Acak</span>
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
          {#each topPlayed as track (track.id)}
            {@const isCurrent = $playerStore.currentTrack?.id === track.id}
            {@const isPlaying = isCurrent && $playerStore.isPlaying}
            <div
              class="group relative flex items-center gap-3 p-2.5 rounded-2xl bg-white/[0.03] hover:bg-white/[0.08] border border-white/[0.06] hover:border-crimson-500/30 transition-all cursor-pointer shadow-sm overflow-hidden {isCurrent
                ? 'bg-crimson-500/10 border-crimson-500/30'
                : ''}"
              on:click={() => {
                queueStore.setQueue(
                  topPlayed,
                  topPlayed.findIndex((t) => t.id === track.id),
                );
                playerStore.playTrack(track);
              }}
              role="button"
              tabindex="0"
              on:keydown={(e) => e.key === "Enter"}
            >
              <div
                class="relative w-12 h-12 shrink-0 rounded-xl overflow-hidden bg-black/40 shadow-xs"
              >
                {#if track.thumbnailUrl}
                  <img
                    src={track.thumbnailUrl}
                    alt={track.title}
                    class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    loading="lazy"
                  />
                {:else}
                  <div
                    class="w-full h-full flex items-center justify-center bg-zinc-800 text-zinc-500"
                  >
                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"
                      ><path
                        d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"
                      /></svg
                    >
                  </div>
                {/if}

                <div
                  class="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  <div
                    class="w-6 h-6 rounded-full bg-crimson-500 text-white flex items-center justify-center shadow-md"
                  >
                    {#if isPlaying}
                      <svg class="w-3 h-3 fill-current" viewBox="0 0 24 24"
                        ><rect x="6" y="4" width="4" height="16" /><rect
                          x="14"
                          y="4"
                          width="4"
                          height="16"
                        /></svg
                      >
                    {:else}
                      <svg
                        class="w-3 h-3 fill-current ml-0.5"
                        viewBox="0 0 24 24"
                        ><polygon points="5 3 19 12 5 21 5 3" /></svg
                      >
                    {/if}
                  </div>
                </div>
              </div>

              <div class="min-w-0 flex-1">
                <p
                  class="text-xs font-bold truncate {isCurrent
                    ? 'text-crimson-400'
                    : 'text-zinc-200 group-hover:text-white'}"
                >
                  {track.title}
                </p>
                <p class="text-[11px] text-zinc-400 truncate mt-0.5">
                  {track.uploaderName}
                </p>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Section 2: Rekomendasi Lagu Populer & Musisi Favorit -->
    <div>
      <TrackList
        title="Rekomendasi Spesial"
        tracks={recommended}
        {isLoading}
        flush={true}
        emptyMessage="Mulai putar beberapa lagu favoritmu untuk melatih rekomendasi cerdas di tab Untuk Anda ini!"
      />
    </div>
  {/if}
</div>
