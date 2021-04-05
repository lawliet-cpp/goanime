package main

import (
	"rest_api/scrapper"

	"github.com/gofiber/fiber/v2"
)

func Download(c *fiber.Ctx) error {
	link, err := scrapper.DownloadLinkAr(c.Query("anime"), c.Query("episode"))
	if err != nil {
		return c.JSON(&fiber.Map{
			"error": "cannot get the episode",
		})
	}

	return c.JSON(&fiber.Map{
		"link": link,
	})
}

func Watch(c *fiber.Ctx) error {

	link := scrapper.EpisodeLinkAr(c.Query("anime"), c.Query("episode"))

	return c.JSON(&fiber.Map{
		"link": link,
	})

}

func Handler(app *fiber.App) {
	app.Get("/download", Download)
	app.Get("/watch", Watch)
}

func main() {
	app := fiber.New()
	Handler(app)
	app.Listen(":3000")
}
