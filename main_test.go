package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func doRequest(t *testing.T, app interface {
	Test(*http.Request, ...int) (*http.Response, error)
}, method, path, body string) *http.Response {
	t.Helper()
	var rdr io.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, rdr)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	return resp
}

func readBody(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("baca body: %v", err)
	}
	if len(b) == 0 {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("body bukan JSON (%s): %v", string(b), err)
	}
	return m
}

func TestHealth(t *testing.T) {
	resp := doRequest(t, newApp(NewStore()), http.MethodGet, "/health", "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health: status = %d, mau 200", resp.StatusCode)
	}
	if got := readBody(t, resp)["status"]; got != "ok" {
		t.Fatalf("health: status = %v, mau ok", got)
	}
}

func TestListKos(t *testing.T) {
	resp := doRequest(t, newApp(NewStore()), http.MethodGet, "/api/kos", "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list: status = %d, mau 200", resp.StatusCode)
	}
	data, ok := readBody(t, resp)["data"].([]any)
	if !ok || len(data) != 2 {
		t.Fatalf("list: jumlah data = %d, mau 2", len(data))
	}
}

func TestListKosFilter(t *testing.T) {
	app := newApp(NewStore())

	resp := doRequest(t, app, http.MethodGet, "/api/kos?lokasi=kampus", "")
	data := readBody(t, resp)["data"].([]any)
	if len(data) != 1 {
		t.Fatalf("filter lokasi=kampus: %d hasil, mau 1", len(data))
	}

	resp = doRequest(t, app, http.MethodGet, "/api/kos?max_harga=1000000", "")
	data = readBody(t, resp)["data"].([]any)
	if len(data) != 1 {
		t.Fatalf("filter max_harga=1000000: %d hasil, mau 1", len(data))
	}

	resp = doRequest(t, app, http.MethodGet, "/api/kos?min_harga=2000000", "")
	data = readBody(t, resp)["data"].([]any)
	if len(data) != 0 {
		t.Fatalf("filter min_harga=2000000: %d hasil, mau 0", len(data))
	}

	resp = doRequest(t, app, http.MethodGet, "/api/kos?max_harga=bukan-angka", "")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("filter max_harga invalid: status = %d, mau 400", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestGetKosByID(t *testing.T) {
	app := newApp(NewStore())

	resp := doRequest(t, app, http.MethodGet, "/api/kos/1", "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get id 1: status = %d, mau 200", resp.StatusCode)
	}
	if got := readBody(t, resp)["data"].(map[string]any)["nama"]; got != "Kos Mentari Hype" {
		t.Fatalf("get id 1: nama = %v", got)
	}

	resp = doRequest(t, app, http.MethodGet, "/api/kos/99", "")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("get id 99: status = %d, mau 404", resp.StatusCode)
	}
	resp.Body.Close()

	resp = doRequest(t, app, http.MethodGet, "/api/kos/abc", "")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("get id abc: status = %d, mau 400", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestCreateKos(t *testing.T) {
	app := newApp(NewStore())

	resp := doRequest(t, app, http.MethodPost, "/api/kos",
		`{"nama":"Kos Pelajar","harga":600000,"lokasi":"Dekat Kampus","fasilitas":"Wi-Fi"}`)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: status = %d, mau 201", resp.StatusCode)
	}
	if got := readBody(t, resp)["data"].(map[string]any)["id"]; got != float64(3) {
		t.Fatalf("create: id = %v, mau 3", got)
	}

	// Validasi: nama kosong dan harga 0 ditolak.
	resp = doRequest(t, app, http.MethodPost, "/api/kos", `{"nama":"","harga":500000}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("create nama kosong: status = %d, mau 400", resp.StatusCode)
	}
	resp.Body.Close()

	resp = doRequest(t, app, http.MethodPost, "/api/kos", `{"nama":"Kos X","harga":0}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("create harga 0: status = %d, mau 400", resp.StatusCode)
	}
	resp.Body.Close()

	// Body bukan JSON ditolak.
	resp = doRequest(t, app, http.MethodPost, "/api/kos", `{bukan json`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("create body invalid: status = %d, mau 400", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestUpdateAndDeleteKos(t *testing.T) {
	app := newApp(NewStore())

	resp := doRequest(t, app, http.MethodPut, "/api/kos/2",
		`{"nama":"Kos Syariah Berkah","harga":850000,"lokasi":"Dekat Kampus","fasilitas":"Kipas Angin, Wi-Fi, Kasur"}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update: status = %d, mau 200", resp.StatusCode)
	}
	if got := readBody(t, resp)["data"].(map[string]any)["harga"]; got != float64(850000) {
		t.Fatalf("update: harga = %v, mau 850000", got)
	}

	resp = doRequest(t, app, http.MethodPut, "/api/kos/99",
		`{"nama":"Kos X","harga":100000}`)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("update id 99: status = %d, mau 404", resp.StatusCode)
	}
	resp.Body.Close()

	resp = doRequest(t, app, http.MethodDelete, "/api/kos/2", "")
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete: status = %d, mau 204", resp.StatusCode)
	}
	resp.Body.Close()

	resp = doRequest(t, app, http.MethodGet, "/api/kos/2", "")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("get setelah delete: status = %d, mau 404", resp.StatusCode)
	}
	resp.Body.Close()

	// ID baru setelah hapus tidak boleh tabrakan: max ID yang ada + 1.
	resp = doRequest(t, app, http.MethodPost, "/api/kos",
		`{"nama":"Kos Baru","harga":700000}`)
	if got := readBody(t, resp)["data"].(map[string]any)["id"]; got != float64(2) {
		t.Fatalf("create setelah delete: id = %v, mau 2", got)
	}
}

func TestStoreFilterCaseInsensitive(t *testing.T) {
	s := NewStore()
	if got := len(s.Filter("BANDA", nil, nil)); got != 1 {
		t.Fatalf("filter BANDA: %d hasil, mau 1", got)
	}
}
