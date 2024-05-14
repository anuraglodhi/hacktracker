package main

import (
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

type CodeforcesPlatform struct{}

func (c CodeforcesPlatform) GetName() string {
	return "codeforces"
}

func (c CodeforcesPlatform) GetContests() ([]Contest, error) {
	url := "https://codeforces.com/contests/"

	res, err := soup.Get(url)
	if err != nil {
		return nil, err
	}
	
	doc := soup.HTMLParse(res)
	
	contests := make([]Contest, 0)
	
	// Selector:
	// #pageContent > div.contestList > div.datatable > div:nth-child(6)
	// > table > tbody > tr -> tablerow for a contest
	trs := doc.
		Find("div", "id", "pageContent").
		Find("div", "class", "contestList").
		Find("div", "class", "datatable").
		FindAll("div")[6].
		Find("table").
		Find("tbody").
		FindAll("tr")

	for _, tr := range trs[1:] {
		tds := tr.FindAll("td")

		// Text
		title := strings.Trim(tds[0].Text(), " \n")

		// Date
		dateText := tds[2].
			Find("a").
			Find("span").Text()

		format := "Jan/02/2006 15:04"
		moscow := time.FixedZone("Moscow Time", int((3 * time.Hour).Seconds()))

		date, err := time.ParseInLocation(format, dateText, moscow)
		if err != nil {
			return nil, err
		}

		// Id
		id := tr.Attrs()["data-contestid"]

		contests = append(contests, Contest{
			Id: "codeforces-"+id,
			Title: title,
			Date: date.UTC(),
			ContestUrl: "https://codeforces.com/contest/"+id,
		})
	}

	return contests, nil
}


