# Deployment — go-kos-api

## Menjalankan lokal

```bash
go run .
# atau build dulu:
go build -o go-kos-api .
PORT=3000 ./go-kos-api
```

Server mendengarkan di `0.0.0.0:$PORT` (default 8080).

## Halaman dokumentasi (Cloudflare Pages)

Folder `landing/` adalah situs statis (tanpa build step) yang di-deploy ke
Cloudflare Pages sebagai halaman dokumentasi + demo interaktif API:

```bash
# lewat wrangler (butuh login Cloudflare)
npx wrangler pages deploy landing --project-name go-kos-api
```

URL live: https://go-kos-api.pages.dev

Halaman ini statis — demo interaktifnya memakai data contoh yang sama
dengan seed API, bukan memanggil API live (API-nya sendiri tidak di-deploy,
hanya dokumentasinya).

## Men-deploy API-nya (belum dilakukan)

API ini belum di-deploy ke mana pun. Pilihan yang masuk akal saat nanti
dibutuhkan:

1. **VPS / Coolify** — build biner Go lalu jalankan sebagai service
   (`PORT` diatur lewat env). Paling sederhana.
2. **Fly.io / Railway** — cocok untuk API kecil; butuh `Dockerfile`
   (belum ada di repo — dibuat saat keputusan deploy diambil).
3. **Cloudflare Workers** — tidak cocok langsung karena ini Go + Fiber;
   butuh rewrite.

Catatan: sebelum API di-expose ke internet, tambahkan autentikasi
(atau batasi lewat firewall/VPN), karena saat ini siapa pun bisa
mengubah data lewat POST/PUT/DELETE.

## CI

Setiap push ke `main` dan `feat/*` serta setiap PR menjalankan
`.github/workflows/ci.yml`: gofmt, go vet, go test (race + coverage),
go build, golangci-lint, dan html-validate untuk `landing/`.
