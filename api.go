package main

import (
	"github.com/gofiber/fiber/v2"
)



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
	
	app.Get("/anime/:anime", AnimeByName)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"List":        "/anime/search/:query",
			"Single Result":     "/anime/:query",
			
		})
	})
	app.Get("/anime/search/:anime", Search)
}

func main() {

	app := fiber.New()
	Handler(app)
	app.Listen(":3000")
}
