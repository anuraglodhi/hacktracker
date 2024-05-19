package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type LeetcodePlatform struct {
}

func (l LeetcodePlatform) GetName() string {
	return "leetcode"
}

func (l LeetcodePlatform) GetContests() ([]Contest, error) {
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

	_, err = page.Goto("https://leetcode.com/contest")
	if err != nil {
		log.Println("failed to go to contest page")
	}

	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{})
	cardSelector := `#__next > div.flex.min-h-screen.min-w-\[360px\].flex-col.text-label-1.dark\:text-dark-label-1.bg-layer-bg.dark\:bg-dark-layer-bg > div.mx-auto.w-full.grow.p-0.md\:max-w-none.md\:p-0.lg\:max-w-none.bg-layer-bg.dark\:bg-dark-layer-bg > div > div > div.lc-md\:mt-\[284px\].lc-lg\:mt-\[346px\].mt-\[238px\] > div > div > div.relative.w-full.px-4.lc-lg\:overflow-hidden.-my-1.-mx-4 > div > div > div`

	cards, err := page.Locator(cardSelector).All()

	titleAndDateSelector := `div > a > div.flex.items-center.lc-md\:min-h-\[84px\].min-h-\[80px\].px-4 > div > div`
	titleSelector := `div > span`
	urlSelector := `div > a`

	contests := make([]Contest, 0)
	for _, card := range cards {
		err := card.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateAttached})
		if err != nil {
			continue
		}

		titleAndDates, err := card.Locator(titleAndDateSelector).All()
		if err != nil {
			continue
		}

		title, err := titleAndDates[0].Locator(titleSelector).InnerText()
		if err != nil {
			continue
		}

		dateText, err := titleAndDates[1].InnerText()
		if err != nil {
			continue
		}
		date, err := l.formatDate(dateText)
		if err != nil {
			continue
		}

		a, err := card.Locator(urlSelector).GetAttribute("href")
		if err != nil {
			continue
		}
		a = "https://leetcode.com" + a

		contests = append(contests, Contest{
			Id: "leetcode-1",
			Title: title,
			Date: date,
			ContestUrl: a,
		})
	}

	return contests, nil
}

func (l LeetcodePlatform) formatDate(dateText string) (time.Time, error) {
	parts := strings.Split(dateText, " ")

	date, err := DateOnNextDay(parts[0])
	if err != nil {
		return time.Now(), err
	}

	tz := parts[3]
	var p string
	if strings.Contains(tz, "+") {
		p = strings.Split(tz, "+")[1]
	} else if strings.Contains(tz, "-") {
		p = strings.Split(tz, "-")[1]
	}
	pp := strings.Split(p, ":")
	tzhour, err := strconv.ParseInt(pp[0], 10, 32)
	if err != nil {
		return time.Now(), err
	}
	tzmin, err := strconv.ParseInt(pp[1], 10, 32)
	if err != nil {
		return time.Now(), err
	}
	location := time.FixedZone("leetcode", int(tzhour*3600+tzmin*60))

	t, err := time.ParseInLocation("3:04 PM", parts[1]+" "+parts[2], location)
	t = t.UTC()

	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC), nil
}
