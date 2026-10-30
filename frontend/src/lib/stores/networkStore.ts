import { writable, get } from 'svelte/store'

export type NetworkStatus = 'online' | 'offline' | 'weak'

export interface NetworkState {
  status: NetworkStatus
  isOnline: boolean
  isWeak: boolean
  latencyMs: number | null
  effectiveType: string | null
  lastChecked: number
}

function createNetworkStore() {
  const isClient = typeof window !== 'undefined'
  const initialOnline = isClient ? navigator.onLine : true

  const { subscribe, update } = writable<NetworkState>({
    status: initialOnline ? 'online' : 'offline',
    isOnline: initialOnline,
    isWeak: false,
    latencyMs: null,
    effectiveType: null,
    lastChecked: Date.now(),
  })

  let reconnectListeners: Array<() => void> = []
  let offlineListeners: Array<() => void> = []
  let pingInterval: any = null
  let isChecking = false

  async function checkConnectivity(): Promise<boolean> {
    if (!isClient || isChecking) return get({ subscribe }).isOnline
    isChecking = true

    // Fast preliminary check with browser API
    if (typeof navigator !== 'undefined' && !navigator.onLine) {
      handleOffline()
      isChecking = false
      return false
    }

    const start = performance.now()
    try {
      const controller = new AbortController()
      const timeoutId = setTimeout(() => controller.abort(), 2500)

      // Zero-byte captive portal endpoint from Google (universal CORS-free endpoint)
      await fetch('https://www.gstatic.com/generate_204', {
        method: 'HEAD',
        mode: 'no-cors',
        cache: 'no-store',
        signal: controller.signal,
      })
      clearTimeout(timeoutId)

      const latency = Math.round(performance.now() - start)
      const conn = (navigator as any).connection
      const effType = conn?.effectiveType || null

      // Weak if latency > 1800ms or 2G/slow-2G
      const isWeak = latency > 1800 || effType === '2g' || effType === 'slow-2g'
      const status: NetworkStatus = isWeak ? 'weak' : 'online'

      const prev = get({ subscribe })
      const wasOffline = !prev.isOnline

      update(st => ({
        ...st,
        status,
        isOnline: true,
        isWeak,
        latencyMs: latency,
        effectiveType: effType,
        lastChecked: Date.now(),
      }))

      if (wasOffline) {
        console.log(`[NetworkStore] Connection restored (RTT: ${latency}ms)`)
        reconnectListeners.forEach(fn => {
          try { fn() } catch (err) { console.error('Reconnect listener error:', err) }
        })
      }

      isChecking = false
      return true
    } catch {
      handleOffline()
      isChecking = false
      return false
    }
  }

  function handleOffline() {
    const prev = get({ subscribe })
    const wasOnline = prev.isOnline

    update(st => ({
      ...st,
      status: 'offline',
      isOnline: false,
      isWeak: false,
      latencyMs: null,
      lastChecked: Date.now(),
    }))

    if (wasOnline) {
      console.warn('[NetworkStore] Device went offline')
      offlineListeners.forEach(fn => {
        try { fn() } catch (err) { console.error('Offline listener error:', err) }
      })
    }
  }

  function init() {
    if (!isClient) return

    window.addEventListener('online', () => {
      checkConnectivity()
    })

    window.addEventListener('offline', () => {
      handleOffline()
    })

    const conn = (navigator as any).connection
    if (conn) {
      conn.addEventListener('change', () => {
        checkConnectivity()
      })
    }

    // Initial probe
    checkConnectivity()

    // Periodic heartbeat every 20 seconds
    clearInterval(pingInterval)
    pingInterval = setInterval(() => {
      checkConnectivity()
    }, 20000)
  }

  init()

  return {
    subscribe,
    checkNow: checkConnectivity,
    onReconnect: (cb: () => void) => {
      reconnectListeners.push(cb)
      return () => {
        reconnectListeners = reconnectListeners.filter(f => f !== cb)
      }
    },
    onOffline: (cb: () => void) => {
      offlineListeners.push(cb)
      return () => {
        offlineListeners = offlineListeners.filter(f => f !== cb)
      }
    },
  }
}

export const networkStore = createNetworkStore()

