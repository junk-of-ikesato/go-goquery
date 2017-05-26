package main


import (
  "fmt"
	"strings"
  "net/http"
	"os"
  "github.com/PuerkitoBio/goquery"
	"bufio"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"regexp"
)


func main() {
	res, err := http.Get("https://krsw.2ch.net/test/read.cgi/ghard/1495799870")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer res.Body.Close()

	utfBody := transform.NewReader(bufio.NewReader(res.Body), japanese.ShiftJIS.NewDecoder())
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	doc, err := goquery.NewDocumentFromReader(utfBody)
  if err != nil {
		fmt.Println(err)
		os.Exit(3)
  }


  doc.Find("div.post").Each(func(i int, s *goquery.Selection) {
		no := strings.TrimSpace(s.Find("div.meta span.number").Text())
		name := strings.TrimSpace(s.Find("div.meta span.name").Text())
		date := strings.TrimSpace(s.Find("div.meta span.date").Text())
		uid := strings.TrimSpace(s.Find("div.meta span.uid").Text())

		m := s.Find("div.message span")
		m.Find("a").Each(func(_ int, link *goquery.Selection) {
			text := strings.TrimSpace(link.Text())
			link.ReplaceWithHtml(text)
		})

		text,_ := m.Html()
    re, _ := regexp.Compile("\\s*\\<br/\\>\\s*")
    text = re.ReplaceAllString(text, "\n")
		text = strings.TrimSpace(text)
		re, _ = regexp.Compile("\\<[\\S\\s]+?\\>") // remove all tags
    text = re.ReplaceAllString(text, "")
		text = strings.TrimSpace(text)
    fmt.Println("-------------------------------------------")
    fmt.Printf("%v %v: %v (%v)\n", date, no, name, uid)
    fmt.Printf("%v\n", text)
  })
}
