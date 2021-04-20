package main

import (
	"github.com/gofiber/fiber/v2"
)

func Download(c *fiber.Ctx) error {

	link, err := DownloadLinkAr(c.Query("anime"), c.Query("episode"))

	if err != nil {
		return c.JSON(&fiber.Map{
			"error": "cannot get the episode",
		})
	}

	image, name := AnimeAvarar(c.Query("anime"))
	return c.JSON(&fiber.Map{
		"link":  link,
		"image": image,
		"name":  name,
	})
}

func Watch(c *fiber.Ctx) error {

	link := EpisodeLinkAr(c.Query("anime"), c.Query("episode"))
	image, name := AnimeAvarar(c.Query("anime"))

	return c.JSON(&fiber.Map{
		"link":  link,
		"image": image,
		"name":  name,
	})

}

func Search(c *fiber.Ctx) error {
	name := c.Params("anime")
	animes := SearchAnime(name)
	return c.JSON(animes)
}
func AnimeByName(c *fiber.Ctx) error {
	name := c.Params("anime")
	anime := SearchFetch(name)
	return c.JSON(anime)
}
func Handler(app *fiber.App) {
	app.Get("/anime/download", Download)
	app.Get("/anime/watch", Watch)
	app.Get("/anime/:anime", AnimeByName)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"Watch":        "/anime/watch",
			"Download":     "/anime/downloads",
			"Query params": "anime:String , Episode:Integer",
		})
	})
	app.Get("/anime/search/:anime", Search)
}

func main() {

	app := fiber.New()
	Handler(app)
	app.Listen(":3000")
}
