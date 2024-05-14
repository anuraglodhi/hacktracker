package main

import (
	"log"
	"time"
)

type Contest struct {
	Id         string    `json:"id"`
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	ContestUrl string    `json:"url"`
}

func GetContestsFromPlatforms(ps ...Platform) []Contest {
	contests := make([]Contest, 0)

	for _, p := range ps {
		c, err := p.GetContests()
		if err != nil {
			log.Printf("error in getting contests from %s: %s", p.GetName(), err.Error())
		}
		contests = append(contests, c...)
	}

	return contests
}
