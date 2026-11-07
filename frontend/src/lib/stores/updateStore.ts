import { writable, get } from 'svelte/store'
import { CheckForUpdates, DownloadUpdate, OpenInFolder, ApplyUpdateAndRestart } from '../../../wailsjs/go/main/App'
import { BrowserOpenURL, EventsOn } from '../../../wailsjs/runtime/runtime.js'
import { showToast } from './settingsStore'

export interface UpdateAsset {
  name: string
  downloadUrl: string
  size: number
}

export interface DownloadProgress {
  downloaded: number
  total: number
  percent: number
  speedText: string
  fileName: string
}

export interface UpdateInfo {
  currentVersion: string
  latestVersion: string
  hasUpdate: boolean
  releaseUrl: string
  releaseNotes: string
  publishedAt: string
  downloadUrl: string
  assetFileName: string
  assetSize: number
  assets?: UpdateAsset[]
  isInstalled: boolean
  installerUrl?: string
  installerFileName?: string
  installerSize?: number
  portableUrl?: string
  portableFileName?: string
  portableSize?: number
}

export const updateState = writable<UpdateInfo | null>(null)
export const isCheckingUpdate = writable<boolean>(false)
export const isDownloadingUpdate = writable<boolean>(false)
export const downloadProgress = writable<string>('')
export const downloadProgressData = writable<DownloadProgress | null>(null)
export const downloadedFilePath = writable<string | null>(null)
export const showUpdateBanner = writable<boolean>(false)

if (typeof window !== 'undefined') {
  try {
    EventsOn('update:download-progress', (data: DownloadProgress) => {
      downloadProgressData.set(data)
    })
  } catch (e) {
    console.warn('Wails EventsOn update:download-progress unavailable:', e)
  }
}

export function formatBytes(bytes: number): string {
  if (!bytes || bytes <= 0) return '0 MB'
  const mb = bytes / (1024 * 1024)
  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(1)} GB`
  }
  return `${mb.toFixed(1)} MB`
}

/**
 * Otomatis periksa pembaruan di latar belakang saat aplikasi pertama kali dibuka.
 * Menunggu jeda 2.5 detik agar inisialisasi awal (search & audio) selesai terlebih dahulu.
 */
export async function checkForUpdatesInBackground() {
  await new Promise(r => setTimeout(r, 2500))

  try {
    isCheckingUpdate.set(true)
    const res = await CheckForUpdates()
    if (res) {
      updateState.set(res)
      if (res.hasUpdate) {
        showUpdateBanner.set(true)
      }
    }
  } catch (err) {
    console.warn('Background update check failed (offline or rate limited):', err)
  } finally {
    isCheckingUpdate.set(false)
  }
}

/**
 * Periksa pembaruan secara manual dari menu Pengaturan.
 */
export async function checkForUpdatesManual(): Promise<UpdateInfo | null> {
  isCheckingUpdate.set(true)
  try {
    const res = await CheckForUpdates()
    if (res) {
      updateState.set(res)
      if (res.hasUpdate) {
        showToast(`Versi baru ${res.latestVersion} tersedia!`)
      } else {
        showToast(`Aplikasi sudah versi terbaru (${res.currentVersion})`)
      }
      return res
    }
    return null
  } catch (err: any) {
    console.error('Manual update check failed:', err)
    showToast('Gagal memeriksa pembaruan. Periksa koneksi internet.')
    return null
  } finally {
    isCheckingUpdate.set(false)
  }
}

/**
 * Unduh paket pembaruan biner (.zip) langsung ke folder Downloads pengguna.
 */
export async function startDownloadUpdate(customUrl?: string, customFileName?: string) {
  const current = get(updateState)
  const targetUrl = customUrl || current?.downloadUrl
  const targetName = customFileName || current?.assetFileName || `dengerin-${current?.latestVersion || 'latest'}-windows-amd64.zip`

  if (!targetUrl) {
    if (current?.releaseUrl) {
      BrowserOpenURL(current.releaseUrl)
    } else {
      showToast('URL unduhan pembaruan tidak ditemukan.')
    }
    return
  }

  isDownloadingUpdate.set(true)
  downloadProgressData.set(null)
  downloadProgress.set('Mengunduh berkas pembaruan...')

  try {
    const savedPath = await DownloadUpdate(targetUrl, targetName)
    downloadedFilePath.set(savedPath)
    downloadProgressData.set(null)
    showToast('Pembaruan selesai diunduh! Siap dipasang.')
  } catch (err: any) {
    console.error('Download update failed:', err)
    showToast(`Gagal mengunduh: ${err?.message || err}`)
    downloadProgressData.set(null)
  } finally {
    isDownloadingUpdate.set(false)
    downloadProgress.set('')
  }
}

/**
 * Membuka Windows Explorer dan menyorot file .zip yang telah diunduh.
 */
export async function revealDownloadedFile() {
  const path = get(downloadedFilePath)
  if (!path) {
    showToast('Belum ada file yang diunduh.')
    return
  }

  try {
    await OpenInFolder(path)
  } catch (err: any) {
    console.error('OpenInFolder error:', err)
    showToast(`Gagal membuka folder: ${err?.message || err}`)
  }
}

/**
 * Terapkan pembaruan biner atau jalankan installer secara otomatis dan mulai ulang aplikasi.
 */
export async function installAndRestartUpdate() {
  const path = get(downloadedFilePath)
  if (!path) {
    showToast('Belum ada file yang diunduh.')
    return
  }

  try {
    showToast('Memasang pembaruan dan memulai ulang aplikasi...')
    await ApplyUpdateAndRestart(path)
  } catch (err: any) {
    console.error('ApplyUpdateAndRestart error:', err)
    showToast(`Gagal memasang otomatis: ${err?.message || err}. Buka folder untuk memasang manual.`)
    await OpenInFolder(path)
  }
}

