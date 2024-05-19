package main

import "log"

type ContestManager struct {
	platforms map[string]Platform
}

func NewContestManager() ContestManager {
	m := ContestManager{}
	m.platforms = make(map[string]Platform)

	m.platforms["codeforces"] = CodeforcesPlatform{}
	m.platforms["leetcode"] = LeetcodePlatform{}

	return m
}

func (m *ContestManager) GetAllContests() map[string][]Contest {
	contests := make(map[string][]Contest)

	for key, p := range m.platforms {
		c, err := p.GetContests()
		if err != nil {
			log.Printf("error in getting contests from %s: %s", p.GetName(), err.Error())
			continue
		}
		contests[key] = c
	}

	return contests
}

func (m *ContestManager) GetContestsOnPlatforms(ps []string) map[string][]Contest {
	contests := make(map[string][]Contest)

	for _, p := range ps {
		plat, ok := m.platforms[p]
		if !ok {
			log.Println("invalid platform requested:", p)
			continue
		}

		c, err := plat.GetContests()
		if err != nil {
			log.Printf("error in getting contests from %s: %s", m.platforms[p].GetName(), err.Error())
			continue
		}
		if len(c) <= 0 {
			log.Printf("no contests found on '%s'", m.platforms[p].GetName())
			continue
		}
		contests[p] = c
	}

	return contests
}
