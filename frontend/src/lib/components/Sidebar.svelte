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
      <svg width="20" height="20" viewBox="0 0 269.22 267.55" fill="none">
        <path fill="white" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.88l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
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
