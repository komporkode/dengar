<script lang="ts">
  import type { NavigationTab } from '../types'
  import { userDataStore } from '../stores/userDataStore'

  export let activeTab: NavigationTab = 'search'

  $: favCount = $userDataStore.favorites.length
  $: histCount = $userDataStore.history.length
</script>

<aside class="sidebar">
  <div class="logo-area">
    <div class="logo-icon">
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
        <path d="M9 18V5l12-2v13"></path>
        <circle cx="6" cy="18" r="3" fill="currentColor"></circle>
        <circle cx="18" cy="16" r="3" fill="currentColor"></circle>
      </svg>
    </div>
    <div class="logo-text">
      <span class="brand-title">Dengar</span>
      <span class="badge">PROD</span>
    </div>
  </div>

  <nav class="nav-links">
    <button
      class="nav-item"
      class:active={activeTab === 'search'}
      on:click={() => (activeTab = 'search')}
    >
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8"></circle>
        <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
      </svg>
      <span>Eksplorasi & Cari</span>
    </button>

    <button
      class="nav-item"
      class:active={activeTab === 'favorites'}
      on:click={() => (activeTab = 'favorites')}
    >
      <svg width="20" height="20" viewBox="0 0 24 24" fill={activeTab === 'favorites' ? '#ff0000' : 'none'} stroke="currentColor" stroke-width="2">
        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
      </svg>
      <span>Lagu Favorit</span>
      {#if favCount > 0}
        <span class="counter">{favCount}</span>
      {/if}
    </button>

    <button
      class="nav-item"
      class:active={activeTab === 'history'}
      on:click={() => (activeTab = 'history')}
    >
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <polyline points="12 6 12 12 16 14"></polyline>
      </svg>
      <span>Riwayat Putar</span>
      {#if histCount > 0}
        <span class="counter">{histCount}</span>
      {/if}
    </button>
  </nav>

  <div class="sidebar-footer">
    <div class="network-badge">
      <span class="pulse-dot"></span>
      <span>Instance Publik OK</span>
    </div>
  </div>
</aside>

<style>
  .sidebar {
    width: var(--sidebar-width);
    height: calc(100vh - var(--player-height));
    background: var(--bg-secondary);
    border-right: 1px solid var(--border-subtle);
    display: flex;
    flex-direction: column;
    padding: 20px 14px;
    user-select: none;
    flex-shrink: 0;
  }

  .logo-area {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 6px 10px 24px 10px;
  }

  .logo-icon {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: var(--accent-red);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    box-shadow: 0 0 14px var(--accent-glow);
  }

  .logo-text {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .brand-title {
    font-size: 20px;
    font-weight: 700;
    letter-spacing: -0.5px;
    color: var(--text-primary);
  }

  .badge {
    font-size: 9px;
    font-weight: 700;
    padding: 2px 5px;
    border-radius: 4px;
    background: rgba(255, 255, 255, 0.1);
    color: var(--text-secondary);
    letter-spacing: 0.5px;
  }

  .nav-links {
    display: flex;
    flex-direction: column;
    gap: 4px;
    flex: 1;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 12px 14px;
    border-radius: 10px;
    color: var(--text-secondary);
    font-weight: 500;
    transition: all 0.15s ease;
    text-align: left;
    width: 100%;
  }

  .nav-item:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .nav-item.active {
    background: var(--bg-card);
    color: var(--text-primary);
    font-weight: 600;
    border-left: 3px solid var(--accent-red);
  }

  .nav-item.active svg {
    color: var(--accent-red);
  }

  .counter {
    margin-left: auto;
    font-size: 11px;
    padding: 2px 7px;
    background: rgba(255, 255, 255, 0.08);
    border-radius: 10px;
    color: var(--text-muted);
  }

  .sidebar-footer {
    padding-top: 14px;
    border-top: 1px solid var(--border-subtle);
  }

  .network-badge {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 11px;
    color: #4ade80;
    padding: 6px 10px;
    border-radius: 8px;
    background: rgba(74, 222, 128, 0.08);
  }

  .pulse-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #4ade80;
    box-shadow: 0 0 8px #4ade80;
  }
</style>
