<script lang="ts">
  import { onMount } from 'svelte'
  import type { SearchResultItem } from '../types'
  import { GetHomeData } from '../../../wailsjs/go/main/App'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { openAddToPlaylistModal } from '../stores/playlistStore'
  import { t } from '../stores/i18nStore'

  let topPlayed: SearchResultItem[] = []
  let recommended: SearchResultItem[] = []
  let isLoading = true

  onMount(async () => {
    await loadHomeData()
  })

  async function loadHomeData() {
    isLoading = true
    try {
      const data = await GetHomeData()
      topPlayed = data.topPlayed || []
      recommended = data.recommended || []
    } catch (e) {
      console.warn('Failed to load home data:', e)
    } finally {
      isLoading = false
    }
  }

  function handlePlayTrack(track: SearchResultItem, context: SearchResultItem[]) {
    queueStore.setQueue(context, context.findIndex(t => t.id === track.id))
    playerStore.playTrack(track)
  }
</script>

{#if !isLoading && (topPlayed.length > 0 || recommended.length > 0)}
  <div class="px-4 sm:px-8 pt-4 pb-2 space-y-6 select-none">
    <!-- Section 1: Sering Kamu Putar -->
    {#if topPlayed.length > 0}
      <div>
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 rounded-full bg-crimson-500 shadow-[0_0_8px_#E61E38]"></div>
            <h2 class="text-sm sm:text-base font-bold text-white tracking-tight">
              {$t('home_top_played')}
            </h2>
          </div>
        </div>

        <!-- Horizontal scrollable or flex wrap of compact cards -->
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2.5 sm:gap-3">
          {#each topPlayed.slice(0, 8) as track (track.id)}
            {@const isCurrent = $playerStore.currentTrack?.id === track.id}
            {@const isPlaying = isCurrent && $playerStore.isPlaying}
            <div
              class="group relative flex items-center gap-3 p-2 rounded-xl bg-white/[0.03] hover:bg-white/[0.08] border border-white/[0.06] hover:border-crimson-500/30 transition-all cursor-pointer shadow-sm overflow-hidden {isCurrent ? 'bg-crimson-500/10 border-crimson-500/30' : ''}"
              on:click={() => handlePlayTrack(track, topPlayed)}
              role="button"
              tabindex="0"
              on:keydown={(e) => e.key === 'Enter' && handlePlayTrack(track, topPlayed)}
            >
              <div class="relative w-12 h-12 shrink-0 rounded-lg overflow-hidden bg-black/40 shadow-xs">
                {#if track.thumbnailUrl}
                  <img src={track.thumbnailUrl} alt={track.title} class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" loading="lazy" />
                {:else}
                  <div class="w-full h-full flex items-center justify-center bg-zinc-800 text-zinc-500">
                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
                  </div>
                {/if}

                <!-- Play Overlay Button -->
                <div class="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                  <div class="w-6 h-6 rounded-full bg-crimson-500 text-white flex items-center justify-center shadow-md">
                    {#if isPlaying}
                      <svg class="w-3 h-3 fill-current" viewBox="0 0 24 24"><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></svg>
                    {:else}
                      <svg class="w-3 h-3 fill-current ml-0.5" viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                    {/if}
                  </div>
                </div>
              </div>

              <div class="min-w-0 flex-1 pr-1">
                <p class="text-xs font-semibold truncate {isCurrent ? 'text-crimson-400' : 'text-zinc-200 group-hover:text-white'}">
                  {track.title}
                </p>
                <p class="text-[11px] text-zinc-400 truncate mt-0.5">
                  {track.uploaderName}
                </p>
              </div>

              <!-- Quick Add to Playlist Button -->
              <button
                class="shrink-0 p-1 text-zinc-500 hover:text-white opacity-0 group-hover:opacity-100 transition rounded-lg hover:bg-white/[0.1] cursor-pointer"
                on:click|stopPropagation={() => openAddToPlaylistModal(track)}
                title={$t('add_to_playlist')}
              >
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <circle cx="12" cy="12" r="1.5"></circle>
                  <circle cx="12" cy="5" r="1.5"></circle>
                  <circle cx="12" cy="19" r="1.5"></circle>
                </svg>
              </button>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Section 2: Rekomendasi Untuk Kamu -->
    {#if recommended.length > 0}
      <div>
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 rounded-full bg-rose-400 shadow-[0_0_8px_#F43F5E]"></div>
            <h2 class="text-sm sm:text-base font-bold text-white tracking-tight">
              {$t('home_recommended')}
            </h2>
          </div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
          {#each recommended.slice(0, 6) as track (track.id)}
            {@const isCurrent = $playerStore.currentTrack?.id === track.id}
            {@const isPlaying = isCurrent && $playerStore.isPlaying}
            <div
              class="group flex flex-col p-2.5 rounded-2xl bg-white/[0.02] hover:bg-white/[0.06] border border-white/[0.04] hover:border-crimson-500/20 transition-all cursor-pointer shadow-sm relative {isCurrent ? 'border-crimson-500/40 bg-crimson-500/[0.06]' : ''}"
              on:click={() => handlePlayTrack(track, recommended)}
              role="button"
              tabindex="0"
              on:keydown={(e) => e.key === 'Enter' && handlePlayTrack(track, recommended)}
            >
              <div class="relative w-full aspect-square rounded-xl overflow-hidden bg-black/40 shadow-sm mb-2.5">
                {#if track.thumbnailUrl}
                  <img src={track.thumbnailUrl} alt={track.title} class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" loading="lazy" />
                {:else}
                  <div class="w-full h-full flex items-center justify-center bg-zinc-800 text-zinc-500">
                    <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
                  </div>
                {/if}

                <!-- Hover Play Badge Button -->
                <button
                  class="absolute bottom-2 right-2 w-9 h-9 rounded-full bg-crimson-500 text-white flex items-center justify-center shadow-lg shadow-crimson-500/30 opacity-0 translate-y-2 group-hover:opacity-100 group-hover:translate-y-0 transition-all duration-200 cursor-pointer hover:scale-105"
                  on:click|stopPropagation={() => handlePlayTrack(track, recommended)}
                >
                  {#if isPlaying}
                    <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24"><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></svg>
                  {:else}
                    <svg class="w-4 h-4 fill-current ml-0.5" viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                  {/if}
                </button>
              </div>

              <div class="flex items-center justify-between gap-1">
                <div class="min-w-0 flex-1">
                  <p class="text-xs font-semibold truncate {isCurrent ? 'text-crimson-400' : 'text-zinc-200 group-hover:text-white'}">
                    {track.title}
                  </p>
                  <p class="text-[11px] text-zinc-400 truncate mt-0.5">
                    {track.uploaderName}
                  </p>
                </div>
                <button
                  class="shrink-0 p-1 text-zinc-500 hover:text-white opacity-0 group-hover:opacity-100 transition rounded-lg hover:bg-white/[0.1] cursor-pointer"
                  on:click|stopPropagation={() => openAddToPlaylistModal(track)}
                  title={$t('add_to_playlist')}
                >
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <circle cx="12" cy="12" r="1.5"></circle>
                    <circle cx="12" cy="5" r="1.5"></circle>
                    <circle cx="12" cy="19" r="1.5"></circle>
                  </svg>
                </button>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </div>
{/if}
