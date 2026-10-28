<script lang="ts">
  import { fade, scale } from 'svelte/transition'
  import { cubicOut } from 'svelte/easing'
  import { isCreatePlaylistOpen, editingPlaylist, activePlaylistId } from '../stores/playlistStore'
  import { userDataStore } from '../stores/userDataStore'
  import { showToast } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'

  let name = ''
  let description = ''
  let isSubmitting = false

  $: if ($isCreatePlaylistOpen) {
    if ($editingPlaylist) {
      name = $editingPlaylist.name
      description = $editingPlaylist.description || ''
    } else {
      name = ''
      description = ''
    }
  }

  function handleClose() {
    isCreatePlaylistOpen.set(false)
    editingPlaylist.set(null)
  }

  async function handleSave() {
    const trimmed = name.trim()
    if (!trimmed) return

    isSubmitting = true
    try {
      if ($editingPlaylist) {
        await userDataStore.updatePlaylist($editingPlaylist.id, trimmed, description.trim())
        showToast('Playlist berhasil diperbarui')
      } else {
        const pl = await userDataStore.createPlaylist(trimmed, description.trim())
        if (pl) {
          activePlaylistId.set(pl.id)
        }
        showToast('Playlist berhasil dibuat')
      }
      handleClose()
    } catch (e: any) {
      showToast(e?.message || 'Gagal menyimpan playlist')
    } finally {
      isSubmitting = false
    }
  }
</script>

{#if $isCreatePlaylistOpen}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4 select-none"
    transition:fade={{ duration: 150 }}
    on:click={handleClose}
    role="presentation"
  >
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class="w-full max-w-md bg-[#11131E] border border-white/[0.12] rounded-3xl shadow-2xl p-6 flex flex-col gap-5 overflow-hidden"
      transition:scale={{ start: 0.95, duration: 180, easing: cubicOut }}
      on:click|stopPropagation
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <div class="flex items-center justify-between border-b border-white/[0.08] pb-3.5">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-xl bg-crimson-500/20 text-crimson-400 flex items-center justify-center">
            <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 14H9v-2h2v2zm0-4H9V7h2v5zm4 4h-2v-2h2v2zm0-4h-2V7h2v5z" opacity="0.3"/>
              <path d="M15 6H3v2h12V6zm0 4H3v2h12v-2zM3 16h8v-2H3v2zM17 6v8.18c-.31-.11-.65-.18-1-.18-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3V8h3V6h-5z"/>
            </svg>
          </div>
          <h3 class="text-base font-bold text-white">
            {$editingPlaylist ? $t('edit_playlist') : $t('create_playlist')}
          </h3>
        </div>
        <button
          class="text-zinc-400 hover:text-white p-1 rounded-lg hover:bg-white/[0.06] transition cursor-pointer"
          on:click={handleClose}
          aria-label="Tutup"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <form on:submit|preventDefault={handleSave} class="flex flex-col gap-4">
        <div>
          <label class="block text-xs font-semibold text-zinc-300 mb-1.5" for="playlist-name">
            {$t('playlist_name')} <span class="text-crimson-400">*</span>
          </label>
          <input
            id="playlist-name"
            type="text"
            bind:value={name}
            placeholder="Misal: My Daily Beats"
            maxlength="60"
            required
            class="w-full bg-[#181A24] border border-white/[0.1] rounded-xl px-3.5 py-2.5 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-crimson-500/60 focus:ring-1 focus:ring-crimson-500/60 transition"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-zinc-300 mb-1.5" for="playlist-desc">
            {$t('playlist_desc')}
          </label>
          <textarea
            id="playlist-desc"
            bind:value={description}
            placeholder="Deskripsi singkat tentang playlist ini..."
            rows="3"
            maxlength="200"
            class="w-full bg-[#181A24] border border-white/[0.1] rounded-xl px-3.5 py-2.5 text-sm text-white placeholder-zinc-500 focus:outline-none focus:border-crimson-500/60 focus:ring-1 focus:ring-crimson-500/60 transition resize-none"
          ></textarea>
        </div>

        <div class="flex items-center justify-end gap-2.5 pt-2">
          <button
            type="button"
            class="px-4 py-2 text-xs font-semibold text-zinc-400 hover:text-white rounded-xl hover:bg-white/[0.06] transition cursor-pointer"
            on:click={handleClose}
          >
            {$t('close')}
          </button>
          <button
            type="submit"
            disabled={!name.trim() || isSubmitting}
            class="px-5 py-2 text-xs font-semibold bg-crimson-500 hover:bg-crimson-600 disabled:opacity-50 text-white rounded-xl shadow-lg shadow-crimson-500/25 transition cursor-pointer flex items-center gap-2"
          >
            {#if isSubmitting}
              <span class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            {/if}
            <span>Simpan</span>
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
