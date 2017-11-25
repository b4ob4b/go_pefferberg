package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"gopkg.in/xmlpath.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	url := `http://www.restaurant-pfefferberg.de/index.php/speisekarten/tagesgerichte`

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Bitte gib einen Wochentag (Montag - Freitag an.")
		os.Exit(1)
	}

	convert_html_body(url, args[0])

}

func regexpIndex(vs []string, t string) int {

	for i, v := range vs {
		match, _ := regexp.MatchString(t, v)
		if match {
			return i
		}
	}
	return -1
}

func convert_html_body(uri string, weekday string) {
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	get_pefferberg_menu(string(body), weekday)
}

func get_pefferberg_menu(s string, weekday string) {
	stringHtml := string(s)

	reader := strings.NewReader(stringHtml)
	root, err := html.Parse(reader)

	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	html.Render(&b, root)
	fixedHtml := b.String()

	reader = strings.NewReader(fixedHtml)
	xmlroot, xmlerr := xmlpath.ParseHTML(reader)

	if xmlerr != nil {
		log.Fatal(xmlerr)
	}

	xpath := []string{}
	for i := 0; i < 25; i++ {
		xpath = append(xpath,
			`/html/body/div[2]/div[2]/div/div/div[2]/p[`+strconv.Itoa(i+23)+`]`)
	}
	value_list := []string{}
	for i := 0; i < len(xpath); i++ {
		path := xmlpath.MustCompile(xpath[i])
		if value, ok := path.String(xmlroot); ok {
			// fmt.Println(value)
			value_list = append(value_list, value)
		}
	}
	position := regexpIndex(value_list, weekday)
	fmt.Println(value_list[position:(position + 5)])

}
