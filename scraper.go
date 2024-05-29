package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"sync"
	"time"
)

type Product struct {
	Title string
	Price string
}

func getNumberOfPages(url string) int {
	c := colly.NewCollector()

	var maxPage int

	c.OnHTML("ul.page-numbers", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, element *colly.HTMLElement) {

			pageNumber, err := strconv.Atoi(element.Text) /// converted string to the integer to find max number

			if err == nil && pageNumber > maxPage {
				maxPage = pageNumber
			}

		})
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	return maxPage
}

func getPageData(page int, wg *sync.WaitGroup, productsChan chan<- Product) {

	defer wg.Done()

	url := "https://scrapeme.live/shop/page/" + strconv.Itoa(page) + "/"

	c := colly.NewCollector()

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		title := e.ChildText("h2.woocommerce-loop-product__title")
		price := e.ChildText("span.price")
		productsChan <- Product{
			title,
			price,
		}
	})

	err := c.Visit(url)

	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	startTime := time.Now() // Record start time

	var wg sync.WaitGroup

	url := "https://scrapeme.live/shop/"

	lastPageVal := getNumberOfPages(url)

	productsChan := make(chan Product)

	for page := 1; page <= lastPageVal; page++ {
		wg.Add(1)
		go getPageData(page, &wg, productsChan)
	}

	go func() {
		wg.Wait()
		defer close(productsChan)
	}()

	var products []Product
	for product := range productsChan {
		products = append(products, product)
	}

	var counter int
	for _, product := range products {
		counter++
		log.Printf("Title: %s, Price: %s", product.Title, product.Price)
		fmt.Println(counter)
	}
	fmt.Println("counter", counter)

	fmt.Println("Total number of pages:", lastPageVal)
	fmt.Println("Total number of results", len(products))

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Printf("Time taken: %s", duration)

}
