<script lang="ts">
  import { WindowHide, Quit } from "../../../wailsjs/runtime/runtime.js";
  import { closeActionStore, isCloseConfirmOpen, showToast, type CloseAction } from "../stores/settingsStore";
  import { t } from "../stores/i18nStore";

  let rememberChoice = false;

  function handleMinimizeToTray() {
    if (rememberChoice) {
      closeActionStore.set("tray");
    }
    isCloseConfirmOpen.set(false);
    try {
      WindowHide();
      showToast($t("toast_playing_in_background"));
    } catch (e) {
      console.warn("Wails WindowHide error:", e);
    }
  }

  function handleQuit() {
    if (rememberChoice) {
      closeActionStore.set("quit");
    }
    isCloseConfirmOpen.set(false);
    try {
      Quit();
    } catch (e) {
      console.warn("Wails Quit error:", e);
    }
  }

  function handleCancel() {
    isCloseConfirmOpen.set(false);
  }
</script>

<svelte:window
  on:keydown={(e) => $isCloseConfirmOpen && e.key === "Escape" && handleCancel()}
/>

{#if $isCloseConfirmOpen}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-sm animate-fade-in"
    on:click|self={handleCancel}
    on:keydown={(e) => e.key === "Escape" && handleCancel()}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
  >
    <div
      class="w-full max-w-md bg-[#12141A] border border-white/10 rounded-2xl shadow-2xl p-6 text-white relative animate-scale-up"
    >
      <div class="flex items-center gap-3 mb-4">
        <div class="w-10 h-10 rounded-xl bg-crimson-600/20 text-crimson-400 flex items-center justify-center shrink-0">
          <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
            <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
          </svg>
        </div>
        <div>
          <h3 class="text-lg font-bold text-white tracking-tight">{$t("close_dialog_title")}</h3>
          <p class="text-xs text-white/50">{$t("close_dialog_subtitle")}</p>
        </div>
      </div>

      <p class="text-sm text-white/70 mb-5 leading-relaxed">
        {$t("close_dialog_prompt")}
      </p>

      <label class="flex items-center gap-2 mb-6 cursor-pointer select-none">
        <input
          type="checkbox"
          id="rememberCloseChoice"
          bind:checked={rememberChoice}
          class="w-4 h-4 rounded border-white/20 bg-white/5 text-crimson-600 focus:ring-0 focus:ring-offset-0 cursor-pointer"
        />
        <span class="text-xs text-white/60 cursor-pointer">
          {$t("close_dialog_remember")}
        </span>
      </label>

      <div class="flex flex-col sm:flex-row gap-2.5">
        <button
          type="button"
          on:click={handleMinimizeToTray}
          class="flex-1 py-2.5 px-4 rounded-xl bg-crimson-600 hover:bg-crimson-500 active:scale-[0.98] text-white font-medium text-xs sm:text-sm transition-all shadow-lg shadow-crimson-600/20 flex items-center justify-center gap-2 cursor-pointer"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M19 14l-7 7m0 0l-7-7m7 7V3"/>
          </svg>
          <span>{$t("btn_play_in_background")}</span>
        </button>

        <button
          type="button"
          on:click={handleQuit}
          class="py-2.5 px-4 rounded-xl bg-white/5 hover:bg-white/10 active:scale-[0.98] text-white/80 hover:text-white font-medium text-xs sm:text-sm transition-all border border-white/5 cursor-pointer"
        >
          {$t("btn_quit_completely")}
        </button>
      </div>

      <button
        type="button"
        on:click={handleCancel}
        class="absolute top-4 right-4 w-7 h-7 rounded-lg bg-white/5 hover:bg-white/10 text-white/60 hover:text-white flex items-center justify-center transition-colors cursor-pointer"
        aria-label={$t("close")}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>
  </div>
{/if}
