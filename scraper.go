package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
		return nil, errors.New("an error ocured")
	}
	return doc, nil

}

func EpisodeLinkAr(name string, episode string) string {

	url := AnimeLinkAr(name, episode)
	src := EmbedLink(url)
	return src

}

func EmbedLink(ep_url string) string {
	doc, err := CreateDoc(ep_url)
	if err != nil {
		return "Please lookup with a valid query"
	}
	iframe := doc.Find("iframe:first-of-type")
	src, _ := iframe.Attr("src")

	return src
}

func AnimeLinkAr(name string, episode string) string {
	var link string
	query := strings.ReplaceAll(name, " ", "+")
	doc, err := CreateDoc(fmt.Sprintf("https://animeblkom.net/search?query=%s", query))
	if err != nil {
		return ""
	}
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
	doc, err := CreateDoc(fmt.Sprintf("https://animeblkom.net/search?query=%s", query))
	if err != nil {
		return ""
	}
	sypnonis := doc.Find(".story-text p").First()
	return sypnonis.Text()
}

func (anime *Anime) GenresAr() ([]string, error) {
	var genres []string

	query := strings.ReplaceAll(anime.Name, " ", "+")
	doc, err := CreateDoc(fmt.Sprintf("https://animeblkom.net/search?query=%s", query))
	if err != nil {
		return nil, errors.New("cannot get the genres")
	}
	class := doc.Find(".genres").First()
	class.Find("a").Each(func(i int, s *goquery.Selection) {
		genres = append(genres, s.Text())
	})
	return genres, nil

}

func AnimeAvarar(name string) (string, string) {
	query := strings.ReplaceAll(name, " ", "+")
	doc, err := CreateDoc(fmt.Sprintf("https://animeblkom.net/search?query=%s", query))
	if err != nil {
		return "", ""
	}
	poster := doc.Find(".poster").First()
	img := poster.Find(".lazy").First()
	src, _ := img.Attr("data-original")
	alt, _ := img.Attr("alt")
	full_name := strings.ReplaceAll(alt, "poster", "")
	link := fmt.Sprintf("https://animeblkom.net%s", src)
	return link, full_name

}
