package main

import (
	"sort"
	"strings"
	"sync"
)

// Kos adalah satu listing kost.
type Kos struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama"`
	Harga     int    `json:"harga"` // harga sewa per bulan, dalam rupiah
	Lokasi    string `json:"lokasi"`
	Fasilitas string `json:"fasilitas"`
}

// Store menyimpan daftar kos di memori (tanpa database).
// Data hilang saat proses berhenti — ini disengaja untuk v1.
type Store struct {
	mu    sync.RWMutex
	items []Kos
}

// NewStore membuat store berisi 2 data contoh bawaan repo.
func NewStore() *Store {
	return &Store{items: []Kos{
		{ID: 1, Nama: "Kos Mentari Hype", Harga: 1500000, Lokasi: "Sekitar Banda", Fasilitas: "AC, Wi-Fi, Kamar Mandi Dalam"},
		{ID: 2, Nama: "Kos Syariah Berkah", Harga: 800000, Lokasi: "Dekat Kampus", Fasilitas: "Kipas Angin, Wi-Fi, Kasur"},
	}}
}

// List mengembalikan salinan semua kos, urut berdasarkan ID.
func (s *Store) List() []Kos {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Kos, len(s.items))
	copy(out, s.items)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Filter mengembalikan kos yang cocok dengan semua filter yang diisi.
// lokasi dicocokkan sebagai substring (case-insensitive).
func (s *Store) Filter(lokasi string, minHarga, maxHarga *int) []Kos {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Kos
	for _, k := range s.items {
		if lokasi != "" && !strings.Contains(strings.ToLower(k.Lokasi), strings.ToLower(lokasi)) {
			continue
		}
		if minHarga != nil && k.Harga < *minHarga {
			continue
		}
		if maxHarga != nil && k.Harga > *maxHarga {
			continue
		}
		out = append(out, k)
	}
	if out == nil {
		out = []Kos{}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Get mencari kos berdasarkan ID.
func (s *Store) Get(id int) (Kos, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, k := range s.items {
		if k.ID == id {
			return k, true
		}
	}
	return Kos{}, false
}

// Create menambahkan kos baru dengan ID = max ID yang ada + 1.
func (s *Store) Create(k Kos) Kos {
	s.mu.Lock()
	defer s.mu.Unlock()
	maxID := 0
	for _, item := range s.items {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	k.ID = maxID + 1
	s.items = append(s.items, k)
	return k
}

// Update mengganti data kos dengan ID yang diberikan.
func (s *Store) Update(id int, k Kos) (Kos, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, item := range s.items {
		if item.ID == id {
			k.ID = id
			s.items[i] = k
			return k, true
		}
	}
	return Kos{}, false
}

// Delete menghapus kos berdasarkan ID.
func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, item := range s.items {
		if item.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}
