# Situs Web Dengerin Desktop

Folder ini berisi situs statis resmi untuk Dengerin Desktop (`https://github.com/komporkode/dengar`), terpisah dari aplikasi desktop Wails pada direktori root.

## Opsi 1: Deploy di Vercel (Direkomendasikan)

### Cara A: Via Dasbor Vercel (Auto-Deploy dari GitHub)

1. Commit dan push folder `website/` ke repositori `komporkode/dengar`.
2. Buka dasbor Vercel, klik **Add New > Project**, lalu impor repositori `komporkode/dengar`.
3. Atur konfigurasi proyek:
   - **Project Name:** isi bebas (misalnya `dengerin` untuk URL `dengerin.vercel.app`)
   - **Root Directory:** klik **Edit** dan pilih `website`
   - **Framework Preset:** `Other`
   - **Build Command:** *(kosongkan)*
   - **Output Directory:** `.`
4. Opsional: pada **Settings > Build and Deployment > Ignored Build Step**, pilih **Custom Command** lalu isi `git diff HEAD^ HEAD --quiet -- .` agar perubahan aplikasi desktop di luar folder `website/` tidak memicu deploy ulang situs web.
5. Klik **Deploy**.

### Cara B: Via Vercel CLI (Langsung dari Lokal)

Untuk deploy cepat dari folder `website` tanpa perlu push ke GitHub terlebih dahulu:

```bash
cd website
npx vercel --prod
```

Pengaturan `"framework": null` dan `"outputDirectory": "."` sudah disiapkan di `website/vercel.json`.

Vercel akan menyajikan `website/index.html`, `website/404.html`, serta header keamanan dan cache dari `website/vercel.json`.

## Opsi 2: Deploy di Cloudflare Pages

1. Buka dasbor **Workers & Pages > Create application**, lalu pastikan masuk ke menu **Pages** (klik tab **Pages** atau tautan **"Looking to deploy Pages? Get started"** di bagian bawah halaman, bukan form **Workers** default agar opsi *Build output directory* muncul).
2. Pilih **Import an existing Git repository** dan hubungkan repositori `komporkode/dengar`.
3. Gunakan nama proyek bebas yang diinginkan (misalnya `dengerin` untuk menghasilkan domain `dengerin.pages.dev`).
4. Konfigurasi Build:
   - **Framework preset:** `None`
   - **Build command:** *(kosongkan)*
   - **Build output directory:** `website`
5. Klik **Save and Deploy**.

Jika ingin upload langsung ke Cloudflare Pages tanpa konfigurasi build dari GitHub, Anda juga bisa memilih **Upload assets** di menu Pages atau menjalankan:

```bash
npx wrangler pages deploy website --project-name=dengerin
```

Cloudflare Pages akan menyajikan `website/index.html`, `website/404.html`, serta header keamanan dan cache dari `website/_headers`.
