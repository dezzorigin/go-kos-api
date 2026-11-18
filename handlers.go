package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// newApp merakit aplikasi Fiber dengan semua route.
// Dipisah dari main() supaya bisa diuji lewat app.Test tanpa menyalakan server.
func newApp(store *Store) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/api/kos", func(c *fiber.Ctx) error {
		lokasi := strings.TrimSpace(c.Query("lokasi"))
		var minHarga, maxHarga *int
		if v := c.Query("min_harga"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "min_harga harus angka"})
			}
			minHarga = &n
		}
		if v := c.Query("max_harga"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "max_harga harus angka"})
			}
			maxHarga = &n
		}
		return c.JSON(fiber.Map{"status": "success", "data": store.Filter(lokasi, minHarga, maxHarga)})
	})

	app.Get("/api/kos/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id harus angka"})
		}
		k, ok := store.Get(id)
		if !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kos tidak ditemukan"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": k})
	})

	app.Post("/api/kos", func(c *fiber.Ctx) error {
		var newKos Kos
		if err := c.BodyParser(&newKos); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "body harus JSON yang valid"})
		}
		if msg := validateKos(newKos); msg != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
		}
		created := store.Create(newKos)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": created})
	})

	app.Put("/api/kos/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id harus angka"})
		}
		var upd Kos
		if err := c.BodyParser(&upd); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "body harus JSON yang valid"})
		}
		if msg := validateKos(upd); msg != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
		}
		k, ok := store.Update(id, upd)
		if !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kos tidak ditemukan"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": k})
	})

	app.Delete("/api/kos/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id harus angka"})
		}
		if !store.Delete(id) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kos tidak ditemukan"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	return app
}

// validateKos memeriksa field wajib. Mengembalikan pesan error, atau "" bila valid.
func validateKos(k Kos) string {
	if strings.TrimSpace(k.Nama) == "" {
		return "nama wajib diisi"
	}
	if k.Harga <= 0 {
		return "harga harus lebih dari 0"
	}
	return ""
}
