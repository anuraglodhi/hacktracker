package main

type Platform interface {
	GetContests() ([]Contest, error)
	GetName() string
}
