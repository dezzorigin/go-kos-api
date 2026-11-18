# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-11-18

### Added
- `GET /health` health-check endpoint
- Filter `GET /api/kos?lokasi=`, `?min_harga=`, `?max_harga=`
- `PUT /api/kos/:id` (update) dan `DELETE /api/kos/:id` (hapus) — CRUD lengkap
- Validasi input: nama wajib diisi, harga harus > 0, id harus angka,
  body harus JSON valid (400 bila tidak)
- Port bisa diatur lewat env `PORT` (default 8080)
- Tes otomatis endpoint (`main_test.go`, coverage 80,6%, `go test -race`)
- CI GitHub Actions: gofmt, go vet, go test race+coverage, go build,
  golangci-lint, html-validate untuk landing
- `.gitignore` Go + konfigurasi golangci-lint
- Dokumentasi: `docs/api.md`, `docs/architecture.md`, `docs/deployment.md`,
  `docs/diagram.svg`
- Halaman dokumentasi statis (`landing/`) di Cloudflare Pages
- README + PRD ditulis ulang secara jujur dari kondisi kode sebenarnya

### Changed
- Kode dipecah: `main.go` (entrypoint), `handlers.go` (route + validasi),
  `store.go` (penyimpanan in-memory thread-safe)
- ID kos baru = max ID + 1 (sebelumnya len + 1, bisa tabrakan setelah hapus)
- Penanganan error terpusat lewat Fiber `ErrorHandler`

### Fixed
- `GET /api/kos/abc` sekarang 400 (sebelumnya 404 karena id gagal parse
  lalu dicari sebagai 0)
- `POST /api/kos` dengan body bukan JSON sekarang 400 (sebelumnya error
  diabaikan diam-diam)

### Removed
- Klaim fiktif di README/PRD/CHANGELOG lama (search, rate limiting, Docker,
  health check, database, `npm install`) yang tidak pernah ada di kode

## [1.0.0] - 2026-11-04

### Added
- Proyek Go + Fiber awal
- `GET /api/kos` — daftar 2 kos contoh
- `GET /api/kos/:id` — satu kos berdasarkan ID
- `POST /api/kos` — tambah kos (tanpa validasi)
- README, PRD, LICENSE MIT
