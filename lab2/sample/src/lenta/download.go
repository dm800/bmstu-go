package main

import (
	"net/http"
	"strconv"

	log "github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

type Item struct {
	Ref, Title string
}

func readItem(item *html.Node) *Item {
	if a := item.FirstChild; isElem(a, "a") {
		return &Item{
			Ref:   getAttr(a, "href"),
			Title: getAttr(a, "aria-label"),
		}
	}
	return nil
}

func search(node *html.Node) []*Item {
	log.Info("step 0")
	if isElem(node, "div") && getAttr(node, "id") == "app" {
		var items []*Item
		log.Info("step 1")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isDiv(c, "dTSkA6xB commercial-branding") {
				log.Info("step 2")
				for d := c.FirstChild; d != nil; d = d.NextSibling {
					if isDiv(d, "AuRBdDZg") {
						log.Info("step 3")
						for e := d.FirstChild; e != nil; e = e.NextSibling {
							if isElem(e, "section") {
								log.Info("step 4")
								for f := e.FirstChild; f != nil; f = f.NextSibling {
									if isDiv(f, "cGZPyk4_") {
										log.Info("step 5")
										for g := f.FirstChild; g != nil; g = g.NextSibling {
											if isDiv(g, "zT5wwAPN fQtJ19Ei") {
												log.Info("step 6")
												for h := g.FirstChild; h != nil; h = h.NextSibling {
													if isDiv(h, "XSvLK2D0 abGoxuyb") {
														log.Info("step 7 (must be multiple)")
														items = append(items, readItem(h))
													}
												}
												log.Info("returning total of " + strconv.Itoa(len(items)))
												return items
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}

	return nil
}

func downloadNews() []*Item {
	log.Info("sending request to rambler.ru")
	if response, err := http.Get("https://news.rambler.ru/latest/"); err != nil {
		log.Error("request to rambler.ru failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response from rambler.ru", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from rambler.ru", "error", err)
			} else {
				log.Info("HTML from rambler.ru parsed successfully")
				return search(doc)
			}
		}
	}
	return nil
}
