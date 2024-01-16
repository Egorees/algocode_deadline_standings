package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	GroupShort string `json:"group_short"`
}

type Problem struct {
	Id    string `json:"id"`
	Long  string `json:"long"`
	Short string `json:"short"`
	Index int    `json:"index"`
}

type UserSubmit struct {
	Score   int    `json:"score"`
	Penalty int    `json:"penalty"`
	Verdict string `json:"verdict"`
	Time    int    `json:"time"`
}

type Contest struct {
	Id          int                     `json:"id"`
	Date        string                  `json:"date"`
	EjudgeId    int                     `json:"ejudge_id"`
	Title       string                  `json:"title"`
	Coefficient float64                 `json:"coefficient"`
	Problems    []Problem               `json:"problems"`
	Users       map[string][]UserSubmit `json:"users"`
}

type SubmitsData struct {
	Users    []User    `json:"users"`
	Contests []Contest `json:"contests"`
}

func getSubmitsData(url string) SubmitsData {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error while quering algocode: %v\n", err.Error())
	}
	parser := json.NewDecoder(res.Body)
	var data SubmitsData
	if err = parser.Decode(&data); err != nil {
		log.Fatalf("Error while parsing json from algocode: %v\n", err.Error())
	}
	return data
}
