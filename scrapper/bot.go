package scrapper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Anime struct {
	Name      string
	episodes  string
	status    string
	premiered string
	source    string
	genres    []string
	rating    string
	ranking   string
	poster    string
	sypnonis  string
}

type Character struct {
	Name      string
	MainImage string
	Images    []string
}

func main() {

	anime := Anime{Name: "one piece"}
	fmt.Println(anime.GenresAr())
}

//getting the link of thre anime page
func (anime *Anime) Link() string {
	var links []string
	//replacing the string white spaces to get the query
	query := strings.ReplaceAll(anime.Name, " ", "%20")
	url := fmt.Sprintf("https://myanimelist.net/search/all?q=%s", query)
	//sending a request to the url
	doc, _ := CreateDoc(url)
	//looking for the element with the class information
	doc.Find(".information").Each(func(i int, s *goquery.Selection) {
		//gettinh the link and appending the links slice
		a := s.Find("a")
		url, _ := a.Attr("href")
		links = append(links, url)

	})
	//returning the link of the first anime in the page
	return links[0]
}

//getting the sypnonis of the anime
func (anime *Anime) Sypnonis() string {
	var sypnonis string
	link := anime.Link()

	doc, _ := CreateDoc(link)
	//loop through every p tag
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		//getting the itemprop attribute
		itemprop, _ := s.Attr("itemprop")
		//if the itemprop is equal to description the it's what we're looking for
		if itemprop == "description" {
			sypnonis = s.Text()
		}
	})
	//function return
	return sypnonis

}

func (anime *Anime) Info() {
	var rating string
	var ranking string
	link := anime.Link()
	doc, _ := CreateDoc(link)
	go doc.Find(".score-label").Each(func(i int, s *goquery.Selection) {
		rating = s.Text()
		anime.rating = rating
	})
	go doc.Find("span.numbers.ranked").Each(func(i int, s *goquery.Selection) {
		rank := s.Find("strong")
		ranking = rank.Text()
		anime.ranking = ranking
	})
	go doc.Find("span").Each(func(i int, s *goquery.Selection) {
		itemprop, _ := s.Attr("itemprop")
		if itemprop == "genre" {
			anime.genres = append(anime.genres, s.Text())
		}
	})
	go doc.Find(".spaceit").Each(func(i int, s *goquery.Selection) {
		span := s.Find(".dark_text")
		if span.Text() == "Source:" {
			source_spaces := strings.ReplaceAll(s.Text(), "Source:", "")
			source_w := strings.ReplaceAll(source_spaces, "\n", "")
			source := strings.ReplaceAll(source_w, " ", "")

			anime.source = source
		}

	})
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		span := s.Find(".dark_text")
		if span.Text() == "Status:" {
			source_spaces := strings.ReplaceAll(s.Text(), "Status:", "")
			source_w := strings.ReplaceAll(source_spaces, "\n", "")
			source := strings.ReplaceAll(source_w, " ", "")

			anime.status = source
		}
		if span.Text() == "Premiered:" {
			source_spaces := strings.ReplaceAll(s.Text(), "Premiered:", "")
			source_w := strings.ReplaceAll(source_spaces, "\n", "")
			source := strings.ReplaceAll(source_w, " ", "")

			anime.premiered = source
		}
		if span.Text() == "Episodes:" {
			source_spaces := strings.ReplaceAll(s.Text(), "Episodes:", "")
			source_w := strings.ReplaceAll(source_spaces, "\n", "")
			source := strings.ReplaceAll(source_w, " ", "")

			anime.episodes = source
		}

	})
	poster := anime.Poster()
	anime.poster = poster
	sypnonis := anime.Sypnonis()
	anime.sypnonis = sypnonis

}

func (anime *Anime) Poster() string {

	query := strings.ReplaceAll(anime.Name, " ", "+")
	doc, _ := CreateDoc(fmt.Sprintf("https://animeblkom.net/search?query=%s", query))
	poster := doc.Find(".poster").First()
	a := poster.Find("img")
	href, _ := a.Attr("src")

	url := fmt.Sprintf("https://animeblkom.net%s", href)

	return url

}

func (char *Character) Pictures() {
	name := char.Name
	var links []string
	query := strings.ReplaceAll(name, " ", "%20")
	url := fmt.Sprintf("https://myanimelist.net/character.php?q=%s&cat=character", query)
	doc, _ := CreateDoc(url)
	doc.Find(".borderClass").Each(func(i int, s *goquery.Selection) {
		width, _ := s.Attr("width")
		href, _ := s.Find("a").Attr("href")
		if width == "175" {

			links = append(links, href)

		}

	})
	var images_links []string

	doc_, _ := CreateDoc(links[0])
	doc_.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("data-src")
		if src != "" {
			images_links = append(images_links, src)

		}

	})

	char.MainImage = images_links[0]
	image_link := fmt.Sprintf("%s/pictures", links[0])
	images := getImages(image_link)
	char.Images = images

}

func getImages(link string) []string {
	var images []string
	doc, _ := CreateDoc(link)
	doc.Find("table img").Each(func(i int, s *goquery.Selection) {
		img_data, _ := s.Attr("data-src")

		if strings.HasPrefix(img_data, "https://cdn.myanimelist.net/images/characters/") {
			images = append(images, img_data)
		}
	})
	return images
}

func NewCharacter(name string) Character {
	character := Character{}
	character.Name = name
	character.Pictures()
	return character

}

func NewAnime(name string) Anime {
	anime := Anime{}
	anime.Name = name
	anime.Info()
	return anime

}

func EpisodeLinkAr(name string, episode string) string {

	url := AnimeLinkAr(name, episode)
	src := EmbedLink(url)
	return src

}

func EmbedLink(ep_url string) string {
	doc, _ := CreateDoc(ep_url)
	iframe := doc.Find("iframe:first-of-type")
	src, _ := iframe.Attr("src")

	return src
}

func AnimeLinkAr(name string, episode string) string {
	var link string
	query := strings.ReplaceAll(name, " ", "+")
	doc, _ := CreateDoc(fmt.Sprintf("https://animeblkom.net/search?query=%s", query))
	poster := doc.Find(".poster").First()
	a := poster.Find("a")
	href, _ := a.Attr("href")
	link = href
	url := fmt.Sprintf("https://animeblkom.net%s/%s", link, episode)
	return url

}

func DownloadLinkAr(name string, episode string) (string, error) {
	ep_url := AnimeLinkAr(name, episode)
	doc, err := CreateDoc(ep_url)
	if err != nil {
		return "", err
	}
	iframe := doc.Find("iframe:first-of-type")
	src, _ := iframe.Attr("src")
	res_, _ := http.Get(src)
	doc_, _ := goquery.NewDocumentFromReader(res_.Body)
	video := doc_.Find("video source").First()

	link, _ := video.Attr("src")
	link_download := strings.ReplaceAll(link, "watch", "download")
	return link_download, nil
}

func (anime *Anime) SypnonisAr() string {

	query := strings.ReplaceAll(anime.Name, " ", "+")
	doc, _ := CreateDoc(fmt.Sprintf("https://animeblkom.net/search?query=%s", query))
	sypnonis := doc.Find(".story-text p").First()
	return sypnonis.Text()
}

func (anime *Anime) GenresAr() ([]string, error) {
	var genres []string

	query := strings.ReplaceAll(anime.Name, " ", "+")
	doc, _ := CreateDoc(fmt.Sprintf("https://animeblkom.net/search?query=%s", query))
	class := doc.Find(".genres").First()
	class.Find("a").Each(func(i int, s *goquery.Selection) {
		genres = append(genres, s.Text())
	})
	return genres, nil

}

func CreateDoc(link string) (*goquery.Document, error) {
	res, err := http.Get(link)
	if err != nil {
		return nil, errors.New("cannot get the page")
	}
	if res.StatusCode != 200 {
		return nil, errors.New("page not found")
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}
	return doc, nil

}
