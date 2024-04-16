package crawl

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func FindUrls(url string, urls []string) (result []string, err error) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error at findUrl %v \n", err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Status code not okay: %v \n", res.StatusCode)
		return
	}
	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Printf("Parsing error could not gt data in right format %v \n", err.Error())
		return
	}
	
	return Visit(doc, urls), err
}

func Visit(doc *html.Node, urls []string) []string {
	if doc.Type == html.ElementNode && doc.Data == "a" {
		var title string
		var url string
		for _, a := range doc.Attr {
			if a.Key == "href" {
				url = a.Val
			}
			if a.Key == "title" {
				title = a.Val
			}
			result := fmt.Sprintf("%v&%v", title, url)
			var exist = false
			for _, v := range urls {
				if v == result {
					exist = true
				}
			}
			if !exist {
				urls = append(urls, result)
			}
			
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		urls = Visit(c, urls)
	}
	return urls
}
