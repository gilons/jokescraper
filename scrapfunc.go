package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var RegexPatternInfo *regexp.Regexp
var RegexPatternCategories *regexp.Regexp
var RegexPatternTiltle *regexp.Regexp
var RegexPatternJoke *regexp.Regexp
var RegexPatternAuthor *regexp.Regexp
var RegexPatternJoke2 *regexp.Regexp

func init() {
	regexString := `(?m)<span property="dc:date dc:created" content=".+" ` +
		`datatype="xsd:dateTime" rel="sioc:has_creator">(?P<creator>.+)<span class="username"` +
		` xml:lang="".+typeof="sioc:UserAccount" property="foaf:name"` +
		` datatype="">(?P<submitor>.+)<\/span>(?P<date>.+)<\/span>`
	RegexPatternInfo = regexp.MustCompile(regexString)

	regexString = `(?m)<li class="views.+leaf">  
  <div>        <span><a href="(?P<categoryURL>.+)">(?P<categoryName>.+)</a></span>  </div></li>`
	RegexPatternCategories = regexp.MustCompile(regexString)
	RegexPatternTiltle = regexp.MustCompile(`(?m)<div>        <h2 class="title"><a href="(.+)">(.+)</a></h2>  </div>`)
	RegexPatternJoke = regexp.MustCompile(`(?s)<div class="field field-name-body field-type-text-with-summary field-label-hidden"><div class="field-items"><div class="field-item even" property="content:encoded">(.+)</div></div></div><div class="field field-name-field-author field-type-text field-label-inline clearfix"><div class="field-label">Author:&nbsp;</div><div class="field-items">`)
	RegexPatternJoke2 = regexp.MustCompile(`(?s)<div class="field field-name-body field-type-text-with-summary field-label-hidden"><div class="field-items"><div class="field-item even" property="content:encoded">(.+)</div></div></div><ul class="flippy">`)
	RegexPatternAuthor = regexp.MustCompile(`<div class="field field-name-field-author field-type-text field-label-inline clearfix"><div class="field-label">Author:&nbsp;</div><div class="field-items"><div class="field-item even">(.+)</div></div></div>`)

}

//FormatHTML func is a function that puts the html content is a form sweetable for pattern matching
func FormatHTML(html []byte) string {
	htmlArray := strings.Split(string(html), "\n")

	spaceTrimedHTMLArray := make([]string, len(htmlArray))
	for _, element := range htmlArray {
		spaceTrimedHTMLArray = append(spaceTrimedHTMLArray, strings.TrimSpace(element))
	}
	formatedHTML := strings.Join(spaceTrimedHTMLArray, " ")
	return strings.TrimSpace(formatedHTML)
}

//GetURLContent function is used to download the the html content of a url
func GetURLContent(url string) ([]byte, error) {
	client := http.Client{}
	client.Timeout = 0
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	html, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return html, nil
}

//ExtractSubmitionInfo returns ownership information about the joke.
func ExtractSubmitionInfo(html []byte) string {
	names := RegexPatternInfo.SubexpNames()
	temp := make([]string, len(names))
	//fmt.Println(names)
	matches := RegexPatternInfo.FindAllStringSubmatch(string(html), -1)
	//fmt.Println(matches)
	submatchArray := make(map[string]string, len(names))
	for i, name := range matches[0] {
		submatchArray[names[i]] = name
		temp[i] = name
	}
	return strings.Replace(strings.Join(temp[1:], ""), "</span>", "", 1)
}

//ExtractCategories extract all the categories from the start point url
func ExtractCategories(url string) map[string]string {
	html, err := GetURLContent(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	categories := RegexPatternCategories.FindAllStringSubmatch(string(html), -1)
	categoriesArray := make(map[string]string, len(categories))
	for _, category := range categories {
		categoriesArray[category[2]] = category[1]
	}
	return categoriesArray
}

//ExtractTitle Extracts out the titles of a given category in a given pagination
func ExtractTitle(url string) map[string]string {
	html, err := GetURLContent(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	Titles := RegexPatternTiltle.FindAllStringSubmatch(string(html), -1)
	titleArray := make(map[string]string, len(Titles))
	for _, Title := range Titles {
		titleArray[Title[2]] = Title[1]
	}
	return titleArray
}

//ExtractJoke Extracts the given Joke at a particular url
func ExtractJoke(url string) (string, string, string) {
	htlm, err := GetURLContent(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(url)
	Joke := RegexPatternJoke.FindAllStringSubmatch(string(htlm), -1)
	if len(Joke) == 0 {
		fmt.Println("rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr")
		Joke = RegexPatternJoke2.FindAllStringSubmatch(string(htlm), -1)
	}
	fmt.Println(Joke)
	author := ExtractAuthor(string(htlm))
	fmt.Println(author)
	info := ExtractSubmitionInfo(htlm)
	fmt.Println(info)
	return strings.Replace(Joke[0][1], "\"", "\\\"", -1), author, info
}

//ExtractAuthor function is used to extract the autho of a joke
func ExtractAuthor(html string) string {
	author := RegexPatternAuthor.FindAllStringSubmatch(html, -1)
	if len(author) == 0 {
		return "none"
	}
	return strings.Replace(author[0][1], "</div>", "", -1)
}
