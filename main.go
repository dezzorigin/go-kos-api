package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Kos struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama"`
	Harga     int    `json:"harga"`
	Lokasi    string `json:"lokasi"`
	Fasilitas string `json:"fasilitas"`
}

var daftarKos = []Kos{
	{ID: 1, Nama: "Kos Mentari Hype", Harga: 1500000, Lokasi: "Sekitar Banda", Fasilitas: "AC, Wi-Fi, Kamar Mandi Dalam"},
	{ID: 2, Nama: "Kos Syariah Berkah", Harga: 800000, Lokasi: "Dekat Kampus", Fasilitas: "Kipas Angin, Wi-Fi, Kasur"},
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/api/kos", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "success", "data": daftarKos})
	})

	app.Get("/api/kos/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		for _, k := range daftarKos {
			if k.ID == id {
				return c.JSON(fiber.Map{"status": "success", "data": k})
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Kos tidak ditemukan"})
	})

	app.Post("/api/kos", func(c *fiber.Ctx) error {
		var newKos Kos
		c.BodyParser(&newKos)
		newKos.ID = len(daftarKos) + 1
		daftarKos = append(daftarKos, newKos)
		return c.Status(201).JSON(fiber.Map{"status": "success", "data": newKos})
	})

	app.Listen(":8080")
}
