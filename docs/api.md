# Referensi API — go-kos-api

Base URL lokal: `http://localhost:8080` (port bisa diubah lewat env `PORT`).

Semua respons sukses berbentuk:

```json
{ "status": "success", "data": ... }
```

Semua error berbentuk:

```json
{ "error": "pesan" }
```

## Objek Kos

| Field | Tipe | Keterangan |
|---|---|---|
| `id` | int | Diisi server, tidak boleh dikirim saat create |
| `nama` | string | Wajib, tidak boleh kosong |
| `harga` | int | Wajib, harus > 0 (rupiah per bulan) |
| `lokasi` | string | Opsional |
| `fasilitas` | string | Opsional |

## GET /health

Cek server hidup.

```bash
curl localhost:8080/health
```

```json
{ "status": "ok" }
```

## GET /api/kos

Daftar semua kos, urut berdasarkan ID.

```bash
curl localhost:8080/api/kos
```

### Filter

| Parameter | Contoh | Arti |
|---|---|---|
| `lokasi` | `?lokasi=kampus` | Substring pada field lokasi, case-insensitive |
| `min_harga` | `?min_harga=500000` | Harga ≥ nilai |
| `max_harga` | `?max_harga=1000000` | Harga ≤ nilai |

Parameter bisa digabung: `?lokasi=kampus&max_harga=1000000`.

```bash
curl "localhost:8080/api/kos?lokasi=kampus&max_harga=1000000"
```

Nilai `min_harga`/`max_harga` yang bukan angka → `400 {"error":"min_harga harus angka"}`.

## GET /api/kos/:id

```bash
curl localhost:8080/api/kos/1
```

- `200` — kos ditemukan.
- `404 {"error":"Kos tidak ditemukan"}` — id tidak ada.
- `400 {"error":"id harus angka"}` — id bukan angka (mis. `/api/kos/abc`).

## POST /api/kos

```bash
curl -X POST localhost:8080/api/kos \
  -H 'Content-Type: application/json' \
  -d '{"nama":"Kos Pelajar","harga":600000,"lokasi":"Dekat Kampus","fasilitas":"Wi-Fi"}'
```

- `201` — kos dibuat, respons berisi data baru termasuk `id`.
- `400` — body bukan JSON valid, nama kosong, atau harga ≤ 0.

ID baru = ID terbesar yang ada + 1.

## PUT /api/kos/:id

Mengganti seluruh data kos. Body dan validasi sama seperti POST.

```bash
curl -X PUT localhost:8080/api/kos/2 \
  -H 'Content-Type: application/json' \
  -d '{"nama":"Kos Syariah Berkah","harga":850000,"lokasi":"Dekat Kampus","fasilitas":"Kipas Angin, Wi-Fi, Kasur"}'
```

- `200` — kos diubah.
- `404` — id tidak ada.
- `400` — id bukan angka atau body tidak valid.

## DELETE /api/kos/:id

```bash
curl -X DELETE localhost:8080/api/kos/3
```

- `204` — kos dihapus (tanpa body).
- `404` — id tidak ada.
- `400` — id bukan angka.

## Kode status ringkas

| Kode | Arti |
|---|---|
| 200 | OK |
| 201 | Dibuat (POST) |
| 204 | Dihapus, tanpa body (DELETE) |
| 400 | Input tidak valid |
| 404 | Kos / route tidak ditemukan |
| 500 | Error server (jarang; lihat log server) |

## Batasan

- Data hanya di memori: semua perubahan hilang saat server berhenti.
- Tanpa autentikasi dan tanpa pagination.
