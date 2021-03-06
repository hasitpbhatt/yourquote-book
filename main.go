package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

// Main is complete json for /basic request
type Main struct {
	Posts []Post `json:"posts"`
	Page  Page   `json:"page"`
	User  User   `json:"user"`
}

// User contains user related details
type User struct {
	ID         string `json:"id"`
	ImageLarge string `json:"image_large"`
	Name       string `json:"name"`
	PenName    string `json:"pen_name"`
}

// Page contains pagination information
type Page struct {
	HasNext bool `json:"has_next"`
}

// Post contains the post
type Post struct {
	Text           string `json:"text"`
	Caption        string `json:"caption"`
	LikeCount      int64  `json:"like_count"`
	DateTimeString string `json:"published_datetime"`
}

func main() {
	var user string
	var start, end int
	// ccq : Hasit Bhatt
	flag.StringVar(&user, "user", "ccq", "userid")
	flag.IntVar(&start, "start", 0, "page start")
	flag.IntVar(&end, "end", 0, "page end")

	flag.Parse()

	posts, userObj := getListOfPosts(user, start, end)
	writeUser(posts, userObj)
}

func writeUser(posts []Post, user User) {
	f, err := os.Create(user.ID + ".html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(`<html><head>
	<script src="https://unpkg.com/pagedjs/dist/paged.polyfill.js"></script>
	<style type="text/css">
	@page {
		size: 148mm 210mm;
		margin-top: 10mm;
		margin-right: 20mm;
		margin-bottom: 25mm;
		margin-left: 15mm;
	
		@bottom-left {
			content: counter(page);
		}
	
		@bottom-center {
			content: string(title);
			text-transform: uppercase;
		}
	
	}
	quote {
		break-before: page;
		font-size: 1.25em;
	}
	</style>`)
	f.WriteString("<title>" + user.Name + "</title>")
	f.WriteString(`</head>`)
	f.WriteString("\n<body>\n")
	f.WriteString(fmt.Sprintf("<h1>%s</h1>", user.Name))
	for _, p := range posts {
		text := strings.ReplaceAll(p.Text, "\n", "<br/>")
		f.WriteString("<quote>")
		f.WriteString(text)
		f.WriteString("<!--<hr/>-->")
		f.WriteString("</quote>")
	}
	f.WriteString("</body></html>")
}

func getListOfPosts(user string, start, end int) ([]Post, User) {
	if start < 1 {
		start = 1
	}
	if end < 1 {
		end = math.MaxInt64
	}
	posts := []Post{}
	s := surf.NewBrowser()
	i := start
	userObj := User{}
	for i <= end {
		url := fmt.Sprintf("https://www.yourquote.in/yourquote-web/web/basic?sort=latest&userId=%s&page=%d", user, i)
		fmt.Println("Processing page", i)
		page, err := getPage(s, url)
		if err != nil {
			log.Fatal(err)
		}
		userObj = page.User
		posts = append(posts, page.Posts...)
		i++
		if !page.Page.HasNext {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return posts, userObj
}

func getPage(b *browser.Browser, url string) (*Main, error) {
	err := b.Open(url)
	if err != nil {
		return nil, err
	}
	jsonText := b.Dom().Text()
	m := Main{}
	err = json.Unmarshal([]byte(jsonText), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
