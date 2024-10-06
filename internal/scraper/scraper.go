package scraper

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/charmbracelet/log"
	md "github.com/JohannesKaufmann/html-to-markdown"
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
	CategoryId string `json:"category_id"`
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
	c.OnHTML(" .post-block__link" + additionalClass, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" {
			*links = append(*links, link)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		if len(*links) == 0 {
			log.Error("Error: No elements found with the specified selector.")
		}
	})


	c.OnHTML("a:contains('Más noticias')", func(e *colly.HTMLElement) {
		*canContinue = true
	})

	c.OnHTML("a:contains('Volver a la página anterior')", func(e *colly.HTMLElement) {
		*canGoBack = true
	})
}

func setupContentCollector(c *colly.Collector, articles *[]Article, mu *sync.Mutex) *colly.Collector {
	contentCollector := c.Clone()

	contentCollector.OnHTML("article.post", func(e *colly.HTMLElement) {
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
	converter := md.NewConverter("", true, nil)

	categoryId := e.DOM.AttrOr("class", "")
	// find the string that contains "category-" in the categoryId string
	categoryId = strings.Split(categoryId, "category-")[1]
	categoryName := strings.ToUpper(categoryId[:1]) + categoryId[1:]

	articleType := e.ChildText("header .post-block__volanta")

	title := e.ChildText("header h1")

	subtitle := e.ChildText("header .article__excerpt")

	e.DOM.Find(".container>div>aside.border-t>a").Remove().End().Text()
	date := e.ChildText(".container>div>aside.border-t")
	date = strings.ReplaceAll(date, "|", "")
	date = strings.TrimSpace(date)

	// remove the elements that match this classes:
	ignoredElements := []string{
		".author",
		"aside",
		".subscribe-to-whatsapp",
		"figure",
	}

	for _, element := range ignoredElements {
		e.DOM.Find(".article__content "+element).Remove().End()
	}

	articleBody := e.ChildText(".article__content")

	// create the following markdown structure:
	/*
		Title
		Subtitle
		Date
		Category
		Body
		Link
	*/

	finalBody := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", title, subtitle, date, categoryName, articleBody, e.Request.URL.String())

	markdown, err := converter.ConvertString(finalBody)
	if err != nil {
		log.Error("Error converting html to markdown: %v", err)
		return nil
	}


	log.Infof("title: %s", title)
	log.Infof("subtitle: %s, articleType: %s", subtitle, articleType)
	log.Infof("date: %s", date)
	log.Infof("categoryName: %s", categoryName)
	log.Infof("articleBody: %s", articleBody)
	log.Infof("link: %s", e.Request.URL.String())

	return &Article{
		Title:      title,
		Date:       date,
		Category:   categoryName,
		CategoryId: categoryId,
		Body:       markdown,
		Link:       e.Request.URL.String(),
	}





	/*

	markdown, err := converter.ConvertString(articleBody)

	if err != nil {
		log.Error("Error converting html to markdown: %v", err)
		return nil
	}

	t, err := parseSpanishDate(info)
	if err != nil {
		log.Error("Error parsing date: %v", err)
		return nil
	}

	return &Article{
		Title:      title,
		Date:       t.Format("2006-01-02 15:04:05"),
		Category:   category,
		CategoryId: categoryId,
		Body:       markdown,
		Link:       e.Request.URL.String(),
	}
	*/
}

func setupOnScrapedCallback(c *colly.Collector, contentCollector *colly.Collector, links *[]string, wg *sync.WaitGroup) {
	c.OnScraped(func(r *colly.Response) {

		for _, link := range *links {

			wg.Add(1)

			go func(url string) {

				defer wg.Done()
				if err := contentCollector.Visit(url); err != nil {
					log.Error("Error visiting %s: %v", url, err)
				}
			}(link)
		}
		*links = []string{}
	})
}

func ScrapePage(page int) ([]Article, bool, bool, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("chacabucoenred.com"),
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

	url := ""
	if (page < 2) {
		url = "https://chacabucoenred.com/"
	} else {
		url = fmt.Sprintf("https://chacabucoenred.com/page/%d/", page)
	}
	err := c.Visit(url)
	if err != nil {
		return nil, false, false, err
	}

	c.Wait()

	wg.Wait()

	return articles, canContinue, canGoBack, nil
}
