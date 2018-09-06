package main

import (
	"fmt"
	"sync"
)

//Jokes is a constant holding the Jokes websit URL
const Jokes = "http://www.akposjokes.com"

var waitGroup sync.WaitGroup

func main() {
	categories := ExtractCategories("http://www.akposjokes.com/category/funny-sayings/jokes?page=1")
	waitGroup.Add(len(categories))
	for categoName, categoURL := range categories {
		go scraper(categoName, categoURL)
	}
	//ExtractJoke("http://www.akposjokes.com/joke/frozen-windows")
	waitGroup.Wait()
}

func scraper(categoName, categoURL string) {
	defer waitGroup.Done()
	i := 0
	for {
		url := Jokes + categoURL + "?page=" + string(i)
		titles := ExtractTitle(url)
		for titleName, titleURL := range titles {
			url := Jokes + titleURL
			joke, author, info := ExtractJoke(url)
			PushJokesToDB(titleName, joke, author, categoName, info)
		}
		if len(titles) == 0 {
			fmt.Println(i, "Tour of category ", categoName, "is completed")
		}
		i++
	}
}
