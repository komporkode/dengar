<script lang="ts">
  import { fade, scale } from 'svelte/transition'
  import { cubicOut } from 'svelte/easing'
  import { sharingPlaylist, isImportPlaylistOpen, activePlaylistId } from '../stores/playlistStore'
  import { userDataStore } from '../stores/userDataStore'
  import { showToast } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'

  let shareCode = ''
  let importPayload = ''
  let isCopied = false
  let isProcessing = false

  $: if ($sharingPlaylist) {
    loadShareCode($sharingPlaylist.id)
  }

  async function loadShareCode(id: string) {
    try {
      shareCode = await userDataStore.exportPlaylist(id)
    } catch (e: any) {
      console.error('Failed to export playlist:', e)
      showToast('Gagal membuat kode berbagi')
    }
  }

  function handleClose() {
    sharingPlaylist.set(null)
    isImportPlaylistOpen.set(false)
    shareCode = ''
    importPayload = ''
    isCopied = false
  }

  async function handleCopy() {
    if (!shareCode) return
    try {
      await navigator.clipboard.writeText(shareCode)
      isCopied = true
      showToast($t('copied'))
      setTimeout(() => (isCopied = false), 2000)
    } catch (e) {
      showToast('Gagal menyalin ke clipboard')
    }
  }

  function handleDownloadJSON() {
    if (!$sharingPlaylist) return
    const exportData = {
      name: $sharingPlaylist.name,
      description: $sharingPlaylist.description,
      tracks: $sharingPlaylist.tracks || [],
    }
    const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${$sharingPlaylist.name.replace(/[^\w\s-]/gi, '') || 'playlist'}.json`
    a.click()
    URL.revokeObjectURL(url)
    showToast('File JSON berhasil diunduh')
  }

  function handleFileSelect(e: Event) {
    const input = e.target as HTMLInputElement
    if (!input.files || input.files.length === 0) return
    const file = input.files[0]
    const reader = new FileReader()
    reader.onload = (event) => {
      importPayload = (event.target?.result as string) || ''
    }
    reader.readAsText(file)
  }

  async function handleImport() {
    const trimmed = importPayload.trim()
    if (!trimmed) return

    isProcessing = true
    try {
      const pl = await userDataStore.importPlaylist(trimmed)
      if (pl) {
        activePlaylistId.set(pl.id)
        showToast(`Berhasil mengimpor playlist "${pl.name}"!`)
      }
      handleClose()
    } catch (e: any) {
      showToast(e?.message || 'Format kode atau file tidak valid')
    } finally {
      isProcessing = false
    }
  }
</script>

{#if $sharingPlaylist || $isImportPlaylistOpen}
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
            <svg class="w-4 h-4 fill-none stroke-current stroke-2" viewBox="0 0 24 24">
              {#if $sharingPlaylist}
                <circle cx="18" cy="5" r="3"></circle>
                <circle cx="6" cy="12" r="3"></circle>
                <circle cx="18" cy="19" r="3"></circle>
                <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line>
                <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line>
              {:else}
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="7 10 12 15 17 10"></polyline>
                <line x1="12" y1="15" x2="12" y2="3"></line>
              {/if}
            </svg>
          </div>
          <h3 class="text-base font-bold text-white">
            {$sharingPlaylist ? $t('share_playlist') : $t('import_playlist')}
          </h3>
        </div>
        <button
          class="text-zinc-400 hover:text-white p-1 rounded-lg hover:bg-white/[0.06] transition cursor-pointer"
          on:click={handleClose}
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      {#if $sharingPlaylist}
        <!-- Export / Share Mode -->
        <div class="flex flex-col gap-4">
          <div>
            <label class="block text-xs font-semibold text-zinc-300 mb-1.5" for="share-code-input">
              {$t('share_code')}
            </label>
            <div class="relative">
              <input
                id="share-code-input"
                type="text"
                readonly
                value={shareCode}
                class="w-full bg-[#181A24] border border-white/[0.1] rounded-xl px-3.5 py-2.5 text-xs font-mono text-zinc-300 pr-24 focus:outline-none"
              />
              <button
                class="absolute right-1.5 top-1.5 bottom-1.5 px-3 bg-crimson-500 hover:bg-crimson-600 text-white rounded-lg text-xs font-semibold transition cursor-pointer flex items-center gap-1.5"
                on:click={handleCopy}
              >
                {#if isCopied}
                  <svg class="w-3.5 h-3.5 fill-none stroke-current stroke-2" viewBox="0 0 24 24">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                  <span>{$t('copied')}</span>
                {:else}
                  <svg class="w-3.5 h-3.5 fill-none stroke-current stroke-2" viewBox="0 0 24 24">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                  </svg>
                  <span>{$t('copy_code')}</span>
                {/if}
              </button>
            </div>
            <p class="text-[11px] text-zinc-400 mt-1.5">
              Bagikan kode ini kepada teman Anda untuk diimpor langsung ke aplikasi Dengar mereka.
            </p>
          </div>

          <div class="border-t border-white/[0.08] pt-3 flex items-center justify-between">
            <span class="text-xs text-zinc-400">Atau simpan sebagai berkas cadangan:</span>
            <button
              class="px-3.5 py-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.1] text-xs font-semibold text-zinc-200 transition cursor-pointer flex items-center gap-2"
              on:click={handleDownloadJSON}
            >
              <svg class="w-3.5 h-3.5 fill-none stroke-current stroke-2" viewBox="0 0 24 24">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="7 10 12 15 17 10"></polyline>
                <line x1="12" y1="15" x2="12" y2="3"></line>
              </svg>
              <span>{$t('export_file')}</span>
            </button>
          </div>
        </div>
      {:else}
        <!-- Import Mode -->
        <div class="flex flex-col gap-4">
          <div>
            <label class="block text-xs font-semibold text-zinc-300 mb-1.5" for="import-payload">
              {$t('import_code_or_file')}
            </label>
            <textarea
              id="import-payload"
              bind:value={importPayload}
              placeholder="Tempel kode Base64 atau isi JSON playlist di sini..."
              rows="4"
              class="w-full bg-[#181A24] border border-white/[0.1] rounded-xl px-3.5 py-2.5 text-xs font-mono text-white placeholder-zinc-500 focus:outline-none focus:border-crimson-500/60 focus:ring-1 focus:ring-crimson-500/60 transition resize-none"
            ></textarea>
          </div>

          <div class="flex items-center justify-between border-t border-white/[0.08] pt-3">
            <label class="px-3.5 py-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.1] text-xs font-semibold text-zinc-300 hover:text-white transition cursor-pointer flex items-center gap-2">
              <svg class="w-3.5 h-3.5 fill-none stroke-current stroke-2" viewBox="0 0 24 24">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="17 8 12 3 7 8"></polyline>
                <line x1="12" y1="3" x2="12" y2="15"></line>
              </svg>
              <span>Pilih File .JSON</span>
              <input type="file" accept=".json" on:change={handleFileSelect} class="hidden" />
            </label>

            <button
              disabled={!importPayload.trim() || isProcessing}
              class="px-5 py-2 text-xs font-semibold bg-crimson-500 hover:bg-crimson-600 disabled:opacity-50 text-white rounded-xl shadow-lg shadow-crimson-500/25 transition cursor-pointer flex items-center gap-2"
              on:click={handleImport}
            >
              {#if isProcessing}
                <span class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
              {/if}
              <span>{$t('import_btn')}</span>
            </button>
          </div>
        </div>
      {/if}
    </div>
  </div>
{/if}
