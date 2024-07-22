package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

var months []string = []string{"styczen", "luty", "marzec", "kwiecien", "maj", "czerwiec", "lipiec", "sierpien", "wrzesien", "pazdziernik", "listopad", "grudzien"}

func main() {
	c := colly.NewCollector(
		colly.CacheDir("./.cache/"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.OnHTML(".msgbox-arkusz", func(e *colly.HTMLElement) {
		link, ok := e.DOM.Children().Attr("href")
		if !ok {
			return
		}
		if !strings.Contains(link, ".pdf") {
			return
		}

		linkEls := strings.Split(link, "/")
		originalFileName := strings.Split(linkEls[len(linkEls)-1], ".")[0]

		originalFileNameSplit := strings.Split(originalFileName, "-")

		attributes := []string{}
		subject := ""
		month := -1
		year := -1

		// normalization
		for _, part := range originalFileNameSplit {
			if part == "matematyczne" {
				part = "matematyka"
			} else if part == "fizykochemiczne" {
				part = "fizykochemia"
			} else if part == "fizyczne" {
				part = "fizyka"
			} else if part == "chemiczne" {
				part = "chemia"
			}

			// ignore/use for directory picking
			if slices.Contains([]string{"poziom", "nowa", "odpowiedzi", "transkrypcja", "maturalny", "matura", "podstawowa", "podstawowy", "arkusz", "przykladowy", "cke", "egzamin", "wstepny", "na", "i", "stale"}, part) {
				continue
			}

			if month == -1 {
				for i, m := range months {
					if m == part {
						month = i
						break
					}
				}
				if month != -1 {
					continue
				}
			}

			// attributes
			if slices.Contains([]string{"stara", "probna", "rozszerzona", "rozszerzony", "poprawkowa", "era", "operon", "aneks", "studia", "dwujezyczna", "wzory", "tablice", "mapa", "informator"}, part) {
				if part == "era" {
					attributes = append(attributes, "nowa_era")
				} else if part == "roszerzony" {
					attributes = append(attributes, "rozszerzona")
				} else {
					attributes = append(attributes, part)
				}

				continue
			}

			y, err := strconv.ParseInt(part, 10, 32)
			if err == nil {
				year = int(y)
				continue
			}
			subject += part
			subject += "_"
		}

		if subject != "" {
			subject = subject[0 : len(subject)-1]
		}

		slices.Sort(attributes)
		if month != -1 {
			attributes = append([]string{fmt.Sprintf("%02d", month)}, attributes...)
		}
		if year != -1 {
			attributes = append([]string{strconv.Itoa(year)}, attributes...)
		}

		fileName := strings.Join(attributes, "-")

		destination := "arkusze/%s/pytania"
		if strings.Contains(link, "odpowiedzi") {
			destination = "arkusze/%s/odpowiedzi"
		} else if strings.Contains(link, "transkrypcja") {
			destination = "arkusze/%s/transkrypcje"
		} else if strings.Contains(link, "informator") || strings.Contains(link, "wzory") || strings.Contains(link, "tablice") || strings.Contains(link, "mapa") {
			destination = "arkusze/%s/dodatki"
		}
		destination = fmt.Sprintf(destination, subject)
		err := os.MkdirAll(destination, 0700)
		if err != nil {
			panic(err)
		}

		fileName += ".pdf"

		destinationFile := path.Join(destination, fileName)
		_, err = os.Stat(destinationFile)

		if os.IsNotExist(err) {
			err = DownloadFile(destinationFile, link)
			if err != nil {
				panic(err)
			}

			fmt.Printf("downloaded %s\n", destinationFile)
		}
	})

	c.OnHTML(".post-right", func(e *colly.HTMLElement) {
		link, ok := e.DOM.Children().Children().Attr("href")
		if ok {
			if strings.Contains(link, "matura") {
				c.Visit(link)
			}
		}
	})

	c.OnHTML(".wt-btn", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" {
			c.Visit(link)
		}
	})

	c.Visit("https://arkusze.pl/")

	c.Wait()

}

func DownloadFile(path string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
