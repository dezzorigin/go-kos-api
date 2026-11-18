# PRD — go-kos-api

## 1. Informasi Proyek

| Field | Nilai |
|---|---|
| Nama Proyek | go-kos-api |
| Tech Stack | Go 1.26+, Fiber v2 |
| Status | v1.1.0 — berjalan, teruji |
| Penanggung Jawab | Nabil Aulia Dika (dezzorigin) |
| Versi Dokumen | 1.1 |
| Tanggal | 2026-11-18 |

## 2. Latar Belakang

Mencari kos biasanya berarti membuka banyak grup chat dan marketplace lalu
membandingkan harga satu per satu. API ini menyediakan satu sumber daftar kos
yang bisa di-query berdasarkan lokasi dan rentang harga, supaya aplikasi
pencari kos (web/mobile) punya backend yang sederhana dan konsisten.

## 3. Masalah yang Diselesaikan

- Daftar kos tersebar dan sulit difilter berdasarkan harga/lokasi.
- Belum ada API kecil yang bisa dipakai langsung untuk prototipe aplikasi
  pencari kos tanpa menyiapkan database.

## 4. Tujuan & Sasaran

- Menyediakan CRUD listing kos lewat REST API yang bisa dijalankan dengan
  satu perintah (`go run .`).
- Filter lokasi + harga bekerja benar dan teruji.
- Respons dan error konsisten sehingga mudah dipakai client mana pun.

## 5. Ruang Lingkup

### In Scope (v1.1.0)
- CRUD listing kos: list, get by id, create, update, delete.
- Filter `lokasi` (substring, case-insensitive), `min_harga`, `max_harga`.
- Health check `GET /health`.
- Validasi input: nama wajib, harga > 0, id harus angka, body harus JSON valid.
- Tes otomatis endpoint (coverage 80,6%), CI hijau.
- Halaman dokumentasi statis di Cloudflare Pages.

### Out of Scope
- Database persisten (data masih di memori, hilang saat server berhenti).
- Autentikasi/otorisasi.
- Upload foto kos.
- Rate limiting, pagination.
- Deployment API-nya sendiri (hanya halaman dokumentasi yang di-deploy).

## 6. Fitur & Persyaratan Fungsional

| # | Fitur | Kriteria |
|---|---|---|
| 1 | `GET /health` | 200 + `{"status":"ok"}` |
| 2 | `GET /api/kos` | 200, daftar urut ID, `{"status":"success","data":[...]}` |
| 3 | Filter lokasi | `?lokasi=kampus` cocok substring tanpa peduli huruf besar/kecil |
| 4 | Filter harga | `?min_harga=` / `?max_harga=`; nilai bukan angka → 400 |
| 5 | `GET /api/kos/:id` | 200 bila ada; 404 bila tidak ada; 400 bila id bukan angka |
| 6 | `POST /api/kos` | 201 + data baru; nama kosong / harga ≤ 0 / body invalid → 400 |
| 7 | `PUT /api/kos/:id` | 200 bila ada; 404 bila tidak ada; validasi sama seperti POST |
| 8 | `DELETE /api/kos/:id` | 204 bila ada; 404 bila tidak ada |
| 9 | Penomoran ID | ID baru = max ID + 1, tidak tabrakan setelah ada yang dihapus |
| 10 | Konfigurasi port | Env `PORT`, default 8080 |

## 7. Persyaratan Non-Fungsional

- **Performa:** respons < 5 ms untuk dataset kecil di localhost.
- **Keamanan:** tidak ada secret di repo; `.env` di-`.gitignore`.
- **Maintainability:** handler, store, dan entrypoint terpisah; `gofmt` +
  `go vet` + `golangci-lint` bersih di CI.
- **Portabilitas:** jalan di Linux/macOS/Windows dengan Go 1.26+.

## 8. Arsitektur & Tech Stack

Go + Fiber v2, penyimpanan in-memory dengan `sync.RWMutex`.
Lihat [docs/architecture.md](docs/architecture.md) dan
[docs/diagram.svg](docs/diagram.svg).

## 9. Milestone / Roadmap

| Fase | Status |
|---|---|
| v1.0.0 — 3 endpoint dasar (list, get, create), data contoh | ✅ |
| v1.1.0 — CRUD lengkap, filter, validasi, tes, CI, docs, landing | ✅ (2026-11-18) |
| Berikutnya — database persisten (SQLite/Postgres) | ⏳ belum diputuskan |
| Berikutnya — autentikasi + deploy API live | ⏳ belum diputuskan |

## 10. Risiko & Mitigasi

| Risiko | Mitigasi |
|---|---|
| Data hilang saat restart | Dinyatakan jujur di README; database masuk roadmap |
| Tanpa auth, data bisa diubah siapa saja | Jangan expose ke internet tanpa auth di depan |
| Dependency Fiber usang | `go get -u` berkala + CI build tiap push |

## 11. Kriteria Keberhasilan

- Semua endpoint di §6 lolos play-test curl end-to-end.
- `go test -race` hijau, coverage ≥ 80%.
- CI hijau (fmt, vet, test, build, golangci-lint, html-validate).
- Halaman dokumentasi live di Cloudflare Pages.

## 12. Lampiran

- Struktur: lihat README.md
- Referensi API: docs/api.md
- Lisensi: MIT (LICENSE)
