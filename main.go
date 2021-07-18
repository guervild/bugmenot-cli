package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
)

const (
	IS_BARRED = "This site has been barred from the bugmenot system."
)

type LoginResult struct {
	Domain string
	Logins []Login
}

type Login struct {
	Username string
	Password string
	Other    string
	Rate     string
}

func newLoginResult(domain string) *LoginResult {
	return &LoginResult{
		Domain: domain,
	}
}

func (l *LoginResult) PrintIntoTable() string {

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	data := [][]string{}

	for _, result := range l.Logins {

		data = append(data, []string{result.Username, result.Password, result.Other, result.Rate})
	}

	table.SetHeader([]string{"Username", "Password", "Other", "Rating"})
	table.SetRowLine(true)
	table.SetRowSeparator("-")

	for _, v := range data {
		table.Append(v)
	}

	table.Render()

	return fmt.Sprintf("Results for %s:\n%s", l.Domain, tableString.String())
}

func (l *LoginResult) PrintIntoJSON() (string, error) {
	data, err := json.Marshal(l)

	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (l *LoginResult) Filter(filter int) {

	var newLogins []Login

	reg := "[0-9]+"
	pattern := regexp.MustCompile(reg)

	for _, v := range l.Logins {

		rate, _ := strconv.Atoi(pattern.FindString(v.Rate))

		if rate >= filter {
			newLogins = append(newLogins, v)
		}
	}

	l.Logins = newLogins
}

func BugMeNotScrape(domain string) (*LoginResult, error) {

	url := fmt.Sprintf("http://bugmenot.com/view/%s", domain)

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	isBarred := doc.FindMatcher(goquery.Single("div #content p"))

	if isBarred.Text() == IS_BARRED {
		return nil, fmt.Errorf(IS_BARRED)
	}

	var logins []Login

	lResult := newLoginResult(domain)

	doc.Find("article").Each(func(i int, s *goquery.Selection) {

		current := s.Find("kbd")

		login := Login{
			Username: current.Eq(0).Text(),
			Password: current.Eq(1).Text(),
		}

		if len(current.Nodes) == 3 {
			login.Other = current.Eq(2).Text()
		}

		login.Rate = s.FindMatcher(goquery.Single("ul li")).Text()

		logins = append(logins, login)
	})

	lResult.Logins = logins

	return lResult, nil
}

func main() {
	var domain string
	var json bool
	var filter int

	flag.StringVar(&domain, "domain", "", "Domain for which we search shared identifiers (Required)")
	flag.BoolVar(&json, "json", false, "Return results in JSON format. Default is false")
	flag.IntVar(&filter, "filter", 0, "Filter the minimum rate success (greater than or equal to the provided value). Default is 0.")

	flag.Usage = func() {
		h := "A command line tools to request logins shared on bugmenot.com\n\n"

		h += "Usage:\n"
		h += "bugmenot-cli -domain <your domain>\n\n"

		h += "Options:\n"
		h += "  -domain <your domain>  Domain for which we search shared identifiers (Required)\n"
		h += "  -json                  Return results in JSON format. Default is false\n"
		h += "  -filter <value>        Filter the minimum rate success (greater than or equal to the provided value). Default is 0.\n\n"

		h += "Examples:\n"
		h += "  bugmenot-cli -domain <your domain> -json\n"
		h += "  bugmenot-cli -domain <your domain> -json -filter 50\n\n"

		fmt.Fprintf(os.Stderr, h)
	}

	flag.Parse()

	if domain == "" {
		log.Println("-domain is required. Here is the help command:")
		flag.Usage()
		return
	}

	lResult, err := BugMeNotScrape(domain)

	if err != nil {
		log.Fatal(err)
	}

	if len(lResult.Logins) < 1 {
		log.Printf("There is no logins for this domain")
		return

	} else {

		if filter > 0 {
			lResult.Filter(filter)
		}

		if json {
			jsonFormat, err := lResult.PrintIntoJSON()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(jsonFormat)

		} else {

			fmt.Println(lResult.PrintIntoTable())
		}
	}
}
