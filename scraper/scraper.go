package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

var (
	spanishMonths = map[string]string{
		"Enero":      "January",
		"Febrero":    "February",
		"Marzo":      "March",
		"Abril":      "April",
		"Mayo":       "May",
		"Junio":      "June",
		"Julio":      "July",
		"Agosto":     "August",
		"Septiembre": "September",
		"Octubre":    "October",
		"Noviembre":  "November",
		"Diciembre":  "December",
	}

	spanishDays = map[string]string{
		"Lunes":     "Monday",
		"Martes":    "Tuesday",
		"Miercoles": "Wednesday",
		"Jueves":    "Thursday",
		"Viernes":   "Friday",
		"Sabado":    "Saturday",
		"Domingo":   "Sunday",
	}
)

type Article struct {
	Title      string `json:"title"`
	Date       string `json:"date"`
	Category   string `json:"category"`
	CategoryId int    `json:"category_id"`
	Body       string `json:"body"`
	Link       string `json:"link"`
}

func parseSpanishDate(dateStr string) (time.Time, error) {

	for spanish, english := range spanishMonths {

		dateStr = strings.ReplaceAll(dateStr, spanish, english)
	}
	for spanish, english := range spanishDays {
		dateStr = strings.ReplaceAll(dateStr, spanish, english)
	}

	dateStr = strings.NewReplacer("de", "", ".", "").Replace(dateStr)
	return time.Parse("Monday, 02 January 2006 15:04 Hs", strings.TrimSpace(dateStr))
}

func setupCollectors(c *colly.Collector, links *[]string, articles *[]Article, mu *sync.Mutex, wg *sync.WaitGroup, canContinue *bool, canGoBack *bool) {
	setupMainCollector(c, links, "", canContinue, canGoBack)
	contentCollector := setupContentCollector(c, articles, mu)
	setupOnScrapedCallback(c, contentCollector, links, wg)
}

func setupMainCollector(c *colly.Collector, links *[]string, additionalClass string, canContinue *bool, canGoBack *bool) {
	class := " .noticia1" + additionalClass

	c.OnHTML("[data-link]"+class, func(e *colly.HTMLElement) {
		parent := e.DOM.Parent()

		if link := parent.AttrOr("data-link", ""); link != "" {
			*links = append(*links, link)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		if len(*links) == 0 {
			log.Println("Error: No elements found with the specified selector.")
		}
	})

	c.OnHTML(".pagination a", func(e *colly.HTMLElement) {
		if e.Text == "Siguiente" {
			*canContinue = true
		}
		if e.Text == "Anterior" {
			*canGoBack = true
		}
	})
}

func setupContentCollector(c *colly.Collector, articles *[]Article, mu *sync.Mutex) *colly.Collector {
	contentCollector := c.Clone()

	contentCollector.OnHTML(".noticia-detalle", func(e *colly.HTMLElement) {
		article := parseArticle(e)

		if article != nil {
			mu.Lock()
			*articles = append(*articles, *article)
			mu.Unlock()
		}
	})
	return contentCollector
}

func parseArticle(e *colly.HTMLElement) *Article {
	info := e.ChildText(".noticia-detalle-info")
	title := e.ChildText(".titulo2 ")
	category := e.ChildText(".titulo")
	classes := e.DOM.AttrOr("class", "")

	re := regexp.MustCompile(`categoria_(\d+)`)
	matches := re.FindStringSubmatch(classes)
	if len(matches) < 2 {
		log.Printf("Error: the categoryId was not found for the article: %s", title)
		return nil
	}

	categoryId, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Printf("Error parsing categoryId: %v", err)
		return nil
	}

	body := strings.Join(e.ChildTexts("p"), "\n\n")

	t, err := parseSpanishDate(info)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		return nil
	}

	return &Article{
		Title:      title,
		Date:       t.Format("2006-01-02 15:04:05"),
		Category:   category,
		CategoryId: categoryId,
		Body:       body,
		Link:       e.Request.URL.String(),
	}
}

func setupOnScrapedCallback(c *colly.Collector, contentCollector *colly.Collector, links *[]string, wg *sync.WaitGroup) {
	c.OnScraped(func(r *colly.Response) {

		for _, link := range *links {

			wg.Add(1)

			go func(url string) {

				defer wg.Done()
				if err := contentCollector.Visit(url); err != nil {
					log.Printf("Error visiting %s: %v", url, err)
				}
			}(link)
		}
		*links = []string{}
	})
}

func ScrapePage(page int) ([]Article, bool, bool, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.quepensaschacabuco.com"),
	)

	var (
		links       []string
		articles    []Article
		canContinue bool
		canGoBack   bool

		mu sync.Mutex

		wg sync.WaitGroup
	)

	setupCollectors(c, &links, &articles, &mu, &wg, &canContinue, &canGoBack)

	err := c.Visit(fmt.Sprintf("https://www.quepensaschacabuco.com/entradas/%d/", page))
	if err != nil {
		return nil, false, false, err
	}

	c.Wait()

	wg.Wait()

	return articles, canContinue, canGoBack, nil
}
