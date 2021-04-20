package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type Anime struct {
	Name       string
	Episodes   string
	ImageUrl   string
	Status     string
	Season     string
	Synopsis   string
	Source     string
	Genres     []string
	Score      string
	Ranked     string
	Popularity string
	Studios    string
	Type       string
	Characters []string
}

func Results(name string) ([]string, []string, []string) {
	var episodes []string
	var names []string
	var stories []string
	c := colly.NewCollector(
		colly.Async(true),
	)

	query := strings.ReplaceAll(name, " ", "+")

	c.OnHTML("strong", func(e *colly.HTMLElement) {
		names = append(names, e.Text)

	})
	c.OnHTML(".pt4", func(e *colly.HTMLElement) {

		stories = append(stories, e.Text)

	})
	c.OnHTML(".borderClass", func(e *colly.HTMLElement) {
		if e.Attr("width") == "40" {
			episodes_num := strings.ReplaceAll(e.Text, "\n", "")
			episodes = append(episodes, episodes_num)

		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(fmt.Sprintf("https://myanimelist.net/anime.php?q=%s&cat=anime", query))
	c.Wait()
	return episodes, names, stories

}

func SearchAnime(name string) []Anime {
	var animes []Anime

	episodes, names, stories := Results(name)
	for i := 0; i < len(names); i++ {
		var anime Anime
		anime.Name = names[i]

		anime.Episodes = strings.ReplaceAll(episodes[i+1], " ", "")
		anime.Synopsis = stories[i]
		animes = append(animes, anime)

	}
	return animes

}

func GetFirstLink(name string) string {
	var links []string
	c := colly.NewCollector(
		colly.Async(true),
	)
	c.OnHTML(".hoverinfo_trigger", func(e *colly.HTMLElement) {
		links = append(links, e.Attr("href"))
	})

	query := strings.ReplaceAll(name, " ", "+")
	c.Visit(fmt.Sprintf("https://myanimelist.net/anime.php?q=%s&cat=anime", query))
	c.Wait()
	return links[0]
}

func SearchFetch(name string) Anime {
	var anime Anime
	url := GetFirstLink(name)
	var characters []string
	c := colly.NewCollector(
		colly.Async(true),
	)
	var images []string
	c.OnHTML(".lazyload", func(e *colly.HTMLElement) {
		image_url := e.Attr("data-src")
		if strings.Contains(image_url, "/characters/") {
			name := strings.ReplaceAll(e.Attr("alt"), ",", "")
			characters = append(characters, name)
		}
		images = append(images, image_url)

	})
	c.OnHTML("p", func(e *colly.HTMLElement) {
		if e.Attr("itemprop") == "description" {

			anime.Synopsis = e.Text
		}

	})
	c.OnHTML("div", func(e *colly.HTMLElement) {

		if e.ChildText(".dark_text") == "Episodes:" {
			_info := strings.ReplaceAll(e.Text, "Episodes:", "")
			__info := strings.ReplaceAll(_info, "\n", "")
			info := strings.ReplaceAll(__info, " ", "")
			anime.Episodes = info
		}
		if e.ChildText(".dark_text") == "Source:" {
			_info := strings.ReplaceAll(e.Text, "Source:", "")
			__info := strings.ReplaceAll(_info, "\n", "")
			info := strings.ReplaceAll(__info, " ", "")
			anime.Source = info
		}
		if e.ChildText(".dark_text") == "Premiered:" {
			_info := strings.ReplaceAll(e.Text, "Premiered:", "")
			__info := strings.ReplaceAll(_info, "\n", "")
			info := strings.ReplaceAll(__info, " ", "")
			anime.Season = info
		}
		if e.ChildText(".dark_text") == "Studios:" {
			_info := strings.ReplaceAll(e.Text, "Studios:", "")
			__info := strings.ReplaceAll(_info, "\n", "")
			info := strings.ReplaceAll(__info, " ", "")
			anime.Studios = info
		}
		if e.ChildText(".dark_text") == "Genres:" {
			_info := strings.ReplaceAll(e.Text, "Genres:", "")
			__info := strings.ReplaceAll(_info, "\n", "")
			info := strings.ReplaceAll(__info, " ", "")
			genres := strings.Split(info, ",")
			anime.Genres = genres
		}

	})
	c.OnHTML("span", func(e *colly.HTMLElement) {
		if e.Attr("itemprop") == "ratingValue" {
			anime.Score = e.Text
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL)
	})

	c.Visit(url)
	c.Wait()
	image_url := images[0]
	anime.ImageUrl = image_url
	anime.Characters = characters
	return anime

}
