package main

import (
	"log"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type CodechefPlatform struct {
}

func (l CodechefPlatform) GetName() string {
	return "codechef"
}

func (l CodechefPlatform) GetContests() ([]Contest, error) {
	pw, err := playwright.Run()
	if err != nil {
		log.Println("failed to start playwright")
		return nil, err
	}

	browser, err := pw.Firefox.Launch()
	if err != nil {
		log.Println("failed to start browser")
		return nil, err
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Println("failed to create a new page")
		return nil, err
	}

	_, err = page.Goto("https://codechef.com/contests")
	if err != nil {
		log.Println("failed to go to contest page")
		return nil, err
	}

	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{})
	// #root > div > div._pageContainer_1se0b_3._dark_1se0b_9 > div > div > div._contest-tables__container_1idej_225 > div:nth-child(1) > div._dataTable__container_1idej_417 > div > div.jss7 > table > tbody
	upcomingContestsSelector := `#root > div > div._pageContainer_1se0b_3._dark_1se0b_9 > div > div > div._contest-tables__container_1idej_225 > div:nth-child(1) > div._dataTable__container_1idej_417 > div > div.jss7 > table > tbody > tr`

	upcomingContestsBody, err := page.Locator(upcomingContestsSelector).All()
	if err != nil || len(upcomingContestsBody) <= 0 {
		log.Println("upcoming contests not found")
		return nil, err
	}

	contests := make([]Contest, 0)
	for _, contestListing := range upcomingContestsBody {
		err := contestListing.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateAttached})
		if err != nil {
			continue
		}

		fields, err := contestListing.Locator("td").All()
		if err != nil  || len(fields) <= 2 {
			log.Println("codechef invalid contest listing")
			continue
		}

		idText, err := fields[0].Locator("div:nth-child(2) > p").InnerText()
		titleElem := fields[1].Locator("div:nth-child(2) > a")
		title, err := titleElem.InnerText()
		a, err := titleElem.GetAttribute("href")
		dateTextA, err := fields[2].Locator("div:nth-child(2) > div > div > p:nth-child(1)").InnerText()
		dateTextB, err := fields[2].Locator("div:nth-child(2) > div > div > p:nth-child(2)").InnerText()
		if err != nil {
			log.Println("codechef parsing failed")
			continue
		}

		dateText := dateTextA + " " + dateTextB
		date, err := l.formatDate(dateText)
		if err != nil {
			continue
		}

		id := "codechef-" + strings.ToLower(idText)

		contests = append(contests, Contest{
			Id: id,
			Title: title,
			Date: date,
			ContestUrl: a,
		})
	}

	return contests, nil
}

func (l CodechefPlatform) formatDate(dateText string) (time.Time, error) {
	format := "_2 Jan 2006 Mon 15:04"
	date, err := time.ParseInLocation(format, dateText, time.Local) 
	if err != nil {
		log.Println("codechef date parse error")
		return time.Now(), err
	}
	return date.UTC(), nil
}
