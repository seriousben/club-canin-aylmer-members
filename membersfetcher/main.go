package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseMembersPage(page string) []member {
	// Request the HTML page.
	res, err := http.Get(page)
	if err != nil {
		log.Fatal("error getting member page", page, err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ms []member
	// Find the review items
	doc.Find("div.galleryInnerImageHolder").Each(func(i int, s *goquery.Selection) {
		text := s.Find("div.galleryCaptionInnerText").Text()
		img, exists := s.Find("a").Attr("href")
		if !exists {
			return
		}
		req, err := http.Head("https://www.clubcaninaylmer.com" + img)
		if err != nil {
			log.Printf("error head img: %v", err)
			return
		}
		url := "https://www.clubcaninaylmer.com" + img
		m := member{Name: text, ImageURL: url, ImageUploadedAt: req.Header.Get("Last-Modified")}
		// fmt.Println(m)
		ms = append(ms, m)

	})
	return ms
}

type sitemap struct {
	URLS []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc string `xml:"loc"`
}

type data struct {
	Members []member `json:"members"`
}

type member struct {
	Name            string `json:"name"`
	ImageURL        string `json:"imageUrl"`
	ImageUploadedAt string `json:"imageUploadedAt"`
}

func getMembersPageURLs() ([]string, error) {
	// Request the HTML page.
	res, err := http.Get("https://www.clubcaninaylmer.com/sitemap.xml")
	if err != nil {
		log.Fatal("error getting site map", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	var s sitemap
	err = xml.Unmarshal(b, &s)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}

	var us []string
	for _, u := range s.URLS {
		if strings.Contains(u.Loc, "membresmembers") {
			us = append(us, u.Loc)
		}
	}
	return us, nil
}

func main() {
	pageMembers, err := getMembersPageURLs()
	if err != nil {
		log.Fatal("error getting member pages", err)
	}

	var ms []member
	for _, p := range pageMembers {
		ms = append(ms, parseMembersPage(p)...)
	}

	d := data{Members: ms}

	e := json.NewEncoder(os.Stdout)
	if err := e.Encode(d); err != nil {
		log.Fatal("error encoding data pages", err)
	}
}
