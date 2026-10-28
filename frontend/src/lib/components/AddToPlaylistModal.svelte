<script lang="ts">
  import { fade, scale } from 'svelte/transition'
  import { cubicOut } from 'svelte/easing'
  import { trackToAddToPlaylist, openCreatePlaylistModal } from '../stores/playlistStore'
  import { userDataStore } from '../stores/userDataStore'
  import { showToast } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'

  function handleClose() {
    trackToAddToPlaylist.set(null)
  }

  async function handleToggleTrack(playlistId: string, isInPlaylist: boolean) {
    if (!$trackToAddToPlaylist) return
    const track = $trackToAddToPlaylist

    if (isInPlaylist) {
      await userDataStore.removeTrackFromPlaylist(playlistId, track.id)
      showToast('Lagu dihapus dari playlist')
    } else {
      await userDataStore.addTrackToPlaylist(playlistId, track)
      showToast('Lagu ditambahkan ke playlist')
    }
  }

  function handleCreateNew() {
    handleClose()
    openCreatePlaylistModal()
  }
</script>

{#if $trackToAddToPlaylist}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4 select-none"
    transition:fade={{ duration: 150 }}
    on:click={handleClose}
    role="presentation"
  >
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class="w-full max-w-sm bg-[#11131E] border border-white/[0.12] rounded-3xl shadow-2xl p-5 flex flex-col gap-4 overflow-hidden"
      transition:scale={{ start: 0.95, duration: 180, easing: cubicOut }}
      on:click|stopPropagation
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <div class="flex items-center justify-between border-b border-white/[0.08] pb-3">
        <div class="min-w-0 pr-2">
          <h3 class="text-sm font-bold text-white truncate">{$t('add_to_playlist')}</h3>
          <p class="text-xs text-zinc-400 truncate mt-0.5">{$trackToAddToPlaylist.title}</p>
        </div>
        <button
          class="text-zinc-400 hover:text-white p-1 rounded-lg hover:bg-white/[0.06] transition cursor-pointer shrink-0"
          on:click={handleClose}
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <!-- Playlists list -->
      <div class="flex flex-col gap-1.5 max-h-60 overflow-y-auto pr-1">
        {#if $userDataStore.playlists.length === 0}
          <div class="py-6 text-center text-xs text-zinc-500">
            Belum ada playlist. Buat playlist baru untuk memulai.
          </div>
        {:else}
          {#each $userDataStore.playlists as pl (pl.id)}
            {@const isIn = (pl.tracks || []).some(t => t.id === $trackToAddToPlaylist?.id)}
            <button
              class="w-full flex items-center justify-between p-2.5 rounded-xl border transition text-left cursor-pointer {isIn ? 'bg-crimson-500/10 border-crimson-500/30 text-white' : 'hover:bg-white/[0.05] border-white/[0.06] text-zinc-300'}"
              on:click={() => handleToggleTrack(pl.id, isIn)}
            >
              <div class="min-w-0 pr-2">
                <div class="text-xs font-semibold truncate">{pl.name}</div>
                <div class="text-[10px] text-zinc-400 mt-0.5">{(pl.tracks || []).length} {$t('items_count')}</div>
              </div>
              <div class="shrink-0 w-5 h-5 rounded-md flex items-center justify-center border {isIn ? 'bg-crimson-500 border-crimson-500 text-white' : 'border-white/20 text-transparent'}">
                <svg class="w-3.5 h-3.5 fill-none stroke-current stroke-[2.5]" viewBox="0 0 24 24">
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </div>
            </button>
          {/each}
        {/if}
      </div>

      <!-- Create New Playlist button -->
      <button
        class="w-full flex items-center justify-center gap-2 py-2.5 rounded-xl border border-dashed border-white/20 hover:border-crimson-500/50 hover:bg-crimson-500/[0.06] text-xs font-medium text-zinc-300 hover:text-white transition cursor-pointer"
        on:click={handleCreateNew}
      >
        <svg class="w-4 h-4 text-crimson-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        <span>{$t('create_playlist')}</span>
      </button>
    </div>
  </div>
{/if}
