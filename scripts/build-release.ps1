param (
    [string]$Version = "1.1.2",
    [switch]$SkipFrontendBuild = $false
)

$ErrorActionPreference = "Stop"
$rootDir = (Get-Item $PSScriptRoot).Parent.FullName
Set-Location $rootDir

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host " Dengerin Desktop - Release Packaging" -ForegroundColor Cyan
Write-Host " Version: v$Version" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# 1. Ensure NSIS is in PATH if present in .tools
$localNsis = Join-Path $rootDir ".tools\nsis\tools"
if (Test-Path $localNsis) {
    Write-Host "[OK] NSIS ditemukan di .tools, menambah ke PATH..." -ForegroundColor Green
    $env:PATH = $localNsis + ";" + $env:PATH
}

# 2. Prepare Output Directory
$distDir = Join-Path $rootDir "dist"
if (-not (Test-Path $distDir)) {
    New-Item -ItemType Directory -Path $distDir | Out-Null
}

# 3. Build Frontend
if (-not $SkipFrontendBuild) {
    Write-Host "`n[1/4] Membangun Frontend Svelte..." -ForegroundColor Yellow
    npm run build --prefix frontend
    if ($LASTEXITCODE -ne 0) {
        throw "Gagal membangun frontend!"
    }
}

# 4. Build Executable and NSIS Installer via Wails
Write-Host "`n[2/4] Mengompilasi Biner dan Installer Wails..." -ForegroundColor Yellow
wails build -platform windows/amd64 -clean -nsis -o dengerin.exe
if ($LASTEXITCODE -ne 0) {
    throw "Gagal mengompilasi biner Wails!"
}

$binDir = Join-Path $rootDir "build\bin"
$rawExe = Join-Path $binDir "dengerin.exe"
if (-not (Test-Path $rawExe)) {
    $rawExe = Join-Path $binDir "dengar.exe"
}

if (-not (Test-Path $rawExe)) {
    throw "File biner hasil kompilasi tidak ditemukan di build/bin!"
}

# 5. Package Portable ZIP
Write-Host "`n[3/4] Mengemas Versi Portabel (.zip)..." -ForegroundColor Yellow
$portableZip = Join-Path $distDir "dengerin-v$Version-windows-amd64-portable.zip"
if (Test-Path $portableZip) {
    Remove-Item $portableZip -Force
}

$stageDir = Join-Path $distDir "stage_portable"
if (Test-Path $stageDir) {
    Remove-Item $stageDir -Recurse -Force
}
New-Item -ItemType Directory -Path $stageDir | Out-Null
Copy-Item $rawExe (Join-Path $stageDir "dengerin.exe")

Compress-Archive -Path "$stageDir\*" -DestinationPath $portableZip -Force
Remove-Item $stageDir -Recurse -Force
Write-Host "[OK] Paket Portabel siap: $portableZip" -ForegroundColor Green

# 6. Package Installer EXE
Write-Host "`n[4/4] Mengemas Versi Installer (.exe)..." -ForegroundColor Yellow
$nsisCandidates = Get-ChildItem -Path $binDir -Filter "*installer.exe"
$installerExe = Join-Path $distDir "dengerin-v$Version-windows-amd64-installer.exe"

if ($nsisCandidates.Count -gt 0) {
    Copy-Item $nsisCandidates[0].FullName $installerExe -Force
    Write-Host "[OK] Paket Installer siap: $installerExe" -ForegroundColor Green
} else {
    Write-Host "[!] Installer NSIS tidak terdeteksi, hanya paket portabel yang dibuat." -ForegroundColor DarkYellow
}

# 7. Summary
Write-Host "`n==========================================" -ForegroundColor Cyan
Write-Host " Ringkasan Rilis v$Version Selesai!" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Get-ChildItem -Path $distDir -Filter "dengerin-v$Version-*" | ForEach-Object {
    $hash = (Get-FileHash $_.FullName -Algorithm SHA256).Hash
    $sizeMB = [math]::Round($_.Length / 1MB, 2)
    Write-Host ("- {0} ({1} MB)" -f $_.Name, $sizeMB) -ForegroundColor White
    Write-Host ("  SHA256: {0}" -f $hash) -ForegroundColor Gray
}
